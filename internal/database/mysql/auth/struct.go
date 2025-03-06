package databasesmysql

import (
	"database/sql"

	protoc "github.com/lunyashon/protoc/auth/gen/go/sso"
)

type StructDatabase struct {
	dataDB DataDB
	db     DB
}

type DataDB interface {
	GetData() (string, string)
}

type DB interface {
	DBJoin() (*sql.DB, error)
	DBClose(dbOpen *sql.DB) error
	CreateInDB(data *protoc.RegisterRequest) (int64, error)
	GetUsers(data *protoc.LoginRequest) (bool, error)
	Validate(*protoc.RegisterRequest) error
	ValidateToken(*protoc.TokenRequest) error
	CreateToken(data *protoc.TokenRequest) (bool, error)
}

func (s *StructDatabase) DataDB() (string, string) {
	return "mysql", "admin:2024T,bnt,f,yfcdt;tvctyt@/web_services"
}
