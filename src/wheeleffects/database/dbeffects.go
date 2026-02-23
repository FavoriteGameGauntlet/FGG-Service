package dbeffects

import (
	"FGG-Service/src/common"
	"FGG-Service/src/dbaccess"
	"time"
)

type Database struct {
}

const GetAvailableRollsCountQuery = `
	SELECT AvailableRolls
	FROM UserStats
	WHERE UserId = ?
`

func (db *Database) GetAvailableRollsCountCommand(userId int) (count int, err error) {
	row := dbaccess.QueryRow(GetAvailableRollsCountQuery, userId)

	err = row.Scan(&count)

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

func (db *Database) GetAvailableEffectsCommand(userId int) (effects common.Effects, err error) {
	rows, err := dbaccess.Query(GetAvailableEffectsQuery, userId)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		effect := common.Effect{}
		err = rows.Scan(&effect.Id, &effect.Name, &effect.Description)

		if err != nil {
			_ = rows.Close()
			return
		}

		effects = append(effects, effect)
	}

	_ = rows.Close()
	return
}

const GetEffectHistoryQuery = `
	SELECT we.Name, we.Description, weh.RollDate
	FROM WheelEffectHistory weh
		INNER JOIN WheelEffects we ON weh.WheelEffectId = we.Id
	WHERE weh.UserId = ?
`

func (db *Database) GetEffectHistoryCommand(userId int) (effects common.RolledEffects, err error) {
	rows, err := dbaccess.Query(GetEffectHistoryQuery, userId)

	if err != nil {
		return
	}

	for rows.Next() {
		effect := common.RolledEffect{}
		var rollDateString string
		err = rows.Scan(&effect.Name, &effect.Description, &rollDateString)

		if err != nil {
			_ = rows.Close()
			return
		}

		var rollDate time.Time
		rollDate, err = dbaccess.ConvertToDate(rollDateString)

		if err != nil {
			_ = rows.Close()
			return
		}

		effect.RollDate = rollDate

		effects = append(effects, effect)
	}

	_ = rows.Close()
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

func (db *Database) MakeEffectRollCommand(userId int) (effects common.Effects, err error) {
	rows, err := dbaccess.Query(MakeEffectRollQuery, userId)

	if err != nil {
		return
	}

	for rows.Next() {
		effect := common.Effect{}
		err = rows.Scan(&effect.Id, &effect.Name, &effect.Description)

		if err != nil {
			_ = rows.Close()
			return
		}

		effects = append(effects, effect)
	}

	_ = rows.Close()
	return
}

const DecreaseAvailableRollsValueQuery = `
	UPDATE UserStats
	SET AvailableRolls = AvailableRolls - 1
	WHERE UserId = ?
`

func (db *Database) DecreaseAvailableRollsValueCommand(userId int) error {
	_, err := dbaccess.Exec(DecreaseAvailableRollsValueQuery, userId)

	return err
}
