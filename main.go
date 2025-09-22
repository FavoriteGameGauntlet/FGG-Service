package main

import (
	"FavoriteGameGauntlet/api"
	"FavoriteGameGauntlet/controller"

	"embed"
	"net/http"

	"github.com/labstack/echo/v4"
)

// go:embed api/index.html
// go:embed api/openapi.yaml
var swaggerUI embed.FS

func main() {
	s := controller.NewServer()
	e := echo.New()

	api.RegisterHandlers(e, api.NewStrictHandler(s, []api.StrictMiddlewareFunc{}))

	e.GET("/swagger/*", echo.WrapHandler(http.StripPrefix("/swagger/", http.FileServer(http.FS(swaggerUI)))))

	e.Logger.Fatal(e.Start("localhost:8080"))
}
