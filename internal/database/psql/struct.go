package database

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lunyashon/auth/internal/config"
)

type StructDatabase struct {
	Validator ValidateProvider
	Base      BaseProvider
	User      UserProvider
	Token     TokenProvider
	Cfg       config.ConfigEnv
}

type DatabaseProvider struct {
	db  *sqlx.DB
	cfg config.ConfigEnv
	log *slog.Logger
}

type ReturnOnceParam interface {
	string | int | bool | time.Time |
		sql.NullString | sql.NullInt64 | sql.NullBool | sql.NullTime
}
