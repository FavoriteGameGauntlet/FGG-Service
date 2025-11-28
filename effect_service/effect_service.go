package effect_service

import (
	"FGG-Service/db_access"
)

const (
	CheckIfAvailableRollExistsCommand = `
		SELECT CASE 
        	WHEN EXISTS (
				SELECT 1
				FROM AvailableRolls
				WHERE UserId = $userId)
         	THEN true
         	ELSE false
       	END AS 'DoesExist'`
	GetAvailableEffectsCommand = `
		SELECT e.Name, e.Description
		FROM Effects e
			LEFT JOIN (
				SELECT EffectId
				FROM EffectHistory
				WHERE UserId = $userId
			) eh ON e.Id = eh.EffectId
		WHERE eh.EffectId IS NULL`
	GetEffectHistoryCommand = `
		SELECT e.Name, e.Description, eh.CreateDate
		FROM EffectHistory eh
			INNER JOIN Effects e ON eh.EffectId = e.Id
		WHERE eh.UserId = $userId`
	MakeEffectRollCommand = `
		SELECT Id, Name, Description
		FROM Effects
		WHERE Id IN (
			SELECT id
			FROM Effects e
		  		LEFT JOIN (
					SELECT EffectId
					FROM EffectHistory
					WHERE UserId = $userId
				) eh ON e.Id = eh.EffectId
			WHERE eh.EffectId IS NULL
			ORDER BY RANDOM()
			LIMIT 5)`
	AddEffectHistoryCommand = `
		INSERT INTO EffectHistory (EffectId, UserId)
		VALUES ($effectId, $userId)`
)

func CheckIfAvailableRollExists(userId int) (bool, error) {
	row := db_access.QueryRow(CheckIfAvailableRollExistsCommand, userId)

	var doesExist bool
	err := row.Scan(&doesExist)

	if err != nil {
		return doesExist, err
	}

	return doesExist, nil
}

func GetAvailableEffects(userId int) (*Effects, error) {
	rows, err := db_access.Query(GetAvailableEffectsCommand, userId)

	if err != nil {
		return nil, err
	}

	effectCount := 0
	errorCount := 0
	effects := Effects{}
	for rows.Next() {
		effectCount++

		effect := Effect{}
		err = rows.Scan(&effect.Name, &effect.Description)

		if err != nil {
			errorCount++
			continue
		}

		effects = append(effects, effect)
	}

	if errorCount > 0 && errorCount == effectCount {
		return nil, err
	}

	return &effects, nil
}

func GetEffectHistory(userId int) (*RolledEffects, error) {
	rows, err := db_access.Query(GetEffectHistoryCommand, userId)

	if err != nil {
		return nil, err
	}

	effectCount := 0
	errorCount := 0
	effects := RolledEffects{}
	for rows.Next() {
		effectCount++

		effect := RolledEffect{}
		err = rows.Scan(&effect.Name, &effect.Description)

		if err != nil {
			errorCount++
			continue
		}

		effects = append(effects, effect)
	}

	if errorCount > 0 && errorCount == effectCount {
		return nil, err
	}

	return &effects, nil
}

func MakeEffectRoll(userId int) (*Effects, error) {
	rows, err := db_access.Query(MakeEffectRollCommand, userId)

	if err != nil {
		return nil, err
	}

	effectCount := 0
	errorCount := 0
	effects := Effects{}
	for rows.Next() {
		effectCount++

		effect := Effect{}
		err = rows.Scan(&effect.Id, &effect.Name, &effect.Description)

		if err != nil {
			errorCount++
			continue
		}

		effects = append(effects, effect)
	}

	if errorCount > 0 && errorCount == effectCount {
		return nil, err
	}

	chosenEffectId := effects[2].Id

	_, err = db_access.Exec(AddEffectHistoryCommand, chosenEffectId, userId)

	if err != nil {
		return nil, err
	}

	return &effects, nil
}
