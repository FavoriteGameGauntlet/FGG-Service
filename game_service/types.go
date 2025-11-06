package game_service

import (
	"time"

	"github.com/google/uuid"
)

type Game struct {
	Id         uuid.UUID
	Link       *string
	Name       string
	State      GameState
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
	GameId uuid.UUID
	Link   *string
	Name   string
}

type UnplayedGames = []UnplayedGame
