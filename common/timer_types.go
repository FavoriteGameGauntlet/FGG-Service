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
	Action           TimerActionAction
	Id               int
	RemainingTimeInS int
}

type TimerActionAction string

const (
	TimerActionPause TimerActionAction = "pause"
	TimerActionStart TimerActionAction = "start"
	TimerActionStop  TimerActionAction = "stop"
)
