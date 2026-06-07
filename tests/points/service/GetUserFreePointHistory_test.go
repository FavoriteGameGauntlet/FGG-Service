package srvpoints_test

import (
	"FGG-Service/src/points/service"
	typepoints "FGG-Service/src/points/type"
	"FGG-Service/tests/points/mock"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type GetUserFreePointHistoryTestCase struct {
	Name            string
	UserId          int
	SetupMock       func() *dbpointsmock.DatabaseMock
	ExpectedHistory typepoints.FreePointChangeHistories
	ExpectedErrorIs error
}

var freePointChangeHistoryEntry = typepoints.FreePointChangeHistory{
	ActualChangeValue:  10,
	ChangeDate:         time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC),
	ChangeSource:       typepoints.FreePointsChangeSourceQuestCompletion,
	DesiredChangeValue: 10,
	FinalValue:         20,
}

var GetUserFreePointHistoryTestCases = []GetUserFreePointHistoryTestCase{
	{
		// GetFreePointHistoryCommand returns a database error. The error will return.
		Name:   "DatabaseError",
		UserId: 1,
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)

			databaseMock.
				On("GetFreePointHistoryCommand", 1).
				Return(typepoints.FreePointChangeHistories{}, dbError)

			return databaseMock
		},
		ExpectedErrorIs: dbError,
	},
	{
		// GetFreePointHistoryCommand returns an empty list. An empty list will return.
		Name:   "EmptyList",
		UserId: 1,
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)

			databaseMock.
				On("GetFreePointHistoryCommand", 1).
				Return(typepoints.FreePointChangeHistories{}, nil)

			return databaseMock
		},
		ExpectedHistory: typepoints.FreePointChangeHistories{},
	},
	{
		// GetFreePointHistoryCommand succeeds. The list of free point change history entries will return.
		Name:   "SuccessReturn",
		UserId: 1,
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)

			databaseMock.
				On("GetFreePointHistoryCommand", 1).
				Return(typepoints.FreePointChangeHistories{freePointChangeHistoryEntry}, nil)

			return databaseMock
		},
		ExpectedHistory: typepoints.FreePointChangeHistories{freePointChangeHistoryEntry},
	},
}

func TestSrvPoints_GetUserFreePointHistory(test *testing.T) {
	for _, testCase := range GetUserFreePointHistoryTestCases {
		test.Run(testCase.Name, func(test *testing.T) {
			// Arrange
			databaseMock := testCase.SetupMock()
			sut := srvpoints.Service{Database: databaseMock}

			// Act
			history, err := sut.GetUserFreePointHistory(testCase.UserId)

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
