package srvtimersmock

import (
	"FGG-Service/src/timers/types"

	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func (m *ServiceMock) ForceStopCurrentTimer(userId int) (typetimers.Timer, error) {
	args := m.Called(userId)
	return args.Get(0).(typetimers.Timer), args.Error(1)
}
