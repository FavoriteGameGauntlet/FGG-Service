package timer_service

import (
	"FGG-Service/src/common"
	"FGG-Service/src/game/game_db"
	"FGG-Service/src/timer/timer_db"
	"database/sql"
	"errors"
)

func GetOrCreateCurrentTimer(userId int) (timer common.Timer, err error) {
	games, err := game_db.GetCurrentGameCommand(userId)

	if errors.Is(err, sql.ErrNoRows) {
		err = common.NewCurrentGameNotFoundError()
		return
	}

	if err != nil {
		return
	}

	game := games[0]

	timer, err = timer_db.GetCurrentTimerCommand(userId)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return
	}

	err = timer_db.CreateCurrentTimerCommand(userId, game.Id)

	if err != nil {
		return
	}

	timer, err = timer_db.GetCurrentTimerCommand(userId)

	return
}

func StartCurrentTimer(userId int) (common.Timer, error) {
	return ActCurrentTimer(
		userId,
		common.TimerStateRunning,
		[]common.TimerStateType{
			common.TimerStateRunning,
			common.TimerStateFinished,
		})
}

func PauseCurrentTimer(userId int) (common.Timer, error) {
	return ActCurrentTimer(
		userId,
		common.TimerStatePaused,
		[]common.TimerStateType{
			common.TimerStateCreated,
			common.TimerStatePaused,
			common.TimerStateFinished,
		})
}

func StopCurrentTimer(userId int) (common.Timer, error) {
	return ActCurrentTimer(
		userId,
		common.TimerStateFinished,
		[]common.TimerStateType{
			common.TimerStateCreated,
			common.TimerStateFinished,
		})
}

func ForceStopCurrentTimer(userId int) (timer common.Timer, err error) {
	timer, err = ActCurrentTimer(
		userId,
		common.TimerStateFinished,
		[]common.TimerStateType{
			common.TimerStateFinished,
		})

	var notFoundError *common.NotFoundError
	if err != nil && !errors.As(err, &notFoundError) {
		return
	}

	err = nil
	return
}

func ActCurrentTimer(
	userId int,
	timerState common.TimerStateType,
	incorrectStates []common.TimerStateType) (timer common.Timer, err error) {

	timer, err = timer_db.GetCurrentTimerCommand(userId)

	if errors.Is(err, sql.ErrNoRows) {
		err = common.NewCurrentTimerNotFoundError()
		return
	}

	if err != nil {
		return
	}

	for _, state := range incorrectStates {
		if timer.State == state {
			err = common.NewCurrentTimerIncorrectStateError(timer.State)
			return
		}
	}

	remainingTime := timer.RemainingTime
	if remainingTime < 0 {
		remainingTime = 0
	}

	err = timer_db.ActTimerCommand(timer.Id, timerState, remainingTime)

	if err != nil {
		return
	}

	timer, err = timer_db.GetCurrentTimerCommand(userId)

	return
}

func StopAllCompletedTimers() error {
	userIds, err := timer_db.GetCompletedTimerUsersCommand()

	if err != nil {
		return err
	}

	for _, userId := range userIds {
		_, err = StopCurrentTimer(userId)
	}

	return nil
}
