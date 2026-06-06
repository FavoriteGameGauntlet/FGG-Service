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

type TerritoryHoursChange struct {
	ChangeSource       string
	DesiredChangeValue int
	IsSomeones         bool
}

const (
	ExperienceChangeSourceLevelUp = "level-up"
	ExperienceChangeSourceOther   = "other"

	TerritoryHoursChangeSourceSeize = "seize"
	TerritoryHoursChangeSourceOther = "other"

	FreePointsChangeSourceQuestCompletion = "quest-completion"
	FreePointsChangeSourceWheelEffect     = "wheel-effect"
	FreePointsChangeSourceNoPathToBase    = "no-path-to-base"
	FreePointsChangeSourceOther           = "other"
)

var ExperienceChangeSourceSlice = []string{ExperienceChangeSourceLevelUp, ExperienceChangeSourceOther}

var TerritoryHoursChangeSourceSlice = []string{TerritoryHoursChangeSourceSeize, TerritoryHoursChangeSourceOther}

var FreePointsGainSourceSlice = []string{
	FreePointsChangeSourceQuestCompletion,
	FreePointsChangeSourceWheelEffect,
	FreePointsChangeSourceOther,
}

var FreePointsLossSourceSlice = []string{
	FreePointsChangeSourceQuestCompletion,
	FreePointsChangeSourceWheelEffect,
	FreePointsChangeSourceNoPathToBase,
	FreePointsChangeSourceOther,
}
