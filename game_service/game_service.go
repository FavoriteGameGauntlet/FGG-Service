package game_service

import (
	"FGG-Service/api"
	"FGG-Service/database"
	"database/sql"
	"errors"

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
       	END AS 'DoesUnplayedGameExist'`
	CheckIfCurrentGameExistsCommand = `
		SELECT CASE 
        	WHEN EXISTS (
				SELECT 1
				FROM GameHistory gh
					INNER JOIN Games g ON gh.GameId = g.Id
				WHERE gh.UserId = $userId
					AND gh.Result IS NULL)
         	THEN true
         	ELSE false
       	END AS 'DoesCurrentGameExist'`
	AddGameCommand = `
		INSERT INTO Games (Id, Name, Link)
		VALUES ($gameId, $gameName, $gameLink)`
	GetGameCommand = `
		SELECT Id, Name, Link
		FROM Games
		WHERE Name = $gameName`
	AddUnplayedGameCommand = `
		INSERT INTO UnplayedGames (Id, UserId, GameId)
		VALUES ($unplayedGameId, $userId, $gameId)`
	GetUnplayedGamesCommand = `
		SELECT g.Name, g.Link
		FROM UnplayedGames ug
			INNER JOIN Games g ON ug.GameId = g.Id
		WHERE ug.UserId = $userId`
	GetCurrentGameCommand = `
		SELECT g.Id, g.Name, gh.State, g.Link
		FROM GameHistory gh
			INNER JOIN Games g ON gh.GameId = g.Id
		WHERE gh.UserId = $userId
			AND gh.Result IS NULL`
	FinishCurrentGameCommand = `
		UPDATE GameHistory
		SET State = $finishedGameState,
			FinishDate = datetime('now')
		WHERE UserId = $userId
			AND GameId = $gameId;`
)

func FinishCurrentGame(userId uuid.UUID, gameId uuid.UUID, timerId uuid.UUID) error {


	_, err := database.Exec(
		FinishCurrentGameCommand,
		api.GameStateFinished,
		userId,
		gameId,
	)

	return err
}

func AddUnplayedGames(userId uuid.UUID, gamesPtr *api.UnplayedGames) error {
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

func AddUnplayedGame(userId uuid.UUID, unplayedGame *api.UnplayedGame) error {
	doesUnplayedGameExist, err := CheckIfUnplayedGameExists(
		userId,
		unplayedGame.Name,
	)

	if err != nil || *doesUnplayedGameExist {
		return err
	}

	game, err := AddOrGetGame(unplayedGame)

	if err != nil {
		return err
	}

	unplayedGameId := uuid.New().String()
	_, err = database.Exec(
		AddUnplayedGameCommand,
		unplayedGameId,
		userId,
		game.Id,
	)

	return err
}

func CheckIfUnplayedGameExists(userId uuid.UUID, gameName string) (*bool, error) {
	row := database.QueryRow(CheckIfUnplayedGameExistsCommand, userId, gameName)

	var doesUnplayedGameExist bool
	err := row.Scan(&doesUnplayedGameExist)

	if err != nil {
		return nil, err
	}

	return &doesUnplayedGameExist, nil
}

func AddOrGetGame(unplayedGame *api.UnplayedGame) (*api.Game, error) {
	doesGameExist, err := CheckIfGameExists(unplayedGame.Name)

	if err != nil {
		return nil, err
	}

	if *doesGameExist {
		row := database.QueryRow(GetGameCommand, unplayedGame.Name)

		game := api.Game{}
		err = row.Scan(&game.Id, &game.Name, &game.Link)

		if err != nil {
			return nil, err
		}

		return &game, nil
	}

	gameId := uuid.New()
	_, err = database.Exec(
		AddGameCommand,
		gameId.String(),
		unplayedGame.Name,
		unplayedGame.Link,
	)

	if err != nil {
		return nil, err
	}

	return &api.Game{
		Id:   gameId,
		Name: unplayedGame.Name,
		Link: unplayedGame.Link,
	}, nil
}

func CheckIfGameExists(gameName string) (*bool, error) {
	row := database.QueryRow(CheckIfGameExistsCommand, gameName)

	var doesGameExist bool
	err := row.Scan(&doesGameExist)

	if err != nil {
		return nil, err
	}

	return &doesGameExist, nil
}

func GetUnplayedGames(userId uuid.UUID) (*api.UnplayedGames, error) {
	rows, err := database.Query(GetUnplayedGamesCommand, userId)

	if err != nil {
		return nil, err
	}

	gameCount := 0
	errorCount := 0
	games := api.UnplayedGames{}
	for rows.Next() {
		gameCount++

		game := api.UnplayedGame{}
		err := rows.Scan(&game.Name, &game.Link)

		if err != nil {
			errorCount++
			continue
		}

		games = append(games, game)
	}

	if errorCount > 0 && errorCount == gameCount {
		return nil, err
	}

	return &games, nil
}

func GetCurrentGame(userId uuid.UUID) (*api.Game, error) {
	row := database.QueryRow(GetCurrentGameCommand, userId)

	game := api.Game{}
	err := row.Scan(&game.Id, &game.Name, &game.State, &game.Link)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &game, nil
}

func CheckIfCurrentGameExists(userId uuid.UUID) (*bool, error) {
	row := database.QueryRow(CheckIfCurrentGameExistsCommand, userId)

	var doesCurrentGameExist bool
	err := row.Scan(&doesCurrentGameExist)

	if err != nil {
		return nil, err
	}

	return &doesCurrentGameExist, nil
}
