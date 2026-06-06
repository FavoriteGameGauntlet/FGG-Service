package dbgamesmock

import (
	"FGG-Service/src/games/types"
	"time"

	"github.com/stretchr/testify/mock"
)

type DatabaseMock struct {
	mock.Mock
}

func (m *DatabaseMock) DoesWishlistGameExistCommand(userId int, name string) (doesExist bool, err error) {
	args := m.Called(userId, name)
	doesExist = args.Get(0).(bool)
	err = args.Error(1)
	return
}

func (m *DatabaseMock) DoesGameExistCommand(name string) (doesExist bool, err error) {
	args := m.Called(name)
	doesExist = args.Get(0).(bool)
	err = args.Error(1)
	return
}

func (m *DatabaseMock) GetWishlistGameCommand(name string) (game typegames.WishlistGame, err error) {
	args := m.Called(name)
	game = args.Get(0).(typegames.WishlistGame)
	err = args.Error(1)
	return
}

func (m *DatabaseMock) CreateWishlistGameCommand(userId int, gameId int) error {
	args := m.Called(userId, gameId)
	return args.Error(0)
}

func (m *DatabaseMock) CreateGameCommand(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *DatabaseMock) DeleteUnplayedGameCommand(userId int, gameId int) error {
	//TODO implement me
	panic("implement me")
}

func (m *DatabaseMock) GetWishlistGamesCommand(userId int) (games typegames.WishlistGames, err error) {
	args := m.Called(userId)
	games = args.Get(0).(typegames.WishlistGames)
	err = args.Error(1)
	return
}

func (m *DatabaseMock) CreateCurrentGameCommand(userId int, gameId int) error {
	//TODO implement me
	panic("implement me")
}

func (m *DatabaseMock) GetCurrentGameCommand(userId int) (games typegames.CurrentGames, err error) {
	args := m.Called(userId)
	games = args.Get(0).(typegames.CurrentGames)
	err = args.Error(1)
	return
}

func (m *DatabaseMock) GetGameTimeSpentCommand(userId int, gameId int) (timeSpent time.Duration, err error) {
	args := m.Called(userId, gameId)
	timeSpent = args.Get(0).(time.Duration)
	err = args.Error(1)
	return
}

func (m *DatabaseMock) CancelCurrentGameCommand(userId int, gameId int) error {
	args := m.Called(userId, gameId)
	return args.Error(0)
}

func (m *DatabaseMock) FinishCurrentGameCommand(userId int, gameId int) error {
	//TODO implement me
	panic("implement me")
}

func (m *DatabaseMock) GetAllCurrentGamesCommand() (games typegames.CurrentGames, err error) {
	args := m.Called()
	games = args.Get(0).(typegames.CurrentGames)
	err = args.Error(1)
	return
}

func (m *DatabaseMock) GetGameHistoryCommand(userId int) (games typegames.CurrentGames, err error) {
	//TODO implement me
	panic("implement me")
}
