package controller

import (
	"FavoriteGameGauntlet/api"
	"FavoriteGameGauntlet/user"
	"context"
)

type Server struct{}

func NewServer() Server {
	return Server{}
}

// (GET /users/{name})
func (Server) GetUsersName(ctx context.Context, request api.GetUsersNameRequestObject) (api.GetUsersNameResponseObject, error) {
	user, _ := user.FindUser(request.Name)

	if user == nil {
		return api.GetUsersName404Response{}, nil
	}

	return api.GetUsersName200JSONResponse(*user), nil
}

// (POST /users/{name})
func (Server) PostUsersName(ctx context.Context, request api.PostUsersNameRequestObject) (api.PostUsersNameResponseObject, error) {
	return api.PostUsersName200Response{}, nil
}
