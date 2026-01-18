package common

import (
	"time"
)

type Timer struct {
	DurationInS      int
	Id               int
	RemainingTimeInS int
	State            TimerStateType
	TimerActionDate  *time.Time
}

type TimerStateType string

const (
	TimerStateCreated  TimerStateType = "created"
	TimerStateFinished TimerStateType = "finished"
	TimerStatePaused   TimerStateType = "paused"
	TimerStateRunning  TimerStateType = "running"
)

type TimerAction struct {
	Type             TimerActionType
	Id               int
	RemainingTimeInS int
}

type TimerActionType string

const (
	TimerActionPause TimerActionType = "pause"
	TimerActionStart TimerActionType = "start"
	TimerActionStop  TimerActionType = "stop"
)
