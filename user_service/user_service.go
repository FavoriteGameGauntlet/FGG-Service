package user_service

import (
	"FGG-Service/database"

	"github.com/google/uuid"
)

const (
	CheckIfUserExistsByNameCommand = `
		SELECT CASE 
        	WHEN EXISTS (
				SELECT 1
				FROM Users
				WHERE Name = $userName) 
         	THEN true
         	ELSE false
       	END AS 'DoesUserExist'`
	CheckIfUserExistsByIdCommand = `
		SELECT 
			CASE WHEN EXISTS (
				SELECT 1
				FROM Users
				WHERE Id = $userId) 
         	THEN true
         	ELSE false
       	END AS 'DoesUserExist'`
	FindUserCommand = `
		SELECT Id
		FROM Users
		WHERE Name = $userName`
	CreateUserCommand = `
		INSERT INTO Users (Id, Name)
		VALUES ($userId, $userName)`
)

func CheckIfUserExistsByName(userName string) (bool, error) {
	row := database.QueryRow(CheckIfUserExistsByNameCommand, userName)

	var doesUserExist bool
	err := row.Scan(&doesUserExist)

	if err != nil {
		return doesUserExist, err
	}

	return doesUserExist, nil
}

func CheckIfUserExistsById(userId uuid.UUID) (bool, error) {
	row := database.QueryRow(CheckIfUserExistsByIdCommand, userId)

	var doesUserExist bool
	err := row.Scan(&doesUserExist)

	if err != nil {
		return doesUserExist, err
	}

	return doesUserExist, nil
}

func FindUser(userName string) (*User, error) {
	row := database.QueryRow(FindUserCommand, userName)

	var userId uuid.UUID
	err := row.Scan(&userId)

	if err != nil {
		return nil, err
	}

	return &User{Id: userId, Name: userName}, err
}

func CreateUser(userName string) error {
	userId := uuid.New().String()
	_, err := database.Exec(CreateUserCommand, userId, userName)

	return err
}
