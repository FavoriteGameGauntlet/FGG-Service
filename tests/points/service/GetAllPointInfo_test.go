package srvpoints_test

import (
	"FGG-Service/src/points/service"
	typepoints "FGG-Service/src/points/type"
	"FGG-Service/tests/points/mock"
	"testing"

	"github.com/stretchr/testify/require"
)

type GetAllPointInfoTestCase struct {
	Name            string
	SetupMock       func() *dbpointsmock.DatabaseMock
	ExpectedInfos   typepoints.PointInfoByLogins
	ExpectedErrorIs error
}

var GetAllPointInfoTestCases = []GetAllPointInfoTestCase{
	{
		// GetAllPointInfoCommand returns a database error. The error will return.
		Name: "DatabaseError",
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)

			databaseMock.
				On("GetAllPointInfoCommand").
				Return(typepoints.PointInfoByLogins{}, dbError)

			return databaseMock
		},
		ExpectedErrorIs: dbError,
	},
	{
		// GetAllPointInfoCommand returns an empty list. An empty list will return.
		Name: "EmptyList",
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)

			databaseMock.
				On("GetAllPointInfoCommand").
				Return(typepoints.PointInfoByLogins{}, nil)

			return databaseMock
		},
		ExpectedInfos: typepoints.PointInfoByLogins{},
	},
	{
		// GetAllPointInfoCommand succeeds. The list of point info by login will return.
		Name: "SuccessReturn",
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)

			databaseMock.
				On("GetAllPointInfoCommand").
				Return(typepoints.PointInfoByLogins{
					{Login: "alice", PointInfo: pointInfo},
					{Login: "bob", PointInfo: pointInfo},
				}, nil)

			return databaseMock
		},
		ExpectedInfos: typepoints.PointInfoByLogins{
			{Login: "alice", PointInfo: pointInfo},
			{Login: "bob", PointInfo: pointInfo},
		},
	},
}

func TestSrvPoints_GetAllPointInfo(test *testing.T) {
	for _, testCase := range GetAllPointInfoTestCases {
		test.Run(testCase.Name, func(test *testing.T) {
			// Arrange
			databaseMock := testCase.SetupMock()
			sut := srvpoints.Service{Database: databaseMock}

			// Act
			infos, err := sut.GetAllPointInfo()

			// Assert
			if testCase.ExpectedErrorIs != nil {
				require.ErrorIs(test, err, testCase.ExpectedErrorIs)
			}

			if testCase.ExpectedInfos != nil {
				require.NoError(test, err)
				require.Equal(test, testCase.ExpectedInfos, infos)
			}

			databaseMock.AssertExpectations(test)
		})
	}
}
