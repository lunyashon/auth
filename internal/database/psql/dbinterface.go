package database

import (
	"context"
	"time"

	"github.com/lunyashon/auth/internal/services/model"
	sso "github.com/lunyashon/protobuf/auth/gen/go/sso/v1"
)

type ValidateProvider interface {
	ValidateToken(
		ctx context.Context,
		apiKey string,
	) (bool, error)
	ValidateServices(
		ctx context.Context,
		services []string,
	) ([]string, error)
	CheckDuplicateUser(
		ctx context.Context,
		email,
		login string,
	) ([]string, error)
	CheckUsingToken(
		ctx context.Context,
		token string,
	) (int, error)
}

type BaseProvider interface {
	Connect() error
	Close()
	GetData() (string, string)
}

type UserProvider interface {
	CheckUser(
		ctx context.Context,
		login,
		pass string,
	) (*model.UserAuth, error)
	CreateInDB(
		ctx context.Context,
		data *sso.RegisterRequest,
	) (int64, error)
	GetParamsByEmail(
		ctx context.Context,
		email string,
	) (*TypesParams, error)
	ChangePassword(
		ctx context.Context,
		userId int,
		password string,
	) error
	GetParamByUserId(
		ctx context.Context,
		userId int,
	) (*TypesParams, error)
}

type TokenProvider interface {
	CreateToken(
		ctx context.Context,
		data *sso.TokenRequest,
	) error
	InsertRefreshTokens(
		ctx context.Context,
		userId int,
		refreshToken string,
		createdAt time.Time,
		expiresAt time.Duration,
	) error
	CheckActiveToken(
		ctx context.Context,
		refreshToken string,
	) error
	RevokeToken(
		ctx context.Context,
		refreshToken string,
	) error
	CreateForgotToken(
		ctx context.Context,
		token string,
		userId int,
	) error
	CheckForgotToken(
		ctx context.Context,
		token string,
	) (int, error)
	DeleteForgotToken(
		ctx context.Context,
		token string,
	) error
	CreateConfirmToken(
		ctx context.Context,
		code string,
		email string,
	) error
	GetConfirmToken(
		ctx context.Context,
		code string,
	) error
}
