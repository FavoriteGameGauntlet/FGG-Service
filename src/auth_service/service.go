package auth_service

import (
	"FGG-Service/src/common"
	"database/sql"
	"errors"
)

func CreateUser(userName string, userEmail string, userPassword string) error {
	user, err := GetUserByNameCommand(userName)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if user.Name != "" {
		return common.NewUserNameAlreadyExistsError()
	}

	user, err = GetUserByEmailCommand(userEmail)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if user.Email != "" {
		return common.NewUserEmailAlreadyExistsError()
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
