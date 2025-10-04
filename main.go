package main

import (
	"FGG-Service/api"
	"FGG-Service/controller"
	"embed"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:embed api/*
var swaggerUI embed.FS

func main() {
	server := controller.NewServer()
	e := echo.New()

	strictHandler := api.NewStrictHandler(server, []api.StrictMiddlewareFunc{})
	api.RegisterHandlers(e, strictHandler)

	fileSystem := http.FS(swaggerUI)
	fileServer := http.FileServer(fileSystem)
	httpHandler := http.StripPrefix("/swagger/", fileServer)
	e.GET("/swagger/*", echo.WrapHandler(httpHandler))

	err := e.Start("localhost:8080")

	if err != nil {
		panic(err)
	}
}
