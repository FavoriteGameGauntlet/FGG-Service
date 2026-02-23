package ctrlusers

import (
	"FGG-Service/src/users/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	Service srvusers.Service
}

func NewController() Controller {
	return Controller{}
}

// GetUserInfos (GET /users/names)
func (c *Controller) GetUserInfos(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetOwnDisplayName (GET /users/self/display-name)
func (c *Controller) GetOwnDisplayName(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// ChangeOwnDisplayName (POST /users/self/display-name)
func (c *Controller) ChangeOwnDisplayName(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}
