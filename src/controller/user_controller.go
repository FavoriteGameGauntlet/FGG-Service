package controller

import (
	"FGG-Service/api"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetUserInfos (GET /users/infos)
func (Server) GetUserInfos(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetOwnDisplayName (GET /users/self/display-name)
func (Server) GetOwnDisplayName(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// ChangeOwnDisplayName (POST /users/self/display-name)
func (Server) ChangeOwnDisplayName(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetOwnStats (GET /users/self/stats)
func (Server) GetOwnStats(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetUserFreePointHistory (GET /users/{login}/free-point-history)
func (Server) GetUserFreePointHistory(ctx echo.Context, _ api.UserLogin) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetUserGameHistory (GET /users/{login}/game-history)
func (Server) GetUserGameHistory(ctx echo.Context, _ api.UserLogin) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetUserStats (GET /users/{login}/stats)
func (Server) GetUserStats(ctx echo.Context, _ api.UserLogin) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetUserTerritoryPointHistory (GET /users/{login}/territory-point-history)
func (Server) GetUserTerritoryPointHistory(ctx echo.Context, _ api.UserLogin) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetUserWheelEffectHistory (GET /users/{login}/wheel-effect-history)
func (Server) GetUserWheelEffectHistory(ctx echo.Context, _ api.UserLogin) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetUserWishlistGames (GET /users/{login}/wishlist-games)
func (Server) GetUserWishlistGames(ctx echo.Context, _ api.UserLogin) error {
	return ctx.NoContent(http.StatusNotImplemented)
}
