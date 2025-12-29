package game_service

import (
	"FGG-Service/common"
	"FGG-Service/db_access"
)

const CreateGameQuery = `
	INSERT INTO Games (Name, Link)
	VALUES ($gameName, $gameLink)`

func CreateGameCommand(name string, link *string) error {
	_, err := db_access.Exec(
		CreateGameQuery,
		name,
		link,
	)

	return err
}

const GetGameQuery = `
	SELECT Id, Name, Link
	FROM Games
	WHERE Name = $gameName`

func GetGameCommand(name string) (common.Game, error) {
	row := db_access.QueryRow(GetGameQuery, name)

	game := common.Game{}
	err := row.Scan(&game.Id, &game.Name, &game.Link)

	if err != nil {
		return game, err
	}

	return game, nil
}

const CreateUnplayedGameQuery = `
	INSERT INTO UnplayedGames (UserId, GameId)
	VALUES ($userId, $gameId)`

func CreateUnplayedGameCommand(userId int, gameId int) error {
	_, err := db_access.Exec(
		CreateUnplayedGameQuery,
		userId,
		gameId,
	)

	return err
}

const DeleteUnplayedGameQuery = `
	DELETE FROM UnplayedGames
	WHERE UserId = $userId
		AND GameId = $gameId`

func DeleteUnplayedGameCommand(userId int, gameId int) error {
	_, err := db_access.Exec(DeleteUnplayedGameQuery, userId, gameId)

	return err
}

const GetUnplayedGamesQuery = `
	SELECT ug.Id, g.Id, g.Name, g.Link
	FROM UnplayedGames ug
		INNER JOIN Games g ON ug.GameId = g.Id
	WHERE ug.UserId = $userId`

func GetUnplayedGamesCommand(userId int) (common.UnplayedGames, error) {
	rows, err := db_access.Query(GetUnplayedGamesQuery, userId)

	if err != nil {
		return common.UnplayedGames{}, err
	}

	games := common.UnplayedGames{}

	for rows.Next() {
		game := common.UnplayedGame{}
		err = rows.Scan(&game.Id, &game.GameId, &game.Name, &game.Link)

		if err != nil {
			return common.UnplayedGames{}, err
		}

		games = append(games, game)
	}

	return games, nil
}

const CreateCurrentGameQuery = `
	INSERT INTO GameHistory (UserId, GameId)
	VALUES ($userId, $gameId)`

func CreateCurrentGameCommand(userId int, gameId int) error {
	_, err := db_access.Exec(CreateCurrentGameQuery, userId, gameId)

	return err
}

const GetCurrentGameCommand = `
	SELECT
	  g.Id,
	  g.Name,
	  gh.State,
	  g.Link,
	  gh.FinishDate,
	  SUM(t.DurationInS) AS SpentSeconds
	FROM GameHistory gh
		INNER JOIN Games g ON gh.GameId = g.Id
		LEFT JOIN Timers t ON t.UserId = gh.UserId AND t.GameId = gh.GameId
	WHERE gh.UserId = $userId
		AND gh.State NOT IN ($finishedGameState, $cancelledGameState)
	GROUP BY g.Id, g.Name, gh.State, g.Link, gh.FinishDate`

const GetGameSpentSeconds1Command = `
	SELECT SUM(t.DurationInS - ta.RemainingTimeInS) AS SpentSeconds
	FROM GameHistory gh
		INNER JOIN Timers t ON t.UserId = gh.UserId AND t.GameId = gh.GameId
		INNER JOIN (
    		SELECT ta.TimerId, ta.RemainingTimeInS
			FROM TimerActions ta
				INNER JOIN (
					SELECT TimerId, MAX(CreateDate) AS MaxCreateDate
					FROM TimerActions
					GROUP BY TimerId
				) t ON t.TimerId = ta.TimerId AND t.MaxCreateDate = ta.CreateDate
		) ta ON ta.TimerId = t.Id`
