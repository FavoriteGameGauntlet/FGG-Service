package srvpoints_test

import (
	"FGG-Service/src/points/service"
	typepoints "FGG-Service/src/points/type"
	"FGG-Service/tests/points/mock"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type GetUserTerritoryPointHistoryTestCase struct {
	Name            string
	UserId          int
	SetupMock       func() *dbpointsmock.DatabaseMock
	ExpectedHistory typepoints.TerritoryPointChangeHistories
	ExpectedErrorIs error
}

var territoryPointChangeHistoryEntry = typepoints.TerritoryPointChangeHistory{
	ActualChangeValue:  10,
	ChangeDate:         time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC),
	ChangeSource:       typepoints.TerritoryPointChangeSourceObtaining,
	DesiredChangeValue: 10,
	FinalValue:         20,
}

var GetUserTerritoryPointHistoryTestCases = []GetUserTerritoryPointHistoryTestCase{
	{
		// GetTerritoryPointHistoryCommand returns a database error. The error will return.
		Name:   "DatabaseError",
		UserId: 1,
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)

			databaseMock.
				On("GetTerritoryPointHistoryCommand", 1).
				Return(typepoints.TerritoryPointChangeHistories{}, dbError)

			return databaseMock
		},
		ExpectedErrorIs: dbError,
	},
	{
		// GetTerritoryPointHistoryCommand returns an empty list. An empty list will return.
		Name:   "EmptyList",
		UserId: 1,
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)

			databaseMock.
				On("GetTerritoryPointHistoryCommand", 1).
				Return(typepoints.TerritoryPointChangeHistories{}, nil)

			return databaseMock
		},
		ExpectedHistory: typepoints.TerritoryPointChangeHistories{},
	},
	{
		// GetTerritoryPointHistoryCommand succeeds. The list of territory point change history entries will return.
		Name:   "SuccessReturn",
		UserId: 1,
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)

			databaseMock.
				On("GetTerritoryPointHistoryCommand", 1).
				Return(typepoints.TerritoryPointChangeHistories{territoryPointChangeHistoryEntry}, nil)

			return databaseMock
		},
		ExpectedHistory: typepoints.TerritoryPointChangeHistories{territoryPointChangeHistoryEntry},
	},
}

func TestSrvPoints_GetUserTerritoryPointHistory(test *testing.T) {
	for _, testCase := range GetUserTerritoryPointHistoryTestCases {
		test.Run(testCase.Name, func(test *testing.T) {
			// Arrange
			databaseMock := testCase.SetupMock()
			sut := srvpoints.Service{Database: databaseMock}

			// Act
			history, err := sut.GetUserTerritoryPointHistory(testCase.UserId)

			// Assert
			if testCase.ExpectedErrorIs != nil {
				require.ErrorIs(test, err, testCase.ExpectedErrorIs)
			}

			if testCase.ExpectedHistory != nil {
				require.NoError(test, err)
				require.Equal(test, testCase.ExpectedHistory, history)
			}

			databaseMock.AssertExpectations(test)
		})
	}
}
