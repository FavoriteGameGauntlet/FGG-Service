package ctrlusers

import (
	"FGG-Service/src/users/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	Service srvusers.Service
}

func NewController() *Controller {
	return new(Controller)
}

// GetAllUserNames (GET /users/all/names)
func (c *Controller) GetAllUserNames(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// GetDisplayName (GET /users/display-name)
func (c *Controller) GetDisplayName(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

// ChangeDisplayName (POST /users/display-name)
func (c *Controller) ChangeDisplayName(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}
