package timer_service

import (
	"FGG-Service/api"
	"FGG-Service/database"

	"github.com/google/uuid"
)

const (
	CheckIfCurrentTimerExistsCommand = `
		SELECT CASE 
			WHEN EXISTS (
				SELECT 1
				FROM Timers
				WHERE UserId = $userId)
         	THEN true
         	ELSE false
       		END AS 'DoesCurrentTimerExist'`
	GetCurrentTimerCommand = `
		SELECT 
			t.State, 
			t.DurationInS, 
			CASE ta.Action
				WHEN $startTimerAction THEN strftime('%s', 'now') - strftime('%s', ta.Date)
				WHEN $pauseTimerAction THEN ta.RemainingTimeInS
				WHEN $finishTimerAction THEN 0
				ELSE t.DurationInS
			END AS RemainingTimeInS
		FROM Timers t
			LEFT JOIN TimerActions ta ON t.Id = ta.TimerId
		WHERE UserId = $userId
		ORDER BY ta.Date DESC`
	CreateTimerCommand = `
		INSERT INTO Timers (Id, UserId, GameId, DurationInS)
		VALUES ($timerId, $userId, $gameId, $timerDurationInS)`
)

func CheckIfCurrentTimerExists(userId uuid.UUID) (*bool, error) {
	row := database.QueryRow(
		CheckIfCurrentTimerExistsCommand,
		userId,
		api.TimerStateFinished,
	)

	var doesCurrentTimerExist bool
	err := row.Scan(&doesCurrentTimerExist)

	if err != nil {
		return nil, err
	}

	return &doesCurrentTimerExist, nil
}

func GetOrCreateCurrentTimer(userId uuid.UUID, gameId uuid.UUID) (*api.Timer, error) {
	doesCurrentTimerExist, err := CheckIfCurrentTimerExists(userId)

	if err != nil {
		return nil, err
	}

	if *doesCurrentTimerExist {
		timer, err := GetCurrentTimer(userId)

		if err != nil {
			return nil, err
		}

		return timer, nil
	}

	timer, err := CreateCurrentTimer(userId, gameId)

	if err != nil {
		return nil, err
	}

	return timer, nil
}

func GetCurrentTimer(userId uuid.UUID) (*api.Timer, error) {
	row := database.QueryRow(
		GetCurrentTimerCommand,
		userId,
		api.TimerStateFinished,
	)

	timer := api.Timer{}
	err := row.Scan(&timer.DurationInS, &timer.State, &timer.RemainingTimeInS)

	if err != nil {
		return nil, err
	}

	return &timer, nil
}

func CreateCurrentTimer(userId uuid.UUID, gameId uuid.UUID) (*api.Timer, error) {
	timerId := uuid.New().String()
	_, err := database.Exec(
		CreateTimerCommand,
		timerId,
		userId,
		gameId,
		10800,
	)

	if err != nil {
		return nil, err
	}

	return &api.Timer{
		DurationInS:      10800,
		RemainingTimeInS: 10800,
		State:            api.TimerStateCreated,
	}, nil
}
