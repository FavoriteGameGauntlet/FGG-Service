package controller

import (
	"FGG-Service/api"
	"FGG-Service/common"
	"FGG-Service/effect_service"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetAvailableEffects (GET /effects/available)
func (Server) GetAvailableEffects(ctx echo.Context) error {
	userId, err := GetUserId(ctx)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	effects, err := effect_service.GetAvailableEffects(userId)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	effectsDto := ConvertEffectsToDto(effects)

	return ctx.JSON(http.StatusOK, effectsDto)
}

func ConvertEffectsToDto(effects common.Effects) api.Effects {
	effectsDto := make(api.Effects, len(effects))

	for i, effect := range effects {
		effectsDto[i] = api.Effect{
			Name:        effect.Name,
			Description: effect.Description,
		}
	}

	return effectsDto
}

// GetAvailableEffectsCount (GET /effects/available/count)
func (Server) GetAvailableEffectsCount(ctx echo.Context) error {
	userId, err := GetUserId(ctx)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	count, err := effect_service.GetAvailableRollsCount(userId)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, count)
}

// GetEffectHistory (GET /effects/history)
func (Server) GetEffectHistory(ctx echo.Context) error {
	userId, err := GetUserId(ctx)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	effects, err := effect_service.GetEffectHistory(userId)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	effectsDto := ConvertRolledEffectsToDto(effects)

	return ctx.JSON(http.StatusOK, effectsDto)
}

func ConvertRolledEffectsToDto(effects common.RolledEffects) api.RolledEffects {
	effectsDto := make(api.RolledEffects, len(effects))

	for i, effect := range effects {
		effectsDto[i] = api.RolledEffect{
			Name:        effect.Name,
			Description: effect.Description,
			RollDate:    effect.RollDate,
		}
	}

	return effectsDto
}

// MakeEffectRoll (POST /effects/roll)
func (Server) MakeEffectRoll(ctx echo.Context) error {
	userId, err := GetUserId(ctx)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	effects, err := effect_service.MakeEffectRoll(userId)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	effectsDto := ConvertEffectsToDto(effects)

	return ctx.JSON(http.StatusOK, effectsDto)
}
