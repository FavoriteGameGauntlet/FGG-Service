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
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed index.html
//go:embed api/specification
var scalarUI embed.FS

func main() {
	e := echo.New()
	e.HideBanner = true

	registerHandlers(e)
	addScalarRoutes(e)
	fixCORS(e)

	dbCloseFunc := dbaccess.Init()
	defer dbCloseFunc()

	f := createFileAndStartLogger()
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	startLogScheduler()

	defer func(e *echo.Echo) {
		_ = e.Close()
	}(e)

	if err := e.Start(":8080"); err != nil {
		panic(err)
	}
}

func registerHandlers(e *echo.Echo) {
	genauth.RegisterHandlers(e, ctrlauth.NewController())
	gengames.RegisterHandlers(e, ctrlgames.NewController())
	genpoints.RegisterHandlers(e, ctrlpoints.NewController())
	gentimers.RegisterHandlers(e, ctrltimers.NewController())
	genusers.RegisterHandlers(e, ctrlusers.NewController())
	geneffects.RegisterHandlers(e, ctrleffects.NewController())
}

func createFileAndStartLogger() *os.File {
	logsDir := getLogsDir()
	filename := filepath.Join(logsDir, time.Now().Format("2006-01-02")+".txt")
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		panic(err)
	}

	handler := slog.NewJSONHandler(file, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(handler)

	slog.SetDefault(logger)

	return file
}

func getLogsDir() string {
	root := os.Getenv("APP_ROOT")

	if root != "" {
		return filepath.Join(root, "logs")
	}

	return filepath.Join("logs")
}

func startLogScheduler() {
	scheduler, err := gocron.NewScheduler()

	if err != nil {
		panic(err)
	}

	var currentFile *os.File

	_, err = scheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(0, 0, 0))),
		gocron.NewTask(func() {
			if currentFile != nil {
				_ = currentFile.Close()
			}

			currentFile = createFileAndStartLogger()
		}),
	)

	if err != nil {
		panic(err)
	}

	scheduler.Start()
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
