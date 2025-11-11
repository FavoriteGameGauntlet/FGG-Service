package game_service

import (
	"FGG-Service/db_access"
	"FGG-Service/timer_service"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
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
       	END AS 'DoesExist'`
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
       	END AS 'DoesExist'`
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
       	END AS 'DoesExist'`
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
		SELECT g.Id, g.Name, g.Link
		FROM UnplayedGames ug
			INNER JOIN Games g ON ug.GameId = g.Id
		WHERE ug.UserId = $userId`
	GetCurrentGameCommand = `
		SELECT 
		  g.Id,
		  g.Name,
		  gh.State,
		  g.Link,
		  SUM(t.DurationInS) AS SpentSeconds,
		  gh.FinishDate
		FROM GameHistory gh
			INNER JOIN Games g ON gh.GameId = g.Id
			LEFT JOIN Timers t ON t.UserId = gh.UserId AND t.GameId = gh.GameId
		WHERE gh.UserId = $userId
			AND gh.State <> $finishedGameState
		GROUP BY g.Id, g.Name, gh.State, g.Link, gh.FinishDate`
	CancelCurrentGameCommand = `
		UPDATE GameHistory
		SET State = $cancelledGameState,
			FinishDate = datetime('now', 'subsec'),
			ResultPoints = $resultPoints
		WHERE UserId = $userId
			AND GameId = $gameId;`
	FinishCurrentGameCommand = `
		UPDATE GameHistory
		SET State = $finishedGameState,
			FinishDate = datetime('now', 'subsec'),
			ResultPoints = $resultPoints
		WHERE UserId = $userId
			AND GameId = $gameId;`
	CreateEffectRollCommand = `
		INSERT INTO EffectHistory (Id, UserId, GameId)
		VALUES ($effectHistoryId, $gameId, $userId)`
	GetGameHistoryCommand = `
		SELECT 
		  g.Id,
		  g.Name,
		  gh.State,
		  g.Link,
		  SUM(t.DurationInS) AS SpentSeconds,
		  gh.FinishDate
		FROM GameHistory gh
			INNER JOIN Games g ON gh.GameId = g.Id
			LEFT JOIN Timers t ON t.UserId = gh.UserId AND t.GameId = gh.GameId
		WHERE gh.UserId = $userId
		GROUP BY g.Id, g.Name, gh.State, g.Link, gh.FinishDate
		ORDER BY gh.FinishDate NULLS FIRST`
	CreateCurrentGameCommand = `
		INSERT INTO GameHistory (Id, UserId, GameId)
		VALUES ($gameHistoryId, $userId, $gameId)`
	DeleteUnplayedGameCommand = `
		DELETE FROM UnplayedGames
		WHERE UserId = $userId
			AND GameId = $gameId`
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
	doesExist, err := CheckIfUnplayedGameExists(
		userId,
		unplayedGame.Name,
	)

	if err != nil || doesExist {
		return err
	}

	game, err := AddOrGetGame(unplayedGame)

	if err != nil {
		return err
	}

	unplayedGameId := uuid.New().String()
	_, err = db_access.Exec(
		AddUnplayedGameCommand,
		unplayedGameId,
		userId,
		game.Id,
	)

	return err
}

func CheckIfUnplayedGameExists(userId uuid.UUID, gameName string) (bool, error) {
	row := db_access.QueryRow(CheckIfUnplayedGameExistsCommand, userId, gameName)

	var doesExist bool
	err := row.Scan(&doesExist)

	if err != nil {
		return doesExist, err
	}

	return doesExist, nil
}

func AddOrGetGame(unplayedGame *UnplayedGame) (*Game, error) {
	doesExist, err := CheckIfGameExists(unplayedGame.Name)

	if err != nil {
		return nil, err
	}

	if doesExist {
		row := db_access.QueryRow(GetGameCommand, unplayedGame.Name)

		game := Game{}
		err = row.Scan(&game.Id, &game.Name, &game.Link)

		if err != nil {
			return nil, err
		}

		return &game, nil
	}

	gameId := uuid.New()
	_, err = db_access.Exec(
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

func CheckIfGameExists(gameName string) (bool, error) {
	row := db_access.QueryRow(CheckIfGameExistsCommand, gameName)

	var doesExist bool
	err := row.Scan(&doesExist)

	if err != nil {
		return doesExist, err
	}

	return doesExist, nil
}

func GetUnplayedGames(userId uuid.UUID) (*UnplayedGames, error) {
	rows, err := db_access.Query(GetUnplayedGamesCommand, userId)

	if err != nil {
		return nil, err
	}

	gameCount := 0
	errorCount := 0
	games := UnplayedGames{}
	for rows.Next() {
		gameCount++

		game := UnplayedGame{}
		err = rows.Scan(&game.GameId, &game.Name, &game.Link)

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
	row := db_access.QueryRow(GetCurrentGameCommand, userId, GameStateFinished)

	game := Game{}
	var spentSeconds *int
	err := row.Scan(&game.Id, &game.Name, &game.State, &game.Link, &spentSeconds, &game.FinishDate)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if spentSeconds != nil {
		hourCount := *spentSeconds / int(time.Hour/time.Second)

		game.HourCount = &hourCount
	}

	return &game, nil
}

func CheckIfCurrentGameExists(userId uuid.UUID) (bool, error) {
	row := db_access.QueryRow(CheckIfCurrentGameExistsCommand, userId, GameStateFinished)

	var doesExist bool
	err := row.Scan(&doesExist)

	if err != nil {
		return doesExist, err
	}

	return doesExist, nil
}

func CancelCurrentGame(userId uuid.UUID) error {
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

	diceCount := CancellingGamePenaltyDiceCount
	resultPoints, err := RollResultPoints(diceCount)

	if err != nil {
		return err
	}

	_, err = db_access.Exec(CancelCurrentGameCommand, GameStateCancelled, resultPoints, userId, game.Id)

	return err
}

func FinishCurrentGame(userId uuid.UUID) (bool, error) {
	_, err := timer_service.StopCurrentTimer(userId)

	if err != nil {
		return false, err
	}

	game, err := GetCurrentGame(userId)

	if err != nil {
		return false, err
	}

	if game == nil {
		return false, err
	}

	if game.HourCount == nil {
		return false, nil
	}

	additionalDiceCount := *game.HourCount / HourCountForDice
	diceCount := 1 + additionalDiceCount
	resultPoints, err := RollResultPoints(diceCount)

	if err != nil {
		return false, err
	}

	_, err = db_access.Exec(FinishCurrentGameCommand, GameStateFinished, resultPoints, userId, game.Id)

	if err != nil {
		return false, err
	}

	effectHistoryId := uuid.New()
	_, err = db_access.Exec(CreateEffectRollCommand, effectHistoryId.String(), userId, game.Id)

	return true, err
}

func RollResultPoints(diceCount int) (int, error) {
	diceFormula := fmt.Sprintf("%dd6", diceCount)

	roll, _, err := dice.Roll(diceFormula)
	var rolledValue int

	if err != nil {
		return rolledValue, err
	}

	rolledValue = roll.Int()

	return rolledValue, nil
}

func GetGameHistory(userId uuid.UUID) (*Games, error) {
	rows, err := db_access.Query(GetGameHistoryCommand, userId)

	if err != nil {
		return nil, err
	}

	gameCount := 0
	errorCount := 0
	games := Games{}
	for rows.Next() {
		gameCount++

		game := Game{}
		var spentSeconds *int
		var finishDateString *string
		err = rows.Scan(&game.Id, &game.Name, &game.State, &game.Link, &spentSeconds, &finishDateString)

		if err != nil {
			errorCount++
			continue
		}

		if spentSeconds != nil {
			hourCount := *spentSeconds / int(time.Hour/time.Second)

			game.HourCount = &hourCount
		}

		var finishDate *time.Time

		if finishDateString != nil {
			var notNilFinishDate time.Time
			notNilFinishDate, err = time.Parse(db_access.ISO8601, *finishDateString)

			if err != nil {
				errorCount++
				continue
			}

			finishDate = &notNilFinishDate
		}

		game.FinishDate = finishDate

		games = append(games, game)
	}

	if errorCount > 0 && errorCount == gameCount {
		return nil, err
	}

	return &games, nil
}

func MakeGameRoll(userId uuid.UUID) (*Game, error) {
	unplayedGames, err := GetUnplayedGames(userId)

	if err != nil {
		return nil, err
	}

	if unplayedGames == nil || len(*unplayedGames) == 0 {
		return nil, err
	}

	randomNumber := rand.Intn(len(*unplayedGames))
	randomUnplayedGame := (*unplayedGames)[randomNumber]

	gameHistoryId := uuid.New()
	_, err = db_access.Exec(CreateCurrentGameCommand, gameHistoryId.String(), userId, randomUnplayedGame.GameId)

	if err != nil {
		return nil, err
	}

	_, err = db_access.Exec(DeleteUnplayedGameCommand, userId, randomUnplayedGame.GameId)

	if err != nil {
		return nil, err
	}

	return &Game{
		Id:    randomUnplayedGame.GameId,
		Name:  randomUnplayedGame.Name,
		Link:  randomUnplayedGame.Link,
		State: GameStateStarted,
	}, nil
}
