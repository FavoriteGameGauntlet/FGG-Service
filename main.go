package main

import (
	"FavoriteGameGauntlet/api"
	"FavoriteGameGauntlet/controller"

	"github.com/labstack/echo/v4"
)

func main() {
	server := controller.NewServer()

	e := echo.New()

	api.RegisterHandlers(e, api.NewStrictHandler(
		server,
		[]api.StrictMiddlewareFunc{},
	))

	e.Start("127.0.0.1:8080")
}
