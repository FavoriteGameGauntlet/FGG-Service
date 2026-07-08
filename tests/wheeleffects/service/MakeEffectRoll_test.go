package srvwheeleffects_test

import (
	"FGG-Service/src/common"
	srvpoints "FGG-Service/src/points/service"
	"FGG-Service/src/sysparams/types"
	srvwheeleffects "FGG-Service/src/wheeleffects/service"
	typewheeleffects "FGG-Service/src/wheeleffects/types"
	dbpointsmock "FGG-Service/tests/points/mock"
	srvsysparamsmock "FGG-Service/tests/sysparams/srvmock"
	dbwheeleffectsmock "FGG-Service/tests/timers/mock/dbwheeleffects"
	"testing"

	"github.com/stretchr/testify/require"
)

var rolledEffects = typewheeleffects.WheelEffects{
	{Id: 10, Name: "effect-a", Description: ptr("desc a")},
	{Id: 11, Name: "effect-b", Description: ptr("desc b")},
	{Id: 12, Name: "effect-c", Description: ptr("desc c")},
}

type MakeEffectRollTestCase struct {
	Name                 string
	UserId               int
	IsReroll             bool
	SetupWheelEffectMock func() *dbwheeleffectsmock.DatabaseMock
	SetupPointsMock      func() *dbpointsmock.DatabaseMock
	SetupSysParams       func() *srvsysparamsmock.ServiceMock
	ExpectedEffects      typewheeleffects.WheelEffects
	ExpectedError        error
	ExpectedErrorCode    string
}

var MakeEffectRollTestCases = []MakeEffectRollTestCase{
	{
		// GetAvailableRollsCountCommand fails. The error returns and nothing else happens.
		Name:   "GetAvailableRollsCount_Error",
		UserId: 1,
		SetupWheelEffectMock: func() *dbwheeleffectsmock.DatabaseMock {
			databaseMock := new(dbwheeleffectsmock.DatabaseMock)
			databaseMock.On("GetAvailableRollsCountCommand", 1).Return(0, dbError)
			return databaseMock
		},
		SetupPointsMock: func() *dbpointsmock.DatabaseMock { return new(dbpointsmock.DatabaseMock) },
		SetupSysParams:  func() *srvsysparamsmock.ServiceMock { return new(srvsysparamsmock.ServiceMock) },
		ExpectedError:   dbError,
	},
	{
		// GetAvailableRollsCountCommand succeeds but the sys param lookup fails.
		Name:   "GetMinimumRollCount_SysParam_Error",
		UserId: 1,
		SetupWheelEffectMock: func() *dbwheeleffectsmock.DatabaseMock {
			databaseMock := new(dbwheeleffectsmock.DatabaseMock)
			databaseMock.On("GetAvailableRollsCountCommand", 1).Return(5, nil)
			return databaseMock
		},
		SetupPointsMock: func() *dbpointsmock.DatabaseMock { return new(dbpointsmock.DatabaseMock) },
		SetupSysParams: func() *srvsysparamsmock.ServiceMock {
			spSvc := new(srvsysparamsmock.ServiceMock)
			spSvc.On("GetInt", typesysparams.ParamMinimumAvailableRollCountForRoll).Return(0, dbError)
			return spSvc
		},
		ExpectedError: dbError,
	},
	{
		// The user has no available rolls. The AVAILABLE_ROLLS_NOT_FOUND error returns.
		Name:   "NotEnoughRolls_Error",
		UserId: 1,
		SetupWheelEffectMock: func() *dbwheeleffectsmock.DatabaseMock {
			databaseMock := new(dbwheeleffectsmock.DatabaseMock)
			databaseMock.On("GetAvailableRollsCountCommand", 1).Return(0, nil)
			return databaseMock
		},
		SetupPointsMock: func() *dbpointsmock.DatabaseMock { return new(dbpointsmock.DatabaseMock) },
		SetupSysParams: func() *srvsysparamsmock.ServiceMock {
			spSvc := new(srvsysparamsmock.ServiceMock)
			spSvc.On("GetInt", typesysparams.ParamMinimumAvailableRollCountForRoll).Return(1, nil)
			return spSvc
		},
		ExpectedErrorCode: "AVAILABLE_ROLLS_NOT_FOUND",
	},
	{
		// MakeEffectRollCommand fails. The error returns.
		Name:   "MakeEffectRoll_Error",
		UserId: 1,
		SetupWheelEffectMock: func() *dbwheeleffectsmock.DatabaseMock {
			databaseMock := new(dbwheeleffectsmock.DatabaseMock)
			databaseMock.On("GetAvailableRollsCountCommand", 1).Return(1, nil)
			databaseMock.On("MakeEffectRollCommand", 1).Return(typewheeleffects.WheelEffects(nil), dbError)
			return databaseMock
		},
		SetupPointsMock: func() *dbpointsmock.DatabaseMock { return new(dbpointsmock.DatabaseMock) },
		SetupSysParams: func() *srvsysparamsmock.ServiceMock {
			spSvc := new(srvsysparamsmock.ServiceMock)
			spSvc.On("GetInt", typesysparams.ParamMinimumAvailableRollCountForRoll).Return(1, nil)
			return spSvc
		},
		ExpectedError: dbError,
	},
	{
		// MakeEffectRollCommand succeeds but the minimum effects count sys param lookup fails.
		Name:   "GetMinimumEffectsCount_SysParam_Error",
		UserId: 1,
		SetupWheelEffectMock: func() *dbwheeleffectsmock.DatabaseMock {
			databaseMock := new(dbwheeleffectsmock.DatabaseMock)
			databaseMock.On("GetAvailableRollsCountCommand", 1).Return(1, nil)
			databaseMock.On("MakeEffectRollCommand", 1).Return(rolledEffects, nil)
			return databaseMock
		},
		SetupPointsMock: func() *dbpointsmock.DatabaseMock { return new(dbpointsmock.DatabaseMock) },
		SetupSysParams: func() *srvsysparamsmock.ServiceMock {
			spSvc := new(srvsysparamsmock.ServiceMock)
			spSvc.On("GetInt", typesysparams.ParamMinimumAvailableRollCountForRoll).Return(1, nil)
			spSvc.On("GetInt", typesysparams.ParamMinimumAvailableWheelEffectsForRoll).Return(0, dbError)
			return spSvc
		},
		ExpectedError: dbError,
	},
	{
		// The roll returned fewer effects than the configured minimum. The NOT_ENOUGH_AVAILABLE_WHEEL_EFFECTS error returns.
		Name:   "NotEnoughEffects_Error",
		UserId: 1,
		SetupWheelEffectMock: func() *dbwheeleffectsmock.DatabaseMock {
			databaseMock := new(dbwheeleffectsmock.DatabaseMock)
			databaseMock.On("GetAvailableRollsCountCommand", 1).Return(1, nil)
			databaseMock.On("MakeEffectRollCommand", 1).Return(rolledEffects[:1], nil)
			return databaseMock
		},
		SetupPointsMock: func() *dbpointsmock.DatabaseMock { return new(dbpointsmock.DatabaseMock) },
		SetupSysParams: func() *srvsysparamsmock.ServiceMock {
			spSvc := new(srvsysparamsmock.ServiceMock)
			spSvc.On("GetInt", typesysparams.ParamMinimumAvailableRollCountForRoll).Return(1, nil)
			spSvc.On("GetInt", typesysparams.ParamMinimumAvailableWheelEffectsForRoll).Return(3, nil)
			return spSvc
		},
		ExpectedErrorCode: "NOT_ENOUGH_AVAILABLE_WHEEL_EFFECTS",
	},
	{
		// The roll change sys param lookup fails.
		Name:   "GetRollChange_SysParam_Error",
		UserId: 1,
		SetupWheelEffectMock: func() *dbwheeleffectsmock.DatabaseMock {
			databaseMock := new(dbwheeleffectsmock.DatabaseMock)
			databaseMock.On("GetAvailableRollsCountCommand", 1).Return(1, nil)
			databaseMock.On("MakeEffectRollCommand", 1).Return(rolledEffects, nil)
			return databaseMock
		},
		SetupPointsMock: func() *dbpointsmock.DatabaseMock { return new(dbpointsmock.DatabaseMock) },
		SetupSysParams: func() *srvsysparamsmock.ServiceMock {
			spSvc := new(srvsysparamsmock.ServiceMock)
			spSvc.On("GetInt", typesysparams.ParamMinimumAvailableRollCountForRoll).Return(1, nil)
			spSvc.On("GetInt", typesysparams.ParamMinimumAvailableWheelEffectsForRoll).Return(3, nil)
			spSvc.On("GetInt", typesysparams.ParamAvailableRollChangeByRoll).Return(0, dbError)
			return spSvc
		},
		ExpectedError: dbError,
	},
	{
		// ChangeAvailableRolls fails. The error returns and the last effects list is not modified.
		Name:   "ChangeAvailableRolls_Error",
		UserId: 1,
		SetupWheelEffectMock: func() *dbwheeleffectsmock.DatabaseMock {
			databaseMock := new(dbwheeleffectsmock.DatabaseMock)
			databaseMock.On("GetAvailableRollsCountCommand", 1).Return(1, nil)
			databaseMock.On("MakeEffectRollCommand", 1).Return(rolledEffects, nil)
			return databaseMock
		},
		SetupPointsMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("ChangeAvailableRollsCommand", 1, -1).Return(dbError)
			return databaseMock
		},
		SetupSysParams: func() *srvsysparamsmock.ServiceMock {
			spSvc := new(srvsysparamsmock.ServiceMock)
			spSvc.On("GetInt", typesysparams.ParamMinimumAvailableRollCountForRoll).Return(1, nil)
			spSvc.On("GetInt", typesysparams.ParamMinimumAvailableWheelEffectsForRoll).Return(3, nil)
			spSvc.On("GetInt", typesysparams.ParamAvailableRollChangeByRoll).Return(-1, nil)
			return spSvc
		},
		ExpectedError: dbError,
	},
	{
		// ClearLastWheelEffectsCommand fails. The error returns and the new effects are not written.
		Name:   "ClearLastWheelEffects_Error",
		UserId: 1,
		SetupWheelEffectMock: func() *dbwheeleffectsmock.DatabaseMock {
			databaseMock := new(dbwheeleffectsmock.DatabaseMock)
			databaseMock.On("GetAvailableRollsCountCommand", 1).Return(1, nil)
			databaseMock.On("MakeEffectRollCommand", 1).Return(rolledEffects, nil)
			databaseMock.On("ClearLastWheelEffectsCommand", 1).Return(dbError)
			return databaseMock
		},
		SetupPointsMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("ChangeAvailableRollsCommand", 1, -1).Return(nil)
			return databaseMock
		},
		SetupSysParams: func() *srvsysparamsmock.ServiceMock {
			spSvc := new(srvsysparamsmock.ServiceMock)
			spSvc.On("GetInt", typesysparams.ParamMinimumAvailableRollCountForRoll).Return(1, nil)
			spSvc.On("GetInt", typesysparams.ParamMinimumAvailableWheelEffectsForRoll).Return(3, nil)
			spSvc.On("GetInt", typesysparams.ParamAvailableRollChangeByRoll).Return(-1, nil)
			return spSvc
		},
		ExpectedError: dbError,
	},
	{
		// ClearLastWheelEffectsCommand succeeds but AddLastRolledWheelEffectsCommand fails.
		Name:   "AddLastRolledWheelEffects_Error",
		UserId: 1,
		SetupWheelEffectMock: func() *dbwheeleffectsmock.DatabaseMock {
			databaseMock := new(dbwheeleffectsmock.DatabaseMock)
			databaseMock.On("GetAvailableRollsCountCommand", 1).Return(1, nil)
			databaseMock.On("MakeEffectRollCommand", 1).Return(rolledEffects, nil)
			databaseMock.On("ClearLastWheelEffectsCommand", 1).Return(nil)
			databaseMock.On("AddLastRolledWheelEffectsCommand", 1, rolledEffects).Return(dbError)
			return databaseMock
		},
		SetupPointsMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("ChangeAvailableRollsCommand", 1, -1).Return(nil)
			return databaseMock
		},
		SetupSysParams: func() *srvsysparamsmock.ServiceMock {
			spSvc := new(srvsysparamsmock.ServiceMock)
			spSvc.On("GetInt", typesysparams.ParamMinimumAvailableRollCountForRoll).Return(1, nil)
			spSvc.On("GetInt", typesysparams.ParamMinimumAvailableWheelEffectsForRoll).Return(3, nil)
			spSvc.On("GetInt", typesysparams.ParamAvailableRollChangeByRoll).Return(-1, nil)
			return spSvc
		},
		ExpectedError: dbError,
	},
	{
		// All steps succeed. The rolled effects return.
		Name:   "Success",
		UserId: 1,
		SetupWheelEffectMock: func() *dbwheeleffectsmock.DatabaseMock {
			databaseMock := new(dbwheeleffectsmock.DatabaseMock)
			databaseMock.On("GetAvailableRollsCountCommand", 1).Return(1, nil)
			databaseMock.On("MakeEffectRollCommand", 1).Return(rolledEffects, nil)
			databaseMock.On("ClearLastWheelEffectsCommand", 1).Return(nil)
			databaseMock.On("AddLastRolledWheelEffectsCommand", 1, rolledEffects).Return(nil)
			return databaseMock
		},
		SetupPointsMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("ChangeAvailableRollsCommand", 1, -1).Return(nil)
			return databaseMock
		},
		SetupSysParams: func() *srvsysparamsmock.ServiceMock {
			spSvc := new(srvsysparamsmock.ServiceMock)
			spSvc.On("GetInt", typesysparams.ParamMinimumAvailableRollCountForRoll).Return(1, nil)
			spSvc.On("GetInt", typesysparams.ParamMinimumAvailableWheelEffectsForRoll).Return(3, nil)
			spSvc.On("GetInt", typesysparams.ParamAvailableRollChangeByRoll).Return(-1, nil)
			return spSvc
		},
		ExpectedEffects: rolledEffects,
	},
	{
		// A reroll skips the available rolls check and does not decrement available rolls.
		Name:     "Reroll_Success",
		UserId:   1,
		IsReroll: true,
		SetupWheelEffectMock: func() *dbwheeleffectsmock.DatabaseMock {
			databaseMock := new(dbwheeleffectsmock.DatabaseMock)
			databaseMock.On("MakeEffectRollCommand", 1).Return(rolledEffects, nil)
			databaseMock.On("ClearLastWheelEffectsCommand", 1).Return(nil)
			databaseMock.On("AddLastRolledWheelEffectsCommand", 1, rolledEffects).Return(nil)
			return databaseMock
		},
		SetupPointsMock: func() *dbpointsmock.DatabaseMock { return new(dbpointsmock.DatabaseMock) },
		SetupSysParams: func() *srvsysparamsmock.ServiceMock {
			spSvc := new(srvsysparamsmock.ServiceMock)
			spSvc.On("GetInt", typesysparams.ParamMinimumAvailableWheelEffectsForRoll).Return(3, nil)
			return spSvc
		},
		ExpectedEffects: rolledEffects,
	},
}

func TestSrvWheelEffects_MakeEffectRoll(test *testing.T) {
	for _, testCase := range MakeEffectRollTestCases {
		test.Run(testCase.Name, func(test *testing.T) {
			// Arrange
			weDatabaseMock := testCase.SetupWheelEffectMock()
			pointsDatabaseMock := testCase.SetupPointsMock()
			spSvc := testCase.SetupSysParams()

			sut := srvwheeleffects.Service{
				Database:         weDatabaseMock,
				PointService:     srvpoints.Service{Database: pointsDatabaseMock},
				SysParamsService: spSvc,
			}

			// Act
			effects, err := sut.MakeEffectRoll(testCase.UserId, testCase.IsReroll)

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
				require.Equal(test, testCase.ExpectedEffects, effects)
			}

			weDatabaseMock.AssertExpectations(test)
			pointsDatabaseMock.AssertExpectations(test)
			spSvc.AssertExpectations(test)
		})
	}
}
