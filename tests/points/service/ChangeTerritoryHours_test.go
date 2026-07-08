package srvpoints_test

import (
	"FGG-Service/src/common"
	"FGG-Service/src/points/service"
	"FGG-Service/src/points/type"
	"FGG-Service/src/sysparams/types"
	"FGG-Service/tests/points/mock"
	"FGG-Service/tests/sysparams/srvmock"
	"testing"

	"github.com/stretchr/testify/require"
)

type ChangeTerritoryHoursTestCase struct {
	Name                       string
	UserId                     int
	Change                     typepoints.TerritoryHourChange
	TargetUserId               *int
	SetupMock                  func() *dbpointsmock.DatabaseMock
	SetupSysParams             func() *srvsysparamsmock.ServiceMock
	ExpectedActualHoursChange  *int
	ExpectedFinalHours         *int
	ExpectedActualPointsChange *int
	ExpectedFinalPoints        *int
	ExpectedPointsChangeSource string
	ExpectedErrorCode          string
}

func defaultSeizeDecreaseSliceSysParams() *srvsysparamsmock.ServiceMock {
	spSvc := new(srvsysparamsmock.ServiceMock)
	spSvc.On("GetIntSlice", typesysparams.ParamTerritoryHourChangeBySeizeSlice).Return([]int{-2, -4}, nil)
	spSvc.On("GetInt", typesysparams.ParamSeizePenaltyPoints).Return(-1, nil)
	return spSvc
}

func seizeSysParamsWithPointsSlice(pointsSlice []int) *srvsysparamsmock.ServiceMock {
	spSvc := defaultSeizeDecreaseSliceSysParams()
	spSvc.On("GetIntSlice", typesysparams.ParamTerritoryPointChangeBySeizeSlice).Return(pointsSlice, nil)
	return spSvc
}

func findResultByType(results typepoints.PointChangeResultByTypes, pointType string) typepoints.PointChangeResult {
	return results[pointType]
}

var ChangeTerritoryHoursTestCases = []ChangeTerritoryHoursTestCase{
	{
		// Change source is not "seize" or "other". Unprocessable error returns.
		Name:   "InvalidSource_UnprocessableError",
		UserId: 1,
		Change: typepoints.TerritoryHourChange{ChangeSource: "invalid", DesiredChangeValue: -2},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			return new(dbpointsmock.DatabaseMock)
		},
		ExpectedErrorCode: "INCORRECT_CHANGE_SOURCE_VALUE",
	},
	{
		// Source is "seize" with positive desired value. Conflict error returns.
		Name:   "Seize_PositiveValue_ConflictError",
		UserId: 1,
		Change: typepoints.TerritoryHourChange{ChangeSource: typepoints.TerritoryHourChangeSourceSeize, DesiredChangeValue: 2},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			return new(dbpointsmock.DatabaseMock)
		},
		SetupSysParams:    defaultSeizeDecreaseSliceSysParams,
		ExpectedErrorCode: "WRONG_DESIRED_CHANGE_VALUE",
	},
	{
		// Source is "seize" with zero desired value. Conflict error returns.
		Name:   "Seize_ZeroValue_ConflictError",
		UserId: 1,
		Change: typepoints.TerritoryHourChange{ChangeSource: typepoints.TerritoryHourChangeSourceSeize, DesiredChangeValue: 0},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			return new(dbpointsmock.DatabaseMock)
		},
		SetupSysParams:    defaultSeizeDecreaseSliceSysParams,
		ExpectedErrorCode: "WRONG_DESIRED_CHANGE_VALUE",
	},
	{
		// Source is "seize" with wrong decrease amount (not -2 or -4). Conflict error returns.
		Name:   "Seize_WrongDecreaseAmount_ConflictError",
		UserId: 1,
		Change: typepoints.TerritoryHourChange{ChangeSource: typepoints.TerritoryHourChangeSourceSeize, DesiredChangeValue: -3},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			return new(dbpointsmock.DatabaseMock)
		},
		SetupSysParams:    defaultSeizeDecreaseSliceSysParams,
		ExpectedErrorCode: "WRONG_DESIRED_CHANGE_VALUE",
	},
	{
		// Source is "seize" with isSomeones=true, wrong value (not 3 or 5). Conflict error returns.
		Name:   "Seize_IsSomeones_WrongDecreaseAmount_ConflictError",
		UserId: 1,
		Change: typepoints.TerritoryHourChange{ChangeSource: typepoints.TerritoryHourChangeSourceSeize, DesiredChangeValue: -2, IsSomeones: true},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			return new(dbpointsmock.DatabaseMock)
		},
		SetupSysParams:    defaultSeizeDecreaseSliceSysParams,
		ExpectedErrorCode: "WRONG_DESIRED_CHANGE_VALUE",
	},
	{
		// Source is "seize", isSomeones=true, valid amount, but no target login provided. Unprocessable error returns.
		Name:         "Seize_IsSomeones_TargetLoginRequired_UnprocessableError",
		UserId:       1,
		Change:       typepoints.TerritoryHourChange{ChangeSource: typepoints.TerritoryHourChangeSourceSeize, DesiredChangeValue: -3, IsSomeones: true},
		TargetUserId: nil,
		SetupMock: func() *dbpointsmock.DatabaseMock {
			return new(dbpointsmock.DatabaseMock)
		},
		SetupSysParams:    defaultSeizeDecreaseSliceSysParams,
		ExpectedErrorCode: "TARGET_LOGIN_REQUIRED",
	},
	{
		// Source is "seize", isSomeones=true, valid amount, but target login resolves to the caller themselves.
		// Unprocessable error returns.
		Name:         "Seize_IsSomeones_TargetIsSelf_UnprocessableError",
		UserId:       1,
		Change:       typepoints.TerritoryHourChange{ChangeSource: typepoints.TerritoryHourChangeSourceSeize, DesiredChangeValue: -3, IsSomeones: true},
		TargetUserId: ptr(1),
		SetupMock: func() *dbpointsmock.DatabaseMock {
			return new(dbpointsmock.DatabaseMock)
		},
		SetupSysParams:    defaultSeizeDecreaseSliceSysParams,
		ExpectedErrorCode: "CANNOT_TARGET_SELF",
	},
	{
		// Source is "seize", isSomeones=true, valid amount, target doesn't have enough territory points.
		// Conflict error returns and no database writes happen.
		Name:         "Seize_IsSomeones_InsufficientTargetPoints_ConflictError",
		UserId:       1,
		Change:       typepoints.TerritoryHourChange{ChangeSource: typepoints.TerritoryHourChangeSourceSeize, DesiredChangeValue: -3, IsSomeones: true},
		TargetUserId: ptr(2),
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetTerritoryHoursCommand", 1).Return(10, nil)
			databaseMock.On("GetTerritoryPointsCommand", 1).Return(100, nil)
			databaseMock.On("GetTerritoryPointsCommand", 2).Return(5, nil)
			return databaseMock
		},
		SetupSysParams: func() *srvsysparamsmock.ServiceMock {
			return seizeSysParamsWithPointsSlice([]int{10, 20})
		},
		ExpectedErrorCode: "NOT_ENOUGH_CURRENT_POINTS",
	},
	{
		// Source is "seize", value -2 valid but not enough hours. Conflict error returns.
		Name:   "Seize_NotEnoughHours_ConflictError",
		UserId: 1,
		Change: typepoints.TerritoryHourChange{ChangeSource: typepoints.TerritoryHourChangeSourceSeize, DesiredChangeValue: -4},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetTerritoryHoursCommand", 1).Return(2, nil)
			return databaseMock
		},
		SetupSysParams: func() *srvsysparamsmock.ServiceMock {
			return seizeSysParamsWithPointsSlice([]int{10, 20})
		},
		ExpectedErrorCode: "NOT_ENOUGH_CURRENT_POINTS",
	},
	{
		// GetTerritoryHoursCommand returns a database error.
		Name:   "Seize_DatabaseError",
		UserId: 1,
		Change: typepoints.TerritoryHourChange{ChangeSource: typepoints.TerritoryHourChangeSourceSeize, DesiredChangeValue: -2},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetTerritoryHoursCommand", 1).Return(0, dbError)
			return databaseMock
		},
		SetupSysParams: func() *srvsysparamsmock.ServiceMock {
			return seizeSysParamsWithPointsSlice([]int{10, 20})
		},
		ExpectedErrorCode: "",
	},
	{
		// Source is "seize", value -2, enough hours, no target. Success — grants points from the parallel slice.
		Name:   "Seize_DecreaseBy2_Success",
		UserId: 1,
		Change: typepoints.TerritoryHourChange{ChangeSource: typepoints.TerritoryHourChangeSourceSeize, DesiredChangeValue: -2},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetTerritoryHoursCommand", 1).Return(10, nil)
			databaseMock.On("ChangeTerritoryHoursCommand", 1, -2).Return(nil)
			databaseMock.On("GetTerritoryPointsCommand", 1).Return(100, nil)
			databaseMock.On("ChangeTerritoryPointsCommand", 1, 10).Return(nil)
			databaseMock.On("AddTerritoryPointHistoryCommand", 1, 1, typepoints.TerritoryPointChangeSourceObtaining, 10, 10, 110).Return(nil)
			return databaseMock
		},
		SetupSysParams: func() *srvsysparamsmock.ServiceMock {
			return seizeSysParamsWithPointsSlice([]int{10, 20})
		},
		ExpectedActualHoursChange:  ptr(-2),
		ExpectedFinalHours:         ptr(8),
		ExpectedActualPointsChange: ptr(10),
		ExpectedFinalPoints:        ptr(110),
		ExpectedPointsChangeSource: typepoints.TerritoryPointChangeSourceObtaining,
	},
	{
		// Source is "seize", value -4, enough hours, no target. Success — grants points from the parallel slice.
		Name:   "Seize_DecreaseBy4_Success",
		UserId: 1,
		Change: typepoints.TerritoryHourChange{ChangeSource: typepoints.TerritoryHourChangeSourceSeize, DesiredChangeValue: -4},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetTerritoryHoursCommand", 1).Return(10, nil)
			databaseMock.On("ChangeTerritoryHoursCommand", 1, -4).Return(nil)
			databaseMock.On("GetTerritoryPointsCommand", 1).Return(100, nil)
			databaseMock.On("ChangeTerritoryPointsCommand", 1, 20).Return(nil)
			databaseMock.On("AddTerritoryPointHistoryCommand", 1, 1, typepoints.TerritoryPointChangeSourceObtaining, 20, 20, 120).Return(nil)
			return databaseMock
		},
		SetupSysParams: func() *srvsysparamsmock.ServiceMock {
			return seizeSysParamsWithPointsSlice([]int{10, 20})
		},
		ExpectedActualHoursChange:  ptr(-4),
		ExpectedFinalHours:         ptr(6),
		ExpectedActualPointsChange: ptr(20),
		ExpectedFinalPoints:        ptr(120),
		ExpectedPointsChangeSource: typepoints.TerritoryPointChangeSourceObtaining,
	},
	{
		// Source is "seize" with isSomeones=true, value -3. Success — target loses the same amount the caller gains.
		Name:         "Seize_IsSomeones_DecreaseBy3_Success",
		UserId:       1,
		Change:       typepoints.TerritoryHourChange{ChangeSource: typepoints.TerritoryHourChangeSourceSeize, DesiredChangeValue: -3, IsSomeones: true},
		TargetUserId: ptr(2),
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetTerritoryHoursCommand", 1).Return(10, nil)
			databaseMock.On("ChangeTerritoryHoursCommand", 1, -3).Return(nil)
			databaseMock.On("GetTerritoryPointsCommand", 1).Return(100, nil)
			databaseMock.On("GetTerritoryPointsCommand", 2).Return(50, nil)
			databaseMock.On("ChangeTerritoryPointsCommand", 1, 10).Return(nil)
			databaseMock.On("AddTerritoryPointHistoryCommand", 1, 1, typepoints.TerritoryPointChangeSourceObtaining, 10, 10, 110).Return(nil)
			databaseMock.On("ChangeTerritoryPointsCommand", 2, -10).Return(nil)
			databaseMock.On("AddTerritoryPointHistoryCommand", 2, 1, typepoints.TerritoryPointChangeSourceLoss, -10, -10, 40).Return(nil)
			return databaseMock
		},
		SetupSysParams: func() *srvsysparamsmock.ServiceMock {
			return seizeSysParamsWithPointsSlice([]int{10, 20})
		},
		ExpectedActualHoursChange:  ptr(-3),
		ExpectedFinalHours:         ptr(7),
		ExpectedActualPointsChange: ptr(10),
		ExpectedFinalPoints:        ptr(110),
		ExpectedPointsChangeSource: typepoints.TerritoryPointChangeSourceObtaining,
	},
	{
		// Source is "seize" with isSomeones=true, value -5. Success — target loses the same amount the caller gains.
		Name:         "Seize_IsSomeones_DecreaseBy5_Success",
		UserId:       1,
		Change:       typepoints.TerritoryHourChange{ChangeSource: typepoints.TerritoryHourChangeSourceSeize, DesiredChangeValue: -5, IsSomeones: true},
		TargetUserId: ptr(2),
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetTerritoryHoursCommand", 1).Return(10, nil)
			databaseMock.On("ChangeTerritoryHoursCommand", 1, -5).Return(nil)
			databaseMock.On("GetTerritoryPointsCommand", 1).Return(100, nil)
			databaseMock.On("GetTerritoryPointsCommand", 2).Return(50, nil)
			databaseMock.On("ChangeTerritoryPointsCommand", 1, 20).Return(nil)
			databaseMock.On("AddTerritoryPointHistoryCommand", 1, 1, typepoints.TerritoryPointChangeSourceObtaining, 20, 20, 120).Return(nil)
			databaseMock.On("ChangeTerritoryPointsCommand", 2, -20).Return(nil)
			databaseMock.On("AddTerritoryPointHistoryCommand", 2, 1, typepoints.TerritoryPointChangeSourceLoss, -20, -20, 30).Return(nil)
			return databaseMock
		},
		SetupSysParams: func() *srvsysparamsmock.ServiceMock {
			return seizeSysParamsWithPointsSlice([]int{10, 20})
		},
		ExpectedActualHoursChange:  ptr(-5),
		ExpectedFinalHours:         ptr(5),
		ExpectedActualPointsChange: ptr(20),
		ExpectedFinalPoints:        ptr(120),
		ExpectedPointsChangeSource: typepoints.TerritoryPointChangeSourceObtaining,
	},
	{
		// Source is "other" with increase. Success — territory points are untouched and absent from the result.
		Name:   "Other_Increase_Success",
		UserId: 1,
		Change: typepoints.TerritoryHourChange{ChangeSource: typepoints.TerritoryHourChangeSourceOther, DesiredChangeValue: 5},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetTerritoryHoursCommand", 1).Return(10, nil)
			databaseMock.On("ChangeTerritoryHoursCommand", 1, 5).Return(nil)
			return databaseMock
		},
		ExpectedActualHoursChange: ptr(5),
		ExpectedFinalHours:        ptr(15),
	},
	{
		// Source is "other" with decrease. Success — territory points are untouched and absent from the result.
		Name:   "Other_Decrease_Success",
		UserId: 1,
		Change: typepoints.TerritoryHourChange{ChangeSource: typepoints.TerritoryHourChangeSourceOther, DesiredChangeValue: -3},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetTerritoryHoursCommand", 1).Return(10, nil)
			databaseMock.On("ChangeTerritoryHoursCommand", 1, -3).Return(nil)
			return databaseMock
		},
		ExpectedActualHoursChange: ptr(-3),
		ExpectedFinalHours:        ptr(7),
	},
	{
		// Source is "other", decrease exceeds current hours — clamped to zero. Territory points are untouched
		// and absent from the result.
		Name:   "Other_DecreaseClamped_Success",
		UserId: 1,
		Change: typepoints.TerritoryHourChange{ChangeSource: typepoints.TerritoryHourChangeSourceOther, DesiredChangeValue: -20},
		SetupMock: func() *dbpointsmock.DatabaseMock {
			databaseMock := new(dbpointsmock.DatabaseMock)
			databaseMock.On("GetTerritoryHoursCommand", 1).Return(5, nil)
			databaseMock.On("ChangeTerritoryHoursCommand", 1, -5).Return(nil)
			return databaseMock
		},
		ExpectedActualHoursChange: ptr(-5),
		ExpectedFinalHours:        ptr(0),
	},
}

func TestSrvPoints_ChangeTerritoryHours(test *testing.T) {
	for _, testCase := range ChangeTerritoryHoursTestCases {
		test.Run(testCase.Name, func(test *testing.T) {
			// Arrange
			databaseMock := testCase.SetupMock()

			var spSvc *srvsysparamsmock.ServiceMock
			if testCase.SetupSysParams != nil {
				spSvc = testCase.SetupSysParams()
			} else {
				spSvc = new(srvsysparamsmock.ServiceMock)
			}

			sut := srvpoints.Service{Database: databaseMock, SysParamsService: spSvc}

			// Act
			result, err := sut.ChangeTerritoryHours(testCase.UserId, testCase.Change, testCase.TargetUserId)

			// Assert
			if testCase.ExpectedErrorCode != "" {
				require.Error(test, err)
				appErr, ok := err.(common.AppError)
				require.True(test, ok)
				require.Equal(test, testCase.ExpectedErrorCode, appErr.GetCode())
			} else if testCase.ExpectedActualHoursChange != nil {
				require.NoError(test, err)
				hoursResult := findResultByType(result, typepoints.PointTypeTerritoryHours)
				require.Equal(test, *testCase.ExpectedActualHoursChange, hoursResult.ActualChangeValue)
				require.Equal(test, *testCase.ExpectedFinalHours, hoursResult.FinalValue)

				if testCase.ExpectedActualPointsChange != nil {
					pointsResult, ok := result[typepoints.PointTypeTerritoryPoints]
					require.True(test, ok)
					require.Equal(test, *testCase.ExpectedActualPointsChange, pointsResult.ActualChangeValue)
					require.Equal(test, *testCase.ExpectedFinalPoints, pointsResult.FinalValue)
					require.Equal(test, testCase.ExpectedPointsChangeSource, pointsResult.ChangeSource)
				} else {
					_, ok := result[typepoints.PointTypeTerritoryPoints]
					require.False(test, ok)
				}
			} else {
				require.Error(test, err)
			}

			databaseMock.AssertExpectations(test)
			spSvc.AssertExpectations(test)
		})
	}
}
