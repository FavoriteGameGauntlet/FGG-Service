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
		return api.GetCurrentTimer503JSONResponse{api.DATABASEQUERY + "1", err.Error()}, nil
	}

	if !*doesUserExist {
		return api.GetCurrentTimer404JSONResponse{Code: api.USERNOTFOUND}, nil
	}

	doesCurrentGameExist, err := game_service.CheckIfCurrentGameExists(request.UserId)

	if err != nil {
		return api.GetCurrentTimer503JSONResponse{api.DATABASEQUERY + "2", err.Error()}, nil
	}

	if !*doesCurrentGameExist {
		return api.GetCurrentTimer404JSONResponse{Code: api.GAMENOTFOUND}, nil
	}

	timer, err := timer_service.GetOrCreateCurrentTimer(request.UserId)

	if err != nil {
		return api.GetCurrentTimer503JSONResponse{api.DATABASEQUERY + "3", err.Error()}, nil
	}

	return api.GetCurrentTimer200JSONResponse(*timer), nil
}

// (POST /users/{userId}/timers/current/pause)
func (Server) PauseCurrentTimer(ctx context.Context, request api.PauseCurrentTimerRequestObject) (api.PauseCurrentTimerResponseObject, error) {
	doesUserExist, err := user_service.CheckIfUserExistsById(request.UserId)

	if err != nil {
		return api.PauseCurrentTimer503JSONResponse{api.DATABASEQUERY + "1", err.Error()}, nil
	}

	if !*doesUserExist {
		return api.PauseCurrentTimer404JSONResponse{Code: api.USERNOTFOUND}, nil
	}

	doesCurrentGameExist, err := game_service.CheckIfCurrentGameExists(request.UserId)

	if err != nil {
		return api.PauseCurrentTimer503JSONResponse{api.DATABASEQUERY + "2", err.Error()}, nil
	}

	if !*doesCurrentGameExist {
		return api.PauseCurrentTimer404JSONResponse{Code: api.GAMENOTFOUND}, nil
	}

	doesCurrentTimerExist, err := timer_service.CheckIfCurrentTimerExists(request.UserId)

	if err != nil {
		return api.PauseCurrentTimer503JSONResponse{api.DATABASEQUERY + "3", err.Error()}, nil
	}

	if !*doesCurrentTimerExist {
		return api.PauseCurrentTimer404JSONResponse{Code: api.TIMERNOTFOUND}, nil
	}

	timerAction, err := timer_service.PauseCurrentTimer(request.UserId)

	if err != nil {
		return api.PauseCurrentTimer503JSONResponse{api.DATABASEQUERY + "4", err.Error()}, nil
	}

	if timerAction == nil {
		return api.PauseCurrentTimer409JSONResponse{Code: api.TIMERINCORRECTSTATE}, nil
	}

	return api.PauseCurrentTimer200JSONResponse(*timerAction), nil
}

// (POST /users/{userId}/timers/current/start)
func (Server) StartCurrentTimer(ctx context.Context, request api.StartCurrentTimerRequestObject) (api.StartCurrentTimerResponseObject, error) {
	doesUserExist, err := user_service.CheckIfUserExistsById(request.UserId)

	if err != nil {
		return api.StartCurrentTimer503JSONResponse{api.DATABASEQUERY + "1", err.Error()}, nil
	}

	if !*doesUserExist {
		return api.StartCurrentTimer404JSONResponse{Code: api.USERNOTFOUND}, nil
	}

	doesCurrentGameExist, err := game_service.CheckIfCurrentGameExists(request.UserId)

	if err != nil {
		return api.StartCurrentTimer503JSONResponse{api.DATABASEQUERY + "2", err.Error()}, nil
	}

	if !*doesCurrentGameExist {
		return api.StartCurrentTimer404JSONResponse{Code: api.GAMENOTFOUND}, nil
	}

	doesCurrentTimerExist, err := timer_service.CheckIfCurrentTimerExists(request.UserId)

	if err != nil {
		return api.StartCurrentTimer503JSONResponse{api.DATABASEQUERY + "3", err.Error()}, nil
	}

	if !*doesCurrentTimerExist {
		return api.StartCurrentTimer404JSONResponse{Code: api.TIMERNOTFOUND}, nil
	}

	timerAction, err := timer_service.StartCurrentTimer(request.UserId)

	if err != nil {
		return api.StartCurrentTimer503JSONResponse{api.DATABASEQUERY + "4", err.Error()}, nil
	}

	if timerAction == nil {
		return api.StartCurrentTimer409JSONResponse{Code: api.TIMERINCORRECTSTATE}, nil
	}

	return api.StartCurrentTimer200JSONResponse(*timerAction), nil
}
