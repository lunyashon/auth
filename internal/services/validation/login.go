package validate

import (
	"context"

	database "github.com/lunyashon/auth/internal/database/psql"

	sso "github.com/lunyashon/protobuf/auth/gen/go/sso/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Validation of input data
// Returning User Auth data or error
func Auth(
	ctx context.Context,
	data *sso.LoginRequest,
	db *database.StructDatabase,
) error {
	if err := checkLogin(data.Login); err != nil {
		return err
	}
	if err := checkPassword(data.Password); err != nil {
		return err
	}
	return nil
}

// Checking for emptiness login
// Return error or nil
func checkLogin(login string) error {
	if login == "" {
		return status.Error(codes.InvalidArgument, "the login is empty")
	}
	return nil
}

// Checking for emptiness password
// Return error or nil
func checkPassword(password string) error {
	if password == "" {
		return status.Error(codes.InvalidArgument, "the password is empty")
	}
	return nil
}
