package auth

import (
	"FGG-Service/src/auth/service"

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

	return nil
}

// Logout (POST /auth/logout)
func (c *Controller) Logout(ctx echo.Context) error {

	return nil
}

// SignUp (POST /auth/signup)
func (c *Controller) SignUp(ctx echo.Context) error {

	return nil
}

//// Login (POST /auth/login)
//func (controller.Server) Login(ctx echo.Context) error {
//	var user api.LoginUser
//	err := ctx.Bind(&user)
//
//	if err != nil {
//		err = common.NewBadRequestError(err.Error())
//		return controller.SendJSONErrorResponse(ctx, err)
//	}
//
//	doesExist, _ := controller.doesUserSessionExist(ctx)
//
//	if doesExist {
//		err = common.NewSessionAlreadyExistsConflictError()
//		return controller.SendJSONErrorResponse(ctx, err)
//	}
//
//	userSession, err := auth_service.CreateSession(user.Login, user.Password)
//
//	if err != nil {
//		return controller.SendJSONErrorResponse(ctx, err)
//	}
//
//	cookie := createSessionCookie(userSession.Id)
//	ctx.SetCookie(cookie)
//
//	return ctx.NoContent(http.StatusNoContent)
//}
//
//func createSessionCookie(sessionId string) *http.Cookie {
//	cookie := new(http.Cookie)
//	cookie.Name = controller.SessionCookieName
//	cookie.Value = sessionId
//	cookie.Expires = time.Now().Add(24 * time.Hour)
//	cookie.HttpOnly = true
//	cookie.Path = "/"
//	cookie.Secure = false
//
//	return cookie
//}
//
//// Logout (POST /auth/logout)
//func (controller.Server) Logout(ctx echo.Context) error {
//	cookie, err := controller.getSessionCookie(ctx)
//
//	if err != nil {
//		return controller.SendJSONErrorResponse(ctx, err)
//	}
//
//	sessionId := cookie.Value
//	err = auth_service.DeleteUserSession(sessionId)
//
//	if err != nil {
//		return controller.SendJSONErrorResponse(ctx, err)
//	}
//
//	cookie.MaxAge = -1
//	ctx.SetCookie(cookie)
//
//	return ctx.NoContent(http.StatusNoContent)
//}
//
//// SignUp (POST /auth/signup)
//func (controller.Server) SignUp(ctx echo.Context) error {
//	var user api.SignupUser
//	err := ctx.Bind(&user)
//
//	if err != nil {
//		err = common.NewBadRequestError(err.Error())
//		return controller.SendJSONErrorResponse(ctx, err)
//	}
//
//	err = validator.ValidateUserLogin(user.Login)
//
//	if err != nil {
//		return controller.SendJSONErrorResponse(ctx, err)
//	}
//
//	err = validator.ValidateEmail(user.Email)
//
//	if err != nil {
//		return controller.SendJSONErrorResponse(ctx, err)
//	}
//
//	err = validator.ValidatePassword(user.Password)
//
//	if err != nil {
//		return controller.SendJSONErrorResponse(ctx, err)
//	}
//
//	err = auth_service.CreateUser(user.Login, user.Email, user.Password)
//
//	if err != nil {
//		return controller.SendJSONErrorResponse(ctx, err)
//	}
//
//	return ctx.NoContent(http.StatusNoContent)
//}
