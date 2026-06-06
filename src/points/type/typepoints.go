package typepoints

type PointChangeResult struct {
	ActualChangeValue  int
	ChangeSource       string
	DesiredChangeValue int
	FinalValue         int
}

type PointChange struct {
	ChangeSource       string
	DesiredChangeValue int
}

type TerritoryHourChange struct {
	ChangeSource       string
	DesiredChangeValue int
	IsSomeones         bool
}

type FreePointChange struct {
	SourceUserId       int
	ChangeSource       string
	DesiredChangeValue int
	WheelEffectName    *string
}

const (
	ExperienceChangeSourceLevelUp = "level-up"
	ExperienceChangeSourceOther   = "other"

	TerritoryHourChangeSourceSeize = "seize"
	TerritoryHourChangeSourceOther = "other"

	FreePointsChangeSourceQuestCompletion  = "quest"
	FreePointsChangeSourceOwnWheelEffect   = "own-wheel-effect"
	FreePointsChangeSourceOtherWheelEffect = "other-wheel-effect"
	FreePointsChangeSourceBaseTeleport     = "base-teleport"
	FreePointsChangeSourceSandStorm        = "sandstorm"
	FreePointsChangeSourceOther            = "other"
)

var ExperienceChangeSourceSlice = []string{ExperienceChangeSourceLevelUp, ExperienceChangeSourceOther}

var TerritoryHourChangeSourceSlice = []string{TerritoryHourChangeSourceSeize, TerritoryHourChangeSourceOther}

var FreePointsChangeSourceSlice = []string{
	FreePointsChangeSourceQuestCompletion,
	FreePointsChangeSourceOwnWheelEffect,
	FreePointsChangeSourceOtherWheelEffect,
	FreePointsChangeSourceBaseTeleport,
	FreePointsChangeSourceSandStorm,
	FreePointsChangeSourceOther,
}
