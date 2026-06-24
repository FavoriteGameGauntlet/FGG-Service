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

var getAllSystemParametersQuery = dbaccess.Query{Name: "GetAllSystemParametersQuery", SQL: `SELECT * FROM get_all_system_parameters()`}

func (db *Database) GetAllSystemParametersCommand() (parameters []typesysparams.SystemParameter, err error) {
	rows, err := dbaccess.QueryRows(getAllSystemParametersQuery)

	if err != nil {
		return
	}

	for rows.Next() {
		parameter := typesysparams.SystemParameter{}
		err = rows.Scan(&parameter.Id, &parameter.Name, &parameter.Value, &parameter.ShouldShowToApp)

		if err != nil {
			_ = rows.Close()
			return
		}

		parameters = append(parameters, parameter)
	}

	dbaccess.LogDbResult(getAllSystemParametersQuery, parameters, err)

	_ = rows.Close()
	return
}

var getSystemParameterQuery = dbaccess.Query{Name: "GetSystemParameterQuery", SQL: `SELECT * FROM get_system_parameter($1::text)`, IsSilent: true}

func (db *Database) GetSystemParameterCommand(name string) (parameter typesysparams.SystemParameter, err error) {
	row := dbaccess.QueryRow(getSystemParameterQuery, name)

	err = row.Scan(&parameter.Id, &parameter.Name, &parameter.Value, &parameter.ShouldShowToApp)

	dbaccess.LogDbResult(getSystemParameterQuery, parameter, err)

	return
}

var changeSystemParameterValueQuery = dbaccess.Query{Name: "ChangeSystemParameterValueQuery", SQL: `SELECT change_system_parameter_value($1::text, $2::text)`}

func (db *Database) ChangeSystemParameterValueCommand(name string, value string) error {
	_, err := dbaccess.Exec(changeSystemParameterValueQuery, name, value)

	dbaccess.LogDbResult(changeSystemParameterValueQuery, nil, err)

	return err
}
