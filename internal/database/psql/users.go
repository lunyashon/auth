package database

import (
	"context"
	"database/sql"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TypesParams struct {
	Id        int
	Email     string
	Login     string
	Password  string
	CreatedAt time.Time
	Confirmed bool
}

func (s *DatabaseProvider) GetParamsByEmail(
	ctx context.Context,
	email string,
) (*TypesParams, error) {
	var (
		params     TypesParams
		q          string
		methodName = "GetTypesParamsByEmail"
		cancel     context.CancelFunc
	)

	if _, ok := ctx.Deadline(); !ok {
		ctx, cancel = context.WithTimeout(ctx, defaultTokenCreationTimeout)
		defer cancel()
	}

	q = `
		SELECT 
			id, email, login, password, created_at, confirmed 
		FROM 
			users 
		WHERE 
			email = $1`

	if err := s.db.QueryRowContext(ctx, q, email).Scan(
		&params.Id,
		&params.Email,
		&params.Login,
		&params.Password,
		&params.CreatedAt,
		&params.Confirmed,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "email not found")
		}
		errMessage := status.Errorf(codes.Internal, "failed to scan data: %v", err)
		s.log.ErrorContext(
			ctx,
			"ERROR database",
			"method", methodName,
			"point", point,
			"email", email,
			"message", errMessage.Error(),
		)
		return nil, status.Errorf(codes.Internal, "database error")
	}

	return &params, nil
}

func (s *DatabaseProvider) ChangePassword(
	ctx context.Context,
	userId int,
	password string,
) error {
	var (
		methodName = "ChangePassword"
		cancel     context.CancelFunc
	)

	if _, ok := ctx.Deadline(); !ok {
		ctx, cancel = context.WithTimeout(ctx, defaultTokenCreationTimeout)
		defer cancel()
	}

	q := `UPDATE users SET password = $1 WHERE id = $2`

	if _, err := s.db.ExecContext(ctx, q, password, userId); err != nil {
		errMessage := status.Errorf(codes.Internal, "failed to change password: %v", err)
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

func (s *DatabaseProvider) GetParamByUserId(
	ctx context.Context,
	userId int,
) (*TypesParams, error) {
	var (
		params     TypesParams
		q          string
		methodName = "GetParamByUserId"
		cancel     context.CancelFunc
	)

	if _, ok := ctx.Deadline(); !ok {
		ctx, cancel = context.WithTimeout(ctx, defaultTokenCreationTimeout)
		defer cancel()
	}

	q = `
		SELECT 
			id, email, login, password, created_at, confirmed 
		FROM 
			users 
		WHERE 
			id = $1`

	if err := s.db.QueryRowContext(ctx, q, userId).Scan(
		&params.Id,
		&params.Email,
		&params.Login,
		&params.Password,
		&params.CreatedAt,
		&params.Confirmed,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		errMessage := status.Errorf(codes.Internal, "failed to get password: %v", err)
		s.log.ErrorContext(
			ctx,
			"ERROR database",
			"method", methodName,
			"point", point,
			"userId", userId,
			"message", errMessage.Error(),
		)
		return nil, status.Errorf(codes.Internal, "database error")
	}

	return &params, nil
}
