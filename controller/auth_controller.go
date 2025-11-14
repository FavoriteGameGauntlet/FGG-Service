package controller

import (
	"FGG-Service/api"
	"FGG-Service/auth_service"
	"context"
)

// Login (POST /auth/login)
func (Server) Login(ctx context.Context, request api.LoginRequestObject) (api.LoginResponseObject, error) {
	doesExist, err := auth_service.CheckIfUserNameExists(request.Body.Name)

	if err != nil {
		return api.Login500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	if !doesExist {
		return api.Login401JSONResponse{Code: api.WRONGAUTHDATA}, nil
	}

	isSuccess, err := auth_service.CreateSession(request.Body.Name, request.Body.Password)

	if err != nil {
		return api.Login500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	if !isSuccess {
		return api.Login401JSONResponse{Code: api.WRONGAUTHDATA}, nil
	}

	return api.Login200Response{}, nil
}

// Logout (POST /auth/logout)
func (Server) Logout(ctx context.Context, _ api.LogoutRequestObject) (api.LogoutResponseObject, error) {
	sessionId, ok := ctx.Value("session_id").(string)

	if !ok {
		return api.Logout401JSONResponse{Code: api.NOACTIVESESSION}, nil
	}

	doesExist, err := auth_service.CheckIfUserSessionExists(sessionId)

	if err != nil {
		return api.Logout500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	if doesExist {
		return api.Logout401JSONResponse{Code: api.NOACTIVESESSION}, nil
	}

	err = auth_service.DeleteSession(sessionId)

	if err != nil {
		return api.Logout500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	return api.Logout200Response{}, nil
}

// SignUp (POST /auth/signup)
func (Server) SignUp(_ context.Context, request api.SignUpRequestObject) (api.SignUpResponseObject, error) {
	doesExist, err := auth_service.CheckIfUserNameExists(request.Body.Name)

	if err != nil {
		return api.SignUp500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	if doesExist {
		return api.SignUp409JSONResponse{Code: api.USERNAMEALREADYEXISTS}, nil
	}

	doesExist, err = auth_service.CheckIfUserEmailExists(request.Body.Email)

	if err != nil {
		return api.SignUp500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	if doesExist {
		return api.SignUp409JSONResponse{Code: api.EMAILALREADYEXISTS}, nil
	}

	err = auth_service.CreateUser(request.Body.Name, request.Body.Email, request.Body.Password)

	if err != nil {
		return api.SignUp500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	return api.SignUp200Response{}, nil
}
