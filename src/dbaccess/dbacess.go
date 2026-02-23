package dbaccess

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const Path = "file:data/FGG.db"

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

	_, err = db.Exec("PRAGMA journal_mode=WAL;")

	if err != nil {
		panic(err)
	}

	_, err = db.Exec("PRAGMA busy_timeout=5000;")

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
