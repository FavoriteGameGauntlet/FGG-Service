package controller

import (
	"FGG-Service/api"
	"FGG-Service/user_service"
	"context"
)

// GetUser (GET /users/{name})
func (Server) GetUser(_ context.Context, request api.GetUserRequestObject) (api.GetUserResponseObject, error) {
	doesExist, err := user_service.CheckIfUserExistsByName(request.Name)

	if err != nil {
		return api.GetUser500JSONResponse{Code: api.UNEXPECTEDDATABASE, Message: err.Error()}, nil
	}

	if !doesExist {
		return api.GetUser404JSONResponse{Code: api.USERNOTFOUND}, nil
	}

	user, err := user_service.FindUser(request.Name)

	if err != nil {
		return api.GetUser500JSONResponse{Code: api.UNEXPECTEDDATABASE, Message: err.Error()}, nil
	}

	return api.GetUser200JSONResponse(*user), nil
}

// CreateUser (POST /users/{name})
func (Server) CreateUser(_ context.Context, request api.CreateUserRequestObject) (api.CreateUserResponseObject, error) {
	doesExist, err := user_service.CheckIfUserExistsByName(request.Name)

	if err != nil {
		return api.CreateUser500JSONResponse{Code: api.UNEXPECTEDDATABASE, Message: err.Error()}, nil
	}

	if doesExist {
		return api.CreateUser409JSONResponse{Code: api.USERALREADYEXISTS}, nil
	}

	err = user_service.CreateUser(request.Name)

	if err != nil {
		return api.CreateUser500JSONResponse{Code: api.UNEXPECTEDDATABASE, Message: err.Error()}, nil
	}

	user, err := user_service.FindUser(request.Name)

	if err != nil {
		return api.CreateUser500JSONResponse{Code: api.UNEXPECTEDDATABASE, Message: err.Error()}, nil
	}

	return api.CreateUser200JSONResponse(*user), nil
}
