package dbsysparams

import (
	"FGG-Service/src/dbaccess"
	"FGG-Service/src/sysparams/types"
)

type IDatabase interface {
	GetAllSystemParametersCommand() (parameters []typesysparams.SystemParameter, err error)
	GetSystemParameterCommand(name string) (parameter typesysparams.SystemParameter, err error)
	ChangeSystemParameterValueCommand(name string, value string) error
}

type Database struct {
}

const GetAllSystemParametersQuery = `SELECT * FROM get_all_system_parameters()`

func (db *Database) GetAllSystemParametersCommand() (parameters []typesysparams.SystemParameter, err error) {
	queryName := "GetAllSystemParametersQuery"
	rows, err := dbaccess.Query(queryName, GetAllSystemParametersQuery)

	if err != nil {
		return
	}

	for rows.Next() {
		parameter := typesysparams.SystemParameter{}
		err = rows.Scan(&parameter.Id, &parameter.Name, &parameter.Value)

		if err != nil {
			_ = rows.Close()
			return
		}

		parameters = append(parameters, parameter)
	}

	dbaccess.LogDbResult(queryName, parameters, err)

	_ = rows.Close()
	return
}

const GetSystemParameterQuery = `SELECT * FROM get_system_parameter($1::text)`

func (db *Database) GetSystemParameterCommand(name string) (parameter typesysparams.SystemParameter, err error) {
	queryName := "GetSystemParameterQuery"
	row := dbaccess.QueryRow(queryName, GetSystemParameterQuery, name)

	err = row.Scan(&parameter.Id, &parameter.Name, &parameter.Value)

	dbaccess.LogDbResult(queryName, parameter, err)

	return
}

const ChangeSystemParameterValueQuery = `SELECT change_system_parameter_value($1::text, $2::text)`

func (db *Database) ChangeSystemParameterValueCommand(name string, value string) error {
	queryName := "ChangeSystemParameterValueQuery"
	_, err := dbaccess.Exec(queryName, ChangeSystemParameterValueQuery, name, value)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}
