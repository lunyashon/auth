package validation

import (
	databasesmysql "main/internal/database/mysql/auth"

	protoc "github.com/lunyashon/protoc/auth/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (v Validate) RegisterToken(data *protoc.TokenRequest) error {

	// Валидация на пустоту
	if data.Token == "" {
		return status.Error(codes.InvalidArgument, "api is empty")
	}
	if data.Services == "" {
		return status.Error(codes.InvalidArgument, "services is empty")
	}

	// Валидация по всем остальным признакам
	if err := databasesmysql.DB.ValidateToken(&databasesmysql.StructDatabase{}, data); err != nil {
		return err
	}

	return nil
}
