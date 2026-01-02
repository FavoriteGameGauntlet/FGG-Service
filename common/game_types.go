package common

import (
	"time"
)

type Game struct {
	Id         int
	Name       string
	State      GameState
	TimeSpent  time.Duration
	FinishDate *time.Time
}

type Games = []Game

type GameState string

const (
	GameStateCancelled GameState = "cancelled"
	GameStateFinished  GameState = "finished"
	GameStateStarted   GameState = "started"
)

type UnplayedGame struct {
	Id     int
	GameId int
	Name   string
}

type UnplayedGames = []UnplayedGame
