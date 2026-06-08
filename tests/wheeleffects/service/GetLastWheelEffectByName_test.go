package srvwheeleffects_test

import (
	"FGG-Service/src/common"
	"FGG-Service/src/wheeleffects/service"
	"FGG-Service/src/wheeleffects/types"
	"FGG-Service/tests/timers/mock/dbwheeleffects"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

type GetLastWheelEffectByNameTestCase struct {
	Name              string
	UserId            int
	EffectName        string
	SetupMock         func() *dbwheeleffectsmock.DatabaseMock
	ExpectedEffect    *typewheeleffects.RolledWheelEffect
	ExpectedError     error
	ExpectedErrorCode string
}

var GetLastWheelEffectByNameTestCases = []GetLastWheelEffectByNameTestCase{
	{
		// GetLastRolledWheelEffectsCommand returns sql.ErrNoRows. The LAST_WHEEL_EFFECTS_NOT_FOUND error returns.
		Name:       "NoRows_NotFoundError",
		UserId:     1,
		EffectName: "test-effect",
		SetupMock: func() *dbwheeleffectsmock.DatabaseMock {
			databaseMock := new(dbwheeleffectsmock.DatabaseMock)

			databaseMock.
				On("GetLastRolledWheelEffectsCommand", 1).
				Return(typewheeleffects.RolledWheelEffects{}, sql.ErrNoRows)

			return databaseMock
		},
		ExpectedErrorCode: "LAST_WHEEL_EFFECTS_NOT_FOUND",
	},
	{
		// GetLastRolledWheelEffectsCommand returns an empty slice. The LAST_WHEEL_EFFECTS_NOT_FOUND error returns.
		Name:       "Empty_NotFoundError",
		UserId:     1,
		EffectName: "test-effect",
		SetupMock: func() *dbwheeleffectsmock.DatabaseMock {
			databaseMock := new(dbwheeleffectsmock.DatabaseMock)

			databaseMock.
				On("GetLastRolledWheelEffectsCommand", 1).
				Return(typewheeleffects.RolledWheelEffects{}, nil)

			return databaseMock
		},
		ExpectedErrorCode: "LAST_WHEEL_EFFECTS_NOT_FOUND",
	},
	{
		// GetLastRolledWheelEffectsCommand returns a database error. The error returns as-is.
		Name:       "DatabaseError",
		UserId:     1,
		EffectName: "test-effect",
		SetupMock: func() *dbwheeleffectsmock.DatabaseMock {
			databaseMock := new(dbwheeleffectsmock.DatabaseMock)

			databaseMock.
				On("GetLastRolledWheelEffectsCommand", 1).
				Return(typewheeleffects.RolledWheelEffects{}, dbError)

			return databaseMock
		},
		ExpectedError: dbError,
	},
	{
		// No last effect matches the requested name. The WHEEL_EFFECT_NAME_NOT_FOUND error returns.
		Name:       "NameNotFound_Error",
		UserId:     1,
		EffectName: "unknown-effect",
		SetupMock: func() *dbwheeleffectsmock.DatabaseMock {
			databaseMock := new(dbwheeleffectsmock.DatabaseMock)

			databaseMock.
				On("GetLastRolledWheelEffectsCommand", 1).
				Return(typewheeleffects.RolledWheelEffects{
					{Id: 42, Name: "test-effect", IsApplied: false},
				}, nil)

			return databaseMock
		},
		ExpectedErrorCode: "WHEEL_EFFECT_NAME_NOT_FOUND",
	},
	{
		// The matching last effect has already been applied. The WHEEL_EFFECT_ROLL_ALREADY_APPLIED error returns.
		Name:       "AlreadyApplied_ConflictError",
		UserId:     1,
		EffectName: "test-effect",
		SetupMock: func() *dbwheeleffectsmock.DatabaseMock {
			databaseMock := new(dbwheeleffectsmock.DatabaseMock)

			databaseMock.
				On("GetLastRolledWheelEffectsCommand", 1).
				Return(typewheeleffects.RolledWheelEffects{
					{Id: 42, Name: "test-effect", IsApplied: true},
				}, nil)

			return databaseMock
		},
		ExpectedErrorCode: "WHEEL_EFFECT_ROLL_ALREADY_APPLIED",
	},
	{
		// The matching last effect hasn't been applied yet. Success — the effect returns.
		Name:       "Success",
		UserId:     1,
		EffectName: "test-effect",
		SetupMock: func() *dbwheeleffectsmock.DatabaseMock {
			databaseMock := new(dbwheeleffectsmock.DatabaseMock)

			databaseMock.
				On("GetLastRolledWheelEffectsCommand", 1).
				Return(typewheeleffects.RolledWheelEffects{
					{Id: 42, Name: "test-effect", IsApplied: false},
				}, nil)

			return databaseMock
		},
		ExpectedEffect: &typewheeleffects.RolledWheelEffect{Id: 42, Name: "test-effect", IsApplied: false},
	},
}

func TestSrvWheelEffects_GetLastWheelEffectByName(test *testing.T) {
	for _, testCase := range GetLastWheelEffectByNameTestCases {
		test.Run(testCase.Name, func(test *testing.T) {
			// Arrange
			databaseMock := testCase.SetupMock()
			sut := srvwheeleffects.Service{Database: databaseMock}

			// Act
			effect, err := sut.GetLastWheelEffectByName(testCase.UserId, testCase.EffectName)

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
				require.Equal(test, *testCase.ExpectedEffect, effect)
			}

			databaseMock.AssertExpectations(test)
		})
	}
}
