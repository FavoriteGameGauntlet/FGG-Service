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

const GetAllSystemParametersQuery = `
	SELECT Id, Name, Value
	FROM SystemParameters
`

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

const GetSystemParameterQuery = `
	SELECT Id, Name, Value
	FROM SystemParameters
	WHERE Name = ?
`

func (db *Database) GetSystemParameterCommand(name string) (parameter typesysparams.SystemParameter, err error) {
	queryName := "GetSystemParameterQuery"
	row := dbaccess.QueryRow(queryName, GetSystemParameterQuery, name)

	err = row.Scan(&parameter.Id, &parameter.Name, &parameter.Value)

	dbaccess.LogDbResult(queryName, parameter, err)

	return
}

const ChangeSystemParameterValueQuery = `
	UPDATE SystemParameters
	SET Value = ?
	WHERE Name = ?
`

func (db *Database) ChangeSystemParameterValueCommand(name string, value string) error {
	queryName := "ChangeSystemParameterValueQuery"
	_, err := dbaccess.Exec(queryName, ChangeSystemParameterValueQuery, value, name)

	dbaccess.LogDbResult(queryName, nil, err)

	return err
}
