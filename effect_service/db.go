package effect_service

import (
	"FGG-Service/common"
	"FGG-Service/db_access"
	"time"
)

const GetAvailableRollsCountQuery = `
	SELECT COUNT(*)
	FROM AvailableRolls
	WHERE UserId = ?
`

func GetAvailableRollsCountCommand(userId int) (count int, err error) {
	row := db_access.QueryRow(GetAvailableRollsCountQuery, userId)

	err = row.Scan(&count)

	return
}

const GetAvailableEffectsQuery = `
	SELECT e.Id, e.Name, e.Description
	FROM Effects e
		LEFT JOIN (
			SELECT EffectId
			FROM EffectHistory
			WHERE UserId = ?
		) eh ON e.Id = eh.EffectId
	WHERE eh.EffectId IS NULL
`

func GetAvailableEffectsCommand(userId int) (effects common.Effects, err error) {
	rows, err := db_access.Query(GetAvailableEffectsQuery, userId)

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
	SELECT e.Name, e.Description, eh.CreateDate
	FROM EffectHistory eh
		INNER JOIN Effects e ON eh.EffectId = e.Id
	WHERE eh.UserId = ?
`

func GetEffectHistoryCommand(userId int) (effects common.RolledEffects, err error) {
	rows, err := db_access.Query(GetEffectHistoryQuery, userId)

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
		rollDate, err = common.ConvertToDate(rollDateString)

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
	SELECT Id, Name, Description
	FROM Effects
	WHERE Id IN (
		SELECT id
		FROM Effects e
			LEFT JOIN (
				SELECT EffectId
				FROM EffectHistory
				WHERE UserId = ?
			) eh ON e.Id = eh.EffectId
		WHERE eh.EffectId IS NULL
		ORDER BY RANDOM()
		LIMIT 5)
`

func MakeEffectRollCommand(userId int) (effects common.Effects, err error) {
	rows, err := db_access.Query(MakeEffectRollQuery, userId)

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

const DeleteAvailableRollQuery = `
	DELETE FROM AvailableRolls
	WHERE Id IN (
		SELECT Id
		FROM AvailableRolls
		WHERE UserId = ?
		ORDER BY CreateDate
		LIMIT 1)
`

func DeleteAvailableRollCommand(userId int) error {
	_, err := db_access.Exec(DeleteAvailableRollQuery, userId)

	return err
}
