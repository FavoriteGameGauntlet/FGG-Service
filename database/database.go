package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	database *sql.DB
}

const DatabasePath = "database/FavoriteGameGauntlet.db"

func NewDatabase() *Database {
	db, err := sql.Open("sqlite3", DatabasePath)

	if err != nil {
		panic(err)
	}

	db.Exec("SELECT load_extension('uuid')")

	return &Database{
		database: db,
	}
}

func (db *Database) Close() {
	db.database.Close()
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
