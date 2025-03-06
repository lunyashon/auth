package databasesmysql

import (
	"database/sql"
	hashpass "main/internal/services/grpc/hash"

	protoc "github.com/lunyashon/protoc/auth/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Регистрация пользователя
func (db *StructDatabase) CreateInDB(data *protoc.RegisterRequest) (int64, error) {

	dbOpen, err := db.DBJoin()
	if err != nil {
		return 0, err
	}
	id, err := insertDB(data, dbOpen, db)
	if err != nil {
		return 0, err
	}
	db.DBClose(dbOpen)
	return id, nil
}

// Запись в БД нового пользователя
func insertDB(data *protoc.RegisterRequest, dbOpen *sql.DB, db *StructDatabase) (int64, error) {

	hashPass := string(hashpass.HashSHA256([]byte(data.Password)))

	q := `INSERT INTO users ` +
		`SET email = ?, ` +
		`login = ?, ` +
		`password = ?, ` +
		`services = ?`

	result, err := dbOpen.Exec(q, data.Email, data.Login, hashPass, data.Services)
	if err != nil {
		return 0, err
	}
	db.DBClose(dbOpen)

	id, err := result.LastInsertId()
	if err != nil {
		return 0, status.Error(codes.InvalidArgument, "error")
	}
	return id, nil
}
