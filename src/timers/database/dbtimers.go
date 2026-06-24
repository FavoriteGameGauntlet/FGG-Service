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

var getCurrentTimerQuery = dbaccess.Query{Name: "GetCurrentTimerQuery", SQL: `SELECT * FROM get_current_timer($1::integer)`}

func (db *Database) GetCurrentTimerCommand(userId int) (timer typetimers.Timer, err error) {
	row := dbaccess.QueryRow(getCurrentTimerQuery, userId)

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
		dbaccess.LogDbResult(getCurrentTimerQuery, timer, err)

		return
	}

	timer.Duration = time.Duration(durationInS) * time.Second
	timer.RemainingTime = time.Duration(remainingTimeInS) * time.Second

	dbaccess.LogDbResult(getCurrentTimerQuery, timer, err)

	return
}

var createCurrentTimerQuery = dbaccess.Query{Name: "CreateCurrentTimerQuery", SQL: `SELECT create_current_timer($1::integer, $2::integer, $3::integer)`}

func (db *Database) CreateCurrentTimerCommand(userId int, gameId int, durationInS int) error {
	_, err := dbaccess.Exec(
		createCurrentTimerQuery,
		userId,
		gameId,
		durationInS,
	)

	dbaccess.LogDbResult(createCurrentTimerQuery, nil, err)

	return err
}

var actTimerQuery = dbaccess.Query{Name: "ActTimerQuery", SQL: `SELECT act_timer($1::integer, $2::text, $3::integer)`}

func (db *Database) ActTimerCommand(
	timerId int,
	timerState typetimers.TimerStateType,
	remainingTime time.Duration) error {

	_, err := dbaccess.Exec(
		actTimerQuery,
		timerId,
		timerState,
		int(remainingTime.Seconds()),
	)

	dbaccess.LogDbResult(actTimerQuery, nil, err)

	return err
}

var getCompletedTimerUsersQuery = dbaccess.Query{Name: "GetCompletedTimerUsersQuery", SQL: `SELECT * FROM get_completed_timer_users()`, IsSilent: true}

func (db *Database) GetCompletedTimerUsersCommand() (userIds []int, err error) {
	rows, err := dbaccess.QueryRows(getCompletedTimerUsersQuery)

	if err != nil {
		return
	}

	for rows.Next() {
		var userId int
		err = rows.Scan(&userId)

		if err != nil {
			dbaccess.LogDbResult(getCompletedTimerUsersQuery, userIds, err)

			continue
		}

		userIds = append(userIds, userId)
	}

	_ = rows.Close()
	err = nil
	return
}
