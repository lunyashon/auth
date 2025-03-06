package databasesmysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func (db *StructDatabase) DBJoin() (*sql.DB, error) {

	dbName, dbData := db.DataDB()
	dbOpen, err := sql.Open(dbName, dbData)
	if err != nil {
		return nil, err
	}

	return dbOpen, nil
}

func (db *StructDatabase) DBClose(dbOpen *sql.DB) error {
	defer dbOpen.Close()
	return nil
}
