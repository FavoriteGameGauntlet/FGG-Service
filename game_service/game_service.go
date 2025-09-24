package game_service

import (
	"FGG-Service/api"
	"FGG-Service/database"

	"github.com/google/uuid"
)

const (
	CheckIfGameExistsCommand = `
		SELECT CASE 
        	WHEN EXISTS (
				SELECT 1
				FROM Games
				WHERE Name = $gameName)
         	THEN true
         	ELSE false
       	END AS 'DoesGameExist'`
	CheckIfUnplayedGameExistsCommand = `
		SELECT CASE 
        	WHEN EXISTS (
				SELECT 1
				FROM UnplayedGames ug
					INNER JOIN Games g ON ug.GameId = g.Id
				WHERE ug.UserId = $userId 
					AND g.Name = $gameName)
         	THEN true
         	ELSE false
       	END AS 'DoesGameExist'`
	AddGameCommand = `
		INSERT INTO Games (Id, Name, Link)
		VALUES ($gameId, $gameName, $gameLink)`
	GetGameCommand = `
		SELECT Id
		FROM Games
		WHERE Name = $gameName`
	AddUnplayedGameCommand = `
		INSERT INTO UnplayedGames (Id, UserId, GameId)
		VALUES ($unplayedGameId, $userId, $gameId)`
	GetUnplayedGamesCommand = `
		SELECT g.Name, g.Link
		FROM UnplayedGames ug
			INNER JOIN Games g ON ug.GameId = g.Id
		WHERE UserId = $userId`
)

func AddUnplayedGames(userId uuid.UUID, gamesPtr *api.Games) error {
	games := *gamesPtr
	numberOfGames := len(games)
	errorCount := 0

	var err error
	for _, game := range games {
		err = AddUnplayedGame(userId, &game)

		if err != nil {
			errorCount++
		}
	}

	if errorCount == numberOfGames {
		return err
	}

	return nil
}

func AddUnplayedGame(userId uuid.UUID, game *api.Game) error {
	doesUnplayedGameExist, err := CheckIfUnplayedGameExists(userId, game.Name)

	if err != nil || *doesUnplayedGameExist {
		return err
	}

	gameId, err := AddOrGetGame(game)

	if err != nil {
		return err
	}

	unplayedGameId := uuid.New().String()
	_, err = database.Exec(
		AddUnplayedGameCommand,
		unplayedGameId,
		userId,
		gameId)

	return err
}

func CheckIfUnplayedGameExists(userId uuid.UUID, gameName string) (*bool, error) {
	row := database.QueryRow(CheckIfUnplayedGameExistsCommand, userId, gameName)

	var doesUnplayedGameExist bool
	err := row.Scan(&doesUnplayedGameExist)

	if err != nil {
		return nil, err
	}

	return &doesUnplayedGameExist, err
}

func AddOrGetGame(game *api.Game) (*uuid.UUID, error) {
	doesGameExist, err := CheckIfGameExists(game.Name)

	if err != nil {
		return nil, err
	}

	if *doesGameExist {
		row := database.QueryRow(GetGameCommand, game.Name)

		var gameId uuid.UUID
		err = row.Scan(&gameId)

		if err != nil {
			return nil, err
		}

		return &gameId, nil
	}

	gameId := uuid.New()
	_, err = database.Exec(AddGameCommand, gameId.String(), game.Name, game.Link)

	if err != nil {
		return nil, err
	}

	return &gameId, nil
}

func CheckIfGameExists(gameName string) (*bool, error) {
	row := database.QueryRow(CheckIfGameExistsCommand, gameName)

	var doesGameExist bool
	err := row.Scan(&doesGameExist)

	if err != nil {
		return nil, err
	}

	return &doesGameExist, err
}

func GetUnplayedGames(userId uuid.UUID) (*api.Games, error) {
	rows, err := database.Query(GetUnplayedGamesCommand, userId)

	if err != nil {
		return nil, err
	}

	games := api.Games{}
	for rows.Next() {
		g := api.Game{}
		err := rows.Scan(&g.Name, &g.Link)

		if err != nil {
			continue
		}

		games = append(games, g)
	}

	return &games, nil
}
