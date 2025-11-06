package timer_service

import (
	"FGG-Service/database"
	"database/sql"
	"errors"
	"time"

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
			ta.Date,
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

func CheckIfCurrentTimerExists(userId uuid.UUID) (bool, error) {
	row := database.QueryRow(
		CheckIfCurrentTimerExistsCommand,
		userId,
		TimerStateFinished,
	)

	var doesCurrentTimerExist bool
	err := row.Scan(&doesCurrentTimerExist)

	if err != nil {
		return doesCurrentTimerExist, err
	}

	return doesCurrentTimerExist, nil
}

func GetOrCreateCurrentTimer(userId uuid.UUID, gameId uuid.UUID) (*Timer, error) {
	doesCurrentTimerExist, err := CheckIfCurrentTimerExists(userId)

	if err != nil {
		return nil, err
	}

	if doesCurrentTimerExist {
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

func GetCurrentTimer(userId uuid.UUID) (*Timer, error) {
	row := database.QueryRow(
		GetCurrentTimerCommand,
		TimerActionStart,
		TimerActionPause,
		TimerActionStop,
		userId,
		TimerStateRolled,
	)

	timer := Timer{}
	var timerActionDateString string
	err := row.Scan(
		&timer.Id,
		&timer.State,
		&timer.DurationInS,
		&timerActionDateString,
		&timer.RemainingTimeInS,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	timerActionDate, err := time.Parse(database.ISO8601, timerActionDateString)

	if err != nil {
		return nil, err
	}

	timer.TimerActionDate = &timerActionDate

	return &timer, nil
}

func CreateCurrentTimer(userId uuid.UUID, gameId uuid.UUID) (*Timer, error) {
	timerId := uuid.New().String()
	_, err := database.Exec(
		CreateTimerCommand,
		timerId,
		userId,
		gameId,
		DefaultTimerDurationInS,
	)

	if err != nil {
		return nil, err
	}

	return &Timer{
		DurationInS:      DefaultTimerDurationInS,
		State:            TimerStateCreated,
		RemainingTimeInS: DefaultTimerDurationInS,
	}, nil
}

func StartCurrentTimer(userId uuid.UUID) (*TimerAction, error) {
	timer, err := GetCurrentTimer(userId)

	if err != nil {
		return nil, err
	}

	if timer == nil {
		return nil, nil
	}

	if timer.State == TimerStateRunning ||
		timer.State == TimerStateFinished {
		return nil, nil
	}

	timerActionId := uuid.New().String()
	timerAction := TimerActionStart
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

	return &TimerAction{
		Action:           timerAction,
		RemainingTimeInS: remainingTimerInS,
	}, nil
}

func PauseCurrentTimer(userId uuid.UUID) (*TimerAction, error) {
	timer, err := GetCurrentTimer(userId)

	if err != nil {
		return nil, err
	}

	if timer == nil {
		return nil, nil
	}

	if timer.State == TimerStateCreated ||
		timer.State == TimerStatePaused ||
		timer.State == TimerStateFinished {
		return nil, nil
	}

	timerActionId := uuid.New().String()
	timerAction := TimerActionPause
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

	return &TimerAction{
		Action:           timerAction,
		RemainingTimeInS: remainingTimerInS,
	}, nil
}

func StopCurrentTimer(userId uuid.UUID) (*TimerAction, error) {
	timer, err := GetCurrentTimer(userId)

	if err != nil {
		return nil, err
	}

	if timer == nil {
		return nil, nil
	}

	if timer.State == TimerStateCreated ||
		timer.State == TimerStateFinished {
		return nil, nil
	}

	timerActionId := uuid.New().String()
	timerAction := TimerActionStop
	remainingTimerInS := 0

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

	return &TimerAction{
		Action:           timerAction,
		RemainingTimeInS: remainingTimerInS,
	}, nil
}
