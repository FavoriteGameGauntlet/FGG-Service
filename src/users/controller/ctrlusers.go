package users

import (
	"FGG-Service/src/users/service"

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

	return nil
}

// GetOwnDisplayName (GET /users/self/display-name)
func (c *Controller) GetOwnDisplayName(ctx echo.Context) error {

	return nil
}

// ChangeOwnDisplayName (POST /users/self/display-name)
func (c *Controller) ChangeOwnDisplayName(ctx echo.Context) error {

	return nil
}
