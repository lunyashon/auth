package database

import (
	"context"

	"github.com/jmoiron/sqlx"
	passauth "github.com/lunyashon/auth/internal/lib/passauth"
	"github.com/lunyashon/auth/internal/services/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Verifying user credentials
// Return model UserAuth or error
func (s *DatabaseProvider) CheckUser(ctx context.Context, login, pass string) (*model.UserAuth, error) {

	var (
		authData   = make([]userAuthData, 2)
		methodName = "CheckUser"
	)

	q := `SELECT u.password, s.name, u.id 
		  FROM users u 
		  JOIN permission p ON u.id = p.user_id 
		  JOIN services s ON s.id = p.service_id 
		  WHERE u.login = $1`

	if err := sqlx.SelectContext(ctx, s.db, &authData, q, login); err != nil {
		errMessage := status.Errorf(codes.Internal, "failed to query: %v", err)
		s.log.ErrorContext(
			ctx,
			"ERROR database",
			"method", methodName,
			"point", point,
			"login", login,
			"message", errMessage.Error(),
		)
		return nil, status.Error(codes.Internal, "database error")
	}

	var user model.UserAuth

	if len(authData) > 0 {
		if authData[0].Pass != "" {
			pas := passauth.ExecAuthService(&passauth.RealBcrypt{})
			if err := pas.VerifyPassword(pass, authData[0].Pass); err != nil {
				return nil, status.Error(codes.InvalidArgument, "incorect password")
			}
		}

		// user ID
		user.UID = authData[0].Uid

		for _, val := range authData {
			user.Services = append(user.Services, val.Service)
		}

		return &user, nil
	}

	return nil, status.Errorf(codes.NotFound, "login %v not exist", login)
}
