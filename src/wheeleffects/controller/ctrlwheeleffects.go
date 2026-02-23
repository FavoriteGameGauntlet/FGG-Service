package wheeleffects

import (
	"FGG-Service/api/generated/auth"
	"FGG-Service/src/wheeleffects/service"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	Service srveffects.Service
}

func NewController() Controller {
	return Controller{}
}

// RollAvailableWheelEffects (POST /wheel-effects/available/roll)
func (c *Controller) RollAvailableWheelEffects(ctx echo.Context) error {

	return nil
}

// ApplyAvailableWheelEffectRoll (POST /wheel-effects/available/roll/apply)
func (c *Controller) ApplyAvailableWheelEffectRoll(ctx echo.Context) error {

	return nil
}

// GetAvailableWheelEffectRollsCount (GET /wheel-effects/available/roll/count)
func (c *Controller) GetAvailableWheelEffectRollsCount(ctx echo.Context) error {

	return nil
}

// GetOwnAvailableWheelEffects (GET /wheel-effects/self/available)
func (c *Controller) GetOwnAvailableWheelEffects(ctx echo.Context) error {

	return nil
}

// GetOwnWheelEffectHistory (GET /wheel-effects/self/history)
func (c *Controller) GetOwnWheelEffectHistory(ctx echo.Context) error {

	return nil
}

// GetWheelEffectHistoryByLogin (GET /wheel-effects/{login}/history)
func (c *Controller) GetWheelEffectHistoryByLogin(ctx echo.Context, login auth.Login) error {

	return nil
}

//// MakeAvailableWheelEffectRoll (POST /wheel-effects/available/roll)
//func (controller.Server) MakeAvailableWheelEffectRoll(ctx echo.Context) error {
//	userId, err := auth.GetUserId(ctx)
//
//	if err != nil {
//		return auth.SendJSONErrorResponse(ctx, err)
//	}
//
//	effects, err := wheel_effect_service.MakeEffectRoll(userId)
//
//	if err != nil {
//		return auth.SendJSONErrorResponse(ctx, err)
//	}
//
//	effectsDto := convertEffectsToDto(effects)
//
//	return ctx.JSON(http.StatusOK, effectsDto)
//}
//
//// ApplyAvailableWheelEffectRoll (POST /wheel-effects/available/roll/apply)
//func (controller.Server) ApplyAvailableWheelEffectRoll(ctx echo.Context) error {
//	return ctx.NoContent(http.StatusNotImplemented)
//}
//
//// GetWheelEffectHistory (GET /wheel-effects/history)
//func (controller.Server) GetWheelEffectHistory(ctx echo.Context) error {
//	userId, err := auth.GetUserId(ctx)
//
//	if err != nil {
//		return auth.SendJSONErrorResponse(ctx, err)
//	}
//
//	effects, err := wheel_effect_service.GetEffectHistory(userId)
//
//	if err != nil {
//		return auth.SendJSONErrorResponse(ctx, err)
//	}
//
//	effectsDto := convertRolledEffectsToDto(effects)
//
//	return ctx.JSON(http.StatusOK, effectsDto)
//}
//
//func convertRolledEffectsToDto(effects common.RolledEffects) api.RolledWheelEffects {
//	effectsDto := make(api.RolledWheelEffects, len(effects))
//
//	for i, effect := range effects {
//		position := i - 2
//
//		effectsDto[i] = api.RolledWheelEffect{
//			Name:        effect.Name,
//			Description: effect.Description,
//			Position:    &position,
//			RollDate:    effect.RollDate,
//		}
//	}
//
//	return effectsDto
//}
//
//// GetAvailableWheelEffects (GET /wheel-effects/available)
//func (s controller.Server) GetAvailableWheelEffects(ctx echo.Context) error {
//	userId, err := auth.GetUserId(ctx)
//
//	if err != nil {
//		return auth.SendJSONErrorResponse(ctx, err)
//	}
//
//	effects, err := wheel_effect_service.GetAvailableEffects(userId)
//
//	if err != nil {
//		return auth.SendJSONErrorResponse(ctx, err)
//	}
//
//	effectsDto := convertEffectsToDto(effects)
//
//	return ctx.JSON(http.StatusOK, effectsDto)
//}
//
//func convertEffectsToDto(effects common.Effects) api.WheelEffects {
//	effectsDto := make(api.WheelEffects, len(effects))
//
//	for i, effect := range effects {
//		effectsDto[i] = api.WheelEffect{
//			Name:        effect.Name,
//			Description: effect.Description,
//		}
//	}
//
//	return effectsDto
//}
