package timer

type TimerState int

const (
	TimerStateCreated TimerState = iota
	TimerStateRunning
	TimerStatePaused
	TimerStateFinished
)
