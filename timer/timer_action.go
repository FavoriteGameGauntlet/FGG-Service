package timer

type TimerAction int

const (
	TimerActionStart TimerAction = iota
	TimerActionPause
	TimerActionStop
)
