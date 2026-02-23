package games

import (
	"FGG-Service/api/generated/auth"
	"FGG-Service/src/games/service"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	Service srvgames.Service
}

func NewController() Controller {
	return Controller{}
}

// GetOwnCurrentGame (GET /games/self/current)
func (c *Controller) GetOwnCurrentGame(ctx echo.Context) error {

	return nil
}

// CancelOwnCurrentGame (POST /games/self/current/cancel)
func (c *Controller) CancelOwnCurrentGame(ctx echo.Context) error {

	return nil
}

// FinishOwnCurrentGame (POST /games/self/current/finish)
func (c *Controller) FinishOwnCurrentGame(ctx echo.Context) error {

	return nil
}

// RollNewCurrentGame (POST /games/self/current/roll)
func (c *Controller) RollNewCurrentGame(ctx echo.Context) error {

	return nil
}

// GetOwnGameHistory (GET /games/self/history)
func (c *Controller) GetOwnGameHistory(ctx echo.Context) error {

	return nil
}

// GetOwnWishlistGames (GET /games/self/wishlist)
func (c *Controller) GetOwnWishlistGames(ctx echo.Context) error {

	return nil
}

// AddOwnWishlistGame (POST /games/self/wishlist)
func (c *Controller) AddOwnWishlistGame(ctx echo.Context) error {

	return nil
}

// GetCurrentGameByLogin (GET /games/{login}/current)
func (c *Controller) GetCurrentGameByLogin(ctx echo.Context, login auth.Login) error {

	return nil
}

// GetGameHistoryByLogin (GET /games/{login}/history)
func (c *Controller) GetGameHistoryByLogin(ctx echo.Context, login auth.Login) error {

	return nil
}

// GetWishlistGamesByLogin (GET /games/{login}/wishlist)
func (c *Controller) GetWishlistGamesByLogin(ctx echo.Context, login auth.Login) error {

	return nil
}

//// GetCurrentGame (GET /games/current)
//func (controller.Server) GetCurrentGame(ctx echo.Context) error {
//	userId, err := controller.GetUserId(ctx)
//
//	if err != nil {
//		return common.SendJSONErrorResponse(ctx, err)
//	}
//
//	game, err := game_service.GetCurrentGame(userId)
//
//	if err != nil {
//		return common.SendJSONErrorResponse(ctx, err)
//	}
//
//	gameDto := convertGameToDto(game)
//
//	return ctx.JSON(http.StatusOK, gameDto)
//}
//
//func convertGameToDto(game common.Game) api.Game {
//	return api.Game{
//		Name:       game.Name,
//		State:      api.GameState(game.State),
//		TimeSpent:  common.DurationToISO8601(game.TimeSpent),
//		FinishDate: game.FinishDate,
//	}
//}
//
//// CancelCurrentGame (POST /games/current/cancel)
//func (controller.Server) CancelCurrentGame(ctx echo.Context) error {
//	userId, err := controller.GetUserId(ctx)
//
//	if err != nil {
//		return common.SendJSONErrorResponse(ctx, err)
//	}
//
//	err = game_service.CancelCurrentGame(userId)
//
//	if err != nil {
//		return common.SendJSONErrorResponse(ctx, err)
//	}
//
//	return ctx.NoContent(http.StatusNoContent)
//}
//
//// FinishCurrentGame (GET /games/current/finish)
//func (controller.Server) FinishCurrentGame(ctx echo.Context) error {
//	userId, err := controller.GetUserId(ctx)
//
//	if err != nil {
//		return common.SendJSONErrorResponse(ctx, err)
//	}
//
//	err = game_service.FinishCurrentGame(userId)
//
//	if err != nil {
//		return common.SendJSONErrorResponse(ctx, err)
//	}
//
//	return ctx.NoContent(http.StatusNoContent)
//}
//
//// GetGameHistory (GET /games/history)
//func (controller.Server) GetGameHistory(ctx echo.Context) error {
//	userId, err := controller.GetUserId(ctx)
//
//	if err != nil {
//		return common.SendJSONErrorResponse(ctx, err)
//	}
//
//	games, err := game_service.GetGameHistory(userId)
//
//	if err != nil {
//		return common.SendJSONErrorResponse(ctx, err)
//	}
//
//	gamesDto := convertGamesToDto(games)
//
//	return ctx.JSON(http.StatusOK, gamesDto)
//}
//
//func convertGamesToDto(games common.Games) api.Games {
//	gamesDto := make(api.Games, len(games))
//
//	for i, game := range games {
//		gamesDto[i] = convertGameToDto(game)
//	}
//
//	return gamesDto
//}
//
//// MakeGameRoll (GET /games/roll)
//func (controller.Server) MakeGameRoll(ctx echo.Context) error {
//	userId, err := controller.GetUserId(ctx)
//
//	if err != nil {
//		return common.SendJSONErrorResponse(ctx, err)
//	}
//
//	game, err := game_service.MakeGameRoll(userId)
//
//	if err != nil {
//		return common.SendJSONErrorResponse(ctx, err)
//	}
//
//	gameDto := convertGameToDto(game)
//
//	return ctx.JSON(http.StatusOK, gameDto)
//}
//
//// GetWishlistGames (GET /games/wishlist)
//func (controller.Server) GetWishlistGames(ctx echo.Context) error {
//	userId, err := controller.GetUserId(ctx)
//
//	if err != nil {
//		return common.SendJSONErrorResponse(ctx, err)
//	}
//
//	games, err := game_service.GetUnplayedGames(userId)
//
//	if err != nil {
//		return common.SendJSONErrorResponse(ctx, err)
//	}
//
//	gamesDto := convertWishlistGamesToDto(games)
//
//	return ctx.JSON(http.StatusOK, gamesDto)
//}
//
//func convertWishlistGamesToDto(games common.UnplayedGames) api.WishlistGames {
//	gamesDto := make(api.WishlistGames, len(games))
//
//	for i, game := range games {
//		gamesDto[i] = api.WishlistGame{
//			Name: game.Name,
//		}
//	}
//
//	return gamesDto
//}
//
//// AddWishlistGame (POST /games/wishlist)
//func (controller.Server) AddWishlistGame(ctx echo.Context) error {
//	userId, err := controller.GetUserId(ctx)
//
//	if err != nil {
//		return common.SendJSONErrorResponse(ctx, err)
//	}
//
//	var gameDto api.WishlistGame
//	err = ctx.Bind(&gameDto)
//
//	if err != nil {
//		err = common.NewBadRequestError(err.Error())
//		return common.SendJSONErrorResponse(ctx, err)
//	}
//
//	game := convertWishlistGameFromDto(gameDto)
//
//	err = validator.ValidateName(game.Name)
//
//	if err != nil {
//		return common.SendJSONErrorResponse(ctx, err)
//	}
//
//	err = game_service.CreateUnplayedGame(userId, game)
//
//	if err != nil {
//		return common.SendJSONErrorResponse(ctx, err)
//	}
//
//	return ctx.NoContent(http.StatusNoContent)
//}
//
//func convertWishlistGameFromDto(gameDto api.WishlistGame) common.UnplayedGame {
//	return common.UnplayedGame{
//		Name: gameDto.Name,
//	}
//}
