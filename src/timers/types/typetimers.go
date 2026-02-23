package typetimers

import (
	"time"
)

type Timer struct {
	Duration       time.Duration
	Id             int
	RemainingTime  time.Duration
	State          TimerStateType
	LastActionDate time.Time
}

type TimerStateType string

const (
	TimerStateCreated  TimerStateType = "created"
	TimerStateFinished TimerStateType = "finished"
	TimerStatePaused   TimerStateType = "paused"
	TimerStateRunning  TimerStateType = "running"
)
