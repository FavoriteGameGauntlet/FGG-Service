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
	Database       dbgames.IDatabase
	TimerService   srvtimers.IService
	GettingService IGettingService
}

func NewService() *Service {
	db := new(dbgames.Database)
	ts := srvtimers.NewService()
	gs := NewGettingService()

	return &Service{
		Database:       db,
		TimerService:   ts,
		GettingService: gs,
	}
}

func (s *Service) AddWishlistGame(userId int, wishlistGame typegames.WishlistGame) error {
	doesExist, err := s.Database.DoesWishlistGameExistCommand(userId, wishlistGame.Name)

	if err != nil {
		return err
	}

	if doesExist {
		return common.NewWishlistGameAlreadyExistsConflictError(wishlistGame.Name)
	}

	doesExist, err = s.Database.DoesGameExistCommand(wishlistGame.Name)

	if err != nil {
		return err
	}

	game := typegames.WishlistGame{}

	if doesExist {
		game, err = s.Database.GetWishlistGameCommand(wishlistGame.Name)
	} else {
		game, err = s.createAndGetGame(wishlistGame)
	}

	if err != nil {
		return err
	}

	err = s.Database.CreateWishlistGameCommand(userId, game.GameId)

	return err
}

func (s *Service) GetUnplayedGames(userId int) (typegames.WishlistGames, error) {
	return s.GettingService.GetWishlistGames(userId)
}

func (s *Service) createAndGetGame(wishlistGame typegames.WishlistGame) (game typegames.WishlistGame, err error) {
	err = s.Database.CreateGameCommand(wishlistGame.Name)

	if err != nil {
		return
	}

	game, err = s.Database.GetWishlistGameCommand(wishlistGame.Name)

	return
}

func (s *Service) GetCurrentGame(userId int) (typegames.CurrentGame, error) {
	return s.GettingService.GetCurrentGame(userId)
}

func (s *Service) CancelCurrentGame(userId int) error {
	game, err := s.GettingService.GetCurrentGame(userId)

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
	game, err := s.GettingService.GetCurrentGame(userId)

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
	game, err = s.GettingService.GetCurrentGame(userId)

	var notFoundError *common.NotFoundError
	if err != nil && !errors.As(err, &notFoundError) {
		return
	}

	if game.Name != "" {
		err = common.NewCurrentGameAlreadyExistsConflictError()
		return
	}

	unplayedGames, err := s.GettingService.GetWishlistGames(userId)

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

func (s *Service) GetAllCurrentGames() (games typegames.CurrentGames, err error) {
	games, err = s.Database.GetAllCurrentGamesCommand()

	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}

	return
}
