package srvgames

import (
	"FGG-Service/src/common"
	"FGG-Service/src/games/database"
	"FGG-Service/src/games/types"
	"database/sql"
	"errors"
)

type IQueryService interface {
	GetCurrentGame(userId int) (typegames.CurrentGame, error)
}

type QueryService struct {
	Database dbgames.IDatabase
}

func (s *QueryService) GetCurrentGame(userId int) (game typegames.CurrentGame, err error) {
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
