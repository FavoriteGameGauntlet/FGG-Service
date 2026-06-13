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

const DoesGameExistQuery = `SELECT does_game_exist($1)`

func (db *Database) DoesGameExistCommand(gameName string) (doesExist bool, err error) {
	queryName := "DoesGameExistQuery"
	row := dbaccess.QueryRow(queryName, DoesGameExistQuery, gameName)

	err = row.Scan(&doesExist)

	dbaccess.LogDbResult(queryName, doesExist, err)

	return
}

const CreateGameQuery = `SELECT create_game($1)`

func (db *Database) CreateGameCommand(name string) error {
	queryName := "CreateGameQuery"
	_, err := dbaccess.Exec(queryName, CreateGameQuery, name)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const GetWishlistGameQuery = `SELECT * FROM get_wishlist_game($1)`

func (db *Database) GetWishlistGameCommand(name string) (game typegames.WishlistGame, err error) {
	queryName := "GetWishlistGameQuery"
	row := dbaccess.QueryRow(queryName, GetWishlistGameQuery, name)

	err = row.Scan(&game.GameId, &game.Name)

	dbaccess.LogDbResult(queryName, game, err)

	return
}

const DoesUnplayedGameExistQuery = `SELECT does_wishlist_game_exist($1, $2)`

func (db *Database) DoesWishlistGameExistCommand(userId int, gameName string) (doesExist bool, err error) {
	queryName := "DoesUnplayedGameExistQuery"
	row := dbaccess.QueryRow(queryName, DoesUnplayedGameExistQuery, userId, gameName)

	err = row.Scan(&doesExist)

	dbaccess.LogDbResult(queryName, doesExist, err)

	return
}

const CreateUnplayedGameQuery = `SELECT create_wishlist_game($1, $2)`

func (db *Database) CreateWishlistGameCommand(userId int, gameId int) error {
	queryName := "CreateUnplayedGameQuery"
	_, err := dbaccess.Exec(queryName, CreateUnplayedGameQuery, userId, gameId)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const DeleteUnplayedGameQuery = `SELECT delete_wishlist_game($1, $2)`

func (db *Database) DeleteUnplayedGameCommand(userId int, gameId int) error {
	queryName := "DeleteUnplayedGameQuery"
	_, err := dbaccess.Exec(queryName, DeleteUnplayedGameQuery, userId, gameId)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const GetWishlistGamesQuery = `SELECT * FROM get_wishlist_games($1)`

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

const CreateCurrentGameQuery = `SELECT create_current_game($1, $2)`

func (db *Database) CreateCurrentGameCommand(userId int, gameId int) error {
	queryName := "CreateCurrentGameQuery"
	_, err := dbaccess.Exec(queryName, CreateCurrentGameQuery, userId, gameId)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const GetCurrentGameQuery = `SELECT * FROM get_current_game($1)`

func (db *Database) GetCurrentGameCommand(userId int) (games typegames.CurrentGames, err error) {
	queryName := "GetCurrentGameQuery"
	games, err = db.getHistoryGames(queryName, GetCurrentGameQuery, userId)

	return
}

func (db *Database) getHistoryGames(queryName string, query string, userId int) (games typegames.CurrentGames, err error) {
	rows, err := dbaccess.Query(queryName, query, userId)

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

const GetGameSecondsSpentQuery = `SELECT get_game_seconds_spent($1, $2)`

func (db *Database) GetGameTimeSpentCommand(userId int, gameId int) (timeSpent time.Duration, err error) {
	queryName := "GetGameSecondsSpentQuery"
	row := dbaccess.QueryRow(queryName, GetGameSecondsSpentQuery, userId, gameId)

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

const CancelCurrentGameQuery = `SELECT cancel_current_game($1, $2)`

func (db *Database) CancelCurrentGameCommand(userId int, gameId int) error {
	queryName := "CancelCurrentGameQuery"
	_, err := dbaccess.Exec(queryName, CancelCurrentGameQuery, userId, gameId)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const FinishCurrentGameQuery = `SELECT finish_current_game($1, $2)`

func (db *Database) FinishCurrentGameCommand(userId int, gameId int) error {
	queryName := "FinishCurrentGameQuery"
	_, err := dbaccess.Exec(queryName, FinishCurrentGameQuery, userId, gameId)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const GetGameHistoryQuery = `SELECT * FROM get_game_history($1)`

func (db *Database) GetGameHistoryCommand(userId int) (games typegames.CurrentGames, err error) {
	queryName := "GetGameHistoryQuery"
	games, err = db.getHistoryGames(queryName, GetGameHistoryQuery, userId)

	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}

	return
}

const GetAllCurrentGamesQuery = `SELECT * FROM get_all_current_games()`

func (db *Database) GetAllCurrentGamesCommand() (games []typegames.CurrentGameWithLogin, err error) {
	queryName := "GetAllCurrentGamesQuery"
	rows, err := dbaccess.Query(queryName, GetAllCurrentGamesQuery)

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
