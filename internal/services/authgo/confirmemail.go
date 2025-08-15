package authgo

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/lunyashon/auth/internal/lib/csrf"
	jwtsso "github.com/lunyashon/auth/internal/lib/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AuthData) ConfirmEmail(ctx context.Context, email, token string) error {

	if _, err := s.DB.User.GetParamsByEmail(ctx, email); err != nil {
		return err
	}

	code := csrf.GenerateConfirmToken()
	if err := s.DB.ConfirmToken.CreateConfirmToken(
		ctx,
		code,
		email,
	); err != nil {
		return err
	}

	go func() {
		body, err := json.Marshal(map[string]string{
			"email": email,
			"code":  code,
		})
		if err != nil {
			s.Log.ErrorContext(ctx, "Failed to marshal body", "error", err)
		}
		s.Rabbit.Send.SendToEmailWithRetry(ctx, body, s.QueueConfirmEmail)
	}()

	return nil
}

func (s *AuthData) CheckConfirmToken(
	ctx context.Context,
	code string,
	accessToken string,
) error {
	userClaims, err := jwtsso.ValidateAccessToken(
		accessToken,
		"sso.auth",
		s.Yaml.NameSSOService,
		s.KeysStore.PublicKey,
	)
	if err != nil {
		return err
	}

	userId, err := strconv.Atoi(userClaims.Subject)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid user id")
	}

	if err := s.DB.ConfirmToken.ConfirmEmailAndDeleteToken(ctx, code, userId); err != nil {
		return err
	}

	return nil
}
