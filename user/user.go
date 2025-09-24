package user

import (
	"FGG-Service/api"
	"FGG-Service/database"

	"github.com/google/uuid"
)

const (
	CheckIfUserExistsCommand = `
		SELECT CASE 
        	WHEN EXISTS (SELECT 1 FROM Users WHERE Name = $username) 
         	THEN true
         	ELSE false
       	END AS 'DoesUserExist'`
	FindUserCommand = `
		SELECT Id
		FROM Users
		WHERE Name = $username`
	AddUserCommand = `
		INSERT INTO Users (Id, Name)
		VALUES ($userId, $username)`
)

func CheckIfUserExists(username string) (*bool, error) {
	row := database.QueryRow(CheckIfUserExistsCommand, username)

	var doesUserExist bool
	err := row.Scan(&doesUserExist)

	if err != nil {
		return nil, err
	}

	return &doesUserExist, err
}

func FindUser(username string) (*api.User, error) {
	row := database.QueryRow(FindUserCommand, username)

	var userId uuid.UUID
	err := row.Scan(&userId)

	if err != nil {
		return nil, err
	}

	return &api.User{Id: userId, Name: username}, err
}

func AddUser(username string) error {
	userId := uuid.New().String()
	_, err := database.Exec(AddUserCommand, userId, username)

	return err
}
