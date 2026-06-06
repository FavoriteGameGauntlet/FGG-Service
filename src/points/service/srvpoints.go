package srvpoints

import (
	"FGG-Service/src/common"
	"FGG-Service/src/points/database"
	"FGG-Service/src/points/type"
	"slices"
	"strconv"
)

type Service struct {
	Database dbpoints.IDatabase
}

func NewService() *Service {
	return &Service{Database: new(dbpoints.Database)}
}

func (s *Service) GetExperiencePoints(userId int) (int, error) {
	return s.Database.GetExperiencePointsCommand(userId)
}

func (s *Service) GetFreePoints(userId int) (int, error) {
	return s.Database.GetFreePointsCommand(userId)
}

func (s *Service) GetTerritoryHours(userId int) (int, error) {
	return s.Database.GetTerritoryHoursCommand(userId)
}

func (s *Service) GetTerritoryPoints(userId int) (int, error) {
	return s.Database.GetTerritoryPointsCommand(userId)
}

func (s *Service) ChangeFreePoints(userId int, pointChange typepoints.PointChange) (
	result typepoints.PointChangeResult, err error) {

	if pointChange.DesiredChangeValue > 0 {
		if !slices.Contains(typepoints.FreePointsGainSourceSlice, pointChange.ChangeSource) {
			err = common.NewChangeSourceUnprocessableError(typepoints.FreePointsGainSourceSlice)
			return
		}
	} else {
		if !slices.Contains(typepoints.FreePointsLossSourceSlice, pointChange.ChangeSource) {
			err = common.NewChangeSourceUnprocessableError(typepoints.FreePointsLossSourceSlice)
			return
		}
	}

	currentPoints, err := s.Database.GetFreePointsCommand(userId)

	if err != nil {
		return
	}

	// We assume that the current points are greater and subtract the smaller from the larger
	finalValue := currentPoints + pointChange.DesiredChangeValue
	changeValue := pointChange.DesiredChangeValue
	if common.FreePointMinimum != nil {
		pointMinimum := *common.FreePointMinimum

		if finalValue < pointMinimum {
			finalValue = pointMinimum
			// The current points are not enough, we calculate how many points are between the current and the minimum
			changeValue = pointMinimum - currentPoints
		}
	}

	err = s.Database.ChangeFreePointsCommand(userId, changeValue)

	if err != nil {
		return
	}

	err = s.Database.AddFreePointHistoryCommand(userId, pointChange.ChangeSource, pointChange.DesiredChangeValue, changeValue, finalValue)

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

func (s *Service) ChangeTerritoryHours(userId int, pointChange typepoints.TerritoryHoursChange) (
	changeResult typepoints.PointChangeResult, err error) {

	if !slices.Contains(typepoints.TerritoryHoursChangeSourceSlice, pointChange.ChangeSource) {
		err = common.NewChangeSourceUnprocessableError(typepoints.TerritoryHoursChangeSourceSlice)
		return
	}

	if pointChange.ChangeSource == typepoints.TerritoryHoursChangeSourceSeize {
		err = validateSeizeChange(pointChange)
		if err != nil {
			return
		}
	}

	currentHours, err := s.Database.GetTerritoryHoursCommand(userId)
	if err != nil {
		return
	}

	if pointChange.ChangeSource == typepoints.TerritoryHoursChangeSourceSeize {
		if currentHours+pointChange.DesiredChangeValue < 0 {
			err = common.NewNotEnoughCurrentPointsConflictError(pointChange.ChangeSource, -pointChange.DesiredChangeValue)
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

func validateSeizeChange(pointChange typepoints.TerritoryHoursChange) error {
	if pointChange.DesiredChangeValue > 0 {
		return common.NewWrongDesiredChangeValueConflictError(
			typepoints.TerritoryHoursChangeSourceSeize, "zero or less")
	}

	penaltyPoints := 0
	if pointChange.IsSomeones {
		penaltyPoints = 1
	}

	if !slices.Contains(common.DefaultTerritoryHoursSeizeDecreaseSlice, pointChange.DesiredChangeValue+penaltyPoints) {
		decreaseSlice := make([]int, len(common.DefaultTerritoryHoursSeizeDecreaseSlice))
		for i, v := range common.DefaultTerritoryHoursSeizeDecreaseSlice {
			decreaseSlice[i] = v - penaltyPoints
		}

		return common.NewWrongDesiredChangeValueConflictError(
			typepoints.TerritoryHoursChangeSourceSeize,
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

	if pointChange.ChangeSource == typepoints.ExperienceChangeSourceLevelUp {
		err = validateLevelUpChange(pointChange)
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
			-common.DefaultExperiencePointsLevelUp)
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

func validateLevelUpChange(pointChange typepoints.PointChange) error {
	if pointChange.DesiredChangeValue != common.DefaultExperiencePointsLevelUp {
		return common.NewWrongDesiredChangeValueConflictError(
			typepoints.ExperienceChangeSourceLevelUp,
			strconv.Itoa(common.DefaultExperiencePointsLevelUp))
	}
	return nil
}
