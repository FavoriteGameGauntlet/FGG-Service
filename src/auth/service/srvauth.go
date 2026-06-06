package srvauth

import (
	"FGG-Service/src/auth/database"
	"FGG-Service/src/auth/types"
	"FGG-Service/src/common"
	"database/sql"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type IService interface {
	DoesUserSessionExist(ctx echo.Context) (bool, error)
	GetUserId(ctx echo.Context) (int, error)
	GetUserIdByLogin(login string) (int, error)
}

type Service struct {
	Database dbauth.Database
}

func NewService() *Service {
	db := new(dbauth.Database)

	return &Service{
		Database: *db,
	}
}

func (s *Service) DoesUserSessionExist(ctx echo.Context) (doesExist bool, err error) {
	cookie, err := s.GetSessionCookie(ctx)

	if err != nil {
		err = common.NewCookieNotFoundUnauthorizedError()
		return
	}

	sessionId := cookie.Value

	_, err = s.GetUserSessionById(sessionId)

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

func (s *Service) GetSessionCookie(ctx echo.Context) (*http.Cookie, error) {
	cookie, err := ctx.Cookie(common.SessionCookieName)

	if err != nil {
		return nil, common.NewCookieNotFoundUnauthorizedError()
	}

	return cookie, nil
}

func (s *Service) GetUserId(ctx echo.Context) (userId int, err error) {
	cookie, err := ctx.Cookie(common.SessionCookieName)

	if err != nil {
		err = common.NewCookieNotFoundUnauthorizedError()
		return
	}

	sessionId := cookie.Value

	userSession, err := s.GetUserSessionById(sessionId)

	if errors.Is(err, sql.ErrNoRows) {
		err = common.NewActiveSessionNotFoundUnauthorizedError()
		return
	}

	if err != nil {
		return
	}

	userId = userSession.UserId

	return
}

func (s *Service) CreateUser(signupUser typeauth.SignupUser) error {
	user, err := s.Database.GetUserByLoginCommand(signupUser.Login)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if user.Login != "" {
		return common.NewUserNameAlreadyExistsConflictError()
	}

	user, err = s.Database.GetUserByEmailCommand(signupUser.Email)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if user.Email != "" {
		return common.NewUserEmailAlreadyExistsConflictError()
	}

	err = s.Database.CreateUserCommand(signupUser)

	if err != nil {
		return err
	}

	err = s.Database.CreateUserStatsCommand(signupUser.Login)

	return err
}

func (s *Service) GetUserSessionById(sessionId string) (userSession typeauth.UserSession, err error) {
	userSession, err = s.Database.GetUserSessionByIdCommand(sessionId)

	return
}

func (s *Service) CreateSession(loginUser typeauth.LoginUser) (userSession typeauth.UserSession, err error) {
	user, err := s.Database.GetUserByLoginAndPasswordCommand(loginUser)

	if errors.Is(err, sql.ErrNoRows) {
		err = common.NewWrongDataUnprocessableError()
		return
	}

	if err != nil {
		return
	}

	userSession, err = s.Database.CreateUserSessionCommand(user.Id)

	return
}

func (s *Service) DeleteUserSession(userSessionId string) error {
	err := s.Database.DeleteUserSessionCommand(userSessionId)

	return err
}

func (s *Service) GetUserIdByLogin(userLogin string) (userId int, err error) {
	user, err := s.Database.GetUserByLoginCommand(userLogin)

	if errors.Is(err, sql.ErrNoRows) {
		err = common.NewUserLoginNotFoundError(userLogin)
		return
	}

	if err != nil {
		return
	}

	userId = user.Id
	return
}
