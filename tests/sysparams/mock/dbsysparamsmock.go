package dbsysparamsmock

import (
	"FGG-Service/src/sysparams/types"

	"github.com/stretchr/testify/mock"
)

type DatabaseMock struct {
	mock.Mock
}

func (m *DatabaseMock) GetAllSystemParametersCommand() (parameters []typesysparams.SystemParameter, err error) {
	args := m.Called()
	parameters = args.Get(0).([]typesysparams.SystemParameter)
	err = args.Error(1)
	return
}

func (m *DatabaseMock) GetSystemParameterCommand(name string) (parameter typesysparams.SystemParameter, err error) {
	args := m.Called(name)
	parameter = args.Get(0).(typesysparams.SystemParameter)
	err = args.Error(1)
	return
}

func (m *DatabaseMock) ChangeSystemParameterValueCommand(name string, value string) error {
	args := m.Called(name, value)
	return args.Error(0)
}
