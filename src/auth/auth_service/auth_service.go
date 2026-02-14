package auth_service

import (
	"FGG-Service/src/auth/auth_db"
	"FGG-Service/src/common"
	"database/sql"
	"errors"
)

func CreateUser(userName string, userEmail string, userPassword string) error {
	user, err := auth_db.GetUserByNameCommand(userName)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if user.Login != "" {
		return common.NewUserNameAlreadyExistsError()
	}

	user, err = auth_db.GetUserByEmailCommand(userEmail)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if user.Email != "" {
		return common.NewUserEmailAlreadyExistsError()
	}

	err = auth_db.CreateUserCommand(userName, userEmail, userPassword)

	return err
}

func GetUserSessionById(sessionId string) (userSession common.UserSession, err error) {
	userSession, err = auth_db.GetUserSessionByIdCommand(sessionId)

	return
}

func CreateSession(userName string, userPassword string) (userSession common.UserSession, err error) {
	user, err := auth_db.GetUserByNameAndPasswordCommand(userName, userPassword)

	if errors.Is(err, sql.ErrNoRows) {
		err = common.NewWrongDataUnprocessableError()
		return
	}

	if err != nil {
		return
	}

	userSession, err = auth_db.CreateUserSessionCommand(user.Id)

	return
}

func DeleteUserSession(userSessionId string) error {
	err := auth_db.DeleteUserSessionCommand(userSessionId)

	return err
}
