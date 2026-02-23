package main

import (
	genauth "FGG-Service/api/generated/auth"
	gengames "FGG-Service/api/generated/games"
	genpoints "FGG-Service/api/generated/points"
	gentimers "FGG-Service/api/generated/timers"
	genusers "FGG-Service/api/generated/users"
	geneffects "FGG-Service/api/generated/wheel_effects"
	ctrlauth "FGG-Service/src/auth/controller"
	"FGG-Service/src/dbaccess"
	ctrlgames "FGG-Service/src/games/controller"
	ctrlpoints "FGG-Service/src/points/controller"
	ctrltimers "FGG-Service/src/timers/controller"
	ctrlusers "FGG-Service/src/users/controller"
	ctrleffects "FGG-Service/src/wheeleffects/controller"
	"embed"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed index.html
//go:embed api/specification
var scalarUI embed.FS

func main() {
	e := echo.New()

	authController := ctrlauth.NewController()
	genauth.RegisterHandlers(e, &authController)

	gamesController := ctrlgames.NewController()
	gengames.RegisterHandlers(e, &gamesController)

	pointsController := ctrlpoints.NewController()
	genpoints.RegisterHandlers(e, &pointsController)

	timersController := ctrltimers.NewController()
	gentimers.RegisterHandlers(e, &timersController)

	usersController := ctrlusers.NewController()
	genusers.RegisterHandlers(e, &usersController)

	effectsController := ctrleffects.NewController()
	geneffects.RegisterHandlers(e, &effectsController)

	addScalarRoutes(e)
	fixCORS(e)

	dbaccess.Init()
	defer dbaccess.Close()
	//StartTimerFinisherScheduler()

	e.HideBanner = true
	err := e.Start(":8080")

	if err != nil {
		panic(err)
	}

	defer func(e *echo.Echo) {
		_ = e.Close()
	}(e)
}

func addScalarRoutes(e *echo.Echo) {
	fileServer := http.FileServer(http.FS(scalarUI))
	e.GET("/api/specification/*", echo.WrapHandler(fileServer))
	e.GET("/scalar/*", echo.WrapHandler(http.StripPrefix("/scalar", fileServer)))
	e.GET("/scalar", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/scalar/")
	})
}

func fixCORS(e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}}))
}

//func StartTimerFinisherScheduler() {
//	scheduler, err := gocron.NewScheduler()
//
//	if err != nil {
//		panic(err)
//	}
//
//	_, err = scheduler.NewJob(
//		gocron.DurationJob(1*time.Second),
//		gocron.NewTask(timer_service.StopAllCompletedTimers),
//		gocron.WithSingletonMode(gocron.LimitModeReschedule),
//	)
//
//	if err != nil {
//		panic(err)
//	}
//
//	scheduler.Start()
//}
