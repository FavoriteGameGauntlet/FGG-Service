package controller

import (
	"github.com/labstack/echo/v4"
)

type Server struct{}

func (s Server) GetOwnCurrentGame(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) CancelOwnCurrentGame(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) FinishOwnCurrentGame(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) RollNewCurrentGame(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetOwnGameHistory(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetOwnWishlistGames(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) AddOwnWishlistGame(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetCurrentGameByLogin(ctx echo.Context, login string) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetGameHistoryByLogin(ctx echo.Context, login string) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetWishlistGamesByLogin(ctx echo.Context, login string) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetOwnExperiencePoints(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) ChangeOwnExperiencePoints(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetOwnFreePoints(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) ChangeOwnFreePoints(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetOwnFreePointHistory(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetOwnPointInfo(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetOwnTerritoryHours(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) ChangeOwnTerritoryHours(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetOwnTerritoryPoints(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) ChangeOwnTerritoryPoints(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetOwnTerritoryPointHistory(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetFreePointHistoryByLogin(ctx echo.Context, login string) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetPointInfoByLogin(ctx echo.Context, login string) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetTerritoryPointHistoryByLogin(ctx echo.Context, login string) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetOwnCurrentTimer(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) PauseOwnCurrentTimer(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) StartOwnCurrentTimer(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) RollAvailableWheelEffects(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetAvailableWheelEffectRollsCount(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetOwnAvailableWheelEffects(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetOwnWheelEffectHistory(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetWheelEffectHistoryByLogin(ctx echo.Context, login string) error {
	//TODO implement me
	panic("implement me")
}

func NewServer() Server {
	return Server{}
}
