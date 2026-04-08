package srvgamesmock

import (
	typegames "FGG-Service/src/games/types"

	"github.com/stretchr/testify/mock"
)

type CurrentGameServiceMock struct {
	mock.Mock
}

func (m *CurrentGameServiceMock) GetCurrentGame(userId int) (typegames.CurrentGame, error) {
	args := m.Called(userId)
	return args.Get(0).(typegames.CurrentGame), args.Error(1)
}
