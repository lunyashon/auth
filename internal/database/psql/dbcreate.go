package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	passauth "github.com/lunyashon/auth/internal/lib/passauth"

	_ "github.com/lib/pq"
	sso "github.com/lunyashon/protobuf/auth/gen/go/sso/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	defaultTimeoutCreationUser = 20 * time.Second
)

// Create user and permission in database
// Return ID new user or error
func (s *DatabaseProvider) CreateInDB(
	ctx context.Context,
	data *sso.RegisterRequest,
) (int64, error) {
	var (
		cancel     context.CancelFunc
		point      = "create.user.transcaction"
		methodName = "CreateInDB"
	)

	if _, ok := ctx.Deadline(); !ok {
		ctx, cancel = context.WithTimeout(ctx, defaultTimeoutCreationUser)
		defer cancel()
	}

	// Start transaction
	ts, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		errMessage := status.Errorf(codes.Internal, "failed to start transaction %v", err)
		s.log.ErrorContext(
			ctx,
			"ERROR database",
			"method", methodName,
			"point", point,
			"message", errMessage.Error(),
		)
		return 0, status.Errorf(codes.Internal, "database error")
	}
	defer ts.Rollback()

	pas := passauth.ExecAuthService(&passauth.RealBcrypt{})
	hashPass, err := pas.GeneratePassword([]byte(data.Password))
	if err != nil {
		return 0, err
	}

	var id int64
	const q = `INSERT INTO users (email, login, password)
		  VALUES ($1, $2, $3)
		  RETURNING id`

	passString := string(hashPass)
	err = ts.QueryRowContext(
		ctx, q, data.Email, data.Login, passString,
	).Scan(&id)
	if err != nil {
		errMessage := status.Errorf(codes.Internal, "failed to create user in `users`: %v", err)
		s.log.ErrorContext(
			ctx,
			"ERROR database",
			"method", methodName,
			"point", point,
			"email", data.GetEmail(),
			"login", data.GetLogin(),
			"password", passString,
			"message", errMessage.Error(),
		)
		return 0, status.Errorf(codes.Internal, "database error")
	}

	// Create permission
	for _, v := range data.Services {
		if err := s.insertPermission(ctx, ts, id, v, point); err != nil {
			return 0, err
		}
	}

	if err := s.updateUsingToken(ctx, ts, data.Token); err != nil {
		return 0, err
	}

	if err := ts.Commit(); err != nil {
		errMessage := status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
		s.log.ErrorContext(
			ctx,
			"ERROR database",
			"method", methodName,
			"point", point,
			"message", errMessage.Error(),
		)
		return 0, status.Errorf(codes.Internal, "database error")
	}

	return id, nil
}

// Create permission in database
// Return error
func (s *DatabaseProvider) insertPermission(
	ctx context.Context,
	ts *sql.Tx,
	id int64,
	name, point string,
) error {

	var (
		serviceId  int64
		methodName = "insertPermission"
	)
	const qs = `SELECT id FROM services WHERE name = $1`
	err := ts.QueryRowContext(ctx, qs, name).Scan(&serviceId)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return status.Errorf(codes.NotFound, "{`%v`} service not found", point)
	case err != nil:
		errMessage := status.Errorf(codes.Internal, "failed to get service ID: %v", err)
		s.log.ErrorContext(
			ctx,
			"ERROR database",
			"method", methodName,
			"point", point,
			"service", name,
			"message", errMessage.Error(),
		)
		return status.Errorf(codes.Internal, "database error")
	}

	const qi = `INSERT INTO permission (user_id, service_id) VALUES ($1, $2)`

	if _, err = ts.ExecContext(ctx, qi, id, serviceId); err != nil {
		errMessage := status.Errorf(codes.Internal, "failed to insert permission: %v", err)
		s.log.ErrorContext(
			ctx,
			"ERROR database",
			"method", methodName,
			"point", point,
			"userID", id,
			"serviceID", serviceId,
			"message", errMessage.Error(),
		)
		return status.Errorf(codes.Internal, "database error")
	}
	return nil
}

func (s *DatabaseProvider) updateUsingToken(
	ctx context.Context,
	ts *sql.Tx,
	token string,
) error {
	methodName := "updateUsingToken"
	const q = `UPDATE tokens SET is_used = 1 WHERE token = $1`
	if _, err := ts.ExecContext(ctx, q, token); err != nil {
		errMessage := status.Errorf(codes.Internal, "failed to update table tokens: %v", err)
		s.log.ErrorContext(
			ctx,
			"ERROR database",
			"method", methodName,
			"point", point,
			"token", token,
			"message", errMessage.Error(),
		)
		return status.Errorf(codes.Internal, "database error")
	}
	return nil
}
