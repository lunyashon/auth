package authgo

import (
	"context"
	"time"

	jwtsso "github.com/lunyashon/auth/internal/lib/jwt"
	validate "github.com/lunyashon/auth/internal/services/validation"
	protoc "github.com/lunyashon/protobuf/auth/gen/go/sso/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Sign in user and get JWT token
// Return model UserAuth (JWT, Services, Id and time to life token)
func (a AuthData) LoginUser(
	ctx context.Context,
	data *protoc.LoginRequest,
	device, ip string,
) (*jwtsso.TokenPair, error) {

	var (
		method = "LoginUser"
	)

	if err := validate.Auth(ctx, data, a.DB); err != nil {
		return nil, err
	}

	user, err := a.DB.User.CheckUser(ctx, data.Login, data.Password)
	if err != nil {
		return nil, err
	}

	tokens, err := jwtsso.GenerateToken(
		user.UID,
		user.Services,
		a.Yaml.NameSSOService,
		a.Yaml.AccessTokenTTL,
		a.Yaml.RefreshTokenTTL,
		a.KeysStore.PrivateKey,
		ip,
		device,
	)
	if err != nil {
		point := "token"
		a.Log.ErrorContext(
			ctx,
			"failed to generate token",
			"method", method,
			"point", point,
			"message", err.Error(),
		)
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	if err := a.DB.ActiveToken.InsertRefreshTokens(
		ctx,
		int(user.UID),
		tokens.RefreshToken,
		time.Now(),
		a.Yaml.RefreshTokenTTL,
		ip,
		device,
	); err != nil {
		return nil, err
	}

	return tokens, nil
}
