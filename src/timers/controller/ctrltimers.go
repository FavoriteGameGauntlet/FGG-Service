package timers

import (
	"FGG-Service/api/generated/timers"
	"FGG-Service/src/timers/service"
	"context"
)

type Controller struct {
	Service srvtimers.Service
}

// GetOwnCurrentTimer GET /timers/self/current
func (c *Controller) GetOwnCurrentTimer(ctx context.Context) (
	timers.GetOwnCurrentTimerRes, error) {

	return nil, nil
}

// PauseOwnCurrentTimer POST /timers/self/current/pause
func (c *Controller) PauseOwnCurrentTimer(ctx context.Context) (
	timers.PauseOwnCurrentTimerRes, error) {

	return nil, nil
}

// StartOwnCurrentTimer POST /timers/self/current/start
func (c *Controller) StartOwnCurrentTimer(ctx context.Context) (
	timers.StartOwnCurrentTimerRes, error) {

	return nil, nil
}

//// GetCurrentTimer (GET /timers/current)
//func (controller.Server) GetCurrentTimer(ctx echo.Context) error {
//	userId, err := auth.GetUserId(ctx)
//
//	if err != nil {
//		return auth.SendJSONErrorResponse(ctx, err)
//	}
//
//	timer, err := timer_service.GetOrCreateCurrentTimer(userId)
//
//	if err != nil {
//		return auth.SendJSONErrorResponse(ctx, err)
//	}
//
//	timerDto := ConvertTimerToDto(timer)
//
//	return ctx.JSON(http.StatusOK, timerDto)
//}
//
//func ConvertTimerToDto(timer common.Timer) api.Timer {
//	return api.Timer{
//		Duration:       common.DurationToISO8601(timer.Duration),
//		RemainingTime:  common.DurationToISO8601(timer.RemainingTime),
//		State:          api.TimerState(timer.State),
//		LastActionDate: timer.LastActionDate,
//	}
//}
//
//// PauseCurrentTimer (POST /timers/current/pause)
//func (controller.Server) PauseCurrentTimer(ctx echo.Context) error {
//	userId, err := auth.GetUserId(ctx)
//
//	if err != nil {
//		return auth.SendJSONErrorResponse(ctx, err)
//	}
//
//	timer, err := timer_service.PauseCurrentTimer(userId)
//
//	if err != nil {
//		return auth.SendJSONErrorResponse(ctx, err)
//	}
//
//	timerActionDto := ConvertTimerToDto(timer)
//
//	return ctx.JSON(http.StatusOK, timerActionDto)
//}
//
//// StartCurrentTimer (POST /timers/current/start)
//func (controller.Server) StartCurrentTimer(ctx echo.Context) error {
//	userId, err := auth.GetUserId(ctx)
//
//	if err != nil {
//		return auth.SendJSONErrorResponse(ctx, err)
//	}
//
//	timer, err := timer_service.StartCurrentTimer(userId)
//
//	if err != nil {
//		return auth.SendJSONErrorResponse(ctx, err)
//	}
//
//	timerActionDto := ConvertTimerToDto(timer)
//
//	return ctx.JSON(http.StatusOK, timerActionDto)
//}
