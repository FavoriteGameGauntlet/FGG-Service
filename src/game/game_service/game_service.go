package game_service

import (
	"FGG-Service/src/common"
	"FGG-Service/src/game/game_db"
	"FGG-Service/src/timer/timer_service"
	"database/sql"
	"errors"
	"math/rand"
	"time"
)

func AddUnplayedGames(userId int, games common.UnplayedGames) error {
	numberOfGames := len(games)
	errorCount := 0

	var err error
	for _, game := range games {
		err = CreateUnplayedGame(userId, &game)

		if err != nil {
			errorCount++
		}
	}

	if errorCount == numberOfGames {
		return err
	}

	return nil
}

func CreateUnplayedGame(userId int, unplayedGame *common.UnplayedGame) error {
	doesExist, err := game_db.DoesUnplayedGameExistCommand(userId, unplayedGame.Name)

	if err != nil {
		return err
	}

	if doesExist {
		return common.NewUnplayedGameAlreadyExistsError(unplayedGame.Name)
	}

	doesExist, err = game_db.DoesGameExistCommand(unplayedGame.Name)

	if err != nil {
		return err
	}

	game := common.Game{}

	if doesExist {
		game, err = game_db.GetGameCommand(unplayedGame.Name)
	} else {
		game, err = CreateGame(unplayedGame)
	}

	if err != nil {
		return err
	}

	err = game_db.CreateUnplayedGameCommand(userId, game.Id)

	return err
}

func GetUnplayedGames(userId int) (common.UnplayedGames, error) {
	return game_db.GetUnplayedGamesCommand(userId)
}

func CreateGame(unplayedGame *common.UnplayedGame) (game common.Game, err error) {
	err = game_db.CreateGameCommand(unplayedGame.Name)

	if err != nil {
		return
	}

	game, err = game_db.GetGameCommand(unplayedGame.Name)

	return
}

func GetCurrentGame(userId int) (game common.Game, err error) {
	game, err = game_db.GetCurrentGameCommand(userId)

	if err != nil {
		return
	}

	timeSpent, err := game_db.GetGameTimeSpentCommand(userId, game.Id)

	if err != nil {
		return
	}

	game.TimeSpent = timeSpent

	return
}

func CancelCurrentGame(userId int) error {
	game, err := GetCurrentGame(userId)

	if err != nil {
		return err
	}

	_, err = timer_service.StopCurrentTimer(userId)

	if err != nil {
		return err
	}

	err = game_db.CancelCurrentGameCommand(userId, game.Id)

	if err != nil {
		return err
	}

	return nil
}

func FinishCurrentGame(userId int) error {
	game, err := GetCurrentGame(userId)

	if err != nil {
		return err
	}

	if game.TimeSpent == 0 {
		return common.NewCompletedTimersNotFoundError()
	}

	_, err = timer_service.StopCurrentTimer(userId)

	if err != nil {
		return err
	}

	err = game_db.FinishCurrentGameCommand(userId, game.Id)

	if err != nil {
		return err
	}

	return nil
}

func GetGameHistory(userId int) (games common.Games, err error) {
	games, err = game_db.GetGameHistoryCommand(userId)

	for _, game := range games {
		var timeSpent time.Duration
		timeSpent, err = game_db.GetGameTimeSpentCommand(userId, game.Id)

		if err != nil {
			return
		}

		game.TimeSpent = timeSpent
	}

	return
}

func MakeGameRoll(userId int) (game common.Game, err error) {
	game, err = GetCurrentGame(userId)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return game, err
	}

	if game.Name != "" {
		err = common.NewCurrentGameAlreadyExistsError()
		return
	}

	unplayedGames, err := game_db.GetUnplayedGamesCommand(userId)

	if err != nil {
		return
	}

	if unplayedGames == nil || len(unplayedGames) < common.MinimumNumberOfUnplayedGames {
		err = common.NewUnplayedGamesNotFoundError()
		return
	}

	randomNumber := rand.Intn(len(unplayedGames))
	randomUnplayedGame := unplayedGames[randomNumber]

	err = game_db.CreateCurrentGameCommand(userId, randomUnplayedGame.GameId)

	if err != nil {
		return
	}

	err = game_db.DeleteUnplayedGameCommand(userId, randomUnplayedGame.GameId)

	if err != nil {
		return
	}

	game.Id = randomUnplayedGame.GameId
	game.Name = randomUnplayedGame.Name
	game.State = common.GameStateStarted

	return
}
