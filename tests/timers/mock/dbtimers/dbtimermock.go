package dbtimermock

import (
	"FGG-Service/src/timers/types"
	"time"

	"github.com/stretchr/testify/mock"
)

type DatabaseMock struct {
	mock.Mock
}

func (m *DatabaseMock) GetCurrentTimerCommand(userId int) (timer typetimers.Timer, err error) {
	args := m.Called(userId)
	timer = args.Get(0).(typetimers.Timer)
	err = args.Error(1)
	return
}

func (m *DatabaseMock) CreateCurrentTimerCommand(userId int, gameId int) error {
	args := m.Called(userId, gameId)
	return args.Error(0)
}

func (m *DatabaseMock) ActTimerCommand(timerId int, timerState typetimers.TimerStateType, remainingTime time.Duration) error {
	args := m.Called(timerId, timerState, remainingTime)
	return args.Error(0)
}

func (m *DatabaseMock) GetCompletedTimerUsersCommand() (userIds []int, err error) {
	args := m.Called()
	userIds = args.Get(0).([]int)
	err = args.Error(1)
	return
}
