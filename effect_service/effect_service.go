package effect_service

import (
	"FGG-Service/db_access"
)

const (
	CheckIfAvailableRollExistsCommand = `
		SELECT CASE 
        	WHEN EXISTS (
				SELECT 1
				FROM AvailableRolls
				WHERE UserId = $userId)
         	THEN true
         	ELSE false
       	END AS 'DoesExist'`
)

func CheckIfAvailableRollExists(userId int) (bool, error) {
	row := db_access.QueryRow(CheckIfAvailableRollExistsCommand, userId)

	var doesExist bool
	err := row.Scan(&doesExist)

	if err != nil {
		return doesExist, err
	}

	return doesExist, nil
}
