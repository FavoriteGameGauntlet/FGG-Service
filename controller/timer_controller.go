package controller

import (
	"FGG-Service/api"
	"FGG-Service/game_service"
	"FGG-Service/timer_service"
	"FGG-Service/user_service"
	"context"
)

// (GET /users/{userId}/timers/current)
func (Server) GetCurrentTimer(ctx context.Context, request api.GetCurrentTimerRequestObject) (api.GetCurrentTimerResponseObject, error) {
	doesUserExist, err := user_service.CheckIfUserExistsById(request.UserId)

	if err != nil {
		return api.GetCurrentTimer503JSONResponse{api.DATABASEQUERY, err.Error()}, nil
	}

	if !*doesUserExist {
		return api.GetCurrentTimer404JSONResponse{Code: api.USERNOTFOUND}, nil
	}

	game, err := game_service.GetCurrentGame(request.UserId)

	if err != nil {
		return api.GetCurrentTimer503JSONResponse{api.DATABASEQUERY, err.Error()}, nil
	}

	if game == nil {
		return api.GetCurrentTimer404JSONResponse{Code: api.GAMENOTFOUND}, nil
	}

	timer, err := timer_service.GetOrCreateCurrentTimer(request.UserId, game.Id)

	if err != nil {
		return api.GetCurrentTimer503JSONResponse{api.DATABASEQUERY, err.Error()}, nil
	}

	return api.GetCurrentTimer200JSONResponse(*timer), nil
}

// (POST /users/{userId}/timers/current/pause)
func (Server) PauseCurrentTimer(ctx context.Context, request api.PauseCurrentTimerRequestObject) (api.PauseCurrentTimerResponseObject, error) {
	return api.PauseCurrentTimer200JSONResponse(0), nil
}

// (POST /users/{userId}/timers/current/start)
func (Server) StartCurrentTimer(ctx context.Context, request api.StartCurrentTimerRequestObject) (api.StartCurrentTimerResponseObject, error) {
	return api.StartCurrentTimer200JSONResponse(0), nil
}
