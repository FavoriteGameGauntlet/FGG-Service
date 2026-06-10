package srvpoints

import (
	"FGG-Service/src/common"
	"FGG-Service/src/points/database"
	"FGG-Service/src/points/type"
	"FGG-Service/src/sysparams/service"
	typesysparams "FGG-Service/src/sysparams/types"
	"slices"
	"strconv"
)

type Service struct {
	Database         dbpoints.IDatabase
	SysParamsService srvsysparams.IService
}

func NewService() *Service {
	pdb := new(dbpoints.Database)
	sps := srvsysparams.NewService()

	return &Service{
		Database:         pdb,
		SysParamsService: sps,
	}
}

func (s *Service) GetExperiencePoints(userId int) (int, error) {
	return s.Database.GetExperiencePointsCommand(userId)
}

func (s *Service) GetFreePoints(userId int) (int, error) {
	return s.Database.GetFreePointsCommand(userId)
}

func (s *Service) GetUserFreePointHistory(userId int) (typepoints.FreePointChangeHistories, error) {
	return s.Database.GetFreePointHistoryCommand(userId)
}

func (s *Service) GetTerritoryHours(userId int) (int, error) {
	return s.Database.GetTerritoryHoursCommand(userId)
}

func (s *Service) GetTerritoryPoints(userId int) (int, error) {
	return s.Database.GetTerritoryPointsCommand(userId)
}

func (s *Service) GetUserTerritoryPointHistory(userId int) (typepoints.TerritoryPointChangeHistories, error) {
	return s.Database.GetTerritoryPointHistoryCommand(userId)
}

func (s *Service) GetUserPointInfo(userId int) (typepoints.PointInfo, error) {
	return s.Database.GetPointInfoCommand(userId)
}

func (s *Service) GetAllPointInfo() (typepoints.PointInfoByLogins, error) {
	return s.Database.GetAllPointInfoCommand()
}

func (s *Service) ChangeFreePoints(userId int, pointChange typepoints.FreePointChange, effectId *int) (
	result typepoints.PointChangeResult, err error) {

	if !slices.Contains(typepoints.FreePointsChangeSourceSlice, pointChange.ChangeSource) {
		err = common.NewChangeSourceUnprocessableError(typepoints.FreePointsChangeSourceSlice)
		return
	}

	if pointChange.ChangeSource == typepoints.FreePointsChangeSourceWheelEffect {
		err = validateWheelEffectChange(pointChange, effectId)

		if err != nil {
			return
		}
	}

	if pointChange.ChangeSource == typepoints.FreePointsChangeSourceBaseTeleport ||
		pointChange.ChangeSource == typepoints.FreePointsChangeSourceSandStorm {
		err = validateTeleportBaseOrSandstormChange(pointChange)

		if err != nil {
			return
		}
	}

	currentPoints, err := s.Database.GetFreePointsCommand(userId)

	if err != nil {
		return
	}

	freePointMinimum, err := s.SysParamsService.GetInt(typesysparams.ParamFreePointsMinimum)

	if err != nil {
		return
	}

	shouldLimitFreePoints, err := s.SysParamsService.GetBool(typesysparams.ParamShouldLimitFreePoints)

	if err != nil {
		return
	}

	// We assume that the current points are greater and subtract the smaller from the larger
	finalValue := currentPoints + pointChange.DesiredChangeValue
	changeValue := pointChange.DesiredChangeValue
	if shouldLimitFreePoints {
		if finalValue < freePointMinimum {
			finalValue = freePointMinimum
			// The current points are not enough, we calculate how many points are between the current and the minimum
			changeValue = freePointMinimum - currentPoints
		}
	}

	err = s.Database.ChangeFreePointsCommand(userId, changeValue)

	if err != nil {
		return
	}

	err = s.Database.AddFreePointHistoryCommand(
		userId,
		pointChange.SourceUserId,
		pointChange.ChangeSource,
		pointChange.DesiredChangeValue,
		changeValue,
		finalValue,
		effectId)

	if err != nil {
		return
	}

	result = typepoints.PointChangeResult{
		ActualChangeValue:  changeValue,
		ChangeSource:       pointChange.ChangeSource,
		DesiredChangeValue: pointChange.DesiredChangeValue,
		FinalValue:         finalValue,
	}

	return
}

func validateWheelEffectChange(pointChange typepoints.FreePointChange, effectId *int) error {
	if effectId == nil {
		return common.NewWheelEffectNameRequiredUnprocessableError(pointChange.ChangeSource)
	}

	return nil
}

func validateTeleportBaseOrSandstormChange(pointChange typepoints.FreePointChange) error {
	if pointChange.DesiredChangeValue > 0 {
		return common.NewWrongDesiredChangeValueConflictError(
			pointChange.ChangeSource,
			common.ConstraintZeroOrLess)
	}

	return nil
}

func (s *Service) ChangeTerritoryHours(userId int, pointChange typepoints.TerritoryHourChange) (
	changeResult typepoints.PointChangeResult, err error) {

	if !slices.Contains(typepoints.TerritoryHourChangeSourceSlice, pointChange.ChangeSource) {
		err = common.NewChangeSourceUnprocessableError(typepoints.TerritoryHourChangeSourceSlice)
		return
	}

	if pointChange.ChangeSource == typepoints.TerritoryHourChangeSourceSeize {
		var seizeDecreaseSlice []int
		seizeDecreaseSlice, err = s.SysParamsService.GetIntSlice(typesysparams.ParamTerritoryHoursSeizeDecreaseSlice)

		if err != nil {
			return
		}

		var seizePenaltyPoints int
		seizePenaltyPoints, err = s.SysParamsService.GetInt(typesysparams.ParamSeizePenaltyPoints)

		if err != nil {
			return
		}

		err = validateSeizeChange(pointChange, seizeDecreaseSlice, seizePenaltyPoints)

		if err != nil {
			return
		}
	}

	currentHours, err := s.Database.GetTerritoryHoursCommand(userId)
	if err != nil {
		return
	}

	if pointChange.ChangeSource == typepoints.TerritoryHourChangeSourceSeize {
		if currentHours+pointChange.DesiredChangeValue < 0 {
			err = common.NewNotEnoughCurrentPointsConflictError(
				pointChange.ChangeSource,
				-pointChange.DesiredChangeValue)
			return
		}
	}

	actualChangeValue := max(pointChange.DesiredChangeValue, -currentHours)

	err = s.Database.ChangeTerritoryHoursCommand(userId, actualChangeValue)

	if err != nil {
		return
	}

	changeResult = typepoints.PointChangeResult{
		ActualChangeValue:  actualChangeValue,
		ChangeSource:       pointChange.ChangeSource,
		DesiredChangeValue: pointChange.DesiredChangeValue,
		FinalValue:         currentHours + actualChangeValue,
	}

	return
}

func (s *Service) ChangeTerritoryPoints(userId int, pointChange typepoints.TerritoryPointChange) (
	changeResult typepoints.PointChangeResult, err error) {

	if !slices.Contains(typepoints.TerritoryPointChangeSourceSlice, pointChange.ChangeSource) {
		err = common.NewChangeSourceUnprocessableError(typepoints.TerritoryPointChangeSourceSlice)
		return
	}

	if pointChange.ChangeSource == typepoints.TerritoryPointChangeSourceObtaining {
		err = validateTerritoryObtainingChange(pointChange)

		if err != nil {
			return
		}
	}

	if pointChange.ChangeSource == typepoints.TerritoryPointChangeSourceLoss {
		err = validateTerritoryLossChange(pointChange)

		if err != nil {
			return
		}
	}

	currentPoints, err := s.Database.GetTerritoryPointsCommand(userId)

	if err != nil {
		return
	}

	actualChangeValue := max(pointChange.DesiredChangeValue, -currentPoints)
	finalValue := currentPoints + actualChangeValue

	err = s.Database.ChangeTerritoryPointsCommand(userId, actualChangeValue)

	if err != nil {
		return
	}

	err = s.Database.AddTerritoryPointHistoryCommand(
		userId,
		pointChange.SourceUserId,
		pointChange.ChangeSource,
		pointChange.DesiredChangeValue,
		actualChangeValue,
		finalValue)

	if err != nil {
		return
	}

	changeResult = typepoints.PointChangeResult{
		ActualChangeValue:  actualChangeValue,
		ChangeSource:       pointChange.ChangeSource,
		DesiredChangeValue: pointChange.DesiredChangeValue,
		FinalValue:         finalValue,
	}

	return
}

func validateTerritoryObtainingChange(pointChange typepoints.TerritoryPointChange) error {
	if pointChange.DesiredChangeValue < 0 {
		return common.NewWrongDesiredChangeValueConflictError(
			pointChange.ChangeSource,
			common.ConstraintZeroOrMore)
	}

	return nil
}

func validateTerritoryLossChange(pointChange typepoints.TerritoryPointChange) error {
	if pointChange.DesiredChangeValue > 0 {
		return common.NewWrongDesiredChangeValueConflictError(
			pointChange.ChangeSource,
			common.ConstraintZeroOrLess)
	}

	return nil
}

func validateSeizeChange(pointChange typepoints.TerritoryHourChange, seizeDecreaseSlice []int, penaltyPoints int) error {
	if pointChange.DesiredChangeValue > 0 {
		return common.NewWrongDesiredChangeValueConflictError(
			pointChange.ChangeSource,
			common.ConstraintZeroOrLess)
	}

	if !pointChange.IsSomeones {
		penaltyPoints = 0
	}

	if !slices.Contains(seizeDecreaseSlice, pointChange.DesiredChangeValue+penaltyPoints) {
		decreaseSlice := make([]int, len(seizeDecreaseSlice))
		for i, v := range seizeDecreaseSlice {
			decreaseSlice[i] = v - penaltyPoints
		}

		return common.NewWrongDesiredChangeValueConflictError(
			pointChange.ChangeSource,
			"one of: "+common.ConvertIntSliceToString(decreaseSlice))
	}

	return nil
}

func (s *Service) ChangeExperiencePoints(userId int, pointChange typepoints.PointChange) (
	changeResult typepoints.PointChangeResult, err error) {

	if !slices.Contains(typepoints.ExperienceChangeSourceSlice, pointChange.ChangeSource) {
		err = common.NewChangeSourceUnprocessableError(typepoints.ExperienceChangeSourceSlice)
		return
	}

	var experiencePointsLevelUp int
	if pointChange.ChangeSource == typepoints.ExperienceChangeSourceLevelUp {
		experiencePointsLevelUp, err = s.SysParamsService.GetInt(typesysparams.ParamExperiencePointsLevelUp)

		if err != nil {
			return
		}

		err = validateLevelUpChange(pointChange, experiencePointsLevelUp)

		if err != nil {
			return
		}
	}

	currentPoints, err := s.Database.GetExperiencePointsCommand(userId)

	if err != nil {
		return
	}

	actualChangeValue := max(pointChange.DesiredChangeValue, -currentPoints)

	if pointChange.ChangeSource == typepoints.ExperienceChangeSourceLevelUp &&
		currentPoints-pointChange.DesiredChangeValue < 0 {

		err = common.NewNotEnoughCurrentPointsConflictError(
			pointChange.ChangeSource,
			-experiencePointsLevelUp)
		return
	}

	err = s.Database.ChangeExperiencePointsCommand(userId, actualChangeValue)

	if err != nil {
		return
	}

	finalValue := currentPoints + actualChangeValue
	changeResult = typepoints.PointChangeResult{
		ActualChangeValue:  actualChangeValue,
		ChangeSource:       pointChange.ChangeSource,
		DesiredChangeValue: pointChange.DesiredChangeValue,
		FinalValue:         finalValue,
	}

	return
}

func validateLevelUpChange(pointChange typepoints.PointChange, experiencePointsLevelUp int) error {
	if pointChange.DesiredChangeValue != experiencePointsLevelUp {
		return common.NewWrongDesiredChangeValueConflictError(
			pointChange.ChangeSource,
			strconv.Itoa(experiencePointsLevelUp))
	}

	return nil
}
