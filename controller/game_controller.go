package controller

import (
	"FGG-Service/api"
	"FGG-Service/auth_service"
	"FGG-Service/game_service"
	"context"
)

// GetCurrentGame (GET /users/{userId}/games/current)
func (Server) GetCurrentGame(ctx context.Context, _ api.GetCurrentGameRequestObject) (api.GetCurrentGameResponseObject, error) {
	sessionId, ok := ctx.Value("session_id").(string)

	if !ok {
		return api.GetCurrentGame401JSONResponse{Code: api.NOACTIVESESSION}, nil
	}

	userId, err := auth_service.GetUserId(sessionId)

	if err != nil {
		return api.GetCurrentGame500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	game, err := game_service.GetCurrentGame(userId)

	if err != nil {
		return api.GetCurrentGame500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	if game == nil {
		return api.GetCurrentGame404JSONResponse{Code: api.GAMENOTFOUND}, nil
	}

	gameDto := ConvertGameToDto(game)

	return api.GetCurrentGame200JSONResponse(*gameDto), nil
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

// CancelCurrentGame (POST /users/{userId}/games/current/cancel)
func (Server) CancelCurrentGame(ctx context.Context, _ api.CancelCurrentGameRequestObject) (api.CancelCurrentGameResponseObject, error) {
	sessionId, ok := ctx.Value("session_id").(string)

	if !ok {
		return api.CancelCurrentGame401JSONResponse{Code: api.NOACTIVESESSION}, nil
	}

	userId, err := auth_service.GetUserId(sessionId)

	if err != nil {
		return api.CancelCurrentGame500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	doesExist, err := game_service.CheckIfCurrentGameExists(userId)

	if err != nil {
		return api.CancelCurrentGame500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	if !doesExist {
		return api.CancelCurrentGame404JSONResponse{Code: api.GAMENOTFOUND}, nil
	}

	err = game_service.CancelCurrentGame(userId)

	if err != nil {
		return api.CancelCurrentGame500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	return api.CancelCurrentGame200Response{}, nil
}

// FinishCurrentGame (GET /users/{userId}/games/current/finish)
func (Server) FinishCurrentGame(ctx context.Context, _ api.FinishCurrentGameRequestObject) (api.FinishCurrentGameResponseObject, error) {
	sessionId, ok := ctx.Value("session_id").(string)

	if !ok {
		return api.FinishCurrentGame401JSONResponse{Code: api.NOACTIVESESSION}, nil
	}

	userId, err := auth_service.GetUserId(sessionId)

	if err != nil {
		return api.FinishCurrentGame500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	doesExist, err := game_service.CheckIfCurrentGameExists(userId)

	if err != nil {
		return api.FinishCurrentGame500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	if !doesExist {
		return api.FinishCurrentGame404JSONResponse{Code: api.GAMENOTFOUND}, nil
	}

	isSuccess, err := game_service.FinishCurrentGame(userId)

	if err != nil {
		return api.FinishCurrentGame500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	if !isSuccess {
		return api.FinishCurrentGame409JSONResponse{Code: api.NOCOMPLETEDTIMERS}, nil
	}

	return api.FinishCurrentGame200Response{}, nil
}

// GetGameHistory (GET /users/{userId}/games/history)
func (Server) GetGameHistory(ctx context.Context, _ api.GetGameHistoryRequestObject) (api.GetGameHistoryResponseObject, error) {
	sessionId, ok := ctx.Value("session_id").(string)

	if !ok {
		return api.GetGameHistory401JSONResponse{Code: api.NOACTIVESESSION}, nil
	}

	userId, err := auth_service.GetUserId(sessionId)

	if err != nil {
		return api.GetGameHistory500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	games, err := game_service.GetGameHistory(userId)

	if err != nil {
		return api.GetGameHistory500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	gamesDto := ConvertGamesToDto(games)

	return api.GetGameHistory200JSONResponse(*gamesDto), nil
}

func ConvertGamesToDto(games *game_service.Games) *api.Games {
	gamesDto := make(api.Games, len(*games))

	for i, game := range *games {
		gamesDto[i] = *ConvertGameToDto(&game)
	}

	return &gamesDto
}

// MakeGameRoll (GET /users/{userId}/games/roll)
func (Server) MakeGameRoll(ctx context.Context, _ api.MakeGameRollRequestObject) (api.MakeGameRollResponseObject, error) {
	sessionId, ok := ctx.Value("session_id").(string)

	if !ok {
		return api.MakeGameRoll401JSONResponse{Code: api.NOACTIVESESSION}, nil
	}

	userId, err := auth_service.GetUserId(sessionId)

	if err != nil {
		return api.MakeGameRoll500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	doesExist, err := game_service.CheckIfCurrentGameExists(userId)

	if err != nil {
		return api.MakeGameRoll500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	if doesExist {
		return api.MakeGameRoll409JSONResponse{Code: api.CURRENTGAMEALREADYEXISTS}, nil
	}

	game, err := game_service.MakeGameRoll(userId)

	if err != nil {
		return api.MakeGameRoll500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	if game == nil {
		return api.MakeGameRoll409JSONResponse{Code: api.NOUNPLAYEDGAMES}, nil
	}

	gameDto := ConvertGameToDto(game)

	return api.MakeGameRoll200JSONResponse(*gameDto), nil
}

// GetUnplayedGames (GET /users/{userId}/games/unplayed)
func (Server) GetUnplayedGames(ctx context.Context, _ api.GetUnplayedGamesRequestObject) (api.GetUnplayedGamesResponseObject, error) {
	sessionId, ok := ctx.Value("session_id").(string)

	if !ok {
		return api.GetUnplayedGames401JSONResponse{Code: api.NOACTIVESESSION}, nil
	}

	userId, err := auth_service.GetUserId(sessionId)

	if err != nil {
		return api.GetUnplayedGames500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	games, err := game_service.GetUnplayedGames(userId)

	if err != nil {
		return api.GetUnplayedGames500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	gamesDto := ConvertUnplayedGamesToDto(games)

	return api.GetUnplayedGames200JSONResponse(*gamesDto), nil
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

// AddUnplayedGames (POST /users/{userId}/games/unplayed)
func (Server) AddUnplayedGames(ctx context.Context, request api.AddUnplayedGamesRequestObject) (api.AddUnplayedGamesResponseObject, error) {
	sessionId, ok := ctx.Value("session_id").(string)

	if !ok {
		return api.AddUnplayedGames401JSONResponse{Code: api.NOACTIVESESSION}, nil
	}

	userId, err := auth_service.GetUserId(sessionId)

	if err != nil {
		return api.AddUnplayedGames500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	games := ConvertUnplayedGamesFrom(request.Body)

	err = game_service.AddUnplayedGames(userId, games)

	if err != nil {
		return api.AddUnplayedGames500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	return api.AddUnplayedGames200Response{}, nil
}

func ConvertUnplayedGamesFrom(games *api.UnplayedGames) *game_service.UnplayedGames {
	gamesDto := make(game_service.UnplayedGames, len(*games))

	for i, g := range *games {
		gamesDto[i] = game_service.UnplayedGame{
			Link: g.Link,
			Name: g.Name,
		}
	}

	return &gamesDto
}
