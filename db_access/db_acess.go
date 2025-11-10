package db_access

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

const Path = "data/FGG.db?_journal_mode=WAL"

func Init() {
	var err error
	db, err = sql.Open("sqlite3", Path)

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
