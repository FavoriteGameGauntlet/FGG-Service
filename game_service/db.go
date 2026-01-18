package game_service

import (
	"FGG-Service/common"
	"FGG-Service/db_access"
	"database/sql"
	"errors"
	"time"
)

const DoesGameExistQuery = `
	SELECT
		CASE WHEN EXISTS (
			SELECT 1
			FROM Games
			WHERE Name = ?
		)
		THEN true
		ELSE false
	END AS DoesExist`

func DoesGameExistCommand(gameName string) (doesExist bool, err error) {
	row := db_access.QueryRow(DoesGameExistQuery, gameName)

	err = row.Scan(&doesExist)

	return
}

const CreateGameQuery = `
	INSERT INTO Games (Name)
	VALUES (?)
`

func CreateGameCommand(name string) error {
	_, err := db_access.Exec(
		CreateGameQuery,
		name,
	)

	return err
}

const GetGameQuery = `
	SELECT Id, Name
	FROM Games
	WHERE Name = ?
`

func GetGameCommand(name string) (game common.Game, err error) {
	row := db_access.QueryRow(GetGameQuery, name)

	err = row.Scan(&game.Id, &game.Name)

	return
}

const DoesUnplayedGameExistQuery = `
	SELECT 
	    CASE WHEN EXISTS (
			SELECT 1
			FROM UnplayedGames ug
				INNER JOIN Games g ON ug.GameId = g.Id
			WHERE ug.UserId = ?
				AND g.Name = ?
		)
		THEN true
		ELSE false
	END AS DoesExist`

func DoesUnplayedGameExistCommand(userId int, gameName string) (doesExist bool, err error) {
	row := db_access.QueryRow(DoesUnplayedGameExistQuery, userId, gameName)

	err = row.Scan(&doesExist)

	return
}

const CreateUnplayedGameQuery = `
	INSERT INTO UnplayedGames (UserId, GameId)
	VALUES (?, ?)
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
	WHERE UserId = ?
		AND GameId = ?
`

func DeleteUnplayedGameCommand(userId int, gameId int) error {
	_, err := db_access.Exec(DeleteUnplayedGameQuery, userId, gameId)

	return err
}

const GetUnplayedGamesQuery = `
	SELECT ug.Id, g.Id, g.Name
	FROM UnplayedGames ug
		INNER JOIN Games g ON ug.GameId = g.Id
	WHERE ug.UserId = ?
`

func GetUnplayedGamesCommand(userId int) (games common.UnplayedGames, err error) {
	rows, err := db_access.Query(GetUnplayedGamesQuery, userId)

	if err != nil {
		return
	}

	for rows.Next() {
		game := common.UnplayedGame{}
		err = rows.Scan(&game.Id, &game.GameId, &game.Name)

		if err != nil {
			_ = rows.Close()
			return
		}

		games = append(games, game)
	}

	_ = rows.Close()
	return
}

const CreateCurrentGameQuery = `
	INSERT INTO GameHistory (UserId, GameId)
	VALUES (?, ?)
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
		gh.FinishDate
	FROM GameHistory gh
		INNER JOIN Games g ON gh.GameId = g.Id
	WHERE gh.UserId = ?
		AND gh.State NOT IN (?, ?)
`

func GetCurrentGameCommand(userId int) (game common.Game, err error) {
	games, err := GetHistoryGames(userId, GetCurrentGameQuery)

	if len(games) == 0 {
		err = common.NewCurrentGameNotFoundError()
		return
	}

	game = games[0]

	return
}

func GetHistoryGames(userId int, query string) (games common.Games, err error) {
	rows, err := db_access.Query(query, userId, common.GameStateFinished, common.GameStateCancelled)

	if err != nil {
		return
	}

	for rows.Next() {
		game := common.Game{}
		var finishDateString *string
		err = rows.Scan(&game.Id, &game.Name, &game.State, &finishDateString)

		if errors.Is(err, sql.ErrNoRows) {
			_ = rows.Close()
			// TODO: Change the error to the most suitable
			err = common.NewCurrentGameNotFoundError()
			return
		}

		if err != nil {
			_ = rows.Close()
			return
		}

		var finishDate *time.Time
		finishDate, err = common.ConvertToNullableDate(finishDateString)

		if err != nil {
			_ = rows.Close()
			return
		}

		game.FinishDate = finishDate

		games = append(games, game)
	}

	_ = rows.Close()

	return
}

const GetGameSecondsSpentQuery = `
SELECT 
    SUM(
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
WHERE t.UserId = ?
  	AND t.GameId = ?
`

func GetGameTimeSpentCommand(userId int, gameId int) (timeSpent time.Duration, err error) {
	row := db_access.QueryRow(GetGameSecondsSpentQuery, userId, gameId)

	var secondsSpent int
	err = row.Scan(&secondsSpent)

	if errors.Is(err, sql.ErrNoRows) {
		err = nil
		return
	}

	if err != nil {
		return
	}

	timeSpent = time.Duration(secondsSpent) * time.Second

	return
}

const CancelCurrentGameQuery = `
	UPDATE GameHistory
	SET State = ?,
		FinishDate = datetime('now', 'subsec')
	WHERE UserId = ?
		AND GameId = ?;
`

func CancelCurrentGameCommand(userId int, gameId int) error {
	_, err := db_access.Exec(CancelCurrentGameQuery, common.GameStateCancelled, userId, gameId)

	return err
}

const FinishCurrentGameQuery = `
	UPDATE GameHistory
	SET State = ?,
		FinishDate = datetime('now', 'subsec')
	WHERE UserId = ?
		AND GameId = ?;
`

func FinishCurrentGameCommand(userId int, gameId int) error {
	_, err := db_access.Exec(FinishCurrentGameQuery, common.GameStateFinished, userId, gameId)

	return err
}

const GetGameHistoryQuery = `
	SELECT 
		g.Id,
		g.Name,
		gh.State,
		gh.FinishDate
	FROM GameHistory gh
		INNER JOIN Games g ON gh.GameId = g.Id
		LEFT JOIN Timers t ON t.UserId = gh.UserId AND t.GameId = gh.GameId
	WHERE gh.UserId = ?
		AND gh.State IN (?, ?)
	ORDER BY gh.FinishDate NULLS FIRST
`

func GetEndedGamesCommand(userId int) (games common.Games, err error) {
	games, err = GetHistoryGames(userId, GetGameHistoryQuery)

	if err != nil {
		return
	}

	for _, game := range games {
		var timeSpent time.Duration
		timeSpent, err = GetGameTimeSpentCommand(userId, game.Id)

		if err != nil {
			return
		}

		game.TimeSpent = timeSpent
	}

	return
}
