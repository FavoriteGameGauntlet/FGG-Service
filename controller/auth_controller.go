package controller

import (
	"FGG-Service/api"
	"FGG-Service/auth_service"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// Login (POST /auth/login)
func (Server) Login(ctx echo.Context) error {
	var user api.User
	if err := ctx.Bind(&user); err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	doesExist, err := auth_service.CheckIfUserNameExists(user.Name)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	if !doesExist {
		return ctx.JSON(http.StatusUnauthorized, api.NotAuthorizedError{Code: api.WRONGAUTHDATA})
	}

	cookie, err := ctx.Cookie("sessionId")

	if err == nil {
		sessionId := cookie.Value
		doesExist, err = auth_service.CheckIfUserSessionExists(sessionId)

		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
		}

		if doesExist {
			return ctx.JSON(http.StatusConflict, api.ConflictError{Code: api.SESSIONALREADYEXISTS})
		}
	}

	sessionId, err := auth_service.CreateSession(user.Name, user.Password)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	if sessionId == nil {
		return ctx.JSON(http.StatusUnauthorized, api.NotAuthorizedError{Code: api.WRONGAUTHDATA})
	}

	cookie = new(http.Cookie)
	cookie.Name = "sessionId"
	cookie.Value = *sessionId
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HttpOnly = true

	ctx.SetCookie(cookie)

	return ctx.NoContent(http.StatusOK)
}

// Logout (POST /auth/logout)
func (Server) Logout(ctx echo.Context) error {
	cookie, err := ctx.Cookie("sessionId")

	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, api.NotAuthorizedError{Code: api.NOACTIVESESSION})
	}

	sessionId := cookie.Value

	doesExist, err := auth_service.CheckIfUserSessionExists(sessionId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	if !doesExist {
		return ctx.JSON(http.StatusUnauthorized, api.NotAuthorizedError{Code: api.NOACTIVESESSION})
	}

	if err = auth_service.DeleteSession(sessionId); err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	cookie.MaxAge = -1
	ctx.SetCookie(cookie)

	return ctx.NoContent(http.StatusOK)
}

// SignUp (POST /auth/signup)
func (Server) SignUp(ctx echo.Context) error {
	var user api.NewUser
	if err := ctx.Bind(&user); err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	doesExist, err := auth_service.CheckIfUserNameExists(user.Name)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	if doesExist {
		return ctx.JSON(http.StatusConflict, api.ConflictError{Code: api.USERNAMEALREADYEXISTS})
	}

	doesExist, err = auth_service.CheckIfUserEmailExists(user.Email)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	if doesExist {
		return ctx.JSON(http.StatusConflict, api.ConflictError{Code: api.EMAILALREADYEXISTS})
	}

	err = auth_service.CreateUser(user.Name, user.Email, user.Password)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	return ctx.NoContent(http.StatusOK)
}
