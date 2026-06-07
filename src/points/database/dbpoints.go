package dbpoints

import (
	"FGG-Service/src/dbaccess"
	typepoints "FGG-Service/src/points/type"
)

type IDatabase interface {
	GetExperiencePointsCommand(userId int) (points int, err error)
	ChangeExperiencePointsCommand(userId int, changeValue int) error
	GetFreePointsCommand(userId int) (points int, err error)
	ChangeFreePointsCommand(userId int, changeValue int) error
	AddFreePointHistoryCommand(
		userId int,
		sourceUserId int,
		changeSource string,
		changeValue int,
		actualChangeValue int,
		finalValue int,
		wheelEffectId *int) error
	GetTerritoryHoursCommand(userId int) (points int, err error)
	ChangeTerritoryHoursCommand(userId int, changeValue int) error
	GetTerritoryPointsCommand(userId int) (points int, err error)
	GetPointInfoCommand(userId int) (info typepoints.PointInfo, err error)
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

const GetTerritoryHoursQuery = `
	SELECT TerritoryHours
	FROM UserStats
	WHERE UserId = ?
`

func (db *Database) GetTerritoryHoursCommand(userId int) (points int, err error) {
	queryName := "GetTerritoryHoursQuery"
	row := dbaccess.QueryRow(queryName, GetTerritoryHoursQuery, userId)

	err = row.Scan(&points)

	dbaccess.LogDbResult(queryName, points, err)

	return
}

const ChangeTerritoryHoursQuery = `
	UPDATE UserStats
	SET TerritoryHours = TerritoryHours + ?
	WHERE UserId = ?
`

func (db *Database) ChangeTerritoryHoursCommand(userId int, changeValue int) error {
	queryName := "ChangeTerritoryHoursQuery"
	_, err := dbaccess.Exec(queryName, ChangeTerritoryHoursQuery, changeValue, userId)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const GetTerritoryPointsQuery = `
	SELECT TerritoryPoints
	FROM UserStats
	WHERE UserId = ?
`

func (db *Database) GetTerritoryPointsCommand(userId int) (points int, err error) {
	queryName := "GetTerritoryPointsQuery"
	row := dbaccess.QueryRow(queryName, GetTerritoryPointsQuery, userId)

	err = row.Scan(&points)

	dbaccess.LogDbResult(queryName, points, err)

	return
}

const GetPointInfoQuery = `
	SELECT TerritoryPoints, FreePoints, AvailableRolls, TerritoryHours, ExperiencePoints
	FROM UserStats
	WHERE UserId = ?
`

func (db *Database) GetPointInfoCommand(userId int) (info typepoints.PointInfo, err error) {
	queryName := "GetPointInfoQuery"
	row := dbaccess.QueryRow(queryName, GetPointInfoQuery, userId)

	err = row.Scan(
		&info.TerritoryPoints,
		&info.FreePoints,
		&info.AvailableRolls,
		&info.TerritoryHours,
		&info.ExperiencePoints)

	dbaccess.LogDbResult(queryName, info, err)

	return
}

const GetFreePointsQuery = `
	SELECT FreePoints
	FROM UserStats
	WHERE UserId = ?
`

func (db *Database) GetFreePointsCommand(userId int) (points int, err error) {
	queryName := "GetFreePointsQuery"
	row := dbaccess.QueryRow(queryName, GetFreePointsQuery, userId)

	err = row.Scan(&points)

	dbaccess.LogDbResult(queryName, points, err)

	return
}

const ChangeFreePointsQuery = `
	UPDATE UserStats
	SET FreePoints = FreePoints + ?
	WHERE UserId = ?
`

func (db *Database) ChangeFreePointsCommand(userId int, changeValue int) error {
	queryName := "ChangeFreePointsQuery"
	_, err := dbaccess.Exec(queryName, ChangeFreePointsQuery, changeValue, userId)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const AddFreePointHistoryQuery = `
	INSERT INTO FreePointHistory (
		UserId,
	    SourceUserId,
		ChangeSource,
		ChangeValue,
		ActualChangeValue,
		FinalValue,
	  	WheelEffectId
	)
	VALUES (?, ?, ?, ?, ?, ?, ?)
`

func (db *Database) AddFreePointHistoryCommand(
	userId int,
	sourceUserId int,
	changeSource string,
	changeValue int,
	actualChangeValue int,
	finalValue int,
	wheelEffectId *int) error {

	queryName := "AddFreePointHistoryQuery"
	_, err := dbaccess.Exec(
		queryName,
		AddFreePointHistoryQuery,
		userId,
		sourceUserId,
		changeSource,
		changeValue,
		actualChangeValue,
		finalValue,
		wheelEffectId)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}
