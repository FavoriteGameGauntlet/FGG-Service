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
	var badRequestError *common.BadRequestError
	var unauthorizedError *common.UnauthorizedError
	var notFoundError *common.NotFoundError
	var conflictError *common.ConflictError

	apiCode := http.StatusInternalServerError

	switch {
	case errors.As(err, &badRequestError):
		apiCode = http.StatusBadRequest
	case errors.As(err, &unauthorizedError):
		apiCode = http.StatusUnauthorized
	case errors.As(err, &notFoundError):
		apiCode = http.StatusNotFound
	case errors.As(err, &conflictError):
		apiCode = http.StatusConflict
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
		return userId, common.NewCookieNotFoundUnauthorizedError()
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
		return nil, common.NewCookieNotFoundUnauthorizedError()
	}

	return cookie, nil
}

func CheckIfSessionExists(ctx echo.Context) (bool, error) {
	cookie, err := GetSessionCookie(ctx)

	if err != nil {
		return false, common.NewCookieNotFoundUnauthorizedError()
	}

	sessionId := cookie.Value
	doesExist, err := auth_service.CheckIfUserSessionExists(sessionId)

	if err != nil {
		return false, common.NewCookieNotFoundUnauthorizedError()
	}

	return doesExist, nil
}
