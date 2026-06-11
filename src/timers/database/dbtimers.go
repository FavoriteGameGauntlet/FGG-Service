package dbtimers

import (
	"FGG-Service/src/dbaccess"
	"FGG-Service/src/timers/types"
	"time"
)

type IDatabase interface {
	GetCurrentTimerCommand(userId int) (timer typetimers.Timer, err error)
	CreateCurrentTimerCommand(userId int, gameId int, durationInS int) error
	ActTimerCommand(timerId int, timerState typetimers.TimerStateType, remainingTime time.Duration) error
	GetCompletedTimerUsersCommand() (userIds []int, err error)
}

type Database struct {
}

const GetCurrentTimerQuery = `
	SELECT
		t.Id,
		t.State,
		t.DurationInS,
		t.LastActionDate,
		CASE WHEN t.State IN ($1, $2)
	    	THEN t.RemainingTimeInS
	    	ELSE t.DurationInS
		END AS RemainingTime
	FROM Timers t
	WHERE UserId = $3
		AND t.State != $4
`

func (db *Database) GetCurrentTimerCommand(userId int) (timer typetimers.Timer, err error) {
	queryName := "GetCurrentTimerQuery"
	row := dbaccess.QueryRow(
		queryName,
		GetCurrentTimerQuery,
		typetimers.TimerStateRunning,
		typetimers.TimerStatePaused,
		userId,
		typetimers.TimerStateFinished,
	)

	var durationInS int
	var remainingTimeInS int
	err = row.Scan(
		&timer.Id,
		&timer.State,
		&durationInS,
		&timer.LastActionDate,
		&remainingTimeInS,
	)

	if err != nil {
		dbaccess.LogDbResult(queryName, timer, err)

		return
	}

	timer.Duration = time.Duration(durationInS) * time.Second
	timer.RemainingTime = time.Duration(remainingTimeInS) * time.Second

	dbaccess.LogDbResult(queryName, timer, err)

	return
}

const CreateCurrentTimerQuery = `
	INSERT INTO Timers (UserId, GameId, DurationInS, RemainingTimeInS)
	VALUES ($1, $2, $3, $4)
`

func (db *Database) CreateCurrentTimerCommand(userId int, gameId int, durationInS int) error {
	queryName := "CreateCurrentTimerQuery"
	_, err := dbaccess.Exec(
		queryName,
		CreateCurrentTimerQuery,
		userId,
		gameId,
		durationInS,
		durationInS,
	)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const ActTimerQuery = `
	UPDATE Timers
	SET
		State = $1,
		RemainingTimeInS = $2,
		LastActionDate = NOW()
	WHERE Id = $3
`

func (db *Database) ActTimerCommand(
	timerId int,
	timerState typetimers.TimerStateType,
	remainingTime time.Duration) error {

	queryName := "ActTimerQuery"
	_, err := dbaccess.Exec(
		queryName,
		ActTimerQuery,
		timerState,
		int(remainingTime.Seconds()),
		timerId,
	)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const GetCompletedTimerUsersQuery = `
	SELECT DISTINCT t.UserId
	FROM Timers t
	WHERE t.State NOT IN ($1, $2)
	    AND CASE t.State
			WHEN $3 THEN t.RemainingTimeInS - CAST(EXTRACT(EPOCH FROM (NOW() - t.LastActionDate)) AS INTEGER)
			WHEN $4 THEN t.RemainingTimeInS
			ELSE t.DurationInS
		END <= 0
`

func (db *Database) GetCompletedTimerUsersCommand() (userIds []int, err error) {
	queryName := "GetCompletedTimerUsersQuery"
	rows, err := dbaccess.Query(
		queryName,
		GetCompletedTimerUsersQuery,
		typetimers.TimerStateCreated,
		typetimers.TimerStateFinished,
		typetimers.TimerStateRunning,
		typetimers.TimerStatePaused,
	)

	if err != nil {
		return
	}

	for rows.Next() {
		var userId int
		err = rows.Scan(&userId)

		if err != nil {
			dbaccess.LogDbResult(queryName, userIds, err)

			continue
		}

		userIds = append(userIds, userId)
	}

	_ = rows.Close()
	err = nil
	return
}
