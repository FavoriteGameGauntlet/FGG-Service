package ctrlauth

import (
	"FGG-Service/api/generated/auth"
	"FGG-Service/src/auth/service"
	"FGG-Service/src/common"
	"FGG-Service/src/validator"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	Service srvauth.Service
}

func NewController() Controller {
	return Controller{}
}

// Login (POST /auth/login)
func (c *Controller) Login(ctx echo.Context) error {
	var user genauth.LoginUser
	err := ctx.Bind(&user)

	if err != nil {
		err = common.NewBadRequestError(err.Error())
		return common.SendJSONErrorResponse(ctx, err)
	}

	doesExist, _ := c.Service.DoesUserSessionExist(ctx)

	if doesExist {
		err = common.NewSessionAlreadyExistsConflictError()
		return common.SendJSONErrorResponse(ctx, err)
	}

	userSession, err := c.Service.CreateSession(user.Login, user.Password)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	cookie := createSessionCookie(userSession.Id)
	ctx.SetCookie(cookie)

	return ctx.NoContent(http.StatusNoContent)
}

func createSessionCookie(sessionId string) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = common.SessionCookieName
	cookie.Value = sessionId
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HttpOnly = true
	cookie.Path = "/"
	cookie.Secure = false

	return cookie
}

// Logout (POST /auth/logout)
func (c *Controller) Logout(ctx echo.Context) error {
	cookie, err := c.Service.GetSessionCookie(ctx)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	sessionId := cookie.Value
	err = c.Service.DeleteUserSession(sessionId)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	cookie.MaxAge = -1
	ctx.SetCookie(cookie)

	return ctx.NoContent(http.StatusNoContent)
}

// SignUp (POST /auth/signup)
func (c *Controller) SignUp(ctx echo.Context) error {
	var user genauth.SignupUser
	err := ctx.Bind(&user)

	if err != nil {
		err = common.NewBadRequestError(err.Error())
		return common.SendJSONErrorResponse(ctx, err)
	}

	err = validator.ValidateUserLogin(user.Login)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	err = validator.ValidateEmail(user.Email)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	err = validator.ValidatePassword(user.Password)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	err = c.Service.CreateUser(user.Login, user.Email, user.Password)

	if err != nil {
		return common.SendJSONErrorResponse(ctx, err)
	}

	return ctx.NoContent(http.StatusNoContent)
}
