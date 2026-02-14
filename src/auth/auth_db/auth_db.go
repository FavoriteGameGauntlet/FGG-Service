package auth_db

import (
	"FGG-Service/src/common"
	"FGG-Service/src/db_access"

	"github.com/google/uuid"
)

const GetUserByNameQuery = `
	SELECT Id, Login, Nickname, Email
	FROM Users
	WHERE Login = ?
`

func GetUserByNameCommand(userName string) (user common.User, err error) {
	row := db_access.QueryRow(GetUserByNameQuery, userName)

	err = row.Scan(&user.Id, &user.Login, &user.Nickname, &user.Email)

	return
}

const GetUserByEmailQuery = `
	SELECT Id, Login, Nickname, Email
	FROM Users
	WHERE Email = ?
`

func GetUserByEmailCommand(userEmail string) (user common.User, err error) {
	row := db_access.QueryRow(GetUserByEmailQuery, userEmail)

	err = row.Scan(&user.Id, &user.Login, &user.Nickname, &user.Email)

	return
}

const GetUserByNameAndPasswordQuery = `
	SELECT Id, Login, Nickname, Email
	FROM Users
	WHERE Login = ?
		AND Password = ?
`

func GetUserByNameAndPasswordCommand(login string, password string) (user common.User, err error) {
	row := db_access.QueryRow(GetUserByNameAndPasswordQuery, login, password)

	err = row.Scan(&user.Id, &user.Login, &user.Nickname, &user.Email)

	return
}

const CreateUserQuery = `
	INSERT INTO Users (Login, Email, Password)
	VALUES (?, ?, ?)
`

func CreateUserCommand(login string, email string, password string) error {
	_, err := db_access.Exec(CreateUserQuery, login, email, password)

	return err
}

const GetUserSessionByIdQuery = `
	SELECT Id, UserId
	FROM UserSessions
	WHERE Id = ?
`

func GetUserSessionByIdCommand(sessionId string) (userSession common.UserSession, err error) {
	row := db_access.QueryRow(GetUserSessionByIdQuery, sessionId)

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

func DeleteUserSessionCommand(sessionId string) error {
	_, err := db_access.Exec(DeleteUserSessionQuery, sessionId)

	return err
}
