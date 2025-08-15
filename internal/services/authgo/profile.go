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

func (a *AuthData) GetProfile(ctx context.Context, accessToken string) (*sso.ProfileResponse, error) {
	var (
		method = "GetProfile"
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

			if status.Code(err) == codes.Internal {
				return nil, status.Error(codes.Internal, "failed to validate access_token")
			}
		}
		return nil, err
	}

	if userClaims == nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	userId, err := strconv.Atoi(userClaims.Subject)
	if err != nil {
		point = "convert.subject.to.int"
		a.Log.ErrorContext(
			ctx,
			"failed to convert subject to int",
			"method", method,
			"point", point,
			"message", err.Error(),
		)
		return nil, status.Error(codes.Internal, "server error")
	}

	if err := a.Redis.TokenProvider.CheckFromBlackList(
		ctx,
		userId,
		userClaims.IssuedAt.Format(time.RFC3339Nano),
		accessToken,
	); err != nil {
		if status.Code(err) == codes.Unauthenticated {
			return nil, err
		}
		a.Log.ErrorContext(
			ctx,
			"ERROR check token from black list",
			"method", method,
			"point", point,
			"message", err.Error(),
		)
		return nil, status.Errorf(codes.Internal, "server error")
	}

	return a.DB.User.GetProfile(ctx, userId)
}

func (a *AuthData) GetMiniProfile(ctx context.Context, accessToken string) (*sso.MiniProfileResponse, error) {
	var (
		method = "GetMiniProfile"
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
			return nil, status.Error(codes.Internal, "failed to validate access_token")
		}
		return nil, err
	}

	if userClaims == nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	userId, err := strconv.Atoi(userClaims.Subject)
	if err != nil {
		point = "convert.subject.to.int"
		a.Log.ErrorContext(
			ctx,
			"failed to convert subject to int",
			"method", method,
			"point", point,
			"message", err.Error(),
		)
		return nil, status.Error(codes.Internal, "server error")
	}

	if err := a.Redis.TokenProvider.CheckFromBlackList(
		ctx,
		userId,
		userClaims.IssuedAt.Format(time.RFC3339Nano),
		accessToken,
	); err != nil {
		if status.Code(err) == codes.Unauthenticated {
			return nil, err
		}
		a.Log.ErrorContext(
			ctx,
			"ERROR check token from black list",
			"method", method,
			"point", point,
			"message", err.Error(),
		)
		return nil, status.Errorf(codes.Internal, "server error")
	}

	return a.DB.User.GetMiniProfile(ctx, userId)
}
