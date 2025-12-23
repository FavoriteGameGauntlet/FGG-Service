package effect_service

import (
	"FGG-Service/common"
	"FGG-Service/db_access"
	"time"
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
		SELECT e.Id, e.Name, e.Description
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
	CreateAvailableRollCommand = `
		INSERT INTO AvailableRolls (UserId)
		VALUES ($userId)`
	DeleteAvailableRollCommand = `
		DELETE FROM AvailableRolls
		WHERE Id IN (
		  SELECT Id
		  FROM AvailableRolls
		  WHERE UserId = $userId
		  ORDER BY CreateDate
		  LIMIT 1)`
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

func GetAvailableEffects(userId int) (*common.Effects, error) {
	rows, err := db_access.Query(GetAvailableEffectsCommand, userId)

	if err != nil {
		return nil, err
	}

	effectCount := 0
	errorCount := 0
	effects := common.Effects{}
	for rows.Next() {
		effectCount++

		effect := common.Effect{}
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

	return &effects, nil
}

func GetEffectHistory(userId int) (*common.RolledEffects, error) {
	rows, err := db_access.Query(GetEffectHistoryCommand, userId)

	if err != nil {
		return nil, err
	}

	effectCount := 0
	errorCount := 0
	effects := common.RolledEffects{}
	for rows.Next() {
		effectCount++

		effect := common.RolledEffect{}
		var rollDateString string
		err = rows.Scan(&effect.Name, &effect.Description, &rollDateString)

		if err != nil {
			errorCount++
			continue
		}

		var rollDate time.Time
		rollDate, err = time.Parse(db_access.ISO8601, rollDateString)

		if err != nil {
			errorCount++
			continue
		}

		effect.RollDate = rollDate

		effects = append(effects, effect)
	}

	if errorCount > 0 && errorCount == effectCount {
		return nil, err
	}

	return &effects, nil
}

func MakeEffectRoll(userId int) (*common.Effects, error) {
	doesExist, err := CheckIfAvailableRollExists(userId)

	if err != nil {
		return nil, err
	}

	if !doesExist {
		return nil, common.NewAvailableRollsNotFoundError()
	}

	rows, err := db_access.Query(MakeEffectRollCommand, userId)

	if err != nil {
		return nil, err
	}

	err = DeleteAvailableRoll(userId)

	if err != nil {
		return nil, err
	}

	effectCount := 0
	errorCount := 0
	effects := common.Effects{}
	for rows.Next() {
		effectCount++

		effect := common.Effect{}
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

func CreateAvailableRoll(userId int) error {
	_, err := db_access.Exec(CreateAvailableRollCommand, userId)

	return err
}

func DeleteAvailableRoll(userId int) error {
	_, err := db_access.Exec(DeleteAvailableRollCommand, userId)

	return err
}
