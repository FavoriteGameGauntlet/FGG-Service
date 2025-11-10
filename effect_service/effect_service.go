package effect_service

import (
	"FGG-Service/db_access"
	"time"

	"github.com/google/uuid"
)

const (
	CheckIfEffectRollExistsCommand = `
		SELECT CASE 
        	WHEN EXISTS (
				SELECT 1
				FROM EffectHistory
				WHERE UserId = $userId
					AND RollDate IS NULL)
         	THEN true
         	ELSE false
       	END AS 'DoesExist'`
	GetEffectHistoryCommand = `
		SELECT e.Name, e.Description, eh.CreateDate, eh.RollDate, g.Name AS GameName 
		FROM EffectHistory eh
			LEFT JOIN Effects e ON eh.RolledEffectId = e.Id
			LEFT JOIN Games g ON eh.GameId = g.Id
		WHERE eh.UserId = $userId
		ORDER BY eh.CreateDate, eh.RollDate`
)

func CheckIfEffectRollExists(userId uuid.UUID) (bool, error) {
	row := db_access.QueryRow(CheckIfEffectRollExistsCommand, userId)

	var doesExist bool
	err := row.Scan(&doesExist)

	if err != nil {
		return doesExist, err
	}

	return doesExist, nil
}

func GetEffectHistory(userId uuid.UUID) (*Effects, error) {
	rows, err := db_access.Query(GetEffectHistoryCommand, userId)

	if err != nil {
		return nil, err
	}

	effectCount := 0
	errorCount := 0
	effects := Effects{}
	for rows.Next() {
		effectCount++

		effect := Effect{}
		var createDateString string
		var rollDateString *string

		err = rows.Scan(&effect.Name, &effect.Description, &createDateString, &rollDateString, &effect.GameName)

		if err != nil {
			errorCount++
			continue
		}

		var createDate time.Time
		createDate, err = time.Parse(db_access.ISO8601, createDateString)

		if err != nil {
			errorCount++
			continue
		}

		var rollDate *time.Time

		if rollDateString != nil {
			var notNilRollDate time.Time
			notNilRollDate, err = time.Parse(db_access.ISO8601, *rollDateString)

			if err != nil {
				errorCount++
				continue
			}

			rollDate = &notNilRollDate
		}

		effect.CreateDate = createDate
		effect.RollDate = rollDate

		effects = append(effects, effect)
	}

	if errorCount > 0 && errorCount == effectCount {
		return nil, err
	}

	return &effects, nil
}
