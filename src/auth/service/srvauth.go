package srvauth

import (
	"FGG-Service/src/auth/database"
	"FGG-Service/src/common"
	"database/sql"
	"errors"
)

type Service struct {
	Database dbauth.Database
}

func (s *Service) CreateUser(login string, email string, password string) error {
	user, err := s.Database.GetUserByNameCommand(login)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if user.Login != "" {
		return common.NewUserNameAlreadyExistsError()
	}

	user, err = s.Database.GetUserByEmailCommand(email)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if user.Email != "" {
		return common.NewUserEmailAlreadyExistsError()
	}

	err = s.Database.CreateUserCommand(login, email, password)

	return err
}

func (s *Service) GetUserSessionById(sessionId string) (userSession common.UserSession, err error) {
	userSession, err = s.Database.GetUserSessionByIdCommand(sessionId)

	return
}

func (s *Service) CreateSession(userLogin string, userPassword string) (userSession common.UserSession, err error) {
	user, err := s.Database.GetUserByLoginAndPasswordCommand(userLogin, userPassword)

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
