package dbgames

import (
	"FGG-Service/src/common"
	"FGG-Service/src/dbaccess"
	"FGG-Service/src/games/types"
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
	row := dbaccess.QueryRow(DoesGameExistQuery, gameName)

	err = row.Scan(&doesExist)

	return
}

const CreateGameQuery = `
	INSERT INTO Games (Name)
	VALUES (?)
`

func (db *Database) CreateGameCommand(name string) error {
	_, err := dbaccess.Exec(
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

func (db *Database) GetGameCommand(name string) (game typegames.CurrentGame, err error) {
	row := dbaccess.QueryRow(GetGameQuery, name)

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
	row := dbaccess.QueryRow(DoesUnplayedGameExistQuery, userId, gameName)

	err = row.Scan(&doesExist)

	return
}

const CreateUnplayedGameQuery = `
	INSERT INTO UnplayedGames (UserId, GameId)
	VALUES (?, ?)
`

func (db *Database) CreateUnplayedGameCommand(userId int, gameId int) error {
	_, err := dbaccess.Exec(
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
	_, err := dbaccess.Exec(DeleteUnplayedGameQuery, userId, gameId)

	return err
}

const GetUnplayedGamesQuery = `
	SELECT ug.Id, g.Id, g.Name
	FROM UnplayedGames ug
		INNER JOIN Games g ON ug.GameId = g.Id
	WHERE ug.UserId = ?
`

func (db *Database) GetUnplayedGamesCommand(userId int) (games typegames.WishlistGames, err error) {
	rows, err := dbaccess.Query(GetUnplayedGamesQuery, userId)

	if err != nil {
		return
	}

	for rows.Next() {
		game := typegames.WishlistGame{}
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
	_, err := dbaccess.Exec(CreateCurrentGameQuery, userId, gameId)

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

func (db *Database) GetCurrentGameCommand(userId int) (games typegames.CurrentGames, err error) {
	games, err = db.getHistoryGames(userId, GetCurrentGameQuery)

	return
}

func (db *Database) getHistoryGames(userId int, query string) (games typegames.CurrentGames, err error) {
	rows, err := dbaccess.Query(query, userId, typegames.GameStateFinished, typegames.GameStateCancelled)

	if err != nil {
		return
	}

	for rows.Next() {
		game := typegames.CurrentGame{}
		var finishDateString *string
		err = rows.Scan(&game.Id, &game.Name, &game.State, &finishDateString)

		if err != nil {
			_ = rows.Close()
			return
		}

		var finishDate *time.Time
		finishDate, err = dbaccess.ConvertToNullableDate(finishDateString)

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
	row := dbaccess.QueryRow(GetGameSecondsSpentQuery,
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
	_, err := dbaccess.Exec(CancelCurrentGameQuery, typegames.GameStateCancelled, userId, gameId)

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
	_, err := dbaccess.Exec(FinishCurrentGameQuery, typegames.GameStateFinished, userId, gameId)

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

func (db *Database) GetGameHistoryCommand(userId int) (games typegames.CurrentGames, err error) {
	games, err = db.getHistoryGames(userId, GetGameHistoryQuery)

	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}

	return
}
