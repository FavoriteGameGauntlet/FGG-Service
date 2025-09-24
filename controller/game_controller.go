package controller

import (
	"FGG-Service/api"
	"context"
)

// (GET /users/{userId}/games/current)
func (Server) GetCurrentGame(ctx context.Context, request api.GetCurrentGameRequestObject) (api.GetCurrentGameResponseObject, error) {
	return api.GetCurrentGame200JSONResponse{}, nil
}

// (GET /users/{userId}/games/current/finish)
func (Server) FinishCurrentGame(ctx context.Context, request api.FinishCurrentGameRequestObject) (api.FinishCurrentGameResponseObject, error) {
	return api.FinishCurrentGame200Response{}, nil
}

// (GET /users/{userId}/games/roll)
func (Server) MakeGameRoll(ctx context.Context, request api.MakeGameRollRequestObject) (api.MakeGameRollResponseObject, error) {
	return api.MakeGameRoll200JSONResponse{}, nil
}

// (GET /users/{userId}/games/unplayed)
func (Server) GetUnplayedGames(ctx context.Context, request api.GetUnplayedGamesRequestObject) (api.GetUnplayedGamesResponseObject, error) {
	return api.GetUnplayedGames200JSONResponse{}, nil
}

// (POST /users/{userId}/games/unplayed)
func (Server) AddUnplayedGames(ctx context.Context, request api.AddUnplayedGamesRequestObject) (api.AddUnplayedGamesResponseObject, error) {
	return api.AddUnplayedGames200Response{}, nil
}
