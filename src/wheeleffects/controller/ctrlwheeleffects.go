package ctrlwheeleffects

import (
	"FGG-Service/api/generated/games"
	"FGG-Service/api/generated/wheel_effects"
	"FGG-Service/src/auth/service"
	"FGG-Service/src/common"
	"FGG-Service/src/points/type"
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
	s := srvwheeleffects.NewService()
	as := srvauth.NewService()

	return &Controller{
		*s,
		*as,
	}
}

// RollAvailableWheelEffects (POST /wheel-effects/available/roll)
func (c *Controller) RollAvailableWheelEffects(ctx echo.Context) error {
	var rollDto genwheeleffects.WheelEffectRoll
	err := ctx.Bind(&rollDto)

	if err != nil {
		err = common.NewBadRequestError(err.Error())
		return common.SendJSONErrorResponse(ctx, err)
	}

	userId, err := c.AuthService.GetUserId(ctx)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	isReroll := rollDto.IsReroll != nil && *rollDto.IsReroll

	effects, err := c.Service.MakeEffectRoll(userId, isReroll)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	effectsDto := convertWheelEffectsToDto(effects)

	return ctx.JSON(http.StatusOK, effectsDto)
}

func convertWheelEffectsToDto(effects typewheeleffects.WheelEffects) genwheeleffects.WheelEffects {
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
	var rollApplyDto genwheeleffects.WheelEffectRollApply
	err := ctx.Bind(&rollApplyDto)

	if err != nil {
		err = common.NewBadRequestError(err.Error())
		return common.SendJSONErrorResponse(ctx, err)
	}

	sourceUserId, err := c.AuthService.GetUserId(ctx)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	rollApply, err := c.convertDtoToWheelEffectRollApply(sourceUserId, rollApplyDto)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	changeResults, err := c.Service.ApplyWheelEffectRoll(sourceUserId, rollApply)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	changeResultsDto := convertPointChangeResultsToDto(changeResults)

	return ctx.JSON(http.StatusOK, changeResultsDto)
}

func (c *Controller) convertDtoToWheelEffectRollApply(sourceUserId int, rollApplyDto genwheeleffects.WheelEffectRollApply) (
	rollApply typewheeleffects.WheelEffectRollApply, err error) {

	pointChanges := make(typepoints.FreePointChangeByUserIds, len(rollApplyDto.PointChanges))

	for i, pointChangeByLogin := range rollApplyDto.PointChanges {
		var userId int
		userId, err = c.AuthService.GetUserIdByLogin(pointChangeByLogin.Login)

		if err != nil {
			return
		}

		pointChanges[i] = typepoints.FreePointChangeByUserId{
			Login:  pointChangeByLogin.Login,
			UserId: userId,
			PointChange: typepoints.FreePointChange{
				SourceUserId:       sourceUserId,
				ChangeSource:       pointChangeByLogin.PointChange.ChangeSource,
				DesiredChangeValue: pointChangeByLogin.PointChange.DesiredChangeValue,
			},
		}
	}

	rollApply = typewheeleffects.WheelEffectRollApply{
		PointChangeByUserIds: pointChanges,
		WheelEffectName:      rollApplyDto.WheelEffectName,
	}

	return
}

func convertPointChangeResultsToDto(changeResults typepoints.PointChangeResultByUserIds) genwheeleffects.FreePointChangeResultByLogins {
	changeResultsDto := make(genwheeleffects.FreePointChangeResultByLogins, len(changeResults))

	for i, changeResult := range changeResults {
		changeResultsDto[i].Login = changeResult.Login
		changeResultsDto[i].ChangeResult = genwheeleffects.PointChangeResult{
			ActualChangeValue: changeResult.ChangeResult.ActualChangeValue,
			FinalValue:        changeResult.ChangeResult.FinalValue,
		}
	}

	return changeResultsDto
}

// GetLastRolledWheelEffects (POST /wheel-effects/available/roll/last)
func (c *Controller) GetLastRolledWheelEffects(ctx echo.Context) error {
	userId, err := c.AuthService.GetUserId(ctx)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	effects, err := c.Service.GetLastRolledWheelEffects(userId)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	effectsDto := convertRolledWheelEffectsToDto(effects)

	return ctx.JSON(http.StatusOK, effectsDto)
}

func convertRolledWheelEffectsToDto(effects typewheeleffects.RolledWheelEffects) genwheeleffects.RolledWheelEffects {
	effectsDto := make(genwheeleffects.RolledWheelEffects, len(effects))

	for i, effect := range effects {
		effectsDto[i] = genwheeleffects.RolledWheelEffect{
			Name:        effect.Name,
			Description: effect.Description,
			RollDate:    effect.RollDate,
			Position:    effect.Position,
			IsApplied:   effect.IsApplied,
		}
	}

	return effectsDto
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

	effectsDto := convertWheelEffectsToDto(effects)

	return ctx.JSON(http.StatusOK, effectsDto)
}

// GetUserWheelEffectHistory (GET /wheel-effects/{login}/history)
func (c *Controller) GetUserWheelEffectHistory(ctx echo.Context, login gengames.Login) error {
	doesExist, err := c.AuthService.DoesUserSessionExist(ctx)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	if !doesExist {
		err = common.NewActiveSessionNotFoundUnauthorizedError()
		return common.SendJSONErrorResponse(ctx, err)
	}

	userId, err := c.AuthService.GetUserIdByLogin(login)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	effects, err := c.Service.GetEffectHistory(userId)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	effectsDto := convertRolledWheelEffectsToHistoryDto(effects)

	return ctx.JSON(http.StatusOK, effectsDto)
}

func convertRolledWheelEffectsToHistoryDto(effects typewheeleffects.RolledWheelEffectHistories) genwheeleffects.RolledWheelEffectHistories {
	effectsDto := make(genwheeleffects.RolledWheelEffectHistories, len(effects))

	for i, effect := range effects {
		effectsDto[i] = genwheeleffects.RolledWheelEffectHistory{
			Name:        effect.Name,
			Description: effect.Description,
			RollDate:    effect.RollDate,
		}
	}

	return effectsDto
}
