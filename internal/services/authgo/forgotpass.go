package authgo

import (
	"context"
	"encoding/json"

	"github.com/lunyashon/auth/internal/lib/csrf"
)

func (s *AuthData) ForgotPassword(
	ctx context.Context,
	email string,
) error {

	param, err := s.DB.User.GetParamsByEmail(ctx, email)
	if err != nil {
		return err
	}

	token, err := csrf.GenerateResetToken()
	if err != nil {
		return err
	}

	if err := s.DB.Token.CreateForgotToken(ctx, token, param.Id); err != nil {
		return err
	}

	s.Log.InfoContext(
		ctx,
		"Forgot password",
		"email", email,
		"token", token,
	)

	go func() {
		body, err := json.Marshal(map[string]string{
			"email": email,
			"token": token,
		})
		if err != nil {
			s.Log.ErrorContext(ctx, "Failed to marshal body", "error", err)
		}
		s.Rabbit.Send.SendToEmailWithRetry(ctx, body, s.QueueForgotToken)
	}()

	return nil
}

func (s *AuthData) CheckForgotToken(
	ctx context.Context,
	token string,
) error {
	if _, err := s.DB.Token.CheckForgotToken(ctx, token); err != nil {
		return err
	}
	return nil
}
