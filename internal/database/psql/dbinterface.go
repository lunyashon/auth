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
	) (string, bool, error)
	ValidateServices(
		ctx context.Context,
		services []int32,
	) error
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
		services []int32,
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
	GetProfile(
		ctx context.Context,
		userId int,
	) (*sso.ProfileResponse, error)
	GetMiniProfile(
		ctx context.Context,
		userId int,
	) (*sso.MiniProfileResponse, error)
}

type ApiTokenProvider interface {
	CreateToken(
		ctx context.Context,
		data *sso.TokenRequest,
		token string,
	) error
}

type ActiveTokenProvider interface {
	InsertRefreshTokens(
		ctx context.Context,
		userId int,
		refreshToken string,
		createdAt time.Time,
		expiresAt time.Duration,
		ip string,
		device string,
	) error
	CheckActiveToken(
		ctx context.Context,
		refreshToken string,
	) error
	RevokeToken(
		ctx context.Context,
		refreshToken string,
		userId int,
	) error
	RevokeAllTokens(
		ctx context.Context,
		userId int,
	) error
}

type ForgotTokenProvider interface {
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
}

type ConfirmTokenProvider interface {
	CreateConfirmToken(
		ctx context.Context,
		code string,
		email string,
	) error
	ConfirmEmailAndDeleteToken(
		ctx context.Context,
		code string,
		userId int,
	) error
}

type ServicesProvider interface {
	GetServicesList(
		ctx context.Context,
	) ([]*sso.StructureServices, error)
	GetServicesByName(
		ctx context.Context,
		name string,
	) ([]*sso.StructureServices, error)
	GetServiceById(
		ctx context.Context,
		id int32,
	) (*sso.StructureServices, error)
}
