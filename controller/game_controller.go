package controller

import (
	"FGG-Service/api"
	"FGG-Service/auth_service"
	"FGG-Service/game_service"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetCurrentGame (GET /games/current)
func (Server) GetCurrentGame(ctx echo.Context) error {
	cookie, err := ctx.Cookie("sessionId")

	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, api.NotAuthorizedError{Code: api.NOACTIVESESSION})
	}

	sessionId := cookie.Value

	userId, err := auth_service.GetUserId(sessionId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	game, err := game_service.GetCurrentGame(userId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	if game == nil {
		return ctx.JSON(http.StatusNotFound, api.NotFoundError{Code: api.GAMENOTFOUND})
	}

	gameDto := ConvertGameToDto(game)

	return ctx.JSON(http.StatusOK, *gameDto)
}

func ConvertGameToDto(game *game_service.Game) *api.Game {
	return &api.Game{
		Link:       game.Link,
		Name:       game.Name,
		State:      api.GameState(game.State),
		HourCount:  game.HourCount,
		FinishDate: game.FinishDate,
	}
}

// CancelCurrentGame (POST /games/current/cancel)
func (Server) CancelCurrentGame(ctx echo.Context) error {
	cookie, err := ctx.Cookie("sessionId")

	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, api.NotAuthorizedError{Code: api.NOACTIVESESSION})
	}

	sessionId := cookie.Value

	userId, err := auth_service.GetUserId(sessionId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	doesExist, err := game_service.CheckIfCurrentGameExists(userId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	if doesExist {
		return ctx.JSON(http.StatusNotFound, api.NotFoundError{Code: api.GAMENOTFOUND})
	}

	err = game_service.CancelCurrentGame(userId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	return ctx.NoContent(http.StatusOK)
}

// FinishCurrentGame (GET /games/current/finish)
func (Server) FinishCurrentGame(ctx echo.Context) error {
	cookie, err := ctx.Cookie("sessionId")

	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, api.NotAuthorizedError{Code: api.NOACTIVESESSION})
	}

	sessionId := cookie.Value

	userId, err := auth_service.GetUserId(sessionId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	doesExist, err := game_service.CheckIfCurrentGameExists(userId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	if !doesExist {
		return ctx.JSON(http.StatusNotFound, api.NotFoundError{Code: api.GAMENOTFOUND})
	}

	isSuccess, err := game_service.FinishCurrentGame(userId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	if !isSuccess {
		return ctx.JSON(http.StatusConflict, api.ConflictError{Code: api.NOCOMPLETEDTIMERS})
	}

	return ctx.NoContent(http.StatusOK)
}

// GetGameHistory (GET /games/history)
func (Server) GetGameHistory(ctx echo.Context) error {
	cookie, err := ctx.Cookie("sessionId")

	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, api.NotAuthorizedError{Code: api.NOACTIVESESSION})
	}

	sessionId := cookie.Value

	userId, err := auth_service.GetUserId(sessionId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	games, err := game_service.GetGameHistory(userId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	gamesDto := ConvertGamesToDto(games)

	return ctx.JSON(http.StatusOK, *gamesDto)
}

func ConvertGamesToDto(games *game_service.Games) *api.Games {
	gamesDto := make(api.Games, len(*games))

	for i, game := range *games {
		gamesDto[i] = *ConvertGameToDto(&game)
	}

	return &gamesDto
}

// MakeGameRoll (GET /games/roll)
func (Server) MakeGameRoll(ctx echo.Context) error {
	cookie, err := ctx.Cookie("sessionId")

	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, api.NotAuthorizedError{Code: api.NOACTIVESESSION})
	}

	sessionId := cookie.Value

	userId, err := auth_service.GetUserId(sessionId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	doesExist, err := game_service.CheckIfCurrentGameExists(userId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	if doesExist {
		return ctx.JSON(http.StatusConflict, api.ConflictError{Code: api.CURRENTGAMEALREADYEXISTS})
	}

	game, err := game_service.MakeGameRoll(userId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	if game == nil {
		return ctx.JSON(http.StatusConflict, api.ConflictError{Code: api.NOUNPLAYEDGAMES})
	}

	gameDto := ConvertGameToDto(game)

	return ctx.JSON(http.StatusOK, *gameDto)
}

// GetUnplayedGames (GET /games/unplayed)
func (Server) GetUnplayedGames(ctx echo.Context) error {
	cookie, err := ctx.Cookie("sessionId")

	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, api.NotAuthorizedError{Code: api.NOACTIVESESSION})
	}

	sessionId := cookie.Value

	userId, err := auth_service.GetUserId(sessionId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	games, err := game_service.GetUnplayedGames(userId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	gamesDto := ConvertUnplayedGamesToDto(games)

	return ctx.JSON(http.StatusOK, *gamesDto)
}

func ConvertUnplayedGamesToDto(games *game_service.UnplayedGames) *api.UnplayedGames {
	gamesDto := make(api.UnplayedGames, len(*games))

	for i, game := range *games {
		gamesDto[i] = api.UnplayedGame{
			Link: game.Link,
			Name: game.Name,
		}
	}

	return &gamesDto
}

// AddUnplayedGames (POST /games/unplayed)
func (Server) AddUnplayedGames(ctx echo.Context) error {
	cookie, err := ctx.Cookie("sessionId")

	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, api.NotAuthorizedError{Code: api.NOACTIVESESSION})
	}

	sessionId := cookie.Value

	userId, err := auth_service.GetUserId(sessionId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	var gamesDto api.UnplayedGames
	if err = ctx.Bind(&gamesDto); err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	games := ConvertUnplayedGamesFromDto(&gamesDto)

	if err = game_service.AddUnplayedGames(userId, games); err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	return ctx.NoContent(http.StatusOK)
}

func ConvertUnplayedGamesFromDto(games *api.UnplayedGames) *game_service.UnplayedGames {
	gamesDto := make(game_service.UnplayedGames, len(*games))

	for i, g := range *games {
		gamesDto[i] = game_service.UnplayedGame{
			Link: g.Link,
			Name: g.Name,
		}
	}

	return &gamesDto
}
