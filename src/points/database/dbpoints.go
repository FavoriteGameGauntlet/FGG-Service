package dbpoints

import "FGG-Service/src/dbaccess"

type Database struct {
}

const IncreaseAvailableRollsQuery = `
	UPDATE UserStats
	SET AvailableRolls = AvailableRolls + 1
	WHERE UserId = ?
`

func (db *Database) IncreaseAvailableRollsCommand(userId int) error {
	queryName := "IncreaseAvailableRollsQuery"
	_, err := dbaccess.Exec(queryName, IncreaseAvailableRollsQuery, userId)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const IncreaseTerritoryHoursQuery = `
	UPDATE UserStats
	SET TerritoryHours = TerritoryHours + ?
	WHERE UserId = ?
`

func (db *Database) IncreaseTerritoryHoursCommand(userId int, changeValue int) error {
	queryName := "IncreaseTerritoryHoursQuery"
	_, err := dbaccess.Exec(queryName, IncreaseTerritoryHoursQuery, changeValue, userId)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const IncreaseExperiencePointsQuery = `
	UPDATE UserStats
	SET ExperiencePoints = ExperiencePoints + ?
	WHERE UserId = ?
`

func (db *Database) IncreaseExperiencePointsCommand(userId int, changeValue int) error {
	queryName := "IncreaseExperiencePointsQuery"
	_, err := dbaccess.Exec(queryName, IncreaseExperiencePointsQuery, changeValue, userId)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}
