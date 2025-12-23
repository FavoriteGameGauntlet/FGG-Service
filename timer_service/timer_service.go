package timer_service

import (
	"FGG-Service/common"
	"FGG-Service/db_access"
	"FGG-Service/effect_service"
	"FGG-Service/game_service"
	"database/sql"
	"errors"
	"time"
)

const (
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

func GetOrCreateCurrentTimer(userId int) (*common.Timer, error) {
	game, err := game_service.GetCurrentGame(userId)

	if err != nil {
		return nil, err
	}

	if game == nil {
		return nil, common.NewCurrentGameNotFoundError()
	}

	timer, err := GetCurrentTimer(userId)

	var notFoundError *common.CurrentTimerNotFoundError
	if err != nil && !errors.As(err, &notFoundError) {
		return nil, err
	}

	if timer != nil {
		return timer, nil
	}

	timer, err = CreateCurrentTimer(userId, game.Id)

	if err != nil {
		return nil, err
	}

	return timer, nil
}

func GetCurrentTimer(userId int) (*common.Timer, error) {
	row := db_access.QueryRow(
		GetCurrentTimerCommand,
		common.TimerActionStart,
		common.TimerActionPause,
		common.TimerActionStop,
		userId,
		common.TimerStateFinished,
	)

	timer := common.Timer{}
	var timerActionDateString *string
	err := row.Scan(
		&timer.Id,
		&timer.State,
		&timer.DurationInS,
		&timerActionDateString,
		&timer.RemainingTimeInS,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, common.NewCurrentTimerNotFoundError()
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

func CreateCurrentTimer(userId int, gameId int) (*common.Timer, error) {
	_, err := db_access.Exec(
		CreateTimerCommand,
		userId,
		gameId,
		common.DefaultTimerDurationInS,
	)

	if err != nil {
		return nil, err
	}

	return &common.Timer{
		DurationInS:      common.DefaultTimerDurationInS,
		State:            common.TimerStateCreated,
		RemainingTimeInS: common.DefaultTimerDurationInS,
	}, nil
}

func StartCurrentTimer(userId int) (*common.TimerAction, error) {
	game, err := game_service.GetCurrentGame(userId)

	if err != nil {
		return nil, err
	}

	if game == nil {
		return nil, common.NewCurrentGameNotFoundError()
	}

	timer, err := GetCurrentTimer(userId)

	if err != nil {
		return nil, err
	}

	if timer == nil {
		return nil, common.NewCurrentTimerNotFoundError()
	}

	if timer.State == common.TimerStateRunning ||
		timer.State == common.TimerStateFinished {
		return nil, common.NewCurrentTimerIncorrectStateError(timer.State)
	}

	timerAction := common.TimerActionStart
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

	return &common.TimerAction{
		Action:           timerAction,
		RemainingTimeInS: remainingTimerInS,
	}, nil
}

func PauseCurrentTimer(userId int) (*common.TimerAction, error) {
	game, err := game_service.GetCurrentGame(userId)

	if err != nil {
		return nil, err
	}

	if game == nil {
		return nil, common.NewCurrentGameNotFoundError()
	}

	timer, err := GetCurrentTimer(userId)

	if err != nil {
		return nil, err
	}

	if timer == nil {
		return nil, common.NewCurrentTimerNotFoundError()
	}

	if timer.State == common.TimerStateCreated ||
		timer.State == common.TimerStatePaused ||
		timer.State == common.TimerStateFinished {
		return nil, common.NewCurrentTimerIncorrectStateError(timer.State)
	}

	timerAction := common.TimerActionPause
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

	return &common.TimerAction{
		Action:           timerAction,
		RemainingTimeInS: remainingTimerInS,
	}, nil
}

func StopCurrentTimer(userId int) (*common.TimerAction, error) {
	timer, err := GetCurrentTimer(userId)

	if err != nil {
		return nil, err
	}

	if timer == nil {
		return nil, common.NewCurrentTimerNotFoundError()
	}

	if timer.State == common.TimerStateCreated ||
		timer.State == common.TimerStateFinished {
		return nil, common.NewCurrentTimerIncorrectStateError(timer.State)
	}

	timerAction := common.TimerActionStop
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

	return &common.TimerAction{
		Action:           timerAction,
		RemainingTimeInS: remainingTimerInS,
	}, nil
}

func StopAllCompletedTimers() error {
	rows, err := db_access.Query(GetCompletedTimerUsersCommand, common.TimerStateRunning, common.TimerStatePaused)

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

		err = effect_service.CreateAvailableRoll(userId)

		if err != nil {
			errorCount++
		}
	}

	if errorCount > 0 && errorCount == timerCount {
		return err
	}

	return nil
}
