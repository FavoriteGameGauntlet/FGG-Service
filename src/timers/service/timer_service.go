package timer_service

import (
	"FGG-Service/src/common"
	"FGG-Service/src/games/db"
	"FGG-Service/src/timers/db"
	"database/sql"
	"errors"
)

type Service struct {
	Db timer_db.Database
}

func (s *Service) GetOrCreateCurrentTimer(userId int) (timer common.Timer, err error) {
	games, err := game_db.GetCurrentGameCommand(userId)

	if errors.Is(err, sql.ErrNoRows) {
		err = common.NewCurrentGameNotFoundError()
		return
	}

	if err != nil {
		return
	}

	game := games[0]

	timer, err = s.Db.GetCurrentTimerCommand(userId)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return
	}

	err = s.Db.CreateCurrentTimerCommand(userId, game.Id)

	if err != nil {
		return
	}

	timer, err = s.Db.GetCurrentTimerCommand(userId)

	return
}

func (s *Service) StartCurrentTimer(userId int) (common.Timer, error) {
	return s.actCurrentTimer(
		userId,
		common.TimerStateRunning,
		[]common.TimerStateType{
			common.TimerStateRunning,
			common.TimerStateFinished,
		})
}

func (s *Service) PauseCurrentTimer(userId int) (common.Timer, error) {
	return s.actCurrentTimer(
		userId,
		common.TimerStatePaused,
		[]common.TimerStateType{
			common.TimerStateCreated,
			common.TimerStatePaused,
			common.TimerStateFinished,
		})
}

func (s *Service) StopCurrentTimer(userId int) (common.Timer, error) {
	return s.actCurrentTimer(
		userId,
		common.TimerStateFinished,
		[]common.TimerStateType{
			common.TimerStateCreated,
			common.TimerStateFinished,
		})
}

func (s *Service) ForceStopCurrentTimer(userId int) (timer common.Timer, err error) {
	timer, err = s.actCurrentTimer(
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

func (s *Service) actCurrentTimer(
	userId int,
	timerState common.TimerStateType,
	incorrectStates []common.TimerStateType) (timer common.Timer, err error) {

	timer, err = s.Db.GetCurrentTimerCommand(userId)

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

	err = s.Db.ActTimerCommand(timer.Id, timerState, remainingTime)

	if err != nil {
		return
	}

	timer, err = s.Db.GetCurrentTimerCommand(userId)

	return
}

func (s *Service) StopAllCompletedTimers() error {
	userIds, err := s.Db.GetCompletedTimerUsersCommand()

	if err != nil {
		return err
	}

	for _, userId := range userIds {
		_, err = s.StopCurrentTimer(userId)
	}

	return nil
}
