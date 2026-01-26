package timer_service

import (
	"FGG-Service/common"
	"FGG-Service/db_access"
	"database/sql"
	"errors"
	"time"
)

const GetCurrentTimerQuery = `
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
	LIMIT 1
`

func GetCurrentTimerCommand(userId int) (timer common.Timer, err error) {
	row := db_access.QueryRow(
		GetCurrentTimerQuery,
		common.TimerActionStart,
		common.TimerActionPause,
		common.TimerActionStop,
		userId,
		common.TimerStateFinished,
	)

	var timerActionDateString *string
	var remainingTimeInS int
	err = row.Scan(
		&timer.Id,
		&timer.State,
		&timer.Duration,
		&timerActionDateString,
		&remainingTimeInS,
	)

	if errors.Is(err, sql.ErrNoRows) {
		err = common.NewCurrentTimerNotFoundError()
		return
	}

	if err != nil {
		return
	}

	var timerActionDate *time.Time

	if timerActionDateString != nil {
		var notNilDate time.Time
		notNilDate, err = time.Parse(db_access.ISO8601, *timerActionDateString)

		if err != nil {
			return
		}

		timerActionDate = &notNilDate
	}

	timer.TimerActionDate = timerActionDate
	timer.RemainingTime = time.Duration(remainingTimeInS) * time.Second

	return
}

const CreateCurrentTimerQuery = `
	INSERT INTO Timers (UserId, GameId, DurationInS)
	VALUES (?, ?, ?)
`

func CreateCurrentTimerCommand(userId int, gameId int) (timer common.Timer, err error) {
	_, err = db_access.Exec(
		CreateCurrentTimerQuery,
		userId,
		gameId,
		common.DefaultTimerDuration,
	)

	if err != nil {
		return
	}

	timer = common.Timer{
		Duration:      common.DefaultTimerDuration,
		State:         common.TimerStateCreated,
		RemainingTime: common.DefaultTimerDuration,
	}

	return
}

const ActTimerQuery = `
	INSERT INTO TimerActions (TimerId, Action, RemainingTimeInS)
	VALUES (?, ?, ?)
`

func ActTimerCommand(timerId int, actionType common.TimerActionType, remainingTime time.Duration) error {
	_, err := db_access.Exec(
		ActTimerQuery,
		timerId,
		actionType,
		remainingTime.Seconds(),
	)

	return err
}

const GetCompletedTimerUsersQuery = `
	SELECT DISTINCT t.UserId
	FROM Timers t
		LEFT JOIN TimerActions ta ON t.Id = ta.TimerId
	WHERE t.State NOT IN (?, ?)
	    AND CASE ta.Action
			WHEN ? THEN ta.RemainingTimeInS - (strftime('%s', 'now') - strftime('%s', ta.CreateDate))
			WHEN ? THEN ta.RemainingTimeInS
			ELSE t.DurationInS
		END <= 0
`

func GetCompletedTimerUsersCommand() (userIds []int, err error) {
	rows, err := db_access.Query(
		GetCompletedTimerUsersQuery,
		common.TimerStateCreated,
		common.TimerStateFinished,
		common.TimerStateRunning,
		common.TimerStatePaused,
	)

	if err != nil {
		return
	}

	for rows.Next() {
		var userId int
		err = rows.Scan(&userId)

		if err != nil {
			continue
		}

		userIds = append(userIds, userId)
	}

	_ = rows.Close()
	err = nil
	return
}
