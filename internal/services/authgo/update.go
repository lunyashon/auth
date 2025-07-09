package authgo

import (
	"context"
	"fmt"
	"strconv"

	jwtsso "github.com/lunyashon/auth/internal/lib/jwt"
	sso "github.com/lunyashon/protobuf/auth/gen/go/sso/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Realization method UpdateAccessToken (update access token)
// Return access token or error
func (a AuthData) UpdateAccessToken(ctx context.Context, data *sso.AccessTokenRequest) (string, error) {
	claims, err := jwtsso.ValidateRefreshToken(
		data.RefreshToken,
		a.Yaml.NameSSOService,
		a.KeysStore.PublicKey,
	)
	if err != nil {
		if status.Code(err) == codes.Internal {
			errMessage := fmt.Sprintf("failed to validate refresh token: %v", err)
			a.Log.ErrorContext(
				ctx,
				"ERROR validate refresh token",
				"method", "UpdateAccessToken",
				"point", "validate.refresh_token",
				"message", errMessage,
			)
			return "", status.Errorf(codes.Internal, "failed to validate refresh token")
		}
		return "", err
	}

	if err := a.DB.Token.CheckActiveToken(ctx, data.RefreshToken); err != nil {
		return "", err
	}

	claimsId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		errMessage := fmt.Sprintf("failed to convert subject to int: %v", err)
		a.Log.ErrorContext(
			ctx,
			"failed to convert subject to int",
			"method", "UpdateAccessToken",
			"point", "convert.subject.to.int",
			"message", errMessage,
		)
		return "", status.Errorf(codes.Internal, "failed to convert subject to int")
	}

	accessToken, err := jwtsso.GenerateAccessToken(
		int64(claimsId),
		claims.Audience,
		a.Yaml.NameSSOService,
		a.Yaml.AccessTokenTTL,
		a.KeysStore.PrivateKey,
	)
	if err != nil {
		errMessage := fmt.Sprintf("failed to generate access token: %v", err)
		a.Log.ErrorContext(
			ctx,
			"failed to generate access token",
			"method", "UpdateAccessToken",
			"point", "generate.access_token",
			"message", errMessage,
		)
		return "", status.Errorf(codes.Internal, "failed to generate access token")
	}

	return accessToken, nil
}
