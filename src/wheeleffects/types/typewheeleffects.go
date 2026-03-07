package typewheeleffects

import "time"

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
	Position    *int
	IsApplied   *bool
}

type RolledWheelEffects = []RolledWheelEffect
