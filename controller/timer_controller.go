package controller

import (
	"FGG-Service/api"
	"FGG-Service/game_service"
	"FGG-Service/timer_service"
	"FGG-Service/user_service"
	"context"
)

// GetCurrentTimer (GET /users/{userId}/timers/current)
func (Server) GetCurrentTimer(_ context.Context, request api.GetCurrentTimerRequestObject) (api.GetCurrentTimerResponseObject, error) {
	doesUserExist, err := user_service.CheckIfUserExistsById(request.UserId)

	if err != nil {
		return api.GetCurrentTimer503JSONResponse{Code: api.CHECKUSER, Message: err.Error()}, nil
	}

	if !doesUserExist {
		return api.GetCurrentTimer404JSONResponse{Code: api.USERNOTFOUND}, nil
	}

	game, err := game_service.GetCurrentGame(request.UserId)

	if err != nil {
		return api.GetCurrentTimer503JSONResponse{Code: api.CHECKCURRENTGAME, Message: err.Error()}, nil
	}

	if game == nil {
		return api.GetCurrentTimer404JSONResponse{Code: api.GAMENOTFOUND}, nil
	}

	timer, err := timer_service.GetOrCreateCurrentTimer(request.UserId, game.Id)

	if err != nil {
		return api.GetCurrentTimer503JSONResponse{Code: api.GETCURRENTTIMER, Message: err.Error()}, nil
	}

	timerDto := ConvertTimerToDto(timer)

	return api.GetCurrentTimer200JSONResponse(*timerDto), nil
}

func ConvertTimerToDto(timer *timer_service.Timer) *api.TimerDto {
	return &api.TimerDto{
		DurationInS:      timer.DurationInS,
		RemainingTimeInS: timer.RemainingTimeInS,
		State:            api.TimerDtoState(timer.State),
		TimerActionDate:  timer.TimerActionDate,
	}
}

// PauseCurrentTimer (POST /users/{userId}/timers/current/pause)
func (Server) PauseCurrentTimer(_ context.Context, request api.PauseCurrentTimerRequestObject) (api.PauseCurrentTimerResponseObject, error) {
	doesUserExist, err := user_service.CheckIfUserExistsById(request.UserId)

	if err != nil {
		return api.PauseCurrentTimer503JSONResponse{Code: api.CHECKUSER, Message: err.Error()}, nil
	}

	if !doesUserExist {
		return api.PauseCurrentTimer404JSONResponse{Code: api.USERNOTFOUND}, nil
	}

	doesCurrentGameExist, err := game_service.CheckIfCurrentGameExists(request.UserId)

	if err != nil {
		return api.PauseCurrentTimer503JSONResponse{Code: api.CHECKCURRENTGAME, Message: err.Error()}, nil
	}

	if !doesCurrentGameExist {
		return api.PauseCurrentTimer404JSONResponse{Code: api.GAMENOTFOUND}, nil
	}

	doesCurrentTimerExist, err := timer_service.CheckIfCurrentTimerExists(request.UserId)

	if err != nil {
		return api.PauseCurrentTimer503JSONResponse{Code: api.CHECKCURRENTTIMER, Message: err.Error()}, nil
	}

	if !doesCurrentTimerExist {
		return api.PauseCurrentTimer404JSONResponse{Code: api.TIMERNOTFOUND}, nil
	}

	timerAction, err := timer_service.PauseCurrentTimer(request.UserId)

	if err != nil {
		return api.PauseCurrentTimer503JSONResponse{Code: api.PAUSECURRENTTIMER, Message: err.Error()}, nil
	}

	if timerAction == nil {
		return api.PauseCurrentTimer409JSONResponse{Code: api.TIMERINCORRECTSTATE}, nil
	}

	timerActionDto := ConvertTimerActionToDto(timerAction)

	return api.PauseCurrentTimer200JSONResponse(*timerActionDto), nil
}

func ConvertTimerActionToDto(timerAction *timer_service.TimerAction) *api.TimerActionDto {
	return &api.TimerActionDto{
		Action:           api.TimerActionDtoAction(timerAction.Action),
		RemainingTimeInS: timerAction.RemainingTimeInS,
	}
}

// StartCurrentTimer (POST /users/{userId}/timers/current/start)
func (Server) StartCurrentTimer(_ context.Context, request api.StartCurrentTimerRequestObject) (api.StartCurrentTimerResponseObject, error) {
	doesUserExist, err := user_service.CheckIfUserExistsById(request.UserId)

	if err != nil {
		return api.StartCurrentTimer503JSONResponse{Code: api.CHECKUSER, Message: err.Error()}, nil
	}

	if !doesUserExist {
		return api.StartCurrentTimer404JSONResponse{Code: api.USERNOTFOUND}, nil
	}

	doesCurrentGameExist, err := game_service.CheckIfCurrentGameExists(request.UserId)

	if err != nil {
		return api.StartCurrentTimer503JSONResponse{Code: api.CHECKCURRENTGAME, Message: err.Error()}, nil
	}

	if !doesCurrentGameExist {
		return api.StartCurrentTimer404JSONResponse{Code: api.GAMENOTFOUND}, nil
	}

	doesCurrentTimerExist, err := timer_service.CheckIfCurrentTimerExists(request.UserId)

	if err != nil {
		return api.StartCurrentTimer503JSONResponse{Code: api.CHECKCURRENTTIMER, Message: err.Error()}, nil
	}

	if !doesCurrentTimerExist {
		return api.StartCurrentTimer404JSONResponse{Code: api.TIMERNOTFOUND}, nil
	}

	timerAction, err := timer_service.StartCurrentTimer(request.UserId)

	if err != nil {
		return api.StartCurrentTimer503JSONResponse{Code: api.STARTCURRENTTIMER, Message: err.Error()}, nil
	}

	if timerAction == nil {
		return api.StartCurrentTimer409JSONResponse{Code: api.TIMERINCORRECTSTATE}, nil
	}

	timerActionDto := ConvertTimerActionToDto(timerAction)

	return api.StartCurrentTimer200JSONResponse(*timerActionDto), nil
}
