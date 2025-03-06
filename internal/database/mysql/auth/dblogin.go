package databasesmysql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	hashpass "main/internal/services/grpc/hash"

	_ "github.com/go-sql-driver/mysql"
	protoc "github.com/lunyashon/protoc/auth/gen/go/sso"
)

// вывод иинформации по запросу
func (db *StructDatabase) GetUsers(data *protoc.LoginRequest) (bool, error) {

	dbOpen, err := db.DBJoin()
	if err != nil {
		fmt.Println(err)
	}

	success, err := getUsersDB(data, dbOpen)
	if err != nil {
		fmt.Println(err)
	}

	switch {
	case err != nil:
		db.DBClose(dbOpen)
		return false, errors.New("email incorrect")
	case !success:
		db.DBClose(dbOpen)
		return false, nil
	default:
		db.DBClose(dbOpen)
		return true, nil
	}
}

// Запрос к БД, чтобы увидеть есть ли ID с таким email
func getUsersDB(data *protoc.LoginRequest, db *sql.DB) (bool, error) {

	var pass string
	selectCount, err := db.Query("SELECT password FROM web_services.users WHERE login = \"" + data.Login + "\"")

	if err != nil {
		log.Println(err)
	}
	defer selectCount.Close()

	for selectCount.Next() {
		if err := selectCount.Scan(&pass); err != nil {
			log.Println(err)
		}
	}

	// Проврка на правильность пороля с помощью bcrypt
	if pass == "" {
		return false, errors.New("email incorrect")
	}

	return hashpass.VerifyUserPass(data.Password), nil
}
