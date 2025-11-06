package effect_service

import (
	"time"
)

type Effect struct {
	Name        *string
	Description *string
	CreateDate  time.Time
	RollDate    *time.Time
	GameName    *string
}

type Effects = []Effect
