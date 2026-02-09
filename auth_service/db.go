package auth_service

import (
	"FGG-Service/common"
	"FGG-Service/db_access"

	"github.com/google/uuid"
)

const GetUserByNameQuery = `
	SELECT Id, Name, Email
	FROM Users
	WHERE Name = ?
`

func GetUserByNameCommand(userName string) (user common.User, err error) {
	row := db_access.QueryRow(GetUserByNameQuery, userName)

	err = row.Scan(&user.Id, &user.Name, &user.Email)

	return
}

const GetUserByEmailQuery = `
	SELECT Id, Name, Email
	FROM Users
	WHERE Email = ?
`

func GetUserByEmailCommand(userEmail string) (user common.User, err error) {
	row := db_access.QueryRow(GetUserByEmailQuery, userEmail)

	err = row.Scan(&user.Id, &user.Name, &user.Email)

	return
}

const GetUserByNameAndPasswordQuery = `
	SELECT Id, Name, Email
	FROM Users
	WHERE Name = ?
		AND Password = ?
`

func GetUserByNameAndPasswordCommand(userName string, userPassword string) (user common.User, err error) {
	row := db_access.QueryRow(GetUserByNameAndPasswordQuery, userName, userPassword)

	err = row.Scan(&user.Id, &user.Name, &user.Email)

	return
}

const CreateUserQuery = `
	INSERT INTO Users (Name, Email, Password)
	VALUES (?, ?, ?)
`

func CreateUserCommand(userName string, userEmail string, userPassword string) error {
	_, err := db_access.Exec(CreateUserQuery, userName, userEmail, userPassword)

	return err
}

const GetUserSessionByIdQuery = `
	SELECT Id, UserId
	FROM UserSessions
	WHERE Id = ?
`

func GetUserSessionByIdCommand(userSessionId string) (user common.UserSession, err error) {
	row := db_access.QueryRow(GetUserSessionByIdQuery, userSessionId)

	var userSession common.UserSession
	err = row.Scan(&userSession.Id, &userSession.UserId)

	return
}

const CreateUserSessionQuery = `
	INSERT INTO UserSessions (Id, UserId)
	VALUES (?, ?)
`

func CreateUserSessionCommand(userId int) (userSession common.UserSession, err error) {
	sessionId := uuid.New().String()

	_, err = db_access.Exec(CreateUserSessionQuery, sessionId, userId)

	if err != nil {
		return
	}

	userSession = common.UserSession{
		Id:     sessionId,
		UserId: userId,
	}

	return
}

const DeleteUserSessionQuery = `
	DELETE FROM UserSessions
	WHERE Id = ?
`

func DeleteUserSessionCommand(userSessionId string) error {
	_, err := db_access.Exec(DeleteUserSessionQuery, userSessionId)

	return err
}
