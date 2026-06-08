package srvtimers_test

import (
	"FGG-Service/src/common"
	"FGG-Service/src/games/types"
	"FGG-Service/src/sysparams/types"
	"FGG-Service/src/timers/service"
	"FGG-Service/src/timers/types"
	"FGG-Service/tests/games/mock"
	"FGG-Service/tests/sysparams/srvmock"
	"FGG-Service/tests/timers/mock/dbtimers"
	"FGG-Service/tests/timers/mock/dbwheeleffects"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type GetOrCreateCurrentTimerTestCase struct {
	Name            string
	UserId          int
	SetupMocks      func() (*dbtimermock.DatabaseMock, *dbgamesmock.DatabaseMock, *dbwheeleffectsmock.DatabaseMock, *srvsysparamsmock.ServiceMock)
	ExpectedTimer   *typetimers.Timer
	ExpectedErrorAs interface{}
	ExpectedErrorIs error
}

var existingTimer = typetimers.Timer{
	Id:             1,
	Duration:       2 * time.Hour,
	RemainingTime:  90 * time.Minute,
	State:          typetimers.TimerStateRunning,
	LastActionDate: time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC),
}

var newTimer = typetimers.Timer{
	Id:             2,
	Duration:       2 * time.Hour,
	RemainingTime:  2 * time.Hour,
	State:          typetimers.TimerStateCreated,
	LastActionDate: time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC),
}

var currentGame = typegames.CurrentGames{
	{Id: 1, Name: "Half-Life 1", State: typegames.GameStateStarted},
}

var GetOrCreateCurrentTimerTestCases = []GetOrCreateCurrentTimerTestCase{
	{
		// GetCurrentGameCommand returns sql.ErrNoRows. The CurrentGameNotFoundError will return.
		Name:   "NotFound_NoRows",
		UserId: 1,
		SetupMocks: func() (*dbtimermock.DatabaseMock, *dbgamesmock.DatabaseMock, *dbwheeleffectsmock.DatabaseMock, *srvsysparamsmock.ServiceMock) {
			timerDb := new(dbtimermock.DatabaseMock)
			gamesDb := new(dbgamesmock.DatabaseMock)
			wheelDb := new(dbwheeleffectsmock.DatabaseMock)
			spSvc := new(srvsysparamsmock.ServiceMock)

			gamesDb.On("GetCurrentGameCommand", 1).Return(typegames.CurrentGames{}, sql.ErrNoRows)

			return timerDb, gamesDb, wheelDb, spSvc
		},
		ExpectedErrorAs: new(common.NotFoundError),
	},
	{
		// GetCurrentGameCommand returns an empty list. The CurrentGameNotFoundError will return.
		Name:   "NotFound_EmptyList",
		UserId: 1,
		SetupMocks: func() (*dbtimermock.DatabaseMock, *dbgamesmock.DatabaseMock, *dbwheeleffectsmock.DatabaseMock, *srvsysparamsmock.ServiceMock) {
			timerDb := new(dbtimermock.DatabaseMock)
			gamesDb := new(dbgamesmock.DatabaseMock)
			wheelDb := new(dbwheeleffectsmock.DatabaseMock)
			spSvc := new(srvsysparamsmock.ServiceMock)

			gamesDb.On("GetCurrentGameCommand", 1).Return(typegames.CurrentGames{}, nil)

			return timerDb, gamesDb, wheelDb, spSvc
		},
		ExpectedErrorAs: new(common.NotFoundError),
	},
	{
		// GetCurrentGameCommand returns a non-empty list alongside a database error. The error will return.
		Name:   "GamesDatabaseError",
		UserId: 1,
		SetupMocks: func() (*dbtimermock.DatabaseMock, *dbgamesmock.DatabaseMock, *dbwheeleffectsmock.DatabaseMock, *srvsysparamsmock.ServiceMock) {
			timerDb := new(dbtimermock.DatabaseMock)
			gamesDb := new(dbgamesmock.DatabaseMock)
			wheelDb := new(dbwheeleffectsmock.DatabaseMock)
			spSvc := new(srvsysparamsmock.ServiceMock)

			gamesDb.On("GetCurrentGameCommand", 1).Return(currentGame, dbError)

			return timerDb, gamesDb, wheelDb, spSvc
		},
		ExpectedErrorIs: dbError,
	},
	{
		// GetCurrentGameCommand succeeds. GetCurrentTimerCommand returns an existing timer. The existing timer will return.
		Name:   "ExistingTimer",
		UserId: 1,
		SetupMocks: func() (*dbtimermock.DatabaseMock, *dbgamesmock.DatabaseMock, *dbwheeleffectsmock.DatabaseMock, *srvsysparamsmock.ServiceMock) {
			timerDb := new(dbtimermock.DatabaseMock)
			gamesDb := new(dbgamesmock.DatabaseMock)
			wheelDb := new(dbwheeleffectsmock.DatabaseMock)
			spSvc := new(srvsysparamsmock.ServiceMock)

			gamesDb.On("GetCurrentGameCommand", 1).Return(currentGame, nil)
			timerDb.On("GetCurrentTimerCommand", 1).Return(existingTimer, nil)

			return timerDb, gamesDb, wheelDb, spSvc
		},
		ExpectedTimer: &existingTimer,
	},
	{
		// GetCurrentGameCommand succeeds. GetCurrentTimerCommand returns a database error. The error will return.
		Name:   "TimerDatabaseError",
		UserId: 1,
		SetupMocks: func() (*dbtimermock.DatabaseMock, *dbgamesmock.DatabaseMock, *dbwheeleffectsmock.DatabaseMock, *srvsysparamsmock.ServiceMock) {
			timerDb := new(dbtimermock.DatabaseMock)
			gamesDb := new(dbgamesmock.DatabaseMock)
			wheelDb := new(dbwheeleffectsmock.DatabaseMock)
			spSvc := new(srvsysparamsmock.ServiceMock)

			gamesDb.On("GetCurrentGameCommand", 1).Return(currentGame, nil)
			timerDb.On("GetCurrentTimerCommand", 1).Return(typetimers.Timer{}, dbError)

			return timerDb, gamesDb, wheelDb, spSvc
		},
		ExpectedErrorIs: dbError,
	},
	{
		// GetCurrentGameCommand and GetCurrentTimerCommand (no rows) succeed.
		// GetAvailableRollsCountCommand returns a database error. The error will return.
		Name:   "WheelEffectsDatabaseError",
		UserId: 1,
		SetupMocks: func() (*dbtimermock.DatabaseMock, *dbgamesmock.DatabaseMock, *dbwheeleffectsmock.DatabaseMock, *srvsysparamsmock.ServiceMock) {
			timerDb := new(dbtimermock.DatabaseMock)
			gamesDb := new(dbgamesmock.DatabaseMock)
			wheelDb := new(dbwheeleffectsmock.DatabaseMock)
			spSvc := new(srvsysparamsmock.ServiceMock)

			gamesDb.On("GetCurrentGameCommand", 1).Return(currentGame, nil)
			timerDb.On("GetCurrentTimerCommand", 1).Return(typetimers.Timer{}, sql.ErrNoRows)
			wheelDb.On("GetAvailableRollsCountCommand", 1).Return(0, dbError)

			return timerDb, gamesDb, wheelDb, spSvc
		},
		ExpectedErrorIs: dbError,
	},
	{
		// GetCurrentGameCommand, GetCurrentTimerCommand (no rows), and GetAvailableRollsCountCommand succeed.
		// GetInt("DefaultTimerDurationInS") returns a database error. The error will return.
		Name:   "SysParams_DatabaseError",
		UserId: 1,
		SetupMocks: func() (*dbtimermock.DatabaseMock, *dbgamesmock.DatabaseMock, *dbwheeleffectsmock.DatabaseMock, *srvsysparamsmock.ServiceMock) {
			timerDb := new(dbtimermock.DatabaseMock)
			gamesDb := new(dbgamesmock.DatabaseMock)
			wheelDb := new(dbwheeleffectsmock.DatabaseMock)
			spSvc := new(srvsysparamsmock.ServiceMock)

			gamesDb.On("GetCurrentGameCommand", 1).Return(currentGame, nil)
			timerDb.On("GetCurrentTimerCommand", 1).Return(typetimers.Timer{}, sql.ErrNoRows)
			wheelDb.On("GetAvailableRollsCountCommand", 1).Return(0, nil)
			spSvc.On("GetInt", typesysparams.ParamDefaultTimerDurationInS).Return(0, dbError)

			return timerDb, gamesDb, wheelDb, spSvc
		},
		ExpectedErrorIs: dbError,
	},
	{
		// GetCurrentGameCommand, GetCurrentTimerCommand (no rows), and GetAvailableRollsCountCommand succeed.
		// CreateCurrentTimerCommand returns a database error. The error will return.
		Name:   "CreateTimerDatabaseError",
		UserId: 1,
		SetupMocks: func() (*dbtimermock.DatabaseMock, *dbgamesmock.DatabaseMock, *dbwheeleffectsmock.DatabaseMock, *srvsysparamsmock.ServiceMock) {
			timerDb := new(dbtimermock.DatabaseMock)
			gamesDb := new(dbgamesmock.DatabaseMock)
			wheelDb := new(dbwheeleffectsmock.DatabaseMock)
			spSvc := new(srvsysparamsmock.ServiceMock)

			gamesDb.On("GetCurrentGameCommand", 1).Return(currentGame, nil)
			timerDb.On("GetCurrentTimerCommand", 1).Return(typetimers.Timer{}, sql.ErrNoRows)
			wheelDb.On("GetAvailableRollsCountCommand", 1).Return(0, nil)
			spSvc.On("GetInt", typesysparams.ParamDefaultTimerDurationInS).Return(30, nil)
			timerDb.On("CreateCurrentTimerCommand", 1, 1, 30).Return(dbError)

			return timerDb, gamesDb, wheelDb, spSvc
		},
		ExpectedErrorIs: dbError,
	},
	{
		// All commands succeed. A new timer is created and returned.
		Name:   "NewTimerCreated",
		UserId: 1,
		SetupMocks: func() (*dbtimermock.DatabaseMock, *dbgamesmock.DatabaseMock, *dbwheeleffectsmock.DatabaseMock, *srvsysparamsmock.ServiceMock) {
			timerDb := new(dbtimermock.DatabaseMock)
			gamesDb := new(dbgamesmock.DatabaseMock)
			wheelDb := new(dbwheeleffectsmock.DatabaseMock)
			spSvc := new(srvsysparamsmock.ServiceMock)

			gamesDb.On("GetCurrentGameCommand", 1).Return(currentGame, nil)
			timerDb.On("GetCurrentTimerCommand", 1).Once().Return(typetimers.Timer{}, sql.ErrNoRows)
			wheelDb.On("GetAvailableRollsCountCommand", 1).Return(0, nil)
			spSvc.On("GetInt", typesysparams.ParamDefaultTimerDurationInS).Return(30, nil)
			timerDb.On("CreateCurrentTimerCommand", 1, 1, 30).Return(nil)
			timerDb.On("GetCurrentTimerCommand", 1).Once().Return(newTimer, nil)

			return timerDb, gamesDb, wheelDb, spSvc
		},
		ExpectedTimer: &newTimer,
	},
}

func TestSrvTimers_GetOrCreateCurrentTimer(test *testing.T) {
	for _, testCase := range GetOrCreateCurrentTimerTestCases {
		test.Run(testCase.Name, func(test *testing.T) {
			// Arrange
			timerDb, gamesDb, wheelDb, spSvc := testCase.SetupMocks()
			sut := srvtimers.Service{
				Database:             timerDb,
				GamesDatabase:        gamesDb,
				WheelEffectsDatabase: wheelDb,
				SysParamsService:     spSvc,
			}

			// Act
			timer, err := sut.GetOrCreateCurrentTimer(testCase.UserId)

			// Assert
			if testCase.ExpectedErrorAs != nil {
				require.Error(test, err)
				require.ErrorAs(test, err, &testCase.ExpectedErrorAs)
			}

			if testCase.ExpectedErrorIs != nil {
				require.ErrorIs(test, err, testCase.ExpectedErrorIs)
			}

			if testCase.ExpectedTimer != nil {
				require.NoError(test, err)
				require.Equal(test, *testCase.ExpectedTimer, timer)
			}

			timerDb.AssertExpectations(test)
			gamesDb.AssertExpectations(test)
			wheelDb.AssertExpectations(test)
			spSvc.AssertExpectations(test)
		})
	}
}
