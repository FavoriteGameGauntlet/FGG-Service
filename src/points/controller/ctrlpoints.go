package ctrlpoints

import (
	"FGG-Service/api/generated/games"
	"FGG-Service/api/generated/points"
	"FGG-Service/src/auth/service"
	"FGG-Service/src/common"
	"FGG-Service/src/points/service"
	typepoints "FGG-Service/src/points/type"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	Service     srvpoints.Service
	AuthService srvauth.Service
}

func NewController() *Controller {
	s := srvpoints.NewService()
	as := new(srvauth.Service)

	return &Controller{*s, *as}
}

// GetExperiencePoints (GET /points/experience-points)
func (c *Controller) GetExperiencePoints(ctx echo.Context) error {
	userId, err := c.AuthService.GetUserId(ctx)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	points, err := c.Service.GetExperiencePoints(userId)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, points)
}

// ChangeExperiencePoints (POST /points/experience-points)
func (c *Controller) ChangeExperiencePoints(ctx echo.Context) error {
	var pointChangeDto genpoints.PointChange
	err := ctx.Bind(&pointChangeDto)

	if err != nil {
		err = common.NewBadRequestError(err.Error())
		return common.SendJSONErrorResponse(ctx, err)
	}

	userId, err := c.AuthService.GetUserId(ctx)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	pointChange := convertDtoToPointChange(pointChangeDto)

	changeResult, err := c.Service.ChangeExperiencePoints(userId, pointChange)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	changeResultDto := convertChangeResultToDto(changeResult)

	return ctx.JSON(http.StatusOK, changeResultDto)
}

func convertDtoToPointChange(pointChangeDto genpoints.PointChange) typepoints.PointChange {
	return typepoints.PointChange{
		ChangeSource:       pointChangeDto.ChangeSource,
		DesiredChangeValue: pointChangeDto.DesiredChangeValue,
	}
}

func convertChangeResultToDto(changeResult typepoints.PointChangeResult) genpoints.PointChangeResult {
	return genpoints.PointChangeResult{
		ActualChangeValue:  changeResult.ActualChangeValue,
		ChangeDate:         changeResult.ChangeDate,
		ChangeSource:       changeResult.ChangeSource,
		DesiredChangeValue: changeResult.DesiredChangeValue,
		FinalValue:         changeResult.FinalValue,
	}
}

// ChangeFreePoints (POST /points/{login}/free-points)
func (c *Controller) ChangeFreePoints(ctx echo.Context, _ genpoints.Login) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetFreePoints (GET /points/{login}/free-points)
func (c *Controller) GetFreePoints(ctx echo.Context, _ genpoints.Login) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetUserFreePointHistory (GET /points/{login}/free-points/history)
func (c *Controller) GetUserFreePointHistory(ctx echo.Context, _ gengames.Login) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetUserPointInfo (GET /points/{login}/info)
func (c *Controller) GetUserPointInfo(ctx echo.Context, _ gengames.Login) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// ChangeUserTerritoryPoints (POST /points/{login}/territory-points)
func (c *Controller) ChangeUserTerritoryPoints(ctx echo.Context, _ gengames.Login) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetUserTerritoryPoints (GET /points/{login}/territory-points)
func (c *Controller) GetUserTerritoryPoints(ctx echo.Context, _ gengames.Login) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetUserTerritoryPointHistory (GET /points/{login}/territory-points/history)
func (c *Controller) GetUserTerritoryPointHistory(ctx echo.Context, _ gengames.Login) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// ChangeTerritoryHours (POST /points/territory-hours)
func (c *Controller) ChangeTerritoryHours(ctx echo.Context) error {
	var pointChangeDto genpoints.TerritoryHoursChange
	err := ctx.Bind(&pointChangeDto)

	if err != nil {
		err = common.NewBadRequestError(err.Error())
		return common.SendJSONErrorResponse(ctx, err)
	}

	userId, err := c.AuthService.GetUserId(ctx)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	pointChange := convertDtoToTerritoryHoursChange(pointChangeDto)
	changeResult, err := c.Service.ChangeTerritoryHours(userId, pointChange)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	changeResultDto := convertChangeResultToDto(changeResult)

	return ctx.JSON(http.StatusOK, changeResultDto)
}

func convertDtoToTerritoryHoursChange(dto genpoints.TerritoryHoursChange) typepoints.TerritoryHoursChange {
	return typepoints.TerritoryHoursChange{
		ChangeSource:       dto.ChangeSource,
		DesiredChangeValue: dto.DesiredChangeValue,
		IsSomeones:         dto.IsSomeones != nil && *dto.IsSomeones,
	}
}

// GetTerritoryHours (GET /points/territory-hours)
func (c *Controller) GetTerritoryHours(ctx echo.Context) error {
	userId, err := c.AuthService.GetUserId(ctx)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	points, err := c.Service.GetTerritoryHours(userId)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, points)
}

// GetAllPointInfo (GET /points/all/info)
func (c *Controller) GetAllPointInfo(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}
