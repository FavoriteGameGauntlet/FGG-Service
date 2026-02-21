package auth_service

import (
	"FGG-Service/src/auth/db"
	"FGG-Service/src/common"
	"database/sql"
	"errors"
)

type Service struct {
	Db auth_db.Database
}

func (s *Service) CreateUser(login string, email string, password string) error {
	user, err := s.Db.GetUserByNameCommand(login)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if user.Login != "" {
		return common.NewUserNameAlreadyExistsError()
	}

	user, err = s.Db.GetUserByEmailCommand(email)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if user.Email != "" {
		return common.NewUserEmailAlreadyExistsError()
	}

	err = s.Db.CreateUserCommand(login, email, password)

	return err
}

func (s *Service) GetUserSessionById(sessionId string) (userSession common.UserSession, err error) {
	userSession, err = s.Db.GetUserSessionByIdCommand(sessionId)

	return
}

func (s *Service) CreateSession(userLogin string, userPassword string) (userSession common.UserSession, err error) {
	user, err := s.Db.GetUserByLoginAndPasswordCommand(userLogin, userPassword)

	if errors.Is(err, sql.ErrNoRows) {
		err = common.NewWrongDataUnprocessableError()
		return
	}

	if err != nil {
		return
	}

	userSession, err = s.Db.CreateUserSessionCommand(user.Id)

	return
}

func (s *Service) DeleteUserSession(userSessionId string) error {
	err := s.Db.DeleteUserSessionCommand(userSessionId)

	return err
}
