package controller

import (
	"FGG-Service/api"
	"FGG-Service/common"
	"FGG-Service/timer_service"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetCurrentTimer (GET /timers/current)
func (Server) GetCurrentTimer(ctx echo.Context) error {
	userId, err := GetUserId(ctx)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	timer, err := timer_service.GetOrCreateCurrentTimer(userId)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	timerDto := ConvertTimerToDto(timer)

	return ctx.JSON(http.StatusOK, timerDto)
}

func ConvertTimerToDto(timer common.Timer) api.Timer {
	return api.Timer{
		Duration:        common.DurationToISO8601(timer.Duration),
		RemainingTime:   common.DurationToISO8601(timer.RemainingTime),
		State:           api.TimerState(timer.State),
		TimerActionDate: timer.TimerActionDate,
	}
}

// PauseCurrentTimer (POST /timers/current/pause)
func (Server) PauseCurrentTimer(ctx echo.Context) error {
	userId, err := GetUserId(ctx)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	timerAction, err := timer_service.PauseCurrentTimer(userId)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	timerActionDto := ConvertTimerActionToDto(timerAction)

	return ctx.JSON(http.StatusOK, timerActionDto)
}

func ConvertTimerActionToDto(timerAction common.TimerAction) api.TimerAction {
	return api.TimerAction{
		Type:          api.TimerActionType(timerAction.Type),
		RemainingTime: common.DurationToISO8601(timerAction.RemainingTime),
	}
}

// StartCurrentTimer (POST /timers/current/start)
func (Server) StartCurrentTimer(ctx echo.Context) error {
	userId, err := GetUserId(ctx)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	timerAction, err := timer_service.StartCurrentTimer(userId)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	timerActionDto := ConvertTimerActionToDto(timerAction)

	return ctx.JSON(http.StatusOK, timerActionDto)
}
