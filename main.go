package main

import (
	auth_server "FGG-Service/api/generated/auth"
	games_server "FGG-Service/api/generated/games"
	points_server "FGG-Service/api/generated/points"
	timers_server "FGG-Service/api/generated/timers"
	users_server "FGG-Service/api/generated/users"
	wheel_effects_server "FGG-Service/api/generated/wheel_effects"
	"FGG-Service/src/auth/controller"
	"FGG-Service/src/db_access"
	"FGG-Service/src/games/controller"
	"FGG-Service/src/points/controller"
	"FGG-Service/src/timers/controller"
	"FGG-Service/src/users/controller"
	"FGG-Service/src/wheel-effects/controller"
	"embed"
	"io/fs"
	"net/http"
)

//go:embed index.html
var indexHTML []byte

//go:embed api/specification/mains
var scalarUI embed.FS

func main() {
	authController := new(auth.Controller)
	gamesController := new(games.Controller)
	pointsController := new(points.Controller)
	timersController := new(timers.Controller)
	usersController := new(users.Controller)
	wheelEffectsController := new(wheel_effects.Controller)

	authServer, err := auth_server.NewServer(authController)

	if err != nil {
		panic(err)
	}

	gamesServer, err := games_server.NewServer(gamesController)

	if err != nil {
		panic(err)
	}
	pointsServer, err := points_server.NewServer(pointsController)

	if err != nil {
		panic(err)
	}
	timersServer, err := timers_server.NewServer(timersController)

	if err != nil {
		panic(err)
	}
	usersServer, err := users_server.NewServer(usersController)

	if err != nil {
		panic(err)
	}
	wheelEffectsServer, err := wheel_effects_server.NewServer(wheelEffectsController)

	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mux.Handle("/auth", authServer)
	mux.Handle("/games", gamesServer)
	mux.Handle("/points", pointsServer)
	mux.Handle("/timers", timersServer)
	mux.Handle("/users", usersServer)
	mux.Handle("/wheel-effects", wheelEffectsServer)

	sub, _ := fs.Sub(scalarUI, "api/specification/mains")

	mux.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/swagger/" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(indexHTML)
			return
		}
		http.StripPrefix("/swagger/", http.FileServer(http.FS(sub))).ServeHTTP(w, r)
	})

	mux.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})

	db_access.Init()
	defer db_access.Close()
	//StartTimerFinisherScheduler()

	err = http.ListenAndServe(":8080", mux)

	if err != nil {
		panic(err)
	}
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
