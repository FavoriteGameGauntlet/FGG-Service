package srvpoints_test

import (
	"FGG-Service/src/common"
	"FGG-Service/src/points/service"
	typepoints "FGG-Service/src/points/type"
	"FGG-Service/tests/points/mock"
	"testing"

	"github.com/stretchr/testify/require"
)

type ChangeTerritoryHoursTestCase struct {
	Name                    string
	UserId                  int
	Change                  typepoints.TerritoryHoursChange
	SetupMock               func() *dbpointsmock.DatabaseMock
	ExpectedActualChange    *int
	ExpectedFinalValue      *int
	ExpectedErrorCode       string
}

var ChangeTerritoryHoursTestCases = []ChangeTerritoryHoursTestCase{
	{
		// Change source is not "seize" or "other". Unprocessable error returns.
		Name:   "InvalidSource_UnprocessableError",
		UserId: 1,
		Change: typepoints.TerritoryHoursChange{ChangeSource: "invalid", DesiredChangeValue: -2},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			return new(dbpointsmock.DatabaseMock)
		},
		ExpectedErrorCode: "INCORRECT_CHANGE_SOURCE_VALUE",
	},
	{
		// Source is "seize" with positive desired value. Conflict error returns.
		Name:   "Seize_PositiveValue_ConflictError",
		UserId: 1,
		Change: typepoints.TerritoryHoursChange{ChangeSource: "seize", DesiredChangeValue: 2},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			return new(dbpointsmock.DatabaseMock)
		},
		ExpectedErrorCode: "WRONG_DESIRED_CHANGE_VALUE",
	},
	{
		// Source is "seize" with zero desired value. Conflict error returns.
		Name:   "Seize_ZeroValue_ConflictError",
		UserId: 1,
		Change: typepoints.TerritoryHoursChange{ChangeSource: "seize", DesiredChangeValue: 0},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			return new(dbpointsmock.DatabaseMock)
		},
		ExpectedErrorCode: "WRONG_DESIRED_CHANGE_VALUE",
	},
	{
		// Source is "seize" with wrong decrease amount (not -2 or -4). Conflict error returns.
		Name:   "Seize_WrongDecreaseAmount_ConflictError",
		UserId: 1,
		Change: typepoints.TerritoryHoursChange{ChangeSource: "seize", DesiredChangeValue: -3},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			return new(dbpointsmock.DatabaseMock)
		},
		ExpectedErrorCode: "WRONG_DESIRED_CHANGE_VALUE",
	},
	{
		// Source is "seize" with isSomeones=true, wrong value (not 3 or 5). Conflict error returns.
		Name:   "Seize_IsSomeones_WrongDecreaseAmount_ConflictError",
		UserId: 1,
		Change: typepoints.TerritoryHoursChange{ChangeSource: "seize", DesiredChangeValue: -2, IsSomeones: true},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			return new(dbpointsmock.DatabaseMock)
		},
		ExpectedErrorCode: "WRONG_DESIRED_CHANGE_VALUE",
	},
	{
		// Source is "seize", value -2 valid but not enough hours. Conflict error returns.
		Name:   "Seize_NotEnoughHours_ConflictError",
		UserId: 1,
		Change: typepoints.TerritoryHoursChange{ChangeSource: "seize", DesiredChangeValue: -4},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetTerritoryHoursCommand", 1).Return(2, nil)
			return databaseMock
		},
		ExpectedErrorCode: "NOT_ENOUGH_CURRENT_POINTS",
	},
	{
		// GetTerritoryHoursCommand returns a database error.
		Name:   "Seize_DatabaseError",
		UserId: 1,
		Change: typepoints.TerritoryHoursChange{ChangeSource: "seize", DesiredChangeValue: -2},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetTerritoryHoursCommand", 1).Return(0, dbError)
			return databaseMock
		},
		ExpectedErrorCode: "",
	},
	{
		// Source is "seize", value -2, enough hours. Success.
		Name:   "Seize_DecreaseBy2_Success",
		UserId: 1,
		Change: typepoints.TerritoryHoursChange{ChangeSource: "seize", DesiredChangeValue: -2},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetTerritoryHoursCommand", 1).Return(10, nil)
			databaseMock.On("ChangeTerritoryHoursCommand", 1, -2).Return(nil)
			return databaseMock
		},
		ExpectedActualChange: ptr(-2),
		ExpectedFinalValue:   ptr(8),
	},
	{
		// Source is "seize", value -4, enough hours. Success.
		Name:   "Seize_DecreaseBy4_Success",
		UserId: 1,
		Change: typepoints.TerritoryHoursChange{ChangeSource: "seize", DesiredChangeValue: -4},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetTerritoryHoursCommand", 1).Return(10, nil)
			databaseMock.On("ChangeTerritoryHoursCommand", 1, -4).Return(nil)
			return databaseMock
		},
		ExpectedActualChange: ptr(-4),
		ExpectedFinalValue:   ptr(6),
	},
	{
		// Source is "seize" with isSomeones=true, value -3. Success.
		Name:   "Seize_IsSomeones_DecreaseBy3_Success",
		UserId: 1,
		Change: typepoints.TerritoryHoursChange{ChangeSource: "seize", DesiredChangeValue: -3, IsSomeones: true},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetTerritoryHoursCommand", 1).Return(10, nil)
			databaseMock.On("ChangeTerritoryHoursCommand", 1, -3).Return(nil)
			return databaseMock
		},
		ExpectedActualChange: ptr(-3),
		ExpectedFinalValue:   ptr(7),
	},
	{
		// Source is "seize" with isSomeones=true, value -5. Success.
		Name:   "Seize_IsSomeones_DecreaseBy5_Success",
		UserId: 1,
		Change: typepoints.TerritoryHoursChange{ChangeSource: "seize", DesiredChangeValue: -5, IsSomeones: true},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetTerritoryHoursCommand", 1).Return(10, nil)
			databaseMock.On("ChangeTerritoryHoursCommand", 1, -5).Return(nil)
			return databaseMock
		},
		ExpectedActualChange: ptr(-5),
		ExpectedFinalValue:   ptr(5),
	},
	{
		// Source is "other" with increase. Success.
		Name:   "Other_Increase_Success",
		UserId: 1,
		Change: typepoints.TerritoryHoursChange{ChangeSource: "other", DesiredChangeValue: 5},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetTerritoryHoursCommand", 1).Return(10, nil)
			databaseMock.On("ChangeTerritoryHoursCommand", 1, 5).Return(nil)
			return databaseMock
		},
		ExpectedActualChange: ptr(5),
		ExpectedFinalValue:   ptr(15),
	},
	{
		// Source is "other" with decrease. Success.
		Name:   "Other_Decrease_Success",
		UserId: 1,
		Change: typepoints.TerritoryHoursChange{ChangeSource: "other", DesiredChangeValue: -3},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetTerritoryHoursCommand", 1).Return(10, nil)
			databaseMock.On("ChangeTerritoryHoursCommand", 1, -3).Return(nil)
			return databaseMock
		},
		ExpectedActualChange: ptr(-3),
		ExpectedFinalValue:   ptr(7),
	},
	{
		// Source is "other", decrease exceeds current hours — clamped to zero.
		Name:   "Other_DecreaseClamped_Success",
		UserId: 1,
		Change: typepoints.TerritoryHoursChange{ChangeSource: "other", DesiredChangeValue: -20},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetTerritoryHoursCommand", 1).Return(5, nil)
			databaseMock.On("ChangeTerritoryHoursCommand", 1, -5).Return(nil)
			return databaseMock
		},
		ExpectedActualChange: ptr(-5),
		ExpectedFinalValue:   ptr(0),
	},
}

func TestSrvPoints_ChangeTerritoryHours(test *testing.T) {
	for _, testCase := range ChangeTerritoryHoursTestCases {
		test.Run(testCase.Name, func(test *testing.T) {
			// Arrange
			databaseMock := testCase.SetupMock()
			sut := srvpoints.Service{Database: databaseMock}

			// Act
			result, err := sut.ChangeTerritoryHours(testCase.UserId, testCase.Change)

			// Assert
			if testCase.ExpectedErrorCode != "" {
				require.Error(test, err)
				appErr, ok := err.(common.AppError)
				require.True(test, ok)
				require.Equal(test, testCase.ExpectedErrorCode, appErr.GetCode())
			} else if testCase.ExpectedActualChange != nil {
				require.NoError(test, err)
				require.Equal(test, *testCase.ExpectedActualChange, result.ActualChangeValue)
				require.Equal(test, *testCase.ExpectedFinalValue, result.FinalValue)
			} else {
				require.Error(test, err)
			}

			databaseMock.AssertExpectations(test)
		})
	}
}
