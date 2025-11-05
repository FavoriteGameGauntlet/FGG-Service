package timer_service

import (
	"github.com/google/uuid"
	"time"
)

type Timer struct {
	DurationInS      int
	Id               uuid.UUID
	RemainingTimeInS int
	State            TimerState
	TimerActionDate  *time.Time
}

type TimerState string

const (
	TimerStateCreated  TimerState = "created"
	TimerStateFinished TimerState = "finished"
	TimerStatePaused   TimerState = "paused"
	TimerStateRolled   TimerState = "rolled"
	TimerStateRunning  TimerState = "running"
)

type TimerAction struct {
	Action           TimerActionAction
	Id               uuid.UUID
	RemainingTimeInS int
}

type TimerActionAction string

const (
	TimerActionPause TimerActionAction = "pause"
	TimerActionStart TimerActionAction = "start"
	TimerActionStop  TimerActionAction = "stop"
)
