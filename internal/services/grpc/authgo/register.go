package authgo

import (
	"context"
	databasesmysql "main/internal/database/mysql/auth"
	"main/internal/services/grpc/validation"

	protoc "github.com/lunyashon/protoc/auth/gen/go/sso"
)

// Запись нового Юзера в БД
func (a AuthData) RegisterUser(ctx context.Context, data *protoc.RegisterRequest) (int64, error) {
	// Валидация данных
	if err := validation.Validate.Register(validation.Validate{}, data); err != nil {
		return 0, err
	}

	return databasesmysql.DB.CreateInDB(&databasesmysql.StructDatabase{}, data)
}
