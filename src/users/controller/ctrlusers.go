package ctrlusers

import (
	genusers "FGG-Service/api/generated/users"
	srvauth "FGG-Service/src/auth/service"
	"FGG-Service/src/common"
	"FGG-Service/src/users/service"
	"FGG-Service/src/users/types"
	"FGG-Service/src/validator"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	Service     srvusers.Service
	AuthService srvauth.Service
}

func NewController() *Controller {
	s := srvusers.NewService()
	as := new(srvauth.Service)

	return &Controller{*s, *as}
}

// GetAllUserNames (GET /users/all/names)
func (c *Controller) GetAllUserNames(ctx echo.Context) error {
	doesExist, err := c.AuthService.DoesUserSessionExist(ctx)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	if !doesExist {
		err = common.NewActiveSessionNotFoundUnauthorizedError()
		return common.SendJSONErrorResponse(ctx, err)
	}

	users, err := c.Service.GetAllUserNames()

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	usersDto := convertUsersToDto(users)

	return ctx.JSON(http.StatusOK, usersDto)
}

func convertUsersToDto(users typeusers.Users) genusers.UserNames {
	dto := make(genusers.UserNames, len(users))

	for i, user := range users {
		dto[i] = genusers.UserName{
			Login:       user.Login,
			DisplayName: user.DisplayName,
		}
	}

	return dto
}

// GetDisplayName (GET /users/display-name)
func (c *Controller) GetDisplayName(ctx echo.Context) error {
	userId, err := c.AuthService.GetUserId(ctx)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	displayName, err := c.Service.GetDisplayName(userId)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, displayName)
}

// ChangeDisplayName (POST /users/display-name)
func (c *Controller) ChangeDisplayName(ctx echo.Context) error {
	userId, err := c.AuthService.GetUserId(ctx)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	var changeName genusers.ChangeName
	err = ctx.Bind(&changeName)

	if err != nil {
		err = common.NewBadRequestError(err.Error())
		return common.SendJSONErrorResponse(ctx, err)
	}

	err = validator.ValidateName(changeName.Name)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	err = c.Service.ChangeDisplayName(userId, changeName.Name)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	return ctx.NoContent(http.StatusNoContent)
}
