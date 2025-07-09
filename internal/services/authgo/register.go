package authgo

import (
	"context"

	validate "github.com/lunyashon/auth/internal/services/validation"

	protoc "github.com/lunyashon/protobuf/auth/gen/go/sso/v1"
)

// Validate and register user and permission in database
// Return id user or error
func (a *AuthData) RegisterUser(ctx context.Context, data *protoc.RegisterRequest) (int64, error) {
	if err := validate.Register(ctx, data, a.DB); err != nil {
		return 0, err
	}
	return a.DB.User.CreateInDB(ctx, data)
}
