package controller

import (
	"FGG-Service/api"
	"FGG-Service/common"
	"FGG-Service/game_service"
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

	gameDto := ConvertGameToDto(game)

	return ctx.JSON(http.StatusOK, *gameDto)
}

func ConvertGameToDto(game *common.Game) *api.Game {
	return &api.Game{
		Link:       game.Link,
		Name:       game.Name,
		State:      api.GameState(game.State),
		TimeSpent:  game.TimeSpent.String(),
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

	return ctx.NoContent(http.StatusOK)
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

	return ctx.NoContent(http.StatusOK)
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

	gamesDto := ConvertGamesToDto(games)

	return ctx.JSON(http.StatusOK, *gamesDto)
}

func ConvertGamesToDto(games *common.Games) *api.Games {
	gamesDto := make(api.Games, len(*games))

	for i, game := range *games {
		gamesDto[i] = *ConvertGameToDto(&game)
	}

	return &gamesDto
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

	gameDto := ConvertGameToDto(game)

	return ctx.JSON(http.StatusOK, *gameDto)
}

// GetUnplayedGames (GET /games/unplayed)
func (Server) GetUnplayedGames(ctx echo.Context) error {
	userId, err := GetUserId(ctx)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	games, err := game_service.GetUnplayedGamesCommand(userId)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	gamesDto := ConvertUnplayedGamesToDto(games)

	return ctx.JSON(http.StatusOK, *gamesDto)
}

func ConvertUnplayedGamesToDto(games common.UnplayedGames) *api.UnplayedGames {
	gamesDto := make(api.UnplayedGames, len(games))

	for i, game := range games {
		gamesDto[i] = api.UnplayedGame{
			Link: game.Link,
			Name: game.Name,
		}
	}

	return &gamesDto
}

// AddUnplayedGames (POST /games/unplayed)
func (Server) AddUnplayedGames(ctx echo.Context) error {
	userId, err := GetUserId(ctx)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	var gamesDto api.UnplayedGames
	err = ctx.Bind(&gamesDto)

	if err != nil {
		// TODO: Fix a status code for this error, should be 400
		return SendJSONErrorResponse(ctx, err)
	}

	games := ConvertUnplayedGamesFromDto(&gamesDto)
	err = game_service.AddUnplayedGames(userId, games)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	return ctx.NoContent(http.StatusOK)
}

func ConvertUnplayedGamesFromDto(games *api.UnplayedGames) *common.UnplayedGames {
	gamesDto := make(common.UnplayedGames, len(*games))

	for i, g := range *games {
		gamesDto[i] = common.UnplayedGame{
			Link: g.Link,
			Name: g.Name,
		}
	}

	return &gamesDto
}
