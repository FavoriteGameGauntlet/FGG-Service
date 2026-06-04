package dbusers

import (
	"FGG-Service/src/dbaccess"
	"FGG-Service/src/users/types"
)

type IDatabase interface {
	GetAllUserNamesCommand() (users typeusers.Users, err error)
	GetDisplayNameCommand(userId int) (displayName *string, err error)
	ChangeDisplayNameCommand(userId int, displayName string) error
}

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

const GetAllUserNamesQuery = `
	SELECT Login, DisplayName
	FROM Users
`

func (db *Database) GetAllUserNamesCommand() (users typeusers.Users, err error) {
	queryName := "GetAllUserNamesQuery"
	rows, err := dbaccess.Query(queryName, GetAllUserNamesQuery)

	if err != nil {
		return
	}

	for rows.Next() {
		user := typeusers.User{}
		err = rows.Scan(&user.Login, &user.DisplayName)

		if err != nil {
			_ = rows.Close()
			return
		}

		users = append(users, user)
	}

	dbaccess.LogDbResult(queryName, users, err)

	_ = rows.Close()
	return
}
