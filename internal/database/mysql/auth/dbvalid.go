package databasesmysql

import (
	"database/sql"
	"strings"

	protoc "github.com/lunyashon/protoc/auth/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Валидация регистрации юзера
func (db *StructDatabase) Validate(data *protoc.RegisterRequest) error {

	var count int

	dbOpen, err := db.DBJoin()

	if err != nil {
		return err
	}

	if err := validateIsAlready(data.Email, "email", dbOpen); err != nil {
		return err
	}

	if err := validateIsAlready(data.Login, "login", dbOpen); err != nil {
		return err
	}

	count, err = validateItem(data.Api, "token", dbOpen)
	if count == 0 {
		return status.Error(codes.InvalidArgument, "token "+data.Api+" not exist")
	}
	if err != nil {
		return err
	}

	if err := validateServices(data.Services, dbOpen); err != nil {
		return err
	}

	db.DBClose(dbOpen)
	return nil
}

// Валидация данных регистрации пользователя
func (db *StructDatabase) ValidateToken(data *protoc.TokenRequest) error {

	var count int

	dbOpen, err := db.DBJoin()

	if err != nil {
		return err
	}

	if err := validateServices(data.Services, dbOpen); err != nil {
		return err
	}

	count, err = validateItem(data.Token, "token", dbOpen)
	if count > 0 {
		return status.Error(codes.InvalidArgument, "token "+data.Token+" already exist")
	}
	if err != nil {
		return err
	}

	db.DBClose(dbOpen)
	return nil
}

// Проверка Email и Login
func validateIsAlready(item, components string, dbOpen *sql.DB) error {
	var count int
	selectCount, err := dbOpen.Query("SELECT COUNT(id) FROM web_services.users WHERE " + components + " = \"" + item + "\"")

	if err != nil {
		return err
	}
	defer selectCount.Close()

	for selectCount.Next() {
		if err := selectCount.Scan(&count); err != nil {
			return err
		}
	}

	if count > 0 {
		return status.Error(codes.InvalidArgument, components+" already there")
	}
	return nil
}

// Проверка Services, api token
func validateItem(item, components string, dbOpen *sql.DB) (int, error) {

	var (
		count int
		q     string
	)

	if components == "token" {
		q = "SELECT COUNT(token) FROM tokens WHERE token = \"" + item + "\""
	} else {
		q = "SELECT COUNT(id) FROM services WHERE services_name = \"" + item + "\""
	}

	selectCount, err := dbOpen.Query(q)

	if err != nil {
		return -1, err
	}
	defer selectCount.Close()

	for selectCount.Next() {
		if err := selectCount.Scan(&count); err != nil {
			return -1, err
		}
	}

	if components == "services" {
		if count == 0 {
			return -1, status.Error(codes.InvalidArgument, components+": "+item+" not in table")
		}
		return -1, nil
	} else {
		return count, nil
	}
}

// Оболочка под проверку services
func validateServices(item string, dbOpen *sql.DB) error {
	arr := strings.Split(item, ",")
	for _, val := range arr {
		val = strings.TrimSpace(val)
		if _, err := validateItem(val, "services", dbOpen); err != nil {
			return err
		}
	}

	return nil
}
