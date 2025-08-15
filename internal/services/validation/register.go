package validate

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"
	"unicode"

	database "github.com/lunyashon/auth/internal/database/psql"

	sso "github.com/lunyashon/protobuf/auth/gen/go/sso/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	minLoginLength    = 3
	maxLoginLength    = 20
	minPasswordLength = 7
	maxPasswordLength = 72
)

var (
	russianLetters = map[rune]struct{}{
		'а': {}, 'б': {}, 'в': {}, 'г': {}, 'д': {}, 'е': {}, 'ё': {}, 'ж': {}, 'з': {}, 'и': {},
		'й': {}, 'к': {}, 'л': {}, 'м': {}, 'н': {}, 'о': {}, 'п': {}, 'р': {}, 'с': {}, 'т': {},
		'у': {}, 'ф': {}, 'х': {}, 'ц': {}, 'ч': {}, 'ш': {}, 'щ': {}, 'ъ': {}, 'ы': {}, 'ь': {},
		'э': {}, 'ю': {}, 'я': {},
	}
	specialChar = map[rune]struct{}{
		'~': {}, '@': {}, '\\': {}, '#': {}, '№': {}, '$': {}, ';': {},
		'^': {}, ':': {}, '&': {}, '?': {}, '*': {}, '(': {}, ')': {},
		'-': {}, '+': {}, '=': {}, '`': {}, ',': {}, '.': {}, '\'': {},
	}
	emailRegEx = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

// Validate registration input data
// Returning error or nil
func Register(
	ctx context.Context,
	data *sso.RegisterRequest,
	db *database.StructDatabase,
) ([]int32, error) {

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, defaultTimerWithValidate)
		defer cancel()
	}

	if err := loginRegisterValidate(data.Login); err != nil {
		return nil, err
	}
	if err := passwordRegisterValidate(data.Password); err != nil {
		return nil, err
	}
	if err := emailRegisterValidate(data.Email); err != nil {
		return nil, err
	}
	if err := checkDuplicateUser(ctx, data.Email, data.Login, db.Validator); err != nil {
		return nil, err
	}

	if services, err := apikeyRegisterValidate(ctx, data.Token, db.Validator); err != nil {
		return nil, err
	} else {
		return services, nil
	}
}

// Syntax checking login
// Return error or nil
func loginRegisterValidate(login string) error {

	if login == "" {
		return status.Error(codes.InvalidArgument, "login is empty")
	}

	if len(login) < minLoginLength {
		return status.Errorf(codes.InvalidArgument, "the login is less than %d characters long", minLoginLength)
	}
	if len(login) > maxLoginLength {
		return status.Errorf(codes.InvalidArgument, "the login is more than %d characters long", maxLoginLength)
	}

	for _, val := range login {
		if _, exist := russianLetters[val]; exist {
			return status.Error(codes.InvalidArgument, "the login must have only English letters")
		}
	}
	return nil
}

// Syntax checking password
// Return error or nil
func passwordRegisterValidate(pass string) error {

	var (
		char, bigLetters, num bool
	)

	if pass == "" {
		return status.Error(codes.InvalidArgument, "password is empty")
	}

	if len(pass) < minPasswordLength {
		return status.Errorf(codes.InvalidArgument, "the password is less than %d characters long", minLoginLength)
	}

	if len(pass) > maxPasswordLength {
		return status.Errorf(codes.InvalidArgument, "the password is more than %d characters long", maxPasswordLength)
	}

	for _, val := range pass {
		if _, exist := russianLetters[val]; exist {
			return status.Error(codes.InvalidArgument, "the password must have only English letters")
		}
		if _, exist := specialChar[val]; exist {
			char = true
		}
		if unicode.IsUpper(val) {
			bigLetters = true
		}
		if unicode.IsDigit(val) {
			num = true
		}
	}

	switch {
	case !char:
		return status.Error(codes.InvalidArgument, "the password doesn't matter")
	case !bigLetters:
		return status.Error(codes.InvalidArgument, "the password must contain uppercase letters")
	case !num:
		return status.Error(codes.InvalidArgument, "password requires at least one digit")
	}

	return nil
}

// Syntax checking email
// Return error or nil
func emailRegisterValidate(email string) error {

	if email == "" {
		return status.Error(codes.InvalidArgument, "email is empty")
	}

	if len(email) > 254 {
		return status.Error(codes.InvalidArgument, "email too long")
	}

	if !emailRegEx.MatchString(email) {
		return status.Error(codes.InvalidArgument, "invalid email format")
	}

	return nil
}

// Syntax checking services and checking services in database
// Return error or nil
func servicesRegisterValidate(ctx context.Context, services []int32, db database.ValidateProvider) error {
	if len(services) == 0 {
		return status.Error(codes.InvalidArgument, "services is empty")
	}

	if err := db.ValidateServices(ctx, services); err != nil {
		return err
	}

	return nil
}

// Syntax checking API token and checking in database
// Return error or nil
func apikeyRegisterValidate(ctx context.Context, apikey string, db database.ValidateProvider) ([]int32, error) {
	if apikey == "0" {
		return []int32{}, nil
	}

	services, isUsed, err := db.ValidateToken(ctx, apikey)
	if err != nil {
		return nil, err
	}

	if isUsed {
		return nil, status.Errorf(codes.AlreadyExists, "token %v is used", apikey)
	}

	var servicesList []int32
	json.Unmarshal([]byte(services), &servicesList)

	return servicesList, nil
}

// Checking dublicat user in database
// Return error or nil
func checkDuplicateUser(ctx context.Context, email, login string, db database.ValidateProvider) error {
	if dubl, err := db.CheckDuplicateUser(ctx, email, login); err != nil {
		if dubl != nil {
			return status.Errorf(codes.AlreadyExists, "%v already exist", strings.Join(dubl, " and "))
		} else {
			return status.Error(codes.Internal, "database error")
		}
	}
	return nil
}
