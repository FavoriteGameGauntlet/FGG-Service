package srvpoints_test

import (
	"FGG-Service/src/common"
	"FGG-Service/src/points/service"
	typepoints "FGG-Service/src/points/type"
	"FGG-Service/tests/points/mock"
	"testing"

	"github.com/stretchr/testify/require"
)

type ChangeTerritoryPointsTestCase struct {
	Name                 string
	UserId               int
	Change               typepoints.TerritoryPointChange
	SetupMock            func() *dbpointsmock.DatabaseMock
	ExpectedActualChange *int
	ExpectedFinalValue   *int
	ExpectedErrorCode    string
}

var ChangeTerritoryPointsTestCases = []ChangeTerritoryPointsTestCase{
	{
		// Change source is not "territory-obtaining", "territory-loss" or "other". Unprocessable error returns.
		Name:   "InvalidSource_UnprocessableError",
		UserId: 1,
		Change: typepoints.TerritoryPointChange{SourceUserId: 2, ChangeSource: "invalid", DesiredChangeValue: 2},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			return new(dbpointsmock.DatabaseMock)
		},
		ExpectedErrorCode: "INCORRECT_CHANGE_SOURCE_VALUE",
	},
	{
		// Source is "territory-obtaining" with negative desired value. Conflict error returns.
		Name:   "Obtaining_NegativeValue_ConflictError",
		UserId: 1,
		Change: typepoints.TerritoryPointChange{SourceUserId: 2, ChangeSource: typepoints.TerritoryPointChangeSourceObtaining, DesiredChangeValue: -2},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			return new(dbpointsmock.DatabaseMock)
		},
		ExpectedErrorCode: "WRONG_DESIRED_CHANGE_VALUE",
	},
	{
		// Source is "territory-loss" with positive desired value. Conflict error returns.
		Name:   "Loss_PositiveValue_ConflictError",
		UserId: 1,
		Change: typepoints.TerritoryPointChange{SourceUserId: 2, ChangeSource: typepoints.TerritoryPointChangeSourceLoss, DesiredChangeValue: 2},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			return new(dbpointsmock.DatabaseMock)
		},
		ExpectedErrorCode: "WRONG_DESIRED_CHANGE_VALUE",
	},
	{
		// GetTerritoryPointsCommand returns a database error.
		Name:   "DatabaseError",
		UserId: 1,
		Change: typepoints.TerritoryPointChange{SourceUserId: 2, ChangeSource: typepoints.TerritoryPointChangeSourceObtaining, DesiredChangeValue: 2},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetTerritoryPointsCommand", 1).Return(0, dbError)
			return databaseMock
		},
		ExpectedErrorCode: "",
	},
	{
		// Source is "territory-obtaining" with zero desired value. Success.
		Name:   "Obtaining_ZeroValue_Success",
		UserId: 1,
		Change: typepoints.TerritoryPointChange{SourceUserId: 2, ChangeSource: typepoints.TerritoryPointChangeSourceObtaining, DesiredChangeValue: 0},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetTerritoryPointsCommand", 1).Return(10, nil)
			databaseMock.On("ChangeTerritoryPointsCommand", 1, 0).Return(nil)
			databaseMock.On("AddTerritoryPointHistoryCommand", 1, 2, typepoints.TerritoryPointChangeSourceObtaining, 0, 0, 10).Return(nil)
			return databaseMock
		},
		ExpectedActualChange: ptr(0),
		ExpectedFinalValue:   ptr(10),
	},
	{
		// Source is "territory-obtaining" with positive desired value. Success.
		Name:   "Obtaining_PositiveValue_Success",
		UserId: 1,
		Change: typepoints.TerritoryPointChange{SourceUserId: 2, ChangeSource: typepoints.TerritoryPointChangeSourceObtaining, DesiredChangeValue: 5},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetTerritoryPointsCommand", 1).Return(10, nil)
			databaseMock.On("ChangeTerritoryPointsCommand", 1, 5).Return(nil)
			databaseMock.On("AddTerritoryPointHistoryCommand", 1, 2, typepoints.TerritoryPointChangeSourceObtaining, 5, 5, 15).Return(nil)
			return databaseMock
		},
		ExpectedActualChange: ptr(5),
		ExpectedFinalValue:   ptr(15),
	},
	{
		// Source is "territory-loss" with zero desired value. Success.
		Name:   "Loss_ZeroValue_Success",
		UserId: 1,
		Change: typepoints.TerritoryPointChange{SourceUserId: 2, ChangeSource: typepoints.TerritoryPointChangeSourceLoss, DesiredChangeValue: 0},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetTerritoryPointsCommand", 1).Return(10, nil)
			databaseMock.On("ChangeTerritoryPointsCommand", 1, 0).Return(nil)
			databaseMock.On("AddTerritoryPointHistoryCommand", 1, 2, typepoints.TerritoryPointChangeSourceLoss, 0, 0, 10).Return(nil)
			return databaseMock
		},
		ExpectedActualChange: ptr(0),
		ExpectedFinalValue:   ptr(10),
	},
	{
		// Source is "territory-loss" with negative desired value. Success.
		Name:   "Loss_NegativeValue_Success",
		UserId: 1,
		Change: typepoints.TerritoryPointChange{SourceUserId: 2, ChangeSource: typepoints.TerritoryPointChangeSourceLoss, DesiredChangeValue: -3},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetTerritoryPointsCommand", 1).Return(10, nil)
			databaseMock.On("ChangeTerritoryPointsCommand", 1, -3).Return(nil)
			databaseMock.On("AddTerritoryPointHistoryCommand", 1, 2, typepoints.TerritoryPointChangeSourceLoss, -3, -3, 7).Return(nil)
			return databaseMock
		},
		ExpectedActualChange: ptr(-3),
		ExpectedFinalValue:   ptr(7),
	},
	{
		// Source is "territory-loss", decrease exceeds current points — clamped to zero.
		Name:   "Loss_DecreaseClamped_Success",
		UserId: 1,
		Change: typepoints.TerritoryPointChange{SourceUserId: 2, ChangeSource: typepoints.TerritoryPointChangeSourceLoss, DesiredChangeValue: -20},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetTerritoryPointsCommand", 1).Return(5, nil)
			databaseMock.On("ChangeTerritoryPointsCommand", 1, -5).Return(nil)
			databaseMock.On("AddTerritoryPointHistoryCommand", 1, 2, typepoints.TerritoryPointChangeSourceLoss, -20, -5, 0).Return(nil)
			return databaseMock
		},
		ExpectedActualChange: ptr(-5),
		ExpectedFinalValue:   ptr(0),
	},
	{
		// Source is "other" with increase. Success.
		Name:   "Other_Increase_Success",
		UserId: 1,
		Change: typepoints.TerritoryPointChange{SourceUserId: 2, ChangeSource: typepoints.TerritoryPointChangeSourceOther, DesiredChangeValue: 5},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetTerritoryPointsCommand", 1).Return(10, nil)
			databaseMock.On("ChangeTerritoryPointsCommand", 1, 5).Return(nil)
			databaseMock.On("AddTerritoryPointHistoryCommand", 1, 2, typepoints.TerritoryPointChangeSourceOther, 5, 5, 15).Return(nil)
			return databaseMock
		},
		ExpectedActualChange: ptr(5),
		ExpectedFinalValue:   ptr(15),
	},
	{
		// Source is "other" with decrease. Success.
		Name:   "Other_Decrease_Success",
		UserId: 1,
		Change: typepoints.TerritoryPointChange{SourceUserId: 2, ChangeSource: typepoints.TerritoryPointChangeSourceOther, DesiredChangeValue: -3},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetTerritoryPointsCommand", 1).Return(10, nil)
			databaseMock.On("ChangeTerritoryPointsCommand", 1, -3).Return(nil)
			databaseMock.On("AddTerritoryPointHistoryCommand", 1, 2, typepoints.TerritoryPointChangeSourceOther, -3, -3, 7).Return(nil)
			return databaseMock
		},
		ExpectedActualChange: ptr(-3),
		ExpectedFinalValue:   ptr(7),
	},
	{
		// AddTerritoryPointHistoryCommand returns a database error.
		Name:   "AddHistory_DatabaseError",
		UserId: 1,
		Change: typepoints.TerritoryPointChange{SourceUserId: 2, ChangeSource: typepoints.TerritoryPointChangeSourceOther, DesiredChangeValue: 5},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetTerritoryPointsCommand", 1).Return(10, nil)
			databaseMock.On("ChangeTerritoryPointsCommand", 1, 5).Return(nil)
			databaseMock.On("AddTerritoryPointHistoryCommand", 1, 2, typepoints.TerritoryPointChangeSourceOther, 5, 5, 15).Return(dbError)
			return databaseMock
		},
		ExpectedErrorCode: "",
	},
}

func TestSrvPoints_ChangeTerritoryPoints(test *testing.T) {
	for _, testCase := range ChangeTerritoryPointsTestCases {
		test.Run(testCase.Name, func(test *testing.T) {
			// Arrange
			databaseMock := testCase.SetupMock()
			sut := srvpoints.Service{Database: databaseMock}

			// Act
			result, err := sut.ChangeTerritoryPoints(testCase.UserId, testCase.Change)

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
