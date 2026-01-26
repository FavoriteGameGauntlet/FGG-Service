package common

import (
	"time"
)

type Timer struct {
	Duration        time.Duration
	Id              int
	RemainingTime   time.Duration
	State           TimerStateType
	TimerActionDate *time.Time
}

type TimerStateType string

const (
	TimerStateCreated  TimerStateType = "created"
	TimerStateFinished TimerStateType = "finished"
	TimerStatePaused   TimerStateType = "paused"
	TimerStateRunning  TimerStateType = "running"
)

type TimerAction struct {
	Type          TimerActionType
	Id            int
	RemainingTime time.Duration
}

type TimerActionType string

const (
	TimerActionPause TimerActionType = "pause"
	TimerActionStart TimerActionType = "start"
	TimerActionStop  TimerActionType = "stop"
)
