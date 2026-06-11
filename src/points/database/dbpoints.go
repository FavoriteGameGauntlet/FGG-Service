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
	GetFreePointHistoryCommand(userId int) (history typepoints.FreePointChangeHistories, err error)
	GetTerritoryHoursCommand(userId int) (points int, err error)
	ChangeTerritoryHoursCommand(userId int, changeValue int) error
	GetTerritoryPointsCommand(userId int) (points int, err error)
	ChangeTerritoryPointsCommand(userId int, changeValue int) error
	AddTerritoryPointHistoryCommand(
		userId int,
		sourceUserId int,
		changeSource string,
		changeValue int,
		actualChangeValue int,
		finalValue int) error
	GetTerritoryPointHistoryCommand(userId int) (history typepoints.TerritoryPointChangeHistories, err error)
	GetPointInfoCommand(userId int) (info typepoints.PointInfo, err error)
	GetAllPointInfoCommand() (infos typepoints.PointInfoByLogins, err error)
}

type Database struct {
}

const IncreaseAvailableRollsQuery = `
	UPDATE UserStats
	SET AvailableRolls = AvailableRolls + 1
	WHERE UserId = $1
`

func (db *Database) IncreaseAvailableRollsCommand(userId int) error {
	queryName := "IncreaseAvailableRollsQuery"
	_, err := dbaccess.Exec(queryName, IncreaseAvailableRollsQuery, userId)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const IncreaseTerritoryHoursQuery = `
	UPDATE UserStats
	SET TerritoryHours = TerritoryHours + $1
	WHERE UserId = $2
`

func (db *Database) IncreaseTerritoryHoursCommand(userId int, changeValue int) error {
	queryName := "IncreaseTerritoryHoursQuery"
	_, err := dbaccess.Exec(queryName, IncreaseTerritoryHoursQuery, changeValue, userId)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const ChangeExperiencePointsQuery = `
	UPDATE UserStats
	SET ExperiencePoints = ExperiencePoints + $1
	WHERE UserId = $2
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
	WHERE UserId = $1
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
	WHERE UserId = $1
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
	SET TerritoryHours = TerritoryHours + $1
	WHERE UserId = $2
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
	WHERE UserId = $1
`

func (db *Database) GetTerritoryPointsCommand(userId int) (points int, err error) {
	queryName := "GetTerritoryPointsQuery"
	row := dbaccess.QueryRow(queryName, GetTerritoryPointsQuery, userId)

	err = row.Scan(&points)

	dbaccess.LogDbResult(queryName, points, err)

	return
}

const ChangeTerritoryPointsQuery = `
	UPDATE UserStats
	SET TerritoryPoints = TerritoryPoints + $1
	WHERE UserId = $2
`

func (db *Database) ChangeTerritoryPointsCommand(userId int, changeValue int) error {
	queryName := "ChangeTerritoryPointsQuery"
	_, err := dbaccess.Exec(queryName, ChangeTerritoryPointsQuery, changeValue, userId)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const AddTerritoryPointHistoryQuery = `
	INSERT INTO TerritoryPointHistory (
		UserId,
		SourceUserId,
		ChangeSource,
		ChangeValue,
		ActualChangeValue,
		FinalValue
	)
	VALUES ($1, $2, $3, $4, $5, $6)
`

func (db *Database) AddTerritoryPointHistoryCommand(
	userId int,
	sourceUserId int,
	changeSource string,
	changeValue int,
	actualChangeValue int,
	finalValue int) error {

	queryName := "AddTerritoryPointHistoryQuery"
	_, err := dbaccess.Exec(
		queryName,
		AddTerritoryPointHistoryQuery,
		userId,
		sourceUserId,
		changeSource,
		changeValue,
		actualChangeValue,
		finalValue)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const GetTerritoryPointHistoryQuery = `
	SELECT
		tph.ActualChangeValue,
		tph.ChangeDate,
		tph.ChangeSource,
		tph.ChangeValue,
		tph.FinalValue,
		u.Login
	FROM TerritoryPointHistory tph
		LEFT JOIN Users u ON u.Id = tph.SourceUserId
	WHERE tph.UserId = $1
	ORDER BY tph.ChangeDate DESC
`

func (db *Database) GetTerritoryPointHistoryCommand(userId int) (
	history typepoints.TerritoryPointChangeHistories, err error) {

	queryName := "GetTerritoryPointHistoryQuery"
	rows, err := dbaccess.Query(queryName, GetTerritoryPointHistoryQuery, userId)

	if err != nil {
		return
	}

	for rows.Next() {
		entry := typepoints.TerritoryPointChangeHistory{}
		err = rows.Scan(
			&entry.ActualChangeValue,
			&entry.ChangeDate,
			&entry.ChangeSource,
			&entry.DesiredChangeValue,
			&entry.FinalValue,
			&entry.SourceLogin)

		if err != nil {
			_ = rows.Close()
			return
		}

		history = append(history, entry)
	}

	dbaccess.LogDbResult(queryName, history, err)

	_ = rows.Close()
	return
}

const GetPointInfoQuery = `
	SELECT TerritoryPoints, FreePoints, AvailableRolls, TerritoryHours, ExperiencePoints
	FROM UserStats
	WHERE UserId = $1
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

const GetAllPointInfoQuery = `
	SELECT u.Login, us.TerritoryPoints, us.FreePoints, us.AvailableRolls, us.TerritoryHours, us.ExperiencePoints
	FROM Users u
	INNER JOIN UserStats us ON us.UserId = u.Id
`

func (db *Database) GetAllPointInfoCommand() (infos typepoints.PointInfoByLogins, err error) {
	queryName := "GetAllPointInfoQuery"
	rows, err := dbaccess.Query(queryName, GetAllPointInfoQuery)

	if err != nil {
		return
	}

	for rows.Next() {
		info := typepoints.PointInfoByLogin{}
		err = rows.Scan(
			&info.Login,
			&info.PointInfo.TerritoryPoints,
			&info.PointInfo.FreePoints,
			&info.PointInfo.AvailableRolls,
			&info.PointInfo.TerritoryHours,
			&info.PointInfo.ExperiencePoints)

		if err != nil {
			_ = rows.Close()
			return
		}

		infos = append(infos, info)
	}

	dbaccess.LogDbResult(queryName, infos, err)

	_ = rows.Close()
	return
}

const GetFreePointsQuery = `
	SELECT FreePoints
	FROM UserStats
	WHERE UserId = $1
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
	SET FreePoints = FreePoints + $1
	WHERE UserId = $2
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
	VALUES ($1, $2, $3, $4, $5, $6, $7)
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

const GetFreePointHistoryQuery = `
	SELECT
		fph.ActualChangeValue,
		fph.ChangeDate,
		fph.ChangeSource,
		fph.ChangeValue,
		fph.FinalValue,
		u.Login,
		we.Name
	FROM FreePointHistory fph
		LEFT JOIN Users u ON u.Id = fph.SourceUserId
		LEFT JOIN WheelEffects we ON we.Id = fph.WheelEffectId
	WHERE fph.UserId = $1
	ORDER BY fph.ChangeDate DESC
`

func (db *Database) GetFreePointHistoryCommand(userId int) (history typepoints.FreePointChangeHistories, err error) {
	queryName := "GetFreePointHistoryQuery"
	rows, err := dbaccess.Query(queryName, GetFreePointHistoryQuery, userId)

	if err != nil {
		return
	}

	for rows.Next() {
		entry := typepoints.FreePointChangeHistory{}
		err = rows.Scan(
			&entry.ActualChangeValue,
			&entry.ChangeDate,
			&entry.ChangeSource,
			&entry.DesiredChangeValue,
			&entry.FinalValue,
			&entry.SourceLogin,
			&entry.WheelEffectName)

		if err != nil {
			_ = rows.Close()
			return
		}

		history = append(history, entry)
	}

	dbaccess.LogDbResult(queryName, history, err)

	_ = rows.Close()
	return
}
