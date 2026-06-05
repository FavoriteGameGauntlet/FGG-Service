package srvpoints_test

import (
	"FGG-Service/src/points/service"
	"FGG-Service/tests/points/mock"
	"testing"

	"github.com/stretchr/testify/require"
)

type GetTerritoryPointsTestCase struct {
	Name           string
	UserId         int
	SetupMock      func() *dbpointsmock.DatabaseMock
	ExpectedPoints *int
	ExpectedError  error
}

var GetTerritoryPointsTestCases = []GetTerritoryPointsTestCase{
	{
		// GetTerritoryPointsCommand returns a database error. The error will return.
		Name:   "DatabaseError",
		UserId: 1,
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)

			databaseMock.
				On("GetTerritoryPointsCommand", 1).
				Return(0, dbError)

			return databaseMock
		},
		ExpectedError: dbError,
	},
	{
		// GetTerritoryPointsCommand succeeds and returns 0.
		Name:   "SuccessReturn_Zero",
		UserId: 1,
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)

			databaseMock.
				On("GetTerritoryPointsCommand", 1).
				Return(0, nil)

			return databaseMock
		},
		ExpectedPoints: ptr(0),
	},
	{
		// GetTerritoryPointsCommand succeeds and returns a positive value.
		Name:   "SuccessReturn_Positive",
		UserId: 1,
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)

			databaseMock.
				On("GetTerritoryPointsCommand", 1).
				Return(42, nil)

			return databaseMock
		},
		ExpectedPoints: ptr(42),
	},
}

func TestSrvPoints_GetTerritoryPoints(test *testing.T) {
	for _, testCase := range GetTerritoryPointsTestCases {
		test.Run(testCase.Name, func(test *testing.T) {
			// Arrange
			databaseMock := testCase.SetupMock()
			sut := srvpoints.Service{Database: databaseMock}

			// Act
			points, err := sut.GetTerritoryPoints(testCase.UserId)

			// Assert
			if testCase.ExpectedError != nil {
				require.ErrorIs(test, err, testCase.ExpectedError)
			}

			if testCase.ExpectedPoints != nil {
				require.NoError(test, err)
				require.Equal(test, *testCase.ExpectedPoints, points)
			}

			databaseMock.AssertExpectations(test)
		})
	}
}
