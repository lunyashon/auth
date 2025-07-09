package tests

import (
	"context"
	"time"

	"github.com/lunyashon/auth/internal/services/model"
	sso "github.com/lunyashon/protobuf/auth/gen/go/sso/v1"
)

type MockDBValidate struct {
	ValidateItemResult             bool
	ValidateItemError              error
	ValidateServiceServicesinvalid []string
	ValidateServiceError           error
	ValidateDuplicateUserResult    []string
	ValidateDuplicateUserError     error
	ValidateCheckUserError         error
	ValidateCheckUserStruct        *model.UserAuth
	CheckingTokenResult            int
	CheckingTokenError             error
}

func (m *MockDBValidate) ValidateToken(ctx context.Context, apikey string) (bool, error) {
	return m.ValidateItemResult, m.ValidateItemError
}

func (m *MockDBValidate) ValidateServices(ctx context.Context, services []string) ([]string, error) {
	return m.ValidateServiceServicesinvalid, m.ValidateServiceError
}

func (m *MockDBValidate) CheckDuplicateUser(ctx context.Context, email, login string) ([]string, error) {
	return m.ValidateDuplicateUserResult, m.ValidateDuplicateUserError
}

func (m *MockDBValidate) CheckUsingToken(ctx context.Context, token string) (int, error) {
	return m.CheckingTokenResult, m.CheckingTokenError
}

type MockDBToken struct {
	CreateTokenError         error
	InsertRefreshTokensError error
	CheckingTokenResult      int
	CheckingTokenError       error
	RevokeTokenError         error
	CreateForgotTokenError   error
	CheckForgotTokenResult   int
	CheckForgotTokenError    error
}

func (m *MockDBToken) CreateToken(ctx context.Context, data *sso.TokenRequest) error {
	return m.CreateTokenError
}

func (m *MockDBToken) InsertRefreshTokens(ctx context.Context, userId int, refreshToken string, createdAt time.Time, expiresAt time.Duration) error {
	return m.InsertRefreshTokensError
}

func (m *MockDBToken) CheckingToken(ctx context.Context, token string) (int, error) {
	return m.CheckingTokenResult, m.CheckingTokenError
}

func (m *MockDBToken) RevokeToken(ctx context.Context, refreshToken string) error {
	return m.RevokeTokenError
}

func (m *MockDBToken) CreateForgotToken(ctx context.Context, token string, userId int) error {
	return m.CreateForgotTokenError
}

func (m *MockDBToken) CheckForgotToken(ctx context.Context, token string) (int, error) {
	return m.CheckForgotTokenResult, m.CheckForgotTokenError
}

type MockDBUser struct {
	CreateUserResult     int64
	CreateUserError      error
	CheckUserResult      *model.UserAuth
	CheckUserError       error
	GetUserByEmailResult int
	GetUserByEmailError  error
	ChangePasswordError  error
}

func (m *MockDBUser) CreateInDB(ctx context.Context, data *sso.RegisterRequest) (int64, error) {
	return m.CreateUserResult, m.CreateUserError
}

func (m *MockDBUser) CheckUser(ctx context.Context, login, pass string) (*model.UserAuth, error) {
	return m.CheckUserResult, m.CheckUserError
}

func (m *MockDBUser) GetUserByEmail(ctx context.Context, email string) (int, error) {
	return m.GetUserByEmailResult, m.GetUserByEmailError
}

func (m *MockDBUser) ChangePassword(ctx context.Context, userId int, password string) error {
	return m.ChangePasswordError
}

type MockDB struct {
	Validate *MockDBValidate
	Token    *MockDBToken
	User     *MockDBUser
}
