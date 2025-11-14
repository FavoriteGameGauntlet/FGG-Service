package controller

import (
	"FGG-Service/api"
	"FGG-Service/auth_service"
	"FGG-Service/game_service"
	"FGG-Service/timer_service"
	"context"
)

// GetCurrentTimer (GET /users/{userId}/timers/current)
func (Server) GetCurrentTimer(ctx context.Context, _ api.GetCurrentTimerRequestObject) (api.GetCurrentTimerResponseObject, error) {
	sessionId, ok := ctx.Value("session_id").(string)

	if !ok {
		return api.GetCurrentTimer401JSONResponse{Code: api.NOACTIVESESSION}, nil
	}

	userId, err := auth_service.GetUserId(sessionId)

	if err != nil {
		return api.GetCurrentTimer500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	game, err := game_service.GetCurrentGame(userId)

	if err != nil {
		return api.GetCurrentTimer500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	if game == nil {
		return api.GetCurrentTimer404JSONResponse{Code: api.GAMENOTFOUND}, nil
	}

	timer, err := timer_service.GetOrCreateCurrentTimer(userId, game.Id)

	if err != nil {
		return api.GetCurrentTimer500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	timerDto := ConvertTimerToDto(timer)

	return api.GetCurrentTimer200JSONResponse(*timerDto), nil
}

func ConvertTimerToDto(timer *timer_service.Timer) *api.Timer {
	return &api.Timer{
		DurationInS:      timer.DurationInS,
		RemainingTimeInS: timer.RemainingTimeInS,
		State:            api.TimerState(timer.State),
		TimerActionDate:  timer.TimerActionDate,
	}
}

// PauseCurrentTimer (POST /users/{userId}/timers/current/pause)
func (Server) PauseCurrentTimer(ctx context.Context, _ api.PauseCurrentTimerRequestObject) (api.PauseCurrentTimerResponseObject, error) {
	sessionId, ok := ctx.Value("session_id").(string)

	if !ok {
		return api.PauseCurrentTimer401JSONResponse{Code: api.NOACTIVESESSION}, nil
	}

	userId, err := auth_service.GetUserId(sessionId)

	if err != nil {
		return api.PauseCurrentTimer500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	doesExist, err := game_service.CheckIfCurrentGameExists(userId)

	if err != nil {
		return api.PauseCurrentTimer500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	if !doesExist {
		return api.PauseCurrentTimer404JSONResponse{Code: api.GAMENOTFOUND}, nil
	}

	doesExist, err = timer_service.CheckIfCurrentTimerExists(userId)

	if err != nil {
		return api.PauseCurrentTimer500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	if !doesExist {
		return api.PauseCurrentTimer404JSONResponse{Code: api.TIMERNOTFOUND}, nil
	}

	timerAction, err := timer_service.PauseCurrentTimer(userId)

	if err != nil {
		return api.PauseCurrentTimer500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	if timerAction == nil {
		return api.PauseCurrentTimer409JSONResponse{Code: api.TIMERINCORRECTSTATE}, nil
	}

	timerActionDto := ConvertTimerActionTo(timerAction)

	return api.PauseCurrentTimer200JSONResponse(*timerActionDto), nil
}

func ConvertTimerActionTo(timerAction *timer_service.TimerAction) *api.TimerAction {
	return &api.TimerAction{
		Type:             api.TimerActionType(timerAction.Action),
		RemainingTimeInS: timerAction.RemainingTimeInS,
	}
}

// StartCurrentTimer (POST /users/{userId}/timers/current/start)
func (Server) StartCurrentTimer(ctx context.Context, _ api.StartCurrentTimerRequestObject) (api.StartCurrentTimerResponseObject, error) {
	sessionId, ok := ctx.Value("session_id").(string)

	if !ok {
		return api.StartCurrentTimer401JSONResponse{Code: api.NOACTIVESESSION}, nil
	}

	userId, err := auth_service.GetUserId(sessionId)

	if err != nil {
		return api.StartCurrentTimer500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	doesExist, err := game_service.CheckIfCurrentGameExists(userId)

	if err != nil {
		return api.StartCurrentTimer500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	if !doesExist {
		return api.StartCurrentTimer404JSONResponse{Code: api.GAMENOTFOUND}, nil
	}

	doesExist, err = timer_service.CheckIfCurrentTimerExists(userId)

	if err != nil {
		return api.StartCurrentTimer500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	if !doesExist {
		return api.StartCurrentTimer404JSONResponse{Code: api.TIMERNOTFOUND}, nil
	}

	timerAction, err := timer_service.StartCurrentTimer(userId)

	if err != nil {
		return api.StartCurrentTimer500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	if timerAction == nil {
		return api.StartCurrentTimer409JSONResponse{Code: api.TIMERINCORRECTSTATE}, nil
	}

	timerActionDto := ConvertTimerActionTo(timerAction)

	return api.StartCurrentTimer200JSONResponse(*timerActionDto), nil
}
