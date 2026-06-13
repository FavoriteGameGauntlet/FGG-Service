package dbauth

import (
	"FGG-Service/src/auth/types"
	"FGG-Service/src/dbaccess"
)

type Database struct {
}

const GetUserByLoginQuery = `SELECT * FROM get_user_by_login($1)`

func (db *Database) GetUserByLoginCommand(userLogin string) (user typeauth.User, err error) {
	queryName := "GetUserByLoginQuery"
	row := dbaccess.QueryRow(queryName, GetUserByLoginQuery, userLogin)

	err = row.Scan(&user.Id, &user.Login, &user.DisplayName, &user.Email)

	dbaccess.LogDbResult(queryName, user, err)

	return
}

const GetUserByIdQuery = `SELECT * FROM get_user_by_id($1)`

func (db *Database) GetUserByIdCommand(userId int) (user typeauth.User, err error) {
	queryName := "GetUserByIdQuery"
	row := dbaccess.QueryRow(queryName, GetUserByIdQuery, userId)

	err = row.Scan(&user.Id, &user.Login, &user.DisplayName, &user.Email)

	dbaccess.LogDbResult(queryName, user, err)

	return
}

const GetUserByEmailQuery = `SELECT * FROM get_user_by_email($1)`

func (db *Database) GetUserByEmailCommand(userEmail string) (user typeauth.User, err error) {
	queryName := "GetUserByEmailQuery"
	row := dbaccess.QueryRow(queryName, GetUserByEmailQuery, userEmail)

	err = row.Scan(&user.Id, &user.Login, &user.DisplayName, &user.Email)

	dbaccess.LogDbResult(queryName, user, err)

	return
}

const GetUserByLoginAndPasswordQuery = `SELECT * FROM get_user_by_login_and_password($1, $2)`

func (db *Database) GetUserByLoginAndPasswordCommand(loginUser typeauth.LoginUser) (user typeauth.User, err error) {
	queryName := "GetUserByLoginAndPasswordQuery"
	row := dbaccess.QueryRow(queryName, GetUserByLoginAndPasswordQuery, loginUser.Login, loginUser.Password)

	err = row.Scan(&user.Id, &user.Login, &user.DisplayName, &user.Email)

	dbaccess.LogDbResult(queryName, user, err)

	return
}

const CreateUserQuery = `SELECT create_user($1, $2, $3)`

func (db *Database) CreateUserCommand(signupUser typeauth.SignupUser) error {
	queryName := "CreateUserQuery"
	_, err := dbaccess.Exec(queryName, CreateUserQuery, signupUser.Login, signupUser.Email, signupUser.Password)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const CreateUserStatsQuery = `SELECT create_user_stats($1)`

func (db *Database) CreateUserStatsCommand(login string) error {
	queryName := "CreateUserStatsQuery"
	_, err := dbaccess.Exec(queryName, CreateUserStatsQuery, login)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}

const GetUserSessionByIdQuery = `SELECT * FROM get_user_session_by_id($1)`

func (db *Database) GetUserSessionByIdCommand(sessionId string) (userSession typeauth.UserSession, err error) {
	queryName := "GetUserSessionByIdQuery"
	row := dbaccess.QueryRow(queryName, GetUserSessionByIdQuery, sessionId)

	err = row.Scan(&userSession.Id, &userSession.UserId)

	dbaccess.LogDbResult(queryName, userSession, err)

	return
}

const CreateUserSessionQuery = `SELECT * FROM create_user_session($1)`

func (db *Database) CreateUserSessionCommand(userId int) (userSession typeauth.UserSession, err error) {
	queryName := "CreateUserSessionQuery"
	row := dbaccess.QueryRow(queryName, CreateUserSessionQuery, userId)

	err = row.Scan(&userSession.Id, &userSession.UserId)

	dbaccess.LogDbResult(queryName, userSession, err)

	return
}

const DeleteUserSessionQuery = `SELECT delete_user_session($1)`

func (db *Database) DeleteUserSessionCommand(sessionId string) error {
	queryName := "DeleteUserSessionQuery"
	_, err := dbaccess.Exec(queryName, DeleteUserSessionQuery, sessionId)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}
