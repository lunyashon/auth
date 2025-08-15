package authgo

import (
	"context"
	"strconv"
	"time"

	jwtsso "github.com/lunyashon/auth/internal/lib/jwt"
	sso "github.com/lunyashon/protobuf/auth/gen/go/sso/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *AuthData) ValidateToken(ctx context.Context, data *sso.ValidateRequest) (int, error) {

	var (
		method = "ValidateToken"
		point  string
	)

	userClaims, err := jwtsso.ValidateAccessToken(
		data.AccessToken,
		data.Service,
		a.Yaml.NameSSOService,
		a.KeysStore.PublicKey,
	)

	if err != nil {
		if _, ok := errors[status.Code(err)]; ok {
			return 0, err
		}
		a.Log.ErrorContext(
			ctx,
			"ERROR validate token",
			"method", method,
			"point", point,
			"message", err.Error(),
		)
		return 0, status.Errorf(codes.Internal, "internal server error")
	}

	userId, err := strconv.Atoi(userClaims.Subject)
	if err != nil {
		return 0, status.Errorf(codes.InvalidArgument, "invalid user id")
	}

	if err := a.Redis.TokenProvider.CheckFromBlackList(
		ctx,
		userId,
		userClaims.IssuedAt.Format(time.RFC3339Nano),
		data.AccessToken,
	); err != nil {
		if status.Code(err) == codes.Unauthenticated {
			return 0, err
		}
		a.Log.ErrorContext(
			ctx,
			"ERROR validate token",
			"method", method,
			"point", point,
			"message", err.Error(),
		)
		return 0, status.Errorf(codes.Internal, "server error")
	}

	return userId, nil
}
