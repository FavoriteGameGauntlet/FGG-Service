package timer_service

import (
	"FGG-Service/db_access"
	"database/sql"
	"errors"
	"time"
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
			ta.CreateDate,
			CASE ta.Action
				WHEN $startTimerAction THEN ta.RemainingTimeInS - (strftime('%s', 'now') - strftime('%s', ta.CreateDate))
				WHEN $pauseTimerAction THEN ta.RemainingTimeInS
				WHEN $finishTimerAction THEN 0
				ELSE t.DurationInS
			END AS RemainingTimeInS
		FROM Timers t
			LEFT JOIN TimerActions ta ON t.Id = ta.TimerId
		WHERE UserId = $userId
			AND t.State != $finishedTimerState
		ORDER BY ta.CreateDate DESC
		LIMIT 1`
	CreateTimerCommand = `
		INSERT INTO Timers (UserId, GameId, DurationInS)
		VALUES ($userId, $gameId, $timerDurationInS)`
	ActCurrentTimerCommand = `
		INSERT INTO TimerActions (TimerId, Action, RemainingTimeInS)
		VALUES ($timerId, $timerAction, $remainingTimeInS)`
	GetCompletedTimerUsersCommand = `
		SELECT DISTINCT t.UserId
		FROM Timers t
			LEFT JOIN TimerActions ta ON t.Id = ta.TimerId
		WHERE t.State IN ($runningTimerState, $pausedTimerState)
	  	    AND CASE t.State
				WHEN $runningTimerState THEN ta.RemainingTimeInS - (strftime('%s', 'now') - strftime('%s', ta.CreateDate))
				WHEN $pausedTimerState THEN ta.RemainingTimeInS
			END <= 0`
)

func CheckIfCurrentTimerExists(userId int) (bool, error) {
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

func GetOrCreateCurrentTimer(userId int, gameId int) (*Timer, error) {
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

func GetCurrentTimer(userId int) (*Timer, error) {
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

func CreateCurrentTimer(userId int, gameId int) (*Timer, error) {
	_, err := db_access.Exec(
		CreateTimerCommand,
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

func StartCurrentTimer(userId int) (*TimerAction, error) {
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

	timerAction := TimerActionStart
	remainingTimerInS := timer.RemainingTimeInS

	_, err = db_access.Exec(
		ActCurrentTimerCommand,
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

func PauseCurrentTimer(userId int) (*TimerAction, error) {
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

	timerAction := TimerActionPause
	remainingTimerInS := timer.RemainingTimeInS

	_, err = db_access.Exec(
		ActCurrentTimerCommand,
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

func StopCurrentTimer(userId int) (*TimerAction, error) {
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

	timerAction := TimerActionStop
	remainingTimerInS := 0

	_, err = db_access.Exec(
		ActCurrentTimerCommand,
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

		var userId int
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
