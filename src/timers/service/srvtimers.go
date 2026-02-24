package srvtimers

import (
	"FGG-Service/src/common"
	"FGG-Service/src/games/database"
	"FGG-Service/src/timers/database"
	"FGG-Service/src/timers/types"
	dbwheeleffects "FGG-Service/src/wheeleffects/database"
	"database/sql"
	"errors"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type Service struct {
	Database               dbtimers.Database
	GamesDatabase          dbgames.Database
	WheelEffectsDatabase   dbwheeleffects.Database
	TimerFinisherScheduler gocron.Scheduler
}

func NewService() *Service {
	s := new(Service)

	s.StartTimerFinisherScheduler()

	return s
}

func (s *Service) StartTimerFinisherScheduler() {
	scheduler, err := gocron.NewScheduler()

	if err != nil {
		panic(err)
	}

	_, err = scheduler.NewJob(
		gocron.DurationJob(1*time.Second),
		gocron.NewTask(s.StopAllCompletedTimers),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)

	if err != nil {
		panic(err)
	}

	s.TimerFinisherScheduler = scheduler

	scheduler.Start()
}

func (s *Service) GetOrCreateCurrentTimer(userId int) (timer typetimers.Timer, err error) {
	games, err := s.GamesDatabase.GetCurrentGameCommand(userId)

	if errors.Is(err, sql.ErrNoRows) {
		err = common.NewCurrentGameNotFoundError()
		return
	}

	if err != nil {
		return
	}

	game := games[0]

	timer, err = s.Database.GetCurrentTimerCommand(userId)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return
	}

	count, err := s.WheelEffectsDatabase.GetAvailableRollsCountCommand(userId)

	if err != nil {
		return
	}

	if count > 0 {
		err = common.NewAvailableRollsExistError()
		return
	}

	err = s.Database.CreateCurrentTimerCommand(userId, game.Id)

	if err != nil {
		return
	}

	timer, err = s.Database.GetCurrentTimerCommand(userId)

	return
}

func (s *Service) StartCurrentTimer(userId int) (typetimers.Timer, error) {
	return s.actCurrentTimer(
		userId,
		typetimers.TimerStateRunning,
		[]typetimers.TimerStateType{
			typetimers.TimerStateRunning,
			typetimers.TimerStateFinished,
		})
}

func (s *Service) PauseCurrentTimer(userId int) (typetimers.Timer, error) {
	return s.actCurrentTimer(
		userId,
		typetimers.TimerStatePaused,
		[]typetimers.TimerStateType{
			typetimers.TimerStateCreated,
			typetimers.TimerStatePaused,
			typetimers.TimerStateFinished,
		})
}

func (s *Service) StopCurrentTimer(userId int) (typetimers.Timer, error) {
	return s.actCurrentTimer(
		userId,
		typetimers.TimerStateFinished,
		[]typetimers.TimerStateType{
			typetimers.TimerStateCreated,
			typetimers.TimerStateFinished,
		})
}

func (s *Service) ForceStopCurrentTimer(userId int) (timer typetimers.Timer, err error) {
	timer, err = s.actCurrentTimer(
		userId,
		typetimers.TimerStateFinished,
		[]typetimers.TimerStateType{
			typetimers.TimerStateFinished,
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
	timerState typetimers.TimerStateType,
	incorrectStates []typetimers.TimerStateType) (timer typetimers.Timer, err error) {

	timer, err = s.Database.GetCurrentTimerCommand(userId)

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

	err = s.Database.ActTimerCommand(timer.Id, timerState, remainingTime)

	if err != nil {
		return
	}

	timer, err = s.Database.GetCurrentTimerCommand(userId)

	return
}

func (s *Service) StopAllCompletedTimers() error {
	userIds, err := s.Database.GetCompletedTimerUsersCommand()

	if err != nil {
		return err
	}

	for _, userId := range userIds {
		_, err = s.StopCurrentTimer(userId)
	}

	return nil
}
