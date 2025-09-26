package controller

import (
	"FGG-Service/api"
	"FGG-Service/user_service"
	"context"
)

// (GET /users/{name})
func (Server) GetUser(ctx context.Context, request api.GetUserRequestObject) (api.GetUserResponseObject, error) {
	doesUserExist, err := user_service.CheckIfUserExistsByName(request.Name)

	if err != nil {
		return api.GetUser503JSONResponse{api.DATABASEQUERY, err.Error()}, nil
	}

	if !*doesUserExist {
		return api.GetUser404JSONResponse{Code: api.USERNOTFOUND}, nil
	}

	user, err := user_service.FindUser(request.Name)

	if user == nil {
		return api.GetUser404JSONResponse{api.USERNOTFOUND, err.Error()}, nil
	}

	return api.GetUser200JSONResponse(*user), nil
}

// (POST /users/{name})
func (Server) CreateUser(ctx context.Context, request api.CreateUserRequestObject) (api.CreateUserResponseObject, error) {
	doesUserExist, err := user_service.CheckIfUserExistsByName(request.Name)

	if err != nil {
		return api.CreateUser503JSONResponse{api.DATABASEQUERY, err.Error()}, nil
	}

	if *doesUserExist {
		return api.CreateUser409JSONResponse{Code: api.USERALREADYEXISTS}, nil
	}

	err = user_service.CreateUser(request.Name)

	if err != nil {
		return api.CreateUser503JSONResponse{api.DATABASEQUERY, err.Error()}, nil
	}

	user, err := user_service.FindUser(request.Name)

	if user == nil {
		return api.CreateUser503JSONResponse{api.DATABASEQUERY, err.Error()}, nil
	}

	return api.CreateUser200JSONResponse(*user), nil
}
