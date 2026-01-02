package auth_service

import (
	"FGG-Service/common"
	"FGG-Service/db_access"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

const (
	CheckIfUserNameExistsCommand = `
		SELECT CASE 
        	WHEN EXISTS (
				SELECT 1
				FROM Users
				WHERE Name = ?) 
         	THEN true
         	ELSE false
       	END AS 'DoesExist'`
	CheckIfUserEmailExistsCommand = `
		SELECT CASE 
        	WHEN EXISTS (
				SELECT 1
				FROM Users
				WHERE Email = ?) 
         	THEN true
         	ELSE false
       	END AS 'DoesExist'`
	CheckIfUserSessionExistsCommand = `
		SELECT CASE 
        	WHEN EXISTS (
				SELECT 1
				FROM UserSessions
				WHERE Id = ?) 
         	THEN true
         	ELSE false
       	END AS 'DoesExist'`
	CreateUserCommand = `
		INSERT INTO Users (Name, Email, Password)
		VALUES (?, ?, ?)`
	GetUserIdBySessionIdCommand = `
		SELECT UserId
		FROM UserSessions
		WHERE Id = ?`
	GetUserIdByUserNameAndPasswordCommand = `
		SELECT Id
		FROM Users
		WHERE Name = ?
			AND PASSWORD = ?`
	DeleteSessionCommand = `
		DELETE FROM main.UserSessions
		WHERE Id = ?`
	CreateSessionCommand = `
		INSERT INTO UserSessions (Id, UserId)
		VALUES (?, ?)`
)

func CheckIfUserSessionExists(sessionId string) (bool, error) {
	row := db_access.QueryRow(CheckIfUserSessionExistsCommand, sessionId)

	var doesExist bool
	err := row.Scan(&doesExist)

	if err != nil {
		return doesExist, err
	}

	return doesExist, nil
}

func CreateUser(userName string, userEmail string, userPassword string) error {
	err := ThrowIfUserNameExists(userName)

	if err != nil {
		return err
	}

	err = CheckIfUserEmailExists(userEmail)

	if err != nil {
		return err
	}

	_, err = db_access.Exec(CreateUserCommand, userName, userEmail, userPassword)

	return err
}

func ThrowIfUserNameExists(userName string) error {
	row := db_access.QueryRow(CheckIfUserNameExistsCommand, userName)

	var doesExist bool
	err := row.Scan(&doesExist)

	if err != nil {
		return err
	}

	if doesExist {
		return common.NewUserNameAlreadyExistsError()
	}

	return nil
}

func CheckIfUserEmailExists(userEmail string) error {
	row := db_access.QueryRow(CheckIfUserEmailExistsCommand, userEmail)

	var doesExist bool
	err := row.Scan(&doesExist)

	if err != nil {
		return err
	}

	if doesExist {
		return common.NewUserEmailAlreadyExistsError()
	}

	return nil
}

func GetUserId(sessionId string) (int, error) {
	row := db_access.QueryRow(GetUserIdBySessionIdCommand, sessionId)

	var userId int
	err := row.Scan(&userId)

	if errors.Is(err, sql.ErrNoRows) {
		return userId, common.NewActiveSessionNotFoundAuthError()
	}

	if err != nil {
		return userId, err
	}

	return userId, nil
}

func DeleteSession(sessionId string) error {
	_, err := db_access.Exec(DeleteSessionCommand, sessionId)

	return err
}

func CreateSession(userName string, userPassword string) (*string, error) {
	row := db_access.QueryRow(GetUserIdByUserNameAndPasswordCommand, userName, userPassword)

	var userId int
	err := row.Scan(&userId)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, common.NewWrongDataAuthError()
	}

	if err != nil {
		return nil, err
	}

	sessionId := uuid.New().String()
	_, err = db_access.Exec(CreateSessionCommand, sessionId, userId)

	return &sessionId, err
}
