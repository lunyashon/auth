package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/lunyashon/auth/internal/lib/hash"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Token struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Token     string    `db:"token"`
	CreatedAt time.Time `db:"created_at"`
	ExpiresAt time.Time `db:"expires_at"`
	IsActive  bool      `db:"is_active"`
}

type ForgotToken struct {
	UserID    int       `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	ExpiresAt time.Time `db:"expires_at"`
}

func (s *DatabaseProvider) InsertRefreshTokens(
	ctx context.Context,
	userId int,
	refreshToken string,
	createdAt time.Time,
	expiresAt time.Duration,
) error {

	var (
		cancel     context.CancelFunc
		result     any
		methodName = "InsertAuthTokens"
	)

	if _, ok := ctx.Deadline(); !ok {
		ctx, cancel = context.WithTimeout(ctx, defaultTokenCreationTimeout)
		defer cancel()
	}

	q := `
		INSERT INTO active_tokens
			(user_id, refresh_token, created_at, expires_at, is_active)
		VALUES
			($1, $2, $3, $4, true)
		RETURNING 
			id`

	err := s.db.QueryRowContext(
		ctx,
		q,
		userId,
		hash.HashToken(refreshToken),
		createdAt,
		createdAt.Add(expiresAt),
	).Scan(&result)

	if err != nil {
		errMessage := status.Errorf(codes.Internal, "failed to insert auth tokens %v", err)
		s.log.ErrorContext(
			ctx,
			"ERROR database",
			"method", methodName,
			"point", point,
			"message", errMessage.Error(),
		)
		return status.Errorf(codes.Internal, "database error")
	}

	return nil
}

func (s *DatabaseProvider) CheckActiveToken(
	ctx context.Context,
	refreshToken string,
) error {
	var (
		methodName = "CheckActiveToken"
		cancel     context.CancelFunc
		token      Token
	)

	if _, ok := ctx.Deadline(); !ok {
		ctx, cancel = context.WithTimeout(ctx, defaultTokenCreationTimeout)
		defer cancel()
	}

	q := `
		SELECT 
			*
		FROM 
			active_tokens
		WHERE 
			refresh_token = $1
			AND is_active = true`

	err := s.db.QueryRowContext(
		ctx,
		q,
		hash.HashToken(refreshToken),
	).Scan(
		&token.ID,
		&token.UserID,
		&token.Token,
		&token.CreatedAt,
		&token.ExpiresAt,
		&token.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return status.Errorf(codes.NotFound, "token not found")
		}
		errMessage := status.Errorf(codes.Internal, "failed to check active token %v", err)
		s.log.ErrorContext(
			ctx,
			"ERROR database",
			"method", methodName,
			"point", point,
			"message", errMessage.Error(),
		)
		return status.Errorf(codes.Internal, "database error")
	}

	if ok := token.ExpiresAt.Before(time.Now()); ok {
		return status.Errorf(codes.Unauthenticated, "token expired")
	}

	return nil
}

func (s *DatabaseProvider) RevokeToken(
	ctx context.Context,
	refreshToken string,
) error {

	var (
		methodName = "RevokeToken"
		cancel     context.CancelFunc
	)

	if _, ok := ctx.Deadline(); !ok {
		ctx, cancel = context.WithTimeout(ctx, defaultTokenCreationTimeout)
		defer cancel()
	}

	q := `
		UPDATE 
			active_tokens
		SET
			is_active = false
		WHERE
			refresh_token = $1`

	result, err := s.db.ExecContext(
		ctx,
		q,
		hash.HashToken(refreshToken),
	)

	if err != nil {
		errMessage := status.Errorf(codes.Internal, "failed to revoke token %v", err)
		s.log.ErrorContext(
			ctx,
			"ERROR database",
			"method", methodName,
			"point", point,
			"message", errMessage.Error(),
		)
		return status.Errorf(codes.Internal, "database error")
	}

	rowsAffected, err := result.RowsAffected()
	errMessage := status.Errorf(codes.Internal, "failed to get rows affected %v", err)
	if err != nil {
		s.log.ErrorContext(
			ctx,
			"ERROR database",
			"method", methodName,
			"point", point,
			"message", errMessage,
		)
		return status.Errorf(codes.Internal, "database error")
	}

	if rowsAffected == 0 {
		return status.Errorf(codes.NotFound, "token not found")
	}

	return nil
}

func (s *DatabaseProvider) CreateForgotToken(
	ctx context.Context,
	token string,
	userId int,
) error {
	var (
		methodName = "CreateForgotToken"
		cancel     context.CancelFunc
	)

	if _, ok := ctx.Deadline(); !ok {
		ctx, cancel = context.WithTimeout(ctx, defaultTokenCreationTimeout)
		defer cancel()
	}

	q := `
		INSERT INTO 
			forgot_tokens (token, user_id)
		VALUES 
			($1, $2)`

	_, err := s.db.ExecContext(ctx, q, token, userId)

	if err != nil {
		errMessage := status.Errorf(codes.Internal, "failed to create forgot token %v", err)
		s.log.ErrorContext(
			ctx,
			"ERROR database",
			"method", methodName,
			"point", point,
			"message", errMessage.Error(),
		)
		return status.Errorf(codes.Internal, "database error")
	}

	return nil
}

func (s *DatabaseProvider) CheckForgotToken(
	ctx context.Context,
	token string,
) (int, error) {
	var (
		methodName  = "CheckForgotToken"
		cancel      context.CancelFunc
		forgotToken ForgotToken
	)

	if _, ok := ctx.Deadline(); !ok {
		ctx, cancel = context.WithTimeout(ctx, defaultTokenCreationTimeout)
		defer cancel()
	}

	q := `
		SELECT 
			user_id, created_at, expires_at
		FROM 
			forgot_tokens
		WHERE 
			token = $1`

	err := s.db.QueryRowContext(
		ctx,
		q,
		token,
	).Scan(&forgotToken.UserID, &forgotToken.CreatedAt, &forgotToken.ExpiresAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, status.Errorf(codes.NotFound, "token not found")
		}
		errMessage := status.Errorf(codes.Internal, "failed to check forgot token %v", err)
		s.log.ErrorContext(
			ctx,
			"ERROR database",
			"method", methodName,
			"point", point,
			"message", errMessage.Error(),
		)
		return 0, status.Errorf(codes.Internal, "database error")
	}

	if ok := forgotToken.ExpiresAt.Before(time.Now()); ok {
		if err := s.DeleteForgotToken(ctx, token); err != nil {
			return 0, err
		}
		return 0, status.Errorf(codes.Unauthenticated, "token expired")
	}

	return forgotToken.UserID, nil
}

func (s *DatabaseProvider) DeleteForgotToken(
	ctx context.Context,
	token string,
) error {
	var (
		methodName = "DeleteForgotToken"
		cancel     context.CancelFunc
	)

	if _, ok := ctx.Deadline(); !ok {
		ctx, cancel = context.WithTimeout(ctx, defaultTokenCreationTimeout)
		defer cancel()
	}

	q := `DELETE FROM forgot_tokens WHERE token = $1`

	_, err := s.db.ExecContext(ctx, q, token)

	if err != nil {
		errMessage := status.Errorf(codes.Internal, "failed to delete forgot token %v", err)
		s.log.ErrorContext(
			ctx,
			"ERROR database",
			"method", methodName,
			"point", point,
			"message", errMessage.Error(),
		)
		return status.Errorf(codes.Internal, "database error")
	}

	return nil
}

func (s *DatabaseProvider) CreateConfirmToken(
	ctx context.Context,
	code string,
	email string,
) error {
	var (
		methodName = "CreateConfirmToken"
		cancel     context.CancelFunc
	)

	if _, ok := ctx.Deadline(); !ok {
		ctx, cancel = context.WithTimeout(ctx, defaultTokenCreationTimeout)
		defer cancel()
	}

	q := `INSERT INTO confirm_email_tokens (code, email) VALUES ($1, $2)`

	_, err := s.db.ExecContext(ctx, q, code, email)

	if err != nil {
		errMessage := status.Errorf(codes.Internal, "failed to create confirm email token %v", err)
		s.log.ErrorContext(
			ctx,
			"ERROR database",
			"method", methodName,
			"point", point,
			"message", errMessage.Error(),
		)
		return status.Errorf(codes.Internal, "database error")
	}

	return nil
}

func (s *DatabaseProvider) GetConfirmToken(
	ctx context.Context,
	code string,
) error {
	var (
		methodName = "GetConfirmToken"
		cancel     context.CancelFunc
		res        int
	)

	if _, ok := ctx.Deadline(); !ok {
		ctx, cancel = context.WithTimeout(ctx, defaultTokenCreationTimeout)
		defer cancel()
	}

	q := `SELECT 1 FROM confirm_email_tokens WHERE code = $1`

	err := s.db.QueryRowContext(ctx, q, code).Scan(&res)

	if err != nil {
		if err == sql.ErrNoRows {
			return status.Errorf(codes.NotFound, "token not found")
		}
		errMessage := status.Errorf(codes.Internal, "failed to check confirm email token %v", err)
		s.log.ErrorContext(
			ctx,
			"ERROR database",
			"method", methodName,
			"point", point,
			"message", errMessage.Error(),
		)
		return status.Errorf(codes.Internal, "database error")
	}

	return nil
}
