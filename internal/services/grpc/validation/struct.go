package validation

import (
	"log/slog"

	protoc "github.com/lunyashon/protoc/auth/gen/go/sso"
)

type ValidInterface interface {
	RegisterToken(data *protoc.TokenRequest) error
	Register(data *protoc.RegisterRequest) error
	Login(data *protoc.LoginRequest) error
}

type Validate struct {
	V   *ValidInterface
	Log *slog.Logger
}
