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
	_, err := dbaccess.Exec(IncreaseAvailableRollsQuery, userId)

	return err
}

const IncreaseTerritoryHoursQuery = `
	UPDATE UserStats
	SET TerritoryHours = TerritoryHours + ?
	WHERE UserId = ?
`

func (db *Database) IncreaseTerritoryHoursCommand(userId int, changeValue int) error {
	_, err := dbaccess.Exec(IncreaseTerritoryHoursQuery, changeValue, userId)

	return err
}

const IncreaseExperiencePointsQuery = `
	UPDATE UserStats
	SET ExperiencePoints = ExperiencePoints + ?
	WHERE UserId = ?
`

func (db *Database) IncreaseExperiencePointsCommand(userId int, changeValue int) error {
	_, err := dbaccess.Exec(IncreaseExperiencePointsQuery, changeValue, userId)

	return err
}
