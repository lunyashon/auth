package authgo

import (
	"context"
	"encoding/json"

	"github.com/lunyashon/auth/internal/lib/csrf"
)

func (s *AuthData) ConfirmEmail(ctx context.Context, email, token string) error {

	if _, err := s.DB.User.GetParamsByEmail(ctx, email); err != nil {
		return err
	}

	code := csrf.GenerateConfirmToken()
	if err := s.DB.Token.CreateConfirmToken(
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
) error {
	return s.DB.Token.GetConfirmToken(ctx, code)
}
