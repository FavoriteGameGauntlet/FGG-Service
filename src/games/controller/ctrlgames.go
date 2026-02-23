package ctrlgames

import (
	"FGG-Service/api/generated/auth"
	"FGG-Service/api/generated/games"
	"FGG-Service/src/auth/service"
	"FGG-Service/src/common"
	"FGG-Service/src/games/service"
	"FGG-Service/src/games/types"
	"FGG-Service/src/validator"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	Service     srvgames.Service
	AuthService srvauth.Service
}

func NewController() Controller {
	return Controller{}
}

// GetOwnCurrentGame (GET /games/self/current)
func (c *Controller) GetOwnCurrentGame(ctx echo.Context) error {
	userId, err := c.AuthService.GetUserId(ctx)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	game, err := c.Service.GetCurrentGame(userId)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	gameDto := convertGameToDto(game)

	return ctx.JSON(http.StatusOK, gameDto)
}

func convertGameToDto(game typegames.CurrentGame) gengames.CurrentGame {
	return gengames.CurrentGame{
		Name:       game.Name,
		State:      gengames.CurrentGameState(game.State),
		TimeSpent:  common.DurationToISO8601(game.TimeSpent),
		FinishDate: game.FinishDate,
	}
}

// CancelOwnCurrentGame (POST /games/self/current/cancel)
func (c *Controller) CancelOwnCurrentGame(ctx echo.Context) error {
	userId, err := c.AuthService.GetUserId(ctx)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	err = c.Service.CancelCurrentGame(userId)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	return ctx.NoContent(http.StatusNoContent)
}

// FinishOwnCurrentGame (POST /games/self/current/finish)
func (c *Controller) FinishOwnCurrentGame(ctx echo.Context) error {
	userId, err := c.AuthService.GetUserId(ctx)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	err = c.Service.FinishCurrentGame(userId)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	return ctx.NoContent(http.StatusNoContent)
}

// RollNewCurrentGame (POST /games/self/current/roll)
func (c *Controller) RollNewCurrentGame(ctx echo.Context) error {
	userId, err := c.AuthService.GetUserId(ctx)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	game, err := c.Service.MakeGameRoll(userId)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	gameDto := convertGameToDto(game)

	return ctx.JSON(http.StatusOK, gameDto)
}

// GetOwnGameHistory (GET /games/self/history)
func (c *Controller) GetOwnGameHistory(ctx echo.Context) error {
	userId, err := c.AuthService.GetUserId(ctx)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	games, err := c.Service.GetGameHistory(userId)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	gamesDto := convertGamesToDto(games)

	return ctx.JSON(http.StatusOK, gamesDto)
}

func convertGamesToDto(games typegames.CurrentGames) gengames.CurrentGames {
	gamesDto := make(gengames.CurrentGames, len(games))

	for i, game := range games {
		gamesDto[i] = convertGameToDto(game)
	}

	return gamesDto
}

// GetOwnWishlistGames (GET /games/self/wishlist)
func (c *Controller) GetOwnWishlistGames(ctx echo.Context) error {
	userId, err := c.AuthService.GetUserId(ctx)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	games, err := c.Service.GetUnplayedGames(userId)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	gamesDto := convertWishlistGamesToDto(games)

	return ctx.JSON(http.StatusOK, gamesDto)
}

func convertWishlistGamesToDto(games typegames.WishlistGames) gengames.WishlistGames {
	gamesDto := make(gengames.WishlistGames, len(games))

	for i, game := range games {
		gamesDto[i] = gengames.WishlistGame{
			Name: game.Name,
		}
	}

	return gamesDto
}

// AddOwnWishlistGame (POST /games/self/wishlist)
func (c *Controller) AddOwnWishlistGame(ctx echo.Context) error {
	userId, err := c.AuthService.GetUserId(ctx)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	var gameDto gengames.WishlistGame
	err = ctx.Bind(&gameDto)

	if err != nil {
		err = common.NewBadRequestError(err.Error())
		return common.SendJSONErrorResponse(ctx, err)
	}

	game := convertWishlistGameFromDto(gameDto)

	err = validator.ValidateName(game.Name)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	err = c.Service.CreateUnplayedGame(userId, game)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	return ctx.NoContent(http.StatusNoContent)
}

func convertWishlistGameFromDto(gameDto gengames.WishlistGame) typegames.WishlistGame {
	return typegames.WishlistGame{
		Name: gameDto.Name,
	}
}

// GetCurrentGameByLogin (GET /games/{login}/current)
func (c *Controller) GetCurrentGameByLogin(ctx echo.Context, _ genauth.Login) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetGameHistoryByLogin (GET /games/{login}/history)
func (c *Controller) GetGameHistoryByLogin(ctx echo.Context, _ genauth.Login) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetWishlistGamesByLogin (GET /games/{login}/wishlist)
func (c *Controller) GetWishlistGamesByLogin(ctx echo.Context, _ genauth.Login) error {
	return ctx.NoContent(http.StatusNotImplemented)
}
