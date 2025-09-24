package timer_service

import (
	"FGG-Service/database"
	"time"
)

type Timer struct {
	id int64
}

const (
	CreateTimerCommand = `
		INSERT INTO Timers (UserId, DurationInS)
		VALUES ($userId, $duration)`
	DoTimerActionCommand = `
		INSERT INTO TimerActions (TimerId, Action)
		VALUES ($timerId, $action)`
)

func NewTimerFromId(timerId int64) *Timer {
	return &Timer{
		id: timerId,
	}
}

func NewTimer(userId int64) *Timer {
	const threeHours = 3 * time.Hour
	result := database.Exec(CreateTimerCommand, userId, threeHours.Seconds())

	timerId, err := result.LastInsertId()

	if err != nil {
		panic(err)
	}

	return &Timer{
		id: timerId,
	}
}

func (timer *Timer) StartTimer() *Timer {
	database.Exec(DoTimerActionCommand, timer.id, TimerActionStart)

	return timer
}

func (timer *Timer) PauseTimer() *Timer {
	database.Exec(DoTimerActionCommand, timer.id, TimerActionPause)

	return timer
}

func (timer *Timer) GetTimerId() int64 {
	return timer.id
}
