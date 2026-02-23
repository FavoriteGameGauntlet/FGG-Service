package dbtimers

import (
	"FGG-Service/src/common"
	"FGG-Service/src/db_access"
	"time"
)

type Database struct {
}

const GetCurrentTimerQuery = `
	SELECT
		t.Id,
		t.State,
		t.DurationInS,
		t.LastActionDate,
		CASE t.State
			WHEN ? THEN t.RemainingTimeInS - (strftime('%s', 'now') - strftime('%s', t.LastActionDate))
			WHEN ? THEN t.RemainingTimeInS
			ELSE t.DurationInS
		END AS RemainingTimeInS
	FROM Timers t
	WHERE UserId = ?
		AND t.State != ?
`

func (db *Database) GetCurrentTimerCommand(userId int) (timer common.Timer, err error) {
	row := db_access.QueryRow(
		GetCurrentTimerQuery,
		common.TimerStateRunning,
		common.TimerStatePaused,
		userId,
		common.TimerStateFinished,
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
		return
	}

	lastActionDate, err := common.ConvertToDate(lastActionDateString)

	if err != nil {
		return
	}

	timer.Duration = time.Duration(durationInS) * time.Second
	timer.LastActionDate = lastActionDate
	timer.RemainingTime = time.Duration(remainingTimeInS) * time.Second

	return
}

const CreateCurrentTimerQuery = `
	INSERT INTO Timers (UserId, GameId, DurationInS, RemainingTimeInS)
	VALUES (?, ?, ?, ?)
`

func (db *Database) CreateCurrentTimerCommand(userId int, gameId int) error {
	_, err := db_access.Exec(
		CreateCurrentTimerQuery,
		userId,
		gameId,
		common.DefaultTimerDurationInS,
		common.DefaultTimerDurationInS,
	)

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
	timerState common.TimerStateType,
	remainingTime time.Duration) error {

	_, err := db_access.Exec(
		ActTimerQuery,
		timerState,
		remainingTime.Seconds(),
		timerId,
	)

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
