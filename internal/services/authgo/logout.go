package authgo

import (
	"context"
	"fmt"
	"strconv"

	jwtsso "github.com/lunyashon/auth/internal/lib/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Realization method LogoutUser (logout user)
// Return error or nil
func (a *AuthData) LogoutUser(
	ctx context.Context,
	accessToken string,
	refreshToken string,
) error {

	var (
		method = "LogoutUser"
		point  string
	)

	userA, err := jwtsso.ValidateAccessToken(
		accessToken,
		"sso.auth",
		a.Yaml.NameSSOService,
		a.KeysStore.PublicKey,
	)
	if err != nil {

		if st, ok := status.FromError(err); ok {

			point = "validate.access_token"
			a.Log.ErrorContext(
				ctx,
				"failed to validate access_token",
				"method", method,
				"point", point,
				"message", err.Error(),
			)

			if st.Code() == codes.Internal {
				return status.Error(codes.Internal, "failed to validate access_token")
			}
			return err
		}
	}

	userR, err := jwtsso.ValidateRefreshToken(
		refreshToken,
		a.Yaml.NameSSOService,
		a.KeysStore.PublicKey,
	)
	if err != nil {
		if st, ok := status.FromError(err); ok {

			point = "validate.refresh_token"
			a.Log.ErrorContext(
				ctx,
				"failed to validate refresh_token",
				"method", method,
				"point", point,
				"message", err.Error(),
			)

			if st.Code() == codes.Internal {
				return status.Error(codes.Internal, "failed to validate refresh_token")
			}
			return err
		}
	}

	userRId, err := strconv.Atoi(userR.Subject)
	if err != nil {
		errMessage := fmt.Sprintf("failed to convert subject to int: %v", err)
		a.Log.ErrorContext(
			ctx,
			"failed to convert subject to int",
			"method", method,
			"point", point,
			"message", errMessage,
		)
		return status.Errorf(codes.Internal, "failed to convert subject to int")
	}

	if userRId != int(userA) {
		return status.Error(codes.InvalidArgument, "sessions from the token do not match")
	}

	if err := a.DB.Token.RevokeToken(ctx, refreshToken); err != nil {
		return err
	}

	return nil
}
