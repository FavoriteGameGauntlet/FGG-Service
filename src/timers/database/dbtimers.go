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
		CASE WHEN t.State IN (?, ?)
	    	THEN t.RemainingTimeInS
	    	ELSE t.DurationInS
		END AS RemainingTime
	FROM Timers t
	WHERE UserId = ?
		AND t.State != ?
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
	var lastActionDateString string
	var remainingTimeInS int
	err = row.Scan(
		&timer.Id,
		&timer.State,
		&durationInS,
		&lastActionDateString,
		&remainingTimeInS,
	)

	if err != nil {
		dbaccess.LogDbResult(queryName, timer, err)

		return
	}

	lastActionDate, err := dbaccess.ConvertToDate(lastActionDateString)

	if err != nil {
		dbaccess.LogDbResult(queryName, timer, err)

		return
	}

	timer.Duration = time.Duration(durationInS) * time.Second
	timer.LastActionDate = lastActionDate
	timer.RemainingTime = time.Duration(remainingTimeInS) * time.Second

	dbaccess.LogDbResult(queryName, timer, err)

	return
}

const CreateCurrentTimerQuery = `
	INSERT INTO Timers (UserId, GameId, DurationInS, RemainingTimeInS)
	VALUES (?, ?, ?, ?)
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
		State = ?,
		RemainingTimeInS = ?,
		LastActionDate = datetime('now', 'subsec')
	WHERE Id = ?
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
	WHERE t.State NOT IN (?, ?)
	    AND CASE t.State
			WHEN ? THEN t.RemainingTimeInS - (strftime('%s', 'now') - strftime('%s', t.LastActionDate))
			WHEN ? THEN t.RemainingTimeInS
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
