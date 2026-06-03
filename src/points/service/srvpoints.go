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

func (s *Service) ChangeExperiencePoints(userId int, pointChange typepoints.PointChange) (
	changeResult typepoints.PointChangeResult, err error) {

	if !slices.Contains(typepoints.ExperienceChangeSourceSlice, pointChange.ChangeSource) {
		err = common.NewExperienceChangeSourceUnprocessableError(typepoints.ExperienceChangeSourceSlice)
		return
	}

	if pointChange.ChangeSource == typepoints.ExperienceChangeSourceLevelUp &&
		pointChange.DesiredChangeValue != common.DefaultExperiencePointsLevelUp {

		err = common.NewWrongDesiredChangeValueConflictError(
			typepoints.ExperienceChangeSourceLevelUp,
			strconv.Itoa(common.DefaultExperiencePointsLevelUp))
		return
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
