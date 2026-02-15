package controller

import (
	"FGG-Service/api"
	"FGG-Service/src/common"
	"FGG-Service/src/game/game_service"
	"FGG-Service/src/validator"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetCurrentGame (GET /games/current)
func (Server) GetCurrentGame(ctx echo.Context) error {
	userId, err := GetUserId(ctx)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	game, err := game_service.GetCurrentGame(userId)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	gameDto := convertGameToDto(game)

	return ctx.JSON(http.StatusOK, gameDto)
}

func convertGameToDto(game common.Game) api.Game {
	return api.Game{
		Name:       game.Name,
		State:      api.GameState(game.State),
		TimeSpent:  common.DurationToISO8601(game.TimeSpent),
		FinishDate: game.FinishDate,
	}
}

// CancelCurrentGame (POST /games/current/cancel)
func (Server) CancelCurrentGame(ctx echo.Context) error {
	userId, err := GetUserId(ctx)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	err = game_service.CancelCurrentGame(userId)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	return ctx.NoContent(http.StatusNoContent)
}

// FinishCurrentGame (GET /games/current/finish)
func (Server) FinishCurrentGame(ctx echo.Context) error {
	userId, err := GetUserId(ctx)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	err = game_service.FinishCurrentGame(userId)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	return ctx.NoContent(http.StatusNoContent)
}

// GetGameHistory (GET /games/history)
func (Server) GetGameHistory(ctx echo.Context) error {
	userId, err := GetUserId(ctx)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	games, err := game_service.GetGameHistory(userId)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	gamesDto := convertGamesToDto(games)

	return ctx.JSON(http.StatusOK, gamesDto)
}

func convertGamesToDto(games common.Games) api.Games {
	gamesDto := make(api.Games, len(games))

	for i, game := range games {
		gamesDto[i] = convertGameToDto(game)
	}

	return gamesDto
}

// MakeGameRoll (GET /games/roll)
func (Server) MakeGameRoll(ctx echo.Context) error {
	userId, err := GetUserId(ctx)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	game, err := game_service.MakeGameRoll(userId)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	gameDto := convertGameToDto(game)

	return ctx.JSON(http.StatusOK, gameDto)
}

// GetWishlistGames (GET /games/wishlist)
func (Server) GetWishlistGames(ctx echo.Context) error {
	userId, err := GetUserId(ctx)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	games, err := game_service.GetUnplayedGames(userId)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	gamesDto := convertWishlistGamesToDto(games)

	return ctx.JSON(http.StatusOK, gamesDto)
}

func convertWishlistGamesToDto(games common.UnplayedGames) api.WishlistGames {
	gamesDto := make(api.WishlistGames, len(games))

	for i, game := range games {
		gamesDto[i] = api.WishlistGame{
			Name: game.Name,
		}
	}

	return gamesDto
}

// AddWishlistGame (POST /games/wishlist)
func (Server) AddWishlistGame(ctx echo.Context) error {
	userId, err := GetUserId(ctx)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	var gameDto api.WishlistGame
	err = ctx.Bind(&gameDto)

	if err != nil {
		err = common.NewBadRequestError(err.Error())
		return SendJSONErrorResponse(ctx, err)
	}

	game := convertWishlistGameFromDto(gameDto)

	err = validator.ValidateName(game.Name)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	err = game_service.CreateUnplayedGame(userId, game)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	return ctx.NoContent(http.StatusNoContent)
}

func convertWishlistGameFromDto(gameDto api.WishlistGame) common.UnplayedGame {
	return common.UnplayedGame{
		Name: gameDto.Name,
	}
}
