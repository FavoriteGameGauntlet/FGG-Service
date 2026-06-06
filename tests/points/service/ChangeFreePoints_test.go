package srvpoints_test

import (
	"FGG-Service/src/common"
	"FGG-Service/src/points/service"
	typepoints "FGG-Service/src/points/type"
	"FGG-Service/tests/points/mock"
	"testing"

	"github.com/stretchr/testify/require"
)

type ChangeFreePointsTestCase struct {
	Name                 string
	UserId               int
	Change               typepoints.PointChange
	SetupMock            func() *dbpointsmock.DatabaseMock
	SetupMinimum         func() (restore func())
	ExpectedActualChange *int
	ExpectedFinalValue   *int
	ExpectedErrorCode    string
}

var ChangeFreePointsTestCases = []ChangeFreePointsTestCase{
	{
		// Gain source 'no-path-to-base' is not allowed. Unprocessable error returns.
		Name:   "Gain_InvalidSource_UnprocessableError",
		UserId: 1,
		Change: typepoints.PointChange{ChangeSource: "no-path-to-base", DesiredChangeValue: 5},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			return new(dbpointsmock.DatabaseMock)
		},
		ExpectedErrorCode: "INCORRECT_CHANGE_SOURCE_VALUE",
	},
	{
		// Gain with completely invalid source. Unprocessable error returns.
		Name:   "Gain_UnknownSource_UnprocessableError",
		UserId: 1,
		Change: typepoints.PointChange{ChangeSource: "invalid", DesiredChangeValue: 5},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			return new(dbpointsmock.DatabaseMock)
		},
		ExpectedErrorCode: "INCORRECT_CHANGE_SOURCE_VALUE",
	},
	{
		// Loss with invalid source. Unprocessable error returns.
		Name:   "Loss_UnknownSource_UnprocessableError",
		UserId: 1,
		Change: typepoints.PointChange{ChangeSource: "invalid", DesiredChangeValue: -5},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			return new(dbpointsmock.DatabaseMock)
		},
		ExpectedErrorCode: "INCORRECT_CHANGE_SOURCE_VALUE",
	},
	{
		// GetFreePointsCommand returns a database error.
		Name:   "Gain_DatabaseError",
		UserId: 1,
		Change: typepoints.PointChange{ChangeSource: "quest-completion", DesiredChangeValue: 5},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetFreePointsCommand", 1).Return(0, dbError)
			return databaseMock
		},
		ExpectedErrorCode: "",
	},
	{
		// Gain with 'quest-completion'. Success.
		Name:   "Gain_QuestCompletion_Success",
		UserId: 1,
		Change: typepoints.PointChange{ChangeSource: "quest-completion", DesiredChangeValue: 10},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetFreePointsCommand", 1).Return(5, nil)
			databaseMock.On("ChangeFreePointsCommand", 1, 10).Return(nil)
			databaseMock.On("AddFreePointHistoryCommand", 1, "quest-completion", 10, 10, 15).Return(nil)
			return databaseMock
		},
		ExpectedActualChange: ptr(10),
		ExpectedFinalValue:   ptr(15),
	},
	{
		// Gain with 'wheel-effect'. Success.
		Name:   "Gain_WheelEffect_Success",
		UserId: 1,
		Change: typepoints.PointChange{ChangeSource: "wheel-effect", DesiredChangeValue: 3},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetFreePointsCommand", 1).Return(10, nil)
			databaseMock.On("ChangeFreePointsCommand", 1, 3).Return(nil)
			databaseMock.On("AddFreePointHistoryCommand", 1, "wheel-effect", 3, 3, 13).Return(nil)
			return databaseMock
		},
		ExpectedActualChange: ptr(3),
		ExpectedFinalValue:   ptr(13),
	},
	{
		// Gain with 'other'. Success.
		Name:   "Gain_Other_Success",
		UserId: 1,
		Change: typepoints.PointChange{ChangeSource: "other", DesiredChangeValue: 7},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetFreePointsCommand", 1).Return(0, nil)
			databaseMock.On("ChangeFreePointsCommand", 1, 7).Return(nil)
			databaseMock.On("AddFreePointHistoryCommand", 1, "other", 7, 7, 7).Return(nil)
			return databaseMock
		},
		ExpectedActualChange: ptr(7),
		ExpectedFinalValue:   ptr(7),
	},
	{
		// Loss with 'no-path-to-base'. Success.
		Name:   "Loss_NoPathToBase_Success",
		UserId: 1,
		Change: typepoints.PointChange{ChangeSource: "no-path-to-base", DesiredChangeValue: -5},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetFreePointsCommand", 1).Return(10, nil)
			databaseMock.On("ChangeFreePointsCommand", 1, -5).Return(nil)
			databaseMock.On("AddFreePointHistoryCommand", 1, "no-path-to-base", -5, -5, 5).Return(nil)
			return databaseMock
		},
		ExpectedActualChange: ptr(-5),
		ExpectedFinalValue:   ptr(5),
	},
	{
		// Loss exceeds current points — clamped to minimum (0). Success.
		Name:   "Loss_ClampedToMinimum_Success",
		UserId: 1,
		Change: typepoints.PointChange{ChangeSource: "quest-completion", DesiredChangeValue: -20},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetFreePointsCommand", 1).Return(5, nil)
			databaseMock.On("ChangeFreePointsCommand", 1, -5).Return(nil)
			databaseMock.On("AddFreePointHistoryCommand", 1, "quest-completion", -20, -5, 0).Return(nil)
			return databaseMock
		},
		ExpectedActualChange: ptr(-5),
		ExpectedFinalValue:   ptr(0),
	},
	{
		// Minimum disabled — loss can take points below 0. Success.
		Name:   "Loss_MinimumDisabled_Success",
		UserId: 1,
		Change: typepoints.PointChange{ChangeSource: "no-path-to-base", DesiredChangeValue: -20},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetFreePointsCommand", 1).Return(5, nil)
			databaseMock.On("ChangeFreePointsCommand", 1, -20).Return(nil)
			databaseMock.On("AddFreePointHistoryCommand", 1, "no-path-to-base", -20, -20, -15).Return(nil)
			return databaseMock
		},
		SetupMinimum: func() (restore func()) {
			original := common.FreePointMinimum
			common.FreePointMinimum = nil
			return func() { common.FreePointMinimum = original }
		},
		ExpectedActualChange: ptr(-20),
		ExpectedFinalValue:   ptr(-15),
	},
	{
		// Zero change with valid loss source. Success, no actual change.
		Name:   "Zero_QuestCompletion_Success",
		UserId: 1,
		Change: typepoints.PointChange{ChangeSource: "quest-completion", DesiredChangeValue: 0},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetFreePointsCommand", 1).Return(10, nil)
			databaseMock.On("ChangeFreePointsCommand", 1, 0).Return(nil)
			databaseMock.On("AddFreePointHistoryCommand", 1, "quest-completion", 0, 0, 10).Return(nil)
			return databaseMock
		},
		ExpectedActualChange: ptr(0),
		ExpectedFinalValue:   ptr(10),
	},
}

func TestSrvPoints_ChangeFreePoints(test *testing.T) {
	for _, testCase := range ChangeFreePointsTestCases {
		test.Run(testCase.Name, func(test *testing.T) {
			// Arrange
			if testCase.SetupMinimum != nil {
				restore := testCase.SetupMinimum()
				defer restore()
			}

			databaseMock := testCase.SetupMock()
			sut := srvpoints.Service{Database: databaseMock}

			// Act
			result, err := sut.ChangeFreePoints(testCase.UserId, testCase.Change)

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
