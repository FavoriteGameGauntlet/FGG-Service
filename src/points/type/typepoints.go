package typepoints

import "time"

type PointChangeResult struct {
	ActualChangeValue  int
	ChangeDate         time.Time
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
)

var ExperienceChangeSourceSlice = []string{ExperienceChangeSourceLevelUp, ExperienceChangeSourceOther}

var TerritoryHoursChangeSourceSlice = []string{TerritoryHoursChangeSourceSeize, TerritoryHoursChangeSourceOther}
