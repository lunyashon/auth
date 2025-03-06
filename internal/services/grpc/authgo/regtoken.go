package authgo

import (
	"context"
	databasesmysql "main/internal/database/mysql/auth"
	"main/internal/services/grpc/validation"

	protoc "github.com/lunyashon/protoc/auth/gen/go/sso"
)

// Запись токена в БД
func (a AuthData) RegisterToken(ctx context.Context, data *protoc.TokenRequest) (bool, error) {

	// Валидация данных при регистрации
	if err := validation.Validate.RegisterToken(validation.Validate{}, data); err != nil {
		return false, err
	}

	return databasesmysql.DB.CreateToken(&databasesmysql.StructDatabase{}, data)
}
