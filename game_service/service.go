package game_service

import (
	"FGG-Service/common"
	"FGG-Service/db_access"
	"database/sql"
	"errors"
	"math/rand"
	"time"
)

const (
	CheckIfGameExistsCommand = `
		SELECT CASE
        	WHEN EXISTS (
				SELECT 1
				FROM Games
				WHERE Name = $gameName)
         	THEN true
         	ELSE false
       	END AS 'DoesExist'`
	CheckIfUnplayedGameExistsCommand = `
		SELECT CASE
        	WHEN EXISTS (
				SELECT 1
				FROM UnplayedGames ug
					INNER JOIN Games g ON ug.GameId = g.Id
				WHERE ug.UserId = $userId
					AND g.Name = $gameName)
         	THEN true
         	ELSE false
       	END AS 'DoesExist'`
	CheckIfCurrentGameExistsCommand = `
		SELECT CASE
        	WHEN EXISTS (
				SELECT 1
				FROM GameHistory gh
					INNER JOIN Games g ON gh.GameId = g.Id
				WHERE gh.UserId = $userId
					AND gh.State NOT IN ($finishedGameState, $cancelledGameState))
         	THEN true
         	ELSE false
       	END AS 'DoesExist'`
	ActCurrentTimerCommand = `
		INSERT INTO TimerActions (TimerId, Action, RemainingTimeInS)
		VALUES ($timerId, $timerAction, $remainingTimeInS)`
	GetCurrentTimerCommand = `
		SELECT
			t.Id,
			t.State,
			t.DurationInS,
			ta.CreateDate,
			CASE ta.Action
				WHEN $startTimerAction THEN ta.RemainingTimeInS - (strftime('%s', 'now') - strftime('%s', ta.CreateDate))
				WHEN $pauseTimerAction THEN ta.RemainingTimeInS
				WHEN $finishTimerAction THEN 0
				ELSE t.DurationInS
			END AS RemainingTimeInS
		FROM Timers t
			LEFT JOIN TimerActions ta ON t.Id = ta.TimerId
		WHERE UserId = $userId
			AND t.State != $finishedTimerState
		ORDER BY ta.CreateDate DESC
		LIMIT 1`
)

func AddUnplayedGames(userId int, games common.UnplayedGames) error {
	numberOfGames := len(games)
	errorCount := 0

	var err error
	for _, game := range games {
		err = AddUnplayedGame(userId, &game)

		if err != nil {
			errorCount++
		}
	}

	if errorCount == numberOfGames {
		return err
	}

	return nil
}

func AddUnplayedGame(userId int, unplayedGame *common.UnplayedGame) error {
	doesExist, err := CheckIfUnplayedGameExists(
		userId,
		unplayedGame.Name,
	)

	if err != nil || doesExist {
		return err
	}

	game, err := AddOrGetGame(unplayedGame)

	if err != nil {
		return err
	}

	err = CreateUnplayedGameCommand(userId, game.Id)

	return err
}

func GetUnplayedGames(userId int) (common.UnplayedGames, error) {
	return GetUnplayedGamesCommand(userId)
}

func CheckIfUnplayedGameExists(userId int, gameName string) (bool, error) {
	row := db_access.QueryRow(CheckIfUnplayedGameExistsCommand, userId, gameName)

	var doesExist bool
	err := row.Scan(&doesExist)

	if err != nil {
		return doesExist, err
	}

	return doesExist, nil
}

func AddOrGetGame(unplayedGame *common.UnplayedGame) (*common.Game, error) {
	doesExist, err := CheckIfGameExists(unplayedGame.Name)

	if err != nil {
		return nil, err
	}

	if doesExist {
		var game common.Game
		game, err = GetGameCommand(unplayedGame.Name)

		if err != nil {
			return nil, err
		}

		return &game, nil
	}

	err = CreateGameCommand(unplayedGame.Name, unplayedGame.Link)

	if err != nil {
		return nil, err
	}

	var game common.Game
	game, err = GetGameCommand(unplayedGame.Name)

	if err != nil {
		return nil, err
	}

	return &game, nil
}

func CheckIfGameExists(gameName string) (bool, error) {
	row := db_access.QueryRow(CheckIfGameExistsCommand, gameName)

	var doesExist bool
	err := row.Scan(&doesExist)

	if err != nil {
		return doesExist, err
	}

	return doesExist, nil
}

func GetCurrentGame(userId int) (common.Game, error) {
	game, err := GetCurrentGameCommand(userId)

	if err != nil {
		return common.Game{}, err
	}

	secondsSpent, err := GetGameSecondsSpentCommand(userId, game.Id)

	if err != nil {
		return common.Game{}, err
	}

	game.TimeSpent = time.Duration(secondsSpent) * time.Second

	return game, nil
}

func CheckIfCurrentGameExists(userId int) (bool, error) {
	row := db_access.QueryRow(CheckIfCurrentGameExistsCommand, userId, common.GameStateFinished, common.GameStateCancelled)

	var doesExist bool
	err := row.Scan(&doesExist)

	if err != nil {
		return doesExist, err
	}

	return doesExist, nil
}

func CancelCurrentGame(userId int) error {
	game, err := GetCurrentGame(userId)

	if err != nil {
		return err
	}

	_, err = StopCurrentTimer(userId)

	if err != nil {
		return err
	}

	err = CancelCurrentGameCommand(userId, game.Id)

	if err != nil {
		return err
	}

	return nil
}

func FinishCurrentGame(userId int) error {
	game, err := GetCurrentGameCommand(userId)

	if err != nil {
		return err
	}

	if game.TimeSpent < 1*time.Second {
		return common.NewCompletedTimersNotFoundError()
	}

	_, err = StopCurrentTimer(userId)

	if err != nil {
		return err
	}

	err = FinishCurrentGameCommand(userId, game.Id)

	if err != nil {
		return err
	}

	return nil
}

func StopCurrentTimer(userId int) (*common.TimerAction, error) {
	timer, err := GetCurrentTimer(userId)

	if err != nil {
		return nil, err
	}

	if timer == nil {
		return nil, nil
	}

	if timer.State == common.TimerStateFinished {
		return nil, nil
	}

	timerAction := common.TimerActionStop
	remainingTimerInS := 0

	_, err = db_access.Exec(
		ActCurrentTimerCommand,
		timer.Id,
		timerAction,
		remainingTimerInS,
	)

	if err != nil {
		return nil, err
	}

	return &common.TimerAction{
		Action:           timerAction,
		RemainingTimeInS: remainingTimerInS,
	}, nil
}

func GetCurrentTimer(userId int) (*common.Timer, error) {
	row := db_access.QueryRow(
		GetCurrentTimerCommand,
		common.TimerActionStart,
		common.TimerActionPause,
		common.TimerActionStop,
		userId,
		common.TimerStateFinished,
	)

	timer := common.Timer{}
	var timerActionDateString *string
	err := row.Scan(
		&timer.Id,
		&timer.State,
		&timer.DurationInS,
		&timerActionDateString,
		&timer.RemainingTimeInS,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var timerActionDate *time.Time

	if timerActionDateString != nil {
		var notNilDate time.Time
		notNilDate, err = time.Parse(db_access.ISO8601, *timerActionDateString)

		if err != nil {
			return nil, err
		}

		timerActionDate = &notNilDate
	}

	timer.TimerActionDate = timerActionDate

	return &timer, nil
}

func GetGameHistory(userId int) (*common.Games, error) {
	rows, err := db_access.Query(GetGameHistoryCommand, userId)

	if err != nil {
		return nil, err
	}

	gameCount := 0
	errorCount := 0
	games := common.Games{}
	for rows.Next() {
		gameCount++

		game := common.Game{}
		var spentSeconds *int
		var finishDateString *string
		err = rows.Scan(&game.Id, &game.Name, &game.State, &game.Link, &spentSeconds, &finishDateString)

		if err != nil {
			errorCount++
			continue
		}

		if spentSeconds != nil {
			game.TimeSpent = time.Duration(*spentSeconds) * time.Second
		}

		var finishDate *time.Time

		if finishDateString != nil {
			var notNilFinishDate time.Time
			notNilFinishDate, err = time.Parse(db_access.ISO8601, *finishDateString)

			if err != nil {
				errorCount++
				continue
			}

			finishDate = &notNilFinishDate
		}

		game.FinishDate = finishDate

		games = append(games, game)
	}

	if errorCount > 0 && errorCount == gameCount {
		return nil, err
	}

	return &games, nil
}

func MakeGameRoll(userId int) (common.Game, error) {
	game := common.Game{}

	doesExist, err := CheckIfCurrentGameExists(userId)

	if err != nil {
		return game, err
	}

	if doesExist {
		return game, common.NewCurrentGameAlreadyExistsError()
	}

	unplayedGames, err := GetUnplayedGamesCommand(userId)

	if err != nil {
		return game, err
	}

	if unplayedGames == nil || len(unplayedGames) < common.MinimumNumberOfUnplayedGames {
		return game, common.NewUnplayedGamesNotFoundError()
	}

	randomNumber := rand.Intn(len(unplayedGames))
	randomUnplayedGame := unplayedGames[randomNumber]

	err = CreateCurrentGameCommand(userId, randomUnplayedGame.GameId)

	if err != nil {
		return game, err
	}

	err = DeleteUnplayedGameCommand(userId, randomUnplayedGame.GameId)

	if err != nil {
		return game, err
	}

	game.Id = randomUnplayedGame.GameId
	game.Name = randomUnplayedGame.Name
	game.Link = randomUnplayedGame.Link
	game.State = common.GameStateStarted

	return game, nil
}
