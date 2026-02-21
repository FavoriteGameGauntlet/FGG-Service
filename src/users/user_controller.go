package users

import (
	"FGG-Service/api"
	"FGG-Service/src/controller"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetUserInfos (GET /users/infos)
func (controller.Server) GetUserInfos(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetOwnDisplayName (GET /users/self/display-name)
func (controller.Server) GetOwnDisplayName(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// ChangeOwnDisplayName (POST /users/self/display-name)
func (controller.Server) ChangeOwnDisplayName(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetOwnStats (GET /users/self/stats)
func (controller.Server) GetOwnStats(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetUserFreePointHistory (GET /users/{login}/free-point-history)
func (controller.Server) GetUserFreePointHistory(ctx echo.Context, _ api.UserLogin) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetUserGameHistory (GET /users/{login}/game-history)
func (controller.Server) GetUserGameHistory(ctx echo.Context, _ api.UserLogin) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetUserStats (GET /users/{login}/stats)
func (controller.Server) GetUserStats(ctx echo.Context, _ api.UserLogin) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetUserTerritoryPointHistory (GET /users/{login}/territory-point-history)
func (controller.Server) GetUserTerritoryPointHistory(ctx echo.Context, _ api.UserLogin) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetUserWheelEffectHistory (GET /users/{login}/wheel-effect-history)
func (controller.Server) GetUserWheelEffectHistory(ctx echo.Context, _ api.UserLogin) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetUserWishlistGames (GET /users/{login}/wishlist-games)
func (controller.Server) GetUserWishlistGames(ctx echo.Context, _ api.UserLogin) error {
	return ctx.NoContent(http.StatusNotImplemented)
}
