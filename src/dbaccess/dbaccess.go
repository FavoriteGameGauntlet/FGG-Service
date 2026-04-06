package dbaccess

import (
	typeauth "FGG-Service/src/auth/types"
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
	queryArgs, logArgs := getArgs(args...)

	if queryName != "GetCompletedTimerUsersQuery" {
		slog.Info(queryName, "args", logArgs)
	}

	result, err := db.Exec(query, queryArgs...)

	if err != nil {
		slog.Error(queryName, "error", err)
	}

	return result, err
}

func Query(queryName string, query string, args ...any) (*sql.Rows, error) {
	queryArgs, logArgs := getArgs(args...)

	if queryName != "GetCompletedTimerUsersQuery" {
		slog.Info(queryName, "args", logArgs)
	}

	rows, err := db.Query(query, queryArgs...)

	if err != nil {
		slog.Error(queryName, "error", err)
	}

	return rows, err
}

func QueryRow(queryName string, query string, args ...any) *sql.Row {
	queryArgs, logArgs := getArgs(args...)

	if queryName != "GetCompletedTimerUsersQuery" {
		slog.Info(queryName, "args", logArgs)
	}

	return db.QueryRow(query, queryArgs...)
}

func getArgs(args ...any) ([]any, []any) {
	logArgs := make([]any, len(args))
	queryArgs := make([]any, len(args))

	for i, arg := range args {
		if p, ok := arg.(typeauth.Password); ok {
			logArgs[i] = "***"
			queryArgs[i] = p.Value
		} else {
			logArgs[i] = arg
			queryArgs[i] = arg
		}
	}

	return queryArgs, logArgs
}
