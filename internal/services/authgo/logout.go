package authgo

import (
	"context"
	"strconv"

	jwtsso "github.com/lunyashon/auth/internal/lib/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Realization method LogoutUser (logout user)
// Return error or nil
func (a *AuthData) OnceLogoutUser(
	ctx context.Context,
	accessToken string,
	refreshToken string,
) error {

	var (
		method = "LogoutUser"
		point  string
	)

	userClaims, err := jwtsso.ValidateAccessToken(
		accessToken,
		"sso.auth",
		a.Yaml.NameSSOService,
		a.KeysStore.PublicKey,
	)
	if err != nil {

		if status.Code(err) == codes.Internal {

			point = "validate.access_token"
			a.Log.ErrorContext(
				ctx,
				"failed to validate access_token",
				"method", method,
				"point", point,
				"message", err.Error(),
			)

			return status.Error(codes.Internal, "failed to validate access_token")
		}
		return err
	}

	if userClaims == nil {
		return status.Error(codes.Unauthenticated, "invalid token")
	}

	userId, err := strconv.Atoi(userClaims.Subject)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid user id")
	}

	if err := a.Redis.TokenProvider.AddToBlackList(ctx, userId, accessToken, "once"); err != nil {
		return err
	}

	userR, err := jwtsso.ValidateRefreshToken(
		refreshToken,
		a.Yaml.NameSSOService,
		a.KeysStore.PublicKey,
	)
	if err != nil {

		if status.Code(err) == codes.Internal {

			point = "validate.refresh_token"
			a.Log.ErrorContext(
				ctx,
				"failed to validate refresh_token",
				"method", method,
				"point", point,
				"message", err.Error(),
			)

			return status.Error(codes.Internal, "failed to validate refresh_token")
		}
		return err
	}

	userRId, err := strconv.Atoi(userR.Subject)
	if err != nil {
		return status.Errorf(codes.Internal, "invalid user id")
	}

	if userRId != userId {
		return status.Error(codes.InvalidArgument, "sessions from the token do not match")
	}

	if err := a.DB.ActiveToken.RevokeToken(ctx, refreshToken, userId); err != nil {
		return err
	}

	return nil
}

func (a *AuthData) MassLogoutUser(
	ctx context.Context,
	accessToken string,
	refreshToken string,
) error {

	var (
		method = "MassLogoutUser"
		point  string
	)

	userClaims, err := jwtsso.ValidateAccessToken(
		accessToken,
		"sso.auth",
		a.Yaml.NameSSOService,
		a.KeysStore.PublicKey,
	)
	if err != nil {

		if status.Code(err) == codes.Internal {

			point = "validate.access_token"
			a.Log.ErrorContext(
				ctx,
				"failed to validate access_token",
				"method", method,
				"point", point,
				"message", err.Error(),
			)

			return status.Error(codes.Internal, "failed to validate access_token")
		}
		return err
	}

	if userClaims == nil {
		return status.Error(codes.Unauthenticated, "invalid token")
	}

	userId, err := strconv.Atoi(userClaims.Subject)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid user id")
	}

	if err := a.Redis.TokenProvider.AddToBlackList(ctx, userId, accessToken, "mass"); err != nil {
		return err
	}

	userR, err := jwtsso.ValidateRefreshToken(
		refreshToken,
		a.Yaml.NameSSOService,
		a.KeysStore.PublicKey,
	)
	if err != nil {

		if status.Code(err) == codes.Internal {

			point = "validate.refresh_token"
			a.Log.ErrorContext(
				ctx,
				"failed to validate refresh_token",
				"method", method,
				"point", point,
				"message", err.Error(),
			)

			return status.Error(codes.Internal, "failed to validate refresh_token")
		}
		return err
	}

	userRId, err := strconv.Atoi(userR.Subject)
	if err != nil {
		return status.Errorf(codes.Internal, "invalid user id")
	}

	if userRId != userId {
		return status.Error(codes.InvalidArgument, "sessions from the token do not match")
	}

	if err := a.DB.ActiveToken.RevokeAllTokens(ctx, userId); err != nil {
		return err
	}

	return nil
}
