package auth_service

import (
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
				WHERE Name = $userName) 
         	THEN true
         	ELSE false
       	END AS 'DoesExist'`
	CheckIfUserEmailExistsCommand = `
		SELECT CASE 
        	WHEN EXISTS (
				SELECT 1
				FROM Users
				WHERE Email = $userEmail) 
         	THEN true
         	ELSE false
       	END AS 'DoesExist'`
	CheckIfUserSessionExistsCommand = `
		SELECT CASE 
        	WHEN EXISTS (
				SELECT 1
				FROM UserSessions
				WHERE Id = $sessionId) 
         	THEN true
         	ELSE false
       	END AS 'DoesExist'`
	CreateUserCommand = `
		INSERT INTO Users (Name, Email, Password)
		VALUES ($userName, $userEmail, $userPassword)`
	GetUserIdBySessionIdCommand = `
		SELECT UserId
		FROM UserSessions
		WHERE Id = $sessionId`
	GetUserIdByUserNameAndPasswordCommand = `
		SELECT Id
		FROM Users
		WHERE Name = $userName
			AND PASSWORD = $userPassword`
	DeleteSessionCommand = `
		DELETE FROM main.UserSessions
		WHERE Id = $sessionId`
	CreateSessionCommand = `
		INSERT INTO UserSessions (Id, UserId)
		VALUES ($sessionId, $userId)`
)

func CheckIfUserNameExists(userName string) (bool, error) {
	row := db_access.QueryRow(CheckIfUserNameExistsCommand, userName)

	var doesExist bool
	err := row.Scan(&doesExist)

	if err != nil {
		return doesExist, err
	}

	return doesExist, nil
}

func CheckIfUserEmailExists(userEmail string) (bool, error) {
	row := db_access.QueryRow(CheckIfUserEmailExistsCommand, userEmail)

	var doesExist bool
	err := row.Scan(&doesExist)

	if err != nil {
		return doesExist, err
	}

	return doesExist, nil
}

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
	_, err := db_access.Exec(CreateUserCommand, userName, userEmail, userPassword)

	return err
}

func GetUserId(sessionId string) (int, error) {
	row := db_access.QueryRow(GetUserIdBySessionIdCommand, sessionId)

	var userId int
	err := row.Scan(&userId)

	if err != nil {
		return userId, err
	}

	return userId, nil
}

func DeleteSession(sessionId string) error {
	_, err := db_access.Exec(DeleteSessionCommand, sessionId)

	return err
}

func CreateSession(userName string, userPassword string) (bool, error) {
	row := db_access.QueryRow(GetUserIdByUserNameAndPasswordCommand, userName, userPassword)

	var userId int
	err := row.Scan(&userId)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	sessionId := uuid.New()
	_, err = db_access.Exec(CreateSessionCommand, sessionId.String(), userId)

	return true, err
}
