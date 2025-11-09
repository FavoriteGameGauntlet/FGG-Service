package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	database *sql.DB
}

const Path = "FGG.db"

func NewDatabase() *Database {
	db, err := sql.Open("sqlite3", Path)

	if err != nil {
		panic(err)
	}

	return &Database{
		database: db,
	}
}

func (db *Database) Close() {
	err := db.database.Close()

	if err != nil {
		panic(err)
	}
}

func Exec(query string, args ...any) (sql.Result, error) {
	db := NewDatabase()
	defer db.Close()

	return db.database.Exec(query, args...)
}

func Query(query string, args ...any) (*sql.Rows, error) {
	db := NewDatabase()
	defer db.Close()

	return db.database.Query(query, args...)
}

func QueryRow(query string, args ...any) *sql.Row {
	db := NewDatabase()
	defer db.Close()

	row := db.database.QueryRow(query, args...)

	return row
}
