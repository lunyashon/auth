package validation

import (
	databasesmysql "main/internal/database/mysql/auth"

	protoc "github.com/lunyashon/protoc/auth/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (v Validate) Register(data *protoc.RegisterRequest) error {

	// Валидация на пустоту
	if data.Login == "" {
		return status.Error(codes.InvalidArgument, "login is empty")
	}
	if data.Password == "" {
		return status.Error(codes.InvalidArgument, "password is empty")
	}
	if data.Api == "" {
		return status.Error(codes.InvalidArgument, "api is empty")
	}
	if data.Email == "" {
		return status.Error(codes.InvalidArgument, "email is empty")
	}
	if data.Services == "" {
		return status.Error(codes.InvalidArgument, "services is empty")
	}

	// Валидация по всем остальным признакам
	if err := databasesmysql.DB.Validate(&databasesmysql.StructDatabase{}, data); err != nil {
		return err
	}

	return nil
}
