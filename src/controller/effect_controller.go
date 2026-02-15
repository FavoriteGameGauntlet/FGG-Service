package controller

import (
	"FGG-Service/api"
	"FGG-Service/src/common"
	"FGG-Service/src/effect/effect_service"
	"net/http"

	"github.com/labstack/echo/v4"
)

// MakeAvailableWheelEffectRoll (POST /wheel-effects/available/roll
func (Server) MakeAvailableWheelEffectRoll(ctx echo.Context) error {
	userId, err := GetUserId(ctx)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	effects, err := effect_service.MakeEffectRoll(userId)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	effectsDto := convertEffectsToDto(effects)

	return ctx.JSON(http.StatusOK, effectsDto)
}

// ApplyAvailableWheelEffectRoll (POST /wheel-effects/available/roll/apply)
func (Server) ApplyAvailableWheelEffectRoll(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetWheelEffectHistory (GET /wheel-effects/history)
func (Server) GetWheelEffectHistory(ctx echo.Context) error {
	userId, err := GetUserId(ctx)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	effects, err := effect_service.GetEffectHistory(userId)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	effectsDto := convertRolledEffectsToDto(effects)

	return ctx.JSON(http.StatusOK, effectsDto)
}

func convertRolledEffectsToDto(effects common.RolledEffects) api.RolledWheelEffects {
	effectsDto := make(api.RolledWheelEffects, len(effects))

	for i, effect := range effects {
		position := i - 2

		effectsDto[i] = api.RolledWheelEffect{
			Name:        effect.Name,
			Description: effect.Description,
			Position:    &position,
			RollDate:    effect.RollDate,
		}
	}

	return effectsDto
}

// GetAvailableWheelEffects (GET /wheel-effects/available)
func (s Server) GetAvailableWheelEffects(ctx echo.Context) error {
	userId, err := GetUserId(ctx)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	effects, err := effect_service.GetAvailableEffects(userId)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	effectsDto := convertEffectsToDto(effects)

	return ctx.JSON(http.StatusOK, effectsDto)
}

func convertEffectsToDto(effects common.Effects) api.WheelEffects {
	effectsDto := make(api.WheelEffects, len(effects))

	for i, effect := range effects {
		effectsDto[i] = api.WheelEffect{
			Name:        effect.Name,
			Description: effect.Description,
		}
	}

	return effectsDto
}
