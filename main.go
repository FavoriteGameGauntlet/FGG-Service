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
	s := controller.NewServer()
	e := echo.New()

	api.RegisterHandlers(e, api.NewStrictHandler(s, []api.StrictMiddlewareFunc{}))

	fs := http.FS(swaggerUI)
	fileServer := http.FileServer(fs)

	e.GET("/swagger/*", echo.WrapHandler(http.StripPrefix("/swagger/", fileServer)))

	e.Start("127.0.0.1:8080")
}
