package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	_ "github.com/lib/pq"
	sso "github.com/lunyashon/protobuf/auth/gen/go/sso/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	defaultTokenCreationTimeout = 5 * time.Second
)

// Create token in database
// Return error or nil
func (s *DatabaseProvider) CreateToken(ctx context.Context, data *sso.TokenRequest) error {

	var (
		cancel     context.CancelFunc
		point      = "create.token"
		methodName = "CreateToken"
	)

	if _, ok := ctx.Deadline(); !ok {
		ctx, cancel = context.WithTimeout(ctx, defaultTokenCreationTimeout)
		defer cancel()
	}

	jsonServ, err := json.Marshal(data.Services)
	if err != nil {
		errMessage := status.Errorf(codes.Internal, "failed to json initialization %v", err)
		s.log.ErrorContext(
			ctx,
			"ERROR database",
			"method", methodName,
			"point", point,
			"data", data.GetServices(),
			"message", errMessage.Error(),
		)
		return status.Errorf(codes.Internal, "database error")
	}

	const q = `
		INSERT INTO tokens 
			(token, services, is_used) 
		VALUES 
			($1, $2, 0)
		RETURNING 
			token`

	var insertingToken string

	err = s.db.QueryRowContext(ctx, q, data.Token, jsonServ).Scan(&insertingToken)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return status.Error(codes.Canceled, "the token was not created")
	case err != nil:
		errMessage := status.Errorf(codes.Internal, "failed to insert %v", err)
		s.log.ErrorContext(
			ctx,
			"ERROR database",
			"method", methodName,
			"point", point,
			"token", data.GetToken(),
			"json", jsonServ,
			"message", errMessage.Error(),
		)
		return status.Errorf(codes.Internal, "database error")
	}

	return nil
}
