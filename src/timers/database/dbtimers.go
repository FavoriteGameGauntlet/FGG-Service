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

const GetCurrentTimerQuery = `SELECT * FROM get_current_timer($1)`

func (db *Database) GetCurrentTimerCommand(userId int) (timer typetimers.Timer, err error) {
	queryName := "GetCurrentTimerQuery"
	row := dbaccess.QueryRow(queryName, GetCurrentTimerQuery, userId)

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

const CreateCurrentTimerQuery = `SELECT create_current_timer($1, $2, $3)`

func (db *Database) CreateCurrentTimerCommand(userId int, gameId int, durationInS int) error {
	queryName := "CreateCurrentTimerQuery"
	_, err := dbaccess.Exec(
		queryName,
		CreateCurrentTimerQuery,
		userId,
		gameId,
		durationInS,
	)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const ActTimerQuery = `SELECT act_timer($1, $2, $3)`

func (db *Database) ActTimerCommand(
	timerId int,
	timerState typetimers.TimerStateType,
	remainingTime time.Duration) error {

	queryName := "ActTimerQuery"
	_, err := dbaccess.Exec(
		queryName,
		ActTimerQuery,
		timerId,
		timerState,
		int(remainingTime.Seconds()),
	)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const GetCompletedTimerUsersQuery = `SELECT * FROM get_completed_timer_users()`

func (db *Database) GetCompletedTimerUsersCommand() (userIds []int, err error) {
	queryName := "GetCompletedTimerUsersQuery"
	rows, err := dbaccess.Query(queryName, GetCompletedTimerUsersQuery)

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
