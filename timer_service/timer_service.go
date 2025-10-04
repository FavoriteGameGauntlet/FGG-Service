package timer_service

import (
	"FGG-Service/api"
	"FGG-Service/database"
	"FGG-Service/game_service"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

const (
	CheckIfCurrentTimerExistsCommand = `
		SELECT CASE 
			WHEN EXISTS (
				SELECT 1
				FROM Timers
				WHERE UserId = $userId
					AND State != $rolledTimerState)
         	THEN true
         	ELSE false
       		END AS 'DoesCurrentTimerExist'`
	GetCurrentTimerCommand = `
		SELECT
			t.Id,
			t.State, 
			t.DurationInS, 
			CASE ta.Action
				WHEN $startTimerAction THEN ta.RemainingTimeInS - (strftime('%s', 'now') - strftime('%s', ta.Date))
				WHEN $pauseTimerAction THEN ta.RemainingTimeInS
				WHEN $finishTimerAction THEN 0
				ELSE t.DurationInS
			END AS RemainingTimeInS
		FROM Timers t
			LEFT JOIN TimerActions ta ON t.Id = ta.TimerId
		WHERE UserId = $userId
			AND State != $rolledTimerState
		ORDER BY ta.Date DESC`
	CreateTimerCommand = `
		INSERT INTO Timers (Id, UserId, GameId, DurationInS)
		VALUES ($timerId, $userId, $gameId, $timerDurationInS)`
	ActCurrentTimerCommand = `
		INSERT INTO TimerActions (Id, TimerId, Action, RemainingTimeInS)
		VALUES ($timerActionId, $timerId, $timerAction, $remainingTimeInS)`
)

const DefaultTimerDurationInS = 10800

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

func GetOrCreateCurrentTimer(userId uuid.UUID) (*api.Timer, error) {
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

	timer, err := CreateCurrentTimer(userId)

	if err != nil {
		return nil, err
	}

	return timer, nil
}

func GetCurrentTimer(userId uuid.UUID) (*api.Timer, error) {
	row := database.QueryRow(
		GetCurrentTimerCommand,
		api.Start,
		api.Pause,
		api.Stop,
		userId,
		api.TimerStateRolled,
	)

	timer := api.Timer{}
	err := row.Scan(
		&timer.Id,
		&timer.State,
		&timer.DurationInS,
		&timer.RemainingTimeInS,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &timer, nil
}

func CreateCurrentTimer(userId uuid.UUID) (*api.Timer, error) {
	game, err := game_service.GetCurrentGame(userId)

	if err != nil {
		return nil, err
	}

	timerId := uuid.New().String()
	_, err = database.Exec(
		CreateTimerCommand,
		timerId,
		userId,
		game.Id,
		DefaultTimerDurationInS,
	)

	if err != nil {
		return nil, err
	}

	return &api.Timer{
		DurationInS:      DefaultTimerDurationInS,
		State:            api.TimerStateCreated,
		RemainingTimeInS: DefaultTimerDurationInS,
	}, nil
}

func StartCurrentTimer(userId uuid.UUID) (*api.TimerAction, error) {
	timer, err := GetCurrentTimer(userId)

	if err != nil {
		return nil, err
	}

	if timer.State == api.TimerStateRunning ||
		timer.State == api.TimerStateFinished {
		return nil, nil
	}

	timerActionId := uuid.New().String()
	timerAction := api.Start
	remainingTimerInS := timer.RemainingTimeInS

	_, err = database.Exec(
		ActCurrentTimerCommand,
		timerActionId,
		timer.Id,
		timerAction,
		remainingTimerInS,
	)

	if err != nil {
		return nil, err
	}

	return &api.TimerAction{
		Action:           timerAction,
		RemainingTimeInS: remainingTimerInS,
	}, nil
}

func PauseCurrentTimer(userId uuid.UUID) (*api.TimerAction, error) {
	timer, err := GetCurrentTimer(userId)

	if err != nil {
		return nil, err
	}

	if timer.State == api.TimerStateCreated ||
		timer.State == api.TimerStatePaused ||
		timer.State == api.TimerStateFinished {
		return nil, nil
	}

	timerActionId := uuid.New().String()
	timerAction := api.Pause
	remainingTimerInS := timer.RemainingTimeInS

	_, err = database.Exec(
		ActCurrentTimerCommand,
		timerActionId,
		timer.Id,
		timerAction,
		remainingTimerInS,
	)

	if err != nil {
		return nil, err
	}

	return &api.TimerAction{
		Action:           timerAction,
		RemainingTimeInS: remainingTimerInS,
	}, nil
}

func StopCurrentTimer(userId uuid.UUID) (*api.TimerAction, error) {
	timer, err := GetCurrentTimer(userId)

	if err != nil {
		return nil, err
	}

	if timer.State == api.TimerStateCreated ||
		timer.State == api.TimerStateFinished {
		return nil, nil
	}

	timerActionId := uuid.New().String()
	timerAction := api.Stop
	var remainingTimerInS int32 = 0

	_, err = database.Exec(
		ActCurrentTimerCommand,
		timerActionId,
		timer.Id,
		timerAction,
		remainingTimerInS,
	)

	if err != nil {
		return nil, err
	}

	return &api.TimerAction{
		Action:           timerAction,
		RemainingTimeInS: remainingTimerInS,
	}, nil
}
