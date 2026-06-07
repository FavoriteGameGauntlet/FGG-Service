package srvpoints_test

import (
	"FGG-Service/src/points/service"
	typepoints "FGG-Service/src/points/type"
	"FGG-Service/tests/points/mock"
	"testing"

	"github.com/stretchr/testify/require"
)

type GetUserPointInfoTestCase struct {
	Name          string
	UserId        int
	SetupMock     func() *dbpointsmock.DatabaseMock
	ExpectedInfo  *typepoints.PointInfo
	ExpectedError error
}

var pointInfo = typepoints.PointInfo{
	TerritoryPoints:  1,
	FreePoints:       2,
	AvailableRolls:   3,
	TerritoryHours:   4,
	ExperiencePoints: 5,
}

var GetUserPointInfoTestCases = []GetUserPointInfoTestCase{
	{
		// GetPointInfoCommand returns a database error. The error will return.
		Name:   "DatabaseError",
		UserId: 1,
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)

			databaseMock.
				On("GetPointInfoCommand", 1).
				Return(typepoints.PointInfo{}, dbError)

			return databaseMock
		},
		ExpectedError: dbError,
	},
	{
		// GetPointInfoCommand succeeds. The point info will return.
		Name:   "SuccessReturn",
		UserId: 1,
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)

			databaseMock.
				On("GetPointInfoCommand", 1).
				Return(pointInfo, nil)

			return databaseMock
		},
		ExpectedInfo: &pointInfo,
	},
}

func TestSrvPoints_GetUserPointInfo(test *testing.T) {
	for _, testCase := range GetUserPointInfoTestCases {
		test.Run(testCase.Name, func(test *testing.T) {
			// Arrange
			databaseMock := testCase.SetupMock()
			sut := srvpoints.Service{Database: databaseMock}

			// Act
			info, err := sut.GetUserPointInfo(testCase.UserId)

			// Assert
			if testCase.ExpectedError != nil {
				require.ErrorIs(test, err, testCase.ExpectedError)
			}

			if testCase.ExpectedInfo != nil {
				require.NoError(test, err)
				require.Equal(test, *testCase.ExpectedInfo, info)
			}

			databaseMock.AssertExpectations(test)
		})
	}
}
