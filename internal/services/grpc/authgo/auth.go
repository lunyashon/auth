package authgo

import (
	"context"
	"log/slog"

	protoc "github.com/lunyashon/protoc/auth/gen/go/sso"
)

// Интерфейс для реализации методов
type Auth interface {
	LoginUser(
		ctx context.Context,
		data *protoc.LoginRequest,
	) (token string, err error)
	RegisterUser(
		ctx context.Context,
		data *protoc.RegisterRequest,
	) (userId int32, err error)
	RegisterToken(
		ctx context.Context,
		data *protoc.TokenRequest,
	) (result bool, err error)
}

type AuthData struct {
	Auth Auth
	Log  *slog.Logger
}
