package main

import (
	"FGG-Service/api"
	"FGG-Service/controller"
	"FGG-Service/db_access"
	"FGG-Service/timer_service"
	"embed"
	"net/http"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed index.html
//go:embed api/openapi.yaml
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
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}}))

	db_access.Init()
	defer db_access.Close()
	StartTimerFinisherScheduler()

	e.HideBanner = true
	err := e.Start(":8080")

	if err != nil {
		panic(err)
	}
}

func StartTimerFinisherScheduler() {
	scheduler, err := gocron.NewScheduler()

	if err != nil {
		panic(err)
	}

	_, err = scheduler.NewJob(
		gocron.DurationJob(1*time.Second),
		gocron.NewTask(timer_service.StopAllCompletedTimers),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)

	if err != nil {
		panic(err)
	}

	scheduler.Start()
}
