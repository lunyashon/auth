package validate

import (
	"context"
	"fmt"
	"unicode"

	database "github.com/lunyashon/auth/internal/database/psql"

	protoc "github.com/lunyashon/protobuf/auth/gen/go/sso/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	lengthToken = 32
)

var ()

// Verification of the registration token input data
// Return error or nil
func RegisterToken(
	ctx context.Context,
	data *protoc.TokenRequest,
	token string,
	db *database.StructDatabase,
) error {

	if err := validateToken(ctx, token, db.Validator); err != nil {
		return err
	}
	if err := servicesRegisterValidate(ctx, data.Services, db.Validator); err != nil {
		return err
	}
	return nil
}

// Checking syntax token and checking token in database
// Return error or nil
func validateToken(ctx context.Context, token string, db database.ValidateProvider) error {
	if token == "" {
		return status.Error(codes.InvalidArgument, "token is empty")
	}
	if len(token) != lengthToken {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("need to enter %v characters in the token", lengthToken))
	}

	var (
		bigLetters   bool
		smallLetters bool
		numbers      bool
	)

	for _, c := range token {
		if _, exists := specialChar[c]; exists {
			return status.Error(codes.InvalidArgument, "the token must consist only of lowercase, uppercase letters and numbers")
		}
		if _, exists := russianLetters[c]; exists {
			return status.Error(codes.InvalidArgument, "the token must consist only of lowercase, uppercase letters and numbers")
		}
		if unicode.IsDigit(c) {
			numbers = true
		}
		if unicode.IsUpper(c) {
			bigLetters = true
		}
		if unicode.IsLower(c) {
			smallLetters = true
		}
	}

	switch {
	case !numbers:
		return status.Error(codes.InvalidArgument, "the token must contain numbers")
	case !bigLetters:
		return status.Error(codes.InvalidArgument, "the token must contain big letters")
	case !smallLetters:
		return status.Error(codes.InvalidArgument, "the token must contain small letters")
	}

	services, isUsed, err := db.ValidateToken(ctx, token)
	if err != nil {
		if code := status.Code(err); code != codes.NotFound {
			return err
		}
	}
	if isUsed {
		return status.Errorf(codes.AlreadyExists, "token %v already exist", token)
	}
	if services != "" {
		return status.Errorf(codes.AlreadyExists, "token %v already exist", token)
	}
	return nil
}
