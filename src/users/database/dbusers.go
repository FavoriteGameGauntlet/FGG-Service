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
	_, err := dbaccess.Exec(ChangeDisplayNameQuery, displayName, userId)

	return err
}

const GetDisplayNameQuery = `
	SELECT DisplayName
    FROM Users
	WHERE Id = ?
`

func (db *Database) GetDisplayNameCommand(userId int) (displayName string, err error) {
	row := dbaccess.QueryRow(GetDisplayNameQuery, userId)

	err = row.Scan(&displayName)

	return
}
