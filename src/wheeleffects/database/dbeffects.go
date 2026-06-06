package dbwheeleffects

import (
	"FGG-Service/src/dbaccess"
	"FGG-Service/src/wheeleffects/types"
	"database/sql"
	"errors"
	"time"
)

type IDatabase interface {
	GetAvailableRollsCountCommand(userId int) (count int, err error)
	GetAvailableEffectsCommand(userId int) (effects typewheeleffects.WheelEffects, err error)
	GetEffectHistoryCommand(userId int) (effects typewheeleffects.RolledWheelEffects, err error)
	MakeEffectRollCommand(userId int) (effects typewheeleffects.WheelEffects, err error)
	DecreaseAvailableRollsValueCommand(userId int) error
	GetLastRolledWheelEffectsCommand(userId int) (effects typewheeleffects.RolledWheelEffects, err error)
}

type Database struct {
}

const GetAvailableRollsCountQuery = `
	SELECT AvailableRolls
	FROM UserStats
	WHERE UserId = ?
`

func (db *Database) GetAvailableRollsCountCommand(userId int) (count int, err error) {
	queryName := "GetAvailableRollsCountQuery"
	row := dbaccess.QueryRow(queryName, GetAvailableRollsCountQuery, userId)

	err = row.Scan(&count)

	if errors.Is(err, sql.ErrNoRows) {
		count = 0
		err = nil
	}

	dbaccess.LogDbResult(queryName, count, err)

	return
}

const GetAvailableEffectsQuery = `
	SELECT we.Id, we.Name, we.Description
	FROM WheelEffects we
	WHERE NOT EXISTS (
		SELECT 1
		FROM WheelEffectHistory weh
		WHERE weh.WheelEffectId = we.Id
			AND weh.UserId = ?)
	  	AND NOT EXISTS (
			SELECT 1
			FROM LastWheelEffects lwe
			WHERE lwe.WheelEffectId = we.Id
				AND lwe.UserId = ?
				AND Position = 0)
`

func (db *Database) GetAvailableEffectsCommand(userId int) (effects typewheeleffects.WheelEffects, err error) {
	queryName := "GetAvailableEffectsQuery"
	rows, err := dbaccess.Query(queryName, GetAvailableEffectsQuery, userId, userId)

	if err != nil {
		return
	}

	for rows.Next() {
		effect := typewheeleffects.WheelEffect{}
		err = rows.Scan(&effect.Id, &effect.Name, &effect.Description)

		if err != nil {
			dbaccess.LogDbResult(queryName, effects, err)

			_ = rows.Close()
			return
		}

		effects = append(effects, effect)
	}

	dbaccess.LogDbResult(queryName, effects, err)

	_ = rows.Close()
	return
}

const GetEffectHistoryQuery = `
	SELECT we.Name, we.Description, weh.RollDate
	FROM WheelEffectHistory weh
		INNER JOIN WheelEffects we ON weh.WheelEffectId = we.Id
	WHERE weh.UserId = ?
`

func (db *Database) GetEffectHistoryCommand(userId int) (effects typewheeleffects.RolledWheelEffects, err error) {
	queryName := "GetEffectHistoryQuery"
	rows, err := dbaccess.Query(queryName, GetEffectHistoryQuery, userId)

	if err != nil {
		return
	}

	for rows.Next() {
		effect := typewheeleffects.RolledWheelEffect{}
		var rollDateString string
		err = rows.Scan(&effect.Name, &effect.Description, &rollDateString)

		if err != nil {
			dbaccess.LogDbResult(queryName, effects, err)

			_ = rows.Close()
			return
		}

		var rollDate time.Time
		rollDate, err = dbaccess.ConvertToDate(rollDateString)

		if err != nil {
			dbaccess.LogDbResult(queryName, effects, err)

			_ = rows.Close()
			return
		}

		effect.RollDate = rollDate

		effects = append(effects, effect)
	}

	dbaccess.LogDbResult(queryName, effects, err)

	_ = rows.Close()
	return
}

const GetEffectHistoryByEffectNameQuery = `
	SELECT we.Name, we.Description, weh.RollDate
	FROM WheelEffectHistory weh
		INNER JOIN WheelEffects we ON weh.WheelEffectId = we.Id
	WHERE weh.UserId = ?
		AND we.Name = ?
	ORDER BY weh.RollDate DESC
`

func (db *Database) GetEffectHistoryByEffectNameCommand(userId int, effectName string) (effect typewheeleffects.RolledWheelEffect, err error) {
	queryName := "GetEffectHistoryByEffectNameQuery"
	row := dbaccess.QueryRow(queryName, GetEffectHistoryByEffectNameQuery, userId, effectName)

	var rollDateString string
	err = row.Scan(&effect.Name, &effect.Description, &rollDateString)

	if err != nil {
		dbaccess.LogDbResult(queryName, effect, err)

		return
	}

	var rollDate time.Time
	rollDate, err = dbaccess.ConvertToDate(rollDateString)

	if err != nil {
		dbaccess.LogDbResult(queryName, effect, err)

		return
	}

	effect.RollDate = rollDate

	dbaccess.LogDbResult(queryName, effect, err)

	return
}

const MakeEffectRollQuery = `
	SELECT we.Id, we.Name, we.Description
	FROM WheelEffects we
	WHERE NOT EXISTS (
		SELECT 1
		FROM WheelEffectHistory weh
		WHERE weh.WheelEffectId = we.Id
			AND weh.UserId = ?)
		AND NOT EXISTS (
			SELECT 1
			FROM LastWheelEffects lwe
			WHERE lwe.WheelEffectId = we.Id
				AND lwe.UserId = ?
				AND Position = 0)
	ORDER BY RANDOM()
	LIMIT 5
`

func (db *Database) MakeEffectRollCommand(userId int) (effects typewheeleffects.WheelEffects, err error) {
	queryName := "MakeEffectRollQuery"
	rows, err := dbaccess.Query(queryName, MakeEffectRollQuery, userId, userId)

	if err != nil {
		return
	}

	for rows.Next() {
		effect := typewheeleffects.WheelEffect{}
		err = rows.Scan(&effect.Id, &effect.Name, &effect.Description)

		if err != nil {
			dbaccess.LogDbResult(queryName, effects, err)

			_ = rows.Close()
			return
		}

		effects = append(effects, effect)
	}

	dbaccess.LogDbResult(queryName, effects, err)

	_ = rows.Close()
	return
}

const DecreaseAvailableRollsValueQuery = `
	UPDATE UserStats
	SET AvailableRolls = AvailableRolls - 1
	WHERE UserId = ?
`

func (db *Database) DecreaseAvailableRollsValueCommand(userId int) error {
	queryName := "DecreaseAvailableRollsValueQuery"
	_, err := dbaccess.Exec(queryName, DecreaseAvailableRollsValueQuery, userId)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const GetLastRolledWheelEffectsQuery = `
	SELECT we.Id, we.Name, we.Description, lwe.RollDate, lwe.Position, lwe.IsApplied
	FROM LastWheelEffects lwe
		INNER JOIN WheelEffects we ON we.Id = lwe.WheelEffectId
	WHERE lwe.UserId = ?
`

func (db *Database) GetLastRolledWheelEffectsCommand(userId int) (effects typewheeleffects.RolledWheelEffects, err error) {
	queryName := "GetLastRolledWheelEffectsQuery"
	rows, err := dbaccess.Query(queryName, GetLastRolledWheelEffectsQuery, userId)

	if err != nil {
		return
	}

	for rows.Next() {
		effect := typewheeleffects.RolledWheelEffect{}

		var rollDateString string
		err = rows.Scan(
			&effect.Id,
			&effect.Name,
			&effect.Description,
			&rollDateString,
			&effect.Position,
			&effect.IsApplied)

		if err != nil {
			dbaccess.LogDbResult(queryName, effects, err)

			_ = rows.Close()
			return
		}

		var rollDate time.Time
		rollDate, err = dbaccess.ConvertToDate(rollDateString)

		if err != nil {
			dbaccess.LogDbResult(queryName, effects, err)

			_ = rows.Close()
			return
		}

		effect.RollDate = rollDate

		effects = append(effects, effect)
	}

	dbaccess.LogDbResult(queryName, effects, err)

	_ = rows.Close()
	return
}
