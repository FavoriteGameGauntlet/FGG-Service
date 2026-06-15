package srvsysparams

import (
	"FGG-Service/src/common"
	"FGG-Service/src/sysparams/database"
	"FGG-Service/src/sysparams/types"
	"database/sql"
	"errors"
	"strconv"
	"strings"
)

type IService interface {
	GetAll() ([]typesysparams.SystemParameter, error)
	GetAllApp() ([]typesysparams.SystemParameter, error)
	GetParameter(name string) (typesysparams.SystemParameter, error)
	GetAppParameter(name string) (typesysparams.SystemParameter, error)
	GetString(name string) (string, error)
	GetInt(name string) (int, error)
	GetBool(name string) (bool, error)
	GetIntSlice(name string) ([]int, error)
	ChangeValue(name string, value string) error
}

type Service struct {
	Database dbsysparams.IDatabase
}

func NewService() *Service {
	db := new(dbsysparams.Database)

	return &Service{
		Database: db,
	}
}

func (s *Service) GetAll() (parameters []typesysparams.SystemParameter, err error) {
	parameters, err = s.Database.GetAllSystemParametersCommand()

	return
}

func (s *Service) GetAllApp() (parameters []typesysparams.SystemParameter, err error) {
	all, err := s.Database.GetAllSystemParametersCommand()

	if err != nil {
		return
	}

	for _, p := range all {
		if p.ShouldShowToApp {
			parameters = append(parameters, p)
		}
	}

	return
}

func (s *Service) GetParameter(name string) (parameter typesysparams.SystemParameter, err error) {
	parameter, err = s.Database.GetSystemParameterCommand(name)

	if errors.Is(err, sql.ErrNoRows) {
		err = common.NewSystemParameterNotFoundError(name)
		return
	}

	return
}

func (s *Service) GetAppParameter(name string) (parameter typesysparams.SystemParameter, err error) {
	parameter, err = s.Database.GetSystemParameterCommand(name)

	if errors.Is(err, sql.ErrNoRows) {
		err = common.NewSystemParameterNotFoundError(name)
		return
	}

	if !parameter.ShouldShowToApp {
		err = common.NewSystemParameterNotFoundError(name)
		return
	}

	return
}

func (s *Service) GetString(name string) (value string, err error) {
	parameter, err := s.GetParameter(name)

	if err != nil {
		return
	}

	value = parameter.Value

	return
}

func (s *Service) GetInt(name string) (value int, err error) {
	stringValue, err := s.GetString(name)

	if err != nil {
		return
	}

	value, err = strconv.Atoi(stringValue)

	return
}

func (s *Service) GetBool(name string) (value bool, err error) {
	str, err := s.GetString(name)

	if err != nil {
		return
	}

	if str == "" {
		return
	}

	value, err = strconv.ParseBool(str)

	return
}

func (s *Service) GetIntSlice(name string) (values []int, err error) {
	str, err := s.GetString(name)

	if err != nil {
		return
	}

	for _, part := range strings.Split(strings.Trim(str, "[]"), ",") {
		var v int
		v, err = strconv.Atoi(strings.TrimSpace(part))

		if err != nil {
			return
		}

		values = append(values, v)
	}

	return
}

func (s *Service) ChangeValue(name string, value string) (err error) {
	_, err = s.GetParameter(name)

	if err != nil {
		return
	}

	err = s.Database.ChangeSystemParameterValueCommand(name, value)

	return
}
