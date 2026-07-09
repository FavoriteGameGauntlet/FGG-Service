package typewheeleffects

import (
	typepoints "FGG-Service/src/points/type"
	"time"
)

type WheelEffect struct {
	Id          int
	Name        string
	Description *string
}

type WheelEffects = []WheelEffect

type RolledWheelEffect struct {
	Id          int
	Name        string
	Description *string
	RollDate    time.Time
	Position    int
	IsApplied   bool
}

type RolledWheelEffects = []RolledWheelEffect

type RolledWheelEffectHistory struct {
	Id          int
	Name        string
	Description *string
	RollDate    time.Time
}

type RolledWheelEffectHistories = []RolledWheelEffectHistory

type WheelEffectRollApply struct {
	PointChangeByUserIds typepoints.PointChangeByUserIds
	WheelEffectName      string
}
