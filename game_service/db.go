package game_service

import (
	"FGG-Service/common"
	"FGG-Service/db_access"
	"database/sql"
	"errors"
)

const CreateGameQuery = `
	INSERT INTO Games (Name, Link)
	VALUES ($gameName, $gameLink)
`

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
	WHERE Name = $gameName
`

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
	VALUES ($userId, $gameId)
`

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
		AND GameId = $gameId
`

func DeleteUnplayedGameCommand(userId int, gameId int) error {
	_, err := db_access.Exec(DeleteUnplayedGameQuery, userId, gameId)

	return err
}

const GetUnplayedGamesQuery = `
	SELECT ug.Id, g.Id, g.Name, g.Link
	FROM UnplayedGames ug
		INNER JOIN Games g ON ug.GameId = g.Id
	WHERE ug.UserId = $userId
`

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
	VALUES ($userId, $gameId)
`

func CreateCurrentGameCommand(userId int, gameId int) error {
	_, err := db_access.Exec(CreateCurrentGameQuery, userId, gameId)

	return err
}

const GetCurrentGameQuery = `
	SELECT
		g.Id,
		g.Name,
		gh.State,
		g.Link,
		gh.FinishDate
	FROM GameHistory gh
		INNER JOIN Games g ON gh.GameId = g.Id
		LEFT JOIN Timers t ON t.UserId = gh.UserId AND t.GameId = gh.GameId
	WHERE gh.UserId = $userId
		AND gh.State NOT IN ($finishedGameState, $cancelledGameState)
`

func GetCurrentGameCommand(userId int) (common.Game, error) {
	row := db_access.QueryRow(GetCurrentGameQuery, userId, common.GameStateFinished, common.GameStateCancelled)

	game := common.Game{}
	var finishDateString *string
	err := row.Scan(&game.Id, &game.Name, &game.State, &game.Link, &finishDateString)

	if errors.Is(err, sql.ErrNoRows) {
		return common.Game{}, common.NewCurrentGameNotFoundError()
	}

	if err != nil {
		return common.Game{}, err
	}

	finishDate, err := common.ConvertToNullableDate(finishDateString)

	if err != nil {
		return common.Game{}, err
	}

	game.FinishDate = finishDate

	return game, nil
}

const GetGameSecondsSpentQuery = `
SELECT SUM(
        t.DurationInS -
        CASE ta.Action
            WHEN 'start' THEN ta.RemainingTimeInS - (strftime('%s','now') - strftime('%s', ta.CreateDate))
            WHEN 'pause' THEN ta.RemainingTimeInS
            WHEN 'stop'  THEN ta.RemainingTimeInS
            ELSE t.DurationInS
        END
    ) AS SecondsSpent
FROM Timers t
	LEFT JOIN TimerActions ta ON ta.Id = (
        SELECT ta2.Id
        FROM TimerActions ta2
        WHERE ta2.TimerId = t.Id
        ORDER BY ta2.CreateDate DESC
        LIMIT 1
    )
WHERE t.UserId = $userId
  	AND t.GameId = $gameId
`

func GetGameSecondsSpentCommand(userId int, gameId int) (int, error) {
	row := db_access.QueryRow(GetGameSecondsSpentQuery, userId, gameId)

	var secondsSpent int
	err := row.Scan(&secondsSpent)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return secondsSpent, nil
}

const CancelCurrentGameQuery = `
	UPDATE GameHistory
	SET State = $cancelledGameState,
		FinishDate = datetime('now', 'subsec')
	WHERE UserId = $userId
		AND GameId = $gameId;
`

func CancelCurrentGameCommand(userId int, gameId int) error {
	_, err := db_access.Exec(CancelCurrentGameQuery, common.GameStateCancelled, userId, gameId)

	return err
}

const FinishCurrentGameQuery = `
	UPDATE GameHistory
	SET State = $finishedGameState,
		FinishDate = datetime('now', 'subsec')
	WHERE UserId = $userId
		AND GameId = $gameId;
`

func FinishCurrentGameCommand(userId int, gameId int) error {
	_, err := db_access.Exec(FinishCurrentGameQuery, common.GameStateFinished, userId, gameId)

	return err
}

const GetGameHistoryCommand = `
	SELECT 
		g.Id,
		g.Name,
		gh.State,
		g.Link,
		gh.FinishDate
	FROM GameHistory gh
		INNER JOIN Games g ON gh.GameId = g.Id
		LEFT JOIN Timers t ON t.UserId = gh.UserId AND t.GameId = gh.GameId
	WHERE gh.UserId = $userId
	ORDER BY gh.FinishDate NULLS FIRST
`
