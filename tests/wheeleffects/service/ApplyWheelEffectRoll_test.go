package srvwheeleffects_test

import (
	"FGG-Service/src/common"
	"FGG-Service/src/points/service"
	"FGG-Service/src/points/type"
	"FGG-Service/src/wheeleffects/service"
	"FGG-Service/src/wheeleffects/types"
	"FGG-Service/tests/points/mock"
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
			PointChangeByUserIds: typepoints.FreePointChangeByUserIds{
				{Login: "alice", UserId: 2, PointChange: typepoints.FreePointChange{SourceUserId: 1, ChangeSource: typepoints.FreePointsChangeSourceQuestCompletion, DesiredChangeValue: 5}},
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
		// ChangeFreePoints fails for a change. The error returns and the roll isn't marked as applied.
		Name:   "ChangeFreePoints_Error",
		UserId: 1,
		RollApply: typewheeleffects.WheelEffectRollApply{
			WheelEffectName: "test-effect",
			PointChangeByUserIds: typepoints.FreePointChangeByUserIds{
				{Login: "alice", UserId: 2, PointChange: typepoints.FreePointChange{SourceUserId: 1, ChangeSource: typepoints.FreePointsChangeSourceQuestCompletion, DesiredChangeValue: 5}},
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
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetFreePointsCommand", 2).Return(0, dbError)
			return databaseMock
		},
		ExpectedError: dbError,
	},
	{
		// All point changes succeed but marking the roll as applied fails. The error returns and the history isn't recorded.
		Name:   "MarkEffectRollAsApplied_Error",
		UserId: 1,
		RollApply: typewheeleffects.WheelEffectRollApply{
			WheelEffectName: "test-effect",
			PointChangeByUserIds: typepoints.FreePointChangeByUserIds{
				{Login: "alice", UserId: 2, PointChange: typepoints.FreePointChange{SourceUserId: 1, ChangeSource: typepoints.FreePointsChangeSourceQuestCompletion, DesiredChangeValue: 5}},
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
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetFreePointsCommand", 2).Return(10, nil)
			databaseMock.On("ChangeFreePointsCommand", 2, 5).Return(nil)
			databaseMock.On("AddFreePointHistoryCommand", 2, 1, typepoints.FreePointsChangeSourceQuestCompletion, 5, 5, 15, ptr(42)).Return(nil)
			return databaseMock
		},
		ExpectedError: dbError,
	},
	{
		// Marking the roll as applied succeeds but recording the history fails. The error returns.
		Name:   "AddWheelEffectHistory_Error",
		UserId: 1,
		RollApply: typewheeleffects.WheelEffectRollApply{
			WheelEffectName: "test-effect",
			PointChangeByUserIds: typepoints.FreePointChangeByUserIds{
				{Login: "alice", UserId: 2, PointChange: typepoints.FreePointChange{SourceUserId: 1, ChangeSource: typepoints.FreePointsChangeSourceQuestCompletion, DesiredChangeValue: 5}},
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
			databaseMock.On("AddWheelEffectHistoryCommand", 1, 42).Return(dbError)

			return databaseMock
		},
		SetupPointsMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetFreePointsCommand", 2).Return(10, nil)
			databaseMock.On("ChangeFreePointsCommand", 2, 5).Return(nil)
			databaseMock.On("AddFreePointHistoryCommand", 2, 1, typepoints.FreePointsChangeSourceQuestCompletion, 5, 5, 15, ptr(42)).Return(nil)
			return databaseMock
		},
		ExpectedError: dbError,
	},
	{
		// Every change applies successfully and the roll is marked as applied with its history recorded. Results return in request order.
		Name:   "Success",
		UserId: 1,
		RollApply: typewheeleffects.WheelEffectRollApply{
			WheelEffectName: "test-effect",
			PointChangeByUserIds: typepoints.FreePointChangeByUserIds{
				{Login: "alice", UserId: 2, PointChange: typepoints.FreePointChange{SourceUserId: 1, ChangeSource: typepoints.FreePointsChangeSourceQuestCompletion, DesiredChangeValue: 5}},
				{Login: "bob", UserId: 3, PointChange: typepoints.FreePointChange{SourceUserId: 1, ChangeSource: typepoints.FreePointsChangeSourceOther, DesiredChangeValue: -2}},
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
			databaseMock.On("AddWheelEffectHistoryCommand", 1, 42).Return(nil)

			return databaseMock
		},
		SetupPointsMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetFreePointsCommand", 2).Return(10, nil)
			databaseMock.On("ChangeFreePointsCommand", 2, 5).Return(nil)
			databaseMock.On("AddFreePointHistoryCommand", 2, 1, typepoints.FreePointsChangeSourceQuestCompletion, 5, 5, 15, ptr(42)).Return(nil)
			databaseMock.On("GetFreePointsCommand", 3).Return(4, nil)
			databaseMock.On("ChangeFreePointsCommand", 3, -2).Return(nil)
			databaseMock.On("AddFreePointHistoryCommand", 3, 1, typepoints.FreePointsChangeSourceOther, -2, -2, 2, ptr(42)).Return(nil)
			return databaseMock
		},
		ExpectedResults: typepoints.PointChangeResultByUserIds{
			{Login: "alice", UserId: 2, ChangeResult: typepoints.PointChangeResult{ActualChangeValue: 5, ChangeSource: typepoints.FreePointsChangeSourceQuestCompletion, DesiredChangeValue: 5, FinalValue: 15}},
			{Login: "bob", UserId: 3, ChangeResult: typepoints.PointChangeResult{ActualChangeValue: -2, ChangeSource: typepoints.FreePointsChangeSourceOther, DesiredChangeValue: -2, FinalValue: 2}},
		},
	},
}

func TestSrvWheelEffects_ApplyWheelEffectRoll(test *testing.T) {
	for _, testCase := range ApplyWheelEffectRollTestCases {
		test.Run(testCase.Name, func(test *testing.T) {
			// Arrange
			weDatabaseMock := testCase.SetupWheelEffectMock()
			pointsDatabaseMock := testCase.SetupPointsMock()
			sut := srvwheeleffects.Service{
				Database:     weDatabaseMock,
				PointService: srvpoints.Service{Database: pointsDatabaseMock},
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
		})
	}
}
