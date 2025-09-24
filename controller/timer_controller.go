package controller

import (
	"FGG-Service/api"
	"context"
)

// (GET /users/{userId}/timers/current)
func (Server) GetCurrentTimer(ctx context.Context, request api.GetCurrentTimerRequestObject) (api.GetCurrentTimerResponseObject, error) {
	return api.GetCurrentTimer200JSONResponse{}, nil
}

// (POST /users/{userId}/timers/current/pause)
func (Server) PauseCurrentTimer(ctx context.Context, request api.PauseCurrentTimerRequestObject) (api.PauseCurrentTimerResponseObject, error) {
	return api.PauseCurrentTimer200JSONResponse(0), nil
}

// (POST /users/{userId}/timers/current/start)
func (Server) StartCurrentTimer(ctx context.Context, request api.StartCurrentTimerRequestObject) (api.StartCurrentTimerResponseObject, error) {
	return api.StartCurrentTimer200JSONResponse(0), nil
}
