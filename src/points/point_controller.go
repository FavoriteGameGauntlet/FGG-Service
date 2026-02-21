package points

import (
	"FGG-Service/src/controller"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ChangeExperiencePoints (POST /points/experience-points)
func (controller.Server) ChangeExperiencePoints(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetFreePointHistory (GET /points/free-point-history)
func (controller.Server) GetFreePointHistory(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// ChangeFreePoints (POST /points/free-points)
func (controller.Server) ChangeFreePoints(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// ChangeTerritoryHours (POST /points/territory-hours)
func (controller.Server) ChangeTerritoryHours(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetTerritoryPointHistory (GET /points/territory-point-history)
func (controller.Server) GetTerritoryPointHistory(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// ChangeTerritoryPoints (POST /points/territory-points)
func (controller.Server) ChangeTerritoryPoints(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}
