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

func (s *Service) ChangeAvailableRolls(userId int, changeValue int) error {
	return s.Database.ChangeAvailableRollsCommand(userId, changeValue)
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

	if pointChange.ChangeSource == typepoints.FreePointsChangeSourceBaseTeleport {
		var freePointChangeByBaseTeleport int
		freePointChangeByBaseTeleport, err = s.SysParamsService.GetInt(typesysparams.ParamFreePointChangeByBaseTeleport)

		if err != nil {
			return
		}

		err = validateBaseTeleportChange(pointChange, freePointChangeByBaseTeleport)

		if err != nil {
			return
		}
	}

	if pointChange.ChangeSource == typepoints.FreePointsChangeSourceSandStorm {
		var freePointChangeBySandstorm int
		freePointChangeBySandstorm, err = s.SysParamsService.GetInt(typesysparams.ParamFreePointChangeBySandstorm)

		if err != nil {
			return
		}

		err = validateSandstormChange(pointChange, freePointChangeBySandstorm)

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

func validateBaseTeleportChange(pointChange typepoints.FreePointChange, expectedChangeValue int) error {
	if pointChange.DesiredChangeValue != expectedChangeValue {
		return common.NewWrongDesiredChangeValueConflictError(
			pointChange.ChangeSource,
			strconv.Itoa(expectedChangeValue))
	}

	return nil
}

func validateSandstormChange(pointChange typepoints.FreePointChange, expectedChangeValue int) error {
	if pointChange.DesiredChangeValue != expectedChangeValue {
		return common.NewWrongDesiredChangeValueConflictError(
			pointChange.ChangeSource,
			strconv.Itoa(expectedChangeValue))
	}

	return nil
}

func (s *Service) ChangeTerritoryHours(userId int, pointChange typepoints.TerritoryHourChange, targetUserId *int) (
	result typepoints.PointChangeResultByTypes, err error) {

	if !slices.Contains(typepoints.TerritoryHourChangeSourceSlice, pointChange.ChangeSource) {
		err = common.NewChangeSourceUnprocessableError(typepoints.TerritoryHourChangeSourceSlice)
		return
	}

	isSeize := pointChange.ChangeSource == typepoints.TerritoryHourChangeSourceSeize

	seizeIndex := -1
	var territoryPointChangeBySeizeSlice []int

	if isSeize {
		seizeIndex, territoryPointChangeBySeizeSlice, err = s.resolveSeizeConversion(userId, pointChange, targetUserId)
		if err != nil {
			return
		}
	}

	currentHours, actualHoursChange, err := s.getTerritoryHoursBalance(userId, pointChange, isSeize)
	if err != nil {
		return
	}

	currentPoints, territoryPointsGranted, targetCurrentPoints, err := s.resolveTerritoryPointsGrant(
		userId,
		targetUserId,
		pointChange,
		isSeize,
		seizeIndex,
		territoryPointChangeBySeizeSlice)

	if err != nil {
		return
	}

	err = s.Database.ChangeTerritoryHoursCommand(userId, actualHoursChange)
	if err != nil {
		return
	}

	result = typepoints.PointChangeResultByTypes{
		typepoints.PointTypeTerritoryHours: typepoints.PointChangeResult{
			ActualChangeValue:  actualHoursChange,
			ChangeSource:       pointChange.ChangeSource,
			DesiredChangeValue: pointChange.DesiredChangeValue,
			FinalValue:         currentHours + actualHoursChange,
		},
	}

	if isSeize {
		var pointsResult typepoints.PointChangeResult
		pointsResult, err = s.applyTerritorySeize(
			userId,
			targetUserId,
			pointChange,
			currentPoints,
			territoryPointsGranted,
			targetCurrentPoints)

		if err != nil {
			return
		}

		result[typepoints.PointTypeTerritoryPoints] = pointsResult
	}

	return
}

// resolveSeizeConversion validates the requested hour amount against the seize sysparams and
// returns the matching index into, and value of, the territory-points-by-seize sysparam.
func (s *Service) resolveSeizeConversion(userId int, pointChange typepoints.TerritoryHourChange, targetUserId *int) (
	seizeIndex int, territoryPointChangeBySeizeSlice []int, err error) {

	territoryHourChangeBySeizeSlice, err := s.SysParamsService.GetIntSlice(typesysparams.ParamTerritoryHourChangeBySeizeSlice)
	if err != nil {
		return
	}

	seizePenaltyPoints, err := s.SysParamsService.GetInt(typesysparams.ParamSeizePenaltyPoints)
	if err != nil {
		return
	}

	seizeIndex, err = validateSeizeChange(pointChange, territoryHourChangeBySeizeSlice, seizePenaltyPoints)
	if err != nil {
		return
	}

	if pointChange.IsSomeones && targetUserId == nil {
		err = common.NewTargetLoginRequiredUnprocessableError(pointChange.ChangeSource)
		return
	}

	if pointChange.IsSomeones && *targetUserId == userId {
		err = common.NewCannotTargetSelfConflictError(pointChange.ChangeSource)
		return
	}

	territoryPointChangeBySeizeSlice, err = s.SysParamsService.GetIntSlice(typesysparams.ParamTerritoryPointChangeBySeizeSlice)
	if err != nil {
		return
	}

	if seizeIndex >= len(territoryPointChangeBySeizeSlice) {
		err = common.NewSystemParameterNotFoundError(typesysparams.ParamTerritoryPointChangeBySeizeSlice)
		return
	}

	return
}

// getTerritoryHoursBalance fetches the user's current hours and, for seize, validates there are enough.
func (s *Service) getTerritoryHoursBalance(userId int, pointChange typepoints.TerritoryHourChange, isSeize bool) (
	currentHours, actualHoursChange int, err error) {

	currentHours, err = s.Database.GetTerritoryHoursCommand(userId)
	if err != nil {
		return
	}

	if isSeize && currentHours+pointChange.DesiredChangeValue < 0 {
		err = common.NewNotEnoughCurrentPointsConflictError(
			pointChange.ChangeSource,
			-pointChange.DesiredChangeValue)
		return
	}

	actualHoursChange = max(pointChange.DesiredChangeValue, -currentHours)
	return
}

// resolveTerritoryPointsGrant is a no-op unless seizing: it fetches the user's current points,
// computes the granted amount, and validates the target has enough points to lose when seizing
// someone else's territory.
func (s *Service) resolveTerritoryPointsGrant(
	userId int, targetUserId *int, pointChange typepoints.TerritoryHourChange,
	isSeize bool, seizeIndex int, territoryPointChangeBySeizeSlice []int) (
	currentPoints, territoryPointsGranted, targetCurrentPoints int, err error) {

	if !isSeize {
		return
	}

	currentPoints, err = s.Database.GetTerritoryPointsCommand(userId)
	if err != nil {
		return
	}

	territoryPointsGranted = territoryPointChangeBySeizeSlice[seizeIndex]

	if pointChange.IsSomeones {
		targetCurrentPoints, err = s.Database.GetTerritoryPointsCommand(*targetUserId)
		if err != nil {
			return
		}

		if targetCurrentPoints < territoryPointsGranted {
			err = common.NewNotEnoughCurrentPointsConflictError(
				typepoints.TerritoryPointChangeSourceLoss,
				territoryPointsGranted)
			return
		}
	}

	return
}

// applyTerritorySeize writes the self points gain (and, when seizing someone else's territory, the
// target's points loss) and returns the caller's territoryPoints result.
func (s *Service) applyTerritorySeize(
	userId int, targetUserId *int, pointChange typepoints.TerritoryHourChange,
	currentPoints, territoryPointsGranted, targetCurrentPoints int) (
	pointsResult typepoints.PointChangeResult, err error) {

	newSelfFinal := currentPoints + territoryPointsGranted

	err = s.Database.ChangeTerritoryPointsCommand(userId, territoryPointsGranted)
	if err != nil {
		return
	}

	err = s.Database.AddTerritoryPointHistoryCommand(
		userId,
		userId,
		typepoints.TerritoryPointChangeSourceObtaining,
		territoryPointsGranted,
		territoryPointsGranted,
		newSelfFinal)
	if err != nil {
		return
	}

	pointsResult = typepoints.PointChangeResult{
		ActualChangeValue:  territoryPointsGranted,
		ChangeSource:       typepoints.TerritoryPointChangeSourceObtaining,
		DesiredChangeValue: territoryPointsGranted,
		FinalValue:         newSelfFinal,
	}

	if !pointChange.IsSomeones {
		return
	}

	targetFinal := targetCurrentPoints - territoryPointsGranted

	err = s.Database.ChangeTerritoryPointsCommand(*targetUserId, -territoryPointsGranted)
	if err != nil {
		return
	}

	err = s.Database.AddTerritoryPointHistoryCommand(
		*targetUserId,
		userId,
		typepoints.TerritoryPointChangeSourceLoss,
		-territoryPointsGranted,
		-territoryPointsGranted,
		targetFinal)
	if err != nil {
		return
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

func validateSeizeChange(pointChange typepoints.TerritoryHourChange, seizeDecreaseSlice []int, penaltyPoints int) (index int, err error) {
	index = -1

	if pointChange.DesiredChangeValue > 0 {
		err = common.NewWrongDesiredChangeValueConflictError(
			pointChange.ChangeSource,
			common.ConstraintZeroOrLess)
		return
	}

	if !pointChange.IsSomeones {
		penaltyPoints = 0
	}

	decreaseSlice := make([]int, len(seizeDecreaseSlice))
	for i, v := range seizeDecreaseSlice {
		decreaseSlice[i] = v + penaltyPoints
	}

	index = slices.Index(decreaseSlice, pointChange.DesiredChangeValue)

	if index == -1 {
		err = common.NewWrongDesiredChangeValueConflictError(
			pointChange.ChangeSource,
			"one of: "+common.ConvertIntSliceToString(decreaseSlice))
	}

	return
}

func (s *Service) ChangeExperiencePoints(userId int, pointChange typepoints.PointChange) (
	changeResult typepoints.PointChangeResult, err error) {

	if !slices.Contains(typepoints.ExperienceChangeSourceSlice, pointChange.ChangeSource) {
		err = common.NewChangeSourceUnprocessableError(typepoints.ExperienceChangeSourceSlice)
		return
	}

	var experiencePointByLevelUp int
	if pointChange.ChangeSource == typepoints.ExperienceChangeSourceLevelUp {
		experiencePointByLevelUp, err = s.SysParamsService.GetInt(typesysparams.ParamExperiencePointByLevelUp)

		if err != nil {
			return
		}

		err = validateLevelUpChange(pointChange, experiencePointByLevelUp)

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
			-experiencePointByLevelUp)
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
