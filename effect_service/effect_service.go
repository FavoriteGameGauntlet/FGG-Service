package effect_service

import (
	"FGG-Service/database"
	"time"

	"github.com/google/uuid"
)

const (
	GetEffectHistoryCommand = `
		SELECT e.Name, e.Description, eh.CreateDate, eh.RollDate, g.Name AS GameName 
		FROM EffectHistory eh
			LEFT JOIN Effects e ON eh.RolledEffectId = e.Id
			LEFT JOIN Games g ON eh.GameId = g.Id
		WHERE eh.UserId = $userId
		ORDER BY eh.CreateDate, eh.RollDate`
)

func GetEffectHistory(userId uuid.UUID) (*Effects, error) {
	rows, err := database.Query(GetEffectHistoryCommand, userId)

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
		createDate, err = time.Parse(database.ISO8601, createDateString)

		if err != nil {
			errorCount++
			continue
		}

		var rollDate *time.Time

		if rollDateString != nil {
			var notNilRollDate time.Time
			notNilRollDate, err = time.Parse(database.ISO8601, *rollDateString)

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
