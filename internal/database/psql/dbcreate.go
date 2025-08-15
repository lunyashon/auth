package database

import (
	"context"
	"database/sql"
	"slices"
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
	services []int32,
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
			"message", errMessage.Error(),
		)
		return 0, status.Errorf(codes.Internal, "database error")
	}

	servicesStruct, err := s.selectServiceId(ctx, ts, services)
	if err != nil {
		return 0, err
	}

	// Create permission
	for key, v := range servicesStruct {
		if err := s.insertPermission(ctx, ts, id, key, v, point); err != nil {
			return 0, err
		}
	}

	if err := s.updateUsingToken(ctx, ts, data.Token); err != nil {
		return 0, err
	}

	if err := s.createUserProfile(ctx, ts, id); err != nil {
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
	serviceId int32,
	service bool,
	point string,
) error {

	var (
		methodName = "insertPermission"
		expiresAt  time.Time
	)

	if service {
		expiresAt = time.Now().Add(time.Hour * 24 * 30)
	} else {
		expiresAt = time.Now()
	}

	const qi = `INSERT INTO permission (user_id, service_id, active, expires_at) VALUES ($1, $2, $3, $4)`

	if _, err := ts.ExecContext(ctx, qi, id, serviceId, service, expiresAt); err != nil {
		errMessage := status.Errorf(codes.Internal, "failed to insert permission: %v", err)
		s.log.ErrorContext(
			ctx,
			"ERROR database",
			"method", methodName,
			"point", point,
			"userID", id,
			"serviceID", serviceId,
			"active", service,
			"expiresAt", expiresAt,
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

func (s *DatabaseProvider) selectServiceId(
	ctx context.Context,
	ts *sql.Tx,
	services []int32,
) (map[int32]bool, error) {
	var (
		id         int32
		point      = "select.service.id"
		methodName = "selectServiceId"
	)
	const q = `SELECT id FROM services`
	rows, err := ts.QueryContext(ctx, q)
	if err != nil {
		errMessage := status.Errorf(codes.Internal, "failed to select service ID: %v", err)
		s.log.ErrorContext(
			ctx,
			"ERROR database",
			"method", methodName,
			"point", point,
			"message", errMessage.Error(),
		)
		return nil, err
	}
	defer rows.Close()

	var servicesBase = make(map[int32]bool, len(services))

	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			errMessage := status.Errorf(codes.Internal, "failed to scan service ID: %v", err)
			s.log.ErrorContext(
				ctx,
				"ERROR database",
				"method", methodName,
				"point", point,
				"message", errMessage.Error(),
			)
			return nil, err
		}
		if slices.Contains(services, id) {
			servicesBase[id] = true
		} else {
			servicesBase[id] = false
		}
	}

	return servicesBase, nil
}

func (s *DatabaseProvider) createUserProfile(
	ctx context.Context,
	ts *sql.Tx,
	id int64,
) error {
	methodName := "createUserProfile"
	const q = `INSERT INTO users_profile (user_id) VALUES ($1)`
	if _, err := ts.ExecContext(ctx, q, id); err != nil {
		errMessage := status.Errorf(codes.Internal, "failed to create user profile: %v", err)
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
