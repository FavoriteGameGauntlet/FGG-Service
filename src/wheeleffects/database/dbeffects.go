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

var getAvailableRollsCountQuery = dbaccess.Query{Name: "GetAvailableRollsCountQuery", SQL: `SELECT * FROM get_available_rolls_count($1::integer)`}

func (db *Database) GetAvailableRollsCountCommand(userId int) (count int, err error) {
	row := dbaccess.QueryRow(getAvailableRollsCountQuery, userId)

	err = row.Scan(&count)

	if errors.Is(err, sql.ErrNoRows) {
		count = 0
		err = nil
	}

	dbaccess.LogDbResult(getAvailableRollsCountQuery, count, err)

	return
}

var getAvailableEffectsQuery = dbaccess.Query{Name: "GetAvailableEffectsQuery", SQL: `SELECT * FROM get_available_effects($1::integer)`}

func (db *Database) GetAvailableEffectsCommand(userId int) (effects typewheeleffects.WheelEffects, err error) {
	rows, err := dbaccess.QueryRows(getAvailableEffectsQuery, userId)

	if err != nil {
		return
	}

	for rows.Next() {
		effect := typewheeleffects.WheelEffect{}
		err = rows.Scan(&effect.Id, &effect.Name, &effect.Description)

		if err != nil {
			dbaccess.LogDbResult(getAvailableEffectsQuery, effects, err)

			_ = rows.Close()
			return
		}

		effects = append(effects, effect)
	}

	dbaccess.LogDbResult(getAvailableEffectsQuery, effects, err)

	_ = rows.Close()
	return
}

var getEffectHistoryQuery = dbaccess.Query{Name: "GetEffectHistoryQuery", SQL: `SELECT * FROM get_effect_history($1::integer)`}

func (db *Database) GetEffectHistoryCommand(userId int) (effects typewheeleffects.RolledWheelEffectHistories, err error) {
	rows, err := dbaccess.QueryRows(getEffectHistoryQuery, userId)

	if err != nil {
		return
	}

	for rows.Next() {
		effect := typewheeleffects.RolledWheelEffectHistory{}
		err = rows.Scan(&effect.Name, &effect.Description, &effect.RollDate)

		if err != nil {
			dbaccess.LogDbResult(getEffectHistoryQuery, effects, err)

			_ = rows.Close()
			return
		}

		effects = append(effects, effect)
	}

	dbaccess.LogDbResult(getEffectHistoryQuery, effects, err)

	_ = rows.Close()
	return
}

var getEffectHistoryByEffectNameQuery = dbaccess.Query{Name: "GetEffectHistoryByEffectNameQuery", SQL: `SELECT * FROM get_effect_history_by_name($1::integer, $2::text)`}

func (db *Database) GetEffectHistoryByEffectNameCommand(userId int, effectName string) (effect typewheeleffects.RolledWheelEffect, err error) {
	row := dbaccess.QueryRow(getEffectHistoryByEffectNameQuery, userId, effectName)

	err = row.Scan(&effect.Name, &effect.Description, &effect.RollDate)

	dbaccess.LogDbResult(getEffectHistoryByEffectNameQuery, effect, err)

	return
}

var makeEffectRollQuery = dbaccess.Query{Name: "MakeEffectRollQuery", SQL: `SELECT * FROM make_effect_roll($1::integer)`}

func (db *Database) MakeEffectRollCommand(userId int) (effects typewheeleffects.WheelEffects, err error) {
	rows, err := dbaccess.QueryRows(makeEffectRollQuery, userId)

	if err != nil {
		return
	}

	for rows.Next() {
		effect := typewheeleffects.WheelEffect{}
		err = rows.Scan(&effect.Id, &effect.Name, &effect.Description)

		if err != nil {
			dbaccess.LogDbResult(makeEffectRollQuery, effects, err)

			_ = rows.Close()
			return
		}

		effects = append(effects, effect)
	}

	dbaccess.LogDbResult(makeEffectRollQuery, effects, err)

	_ = rows.Close()
	return
}

var decreaseAvailableRollsValueQuery = dbaccess.Query{Name: "DecreaseAvailableRollsValueQuery", SQL: `SELECT decrease_available_rolls($1::integer)`}

func (db *Database) DecreaseAvailableRollsValueCommand(userId int) error {
	_, err := dbaccess.Exec(decreaseAvailableRollsValueQuery, userId)

	dbaccess.LogDbResult(decreaseAvailableRollsValueQuery, nil, err)

	return err
}

var addLastRolledWheelEffectsQuery = dbaccess.Query{Name: "AddLastRolledWheelEffectsQuery", SQL: `SELECT add_last_rolled_wheel_effect($1::integer, $2::integer, $3::integer)`}

func (db *Database) AddLastRolledWheelEffectsCommand(userId int, effects typewheeleffects.WheelEffects) (err error) {
	for i, effect := range effects {
		_, err = dbaccess.Exec(addLastRolledWheelEffectsQuery, userId, effect.Id, i-2)

		if err != nil {
			dbaccess.LogDbResult(addLastRolledWheelEffectsQuery, nil, err)

			return
		}
	}

	dbaccess.LogDbResult(addLastRolledWheelEffectsQuery, effects, err)

	return
}

var getLastRolledWheelEffectsQuery = dbaccess.Query{Name: "GetLastRolledWheelEffectsQuery", SQL: `SELECT * FROM get_last_rolled_wheel_effects($1::integer)`}

func (db *Database) GetLastRolledWheelEffectsCommand(userId int) (effects typewheeleffects.RolledWheelEffects, err error) {
	rows, err := dbaccess.QueryRows(getLastRolledWheelEffectsQuery, userId)

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
			dbaccess.LogDbResult(getLastRolledWheelEffectsQuery, effects, err)

			_ = rows.Close()
			return
		}

		effects = append(effects, effect)
	}

	dbaccess.LogDbResult(getLastRolledWheelEffectsQuery, effects, err)

	_ = rows.Close()
	return
}

var markLastWheelEffectAppliedQuery = dbaccess.Query{Name: "MarkLastWheelEffectAppliedQuery", SQL: `SELECT mark_last_wheel_effect_applied($1::integer, $2::integer)`}

func (db *Database) MarkLastWheelEffectAppliedCommand(userId int, wheelEffectId int) error {
	_, err := dbaccess.Exec(markLastWheelEffectAppliedQuery, userId, wheelEffectId)

	dbaccess.LogDbResult(markLastWheelEffectAppliedQuery, nil, err)

	return err
}

var addWheelEffectHistoryQuery = dbaccess.Query{Name: "AddWheelEffectHistoryQuery", SQL: `SELECT add_wheel_effect_history($1::integer, $2::integer)`}

func (db *Database) AddWheelEffectHistoryCommand(userId int, wheelEffectId int) error {
	_, err := dbaccess.Exec(addWheelEffectHistoryQuery, userId, wheelEffectId)

	dbaccess.LogDbResult(addWheelEffectHistoryQuery, nil, err)

	return err
}
