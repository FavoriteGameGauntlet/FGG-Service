package auth_service

import (
	"FGG-Service/common"
	"database/sql"
	"errors"
)

func CreateUser(userName string, userEmail string, userPassword string) error {
	_, err := GetUserByNameCommand(userName)

	if errors.Is(err, sql.ErrNoRows) {
		return common.NewUserNameAlreadyExistsError()
	}

	if err != nil {
		return err
	}

	_, err = GetUserByEmailCommand(userEmail)

	if errors.Is(err, sql.ErrNoRows) {
		return common.NewUserNameAlreadyExistsError()
	}

	if err != nil {
		return err
	}

	err = CreateUserCommand(userName, userEmail, userPassword)

	return err
}

func GetUserSessionById(sessionId string) (userSession common.UserSession, err error) {
	userSession, err = GetUserSessionByIdCommand(sessionId)

	return
}

func CreateSession(userName string, userPassword string) (userSession common.UserSession, err error) {
	user, err := GetUserByNameAndPasswordCommand(userName, userPassword)

	if errors.Is(err, sql.ErrNoRows) {
		err = common.NewWrongDataUnauthorizedError()
		return
	}

	if err != nil {
		return
	}

	userSession, err = CreateUserSessionCommand(user.Id)

	return
}

func DeleteUserSession(userSessionId string) error {
	err := DeleteUserSessionCommand(userSessionId)

	return err
}
