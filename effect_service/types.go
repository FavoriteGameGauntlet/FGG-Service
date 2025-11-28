package effect_service

import "time"

type Effect struct {
	Id          int
	Name        string
	Description *string
}

type Effects = []Effect

type RolledEffect struct {
	Name        string
	Description *string
	RollDate    time.Time
}

type RolledEffects = []RolledEffect
