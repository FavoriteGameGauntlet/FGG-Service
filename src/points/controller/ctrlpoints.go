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
	as := srvauth.NewService()

	return &Controller{
		*s,
		*as,
	}
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
		ChangeSource:       changeResult.ChangeSource,
		DesiredChangeValue: changeResult.DesiredChangeValue,
		FinalValue:         changeResult.FinalValue,
	}
}

// ChangeFreePoints (POST /points/{login}/free-points)
func (c *Controller) ChangeFreePoints(ctx echo.Context, login genpoints.Login) error {
	var pointChangeDto genpoints.FreePointChange
	err := ctx.Bind(&pointChangeDto)

	if err != nil {
		err = common.NewBadRequestError(err.Error())
		return common.SendJSONErrorResponse(ctx, err)
	}

	sourceUserId, err := c.AuthService.GetUserId(ctx)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	userId, err := c.AuthService.GetUserIdByLogin(login)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	pointChange := convertDtoToFreePointChange(sourceUserId, pointChangeDto)

	changeResult, err := c.Service.ChangeFreePoints(userId, pointChange)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	changeResultDto := convertFreePointChangeResultToDto(changeResult)

	return ctx.JSON(http.StatusOK, changeResultDto)
}

func convertDtoToFreePointChange(sourceUserId int, pointChangeDto genpoints.FreePointChange) typepoints.FreePointChange {
	return typepoints.FreePointChange{
		SourceUserId:       sourceUserId,
		ChangeSource:       pointChangeDto.ChangeSource,
		DesiredChangeValue: pointChangeDto.DesiredChangeValue,
		WheelEffectName:    pointChangeDto.WheelEffectName,
	}
}

func convertFreePointChangeResultToDto(changeResult typepoints.PointChangeResult) genpoints.FreePointChangeResult {
	return genpoints.FreePointChangeResult{
		ActualChangeValue:  changeResult.ActualChangeValue,
		ChangeSource:       changeResult.ChangeSource,
		DesiredChangeValue: changeResult.DesiredChangeValue,
		FinalValue:         changeResult.FinalValue,
	}
}

// GetFreePoints (GET /points/{login}/free-points)
func (c *Controller) GetFreePoints(ctx echo.Context, login genpoints.Login) error {
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

	points, err := c.Service.GetFreePoints(userId)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, points)
}

// GetUserFreePointHistory (GET /points/{login}/free-points/history)
func (c *Controller) GetUserFreePointHistory(ctx echo.Context, login genpoints.Login) error {
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

	history, err := c.Service.GetUserFreePointHistory(userId)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	historyDto := convertFreePointChangeHistoriesToDto(history)

	return ctx.JSON(http.StatusOK, historyDto)
}

func convertFreePointChangeHistoriesToDto(history typepoints.FreePointChangeHistories) genpoints.FreePointChangeHistories {
	historyDto := make(genpoints.FreePointChangeHistories, len(history))

	for i, entry := range history {
		historyDto[i] = genpoints.FreePointChangeHistory{
			ActualChangeValue:  entry.ActualChangeValue,
			ChangeDate:         entry.ChangeDate,
			ChangeSource:       entry.ChangeSource,
			DesiredChangeValue: entry.DesiredChangeValue,
			FinalValue:         entry.FinalValue,
		}
	}

	return historyDto
}

// GetUserPointInfo (GET /points/{login}/info)
func (c *Controller) GetUserPointInfo(ctx echo.Context, login genpoints.Login) error {
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

	info, err := c.Service.GetUserPointInfo(userId)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, convertPointInfoToDto(info))
}

func convertPointInfoToDto(info typepoints.PointInfo) genpoints.PointInfo {
	return genpoints.PointInfo{
		TerritoryPoints:  info.TerritoryPoints,
		FreePoints:       info.FreePoints,
		AvailableRolls:   info.AvailableRolls,
		TerritoryHours:   info.TerritoryHours,
		ExperiencePoints: info.ExperiencePoints,
	}
}

// ChangeUserTerritoryPoints (POST /points/{login}/territory-points)
func (c *Controller) ChangeUserTerritoryPoints(ctx echo.Context, _ gengames.Login) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetUserTerritoryPoints (GET /points/{login}/territory-points)
func (c *Controller) GetUserTerritoryPoints(ctx echo.Context, login genpoints.Login) error {
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

	points, err := c.Service.GetTerritoryPoints(userId)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, points)
}

// GetUserTerritoryPointHistory (GET /points/{login}/territory-points/history)
func (c *Controller) GetUserTerritoryPointHistory(ctx echo.Context, _ gengames.Login) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// ChangeTerritoryHours (POST /points/territory-hours)
func (c *Controller) ChangeTerritoryHours(ctx echo.Context) error {
	var pointChangeDto genpoints.TerritoryHourChange
	err := ctx.Bind(&pointChangeDto)

	if err != nil {
		err = common.NewBadRequestError(err.Error())
		return common.SendJSONErrorResponse(ctx, err)
	}

	userId, err := c.AuthService.GetUserId(ctx)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	pointChange := convertDtoToTerritoryHourChange(pointChangeDto)
	changeResult, err := c.Service.ChangeTerritoryHours(userId, pointChange)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	changeResultDto := convertChangeResultToDto(changeResult)

	return ctx.JSON(http.StatusOK, changeResultDto)
}

func convertDtoToTerritoryHourChange(dto genpoints.TerritoryHourChange) typepoints.TerritoryHourChange {
	return typepoints.TerritoryHourChange{
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
	doesExist, err := c.AuthService.DoesUserSessionExist(ctx)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	if !doesExist {
		err = common.NewActiveSessionNotFoundUnauthorizedError()
		return common.SendJSONErrorResponse(ctx, err)
	}

	infos, err := c.Service.GetAllPointInfo()

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	infosDto := convertPointInfoByLoginsToDto(infos)

	return ctx.JSON(http.StatusOK, infosDto)
}

func convertPointInfoByLoginsToDto(infos typepoints.PointInfoByLogins) genpoints.PointInfoByLogins {
	infosDto := make(genpoints.PointInfoByLogins, len(infos))

	for i, info := range infos {
		login := info.Login
		pointInfo := convertPointInfoToDto(info.PointInfo)
		infosDto[i].Login = &login
		infosDto[i].PointInfo = &pointInfo
	}

	return infosDto
}
