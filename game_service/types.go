package game_service

import (
    "github.com/google/uuid"
)

type Game struct {
    Id    uuid.UUID
    Link  *string
    Name  string
    State GameState
}

type GameState string

const (
    GameStateCancelled GameState = "cancelled"
    GameStateFinished  GameState = "finished"
    GameStateStarted   GameState = "started"
)

type UnplayedGame struct {
    Id   uuid.UUID
    Link *string
    Name string
}

type UnplayedGames = []UnplayedGame
