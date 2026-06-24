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

var changeDisplayNameQuery = dbaccess.Query{Name: "ChangeDisplayNameQuery", SQL: `SELECT change_display_name($1::integer, $2::text)`}

func (db *Database) ChangeDisplayNameCommand(userId int, displayName string) error {
	_, err := dbaccess.Exec(changeDisplayNameQuery, userId, displayName)

	dbaccess.LogDbResult(changeDisplayNameQuery, nil, err)

	return err
}

var getDisplayNameQuery = dbaccess.Query{Name: "GetDisplayNameQuery", SQL: `SELECT * FROM get_display_name($1::integer)`}

func (db *Database) GetDisplayNameCommand(userId int) (displayName *string, err error) {
	row := dbaccess.QueryRow(getDisplayNameQuery, userId)

	err = row.Scan(&displayName)

	dbaccess.LogDbResult(getDisplayNameQuery, displayName, err)

	return
}

var getAllUserNamesQuery = dbaccess.Query{Name: "GetAllUserNamesQuery", SQL: `SELECT * FROM get_all_user_names()`}

func (db *Database) GetAllUserNamesCommand() (users typeusers.Users, err error) {
	rows, err := dbaccess.QueryRows(getAllUserNamesQuery)

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

	dbaccess.LogDbResult(getAllUserNamesQuery, users, err)

	_ = rows.Close()
	return
}
