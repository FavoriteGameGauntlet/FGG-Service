package dbgames

import (
	"FGG-Service/src/dbaccess"
	"FGG-Service/src/games/types"
	"FGG-Service/src/timers/types"
	"database/sql"
	"errors"
	"time"
)

type IDatabase interface {
	DoesGameExistCommand(gameName string) (doesExist bool, err error)
	CreateGameCommand(name string) error
	GetWishlistGameCommand(name string) (game typegames.WishlistGame, err error)
	DoesWishlistGameExistCommand(userId int, gameName string) (doesExist bool, err error)
	CreateWishlistGameCommand(userId int, gameId int) error
	DeleteUnplayedGameCommand(userId int, gameId int) error
	GetWishlistGamesCommand(userId int) (games typegames.WishlistGames, err error)
	CreateCurrentGameCommand(userId int, gameId int) error
	GetCurrentGameCommand(userId int) (games typegames.CurrentGames, err error)
	GetGameTimeSpentCommand(userId int, gameId int) (timeSpent time.Duration, err error)
	CancelCurrentGameCommand(userId int, gameId int) error
	FinishCurrentGameCommand(userId int, gameId int) error
	GetGameHistoryCommand(userId int) (games typegames.CurrentGames, err error)
	GetAllCurrentGamesCommand() (games []typegames.CurrentGameWithLogin, err error)
}

type Database struct {
}

const DoesGameExistQuery = `
	SELECT
		CASE WHEN EXISTS (
			SELECT 1
			FROM Games
			WHERE Name = $1
		)
		THEN true
		ELSE false
	END AS DoesExist`

func (db *Database) DoesGameExistCommand(gameName string) (doesExist bool, err error) {
	queryName := "DoesGameExistQuery"
	row := dbaccess.QueryRow(queryName, DoesGameExistQuery, gameName)

	err = row.Scan(&doesExist)

	dbaccess.LogDbResult(queryName, doesExist, err)

	return
}

const CreateGameQuery = `
	INSERT INTO Games (Name)
	VALUES ($1)
`

func (db *Database) CreateGameCommand(name string) error {
	queryName := "CreateGameQuery"
	_, err := dbaccess.Exec(queryName, CreateGameQuery, name)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const GetWishlistGameQuery = `
	SELECT Id, Name
	FROM Games
	WHERE Name = $1
`

func (db *Database) GetWishlistGameCommand(name string) (game typegames.WishlistGame, err error) {
	queryName := "GetWishlistGameQuery"
	row := dbaccess.QueryRow(queryName, GetWishlistGameQuery, name)

	err = row.Scan(&game.GameId, &game.Name)

	dbaccess.LogDbResult(queryName, game, err)

	return
}

const DoesUnplayedGameExistQuery = `
	SELECT
	    CASE WHEN EXISTS (
			SELECT 1
			FROM UnplayedGames ug
				INNER JOIN Games g ON ug.GameId = g.Id
			WHERE ug.UserId = $1
				AND g.Name = $2
		)
		THEN true
		ELSE false
	END AS DoesExist`

func (db *Database) DoesWishlistGameExistCommand(userId int, gameName string) (doesExist bool, err error) {
	queryName := "DoesUnplayedGameExistQuery"
	row := dbaccess.QueryRow(queryName, DoesUnplayedGameExistQuery, userId, gameName)

	err = row.Scan(&doesExist)

	dbaccess.LogDbResult(queryName, doesExist, err)

	return
}

const CreateUnplayedGameQuery = `
	INSERT INTO UnplayedGames (UserId, GameId)
	VALUES ($1, $2)
`

func (db *Database) CreateWishlistGameCommand(userId int, gameId int) error {
	queryName := "CreateUnplayedGameQuery"
	_, err := dbaccess.Exec(queryName, CreateUnplayedGameQuery, userId, gameId)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const DeleteUnplayedGameQuery = `
	DELETE FROM UnplayedGames
	WHERE UserId = $1
		AND GameId = $2
`

func (db *Database) DeleteUnplayedGameCommand(userId int, gameId int) error {
	queryName := "DeleteUnplayedGameQuery"
	_, err := dbaccess.Exec(queryName, DeleteUnplayedGameQuery, userId, gameId)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const GetWishlistGamesQuery = `
	SELECT ug.Id, g.Id, g.Name
	FROM UnplayedGames ug
		INNER JOIN Games g ON ug.GameId = g.Id
	WHERE ug.UserId = $1
`

func (db *Database) GetWishlistGamesCommand(userId int) (games typegames.WishlistGames, err error) {
	queryName := "GetWishlistGamesQuery"
	rows, err := dbaccess.Query(queryName, GetWishlistGamesQuery, userId)

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

	dbaccess.LogDbResult(queryName, games, err)

	_ = rows.Close()
	return
}

const CreateCurrentGameQuery = `
	INSERT INTO GameHistory (UserId, GameId)
	VALUES ($1, $2)
`

func (db *Database) CreateCurrentGameCommand(userId int, gameId int) error {
	queryName := "CreateCurrentGameQuery"
	_, err := dbaccess.Exec(queryName, CreateCurrentGameQuery, userId, gameId)

	dbaccess.LogDbResult(queryName, nil, err)

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
	WHERE gh.UserId = $1
		AND gh.State NOT IN ($2, $3)
`

func (db *Database) GetCurrentGameCommand(userId int) (games typegames.CurrentGames, err error) {
	queryName := "GetCurrentGameQuery"
	games, err = db.getHistoryGames(queryName, GetCurrentGameQuery, &userId)

	return
}

func (db *Database) getHistoryGames(queryName string, query string, userId *int) (games typegames.CurrentGames, err error) {
	var rows *sql.Rows

	if userId == nil {
		rows, err = dbaccess.Query(queryName, query, typegames.GameStateFinished, typegames.GameStateCancelled)
	} else {
		rows, err = dbaccess.Query(queryName, query, userId, typegames.GameStateFinished, typegames.GameStateCancelled)
	}

	if err != nil {
		return
	}

	for rows.Next() {
		game := typegames.CurrentGame{}
		err = rows.Scan(&game.Id, &game.Name, &game.State, &game.FinishDate)

		if err != nil {
			dbaccess.LogDbResult(queryName, games, err)

			_ = rows.Close()
			return
		}

		games = append(games, game)
	}

	dbaccess.LogDbResult(queryName, games, err)

	_ = rows.Close()
	return
}

const GetGameSecondsSpentQuery = `
	SELECT
		COALESCE(
			SUM(
				t.DurationInS -
				CASE t.State
					WHEN $1 THEN t.RemainingTimeInS - CAST(EXTRACT(EPOCH FROM (NOW() - t.LastActionDate)) AS INTEGER)
					WHEN $2 THEN t.RemainingTimeInS
					WHEN $3 THEN t.RemainingTimeInS
					ELSE t.DurationInS
				END
			),
			0
	    ) AS SecondsSpent
	FROM Timers t
	WHERE t.UserId = $4
		AND t.GameId = $5
`

func (db *Database) GetGameTimeSpentCommand(userId int, gameId int) (timeSpent time.Duration, err error) {
	queryName := "GetGameSecondsSpentQuery"
	row := dbaccess.QueryRow(
		queryName,
		GetGameSecondsSpentQuery,
		typetimers.TimerStateRunning,
		typetimers.TimerStatePaused,
		typetimers.TimerStateFinished,
		userId,
		gameId)

	var secondsSpent int
	err = row.Scan(&secondsSpent)

	if errors.Is(err, sql.ErrNoRows) {
		dbaccess.LogDbResult(queryName, timeSpent, err)

		err = nil
		return
	}

	if err != nil {
		dbaccess.LogDbResult(queryName, timeSpent, err)

		return
	}

	timeSpent = time.Duration(secondsSpent) * time.Second

	dbaccess.LogDbResult(queryName, timeSpent, err)

	return
}

const CancelCurrentGameQuery = `
	UPDATE GameHistory
	SET State = $1,
		FinishDate = NOW()
	WHERE UserId = $2
		AND GameId = $3
`

func (db *Database) CancelCurrentGameCommand(userId int, gameId int) error {
	queryName := "CancelCurrentGameQuery"
	_, err := dbaccess.Exec(queryName, CancelCurrentGameQuery, typegames.GameStateCancelled, userId, gameId)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const FinishCurrentGameQuery = `
	UPDATE GameHistory
	SET State = $1,
		FinishDate = NOW()
	WHERE UserId = $2
		AND GameId = $3
`

func (db *Database) FinishCurrentGameCommand(userId int, gameId int) error {
	queryName := "FinishCurrentGameQuery"
	_, err := dbaccess.Exec(queryName, FinishCurrentGameQuery, typegames.GameStateFinished, userId, gameId)

	dbaccess.LogDbResult(queryName, nil, err)

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
	WHERE gh.UserId = $1
		AND gh.State IN ($2, $3)
	ORDER BY gh.FinishDate NULLS FIRST
`

func (db *Database) GetGameHistoryCommand(userId int) (games typegames.CurrentGames, err error) {
	queryName := "GetGameHistoryQuery"
	games, err = db.getHistoryGames(queryName, GetGameHistoryQuery, &userId)

	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}

	return
}

const GetAllCurrentGamesQuery = `
	SELECT
		g.Id,
		g.Name,
		gh.State,
		gh.FinishDate,
		u.Login
	FROM GameHistory gh
		INNER JOIN Games g ON gh.GameId = g.Id
		INNER JOIN Users u ON gh.UserId = u.Id
	WHERE gh.State NOT IN ($1, $2)
`

func (db *Database) GetAllCurrentGamesCommand() (games []typegames.CurrentGameWithLogin, err error) {
	queryName := "GetAllCurrentGamesQuery"
	rows, err := dbaccess.Query(queryName, GetAllCurrentGamesQuery, typegames.GameStateFinished, typegames.GameStateCancelled)

	if err != nil {
		return
	}

	for rows.Next() {
		game := typegames.CurrentGame{}
		var login string
		err = rows.Scan(&game.Id, &game.Name, &game.State, &game.FinishDate, &login)

		if err != nil {
			dbaccess.LogDbResult(queryName, games, err)
			_ = rows.Close()
			return
		}

		games = append(games, typegames.CurrentGameWithLogin{Login: login, Game: game})
	}

	dbaccess.LogDbResult(queryName, games, err)
	_ = rows.Close()
	return
}
