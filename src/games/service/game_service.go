package game_service

import (
	"FGG-Service/src/common"
	"FGG-Service/src/games/db"
	"FGG-Service/src/timers/service"
	"database/sql"
	"errors"
	"math/rand"
	"time"
)

type Service struct {
	Db game_db.Database
}

func (s *Service) AddUnplayedGames(userId int, games common.UnplayedGames) error {
	numberOfGames := len(games)
	errorCount := 0

	var err error
	for _, game := range games {
		err = s.CreateUnplayedGame(userId, game)

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
	doesExist, err := s.Db.DoesUnplayedGameExistCommand(userId, unplayedGame.Name)

	if err != nil {
		return err
	}

	if doesExist {
		return common.NewUnplayedGameAlreadyExistsError(unplayedGame.Name)
	}

	doesExist, err = s.Db.DoesGameExistCommand(unplayedGame.Name)

	if err != nil {
		return err
	}

	game := common.Game{}

	if doesExist {
		game, err = s.Db.GetGameCommand(unplayedGame.Name)
	} else {
		game, err = s.createGame(unplayedGame)
	}

	if err != nil {
		return err
	}

	err = s.Db.CreateUnplayedGameCommand(userId, game.Id)

	return err
}

func (s *Service) GetUnplayedGames(userId int) (common.UnplayedGames, error) {
	return s.Db.GetUnplayedGamesCommand(userId)
}

func (s *Service) createGame(unplayedGame common.UnplayedGame) (game common.Game, err error) {
	err = s.Db.CreateGameCommand(unplayedGame.Name)

	if err != nil {
		return
	}

	game, err = s.Db.GetGameCommand(unplayedGame.Name)

	return
}

func (s *Service) GetCurrentGame(userId int) (game common.Game, err error) {
	games, err := s.Db.GetCurrentGameCommand(userId)

	if errors.Is(err, sql.ErrNoRows) || len(games) == 0 {
		err = common.NewCurrentGameNotFoundError()
		return
	}

	if err != nil {
		return
	}

	game = games[0]

	timeSpent, err := s.Db.GetGameTimeSpentCommand(userId, game.Id)

	if err != nil {
		return
	}

	game.TimeSpent = timeSpent

	return
}

func (s *Service) CancelCurrentGame(userId int) error {
	game, err := s.GetCurrentGame(userId)

	if err != nil {
		return err
	}

	_, err = timer_service.ForceStopCurrentTimer(userId)

	if err != nil {
		return err
	}

	err = s.Db.CancelCurrentGameCommand(userId, game.Id)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) FinishCurrentGame(userId int) error {
	game, err := s.GetCurrentGame(userId)

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

	err = s.Db.FinishCurrentGameCommand(userId, game.Id)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetGameHistory(userId int) (games common.Games, err error) {
	games, err = s.Db.GetGameHistoryCommand(userId)

	for _, game := range games {
		var timeSpent time.Duration
		timeSpent, err = s.Db.GetGameTimeSpentCommand(userId, game.Id)

		if err != nil {
			return
		}

		game.TimeSpent = timeSpent
	}

	return
}

func (s *Service) MakeGameRoll(userId int) (game common.Game, err error) {
	game, err = s.GetCurrentGame(userId)

	var notFoundError *common.NotFoundError
	if err != nil && !errors.As(err, &notFoundError) {
		return
	}

	if game.Name != "" {
		err = common.NewCurrentGameAlreadyExistsError()
		return
	}

	unplayedGames, err := s.Db.GetUnplayedGamesCommand(userId)

	if err != nil {
		return
	}

	if unplayedGames == nil || len(unplayedGames) < common.MinimumNumberOfUnplayedGames {
		err = common.NewUnplayedGamesNotFoundError()
		return
	}

	randomNumber := rand.Intn(len(unplayedGames))
	randomUnplayedGame := unplayedGames[randomNumber]

	err = s.Db.CreateCurrentGameCommand(userId, randomUnplayedGame.GameId)

	if err != nil {
		return
	}

	err = s.Db.DeleteUnplayedGameCommand(userId, randomUnplayedGame.GameId)

	if err != nil {
		return
	}

	game.Id = randomUnplayedGame.GameId
	game.Name = randomUnplayedGame.Name
	game.State = common.GameStateStarted

	return
}
