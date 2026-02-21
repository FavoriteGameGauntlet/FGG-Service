package game_service

import (
	"FGG-Service/src/common"
	"FGG-Service/src/games/game_db"
	"FGG-Service/src/timers/timer_service"
	"database/sql"
	"errors"
	"math/rand"
	"time"
)

type Service struct {
}

func (s *Service) AddUnplayedGames(userId int, games common.UnplayedGames) error {
	numberOfGames := len(games)
	errorCount := 0

	var err error
	for _, game := range games {
		err = CreateUnplayedGame(userId, game)

		if err != nil {
			errorCount++
		}
	}

	if errorCount == numberOfGames {
		return err
	}

	return nil
}

func (s *Service) CreateUnplayedGame(userId int, unplayedGame common.UnplayedGame) error {
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
		game, err = createGame(unplayedGame)
	}

	if err != nil {
		return err
	}

	err = game_db.CreateUnplayedGameCommand(userId, game.Id)

	return err
}

func (s *Service) GetUnplayedGames(userId int) (common.UnplayedGames, error) {
	return game_db.GetUnplayedGamesCommand(userId)
}

func (s *Service) createGame(unplayedGame common.UnplayedGame) (game common.Game, err error) {
	err = game_db.CreateGameCommand(unplayedGame.Name)

	if err != nil {
		return
	}

	game, err = game_db.GetGameCommand(unplayedGame.Name)

	return
}

func (s *Service) GetCurrentGame(userId int) (game common.Game, err error) {
	games, err := game_db.GetCurrentGameCommand(userId)

	if errors.Is(err, sql.ErrNoRows) || len(games) == 0 {
		err = common.NewCurrentGameNotFoundError()
		return
	}

	if err != nil {
		return
	}

	game = games[0]

	timeSpent, err := game_db.GetGameTimeSpentCommand(userId, game.Id)

	if err != nil {
		return
	}

	game.TimeSpent = timeSpent

	return
}

func (s *Service) CancelCurrentGame(userId int) error {
	game, err := GetCurrentGame(userId)

	if err != nil {
		return err
	}

	_, err = timer_service.ForceStopCurrentTimer(userId)

	if err != nil {
		return err
	}

	err = game_db.CancelCurrentGameCommand(userId, game.Id)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) FinishCurrentGame(userId int) error {
	game, err := GetCurrentGame(userId)

	if err != nil {
		return err
	}

	if game.TimeSpent == 0 {
		return common.NewCompletedTimersNotFoundError()
	}

	_, err = timer_service.ForceStopCurrentTimer(userId)

	if err != nil {
		return err
	}

	err = game_db.FinishCurrentGameCommand(userId, game.Id)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetGameHistory(userId int) (games common.Games, err error) {
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

func (s *Service) MakeGameRoll(userId int) (game common.Game, err error) {
	game, err = GetCurrentGame(userId)

	var notFoundError *common.NotFoundError
	if err != nil && !errors.As(err, &notFoundError) {
		return
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
