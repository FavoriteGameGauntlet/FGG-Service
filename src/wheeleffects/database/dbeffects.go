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

const GetAvailableRollsCountQuery = `SELECT * FROM get_available_rolls_count($1)`

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

const GetAvailableEffectsQuery = `SELECT * FROM get_available_effects($1)`

func (db *Database) GetAvailableEffectsCommand(userId int) (effects typewheeleffects.WheelEffects, err error) {
	queryName := "GetAvailableEffectsQuery"
	rows, err := dbaccess.Query(queryName, GetAvailableEffectsQuery, userId)

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

const GetEffectHistoryQuery = `SELECT * FROM get_effect_history($1)`

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

const GetEffectHistoryByEffectNameQuery = `SELECT * FROM get_effect_history_by_name($1, $2)`

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

const MakeEffectRollQuery = `SELECT * FROM make_effect_roll($1)`

func (db *Database) MakeEffectRollCommand(userId int) (effects typewheeleffects.WheelEffects, err error) {
	queryName := "MakeEffectRollQuery"
	rows, err := dbaccess.Query(queryName, MakeEffectRollQuery, userId)

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

const DecreaseAvailableRollsValueQuery = `SELECT decrease_available_rolls($1)`

func (db *Database) DecreaseAvailableRollsValueCommand(userId int) error {
	queryName := "DecreaseAvailableRollsValueQuery"
	_, err := dbaccess.Exec(queryName, DecreaseAvailableRollsValueQuery, userId)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const AddLastRolledWheelEffectsQuery = `SELECT add_last_rolled_wheel_effect($1, $2, $3)`

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

const GetLastRolledWheelEffectsQuery = `SELECT * FROM get_last_rolled_wheel_effects($1)`

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

const MarkLastWheelEffectAppliedQuery = `SELECT mark_last_wheel_effect_applied($1, $2)`

func (db *Database) MarkLastWheelEffectAppliedCommand(userId int, wheelEffectId int) error {
	queryName := "MarkLastWheelEffectAppliedQuery"
	_, err := dbaccess.Exec(queryName, MarkLastWheelEffectAppliedQuery, userId, wheelEffectId)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const AddWheelEffectHistoryQuery = `SELECT add_wheel_effect_history($1, $2)`

func (db *Database) AddWheelEffectHistoryCommand(userId int, wheelEffectId int) error {
	queryName := "AddWheelEffectHistoryQuery"
	_, err := dbaccess.Exec(queryName, AddWheelEffectHistoryQuery, userId, wheelEffectId)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}
