package srvgames_test

import (
	"FGG-Service/src/games/service"
	"FGG-Service/src/games/types"
	"FGG-Service/tests/games/mock"
	"testing"

	"github.com/stretchr/testify/require"
)

type GetUnplayedGamesTestCase struct {
	Name            string
	UserId          int
	SetupMock       func() *dbgamesmock.DatabaseMock
	ExpectedGames   typegames.WishlistGames
	ExpectedErrorIs error
}

var GetUnplayedGamesTestCases = []GetUnplayedGamesTestCase{
	{
		// GetUnplayedGamesCommand returns a database error. The error will return.
		Name:   "DatabaseError",
		UserId: 1,
		SetupMock: func() *dbgamesmock.DatabaseMock {
			databaseMock := new(dbgamesmock.DatabaseMock)

			databaseMock.
				On("GetUnplayedGamesCommand",
					1).
				Return(typegames.WishlistGames{}, dbError)

			return databaseMock
		},
		ExpectedErrorIs: dbError,
	},
	{
		// GetUnplayedGamesCommand returns an empty list. An empty list will return.
		Name:   "EmptyList",
		UserId: 1,
		SetupMock: func() *dbgamesmock.DatabaseMock {
			databaseMock := new(dbgamesmock.DatabaseMock)

			databaseMock.
				On("GetUnplayedGamesCommand",
					1).
				Return(typegames.WishlistGames{}, nil)

			return databaseMock
		},
		ExpectedGames: typegames.WishlistGames{},
	},
	{
		// GetUnplayedGamesCommand succeeds. The list of unplayed games will return.
		Name:   "SuccessReturn",
		UserId: 1,
		SetupMock: func() *dbgamesmock.DatabaseMock {
			databaseMock := new(dbgamesmock.DatabaseMock)

			databaseMock.
				On("GetUnplayedGamesCommand",
					1).
				Return(typegames.WishlistGames{
					{Id: 1, GameId: 10, Name: "Half-Life 1"},
					{Id: 2, GameId: 11, Name: "Half-Life 2"},
				}, nil)

			return databaseMock
		},
		ExpectedGames: typegames.WishlistGames{
			{Id: 1, GameId: 10, Name: "Half-Life 1"},
			{Id: 2, GameId: 11, Name: "Half-Life 2"},
		},
	},
}

func TestSrvGames_GetUnplayedGames(test *testing.T) {
	for _, testCase := range GetUnplayedGamesTestCases {
		test.Run(testCase.Name, func(test *testing.T) {
			// Arrange
			databaseMock := testCase.SetupMock()
			sut := srvgames.Service{Database: databaseMock}

			// Act
			games, err := sut.GetUnplayedGames(testCase.UserId)

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
