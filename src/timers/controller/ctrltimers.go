package ctrltimers

import (
	gentimers "FGG-Service/api/generated/timers"
	srvauth "FGG-Service/src/auth/service"
	"FGG-Service/src/common"
	"FGG-Service/src/timers/service"
	"FGG-Service/src/timers/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	Service     srvtimers.Service
	AuthService srvauth.Service
}

func NewController() Controller {
	return Controller{}
}

// GetOwnCurrentTimer (GET /timers/self/current)
func (c *Controller) GetOwnCurrentTimer(ctx echo.Context) error {
	userId, err := c.AuthService.GetUserId(ctx)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	timer, err := c.Service.GetOrCreateCurrentTimer(userId)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	timerDto := ConvertTimerToDto(timer)

	return ctx.JSON(http.StatusOK, timerDto)
}

func ConvertTimerToDto(timer typetimers.Timer) gentimers.Timer {
	return gentimers.Timer{
		Duration:       common.DurationToISO8601(timer.Duration),
		RemainingTime:  common.DurationToISO8601(timer.RemainingTime),
		State:          gentimers.TimerState(timer.State),
		LastActionDate: timer.LastActionDate,
	}
}

// PauseOwnCurrentTimer (POST /timers/self/current/pause)
func (c *Controller) PauseOwnCurrentTimer(ctx echo.Context) error {
	userId, err := c.AuthService.GetUserId(ctx)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	timer, err := c.Service.PauseCurrentTimer(userId)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	timerActionDto := ConvertTimerToDto(timer)

	return ctx.JSON(http.StatusOK, timerActionDto)
}

// StartOwnCurrentTimer (POST /timers/self/current/start)
func (c *Controller) StartOwnCurrentTimer(ctx echo.Context) error {
	userId, err := c.AuthService.GetUserId(ctx)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	timer, err := c.Service.StartCurrentTimer(userId)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	timerActionDto := ConvertTimerToDto(timer)

	return ctx.JSON(http.StatusOK, timerActionDto)
}
