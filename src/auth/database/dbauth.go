package dbauth

import (
	"FGG-Service/src/auth/types"
	"FGG-Service/src/dbaccess"

	"github.com/google/uuid"
)

type Database struct {
}

const GetUserByLoginQuery = `
	SELECT Id, Login, DisplayName, Email
	FROM Users
	WHERE Login = ?
`

func (db *Database) GetUserByLoginCommand(userLogin string) (user typeauth.User, err error) {
	row := dbaccess.QueryRow("GetUserByLoginQuery", GetUserByLoginQuery, userLogin)

	err = row.Scan(&user.Id, &user.Login, &user.DisplayName, &user.Email)

	return
}

const GetUserByEmailQuery = `
	SELECT Id, Login, DisplayName, Email
	FROM Users
	WHERE Email = ?
`

func (db *Database) GetUserByEmailCommand(userEmail string) (user typeauth.User, err error) {
	row := dbaccess.QueryRow("GetUserByEmailQuery", GetUserByEmailQuery, userEmail)

	err = row.Scan(&user.Id, &user.Login, &user.DisplayName, &user.Email)

	return
}

const GetUserByLoginAndPasswordQuery = `
	SELECT Id, Login, DisplayName, Email
	FROM Users
	WHERE Login = ?
		AND Password = ?
`

func (db *Database) GetUserByLoginAndPasswordCommand(login string, password string) (user typeauth.User, err error) {
	row := dbaccess.QueryRow("GetUserByLoginAndPasswordQuery", GetUserByLoginAndPasswordQuery, login, password)

	err = row.Scan(&user.Id, &user.Login, &user.DisplayName, &user.Email)

	return
}

const CreateUserQuery = `
	INSERT INTO Users (Login, Email, Password)
	VALUES (?, ?, ?)
`

func (db *Database) CreateUserCommand(login string, email string, password string) error {
	_, err := dbaccess.Exec("CreateUserQuery", CreateUserQuery, login, email, password)

	return err
}

const CreateUserStatsQuery = `
	INSERT INTO UserStats (UserId)
	SELECT Id FROM Users WHERE Login = ?
`

func (db *Database) CreateUserStatsCommand(login string) error {
	_, err := dbaccess.Exec("CreateUserStatsQuery", CreateUserStatsQuery, login)

	return err
}

const GetUserSessionByIdQuery = `
	SELECT Id, UserId
	FROM UserSessions
	WHERE Id = ?
`

func (db *Database) GetUserSessionByIdCommand(sessionId string) (userSession typeauth.UserSession, err error) {
	row := dbaccess.QueryRow("GetUserSessionByIdQuery", GetUserSessionByIdQuery, sessionId)

	err = row.Scan(&userSession.Id, &userSession.UserId)

	return
}

const CreateUserSessionQuery = `
	INSERT INTO UserSessions (Id, UserId)
	VALUES (?, ?)
`

func (db *Database) CreateUserSessionCommand(userId int) (userSession typeauth.UserSession, err error) {
	sessionId := uuid.New().String()

	_, err = dbaccess.Exec("CreateUserSessionQuery", CreateUserSessionQuery, sessionId, userId)

	if err != nil {
		return
	}

	userSession = typeauth.UserSession{
		Id:     sessionId,
		UserId: userId,
	}

	return
}

const DeleteUserSessionQuery = `
	DELETE FROM UserSessions
	WHERE Id = ?
`

func (db *Database) DeleteUserSessionCommand(sessionId string) error {
	_, err := dbaccess.Exec("DeleteUserSessionQuery", DeleteUserSessionQuery, sessionId)

	return err
}
