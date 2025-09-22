package controller

import (
	"FavoriteGameGauntlet/api"
	"context"
)

// (GET /users/{userId}/games/current)
func (Server) GetUsersUserIdGamesCurrent(ctx context.Context, request api.GetUsersUserIdGamesCurrentRequestObject) (api.GetUsersUserIdGamesCurrentResponseObject, error) {
	return api.GetUsersUserIdGamesCurrent200JSONResponse{}, nil
}

// (GET /users/{userId}/games/current/finish)
func (Server) GetUsersUserIdGamesCurrentFinish(ctx context.Context, request api.GetUsersUserIdGamesCurrentFinishRequestObject) (api.GetUsersUserIdGamesCurrentFinishResponseObject, error) {
	return api.GetUsersUserIdGamesCurrentFinish200Response{}, nil
}

// (GET /users/{userId}/games/roll)
func (Server) GetUsersUserIdGamesRoll(ctx context.Context, request api.GetUsersUserIdGamesRollRequestObject) (api.GetUsersUserIdGamesRollResponseObject, error) {
	return api.GetUsersUserIdGamesRoll200JSONResponse{}, nil
}

// (GET /users/{userId}/games/unplayed)
func (Server) GetUsersUserIdGamesUnplayed(ctx context.Context, request api.GetUsersUserIdGamesUnplayedRequestObject) (api.GetUsersUserIdGamesUnplayedResponseObject, error) {
	return api.GetUsersUserIdGamesUnplayed200JSONResponse{}, nil
}

// (POST /users/{userId}/games/unplayed)
func (Server) PostUsersUserIdGamesUnplayed(ctx context.Context, request api.PostUsersUserIdGamesUnplayedRequestObject) (api.PostUsersUserIdGamesUnplayedResponseObject, error) {
	return api.PostUsersUserIdGamesUnplayed200Response{}, nil
}
