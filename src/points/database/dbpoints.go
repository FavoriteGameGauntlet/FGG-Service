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

const IncreaseAvailableRollsQuery = `SELECT increase_available_rolls($1::integer)`

func (db *Database) IncreaseAvailableRollsCommand(userId int) error {
	queryName := "IncreaseAvailableRollsQuery"
	_, err := dbaccess.Exec(queryName, IncreaseAvailableRollsQuery, userId)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const IncreaseTerritoryHoursQuery = `SELECT increase_territory_hours($1::integer, $2::integer)`

func (db *Database) IncreaseTerritoryHoursCommand(userId int, changeValue int) error {
	queryName := "IncreaseTerritoryHoursQuery"
	_, err := dbaccess.Exec(queryName, IncreaseTerritoryHoursQuery, userId, changeValue)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const ChangeExperiencePointsQuery = `SELECT change_experience_points($1::integer, $2::integer)`

func (db *Database) ChangeExperiencePointsCommand(userId int, changeValue int) error {
	queryName := "ChangeExperiencePointsQuery"
	_, err := dbaccess.Exec(queryName, ChangeExperiencePointsQuery, userId, changeValue)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const GetExperiencePointsQuery = `SELECT * FROM get_experience_points($1::integer)`

func (db *Database) GetExperiencePointsCommand(userId int) (points int, err error) {
	queryName := "ExperiencePointsQuery"
	row := dbaccess.QueryRow(queryName, GetExperiencePointsQuery, userId)

	err = row.Scan(&points)

	dbaccess.LogDbResult(queryName, points, err)

	return
}

const GetTerritoryHoursQuery = `SELECT * FROM get_territory_hours($1::integer)`

func (db *Database) GetTerritoryHoursCommand(userId int) (points int, err error) {
	queryName := "GetTerritoryHoursQuery"
	row := dbaccess.QueryRow(queryName, GetTerritoryHoursQuery, userId)

	err = row.Scan(&points)

	dbaccess.LogDbResult(queryName, points, err)

	return
}

const ChangeTerritoryHoursQuery = `SELECT change_territory_hours($1::integer, $2::integer)`

func (db *Database) ChangeTerritoryHoursCommand(userId int, changeValue int) error {
	queryName := "ChangeTerritoryHoursQuery"
	_, err := dbaccess.Exec(queryName, ChangeTerritoryHoursQuery, userId, changeValue)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const GetTerritoryPointsQuery = `SELECT * FROM get_territory_points($1::integer)`

func (db *Database) GetTerritoryPointsCommand(userId int) (points int, err error) {
	queryName := "GetTerritoryPointsQuery"
	row := dbaccess.QueryRow(queryName, GetTerritoryPointsQuery, userId)

	err = row.Scan(&points)

	dbaccess.LogDbResult(queryName, points, err)

	return
}

const ChangeTerritoryPointsQuery = `SELECT change_territory_points($1::integer, $2::integer)`

func (db *Database) ChangeTerritoryPointsCommand(userId int, changeValue int) error {
	queryName := "ChangeTerritoryPointsQuery"
	_, err := dbaccess.Exec(queryName, ChangeTerritoryPointsQuery, userId, changeValue)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const AddTerritoryPointHistoryQuery = `SELECT add_territory_point_history($1::integer, $2::integer, $3::text, $4::integer, $5::integer, $6::integer)`

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

const GetTerritoryPointHistoryQuery = `SELECT * FROM get_territory_point_history($1::integer)`

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

const GetPointInfoQuery = `SELECT * FROM get_point_info($1::integer)`

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

const GetAllPointInfoQuery = `SELECT * FROM get_all_point_info()`

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

const GetFreePointsQuery = `SELECT * FROM get_free_points($1::integer)`

func (db *Database) GetFreePointsCommand(userId int) (points int, err error) {
	queryName := "GetFreePointsQuery"
	row := dbaccess.QueryRow(queryName, GetFreePointsQuery, userId)

	err = row.Scan(&points)

	dbaccess.LogDbResult(queryName, points, err)

	return
}

const ChangeFreePointsQuery = `SELECT change_free_points($1::integer, $2::integer)`

func (db *Database) ChangeFreePointsCommand(userId int, changeValue int) error {
	queryName := "ChangeFreePointsQuery"
	_, err := dbaccess.Exec(queryName, ChangeFreePointsQuery, userId, changeValue)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const AddFreePointHistoryQuery = `SELECT add_free_point_history($1::integer, $2::integer, $3::text, $4::integer, $5::integer, $6::integer, $7::integer)`

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

const GetFreePointHistoryQuery = `SELECT * FROM get_free_point_history($1::integer)`

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
