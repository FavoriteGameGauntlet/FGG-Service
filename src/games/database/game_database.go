package game_database

import (
	"FGG-Service/src/common"
	"FGG-Service/src/db_access"
	"database/sql"
	"errors"
	"time"
)

type Database struct {
}

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

func (db *Database) DoesGameExistCommand(gameName string) (doesExist bool, err error) {
	row := db_access.QueryRow(DoesGameExistQuery, gameName)

	err = row.Scan(&doesExist)

	return
}

const CreateGameQuery = `
	INSERT INTO Games (Name)
	VALUES (?)
`

func (db *Database) CreateGameCommand(name string) error {
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

func (db *Database) GetGameCommand(name string) (game common.Game, err error) {
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

func (db *Database) DoesUnplayedGameExistCommand(userId int, gameName string) (doesExist bool, err error) {
	row := db_access.QueryRow(DoesUnplayedGameExistQuery, userId, gameName)

	err = row.Scan(&doesExist)

	return
}

const CreateUnplayedGameQuery = `
	INSERT INTO UnplayedGames (UserId, GameId)
	VALUES (?, ?)
`

func (db *Database) CreateUnplayedGameCommand(userId int, gameId int) error {
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

func (db *Database) DeleteUnplayedGameCommand(userId int, gameId int) error {
	_, err := db_access.Exec(DeleteUnplayedGameQuery, userId, gameId)

	return err
}

const GetUnplayedGamesQuery = `
	SELECT ug.Id, g.Id, g.Name
	FROM UnplayedGames ug
		INNER JOIN Games g ON ug.GameId = g.Id
	WHERE ug.UserId = ?
`

func (db *Database) GetUnplayedGamesCommand(userId int) (games common.UnplayedGames, err error) {
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

func (db *Database) CreateCurrentGameCommand(userId int, gameId int) error {
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

func (db *Database) GetCurrentGameCommand(userId int) (games common.Games, err error) {
	games, err = db.getHistoryGames(userId, GetCurrentGameQuery)

	return
}

func (db *Database) getHistoryGames(userId int, query string) (games common.Games, err error) {
	rows, err := db_access.Query(query, userId, common.GameStateFinished, common.GameStateCancelled)

	if err != nil {
		return
	}

	for rows.Next() {
		game := common.Game{}
		var finishDateString *string
		err = rows.Scan(&game.Id, &game.Name, &game.State, &finishDateString)

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
		COALESCE(
			SUM(
				t.DurationInS -
				CASE t.State
					WHEN ? THEN t.RemainingTimeInS - (strftime('%s','now') - strftime('%s', t.LastActionDate))
					WHEN ? THEN t.RemainingTimeInS
					WHEN ? THEN t.RemainingTimeInS
					ELSE t.DurationInS
				END
			),
			0
	    ) AS SecondsSpent
	FROM Timers t
	WHERE t.UserId = ?
		AND t.GameId = ?
`

func (db *Database) GetGameTimeSpentCommand(userId int, gameId int) (timeSpent time.Duration, err error) {
	row := db_access.QueryRow(GetGameSecondsSpentQuery,
		common.TimerStateRunning,
		common.TimerStatePaused,
		common.TimerStateFinished,
		userId,
		gameId)

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

func (db *Database) CancelCurrentGameCommand(userId int, gameId int) error {
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

func (db *Database) FinishCurrentGameCommand(userId int, gameId int) error {
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
	WHERE gh.UserId = ?
		AND gh.State IN (?, ?)
	ORDER BY gh.FinishDate NULLS FIRST
`

func (db *Database) GetGameHistoryCommand(userId int) (games common.Games, err error) {
	games, err = db.getHistoryGames(userId, GetGameHistoryQuery)

	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}

	return
}
