package srvusers_test

import (
	"FGG-Service/src/users/service"
	"FGG-Service/src/users/types"
	"FGG-Service/tests/users/mock"
	"testing"

	"github.com/stretchr/testify/require"
)

type GetAllUserNamesTestCase struct {
	Name            string
	SetupMock       func() *dbusersmock.DatabaseMock
	ExpectedUsers   typeusers.Users
	ExpectedErrorIs error
}

func ptr(s string) *string { return &s }

var GetAllUserNamesTestCases = []GetAllUserNamesTestCase{
	{
		// GetAllUserNamesCommand returns a database error. The error will return.
		Name: "DatabaseError",
		SetupMock: func() *dbusersmock.DatabaseMock {
			databaseMock := new(dbusersmock.DatabaseMock)

			databaseMock.
				On("GetAllUserNamesCommand").
				Return(typeusers.Users{}, dbError)

			return databaseMock
		},
		ExpectedErrorIs: dbError,
	},
	{
		// GetAllUserNamesCommand returns an empty list. An empty list will return.
		Name: "EmptyList",
		SetupMock: func() *dbusersmock.DatabaseMock {
			databaseMock := new(dbusersmock.DatabaseMock)

			databaseMock.
				On("GetAllUserNamesCommand").
				Return(typeusers.Users{}, nil)

			return databaseMock
		},
		ExpectedUsers: typeusers.Users{},
	},
	{
		// GetAllUserNamesCommand returns users without display names. The list will return.
		Name: "SuccessNoDisplayNames",
		SetupMock: func() *dbusersmock.DatabaseMock {
			databaseMock := new(dbusersmock.DatabaseMock)

			databaseMock.
				On("GetAllUserNamesCommand").
				Return(typeusers.Users{
					{Login: "alice"},
					{Login: "bob"},
				}, nil)

			return databaseMock
		},
		ExpectedUsers: typeusers.Users{
			{Login: "alice"},
			{Login: "bob"},
		},
	},
	{
		// GetAllUserNamesCommand returns users with display names. The list will return.
		Name: "SuccessWithDisplayNames",
		SetupMock: func() *dbusersmock.DatabaseMock {
			databaseMock := new(dbusersmock.DatabaseMock)

			databaseMock.
				On("GetAllUserNamesCommand").
				Return(typeusers.Users{
					{Login: "alice", DisplayName: ptr("Alice")},
					{Login: "bob"},
				}, nil)

			return databaseMock
		},
		ExpectedUsers: typeusers.Users{
			{Login: "alice", DisplayName: ptr("Alice")},
			{Login: "bob"},
		},
	},
}

func TestSrvUsers_GetAllUserNames(test *testing.T) {
	for _, testCase := range GetAllUserNamesTestCases {
		test.Run(testCase.Name, func(test *testing.T) {
			// Arrange
			databaseMock := testCase.SetupMock()
			sut := srvusers.Service{Database: databaseMock}

			// Act
			users, err := sut.GetAllUserNames()

			// Assert
			if testCase.ExpectedErrorIs != nil {
				require.ErrorIs(test, err, testCase.ExpectedErrorIs)
			}

			if testCase.ExpectedUsers != nil {
				require.NoError(test, err)
				require.Equal(test, testCase.ExpectedUsers, users)
			}

			databaseMock.AssertExpectations(test)
		})
	}
}
