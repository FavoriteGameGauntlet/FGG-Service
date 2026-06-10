package dbwheeleffects

import (
	"FGG-Service/src/dbaccess"
	"FGG-Service/src/wheeleffects/types"
	"database/sql"
	"errors"
)

type IDatabase interface {
	GetAvailableRollsCountCommand(userId int) (count int, err error)
	GetAvailableEffectsCommand(userId int) (effects typewheeleffects.WheelEffects, err error)
	GetEffectHistoryCommand(userId int) (effects typewheeleffects.RolledWheelEffectHistories, err error)
	GetEffectHistoryByEffectNameCommand(userId int, effectName string) (effect typewheeleffects.RolledWheelEffect, err error)
	MakeEffectRollCommand(userId int) (effects typewheeleffects.WheelEffects, err error)
	DecreaseAvailableRollsValueCommand(userId int) error
	AddLastRolledWheelEffectsCommand(userId int, effects typewheeleffects.WheelEffects) (err error)
	GetLastRolledWheelEffectsCommand(userId int) (effects typewheeleffects.RolledWheelEffects, err error)
	MarkLastWheelEffectAppliedCommand(userId int, wheelEffectId int) error
	AddWheelEffectHistoryCommand(userId int, wheelEffectId int) error
}

type Database struct {
}

const GetAvailableRollsCountQuery = `
	SELECT AvailableRolls
	FROM UserStats
	WHERE UserId = $1
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
			AND weh.UserId = $1)
	  	AND NOT EXISTS (
			SELECT 1
			FROM LastWheelEffects lwe
			WHERE lwe.WheelEffectId = we.Id
				AND lwe.UserId = $2
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
	WHERE weh.UserId = $1
`

func (db *Database) GetEffectHistoryCommand(userId int) (effects typewheeleffects.RolledWheelEffectHistories, err error) {
	queryName := "GetEffectHistoryQuery"
	rows, err := dbaccess.Query(queryName, GetEffectHistoryQuery, userId)

	if err != nil {
		return
	}

	for rows.Next() {
		effect := typewheeleffects.RolledWheelEffectHistory{}
		err = rows.Scan(&effect.Name, &effect.Description, &effect.RollDate)

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

const GetEffectHistoryByEffectNameQuery = `
	SELECT we.Name, we.Description, weh.RollDate
	FROM WheelEffectHistory weh
		INNER JOIN WheelEffects we ON weh.WheelEffectId = we.Id
	WHERE weh.UserId = $1
		AND we.Name = $2
	ORDER BY weh.RollDate DESC
`

func (db *Database) GetEffectHistoryByEffectNameCommand(userId int, effectName string) (effect typewheeleffects.RolledWheelEffect, err error) {
	queryName := "GetEffectHistoryByEffectNameQuery"
	row := dbaccess.QueryRow(queryName, GetEffectHistoryByEffectNameQuery, userId, effectName)

	err = row.Scan(&effect.Name, &effect.Description, &effect.RollDate)

	if err != nil {
		dbaccess.LogDbResult(queryName, effect, err)

		return
	}

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
			AND weh.UserId = $1)
		AND NOT EXISTS (
			SELECT 1
			FROM LastWheelEffects lwe
			WHERE lwe.WheelEffectId = we.Id
				AND lwe.UserId = $2
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
	WHERE UserId = $1
`

func (db *Database) DecreaseAvailableRollsValueCommand(userId int) error {
	queryName := "DecreaseAvailableRollsValueQuery"
	_, err := dbaccess.Exec(queryName, DecreaseAvailableRollsValueQuery, userId)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const AddLastRolledWheelEffectsQuery = `
	INSERT INTO LastWheelEffects (UserId, WheelEffectId, Position)
	VALUES ($1, $2, $3)
`

func (db *Database) AddLastRolledWheelEffectsCommand(userId int, effects typewheeleffects.WheelEffects) (err error) {
	queryName := "AddLastRolledWheelEffectsQuery"

	for i, effect := range effects {
		_, err = dbaccess.Exec(queryName, AddLastRolledWheelEffectsQuery, userId, effect.Id, i-2)

		if err != nil {
			dbaccess.LogDbResult(queryName, nil, err)

			return
		}
	}

	dbaccess.LogDbResult(queryName, effects, err)

	return
}

const GetLastRolledWheelEffectsQuery = `
	SELECT we.Id, we.Name, we.Description, lwe.RollDate, lwe.Position, lwe.IsApplied
	FROM LastWheelEffects lwe
		INNER JOIN WheelEffects we ON we.Id = lwe.WheelEffectId
	WHERE lwe.UserId = $1
`

func (db *Database) GetLastRolledWheelEffectsCommand(userId int) (effects typewheeleffects.RolledWheelEffects, err error) {
	queryName := "GetLastRolledWheelEffectsQuery"
	rows, err := dbaccess.Query(queryName, GetLastRolledWheelEffectsQuery, userId)

	if err != nil {
		return
	}

	for rows.Next() {
		effect := typewheeleffects.RolledWheelEffect{}

		err = rows.Scan(
			&effect.Id,
			&effect.Name,
			&effect.Description,
			&effect.RollDate,
			&effect.Position,
			&effect.IsApplied)

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

const MarkLastWheelEffectAppliedQuery = `
	UPDATE LastWheelEffects
	SET IsApplied = 1
	WHERE UserId = $1
		AND WheelEffectId = $2
`

func (db *Database) MarkLastWheelEffectAppliedCommand(userId int, wheelEffectId int) error {
	queryName := "MarkLastWheelEffectAppliedQuery"
	_, err := dbaccess.Exec(queryName, MarkLastWheelEffectAppliedQuery, userId, wheelEffectId)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const AddWheelEffectHistoryQuery = `
	INSERT INTO WheelEffectHistory (UserId, WheelEffectId)
	VALUES ($1, $2)
`

func (db *Database) AddWheelEffectHistoryCommand(userId int, wheelEffectId int) error {
	queryName := "AddWheelEffectHistoryQuery"
	_, err := dbaccess.Exec(queryName, AddWheelEffectHistoryQuery, userId, wheelEffectId)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}
