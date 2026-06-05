package srvpoints

import (
	"FGG-Service/src/common"
	"FGG-Service/src/points/database"
	"FGG-Service/src/points/type"
	"slices"
	"strconv"
	"time"
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
		ChangeDate:         time.Now(),
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

	if !slices.Contains(common.DefaultTerritoryHoursSeizeDecreasingSlice, pointChange.DesiredChangeValue+penaltyPoints) {
		decreasingSlice := make([]int, len(common.DefaultTerritoryHoursSeizeDecreasingSlice))
		for i, v := range common.DefaultTerritoryHoursSeizeDecreasingSlice {
			decreasingSlice[i] = v - penaltyPoints
		}

		return common.NewWrongDesiredChangeValueConflictError(
			typepoints.TerritoryHoursChangeSourceSeize,
			"one of: "+common.ConvertIntSliceToString(decreasingSlice))
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
		ChangeDate:         time.Now(),
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
