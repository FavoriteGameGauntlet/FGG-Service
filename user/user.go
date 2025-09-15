package user

import (
	"FavoriteGameGauntlet/api"
	"FavoriteGameGauntlet/database"
)

const (
	FindUserCommand = `
		SELECT Id
		FROM Users
		WHERE Name = $userName`
)

func FindUser(username string) (*api.User, error) {
	row := database.QueryRow(FindUserCommand, username)

	var userId int64
	err := row.Scan(&userId)

	if err != nil {
		return nil, err
	}

	return &api.User{Id: userId, Name: username}, nil
}
