package controller

import (
	"FavoriteGameGauntlet/api"
	"context"
)

// (GET /users/{userId}/effects/has-roll)
func (Server) GetUsersUserIdEffectsHasRoll(ctx context.Context, request api.GetUsersUserIdEffectsHasRollRequestObject) (api.GetUsersUserIdEffectsHasRollResponseObject, error) {
	return api.GetUsersUserIdEffectsHasRoll200JSONResponse(true), nil
}

// (GET /users/{userId}/effects/history)
func (Server) GetUsersUserIdEffectsHistory(ctx context.Context, request api.GetUsersUserIdEffectsHistoryRequestObject) (api.GetUsersUserIdEffectsHistoryResponseObject, error) {
	return api.GetUsersUserIdEffectsHistory200JSONResponse{}, nil
}

// (POST /users/{userId}/effects/roll)
func (Server) PostUsersUserIdEffectsRoll(ctx context.Context, request api.PostUsersUserIdEffectsRollRequestObject) (api.PostUsersUserIdEffectsRollResponseObject, error) {
	return api.PostUsersUserIdEffectsRoll200JSONResponse{}, nil
}
