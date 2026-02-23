package main

import (
	genauth "FGG-Service/api/generated/auth"
	gengames "FGG-Service/api/generated/games"
	genpoints "FGG-Service/api/generated/points"
	gentimers "FGG-Service/api/generated/timers"
	genusers "FGG-Service/api/generated/users"
	geneffects "FGG-Service/api/generated/wheel_effects"
	ctrlauth "FGG-Service/src/auth/controller"
	"FGG-Service/src/db_access"
	ctrlgames "FGG-Service/src/games/controller"
	ctrlpoints "FGG-Service/src/points/controller"
	ctrltimers "FGG-Service/src/timers/controller"
	ctrlusers "FGG-Service/src/users/controller"
	ctrleffects "FGG-Service/src/wheel-effects/controller"
	"embed"
	"io/fs"
	"net/http"
)

//go:embed index.html
var indexHTML []byte

//go:embed api/specification/mains
var scalarUI embed.FS

func main() {
	authController := new(ctrlauth.Controller)
	gamesController := new(ctrlgames.Controller)
	pointsController := new(ctrlpoints.Controller)
	timersController := new(ctrltimers.Controller)
	usersController := new(ctrlusers.Controller)
	wheelEffectsController := new(ctrleffects.Controller)

	authServer, err := genauth.NewServer(authController)

	if err != nil {
		panic(err)
	}

	gamesServer, err := gengames.NewServer(gamesController)

	if err != nil {
		panic(err)
	}

	pointsServer, err := genpoints.NewServer(pointsController)

	if err != nil {
		panic(err)
	}

	timersServer, err := gentimers.NewServer(timersController)

	if err != nil {
		panic(err)
	}

	usersServer, err := genusers.NewServer(usersController)

	if err != nil {
		panic(err)
	}

	wheelEffectsServer, err := geneffects.NewServer(wheelEffectsController)

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
	mux.HandleFunc("/scalar/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/scalar/" {
			_, _ = w.Write(indexHTML)
			return
		}
		http.StripPrefix("/scalar/", http.FileServer(http.FS(sub))).ServeHTTP(w, r)
	})

	mux.HandleFunc("/scalar", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/scalar/", http.StatusMovedPermanently)
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
