package timer_service

import (
	"FGG-Service/common"
	"FGG-Service/game_service"
	"errors"
)

func GetOrCreateCurrentTimer(userId int) (timer common.Timer, err error) {
	game, err := game_service.GetCurrentGameCommand(userId)

	if err != nil {
		return
	}

	timer, err = GetCurrentTimerCommand(userId)

	var notFoundError *common.NotFoundError
	if err != nil && !errors.As(err, &notFoundError) {
		return
	}

	if errors.As(err, &notFoundError) {
		timer, err = CreateCurrentTimerCommand(userId, game.Id)
	}

	return
}

func StartCurrentTimer(userId int) (common.TimerAction, error) {
	return ActCurrentTimer(
		userId,
		common.TimerActionStart,
		[]common.TimerStateType{
			common.TimerStateRunning,
			common.TimerStateFinished,
		})
}

func PauseCurrentTimer(userId int) (common.TimerAction, error) {
	return ActCurrentTimer(
		userId,
		common.TimerActionPause,
		[]common.TimerStateType{
			common.TimerStateCreated,
			common.TimerStatePaused,
			common.TimerStateFinished,
		})
}

func StopCurrentTimer(userId int) (common.TimerAction, error) {
	return ActCurrentTimer(
		userId,
		common.TimerActionStop,
		[]common.TimerStateType{
			common.TimerStateCreated,
			common.TimerStateFinished,
		})
}

func ActCurrentTimer(
	userId int,
	actionType common.TimerActionType,
	incorrectStates []common.TimerStateType) (timerAction common.TimerAction, err error) {

	timer, err := GetCurrentTimerCommand(userId)

	if err != nil {
		return
	}

	for _, state := range incorrectStates {
		if timer.State == state {
			return common.TimerAction{}, common.NewCurrentTimerIncorrectStateError(timer.State)
		}
	}

	remainingTime := timer.RemainingTime
	if remainingTime < 0 {
		remainingTime = 0
	}

	err = ActTimerCommand(timer.Id, actionType, remainingTime)

	if err != nil {
		return
	}

	timerAction = common.TimerAction{
		Type:          actionType,
		RemainingTime: remainingTime,
	}

	return
}

func StopAllCompletedTimers() error {
	userIds, err := GetCompletedTimerUsersCommand()

	if err != nil {
		return err
	}

	for _, userId := range userIds {
		_, err = StopCurrentTimer(userId)
	}

	return nil
}
