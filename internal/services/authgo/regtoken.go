package authgo

import (
	"context"

	validate "github.com/lunyashon/auth/internal/services/validation"
	protoc "github.com/lunyashon/protobuf/auth/gen/go/sso/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Create token in database
// Return resul (true or false with error)
func (a AuthData) RegisterToken(ctx context.Context, data *protoc.TokenRequest) (bool, error) {

	if err := validate.RegisterToken(ctx, data, a.DB); err != nil {
		return false, err
	}

	if err := a.DB.Token.CreateToken(ctx, data); err != nil {
		a.Log.ErrorContext(ctx, err.Error())
		return false, status.Error(codes.Internal, "database error")
	}

	return true, nil
}
