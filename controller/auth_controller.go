package controller

import (
	"FGG-Service/api"
	"FGG-Service/auth_service"
	"FGG-Service/common"
	"FGG-Service/validator"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// Login (POST /auth/login)
func (Server) Login(ctx echo.Context) error {
	var user api.User
	err := ctx.Bind(&user)

	if err != nil {
		err = common.NewBadRequestError(err.Error())
		return SendJSONErrorResponse(ctx, err)
	}

	doesExist, _ := CheckIfSessionExists(ctx)

	if doesExist {
		return SendJSONErrorResponse(ctx, common.NewSessionAlreadyExistsConflictError())
	}

	sessionId, err := auth_service.CreateSession(user.Name, user.Password)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	cookie := CreateSessionCookie(*sessionId)
	ctx.SetCookie(cookie)

	return ctx.NoContent(http.StatusOK)
}

func CreateSessionCookie(sessionId string) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = "sessionId"
	cookie.Value = sessionId
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HttpOnly = true
	cookie.Path = "/"
	cookie.Secure = false

	return cookie
}

// Logout (POST /auth/logout)
func (Server) Logout(ctx echo.Context) error {
	cookie, err := GetSessionCookie(ctx)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	sessionId := cookie.Value
	err = auth_service.DeleteSession(sessionId)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	cookie.MaxAge = -1
	ctx.SetCookie(cookie)

	return ctx.NoContent(http.StatusOK)
}

// SignUp (POST /auth/signup)
func (Server) SignUp(ctx echo.Context) error {
	var user api.NewUser
	err := ctx.Bind(&user)

	if err != nil {
		err = common.NewBadRequestError(err.Error())
		return SendJSONErrorResponse(ctx, err)
	}

	err = validator.ValidateUserName(user.Name)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	err = validator.ValidateEmail(user.Email)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	err = validator.ValidatePassword(user.Password)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	err = auth_service.CreateUser(user.Name, user.Email, user.Password)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	return ctx.NoContent(http.StatusOK)
}
