package databasesmysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	protoc "github.com/lunyashon/protoc/auth/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// вывод иинформации по запросу
func (db *StructDatabase) CreateToken(data *protoc.TokenRequest) (bool, error) {

	dbOpen, err := db.DBJoin()
	if err != nil {
		fmt.Println(err)
	}

	return createToken(data, dbOpen, db)
}

// Запрос к БД, чтобы увидеть есть ли ID с таким email
func createToken(data *protoc.TokenRequest, dbOpen *sql.DB, db *StructDatabase) (bool, error) {

	fmt.Println(data)

	q := `INSERT INTO tokens (token, services) VALUES(?,?)`

	result, err := dbOpen.Exec(q, data.Token, data.Services)
	if err != nil {
		return false, err
	}
	db.DBClose(dbOpen)

	count, err := result.RowsAffected()
	if err != nil {
		status.Error(codes.InvalidArgument, "error")
	}

	if count > 0 {
		return true, nil
	}

	return false, status.Error(codes.InvalidArgument, "not created")
}
