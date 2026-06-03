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

const (
	ExperienceChangeSourceLevelUp = "level-up"
	ExperienceChangeSourceOther   = "other"
)

var ExperienceChangeSourceSlice = []string{ExperienceChangeSourceLevelUp, ExperienceChangeSourceOther}
