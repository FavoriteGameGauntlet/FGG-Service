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
	CheckIfCurrentGameExistsCommand = `
		SELECT CASE
        	WHEN EXISTS (
				SELECT 1
				FROM GameHistory gh
					INNER JOIN Games g ON gh.GameId = g.Id
				WHERE gh.UserId = ?
					AND gh.State NOT IN (?, ?))
         	THEN true
         	ELSE false
       	END AS 'DoesExist'`
	ActCurrentTimerCommand = `
		INSERT INTO TimerActions (TimerId, Action, RemainingTimeInS)
		VALUES (?, ?, ?)`
	GetCurrentTimerCommand = `
		SELECT
			t.Id,
			t.State,
			t.DurationInS,
			ta.CreateDate,
			CASE ta.Action
				WHEN ? THEN ta.RemainingTimeInS - (strftime('%s', 'now') - strftime('%s', ta.CreateDate))
				WHEN ? THEN ta.RemainingTimeInS
				WHEN ? THEN 0
				ELSE t.DurationInS
			END AS RemainingTimeInS
		FROM Timers t
			LEFT JOIN TimerActions ta ON t.Id = ta.TimerId
		WHERE UserId = ?
			AND t.State != ?
		ORDER BY ta.CreateDate DESC
		LIMIT 1`
)

func AddUnplayedGames(userId int, games common.UnplayedGames) error {
	numberOfGames := len(games)
	errorCount := 0

	var err error
	for _, game := range games {
		err = CreateUnplayedGame(userId, &game)

		if err != nil {
			errorCount++
		}
	}

	if errorCount == numberOfGames {
		return err
	}

	return nil
}

func CreateUnplayedGame(userId int, unplayedGame *common.UnplayedGame) error {
	doesExist, err := DoesUnplayedGameExistCommand(userId, unplayedGame.Name)

	if err != nil {
		return err
	}

	if doesExist {
		return common.NewUnplayedGameAlreadyExistsError(unplayedGame.Name)
	}

	doesExist, err = DoesGameExistCommand(unplayedGame.Name)

	if err != nil {
		return err
	}

	game := common.Game{}

	if doesExist {
		game, err = GetGameCommand(unplayedGame.Name)
	} else {
		game, err = CreateGame(unplayedGame)
	}

	if err != nil {
		return err
	}

	err = CreateUnplayedGameCommand(userId, game.Id)

	return err
}

func GetUnplayedGames(userId int) (common.UnplayedGames, error) {
	return GetUnplayedGamesCommand(userId)
}

func CreateGame(unplayedGame *common.UnplayedGame) (game common.Game, err error) {
	err = CreateGameCommand(unplayedGame.Name)

	if err != nil {
		return
	}

	game, err = GetGameCommand(unplayedGame.Name)

	return
}

func GetCurrentGame(userId int) (game common.Game, err error) {
	game, err = GetCurrentGameCommand(userId)

	if err != nil {
		return
	}

	timeSpent, err := GetGameTimeSpentCommand(userId, game.Id)

	if err != nil {
		return
	}

	game.TimeSpent = timeSpent

	return
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
		Type:             timerAction,
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

func GetGameHistory(userId int) (games common.Games, err error) {
	games, err = GetEndedGamesCommand(userId)

	return
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
	game.State = common.GameStateStarted

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
