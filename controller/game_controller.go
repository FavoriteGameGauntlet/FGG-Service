package controller

import (
	"FGG-Service/api"
	"FGG-Service/game_service"
	"FGG-Service/user_service"
	"context"
)

// GetCurrentGame (GET /users/{userId}/games/current)
func (Server) GetCurrentGame(_ context.Context, request api.GetCurrentGameRequestObject) (api.GetCurrentGameResponseObject, error) {
	doesUserExist, err := user_service.CheckIfUserExistsById(request.UserId)

	if err != nil {
		return api.GetCurrentGame503JSONResponse{Code: api.CHECKUSER, Message: err.Error()}, nil
	}

	if !*doesUserExist {
		return api.GetCurrentGame404JSONResponse{Code: api.USERNOTFOUND}, nil
	}

	game, err := game_service.GetCurrentGame(request.UserId)

	if err != nil {
		return api.GetCurrentGame503JSONResponse{Code: api.GETCURRENTGAME, Message: err.Error()}, nil
	}

	if game == nil {
		return api.GetCurrentGame404JSONResponse{Code: api.GAMENOTFOUND}, nil
	}

	return api.GetCurrentGame200JSONResponse(*game), nil
}

// FinishCurrentGame (GET /users/{userId}/games/current/finish)
func (Server) FinishCurrentGame(_ context.Context, request api.FinishCurrentGameRequestObject) (api.FinishCurrentGameResponseObject, error) {
	return api.FinishCurrentGame200Response{}, nil
}

// MakeGameRoll (GET /users/{userId}/games/roll)
func (Server) MakeGameRoll(_ context.Context, request api.MakeGameRollRequestObject) (api.MakeGameRollResponseObject, error) {
	return api.MakeGameRoll200JSONResponse{}, nil
}

// GetUnplayedGames (GET /users/{userId}/games/unplayed)
func (Server) GetUnplayedGames(_ context.Context, request api.GetUnplayedGamesRequestObject) (api.GetUnplayedGamesResponseObject, error) {
	doesUserExist, err := user_service.CheckIfUserExistsById(request.UserId)

	if err != nil {
		return api.GetUnplayedGames503JSONResponse{Code: api.CHECKUSER, Message: err.Error()}, nil
	}

	if !*doesUserExist {
		return api.GetUnplayedGames404JSONResponse{Code: api.USERNOTFOUND}, nil
	}

	games, err := game_service.GetUnplayedGames(request.UserId)

	if err != nil {
		return api.GetUnplayedGames503JSONResponse{Code: api.GETUNPLAYEDGAMES, Message: err.Error()}, nil
	}

	return api.GetUnplayedGames200JSONResponse(*games), nil
}

// AddUnplayedGames (POST /users/{userId}/games/unplayed)
func (Server) AddUnplayedGames(_ context.Context, request api.AddUnplayedGamesRequestObject) (api.AddUnplayedGamesResponseObject, error) {
	doesUserExist, err := user_service.CheckIfUserExistsById(request.UserId)

	if err != nil {
		return api.AddUnplayedGames503JSONResponse{Code: api.CHECKUSER, Message: err.Error()}, nil
	}

	if !*doesUserExist {
		return api.AddUnplayedGames404JSONResponse{Code: api.USERNOTFOUND}, nil
	}

	err = game_service.AddUnplayedGames(request.UserId, request.Body)

	if err != nil {
		return api.AddUnplayedGames503JSONResponse{Code: api.ADDUNPLAYEDGAMES, Message: err.Error()}, nil
	}

	return api.AddUnplayedGames200Response{}, nil
}
