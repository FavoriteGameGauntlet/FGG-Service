package controller

import (
	"FGG-Service/api"
	"FGG-Service/auth_service"
	"FGG-Service/game_service"
	"FGG-Service/timer_service"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetCurrentTimer (GET /timers/current)
func (Server) GetCurrentTimer(ctx echo.Context) error {
	cookie, err := ctx.Cookie("sessionId")

	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, api.NotAuthorizedError{Code: api.NOACTIVESESSION})
	}

	sessionId := cookie.Value

	userId, err := auth_service.GetUserId(sessionId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	game, err := game_service.GetCurrentGame(userId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	if game == nil {
		return ctx.JSON(http.StatusNotFound, api.NotFoundError{Code: api.GAMENOTFOUND})
	}

	timer, err := timer_service.GetOrCreateCurrentTimer(userId, game.Id)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	timerDto := ConvertTimerToDto(timer)

	return ctx.JSON(http.StatusOK, *timerDto)
}

func ConvertTimerToDto(timer *timer_service.Timer) *api.Timer {
	return &api.Timer{
		DurationInS:      timer.DurationInS,
		RemainingTimeInS: timer.RemainingTimeInS,
		State:            api.TimerState(timer.State),
		TimerActionDate:  timer.TimerActionDate,
	}
}

// PauseCurrentTimer (POST /timers/current/pause)
func (Server) PauseCurrentTimer(ctx echo.Context) error {
	cookie, err := ctx.Cookie("sessionId")

	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, api.NotAuthorizedError{Code: api.NOACTIVESESSION})
	}

	sessionId := cookie.Value

	userId, err := auth_service.GetUserId(sessionId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	doesExist, err := game_service.CheckIfCurrentGameExists(userId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	if !doesExist {
		return ctx.JSON(http.StatusNotFound, api.NotFoundError{Code: api.GAMENOTFOUND})
	}

	doesExist, err = timer_service.CheckIfCurrentTimerExists(userId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	if !doesExist {
		return ctx.JSON(http.StatusNotFound, api.NotFoundError{Code: api.TIMERNOTFOUND})
	}

	timerAction, err := timer_service.PauseCurrentTimer(userId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	if timerAction == nil {
		return ctx.JSON(http.StatusConflict, api.ConflictError{Code: api.TIMERINCORRECTSTATE})
	}

	timerActionDto := ConvertTimerActionToDto(timerAction)

	return ctx.JSON(http.StatusOK, *timerActionDto)
}

func ConvertTimerActionToDto(timerAction *timer_service.TimerAction) *api.TimerAction {
	return &api.TimerAction{
		Type:             api.TimerActionType(timerAction.Action),
		RemainingTimeInS: timerAction.RemainingTimeInS,
	}
}

// StartCurrentTimer (POST /timers/current/start)
func (Server) StartCurrentTimer(ctx echo.Context) error {
	cookie, err := ctx.Cookie("sessionId")

	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, api.NotAuthorizedError{Code: api.NOACTIVESESSION})
	}

	sessionId := cookie.Value

	userId, err := auth_service.GetUserId(sessionId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	doesExist, err := game_service.CheckIfCurrentGameExists(userId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	if !doesExist {
		return ctx.JSON(http.StatusNotFound, api.NotFoundError{Code: api.GAMENOTFOUND})
	}

	doesExist, err = timer_service.CheckIfCurrentTimerExists(userId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	if !doesExist {
		return ctx.JSON(http.StatusNotFound, api.NotFoundError{Code: api.TIMERNOTFOUND})
	}

	timerAction, err := timer_service.StartCurrentTimer(userId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	if timerAction == nil {
		return ctx.JSON(http.StatusConflict, api.ConflictError{Code: api.TIMERINCORRECTSTATE})
	}

	timerActionDto := ConvertTimerActionToDto(timerAction)

	return ctx.JSON(http.StatusOK, *timerActionDto)
}
