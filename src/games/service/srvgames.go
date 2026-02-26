package srvgames

import (
	"FGG-Service/src/common"
	"FGG-Service/src/games/database"
	"FGG-Service/src/games/types"
	"FGG-Service/src/timers/service"
	"database/sql"
	"errors"
	"math/rand"
	"time"
)

type Service struct {
	Database     dbgames.IDatabase
	TimerService srvtimers.Service
}

func (s *Service) AddUnplayedGames(userId int, games typegames.WishlistGames) error {
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

func (s *Service) CreateUnplayedGame(userId int, unplayedGame typegames.WishlistGame) error {
	doesExist, err := s.Database.DoesUnplayedGameExistCommand(userId, unplayedGame.Name)

	if err != nil {
		return err
	}

	if doesExist {
		return common.NewUnplayedGameAlreadyExistsError(unplayedGame.Name)
	}

	doesExist, err = s.Database.DoesGameExistCommand(unplayedGame.Name)

	if err != nil {
		return err
	}

	game := typegames.CurrentGame{}

	if doesExist {
		game, err = s.Database.GetGameCommand(unplayedGame.Name)
	} else {
		game, err = s.createGame(unplayedGame)
	}

	if err != nil {
		return err
	}

	err = s.Database.CreateUnplayedGameCommand(userId, game.Id)

	return err
}

func (s *Service) GetUnplayedGames(userId int) (typegames.WishlistGames, error) {
	return s.Database.GetUnplayedGamesCommand(userId)
}

func (s *Service) createGame(unplayedGame typegames.WishlistGame) (game typegames.CurrentGame, err error) {
	err = s.Database.CreateGameCommand(unplayedGame.Name)

	if err != nil {
		return
	}

	game, err = s.Database.GetGameCommand(unplayedGame.Name)

	return
}

func (s *Service) GetCurrentGame(userId int) (game typegames.CurrentGame, err error) {
	games, err := s.Database.GetCurrentGameCommand(userId)

	if errors.Is(err, sql.ErrNoRows) || len(games) == 0 {
		err = common.NewCurrentGameNotFoundError()
		return
	}

	if err != nil {
		return
	}

	game = games[0]

	timeSpent, err := s.Database.GetGameTimeSpentCommand(userId, game.Id)

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

	_, err = s.TimerService.ForceStopCurrentTimer(userId)

	if err != nil {
		return err
	}

	err = s.Database.CancelCurrentGameCommand(userId, game.Id)

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

	_, err = s.TimerService.ForceStopCurrentTimer(userId)

	if err != nil {
		return err
	}

	err = s.Database.FinishCurrentGameCommand(userId, game.Id)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetGameHistory(userId int) (games typegames.CurrentGames, err error) {
	games, err = s.Database.GetGameHistoryCommand(userId)

	for _, game := range games {
		var timeSpent time.Duration
		timeSpent, err = s.Database.GetGameTimeSpentCommand(userId, game.Id)

		if err != nil {
			return
		}

		game.TimeSpent = timeSpent
	}

	return
}

func (s *Service) MakeGameRoll(userId int) (game typegames.CurrentGame, err error) {
	game, err = s.GetCurrentGame(userId)

	var notFoundError *common.NotFoundError
	if err != nil && !errors.As(err, &notFoundError) {
		return
	}

	if game.Name != "" {
		err = common.NewCurrentGameAlreadyExistsError()
		return
	}

	unplayedGames, err := s.Database.GetUnplayedGamesCommand(userId)

	if err != nil {
		return
	}

	if unplayedGames == nil || len(unplayedGames) < common.MinimumNumberOfUnplayedGames {
		err = common.NewUnplayedGamesNotFoundError()
		return
	}

	randomNumber := rand.Intn(len(unplayedGames))
	randomUnplayedGame := unplayedGames[randomNumber]

	err = s.Database.CreateCurrentGameCommand(userId, randomUnplayedGame.GameId)

	if err != nil {
		return
	}

	err = s.Database.DeleteUnplayedGameCommand(userId, randomUnplayedGame.GameId)

	if err != nil {
		return
	}

	game.Id = randomUnplayedGame.GameId
	game.Name = randomUnplayedGame.Name
	game.State = typegames.GameStateStarted

	return
}
