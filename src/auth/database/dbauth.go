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
	queryName := "GetUserByLoginQuery"
	row := dbaccess.QueryRow(queryName, GetUserByLoginQuery, userLogin)

	err = row.Scan(&user.Id, &user.Login, &user.DisplayName, &user.Email)

	dbaccess.LogDbResult(queryName, user, err)

	return
}

const GetUserByIdQuery = `
	SELECT Id, Login, DisplayName, Email
	FROM Users
	WHERE Id = ?
`

func (db *Database) GetUserByIdCommand(userId int) (user typeauth.User, err error) {
	queryName := "GetUserByIdQuery"
	row := dbaccess.QueryRow(queryName, GetUserByIdQuery, userId)

	err = row.Scan(&user.Id, &user.Login, &user.DisplayName, &user.Email)

	dbaccess.LogDbResult(queryName, user, err)

	return
}

const GetUserByEmailQuery = `
	SELECT Id, Login, DisplayName, Email
	FROM Users
	WHERE Email = ?
`

func (db *Database) GetUserByEmailCommand(userEmail string) (user typeauth.User, err error) {
	queryName := "GetUserByEmailQuery"
	row := dbaccess.QueryRow(queryName, GetUserByEmailQuery, userEmail)

	err = row.Scan(&user.Id, &user.Login, &user.DisplayName, &user.Email)

	dbaccess.LogDbResult(queryName, user, err)

	return
}

const GetUserByLoginAndPasswordQuery = `
	SELECT Id, Login, DisplayName, Email
	FROM Users
	WHERE Login = ?
		AND Password = ?
`

func (db *Database) GetUserByLoginAndPasswordCommand(loginUser typeauth.LoginUser) (user typeauth.User, err error) {
	queryName := "GetUserByLoginAndPasswordQuery"
	row := dbaccess.QueryRow(queryName, GetUserByLoginAndPasswordQuery, loginUser.Login, loginUser.Password)

	err = row.Scan(&user.Id, &user.Login, &user.DisplayName, &user.Email)

	dbaccess.LogDbResult(queryName, user, err)

	return
}

const CreateUserQuery = `
	INSERT INTO Users (Login, Email, Password)
	VALUES (?, ?, ?)
`

func (db *Database) CreateUserCommand(signupUser typeauth.SignupUser) error {
	queryName := "CreateUserQuery"
	_, err := dbaccess.Exec(queryName, CreateUserQuery, signupUser.Login, signupUser.Email, signupUser.Password)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const CreateUserStatsQuery = `
	INSERT INTO UserStats (UserId)
	SELECT Id FROM Users WHERE Login = ?
`

func (db *Database) CreateUserStatsCommand(login string) error {
	queryName := "CreateUserStatsQuery"
	_, err := dbaccess.Exec(queryName, CreateUserStatsQuery, login)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const GetUserSessionByIdQuery = `
	SELECT Id, UserId
	FROM UserSessions
	WHERE Id = ?
`

func (db *Database) GetUserSessionByIdCommand(sessionId string) (userSession typeauth.UserSession, err error) {
	queryName := "GetUserSessionByIdQuery"
	row := dbaccess.QueryRow(queryName, GetUserSessionByIdQuery, sessionId)

	err = row.Scan(&userSession.Id, &userSession.UserId)

	dbaccess.LogDbResult(queryName, userSession, err)

	return
}

const CreateUserSessionQuery = `
	INSERT INTO UserSessions (Id, UserId)
	VALUES (?, ?)
`

func (db *Database) CreateUserSessionCommand(userId int) (userSession typeauth.UserSession, err error) {
	queryName := "CreateUserSessionQuery"
	sessionId := uuid.New().String()
	_, err = dbaccess.Exec(queryName, CreateUserSessionQuery, sessionId, userId)

	if err != nil {
		return
	}

	userSession = typeauth.UserSession{
		Id:     sessionId,
		UserId: userId,
	}

	dbaccess.LogDbResult(queryName, userSession, err)

	return
}

const DeleteUserSessionQuery = `
	DELETE FROM UserSessions
	WHERE Id = ?
`

func (db *Database) DeleteUserSessionCommand(sessionId string) error {
	queryName := "DeleteUserSessionQuery"
	_, err := dbaccess.Exec(queryName, DeleteUserSessionQuery, sessionId)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}
