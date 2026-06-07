package typepoints

import "time"

type PointChangeResult struct {
	ActualChangeValue  int
	ChangeSource       string
	DesiredChangeValue int
	FinalValue         int
}

type PointChangeResultByUserId struct {
	Login        string
	UserId       int
	ChangeResult PointChangeResult
}

type PointChangeResultByUserIds = []PointChangeResultByUserId

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
}

type FreePointChangeByUserId struct {
	Login       string
	UserId      int
	PointChange FreePointChange
}

type FreePointChangeByUserIds = []FreePointChangeByUserId

type TerritoryPointChange struct {
	SourceUserId       int
	ChangeSource       string
	DesiredChangeValue int
}

type PointChangeHistory struct {
	ActualChangeValue  int
	ChangeDate         time.Time
	ChangeSource       string
	DesiredChangeValue int
	FinalValue         int
}

type PointChangeHistories = []PointChangeHistory

type PointInfo struct {
	TerritoryPoints  int
	FreePoints       int
	AvailableRolls   int
	TerritoryHours   int
	ExperiencePoints int
}

type PointInfoByLogin struct {
	Login     string
	PointInfo PointInfo
}

type PointInfoByLogins = []PointInfoByLogin

const (
	ExperienceChangeSourceLevelUp = "level-up"
	ExperienceChangeSourceOther   = "other"

	TerritoryHourChangeSourceSeize = "seize"
	TerritoryHourChangeSourceOther = "other"

	FreePointsChangeSourceQuestCompletion = "quest"
	FreePointsChangeSourceWheelEffect     = "wheel-effect"
	FreePointsChangeSourceBaseTeleport    = "base-teleport"
	FreePointsChangeSourceSandStorm       = "sandstorm"
	FreePointsChangeSourceOther           = "other"

	TerritoryPointChangeSourceObtaining = "territory-obtaining"
	TerritoryPointChangeSourceLoss      = "territory-loss"
	TerritoryPointChangeSourceOther     = "other"
)

var ExperienceChangeSourceSlice = []string{ExperienceChangeSourceLevelUp, ExperienceChangeSourceOther}

var TerritoryHourChangeSourceSlice = []string{TerritoryHourChangeSourceSeize, TerritoryHourChangeSourceOther}

var FreePointsChangeSourceSlice = []string{
	FreePointsChangeSourceQuestCompletion,
	FreePointsChangeSourceWheelEffect,
	FreePointsChangeSourceBaseTeleport,
	FreePointsChangeSourceSandStorm,
	FreePointsChangeSourceOther,
}

var TerritoryPointChangeSourceSlice = []string{
	TerritoryPointChangeSourceObtaining,
	TerritoryPointChangeSourceLoss,
	TerritoryPointChangeSourceOther,
}
