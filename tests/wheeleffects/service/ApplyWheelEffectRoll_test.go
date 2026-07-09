package srvwheeleffects_test

import (
	"FGG-Service/src/common"
	"FGG-Service/src/points/service"
	"FGG-Service/src/points/type"
	"FGG-Service/src/sysparams/types"
	"FGG-Service/src/wheeleffects/service"
	"FGG-Service/src/wheeleffects/types"
	"FGG-Service/tests/points/mock"
	"FGG-Service/tests/sysparams/srvmock"
	"FGG-Service/tests/timers/mock/dbwheeleffects"
	"testing"

	"github.com/stretchr/testify/require"
)

type ApplyWheelEffectRollTestCase struct {
	Name                 string
	UserId               int
	RollApply            typewheeleffects.WheelEffectRollApply
	SetupWheelEffectMock func() *dbwheeleffectsmock.DatabaseMock
	SetupPointsMock      func() *dbpointsmock.DatabaseMock
	SetupSysParams       func() *srvsysparamsmock.ServiceMock
	ExpectedResults      typepoints.PointChangeResultByUserIds
	ExpectedError        error
	ExpectedErrorCode    string
}

var ApplyWheelEffectRollTestCases = []ApplyWheelEffectRollTestCase{
	{
		// The wheel effect name doesn't match any of the user's last rolled effects. The error returns and nothing else happens.
		Name:   "EffectRoll_NotFound_Error",
		UserId: 1,
		RollApply: typewheeleffects.WheelEffectRollApply{
			WheelEffectName: "unknown-effect",
			PointChangeByUserIds: typepoints.PointChangeByUserIds{
				{Login: "alice", UserId: 2, FreePointChange: typepoints.FreePointChange{SourceUserId: 1, ChangeSource: typepoints.FreePointsChangeSourceQuestCompletion, DesiredChangeValue: 5}},
			},
		},
		SetupWheelEffectMock: func() *dbwheeleffectsmock.DatabaseMock {
			databaseMock := new(dbwheeleffectsmock.DatabaseMock)

			databaseMock.
				On("GetLastRolledWheelEffectsCommand", 1).
				Return(typewheeleffects.RolledWheelEffects{
					{Id: 42, Name: "test-effect", IsApplied: false},
				}, nil)

			return databaseMock
		},
		SetupPointsMock: func() *dbpointsmock.DatabaseMock {
			return new(dbpointsmock.DatabaseMock)
		},
		ExpectedErrorCode: "WHEEL_EFFECT_NAME_NOT_FOUND",
	},
	{
		// The roll is marked as applied and history is recorded, but ChangeFreePoints fails. The error returns.
		Name:   "ChangeFreePoints_Error",
		UserId: 1,
		RollApply: typewheeleffects.WheelEffectRollApply{
			WheelEffectName: "test-effect",
			PointChangeByUserIds: typepoints.PointChangeByUserIds{
				{Login: "alice", UserId: 2, FreePointChange: typepoints.FreePointChange{SourceUserId: 1, ChangeSource: typepoints.FreePointsChangeSourceQuestCompletion, DesiredChangeValue: 5}},
			},
		},
		SetupWheelEffectMock: func() *dbwheeleffectsmock.DatabaseMock {
			databaseMock := new(dbwheeleffectsmock.DatabaseMock)

			databaseMock.
				On("GetLastRolledWheelEffectsCommand", 1).
				Return(typewheeleffects.RolledWheelEffects{
					{Id: 42, Name: "test-effect", IsApplied: false},
				}, nil)
			databaseMock.On("MarkLastWheelEffectAppliedCommand", 1, 42).Return(nil)
			databaseMock.On("AddWheelEffectHistoryCommand", 1, 42).Return(100, nil)

			return databaseMock
		},
		SetupPointsMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetFreePointsCommand", 2).Return(0, dbError)
			return databaseMock
		},
		ExpectedError: dbError,
	},
	{
		// Marking the roll as applied fails. The error returns and neither the history nor ChangeFreePoints are called.
		Name:   "MarkEffectRollAsApplied_Error",
		UserId: 1,
		RollApply: typewheeleffects.WheelEffectRollApply{
			WheelEffectName: "test-effect",
			PointChangeByUserIds: typepoints.PointChangeByUserIds{
				{Login: "alice", UserId: 2, FreePointChange: typepoints.FreePointChange{SourceUserId: 1, ChangeSource: typepoints.FreePointsChangeSourceQuestCompletion, DesiredChangeValue: 5}},
			},
		},
		SetupWheelEffectMock: func() *dbwheeleffectsmock.DatabaseMock {
			databaseMock := new(dbwheeleffectsmock.DatabaseMock)

			databaseMock.
				On("GetLastRolledWheelEffectsCommand", 1).
				Return(typewheeleffects.RolledWheelEffects{
					{Id: 42, Name: "test-effect", IsApplied: false},
				}, nil)
			databaseMock.On("MarkLastWheelEffectAppliedCommand", 1, 42).Return(dbError)

			return databaseMock
		},
		SetupPointsMock: func() *dbpointsmock.DatabaseMock {
			return new(dbpointsmock.DatabaseMock)
		},
		ExpectedError: dbError,
	},
	{
		// Marking the roll as applied succeeds but recording the history fails. The error returns and ChangeFreePoints is never called.
		Name:   "AddWheelEffectHistory_Error",
		UserId: 1,
		RollApply: typewheeleffects.WheelEffectRollApply{
			WheelEffectName: "test-effect",
			PointChangeByUserIds: typepoints.PointChangeByUserIds{
				{Login: "alice", UserId: 2, FreePointChange: typepoints.FreePointChange{SourceUserId: 1, ChangeSource: typepoints.FreePointsChangeSourceQuestCompletion, DesiredChangeValue: 5}},
			},
		},
		SetupWheelEffectMock: func() *dbwheeleffectsmock.DatabaseMock {
			databaseMock := new(dbwheeleffectsmock.DatabaseMock)

			databaseMock.
				On("GetLastRolledWheelEffectsCommand", 1).
				Return(typewheeleffects.RolledWheelEffects{
					{Id: 42, Name: "test-effect", IsApplied: false},
				}, nil)
			databaseMock.On("MarkLastWheelEffectAppliedCommand", 1, 42).Return(nil)
			databaseMock.On("AddWheelEffectHistoryCommand", 1, 42).Return(0, dbError)

			return databaseMock
		},
		SetupPointsMock: func() *dbpointsmock.DatabaseMock {
			return new(dbpointsmock.DatabaseMock)
		},
		ExpectedError: dbError,
	},
	{
		// Every change applies successfully and the roll is marked as applied with its history recorded. Results return in request order.
		Name:   "Success",
		UserId: 1,
		RollApply: typewheeleffects.WheelEffectRollApply{
			WheelEffectName: "test-effect",
			PointChangeByUserIds: typepoints.PointChangeByUserIds{
				{Login: "alice", UserId: 2, FreePointChange: typepoints.FreePointChange{SourceUserId: 1, ChangeSource: typepoints.FreePointsChangeSourceQuestCompletion, DesiredChangeValue: 5}},
				{Login: "bob", UserId: 3, FreePointChange: typepoints.FreePointChange{SourceUserId: 1, ChangeSource: typepoints.FreePointsChangeSourceOther, DesiredChangeValue: -2}},
			},
		},
		SetupWheelEffectMock: func() *dbwheeleffectsmock.DatabaseMock {
			databaseMock := new(dbwheeleffectsmock.DatabaseMock)

			databaseMock.
				On("GetLastRolledWheelEffectsCommand", 1).
				Return(typewheeleffects.RolledWheelEffects{
					{Id: 42, Name: "test-effect", IsApplied: false},
				}, nil)
			databaseMock.On("MarkLastWheelEffectAppliedCommand", 1, 42).Return(nil)
			databaseMock.On("AddWheelEffectHistoryCommand", 1, 42).Return(100, nil)

			return databaseMock
		},
		SetupPointsMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetFreePointsCommand", 2).Return(10, nil)
			databaseMock.On("ChangeFreePointsCommand", 2, 5).Return(nil)
			databaseMock.On("AddFreePointHistoryCommand", 2, 1, typepoints.FreePointsChangeSourceQuestCompletion, 5, 5, 15, ptr(100)).Return(nil)
			databaseMock.On("GetFreePointsCommand", 3).Return(4, nil)
			databaseMock.On("ChangeFreePointsCommand", 3, -2).Return(nil)
			databaseMock.On("AddFreePointHistoryCommand", 3, 1, typepoints.FreePointsChangeSourceOther, -2, -2, 2, ptr(100)).Return(nil)
			return databaseMock
		},
		SetupSysParams: func() *srvsysparamsmock.ServiceMock {
			spSvc := new(srvsysparamsmock.ServiceMock)
			spSvc.On("GetInt", typesysparams.ParamFreePointsMinimum).Return(0, nil)
			spSvc.On("GetBool", typesysparams.ParamShouldLimitFreePoints).Return(true, nil)
			return spSvc
		},
		ExpectedResults: typepoints.PointChangeResultByUserIds{
			{Login: "alice", UserId: 2, ChangeResults: typepoints.PointChangeResultByTypes{
				typepoints.PointTypeFreePoints: {ActualChangeValue: 5, ChangeSource: typepoints.FreePointsChangeSourceQuestCompletion, DesiredChangeValue: 5, FinalValue: 15},
			}},
			{Login: "bob", UserId: 3, ChangeResults: typepoints.PointChangeResultByTypes{
				typepoints.PointTypeFreePoints: {ActualChangeValue: -2, ChangeSource: typepoints.FreePointsChangeSourceOther, DesiredChangeValue: -2, FinalValue: 2},
			}},
		},
	},
	{
		// A per-user roll change applies via ChangeAvailableRolls right after that user's free points change succeeds.
		Name:   "Success_WithRollChanges",
		UserId: 1,
		RollApply: typewheeleffects.WheelEffectRollApply{
			WheelEffectName: "test-effect",
			PointChangeByUserIds: typepoints.PointChangeByUserIds{
				{Login: "alice", UserId: 2, FreePointChange: typepoints.FreePointChange{SourceUserId: 1, ChangeSource: typepoints.FreePointsChangeSourceQuestCompletion, DesiredChangeValue: 5}, AvailableRollChange: &typepoints.PointChange{ChangeSource: typepoints.FreePointsChangeSourceQuestCompletion, DesiredChangeValue: 1}},
				{Login: "bob", UserId: 3, FreePointChange: typepoints.FreePointChange{SourceUserId: 1, ChangeSource: typepoints.FreePointsChangeSourceOther, DesiredChangeValue: 0}, AvailableRollChange: &typepoints.PointChange{ChangeSource: typepoints.FreePointsChangeSourceOther, DesiredChangeValue: 2}},
			},
		},
		SetupWheelEffectMock: func() *dbwheeleffectsmock.DatabaseMock {
			databaseMock := new(dbwheeleffectsmock.DatabaseMock)

			databaseMock.
				On("GetLastRolledWheelEffectsCommand", 1).
				Return(typewheeleffects.RolledWheelEffects{
					{Id: 42, Name: "test-effect", IsApplied: false},
				}, nil)
			databaseMock.On("MarkLastWheelEffectAppliedCommand", 1, 42).Return(nil)
			databaseMock.On("AddWheelEffectHistoryCommand", 1, 42).Return(100, nil)
			databaseMock.On("GetAvailableRollsCountCommand", 2).Return(3, nil)
			databaseMock.On("GetAvailableRollsCountCommand", 3).Return(2, nil)

			return databaseMock
		},
		SetupPointsMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetFreePointsCommand", 2).Return(10, nil)
			databaseMock.On("ChangeFreePointsCommand", 2, 5).Return(nil)
			databaseMock.On("AddFreePointHistoryCommand", 2, 1, typepoints.FreePointsChangeSourceQuestCompletion, 5, 5, 15, ptr(100)).Return(nil)
			databaseMock.On("ChangeAvailableRollsCommand", 2, 1).Return(nil)
			databaseMock.On("GetFreePointsCommand", 3).Return(4, nil)
			databaseMock.On("ChangeFreePointsCommand", 3, 0).Return(nil)
			databaseMock.On("AddFreePointHistoryCommand", 3, 1, typepoints.FreePointsChangeSourceOther, 0, 0, 4, ptr(100)).Return(nil)
			databaseMock.On("ChangeAvailableRollsCommand", 3, 2).Return(nil)
			return databaseMock
		},
		SetupSysParams: func() *srvsysparamsmock.ServiceMock {
			spSvc := new(srvsysparamsmock.ServiceMock)
			spSvc.On("GetInt", typesysparams.ParamFreePointsMinimum).Return(0, nil)
			spSvc.On("GetBool", typesysparams.ParamShouldLimitFreePoints).Return(true, nil)
			return spSvc
		},
		ExpectedResults: typepoints.PointChangeResultByUserIds{
			{Login: "alice", UserId: 2, ChangeResults: typepoints.PointChangeResultByTypes{
				typepoints.PointTypeFreePoints:     {ActualChangeValue: 5, ChangeSource: typepoints.FreePointsChangeSourceQuestCompletion, DesiredChangeValue: 5, FinalValue: 15},
				typepoints.PointTypeAvailableRolls: {ActualChangeValue: 1, ChangeSource: typepoints.FreePointsChangeSourceQuestCompletion, DesiredChangeValue: 1, FinalValue: 4},
			}},
			{Login: "bob", UserId: 3, ChangeResults: typepoints.PointChangeResultByTypes{
				typepoints.PointTypeFreePoints:     {ActualChangeValue: 0, ChangeSource: typepoints.FreePointsChangeSourceOther, DesiredChangeValue: 0, FinalValue: 4},
				typepoints.PointTypeAvailableRolls: {ActualChangeValue: 2, ChangeSource: typepoints.FreePointsChangeSourceOther, DesiredChangeValue: 2, FinalValue: 4},
			}},
		},
	},
	{
		// A negative roll change is rejected outright; ChangeAvailableRolls is never called.
		Name:   "AvailableRollChange_Negative_Error",
		UserId: 1,
		RollApply: typewheeleffects.WheelEffectRollApply{
			WheelEffectName: "test-effect",
			PointChangeByUserIds: typepoints.PointChangeByUserIds{
				{Login: "alice", UserId: 2, FreePointChange: typepoints.FreePointChange{SourceUserId: 1, ChangeSource: typepoints.FreePointsChangeSourceOther, DesiredChangeValue: 0}, AvailableRollChange: &typepoints.PointChange{ChangeSource: typepoints.FreePointsChangeSourceOther, DesiredChangeValue: -1}},
			},
		},
		SetupWheelEffectMock: func() *dbwheeleffectsmock.DatabaseMock {
			databaseMock := new(dbwheeleffectsmock.DatabaseMock)

			databaseMock.
				On("GetLastRolledWheelEffectsCommand", 1).
				Return(typewheeleffects.RolledWheelEffects{
					{Id: 42, Name: "test-effect", IsApplied: false},
				}, nil)
			databaseMock.On("MarkLastWheelEffectAppliedCommand", 1, 42).Return(nil)
			databaseMock.On("AddWheelEffectHistoryCommand", 1, 42).Return(100, nil)

			return databaseMock
		},
		SetupPointsMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetFreePointsCommand", 2).Return(0, nil)
			databaseMock.On("ChangeFreePointsCommand", 2, 0).Return(nil)
			databaseMock.On("AddFreePointHistoryCommand", 2, 1, typepoints.FreePointsChangeSourceOther, 0, 0, 0, ptr(100)).Return(nil)
			return databaseMock
		},
		SetupSysParams: func() *srvsysparamsmock.ServiceMock {
			spSvc := new(srvsysparamsmock.ServiceMock)
			spSvc.On("GetInt", typesysparams.ParamFreePointsMinimum).Return(0, nil)
			spSvc.On("GetBool", typesysparams.ParamShouldLimitFreePoints).Return(true, nil)
			return spSvc
		},
		ExpectedErrorCode: "WRONG_DESIRED_CHANGE_VALUE",
	},
	{
		// The roll change fails after that user's free points change succeeds. The error returns.
		Name:   "ChangeAvailableRolls_Error",
		UserId: 1,
		RollApply: typewheeleffects.WheelEffectRollApply{
			WheelEffectName: "test-effect",
			PointChangeByUserIds: typepoints.PointChangeByUserIds{
				{Login: "alice", UserId: 2, FreePointChange: typepoints.FreePointChange{SourceUserId: 1, ChangeSource: typepoints.FreePointsChangeSourceOther, DesiredChangeValue: 0}, AvailableRollChange: &typepoints.PointChange{ChangeSource: typepoints.FreePointsChangeSourceOther, DesiredChangeValue: 1}},
			},
		},
		SetupWheelEffectMock: func() *dbwheeleffectsmock.DatabaseMock {
			databaseMock := new(dbwheeleffectsmock.DatabaseMock)

			databaseMock.
				On("GetLastRolledWheelEffectsCommand", 1).
				Return(typewheeleffects.RolledWheelEffects{
					{Id: 42, Name: "test-effect", IsApplied: false},
				}, nil)
			databaseMock.On("MarkLastWheelEffectAppliedCommand", 1, 42).Return(nil)
			databaseMock.On("AddWheelEffectHistoryCommand", 1, 42).Return(100, nil)
			databaseMock.On("GetAvailableRollsCountCommand", 2).Return(3, nil)

			return databaseMock
		},
		SetupPointsMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetFreePointsCommand", 2).Return(0, nil)
			databaseMock.On("ChangeFreePointsCommand", 2, 0).Return(nil)
			databaseMock.On("AddFreePointHistoryCommand", 2, 1, typepoints.FreePointsChangeSourceOther, 0, 0, 0, ptr(100)).Return(nil)
			databaseMock.On("ChangeAvailableRollsCommand", 2, 1).Return(dbError)
			return databaseMock
		},
		SetupSysParams: func() *srvsysparamsmock.ServiceMock {
			spSvc := new(srvsysparamsmock.ServiceMock)
			spSvc.On("GetInt", typesysparams.ParamFreePointsMinimum).Return(0, nil)
			spSvc.On("GetBool", typesysparams.ParamShouldLimitFreePoints).Return(true, nil)
			return spSvc
		},
		ExpectedError: dbError,
	},
}

func TestSrvWheelEffects_ApplyWheelEffectRoll(test *testing.T) {
	for _, testCase := range ApplyWheelEffectRollTestCases {
		test.Run(testCase.Name, func(test *testing.T) {
			// Arrange
			weDatabaseMock := testCase.SetupWheelEffectMock()
			pointsDatabaseMock := testCase.SetupPointsMock()

			var spSvc *srvsysparamsmock.ServiceMock
			if testCase.SetupSysParams != nil {
				spSvc = testCase.SetupSysParams()
			} else {
				spSvc = new(srvsysparamsmock.ServiceMock)
			}

			sut := srvwheeleffects.Service{
				Database:     weDatabaseMock,
				PointService: srvpoints.Service{Database: pointsDatabaseMock, SysParamsService: spSvc},
			}

			// Act
			results, err := sut.ApplyWheelEffectRoll(testCase.UserId, testCase.RollApply)

			// Assert
			if testCase.ExpectedErrorCode != "" {
				require.Error(test, err)
				appErr, ok := err.(common.AppError)
				require.True(test, ok)
				require.Equal(test, testCase.ExpectedErrorCode, appErr.GetCode())
			} else if testCase.ExpectedError != nil {
				require.ErrorIs(test, err, testCase.ExpectedError)
			} else {
				require.NoError(test, err)
				require.Equal(test, testCase.ExpectedResults, results)
			}

			weDatabaseMock.AssertExpectations(test)
			pointsDatabaseMock.AssertExpectations(test)
			spSvc.AssertExpectations(test)
		})
	}
}
