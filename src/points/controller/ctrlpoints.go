package ctrlpoints

import (
	gengames "FGG-Service/api/generated/games"
	"FGG-Service/src/points/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	Service srvpoints.Service
}

func NewController() *Controller {
	return new(Controller)
}

// GetExperiencePoints (GET /points/experience-points)
func (c *Controller) GetExperiencePoints(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// ChangeExperiencePoints (POST /points/experience-points)
func (c *Controller) ChangeExperiencePoints(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// ChangeFreePoints (POST /points/free-points)
func (c *Controller) ChangeFreePoints(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetFreePoints (GET /points/free-points)
func (c *Controller) GetFreePoints(ctx echo.Context) error {
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

// GetUserTerritoryPoints (GET /points/{login}/territory-hours)
func (c *Controller) GetUserTerritoryPoints(ctx echo.Context, _ gengames.Login) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetUserTerritoryPointHistory (GET /points/{login}/territory-points/history)
func (c *Controller) GetUserTerritoryPointHistory(ctx echo.Context, _ gengames.Login) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// ChangeTerritoryHours (POST /points/territory-hours)
func (c *Controller) ChangeTerritoryHours(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetTerritoryHours (GET /points/territory-points)
func (c *Controller) GetTerritoryHours(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetAllPointInfo (GET /points/all/info)
func (c *Controller) GetAllPointInfo(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}
