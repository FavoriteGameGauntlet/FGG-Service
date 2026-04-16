package dbusers

import (
	"FGG-Service/src/dbaccess"
)

type Database struct {
}

const ChangeDisplayNameQuery = `
	UPDATE Users
	SET DisplayName = ?
	WHERE Id = ?
`

func (db *Database) ChangeDisplayNameCommand(userId int, displayName string) error {
	queryName := "ChangeDisplayNameQuery"
	_, err := dbaccess.Exec(queryName, ChangeDisplayNameQuery, displayName, userId)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const GetDisplayNameQuery = `
	SELECT DisplayName
    FROM Users
	WHERE Id = ?
`

func (db *Database) GetDisplayNameCommand(userId int) (displayName *string, err error) {
	queryName := "GetDisplayNameQuery"
	row := dbaccess.QueryRow(queryName, GetDisplayNameQuery, userId)

	err = row.Scan(&displayName)

	dbaccess.LogDbResult(queryName, displayName, err)

	return
}
