package controller

import (
	"FGG-Service/api"
	"FGG-Service/game_service"
	"FGG-Service/user_service"
	"context"
)

// GetCurrentGame (GET /users/{userId}/games/current)
func (Server) GetCurrentGame(_ context.Context, request api.GetCurrentGameRequestObject) (api.GetCurrentGameResponseObject, error) {
	doesExist, err := user_service.CheckIfUserExistsById(request.UserId)

	if err != nil {
		return api.GetCurrentGame503JSONResponse{Code: api.CHECKUSER, Message: err.Error()}, nil
	}

	if !doesExist {
		return api.GetCurrentGame404JSONResponse{Code: api.USERNOTFOUND}, nil
	}

	game, err := game_service.GetCurrentGame(request.UserId)

	if err != nil {
		return api.GetCurrentGame503JSONResponse{Code: api.GETCURRENTGAME, Message: err.Error()}, nil
	}

	if game == nil {
		return api.GetCurrentGame404JSONResponse{Code: api.GAMENOTFOUND}, nil
	}

	gameDto := ConvertGameToDto(game)

	return api.GetCurrentGame200JSONResponse(*gameDto), nil
}

func ConvertGameToDto(game *game_service.Game) *api.GameDto {
	return &api.GameDto{
		Link:       game.Link,
		Name:       game.Name,
		State:      api.GameDtoState(game.State),
		FinishDate: game.FinishDate,
	}
}

// CancelCurrentGame (POST /users/{userId}/games/current/cancel)
func (Server) CancelCurrentGame(_ context.Context, request api.CancelCurrentGameRequestObject) (api.CancelCurrentGameResponseObject, error) {
	doesExist, err := user_service.CheckIfUserExistsById(request.UserId)

	if err != nil {
		return api.CancelCurrentGame503JSONResponse{Code: api.CHECKUSER, Message: err.Error()}, nil
	}

	if !doesExist {
		return api.CancelCurrentGame404JSONResponse{Code: api.USERNOTFOUND}, nil
	}

	doesExist, err = game_service.CheckIfCurrentGameExists(request.UserId)

	if err != nil {
		return api.CancelCurrentGame503JSONResponse{Code: api.CHECKCURRENTGAME, Message: err.Error()}, nil
	}

	if !doesExist {
		return api.CancelCurrentGame404JSONResponse{Code: api.GAMENOTFOUND}, nil
	}

	err = game_service.CancelCurrentGame(request.UserId)

	if err != nil {
		return api.CancelCurrentGame503JSONResponse{Code: api.CANCELCURRENTGAME, Message: err.Error()}, nil
	}

	return api.CancelCurrentGame200Response{}, nil
}

// FinishCurrentGame (GET /users/{userId}/games/current/finish)
func (Server) FinishCurrentGame(_ context.Context, request api.FinishCurrentGameRequestObject) (api.FinishCurrentGameResponseObject, error) {
	doesExist, err := user_service.CheckIfUserExistsById(request.UserId)

	if err != nil {
		return api.FinishCurrentGame503JSONResponse{Code: api.CHECKUSER, Message: err.Error()}, nil
	}

	if !doesExist {
		return api.FinishCurrentGame404JSONResponse{Code: api.USERNOTFOUND}, nil
	}

	doesExist, err = game_service.CheckIfCurrentGameExists(request.UserId)

	if err != nil {
		return api.FinishCurrentGame503JSONResponse{Code: api.CHECKCURRENTGAME, Message: err.Error()}, nil
	}

	if !doesExist {
		return api.FinishCurrentGame404JSONResponse{Code: api.GAMENOTFOUND}, nil
	}

	err = game_service.FinishCurrentGame(request.UserId)

	if err != nil {
		return api.FinishCurrentGame503JSONResponse{Code: api.FINISHCURRENTGAME, Message: err.Error()}, nil
	}

	return api.FinishCurrentGame200Response{}, nil
}

// GetGameHistory (GET /users/{userId}/games/history)
func (Server) GetGameHistory(_ context.Context, request api.GetGameHistoryRequestObject) (api.GetGameHistoryResponseObject, error) {
	doesExist, err := user_service.CheckIfUserExistsById(request.UserId)

	if err != nil {
		return api.GetGameHistory503JSONResponse{Code: api.CHECKUSER, Message: err.Error()}, nil
	}

	if !doesExist {
		return api.GetGameHistory404JSONResponse{Code: api.USERNOTFOUND}, nil
	}

	games, err := game_service.GetGameHistory(request.UserId)

	if err != nil {
		return api.GetGameHistory503JSONResponse{Code: api.GETGAMEHISTORY, Message: err.Error()}, nil
	}

	gamesDto := ConvertGamesToDto(games)

	return api.GetGameHistory200JSONResponse(*gamesDto), nil
}

func ConvertGamesToDto(games *game_service.Games) *api.GamesDto {
	gamesDto := make(api.GamesDto, len(*games))

	for i, game := range *games {
		gamesDto[i] = *ConvertGameToDto(&game)
	}

	return &gamesDto
}

// MakeGameRoll (GET /users/{userId}/games/roll)
func (Server) MakeGameRoll(_ context.Context, request api.MakeGameRollRequestObject) (api.MakeGameRollResponseObject, error) {
	doesExist, err := user_service.CheckIfUserExistsById(request.UserId)

	if err != nil {
		return api.MakeGameRoll503JSONResponse{Code: api.CHECKUSER, Message: err.Error()}, nil
	}

	if !doesExist {
		return api.MakeGameRoll404JSONResponse{Code: api.USERNOTFOUND}, nil
	}

	doesExist, err = game_service.CheckIfCurrentGameExists(request.UserId)

	if err != nil {
		return api.MakeGameRoll503JSONResponse{Code: api.CHECKCURRENTGAME, Message: err.Error()}, nil
	}

	if doesExist {
		return api.MakeGameRoll409JSONResponse{Code: api.CURRENTGAMEALREADYEXISTS}, nil
	}

	game, err := game_service.MakeGameRoll(request.UserId)

	if err != nil {
		return api.MakeGameRoll503JSONResponse{Code: api.MAKEGAMEROLL, Message: err.Error()}, nil
	}

	gameDto := ConvertGameToDto(game)

	return api.MakeGameRoll200JSONResponse(*gameDto), nil
}

// GetUnplayedGames (GET /users/{userId}/games/unplayed)
func (Server) GetUnplayedGames(_ context.Context, request api.GetUnplayedGamesRequestObject) (api.GetUnplayedGamesResponseObject, error) {
	doesExist, err := user_service.CheckIfUserExistsById(request.UserId)

	if err != nil {
		return api.GetUnplayedGames503JSONResponse{Code: api.CHECKUSER, Message: err.Error()}, nil
	}

	if !doesExist {
		return api.GetUnplayedGames404JSONResponse{Code: api.USERNOTFOUND}, nil
	}

	games, err := game_service.GetUnplayedGames(request.UserId)

	if err != nil {
		return api.GetUnplayedGames503JSONResponse{Code: api.GETUNPLAYEDGAMES, Message: err.Error()}, nil
	}

	gamesDto := ConvertUnplayedGamesToDto(games)

	return api.GetUnplayedGames200JSONResponse(*gamesDto), nil
}

func ConvertUnplayedGamesToDto(games *game_service.UnplayedGames) *api.UnplayedGamesDto {
	gamesDto := make(api.UnplayedGamesDto, len(*games))

	for i, game := range *games {
		gamesDto[i] = api.UnplayedGameDto{
			Link: game.Link,
			Name: game.Name,
		}
	}

	return &gamesDto
}

// AddUnplayedGames (POST /users/{userId}/games/unplayed)
func (Server) AddUnplayedGames(_ context.Context, request api.AddUnplayedGamesRequestObject) (api.AddUnplayedGamesResponseObject, error) {
	doesExist, err := user_service.CheckIfUserExistsById(request.UserId)

	if err != nil {
		return api.AddUnplayedGames503JSONResponse{Code: api.CHECKUSER, Message: err.Error()}, nil
	}

	if !doesExist {
		return api.AddUnplayedGames404JSONResponse{Code: api.USERNOTFOUND}, nil
	}

	games := ConvertUnplayedGamesFromDto(request.Body)

	err = game_service.AddUnplayedGames(request.UserId, games)

	if err != nil {
		return api.AddUnplayedGames503JSONResponse{Code: api.ADDUNPLAYEDGAMES, Message: err.Error()}, nil
	}

	return api.AddUnplayedGames200Response{}, nil
}

func ConvertUnplayedGamesFromDto(gamesDto *api.UnplayedGamesDto) *game_service.UnplayedGames {
	games := make(game_service.UnplayedGames, len(*gamesDto))

	for i, g := range *gamesDto {
		games[i] = game_service.UnplayedGame{
			Link: g.Link,
			Name: g.Name,
		}
	}

	return &games
}
