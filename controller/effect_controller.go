package controller

import (
	"FGG-Service/api"
	"FGG-Service/auth_service"
	"FGG-Service/effect_service"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetAvailableEffects (POST /effects/available)
func (Server) GetAvailableEffects(ctx echo.Context) error {
	cookie, err := ctx.Cookie("sessionId")

	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, api.NotAuthorizedError{Code: api.NOACTIVESESSION})
	}

	sessionId := cookie.Value

	userId, err := auth_service.GetUserId(sessionId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	effects, err := effect_service.GetAvailableEffects(userId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	effectsDto := ConvertEffectsToDto(effects)

	return ctx.JSON(http.StatusOK, *effectsDto)
}

func ConvertEffectsToDto(effects *effect_service.Effects) *api.Effects {
	effectsDto := make(api.Effects, len(*effects))

	for i, effect := range *effects {
		effectsDto[i] = api.Effect{
			Name:        effect.Name,
			Description: effect.Description,
		}
	}

	return &effectsDto
}

// CheckAvailableEffectRoll CheckEffectRoll (GET /effects/has-roll)
func (Server) CheckAvailableEffectRoll(ctx echo.Context) error {
	cookie, err := ctx.Cookie("sessionId")

	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, api.NotAuthorizedError{Code: api.NOACTIVESESSION})
	}

	sessionId := cookie.Value

	userId, err := auth_service.GetUserId(sessionId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	doesExist, err := effect_service.CheckIfAvailableRollExists(userId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, doesExist)
}

// GetEffectHistory (GET /effects/history)
func (Server) GetEffectHistory(ctx echo.Context) error {
	cookie, err := ctx.Cookie("sessionId")

	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, api.NotAuthorizedError{Code: api.NOACTIVESESSION})
	}

	sessionId := cookie.Value

	userId, err := auth_service.GetUserId(sessionId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	effects, err := effect_service.GetEffectHistory(userId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	effectsDto := ConvertRolledEffectsToDto(effects)

	return ctx.JSON(http.StatusOK, *effectsDto)
}

func ConvertRolledEffectsToDto(effects *effect_service.RolledEffects) *api.RolledEffects {
	effectsDto := make(api.RolledEffects, len(*effects))

	for i, effect := range *effects {
		effectsDto[i] = api.RolledEffect{
			Name:        effect.Name,
			Description: effect.Description,
			RollDate:    effect.RollDate,
		}
	}

	return &effectsDto
}

// MakeEffectRoll (POST /effects/roll)
func (Server) MakeEffectRoll(ctx echo.Context) error {
	cookie, err := ctx.Cookie("sessionId")

	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, api.NotAuthorizedError{Code: api.NOACTIVESESSION})
	}

	sessionId := cookie.Value

	userId, err := auth_service.GetUserId(sessionId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	doesExist, err := effect_service.CheckIfAvailableRollExists(userId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	if !doesExist {
		return ctx.JSON(http.StatusConflict, api.ConflictError{Code: api.NOAVAILABLEROLLS})
	}

	effects, err := effect_service.MakeEffectRoll(userId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	effectsDto := ConvertEffectsToDto(effects)

	return ctx.JSON(http.StatusOK, *effectsDto)
}
