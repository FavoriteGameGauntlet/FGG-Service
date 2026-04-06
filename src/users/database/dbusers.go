package dbusers

import "FGG-Service/src/dbaccess"

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
