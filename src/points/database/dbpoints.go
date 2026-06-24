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

var changeAvailableRollsQuery = dbaccess.Query{Name: "ChangeAvailableRollsQuery", SQL: `SELECT change_available_rolls($1::integer, $2::integer)`}

func (db *Database) ChangeAvailableRollsCommand(userId int, changeValue int) error {
	_, err := dbaccess.Exec(changeAvailableRollsQuery, userId, changeValue)

	dbaccess.LogDbResult(changeAvailableRollsQuery, nil, err)

	return err
}

var changeExperiencePointsQuery = dbaccess.Query{Name: "ChangeExperiencePointsQuery", SQL: `SELECT change_experience_points($1::integer, $2::integer)`}

func (db *Database) ChangeExperiencePointsCommand(userId int, changeValue int) error {
	_, err := dbaccess.Exec(changeExperiencePointsQuery, userId, changeValue)

	dbaccess.LogDbResult(changeExperiencePointsQuery, nil, err)

	return err
}

var getExperiencePointsQuery = dbaccess.Query{Name: "ExperiencePointsQuery", SQL: `SELECT * FROM get_experience_points($1::integer)`}

func (db *Database) GetExperiencePointsCommand(userId int) (points int, err error) {
	row := dbaccess.QueryRow(getExperiencePointsQuery, userId)

	err = row.Scan(&points)

	dbaccess.LogDbResult(getExperiencePointsQuery, points, err)

	return
}

var getTerritoryHoursQuery = dbaccess.Query{Name: "GetTerritoryHoursQuery", SQL: `SELECT * FROM get_territory_hours($1::integer)`}

func (db *Database) GetTerritoryHoursCommand(userId int) (points int, err error) {
	row := dbaccess.QueryRow(getTerritoryHoursQuery, userId)

	err = row.Scan(&points)

	dbaccess.LogDbResult(getTerritoryHoursQuery, points, err)

	return
}

var changeTerritoryHoursQuery = dbaccess.Query{Name: "ChangeTerritoryHoursQuery", SQL: `SELECT change_territory_hours($1::integer, $2::integer)`}

func (db *Database) ChangeTerritoryHoursCommand(userId int, changeValue int) error {
	_, err := dbaccess.Exec(changeTerritoryHoursQuery, userId, changeValue)

	dbaccess.LogDbResult(changeTerritoryHoursQuery, nil, err)

	return err
}

var getTerritoryPointsQuery = dbaccess.Query{Name: "GetTerritoryPointsQuery", SQL: `SELECT * FROM get_territory_points($1::integer)`}

func (db *Database) GetTerritoryPointsCommand(userId int) (points int, err error) {
	row := dbaccess.QueryRow(getTerritoryPointsQuery, userId)

	err = row.Scan(&points)

	dbaccess.LogDbResult(getTerritoryPointsQuery, points, err)

	return
}

var changeTerritoryPointsQuery = dbaccess.Query{Name: "ChangeTerritoryPointsQuery", SQL: `SELECT change_territory_points($1::integer, $2::integer)`}

func (db *Database) ChangeTerritoryPointsCommand(userId int, changeValue int) error {
	_, err := dbaccess.Exec(changeTerritoryPointsQuery, userId, changeValue)

	dbaccess.LogDbResult(changeTerritoryPointsQuery, nil, err)

	return err
}

var addTerritoryPointHistoryQuery = dbaccess.Query{Name: "AddTerritoryPointHistoryQuery", SQL: `SELECT add_territory_point_history($1::integer, $2::integer, $3::text, $4::integer, $5::integer, $6::integer)`}

func (db *Database) AddTerritoryPointHistoryCommand(
	userId int,
	sourceUserId int,
	changeSource string,
	changeValue int,
	actualChangeValue int,
	finalValue int) error {

	_, err := dbaccess.Exec(
		addTerritoryPointHistoryQuery,
		userId,
		sourceUserId,
		changeSource,
		changeValue,
		actualChangeValue,
		finalValue)

	dbaccess.LogDbResult(addTerritoryPointHistoryQuery, nil, err)

	return err
}

var getTerritoryPointHistoryQuery = dbaccess.Query{Name: "GetTerritoryPointHistoryQuery", SQL: `SELECT * FROM get_territory_point_history($1::integer)`}

func (db *Database) GetTerritoryPointHistoryCommand(userId int) (
	history typepoints.TerritoryPointChangeHistories, err error) {

	rows, err := dbaccess.QueryRows(getTerritoryPointHistoryQuery, userId)

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

	dbaccess.LogDbResult(getTerritoryPointHistoryQuery, history, err)

	_ = rows.Close()
	return
}

var getPointInfoQuery = dbaccess.Query{Name: "GetPointInfoQuery", SQL: `SELECT * FROM get_point_info($1::integer)`}

func (db *Database) GetPointInfoCommand(userId int) (info typepoints.PointInfo, err error) {
	row := dbaccess.QueryRow(getPointInfoQuery, userId)

	err = row.Scan(
		&info.TerritoryPoints,
		&info.FreePoints,
		&info.AvailableRolls,
		&info.TerritoryHours,
		&info.ExperiencePoints)

	dbaccess.LogDbResult(getPointInfoQuery, info, err)

	return
}

var getAllPointInfoQuery = dbaccess.Query{Name: "GetAllPointInfoQuery", SQL: `SELECT * FROM get_all_point_info()`}

func (db *Database) GetAllPointInfoCommand() (infos typepoints.PointInfoByLogins, err error) {
	rows, err := dbaccess.QueryRows(getAllPointInfoQuery)

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

	dbaccess.LogDbResult(getAllPointInfoQuery, infos, err)

	_ = rows.Close()
	return
}

var getFreePointsQuery = dbaccess.Query{Name: "GetFreePointsQuery", SQL: `SELECT * FROM get_free_points($1::integer)`}

func (db *Database) GetFreePointsCommand(userId int) (points int, err error) {
	row := dbaccess.QueryRow(getFreePointsQuery, userId)

	err = row.Scan(&points)

	dbaccess.LogDbResult(getFreePointsQuery, points, err)

	return
}

var changeFreePointsQuery = dbaccess.Query{Name: "ChangeFreePointsQuery", SQL: `SELECT change_free_points($1::integer, $2::integer)`}

func (db *Database) ChangeFreePointsCommand(userId int, changeValue int) error {
	_, err := dbaccess.Exec(changeFreePointsQuery, userId, changeValue)

	dbaccess.LogDbResult(changeFreePointsQuery, nil, err)

	return err
}

var addFreePointHistoryQuery = dbaccess.Query{Name: "AddFreePointHistoryQuery", SQL: `SELECT add_free_point_history($1::integer, $2::integer, $3::text, $4::integer, $5::integer, $6::integer, $7::integer)`}

func (db *Database) AddFreePointHistoryCommand(
	userId int,
	sourceUserId int,
	changeSource string,
	changeValue int,
	actualChangeValue int,
	finalValue int,
	wheelEffectId *int) error {

	_, err := dbaccess.Exec(
		addFreePointHistoryQuery,
		userId,
		sourceUserId,
		changeSource,
		changeValue,
		actualChangeValue,
		finalValue,
		wheelEffectId)

	dbaccess.LogDbResult(addFreePointHistoryQuery, nil, err)

	return err
}

var getFreePointHistoryQuery = dbaccess.Query{Name: "GetFreePointHistoryQuery", SQL: `SELECT * FROM get_free_point_history($1::integer)`}

func (db *Database) GetFreePointHistoryCommand(userId int) (history typepoints.FreePointChangeHistories, err error) {
	rows, err := dbaccess.QueryRows(getFreePointHistoryQuery, userId)

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

	dbaccess.LogDbResult(getFreePointHistoryQuery, history, err)

	_ = rows.Close()
	return
}
