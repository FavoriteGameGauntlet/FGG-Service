package typegames

import "time"

type CurrentGame struct {
	Id         int
	Name       string
	State      CurrentGameState
	TimeSpent  time.Duration
	FinishDate *time.Time
}

type CurrentGames = []CurrentGame

type CurrentGameState string

const (
	GameStateCancelled CurrentGameState = "cancelled"
	GameStateFinished  CurrentGameState = "finished"
	GameStateStarted   CurrentGameState = "started"
)

type WishlistGame struct {
	Id     int
	GameId int
	Name   string
}

type WishlistGames = []WishlistGame
