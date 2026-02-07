package db_access

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const Path = "file:data/FGG.db?journal_mode=WAL&busy_timeout=5000"

func Init() {
	var err error
	db, err = sql.Open("sqlite", Path)

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}
}

func Close() {
	if db != nil {
		_ = db.Close()
	}
}

func Exec(query string, args ...any) (sql.Result, error) {
	return db.Exec(query, args...)
}

func Query(query string, args ...any) (*sql.Rows, error) {
	return db.Query(query, args...)
}

func QueryRow(query string, args ...any) *sql.Row {
	return db.QueryRow(query, args...)
}
