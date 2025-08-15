package tests

import (
	"context"
	"time"

	database "github.com/lunyashon/auth/internal/database/psql"
	"github.com/lunyashon/auth/internal/services/model"
	sso "github.com/lunyashon/protobuf/auth/gen/go/sso/v1"
)

// TypesParams - копия структуры из database/psql/users.go для тестов
type TypesParams struct {
	Id        int
	Email     string
	Login     string
	Password  string
	CreatedAt time.Time
	Confirmed bool
}

type MockDBValidate struct {
	ValidateItemResult          string
	ValidateItemIsUsed          bool
	ValidateItemError           error
	ValidateServiceError        error
	ValidateDuplicateUserResult []string
	ValidateDuplicateUserError  error
	ValidateCheckUserError      error
	ValidateCheckUserStruct     *model.UserAuth
	CheckingTokenResult         int
	CheckingTokenError          error
}

func (m *MockDBValidate) ValidateToken(ctx context.Context, apikey string) (string, bool, error) {
	return m.ValidateItemResult, m.ValidateItemIsUsed, m.ValidateItemError
}

func (m *MockDBValidate) ValidateServices(ctx context.Context, services []int32) error {
	return m.ValidateServiceError
}

func (m *MockDBValidate) CheckDuplicateUser(ctx context.Context, email, login string) ([]string, error) {
	return m.ValidateDuplicateUserResult, m.ValidateDuplicateUserError
}

func (m *MockDBValidate) CheckUsingToken(ctx context.Context, token string) (int, error) {
	return m.CheckingTokenResult, m.CheckingTokenError
}

type MockDBApiToken struct {
	CreateTokenError error
}

func (m *MockDBApiToken) CreateToken(ctx context.Context, data *sso.TokenRequest) error {
	return m.CreateTokenError
}

type MockDBActiveToken struct {
	InsertRefreshTokensError error
	CheckActiveTokenError    error
	RevokeTokenError         error
}

func (m *MockDBActiveToken) InsertRefreshTokens(ctx context.Context, userId int, refreshToken string, createdAt time.Time, expiresAt time.Duration) error {
	return m.InsertRefreshTokensError
}

func (m *MockDBActiveToken) CheckActiveToken(ctx context.Context, refreshToken string) error {
	return m.CheckActiveTokenError
}

func (m *MockDBActiveToken) RevokeToken(ctx context.Context, refreshToken string) error {
	return m.RevokeTokenError
}

type MockDBForgotToken struct {
	CreateForgotTokenError error
	CheckForgotTokenResult int
	CheckForgotTokenError  error
	DeleteForgotTokenError error
}

func (m *MockDBForgotToken) CreateForgotToken(ctx context.Context, token string, userId int) error {
	return m.CreateForgotTokenError
}

func (m *MockDBForgotToken) CheckForgotToken(ctx context.Context, token string) (int, error) {
	return m.CheckForgotTokenResult, m.CheckForgotTokenError
}

func (m *MockDBForgotToken) DeleteForgotToken(ctx context.Context, token string) error {
	return m.DeleteForgotTokenError
}

type MockDBConfirmToken struct {
	CreateConfirmTokenError error
	GetConfirmTokenResult   int
	GetConfirmTokenError    error
	DeleteConfirmTokenError error
}

func (m *MockDBConfirmToken) CreateConfirmToken(ctx context.Context, code string, email string) error {
	return m.CreateConfirmTokenError
}

func (m *MockDBConfirmToken) GetConfirmToken(ctx context.Context, code string, userId int) (int, error) {
	return m.GetConfirmTokenResult, m.GetConfirmTokenError
}

func (m *MockDBConfirmToken) DeleteConfirmToken(ctx context.Context, code string) error {
	return m.DeleteConfirmTokenError
}

type MockDBUser struct {
	CreateUserResult       int64
	CreateUserError        error
	CheckUserResult        *model.UserAuth
	CheckUserError         error
	GetUserByEmailResult   int
	GetUserByEmailError    error
	ChangePasswordError    error
	GetParamsByEmailResult *TypesParams
	GetParamsByEmailError  error
	GetParamByUserIdResult *TypesParams
	GetParamByUserIdError  error
}

func (m *MockDBUser) CreateInDB(ctx context.Context, data *sso.RegisterRequest, services []string) (int64, error) {
	return m.CreateUserResult, m.CreateUserError
}

func (m *MockDBUser) CheckUser(ctx context.Context, login, pass string) (*model.UserAuth, error) {
	return m.CheckUserResult, m.CheckUserError
}

func (m *MockDBUser) GetUserByEmail(ctx context.Context, email string) (int, error) {
	return m.GetUserByEmailResult, m.GetUserByEmailError
}

func (m *MockDBUser) GetParamsByEmail(ctx context.Context, email string) (*TypesParams, error) {
	return m.GetParamsByEmailResult, m.GetParamsByEmailError
}

func (m *MockDBUser) ChangePassword(ctx context.Context, userId int, password string) error {
	return m.ChangePasswordError
}

func (m *MockDBUser) GetParamByUserId(ctx context.Context, userId int) (*TypesParams, error) {
	return m.GetParamByUserIdResult, m.GetParamByUserIdError
}

type MockDBServices struct {
	GetServicesListResult   []database.ServicesList
	GetServicesListError    error
	GetServicesByNameResult []database.ServicesList
	GetServicesByNameError  error
	GetServiceByIdResult    database.ServicesList
	GetServiceByIdError     error
}

func (m *MockDBServices) GetServicesList(ctx context.Context) ([]database.ServicesList, error) {
	return m.GetServicesListResult, m.GetServicesListError
}

func (m *MockDBServices) GetServicesByName(ctx context.Context, name string) ([]database.ServicesList, error) {
	return m.GetServicesByNameResult, m.GetServicesByNameError
}

func (m *MockDBServices) GetServiceById(ctx context.Context, id int32) (database.ServicesList, error) {
	return m.GetServiceByIdResult, m.GetServiceByIdError
}

type MockDB struct {
	Validate     *MockDBValidate
	Token        *MockDBApiToken
	ActiveToken  *MockDBActiveToken
	ForgotToken  *MockDBForgotToken
	ConfirmToken *MockDBConfirmToken
	Services     *MockDBServices
	User         *MockDBUser
}
