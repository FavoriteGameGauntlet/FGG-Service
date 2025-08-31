package user

import (
	"FavoriteGameGauntlet/database"
	"FavoriteGameGauntlet/timer"
)

type User struct {
	id int64
}

const (
	FindUserCommand = `
		SELECT Id
		FROM Users
		WHERE Name = $userName`
	GetCurrentTimer = `
		SELECT Id
		FROM Timers
		WHERE UserId = $userId AND State != 3`
)

func FindUser(userName string) *User {
	row := database.QueryRow(FindUserCommand, userName)

	var userId int64
	err := row.Scan(&userId)

	if err != nil {
		panic(err)
	}

	return &User{
		id: userId,
	}
}

func (user *User) GetCurrentTimer() *timer.Timer {
	row := database.QueryRow(GetCurrentTimer, user.id)

	var timerId int64
	err := row.Scan(&timerId)

	if err == nil {
		return timer.NewTimerFromId(timerId)
	}

	return timer.NewTimer(user.id)
}
