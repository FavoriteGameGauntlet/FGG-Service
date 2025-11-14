package game_service

import (
	"time"
)

type Game struct {
	Id         int
	Link       *string
	Name       string
	State      GameState
	HourCount  int
	FinishDate *time.Time
}

type Games = []Game

type GameState string

const (
	GameStateCancelled GameState = "cancelled"
	GameStateFinished  GameState = "finished"
	GameStateStarted   GameState = "started"
)

const (
	CancellingGamePenaltyDiceCount = 2
	HourCountForDice               = 6
)

type UnplayedGame struct {
	GameId int
	Link   *string
	Name   string
}

type UnplayedGames = []UnplayedGame
