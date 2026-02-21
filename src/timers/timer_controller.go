package timers

import (
	"FGG-Service/api"
	"FGG-Service/src/common"
	"FGG-Service/src/controller"
	"FGG-Service/src/timers/timer_service"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetCurrentTimer (GET /timers/current)
func (controller.Server) GetCurrentTimer(ctx echo.Context) error {
	userId, err := controller.GetUserId(ctx)

	if err != nil {
		return controller.SendJSONErrorResponse(ctx, err)
	}

	timer, err := timer_service.GetOrCreateCurrentTimer(userId)

	if err != nil {
		return controller.SendJSONErrorResponse(ctx, err)
	}

	timerDto := ConvertTimerToDto(timer)

	return ctx.JSON(http.StatusOK, timerDto)
}

func ConvertTimerToDto(timer common.Timer) api.Timer {
	return api.Timer{
		Duration:       common.DurationToISO8601(timer.Duration),
		RemainingTime:  common.DurationToISO8601(timer.RemainingTime),
		State:          api.TimerState(timer.State),
		LastActionDate: timer.LastActionDate,
	}
}

// PauseCurrentTimer (POST /timers/current/pause)
func (controller.Server) PauseCurrentTimer(ctx echo.Context) error {
	userId, err := controller.GetUserId(ctx)

	if err != nil {
		return controller.SendJSONErrorResponse(ctx, err)
	}

	timer, err := timer_service.PauseCurrentTimer(userId)

	if err != nil {
		return controller.SendJSONErrorResponse(ctx, err)
	}

	timerActionDto := ConvertTimerToDto(timer)

	return ctx.JSON(http.StatusOK, timerActionDto)
}

// StartCurrentTimer (POST /timers/current/start)
func (controller.Server) StartCurrentTimer(ctx echo.Context) error {
	userId, err := controller.GetUserId(ctx)

	if err != nil {
		return controller.SendJSONErrorResponse(ctx, err)
	}

	timer, err := timer_service.StartCurrentTimer(userId)

	if err != nil {
		return controller.SendJSONErrorResponse(ctx, err)
	}

	timerActionDto := ConvertTimerToDto(timer)

	return ctx.JSON(http.StatusOK, timerActionDto)
}
