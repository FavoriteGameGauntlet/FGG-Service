package dbgames

import (
	"FGG-Service/src/dbaccess"
	"FGG-Service/src/games/types"
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

var doesGameExistQuery = dbaccess.Query{Name: "DoesGameExistQuery", SQL: `SELECT does_game_exist($1::text)`}

func (db *Database) DoesGameExistCommand(gameName string) (doesExist bool, err error) {
	row := dbaccess.QueryRow(doesGameExistQuery, gameName)

	err = row.Scan(&doesExist)

	dbaccess.LogDbResult(doesGameExistQuery, doesExist, err)

	return
}

var createGameQuery = dbaccess.Query{Name: "CreateGameQuery", SQL: `SELECT create_game($1::text)`}

func (db *Database) CreateGameCommand(name string) error {
	_, err := dbaccess.Exec(createGameQuery, name)

	dbaccess.LogDbResult(createGameQuery, nil, err)

	return err
}

var getWishlistGameQuery = dbaccess.Query{Name: "GetWishlistGameQuery", SQL: `SELECT * FROM get_wishlist_game($1::text)`}

func (db *Database) GetWishlistGameCommand(name string) (game typegames.WishlistGame, err error) {
	row := dbaccess.QueryRow(getWishlistGameQuery, name)

	err = row.Scan(&game.GameId, &game.Name)

	dbaccess.LogDbResult(getWishlistGameQuery, game, err)

	return
}

var doesUnplayedGameExistQuery = dbaccess.Query{Name: "DoesUnplayedGameExistQuery", SQL: `SELECT does_wishlist_game_exist($1::integer, $2::text)`}

func (db *Database) DoesWishlistGameExistCommand(userId int, gameName string) (doesExist bool, err error) {
	row := dbaccess.QueryRow(doesUnplayedGameExistQuery, userId, gameName)

	err = row.Scan(&doesExist)

	dbaccess.LogDbResult(doesUnplayedGameExistQuery, doesExist, err)

	return
}

var createUnplayedGameQuery = dbaccess.Query{Name: "CreateUnplayedGameQuery", SQL: `SELECT create_wishlist_game($1::integer, $2::integer)`}

func (db *Database) CreateWishlistGameCommand(userId int, gameId int) error {
	_, err := dbaccess.Exec(createUnplayedGameQuery, userId, gameId)

	dbaccess.LogDbResult(createUnplayedGameQuery, nil, err)

	return err
}

var deleteUnplayedGameQuery = dbaccess.Query{Name: "DeleteUnplayedGameQuery", SQL: `SELECT delete_wishlist_game($1::integer, $2::integer)`}

func (db *Database) DeleteUnplayedGameCommand(userId int, gameId int) error {
	_, err := dbaccess.Exec(deleteUnplayedGameQuery, userId, gameId)

	dbaccess.LogDbResult(deleteUnplayedGameQuery, nil, err)

	return err
}

var getWishlistGamesQuery = dbaccess.Query{Name: "GetWishlistGamesQuery", SQL: `SELECT * FROM get_wishlist_games($1::integer)`}

func (db *Database) GetWishlistGamesCommand(userId int) (games typegames.WishlistGames, err error) {
	rows, err := dbaccess.QueryRows(getWishlistGamesQuery, userId)

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

	dbaccess.LogDbResult(getWishlistGamesQuery, games, err)

	_ = rows.Close()
	return
}

var createCurrentGameQuery = dbaccess.Query{Name: "CreateCurrentGameQuery", SQL: `SELECT create_current_game($1::integer, $2::integer)`}

func (db *Database) CreateCurrentGameCommand(userId int, gameId int) error {
	_, err := dbaccess.Exec(createCurrentGameQuery, userId, gameId)

	dbaccess.LogDbResult(createCurrentGameQuery, nil, err)

	return err
}

var getCurrentGameQuery = dbaccess.Query{Name: "GetCurrentGameQuery", SQL: `SELECT * FROM get_current_game($1::integer)`}

func (db *Database) GetCurrentGameCommand(userId int) (games typegames.CurrentGames, err error) {
	games, err = db.getHistoryGames(getCurrentGameQuery, userId)

	return
}

func (db *Database) getHistoryGames(q dbaccess.Query, userId int) (games typegames.CurrentGames, err error) {
	rows, err := dbaccess.QueryRows(q, userId)

	if err != nil {
		return
	}

	for rows.Next() {
		game := typegames.CurrentGame{}
		err = rows.Scan(&game.Id, &game.Name, &game.State, &game.StartDate, &game.FinishDate)

		if err != nil {
			dbaccess.LogDbResult(q, games, err)

			_ = rows.Close()
			return
		}

		games = append(games, game)
	}

	dbaccess.LogDbResult(q, games, err)

	_ = rows.Close()
	return
}

var getGameSecondsSpentQuery = dbaccess.Query{Name: "GetGameSecondsSpentQuery", SQL: `SELECT get_game_seconds_spent($1::integer, $2::integer)`}

func (db *Database) GetGameTimeSpentCommand(userId int, gameId int) (timeSpent time.Duration, err error) {
	row := dbaccess.QueryRow(getGameSecondsSpentQuery, userId, gameId)

	var secondsSpent int
	err = row.Scan(&secondsSpent)

	if errors.Is(err, sql.ErrNoRows) {
		dbaccess.LogDbResult(getGameSecondsSpentQuery, timeSpent, err)

		err = nil
		return
	}

	if err != nil {
		dbaccess.LogDbResult(getGameSecondsSpentQuery, timeSpent, err)

		return
	}

	timeSpent = time.Duration(secondsSpent) * time.Second

	dbaccess.LogDbResult(getGameSecondsSpentQuery, timeSpent, err)

	return
}

var cancelCurrentGameQuery = dbaccess.Query{Name: "CancelCurrentGameQuery", SQL: `SELECT cancel_current_game($1::integer, $2::integer)`}

func (db *Database) CancelCurrentGameCommand(userId int, gameId int) error {
	_, err := dbaccess.Exec(cancelCurrentGameQuery, userId, gameId)

	dbaccess.LogDbResult(cancelCurrentGameQuery, nil, err)

	return err
}

var finishCurrentGameQuery = dbaccess.Query{Name: "FinishCurrentGameQuery", SQL: `SELECT finish_current_game($1::integer, $2::integer)`}

func (db *Database) FinishCurrentGameCommand(userId int, gameId int) error {
	_, err := dbaccess.Exec(finishCurrentGameQuery, userId, gameId)

	dbaccess.LogDbResult(finishCurrentGameQuery, nil, err)

	return err
}

var getGameHistoryQuery = dbaccess.Query{Name: "GetGameHistoryQuery", SQL: `SELECT * FROM get_game_history($1::integer)`}

func (db *Database) GetGameHistoryCommand(userId int) (games typegames.CurrentGames, err error) {
	games, err = db.getHistoryGames(getGameHistoryQuery, userId)

	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}

	return
}

var getAllCurrentGamesQuery = dbaccess.Query{Name: "GetAllCurrentGamesQuery", SQL: `SELECT * FROM get_all_current_games()`}

func (db *Database) GetAllCurrentGamesCommand() (games []typegames.CurrentGameWithLogin, err error) {
	rows, err := dbaccess.QueryRows(getAllCurrentGamesQuery)

	if err != nil {
		return
	}

	for rows.Next() {
		game := typegames.CurrentGame{}
		var login string
		err = rows.Scan(&game.Id, &game.Name, &game.State, &game.StartDate, &game.FinishDate, &login)

		if err != nil {
			dbaccess.LogDbResult(getAllCurrentGamesQuery, games, err)
			_ = rows.Close()
			return
		}

		games = append(games, typegames.CurrentGameWithLogin{Login: login, Game: game})
	}

	dbaccess.LogDbResult(getAllCurrentGamesQuery, games, err)
	_ = rows.Close()
	return
}
