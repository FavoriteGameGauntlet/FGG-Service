package game_service

import (
	"FGG-Service/database"
	"FGG-Service/timer_service"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/justinian/dice"
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
					AND gh.State <> $finishedGameState)
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
		SELECT ug.Id, g.Name, g.Link
		FROM UnplayedGames ug
			INNER JOIN Games g ON ug.GameId = g.Id
		WHERE ug.UserId = $userId`
	GetCurrentGameCommand = `
		SELECT g.Id, g.Name, gh.State, g.Link, gh.FinishDate
		FROM GameHistory gh
			INNER JOIN Games g ON gh.GameId = g.Id
		WHERE gh.UserId = $userId
			AND gh.State <> $finishedGameState`
	GetFinishedTimerCountCommand = `
		SELECT COUNT(*) AS FinishedTimerCount
		FROM Timers
		WHERE UserId = $userId 
			AND GameId = $gameId
			AND State = $finishedTimerState`
	FinishCurrentGameCommand = `
		UPDATE GameHistory
		SET State = $finishedGameState,
			FinishDate = datetime('now'),
			ResultPoints = $resultPoints
		WHERE UserId = $userId
			AND GameId = $gameId;`
	CreateEffectRollCommand = `
		INSERT INTO EffectHistory (Id, UserId, GameId)
		VALUES ($effectHistoryId, $gameId, $userId)`
	GetFinishedGamesCommand = `
		SELECT g.Id, g.Name, gh.State, g.Link, gh.FinishDate
		FROM GameHistory gh
			INNER JOIN Games g ON gh.GameId = g.Id
		WHERE gh.UserId = $userId
			AND gh.State = $finishedGameState
		ORDER BY gh.FinishDate`
)

func AddUnplayedGames(userId uuid.UUID, gamesPtr *UnplayedGames) error {
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

func AddUnplayedGame(userId uuid.UUID, unplayedGame *UnplayedGame) error {
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

func AddOrGetGame(unplayedGame *UnplayedGame) (*Game, error) {
	doesGameExist, err := CheckIfGameExists(unplayedGame.Name)

	if err != nil {
		return nil, err
	}

	if *doesGameExist {
		row := database.QueryRow(GetGameCommand, unplayedGame.Name)

		game := Game{}
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

	return &Game{
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

func GetUnplayedGames(userId uuid.UUID) (*UnplayedGames, error) {
	rows, err := database.Query(GetUnplayedGamesCommand, userId)

	if err != nil {
		return nil, err
	}

	gameCount := 0
	errorCount := 0
	games := UnplayedGames{}
	for rows.Next() {
		gameCount++

		game := UnplayedGame{}
		err = rows.Scan(&game.Id, &game.Name, &game.Link)

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

func GetCurrentGame(userId uuid.UUID) (*Game, error) {
	row := database.QueryRow(GetCurrentGameCommand, userId, GameStateFinished)

	game := Game{}
	err := row.Scan(&game.Id, &game.Name, &game.State, &game.Link, &game.FinishDate)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &game, nil
}

func CheckIfCurrentGameExists(userId uuid.UUID) (*bool, error) {
	row := database.QueryRow(CheckIfCurrentGameExistsCommand, userId, GameStateFinished)

	var doesCurrentGameExist bool
	err := row.Scan(&doesCurrentGameExist)

	if err != nil {
		return nil, err
	}

	return &doesCurrentGameExist, nil
}

func FinishCurrentGame(userId uuid.UUID) error {
	game, err := GetCurrentGame(userId)

	if err != nil {
		return err
	}

	if game == nil {
		return err
	}

	_, err = timer_service.StopCurrentTimer(userId)

	if err != nil {
		return err
	}

	row := database.QueryRow(GetFinishedTimerCountCommand, userId, game.Id, timer_service.TimerStateFinished)

	var finishedTimerCount int
	err = row.Scan(&finishedTimerCount)

	if err != nil {
		return err
	}

	resultPoints, err := RollFinishedGameResultPoints(finishedTimerCount)

	if err != nil {
		return err
	}

	_, err = database.Exec(FinishCurrentGameCommand, GameStateFinished, resultPoints, userId, game.Id)

	if err != nil {
		return err
	}

	effectHistoryId := uuid.New()
	_, err = database.Exec(CreateEffectRollCommand, effectHistoryId.String(), userId, game.Id)

	return err
}

func RollFinishedGameResultPoints(diceCount int) (int, error) {
	diceFormula := fmt.Sprintf("%dd6", diceCount)

	roll, _, err := dice.Roll(diceFormula)
	var rolledValue int

	if err != nil {
		return rolledValue, err
	}

	rolledValue = roll.Int()

	return rolledValue, nil
}

func GetFinishedGames(userId uuid.UUID) (*Games, error) {
	rows, err := database.Query(GetFinishedGamesCommand, userId, GameStateFinished)

	if err != nil {
		return nil, err
	}

	gameCount := 0
	errorCount := 0
	games := Games{}
	for rows.Next() {
		gameCount++

		game := Game{}
		var finishDateString string
		err = rows.Scan(&game.Id, &game.Name, &game.State, &game.Link, &finishDateString)

		if err != nil {
			errorCount++
			continue
		}

		var finishDate time.Time

		if finishDateString != "" {
			finishDate, err = time.Parse(database.ISO8601, finishDateString)

			if err != nil {
				errorCount++
				continue
			}
		}

		game.FinishDate = &finishDate

		games = append(games, game)
	}

	if errorCount > 0 && errorCount == gameCount {
		return nil, err
	}

	return &games, nil
}
