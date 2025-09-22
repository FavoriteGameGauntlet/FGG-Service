package controller

import (
	"FavoriteGameGauntlet/api"
	"context"
)

// (GET /users/{userId}/timers/current)
func (Server) GetUsersUserIdTimersCurrent(ctx context.Context, request api.GetUsersUserIdTimersCurrentRequestObject) (api.GetUsersUserIdTimersCurrentResponseObject, error) {
	return api.GetUsersUserIdTimersCurrent200JSONResponse{}, nil
}

// (POST /users/{userId}/timers/current/pause)
func (Server) PostUsersUserIdTimersCurrentPause(ctx context.Context, request api.PostUsersUserIdTimersCurrentPauseRequestObject) (api.PostUsersUserIdTimersCurrentPauseResponseObject, error) {
	return api.PostUsersUserIdTimersCurrentPause200JSONResponse(0), nil
}

// (POST /users/{userId}/timers/current/start)
func (Server) PostUsersUserIdTimersCurrentStart(ctx context.Context, request api.PostUsersUserIdTimersCurrentStartRequestObject) (api.PostUsersUserIdTimersCurrentStartResponseObject, error) {
	return api.PostUsersUserIdTimersCurrentStart200JSONResponse(0), nil
}
