package timer_service

import (
	"FGG-Service/db_access"
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
					AND State != $finishedTimerState)
         	THEN true
         	ELSE false
       		END AS 'DoesExist'`
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
			AND t.State != $finishedTimerState
		ORDER BY ta.Date
		LIMIT 1`
	CreateTimerCommand = `
		INSERT INTO Timers (Id, UserId, GameId, DurationInS)
		VALUES ($timerId, $userId, $gameId, $timerDurationInS)`
	ActCurrentTimerCommand = `
		INSERT INTO TimerActions (Id, TimerId, Action, RemainingTimeInS)
		VALUES ($timerActionId, $timerId, $timerAction, $remainingTimeInS)`
	GetCompletedTimerUsersCommand = `
		SELECT t.UserId
		FROM Timers t
			LEFT JOIN TimerActions ta ON t.Id = ta.TimerId
		WHERE t.State IN ($runningTimerState, $pausedTimerState)
	  	    AND CASE t.State
				WHEN $runningTimerState THEN ta.RemainingTimeInS - (strftime('%s', 'now') - strftime('%s', ta.Date))
				WHEN $pausedTimerState THEN ta.RemainingTimeInS
			END <= 0`
)

func CheckIfCurrentTimerExists(userId uuid.UUID) (bool, error) {
	row := db_access.QueryRow(
		CheckIfCurrentTimerExistsCommand,
		userId,
		TimerStateFinished,
	)

	var doesExist bool
	err := row.Scan(&doesExist)

	if err != nil {
		return doesExist, err
	}

	return doesExist, nil
}

func GetOrCreateCurrentTimer(userId uuid.UUID, gameId uuid.UUID) (*Timer, error) {
	doesExist, err := CheckIfCurrentTimerExists(userId)

	if err != nil {
		return nil, err
	}

	if doesExist {
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
	row := db_access.QueryRow(
		GetCurrentTimerCommand,
		TimerActionStart,
		TimerActionPause,
		TimerActionStop,
		userId,
		TimerStateFinished,
	)

	timer := Timer{}
	var timerActionDateString *string
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

	var timerActionDate *time.Time

	if timerActionDateString != nil {
		var notNilDate time.Time
		notNilDate, err = time.Parse(db_access.ISO8601, *timerActionDateString)

		if err != nil {
			return nil, err
		}

		timerActionDate = &notNilDate
	}

	timer.TimerActionDate = timerActionDate

	return &timer, nil
}

func CreateCurrentTimer(userId uuid.UUID, gameId uuid.UUID) (*Timer, error) {
	timerId := uuid.New().String()
	_, err := db_access.Exec(
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

	_, err = db_access.Exec(
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

	_, err = db_access.Exec(
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

	_, err = db_access.Exec(
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

func StopAllCompletedTimers() error {
	rows, err := db_access.Query(GetCompletedTimerUsersCommand, TimerStateRunning, TimerStatePaused)

	if err != nil {
		return err
	}

	timerCount := 0
	errorCount := 0
	for rows.Next() {
		timerCount++

		var userId uuid.UUID
		err = rows.Scan(&userId)

		if err != nil {
			errorCount++
			continue
		}

		_, err = StopCurrentTimer(userId)

		if err != nil {
			errorCount++
			continue
		}
	}

	if errorCount > 0 && errorCount == timerCount {
		return err
	}

	return nil
}
