package srvpoints_test

import (
	"FGG-Service/src/common"
	"FGG-Service/src/points/service"
	"FGG-Service/src/points/type"
	"FGG-Service/src/wheeleffects/types"
	"FGG-Service/tests/points/mock"
	"FGG-Service/tests/points/mock/srvwheeleffects"
	"testing"

	"github.com/stretchr/testify/require"
)

type ChangeFreePointsTestCase struct {
	Name                 string
	UserId               int
	Change               typepoints.FreePointChange
	SetupMock            func() *dbpointsmock.DatabaseMock
	SetupWheelEffectMock func() *srvwheeleffectsmock.ServiceMock
	SetupMinimum         func() (restore func())
	ExpectedActualChange *int
	ExpectedFinalValue   *int
	ExpectedErrorCode    string
}

var ChangeFreePointsTestCases = []ChangeFreePointsTestCase{
	{
		// Gain with 'base-teleport' is not allowed — only loss is valid. Conflict error returns.
		Name:   "Gain_BaseTeleport_ConflictError",
		UserId: 1,
		Change: typepoints.FreePointChange{ChangeSource: typepoints.FreePointsChangeSourceBaseTeleport, DesiredChangeValue: 5},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			return new(dbpointsmock.DatabaseMock)
		},
		ExpectedErrorCode: "WRONG_DESIRED_CHANGE_VALUE",
	},
	{
		// Gain with completely invalid source. Unprocessable error returns.
		Name:   "Gain_UnknownSource_UnprocessableError",
		UserId: 1,
		Change: typepoints.FreePointChange{ChangeSource: "invalid", DesiredChangeValue: 5},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			return new(dbpointsmock.DatabaseMock)
		},
		ExpectedErrorCode: "INCORRECT_CHANGE_SOURCE_VALUE",
	},
	{
		// Loss with invalid source. Unprocessable error returns.
		Name:   "Loss_UnknownSource_UnprocessableError",
		UserId: 1,
		Change: typepoints.FreePointChange{ChangeSource: "invalid", DesiredChangeValue: -5},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			return new(dbpointsmock.DatabaseMock)
		},
		ExpectedErrorCode: "INCORRECT_CHANGE_SOURCE_VALUE",
	},
	{
		// GetFreePointsCommand returns a database error.
		Name:   "Gain_DatabaseError",
		UserId: 1,
		Change: typepoints.FreePointChange{ChangeSource: typepoints.FreePointsChangeSourceQuestCompletion, DesiredChangeValue: 5},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetFreePointsCommand", 1).Return(0, dbError)
			return databaseMock
		},
		ExpectedErrorCode: "",
	},
	{
		// Gain with 'quest'. Success.
		Name:   "Gain_Quest_Success",
		UserId: 1,
		Change: typepoints.FreePointChange{ChangeSource: typepoints.FreePointsChangeSourceQuestCompletion, DesiredChangeValue: 10},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetFreePointsCommand", 1).Return(5, nil)
			databaseMock.On("ChangeFreePointsCommand", 1, 10).Return(nil)
			databaseMock.On("AddFreePointHistoryCommand", 1, 0, typepoints.FreePointsChangeSourceQuestCompletion, 10, 10, 15, (*int)(nil)).Return(nil)
			return databaseMock
		},
		ExpectedActualChange: ptr(10),
		ExpectedFinalValue:   ptr(15),
	},
	{
		// Gain with 'own-wheel-effect'. Success.
		Name:   "Gain_OwnWheelEffect_Success",
		UserId: 1,
		Change: typepoints.FreePointChange{ChangeSource: typepoints.FreePointsChangeSourceOwnWheelEffect, DesiredChangeValue: 3, WheelEffectName: ptrStr("test-effect")},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetFreePointsCommand", 1).Return(10, nil)
			databaseMock.On("ChangeFreePointsCommand", 1, 3).Return(nil)
			databaseMock.On("AddFreePointHistoryCommand", 1, 0, typepoints.FreePointsChangeSourceOwnWheelEffect, 3, 3, 13, ptr(42)).Return(nil)
			return databaseMock
		},
		SetupWheelEffectMock: func() *srvwheeleffectsmock.ServiceMock {
			weMock := new(srvwheeleffectsmock.ServiceMock)
			weMock.On("GetEffectHistoryByEffectName", 1, "test-effect").Return(&typewheeleffects.RolledWheelEffect{Id: 42}, nil)
			return weMock
		},
		ExpectedActualChange: ptr(3),
		ExpectedFinalValue:   ptr(13),
	},
	{
		// WheelEffectName is nil for 'own-wheel-effect' — required field missing. Unprocessable error returns.
		Name:   "Gain_OwnWheelEffect_NoName_UnprocessableError",
		UserId: 1,
		Change: typepoints.FreePointChange{ChangeSource: typepoints.FreePointsChangeSourceOwnWheelEffect, DesiredChangeValue: 3},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			return new(dbpointsmock.DatabaseMock)
		},
		ExpectedErrorCode: "WHEEL_EFFECT_NAME_REQUIRED",
	},
	{
		// GetEffectHistoryByEffectName returns a database error for 'own-wheel-effect'.
		Name:   "Gain_OwnWheelEffect_ServiceError",
		UserId: 1,
		Change: typepoints.FreePointChange{ChangeSource: typepoints.FreePointsChangeSourceOwnWheelEffect, DesiredChangeValue: 3, WheelEffectName: ptrStr("test-effect")},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			return new(dbpointsmock.DatabaseMock)
		},
		SetupWheelEffectMock: func() *srvwheeleffectsmock.ServiceMock {
			weMock := new(srvwheeleffectsmock.ServiceMock)
			weMock.On("GetEffectHistoryByEffectName", 1, "test-effect").Return(nil, dbError)
			return weMock
		},
	},
	{
		// GetEffectHistoryByEffectName returns nil — effect not found. Not found error returns.
		Name:   "Gain_OwnWheelEffect_NotFound_Error",
		UserId: 1,
		Change: typepoints.FreePointChange{ChangeSource: typepoints.FreePointsChangeSourceOwnWheelEffect, DesiredChangeValue: 3, WheelEffectName: ptrStr("unknown-effect")},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			return new(dbpointsmock.DatabaseMock)
		},
		SetupWheelEffectMock: func() *srvwheeleffectsmock.ServiceMock {
			weMock := new(srvwheeleffectsmock.ServiceMock)
			weMock.On("GetEffectHistoryByEffectName", 1, "unknown-effect").Return(nil, nil)
			return weMock
		},
		ExpectedErrorCode: "WHEEL_EFFECT_NAME_REQUIRED",
	},
	{
		// Gain with 'other-wheel-effect', non-zero SourceUserId. Success.
		Name:   "Gain_OtherWheelEffect_Success",
		UserId: 1,
		Change: typepoints.FreePointChange{SourceUserId: 2, ChangeSource: typepoints.FreePointsChangeSourceOtherWheelEffect, DesiredChangeValue: 5, WheelEffectName: ptrStr("test-effect")},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetFreePointsCommand", 1).Return(10, nil)
			databaseMock.On("ChangeFreePointsCommand", 1, 5).Return(nil)
			databaseMock.On("AddFreePointHistoryCommand", 1, 2, typepoints.FreePointsChangeSourceOtherWheelEffect, 5, 5, 15, ptr(42)).Return(nil)
			return databaseMock
		},
		SetupWheelEffectMock: func() *srvwheeleffectsmock.ServiceMock {
			weMock := new(srvwheeleffectsmock.ServiceMock)
			weMock.On("GetEffectHistoryByEffectName", 1, "test-effect").Return(&typewheeleffects.RolledWheelEffect{Id: 42}, nil)
			return weMock
		},
		ExpectedActualChange: ptr(5),
		ExpectedFinalValue:   ptr(15),
	},
	{
		// WheelEffectName is nil for 'other-wheel-effect'. Unprocessable error returns.
		Name:   "Gain_OtherWheelEffect_NoName_UnprocessableError",
		UserId: 1,
		Change: typepoints.FreePointChange{SourceUserId: 2, ChangeSource: typepoints.FreePointsChangeSourceOtherWheelEffect, DesiredChangeValue: 5},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			return new(dbpointsmock.DatabaseMock)
		},
		ExpectedErrorCode: "WHEEL_EFFECT_NAME_REQUIRED",
	},
	{
		// Gain with 'other'. Success.
		Name:   "Gain_Other_Success",
		UserId: 1,
		Change: typepoints.FreePointChange{ChangeSource: typepoints.FreePointsChangeSourceOther, DesiredChangeValue: 7},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetFreePointsCommand", 1).Return(0, nil)
			databaseMock.On("ChangeFreePointsCommand", 1, 7).Return(nil)
			databaseMock.On("AddFreePointHistoryCommand", 1, 0, typepoints.FreePointsChangeSourceOther, 7, 7, 7, (*int)(nil)).Return(nil)
			return databaseMock
		},
		ExpectedActualChange: ptr(7),
		ExpectedFinalValue:   ptr(7),
	},
	{
		// Loss with 'base-teleport'. Success.
		Name:   "Loss_BaseTeleport_Success",
		UserId: 1,
		Change: typepoints.FreePointChange{ChangeSource: typepoints.FreePointsChangeSourceBaseTeleport, DesiredChangeValue: -5},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetFreePointsCommand", 1).Return(10, nil)
			databaseMock.On("ChangeFreePointsCommand", 1, -5).Return(nil)
			databaseMock.On("AddFreePointHistoryCommand", 1, 0, typepoints.FreePointsChangeSourceBaseTeleport, -5, -5, 5, (*int)(nil)).Return(nil)
			return databaseMock
		},
		ExpectedActualChange: ptr(-5),
		ExpectedFinalValue:   ptr(5),
	},
	{
		// Loss exceeds current points — clamped to minimum (0). Success.
		Name:   "Loss_ClampedToMinimum_Success",
		UserId: 1,
		Change: typepoints.FreePointChange{ChangeSource: typepoints.FreePointsChangeSourceQuestCompletion, DesiredChangeValue: -20},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetFreePointsCommand", 1).Return(5, nil)
			databaseMock.On("ChangeFreePointsCommand", 1, -5).Return(nil)
			databaseMock.On("AddFreePointHistoryCommand", 1, 0, typepoints.FreePointsChangeSourceQuestCompletion, -20, -5, 0, (*int)(nil)).Return(nil)
			return databaseMock
		},
		ExpectedActualChange: ptr(-5),
		ExpectedFinalValue:   ptr(0),
	},
	{
		// Minimum disabled — loss can take points below 0. Success.
		Name:   "Loss_MinimumDisabled_Success",
		UserId: 1,
		Change: typepoints.FreePointChange{ChangeSource: typepoints.FreePointsChangeSourceSandStorm, DesiredChangeValue: -20},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetFreePointsCommand", 1).Return(5, nil)
			databaseMock.On("ChangeFreePointsCommand", 1, -20).Return(nil)
			databaseMock.On("AddFreePointHistoryCommand", 1, 0, typepoints.FreePointsChangeSourceSandStorm, -20, -20, -15, (*int)(nil)).Return(nil)
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
		Name:   "Zero_Quest_Success",
		UserId: 1,
		Change: typepoints.FreePointChange{ChangeSource: typepoints.FreePointsChangeSourceQuestCompletion, DesiredChangeValue: 0},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetFreePointsCommand", 1).Return(10, nil)
			databaseMock.On("ChangeFreePointsCommand", 1, 0).Return(nil)
			databaseMock.On("AddFreePointHistoryCommand", 1, 0, typepoints.FreePointsChangeSourceQuestCompletion, 0, 0, 10, (*int)(nil)).Return(nil)
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
			var weMock *srvwheeleffectsmock.ServiceMock
			if testCase.SetupWheelEffectMock != nil {
				weMock = testCase.SetupWheelEffectMock()
				sut.WheelEffectService = weMock
			}

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
			if weMock != nil {
				weMock.AssertExpectations(test)
			}
		})
	}
}
