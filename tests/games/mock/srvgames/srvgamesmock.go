package srvgamesmock

import (
	typegames "FGG-Service/src/games/types"

	"github.com/stretchr/testify/mock"
)

type GettingServiceMock struct {
	mock.Mock
}

func (m *GettingServiceMock) GetWishlistGames(userId int) (typegames.WishlistGames, error) {
	args := m.Called(userId)
	return args.Get(0).(typegames.WishlistGames), args.Error(1)
}

func (m *GettingServiceMock) GetCurrentGame(userId int) (typegames.CurrentGame, error) {
	args := m.Called(userId)
	return args.Get(0).(typegames.CurrentGame), args.Error(1)
}
