package controller

import (
	"FGG-Service/api"
	"FGG-Service/auth_service"
	"FGG-Service/common"
	"database/sql"
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

func GetUserId(ctx echo.Context) (userId int, err error) {
	cookie, err := ctx.Cookie(SessionCookieName)

	if err != nil {
		err = common.NewCookieNotFoundUnauthorizedError()
		return
	}

	sessionId := cookie.Value

	userSession, err := auth_service.GetUserSessionById(sessionId)

	if err != nil {
		return
	}

	userId = userSession.UserId

	return
}

func GetSessionCookie(ctx echo.Context) (*http.Cookie, error) {
	cookie, err := ctx.Cookie(SessionCookieName)

	if err != nil {
		return nil, common.NewCookieNotFoundUnauthorizedError()
	}

	return cookie, nil
}

func DoesUserSessionExist(ctx echo.Context) (doesExist bool, err error) {
	cookie, err := GetSessionCookie(ctx)

	if err != nil {
		err = common.NewCookieNotFoundUnauthorizedError()
		return
	}

	sessionId := cookie.Value
	_, err = auth_service.GetUserSessionById(sessionId)

	if errors.Is(err, sql.ErrNoRows) {
		err = nil
		return
	}

	if err != nil {
		return
	}

	doesExist = true
	return
}
