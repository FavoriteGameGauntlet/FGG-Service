package srvgamestest

import (
	"FGG-Service/src/games/service"
	"FGG-Service/src/games/types"
	"testing"

	"github.com/stretchr/testify/mock"
)

type DatabaseMock struct {
	mock.Mock
}

func (m *DatabaseMock) DoesUnplayedGameExistCommand(userId int, gameName string) (doesExist bool, err error) {
	args := m.Called(userId, gameName)
	doesExist = args.Get(0).(bool)
	err = args.Error(1)
	return
}

func (m *DatabaseMock) DoesGameExistCommand(gameName string) (doesExist bool, err error) {
	args := m.Called(gameName)
	doesExist = args.Get(0).(bool)
	err = args.Error(1)
	return
}

func (m *DatabaseMock) GetGameCommand(gameName string) (game typegames.CurrentGame, err error) {
	args := m.Called(gameName)
	game = args.Get(0).(typegames.CurrentGame)
	err = args.Error(1)
	return
}

func (m *DatabaseMock) CreateUnplayedGameCommand(userId int, gameId int) error {
	args := m.Called(userId, gameId)
	return args.Error(0)
}

func TestAddUnplayedGames(t *testing.T) {
	// Arrange
	userId := 1
	games := typegames.WishlistGames{
		typegames.WishlistGame{Name: "Half-Life 1"},
	}

	databaseMock := InitializeMock()

	sut := srvgames.Service{Database: databaseMock}

	// Act
	_ = sut.AddUnplayedGames(userId, games)

	// Assert
	databaseMock.AssertExpectations(t)
}

func InitializeMock() DatabaseMock {
	databaseMock := new(DatabaseMock)

	databaseMock.
		On(
			"DoesUnplayedGameExistCommand",
			1, "Half-Life 1").
		Return(
			false, nil)
	databaseMock.
		On(
			"DoesGameExistCommand",
			"Half-Life 1").
		Return(
			true, nil)
	databaseMock.
		On(
			"GetGameCommand",
			"Half-Life 1").
		Return(
			typegames.CurrentGame{Name: "Half-Life 1"}, nil)
	databaseMock.
		On(
			"CreateUnplayedGameCommand",
			1, 1).
		Return(
			nil)

	return *databaseMock
}
