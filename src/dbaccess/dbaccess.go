package dbaccess

import (
	"database/sql"
	"log/slog"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const Path = "file:data/FGG.db"

func Init() func() {
	var err error
	db, err = sql.Open("sqlite", Path)

	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(1)

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

	return func() {
		_ = db.Close()
	}
}

func Exec(queryName string, query string, args ...any) (sql.Result, error) {
	slog.Info(queryName, "args", args)

	result, err := db.Exec(query, args...)

	if err != nil {
		slog.Error(queryName, "error", err)
	}

	return result, err
}

func Query(queryName string, query string, args ...any) (*sql.Rows, error) {
	slog.Info(queryName, "args", args)

	rows, err := db.Query(query, args...)

	if err != nil {
		slog.Error(queryName, "error", err)
	}

	return rows, err
}

func QueryRow(queryName string, query string, args ...any) *sql.Row {
	slog.Info(queryName, "args", args)

	return db.QueryRow(query, args...)
}
