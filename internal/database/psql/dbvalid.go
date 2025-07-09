package database

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	duplicateResult struct {
		Type  string `db:"type"`
		Value string `db:"value"`
	}

	userAuthData struct {
		Pass    string `db:"password"`
		Service string `db:"name"`
		Uid     int64  `db:"id"`
	}
)

const (
	point = "database"
)

// Checking the entered registration data for a duplicate
// Return empty line on success or (empty line or error on error in database)
// or the found entity in case of an error
func (s *DatabaseProvider) CheckDuplicateUser(ctx context.Context, email, login string) ([]string, error) {

	var (
		duplicates []duplicateResult
		result     = make([]string, 0)
		methodName = "CheckDuplicateUser"
	)

	q := `SELECT 'email' AS type, email AS value FROM users WHERE email = $1
		 UNION ALL
		 SELECT 'login' AS type, login AS value FROM users WHERE login = $2`

	if err := sqlx.SelectContext(ctx, s.db, &duplicates, q, email, login); err != nil {
		errMessage := status.Errorf(codes.Internal, "failed to query: %v", err)
		s.log.ErrorContext(
			ctx,
			"ERROR database",
			"method", methodName,
			"point", point,
			"email", email,
			"login", login,
			"message", errMessage.Error(),
		)
		return nil, status.Error(codes.Internal, "database error")
	}

	for _, val := range duplicates {
		result = append(result, val.Type)
	}

	if len(result) == 0 {
		return nil, nil
	} else {
		return result, status.Errorf(codes.AlreadyExists, "{`%v`} duplicates user:", point)
	}
}

// Checking the access token for availability
// Return count found tokens in Postgree or error
func (s *DatabaseProvider) ValidateToken(ctx context.Context, item string) (bool, error) {

	var (
		exists     bool
		q          string
		methodName = "ValidateToken"
	)

	q = `SELECT EXISTS
		 (SELECT 1 FROM tokens WHERE token = $1)`

	if err := s.db.QueryRowContext(ctx, q, item).Scan(&exists); err != nil {
		errMessage := status.Errorf(codes.Internal, "failed to scan data: %v", err)
		s.log.ErrorContext(
			ctx,
			"ERROR database",
			"method", methodName,
			"point", point,
			"token", item,
			"message", errMessage.Error(),
		)
		return false, status.Error(codes.Internal, "database error")
	}

	return exists, nil
}

// Checking services in database
// Return empty slice on success or (error and undiscovered services for error)
func (s *DatabaseProvider) ValidateServices(ctx context.Context, item []string) ([]string, error) {

	var (
		dbServices    = make(map[string]struct{})
		validServices []string
		res           []string
		methodName    = "ValidateServices"
	)

	services := pq.Array(item)

	if len(item) > 0 {
		q := `
			SELECT name 
			FROM services
			WHERE name = any($1)`

		if err := sqlx.SelectContext(ctx, s.db, &validServices, q, services); err != nil {
			errMessage := status.Errorf(codes.Internal, "failed to query %v", err)
			s.log.ErrorContext(
				ctx,
				"ERROR database",
				"method", methodName,
				"point", point,
				"services", services,
				"message", errMessage.Error(),
			)
			return res, status.Error(codes.Internal, "database error")
		}

		for _, val := range validServices {
			dbServices[val] = struct{}{}
		}
	}

	if res := findMissingServices(dbServices, item); len(res) > 0 {
		return res, status.Errorf(codes.InvalidArgument, "{`%v`} invalid services:", point)
	}

	return res, nil
}

// Comparison of 2 cross-sections of services to identify erroneous ones
// Return slice services
func findMissingServices(dbServices map[string]struct{}, inputServices []string) []string {
	missing := make([]string, 0)
	for _, val := range inputServices {
		if _, exist := dbServices[val]; !exist {
			missing = append(missing, val)
		}
	}

	return missing
}

// Verification of the token for use
// Return numbers of uses or error
func (s *DatabaseProvider) CheckUsingToken(ctx context.Context, token string) (int, error) {
	var (
		exists     = make([]int, 0)
		q          string
		methodName = "CheckUsingToken"
	)

	q = `SELECT is_used FROM tokens WHERE token = $1`

	if err := sqlx.SelectContext(ctx, s.db, &exists, q, token); err != nil {
		errMessage := status.Errorf(codes.Internal, "failed to query token: %v", err)
		s.log.ErrorContext(
			ctx,
			"ERROR database",
			"method", methodName,
			"point", point,
			"token", token,
			"message", errMessage.Error(),
		)
		return -1, status.Error(codes.Internal, "database error")
	}

	return exists[0], nil
}
