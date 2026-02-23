package points

import (
	"FGG-Service/api/generated/auth"
	"FGG-Service/src/points/service"

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

	return nil
}

// ChangeOwnExperiencePoints (POST /points/self/experience-points)
func (c *Controller) ChangeOwnExperiencePoints(ctx echo.Context) error {

	return nil
}

// GetOwnFreePoints (GET /points/self/free-points)
func (c *Controller) GetOwnFreePoints(ctx echo.Context) error {

	return nil
}

// ChangeOwnFreePoints (POST /points/self/free-points)
func (c *Controller) ChangeOwnFreePoints(ctx echo.Context) error {

	return nil
}

// GetOwnFreePointHistory (GET /points/self/free-points/history)
func (c *Controller) GetOwnFreePointHistory(ctx echo.Context) error {

	return nil
}

// GetOwnPointInfo (GET /points/self/info)
func (c *Controller) GetOwnPointInfo(ctx echo.Context) error {

	return nil
}

// GetOwnTerritoryHours (GET /points/self/territory-hours)
func (c *Controller) GetOwnTerritoryHours(ctx echo.Context) error {

	return nil
}

// ChangeOwnTerritoryHours (POST /points/self/territory-hours)
func (c *Controller) ChangeOwnTerritoryHours(ctx echo.Context) error {

	return nil
}

// GetOwnTerritoryPoints (GET /points/self/territory-points)
func (c *Controller) GetOwnTerritoryPoints(ctx echo.Context) error {

	return nil
}

// ChangeOwnTerritoryPoints (POST /points/self/territory-points)
func (c *Controller) ChangeOwnTerritoryPoints(ctx echo.Context) error {

	return nil
}

// GetOwnTerritoryPointHistory (GET /points/self/territory-points/history)
func (c *Controller) GetOwnTerritoryPointHistory(ctx echo.Context) error {

	return nil
}

// GetFreePointHistoryByLogin (GET /points/{login}/free-points/history)
func (c *Controller) GetFreePointHistoryByLogin(ctx echo.Context, login auth.Login) error {

	return nil
}

// GetPointInfoByLogin (GET /points/{login}/info)
func (c *Controller) GetPointInfoByLogin(ctx echo.Context, login auth.Login) error {

	return nil
}

// GetTerritoryPointHistoryByLogin (GET /points/{login}/territory-points/history)
func (c *Controller) GetTerritoryPointHistoryByLogin(ctx echo.Context, login auth.Login) error {

	return nil
}
