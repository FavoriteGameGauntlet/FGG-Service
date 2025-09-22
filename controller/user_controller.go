package controller

import (
	"FavoriteGameGauntlet/api"
	"FavoriteGameGauntlet/user"
	"context"
)

// (GET /users/{name})
func (Server) GetUsersName(ctx context.Context, request api.GetUsersNameRequestObject) (api.GetUsersNameResponseObject, error) {
	doesUserExist, err := user.CheckIfUserExists(request.Name)

	if doesUserExist == nil {
		return api.GetUsersName503JSONResponse{api.DATABASEQUERY, err.Error()}, nil
	}

	if !*doesUserExist {
		return api.GetUsersName404JSONResponse{Code: api.USERNOTFOUND}, nil
	}

	user, err := user.FindUser(request.Name)

	if user == nil {
		return api.GetUsersName404JSONResponse{api.USERNOTFOUND, err.Error()}, nil
	}

	return api.GetUsersName200JSONResponse(*user), nil
}

// (POST /users/{name})
func (Server) PostUsersName(ctx context.Context, request api.PostUsersNameRequestObject) (api.PostUsersNameResponseObject, error) {
	doesUserExist, err := user.CheckIfUserExists(request.Name)

	if doesUserExist == nil {
		return api.PostUsersName503JSONResponse{api.DATABASEQUERY, err.Error()}, nil
	}

	if *doesUserExist {
		return api.PostUsersName409JSONResponse{Code: api.USERALREADYEXISTS}, nil
	}

	err = user.AddUser(request.Name)

	if err != nil {
		return api.PostUsersName503JSONResponse{api.DATABASEQUERY, err.Error()}, nil
	}

	user, err := user.FindUser(request.Name)

	if user == nil {
		return api.PostUsersName503JSONResponse{api.DATABASEQUERY, err.Error()}, nil
	}

	return api.PostUsersName200JSONResponse(*user), nil
}
