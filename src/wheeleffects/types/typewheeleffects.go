package typewheeleffects

import "time"

type WheelEffect struct {
	Id          int
	Name        string
	Description *string
}

type WheelEffects = []WheelEffect

type RolledWheelEffect struct {
	Name        string
	Description *string
	RollDate    time.Time
}

type RolledWheelEffects = []RolledWheelEffect
