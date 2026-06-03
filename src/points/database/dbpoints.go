package dbpoints

import "FGG-Service/src/dbaccess"

type IDatabase interface {
	GetExperiencePointsCommand(userId int) (points int, err error)
	ChangeExperiencePointsCommand(userId int, changeValue int) error
}

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

const ChangeExperiencePointsQuery = `
	UPDATE UserStats
	SET ExperiencePoints = ExperiencePoints + ?
	WHERE UserId = ?
`

func (db *Database) ChangeExperiencePointsCommand(userId int, changeValue int) error {
	queryName := "ChangeExperiencePointsQuery"
	_, err := dbaccess.Exec(queryName, ChangeExperiencePointsQuery, changeValue, userId)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const GetExperiencePointsQuery = `
	SELECT ExperiencePoints
	FROM UserStats
	WHERE UserId = ?
`

func (db *Database) GetExperiencePointsCommand(userId int) (points int, err error) {
	queryName := "ExperiencePointsQuery"
	row := dbaccess.QueryRow(queryName, GetExperiencePointsQuery, userId)

	err = row.Scan(&points)

	dbaccess.LogDbResult(queryName, points, err)

	return
}
