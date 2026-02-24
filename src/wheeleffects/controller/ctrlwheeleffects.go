package ctrlwheeleffects

import (
	gengames "FGG-Service/api/generated/games"
	genwheeleffects "FGG-Service/api/generated/wheel_effects"
	srvauth "FGG-Service/src/auth/service"
	"FGG-Service/src/common"
	"FGG-Service/src/wheeleffects/service"
	"FGG-Service/src/wheeleffects/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	Service     srvwheeleffects.Service
	AuthService srvauth.Service
}

func NewController() *Controller {
	return new(Controller)
}

// RollAvailableWheelEffects (POST /wheel-effects/available/roll)
func (c *Controller) RollAvailableWheelEffects(ctx echo.Context) error {
	userId, err := c.AuthService.GetUserId(ctx)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	effects, err := c.Service.MakeEffectRoll(userId)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	effectsDto := convertEffectsToDto(effects)

	return ctx.JSON(http.StatusOK, effectsDto)
}

func convertEffectsToDto(effects typewheeleffects.WheelEffects) genwheeleffects.WheelEffects {
	effectsDto := make(genwheeleffects.WheelEffects, len(effects))

	for i, effect := range effects {
		effectsDto[i] = genwheeleffects.WheelEffect{
			Name:        effect.Name,
			Description: effect.Description,
		}
	}

	return effectsDto
}

// ApplyAvailableWheelEffectRoll (POST /wheel-effects/available/roll/apply)
func (c *Controller) ApplyAvailableWheelEffectRoll(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetAvailableWheelEffectRollsCount (GET /wheel-effects/available/roll/count)
func (c *Controller) GetAvailableWheelEffectRollsCount(ctx echo.Context) error {
	userId, err := c.AuthService.GetUserId(ctx)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	count, err := c.Service.GetAvailableRollsCount(userId)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, count)
}

// GetAvailableWheelEffects (GET /wheel-effects/available)
func (c *Controller) GetAvailableWheelEffects(ctx echo.Context) error {
	userId, err := c.AuthService.GetUserId(ctx)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	effects, err := c.Service.GetAvailableEffects(userId)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	effectsDto := convertEffectsToDto(effects)

	return ctx.JSON(http.StatusOK, effectsDto)
}

// GetUserWheelEffectHistory (GET /wheel-effects/{login}/history)
func (c *Controller) GetUserWheelEffectHistory(ctx echo.Context, _ gengames.Login) error {
	userId, err := c.AuthService.GetUserId(ctx)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	effects, err := c.Service.GetEffectHistory(userId)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	effectsDto := convertRolledEffectsToDto(effects)

	return ctx.JSON(http.StatusOK, effectsDto)
}

func convertRolledEffectsToDto(effects typewheeleffects.RolledWheelEffects) genwheeleffects.RolledWheelEffects {
	effectsDto := make(genwheeleffects.RolledWheelEffects, len(effects))

	for i, effect := range effects {
		position := i - 2

		effectsDto[i] = genwheeleffects.RolledWheelEffect{
			Name:        effect.Name,
			Description: effect.Description,
			Position:    &position,
			RollDate:    effect.RollDate,
		}
	}

	return effectsDto
}
