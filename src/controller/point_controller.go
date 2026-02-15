package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// ChangeExperiencePoints (POST /points/experience-points)
func (Server) ChangeExperiencePoints(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetFreePointHistory (GET /points/free-point-history)
func (Server) GetFreePointHistory(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// ChangeFreePoints (POST /points/free-points)
func (Server) ChangeFreePoints(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// ChangeTerritoryHours (POST /points/territory-hours)
func (Server) ChangeTerritoryHours(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetTerritoryPointHistory (GET /points/territory-point-history)
func (Server) GetTerritoryPointHistory(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// ChangeTerritoryPoints (POST /points/territory-points)
func (Server) ChangeTerritoryPoints(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}
