package ctrlpoints

import (
	"FGG-Service/api/generated/auth"
	"FGG-Service/src/points/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	Service srvpoints.Service
}

func NewController() Controller {
	return Controller{}
}

// GetOwnExperiencePoints (GET /points/self/experience-points)
func (c *Controller) GetOwnExperiencePoints(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// ChangeOwnExperiencePoints (POST /points/self/experience-points)
func (c *Controller) ChangeOwnExperiencePoints(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetOwnFreePoints (GET /points/self/free-points)
func (c *Controller) GetOwnFreePoints(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// ChangeOwnFreePoints (POST /points/self/free-points)
func (c *Controller) ChangeOwnFreePoints(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetOwnFreePointHistory (GET /points/self/free-points/history)
func (c *Controller) GetOwnFreePointHistory(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetOwnPointInfo (GET /points/self/info)
func (c *Controller) GetOwnPointInfo(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetOwnTerritoryHours (GET /points/self/territory-hours)
func (c *Controller) GetOwnTerritoryHours(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// ChangeOwnTerritoryHours (POST /points/self/territory-hours)
func (c *Controller) ChangeOwnTerritoryHours(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetOwnTerritoryPoints (GET /points/self/territory-points)
func (c *Controller) GetOwnTerritoryPoints(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// ChangeOwnTerritoryPoints (POST /points/self/territory-points)
func (c *Controller) ChangeOwnTerritoryPoints(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetOwnTerritoryPointHistory (GET /points/self/territory-points/history)
func (c *Controller) GetOwnTerritoryPointHistory(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetFreePointHistoryByLogin (GET /points/{login}/free-points/history)
func (c *Controller) GetFreePointHistoryByLogin(ctx echo.Context, _ genauth.Login) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetPointInfoByLogin (GET /points/{login}/info)
func (c *Controller) GetPointInfoByLogin(ctx echo.Context, _ genauth.Login) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetTerritoryPointHistoryByLogin (GET /points/{login}/territory-points/history)
func (c *Controller) GetTerritoryPointHistoryByLogin(ctx echo.Context, _ genauth.Login) error {
	return ctx.NoContent(http.StatusNotImplemented)
}
