package srvgames_test

import (
	"FGG-Service/src/common"
	"FGG-Service/src/games/service"
	"FGG-Service/src/games/types"
	"FGG-Service/tests/games/mock"
	"database/sql"
	"errors"
	"testing"
	"time"
)

// GetCurrentGameCommand returns sql.ErrNoRows. The CurrentGameNotFoundError will return.
func TestSrvGames_GetCurrentGame_NotFound_NoRows(test *testing.T) {
	// Arrange
	userId := 1

	databaseMock := Initialize_GetCurrentGame_DatabaseMock_NotFound_NoRows()
	sut := srvgames.Service{Database: databaseMock}

	// Act
	_, err := sut.GetCurrentGame(userId)

	// Assert
	if err == nil {
		test.Fatalf("No error found")
	}

	var notFoundError *common.NotFoundError
	if !errors.As(err, &notFoundError) {
		test.Fatalf("Unexpected error type: %v", err)
	}

	databaseMock.AssertExpectations(test)
}

func Initialize_GetCurrentGame_DatabaseMock_NotFound_NoRows() *dbgamesmock.DatabaseMock {
	databaseMock := new(dbgamesmock.DatabaseMock)

	databaseMock.
		On(
			"GetCurrentGameCommand",
			1).
		Return(
			typegames.CurrentGames{}, sql.ErrNoRows)

	return databaseMock
}

// GetCurrentGameCommand returns an empty list. The CurrentGameNotFoundError will return.
func TestSrvGames_GetCurrentGame_NotFound_EmptyList(test *testing.T) {
	// Arrange
	userId := 1

	databaseMock := Initialize_GetCurrentGame_DatabaseMock_NotFound_EmptyList()
	sut := srvgames.Service{Database: databaseMock}

	// Act
	_, err := sut.GetCurrentGame(userId)

	// Assert
	if err == nil {
		test.Fatalf("No error found")
	}

	var notFoundError *common.NotFoundError
	if !errors.As(err, &notFoundError) {
		test.Fatalf("Unexpected error type: %v", err)
	}

	databaseMock.AssertExpectations(test)
}

func Initialize_GetCurrentGame_DatabaseMock_NotFound_EmptyList() *dbgamesmock.DatabaseMock {
	databaseMock := new(dbgamesmock.DatabaseMock)

	databaseMock.
		On(
			"GetCurrentGameCommand",
			1).
		Return(
			typegames.CurrentGames{}, nil)

	return databaseMock
}

// GetCurrentGameCommand returns a database error. The error will return.
func TestSrvGames_GetCurrentGame_DatabaseError(test *testing.T) {
	// Arrange
	userId := 1
	dbError := errors.New("database connection lost")

	databaseMock := Initialize_GetCurrentGame_DatabaseMock_DatabaseError(dbError)
	sut := srvgames.Service{Database: databaseMock}

	// Act
	_, err := sut.GetCurrentGame(userId)

	// Assert
	if err == nil {
		test.Fatalf("No error found")
	}

	if !errors.Is(err, dbError) {
		test.Fatalf("Unexpected error: %v", err)
	}

	databaseMock.AssertExpectations(test)
}

func Initialize_GetCurrentGame_DatabaseMock_DatabaseError(dbError error) *dbgamesmock.DatabaseMock {
	databaseMock := new(dbgamesmock.DatabaseMock)

	databaseMock.
		On(
			"GetCurrentGameCommand",
			1).
		Return(
			typegames.CurrentGames{typegames.CurrentGame{Id: 1, Name: "Half-Life 1", State: typegames.GameStateStarted}},
			dbError)

	return databaseMock
}

// GetCurrentGameCommand succeeds. GetGameTimeSpentCommand returns a database error. The error will return.
func TestSrvGames_GetCurrentGame_TimeSpentDatabaseError(test *testing.T) {
	// Arrange
	userId := 1
	dbError := errors.New("database connection lost")

	databaseMock := Initialize_GetCurrentGame_DatabaseMock_TimeSpentDatabaseError(dbError)
	sut := srvgames.Service{Database: databaseMock}

	// Act
	_, err := sut.GetCurrentGame(userId)

	// Assert
	if err == nil {
		test.Fatalf("No error found")
	}

	if !errors.Is(err, dbError) {
		test.Fatalf("Unexpected error: %v", err)
	}

	databaseMock.AssertExpectations(test)
}

func Initialize_GetCurrentGame_DatabaseMock_TimeSpentDatabaseError(dbError error) *dbgamesmock.DatabaseMock {
	databaseMock := new(dbgamesmock.DatabaseMock)

	databaseMock.
		On(
			"GetCurrentGameCommand",
			1).
		Return(
			typegames.CurrentGames{typegames.CurrentGame{Id: 1, Name: "Half-Life 1", State: typegames.GameStateStarted}},
			nil)
	databaseMock.
		On(
			"GetGameTimeSpentCommand",
			1, 1).
		Return(
			time.Duration(0), dbError)

	return databaseMock
}

// GetCurrentGameCommand and GetGameTimeSpentCommand succeed. The current game with TimeSpent will return.
func TestSrvGames_GetCurrentGame_SuccessReturn(test *testing.T) {
	// Arrange
	userId := 1
	expectedTimeSpent := 2 * time.Hour

	databaseMock := Initialize_GetCurrentGame_DatabaseMock_SuccessReturn(expectedTimeSpent)
	sut := srvgames.Service{Database: databaseMock}

	// Act
	game, err := sut.GetCurrentGame(userId)

	// Assert
	if err != nil {
		test.Fatalf("Unexpected error: %v", err)
	}

	if game.Id != 1 {
		test.Fatalf("Unexpected game id: %d", game.Id)
	}

	if game.Name != "Half-Life 1" {
		test.Fatalf("Unexpected game name: %s", game.Name)
	}

	if game.TimeSpent != expectedTimeSpent {
		test.Fatalf("Unexpected time spent: %v", game.TimeSpent)
	}

	databaseMock.AssertExpectations(test)
}

func Initialize_GetCurrentGame_DatabaseMock_SuccessReturn(timeSpent time.Duration) *dbgamesmock.DatabaseMock {
	databaseMock := new(dbgamesmock.DatabaseMock)

	databaseMock.
		On(
			"GetCurrentGameCommand",
			1).
		Return(
			typegames.CurrentGames{typegames.CurrentGame{Id: 1, Name: "Half-Life 1", State: typegames.GameStateStarted}},
			nil)
	databaseMock.
		On(
			"GetGameTimeSpentCommand",
			1, 1).
		Return(
			timeSpent, nil)

	return databaseMock
}
