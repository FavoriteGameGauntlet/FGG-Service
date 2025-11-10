package timer_service

import (
	"time"

	"github.com/google/uuid"
)

type Timer struct {
	DurationInS      int
	Id               uuid.UUID
	RemainingTimeInS int
	State            TimerStateType
	TimerActionDate  *time.Time
}

const DefaultTimerDurationInS = 30

type TimerStateType string

const (
	TimerStateCreated  TimerStateType = "created"
	TimerStateFinished TimerStateType = "finished"
	TimerStatePaused   TimerStateType = "paused"
	TimerStateRunning  TimerStateType = "running"
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
