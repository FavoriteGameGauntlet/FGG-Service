package controller

import (
	"FGG-Service/api"
	"FGG-Service/auth_service"
	"FGG-Service/common"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SendJSONErrorResponse(ctx echo.Context, err error) error {
	var apiCode int
	var authError *common.AuthError
	var notFoundError *common.NotFoundError
	var currentTimerIncorrectStateError *common.CurrentTimerIncorrectStateError
	var databaseError *common.DatabaseError

	switch {
	case errors.As(err, &authError):
		apiCode = http.StatusUnauthorized
	case errors.As(err, &notFoundError):
		apiCode = http.StatusNotFound
	case errors.As(err, &currentTimerIncorrectStateError):
		apiCode = http.StatusConflict
	case errors.As(err, &databaseError):
		apiCode = http.StatusInternalServerError
	default:
		apiCode = http.StatusInternalServerError
	}

	apiError := ConvertToError(err)

	return ctx.JSON(apiCode, apiError)
}

func ConvertToError(err error) api.Error {
	var appError common.AppError
	if errors.As(err, &appError) {
		return api.Error{
			Code:    appError.GetCode(),
			Message: appError.GetMessage(),
		}
	}

	return api.Error{
		Code:    "UNEXPECTED",
		Message: err.Error(),
	}
}

func GetUserId(ctx echo.Context) (int, error) {
	var userId int

	cookie, err := ctx.Cookie("sessionId")

	if err != nil {
		return userId, common.NewCookieNotFoundAuthError()
	}

	sessionId := cookie.Value

	userId, err = auth_service.GetUserId(sessionId)

	if err != nil {
		return userId, err
	}

	return userId, nil
}

func GetSessionCookie(ctx echo.Context) (*http.Cookie, error) {
	cookie, err := ctx.Cookie("sessionId")

	if err != nil {
		return nil, common.NewCookieNotFoundAuthError()
	}

	return cookie, nil
}

func CheckIfSessionExists(ctx echo.Context) (bool, error) {
	cookie, err := GetSessionCookie(ctx)

	if err != nil {
		return false, common.NewCookieNotFoundAuthError()
	}

	sessionId := cookie.Value
	doesExist, err := auth_service.CheckIfUserSessionExists(sessionId)

	if err != nil {
		return false, common.NewCookieNotFoundAuthError()
	}

	return doesExist, nil
}
