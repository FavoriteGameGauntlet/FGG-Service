package controller

import (
	"FGG-Service/api"
	"FGG-Service/auth_service"
	"FGG-Service/effect_service"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetAvailableEffects (POST /effects/available)
func (Server) GetAvailableEffects(ctx echo.Context) error {
	cookie, err := ctx.Cookie("sessionId")

	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, api.NotAuthorizedError{Code: api.NOACTIVESESSION})
	}

	sessionId := cookie.Value

	_, err = auth_service.GetUserId(sessionId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	return ctx.NoContent(http.StatusOK)
}

// CheckAvailableEffectRoll CheckEffectRoll (GET /effects/has-roll)
func (Server) CheckAvailableEffectRoll(ctx echo.Context) error {
	cookie, err := ctx.Cookie("sessionId")

	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, api.NotAuthorizedError{Code: api.NOACTIVESESSION})
	}

	sessionId := cookie.Value

	userId, err := auth_service.GetUserId(sessionId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	doesExist, err := effect_service.CheckIfAvailableRollExists(userId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, doesExist)
}

// GetEffectHistory (GET /effects/history)
func (Server) GetEffectHistory(ctx echo.Context) error {
	cookie, err := ctx.Cookie("sessionId")

	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, api.NotAuthorizedError{Code: api.NOACTIVESESSION})
	}

	sessionId := cookie.Value

	_, err = auth_service.GetUserId(sessionId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	return ctx.NoContent(http.StatusOK)
}

// MakeEffectRoll (POST /effects/roll)
func (Server) MakeEffectRoll(ctx echo.Context) error {
	cookie, err := ctx.Cookie("sessionId")

	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, api.NotAuthorizedError{Code: api.NOACTIVESESSION})
	}

	sessionId := cookie.Value

	_, err = auth_service.GetUserId(sessionId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.InternalServerError{Code: api.UNEXPECTED, Message: err.Error()})
	}

	return ctx.NoContent(http.StatusOK)
}
