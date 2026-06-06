package srvgames

import (
	"FGG-Service/src/common"
	"FGG-Service/src/games/database"
	"FGG-Service/src/games/types"
	"database/sql"
	"errors"
)

type IGettingService interface {
	GetWishlistGames(userId int) (typegames.WishlistGames, error)
	GetCurrentGame(userId int) (typegames.CurrentGame, error)
}

type GettingService struct {
	Database dbgames.IDatabase
}

func NewGettingService() IGettingService {
	db := new(dbgames.Database)

	return &GettingService{
		Database: db,
	}
}

func (s *GettingService) GetCurrentGame(userId int) (game typegames.CurrentGame, err error) {
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

func (s *GettingService) GetWishlistGames(userId int) (typegames.WishlistGames, error) {
	return s.Database.GetWishlistGamesCommand(userId)
}
