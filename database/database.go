package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	database *sql.DB
}

const DatabaseName = "FavoriteGameGauntlet.db"

func NewDatabase() *Database {
	db, err := sql.Open("sqlite3", DatabaseName)

	if err != nil {
		panic(err)
	}

	return &Database{
		database: db,
	}
}

func (db *Database) Close() {
	db.database.Close()
}

func Exec(query string, args ...any) sql.Result {
	db := NewDatabase()
	defer db.Close()

	result, err := db.database.Exec(query, args...)

	if err != nil {
		panic(err)
	}

	return result
}

func Query(query string, args ...any) *sql.Rows {
	db := NewDatabase()
	defer db.Close()

	rows, err := db.database.Query(query, args...)

	if err != nil {
		panic(err)
	}

	return rows
}

func QueryRow(query string, args ...any) *sql.Row {
	db := NewDatabase()
	defer db.Close()

	row := db.database.QueryRow(query, args...)

	return row
}
