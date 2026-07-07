package srvgames_test

import (
	"FGG-Service/src/games/service"
	"FGG-Service/src/games/types"
	"FGG-Service/tests/games/mock"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type GetGameHistoryTestCase struct {
	Name            string
	UserId          int
	SetupMock       func() *dbgamesmock.DatabaseMock
	ExpectedGames   typegames.CurrentGames
	ExpectedErrorIs error
}

var GetGameHistoryTestCases = []GetGameHistoryTestCase{
	{
		// GetGameHistoryCommand returns a database error. The error will return.
		Name:   "GetGameHistoryCommand_DatabaseError",
		UserId: 1,
		SetupMock: func() *dbgamesmock.DatabaseMock {
			databaseMock := new(dbgamesmock.DatabaseMock)

			databaseMock.
				On("GetGameHistoryCommand",
					1).
				Return(typegames.CurrentGames{}, dbError)

			return databaseMock
		},
		ExpectedErrorIs: dbError,
	},
	{
		// GetGameHistoryCommand succeeds. GetGameTimeSpentCommand returns a database error. The error will return.
		Name:   "GetGameTimeSpentCommand_DatabaseError",
		UserId: 1,
		SetupMock: func() *dbgamesmock.DatabaseMock {
			databaseMock := new(dbgamesmock.DatabaseMock)

			databaseMock.
				On("GetGameHistoryCommand",
					1).
				Return(typegames.CurrentGames{
					typegames.CurrentGame{Id: 1, Name: "Half-Life 1", State: typegames.GameStateFinished}},
					nil)
			databaseMock.
				On("GetGameTimeSpentCommand",
					1, 1).
				Return(time.Duration(0), dbError)

			return databaseMock
		},
		ExpectedErrorIs: dbError,
	},
	{
		// GetGameHistoryCommand and GetGameTimeSpentCommand succeed for multiple games.
		// Each game in the returned history must have its own TimeSpent set.
		Name:   "SuccessReturn_MultipleGames",
		UserId: 1,
		SetupMock: func() *dbgamesmock.DatabaseMock {
			databaseMock := new(dbgamesmock.DatabaseMock)

			databaseMock.
				On("GetGameHistoryCommand",
					1).
				Return(typegames.CurrentGames{
					typegames.CurrentGame{Id: 1, Name: "Half-Life 1", State: typegames.GameStateFinished},
					typegames.CurrentGame{Id: 2, Name: "Half-Life 2", State: typegames.GameStateCancelled},
				}, nil)
			databaseMock.
				On("GetGameTimeSpentCommand",
					1, 1).
				Return(2*time.Hour, nil)
			databaseMock.
				On("GetGameTimeSpentCommand",
					1, 2).
				Return(30*time.Minute, nil)

			return databaseMock
		},
		ExpectedGames: typegames.CurrentGames{
			typegames.CurrentGame{Id: 1, Name: "Half-Life 1", State: typegames.GameStateFinished, TimeSpent: 2 * time.Hour},
			typegames.CurrentGame{Id: 2, Name: "Half-Life 2", State: typegames.GameStateCancelled, TimeSpent: 30 * time.Minute},
		},
	},
}

func TestSrvGames_GetGameHistory(test *testing.T) {
	for _, testCase := range GetGameHistoryTestCases {
		test.Run(testCase.Name, func(test *testing.T) {
			// Arrange
			databaseMock := testCase.SetupMock()
			sut := srvgames.Service{Database: databaseMock}

			// Act
			games, err := sut.GetGameHistory(testCase.UserId)

			// Assert
			if testCase.ExpectedErrorIs != nil {
				require.ErrorIs(test, err, testCase.ExpectedErrorIs)
			}

			if testCase.ExpectedGames != nil {
				require.NoError(test, err)
				require.Equal(test, testCase.ExpectedGames, games)
			}

			databaseMock.AssertExpectations(test)
		})
	}
}
