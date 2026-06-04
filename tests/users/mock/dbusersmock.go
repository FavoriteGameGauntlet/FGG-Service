package dbusersmock

import (
	"FGG-Service/src/users/types"

	"github.com/stretchr/testify/mock"
)

type DatabaseMock struct {
	mock.Mock
}

func (m *DatabaseMock) GetAllUserNamesCommand() (users typeusers.Users, err error) {
	args := m.Called()
	users = args.Get(0).(typeusers.Users)
	err = args.Error(1)
	return
}

func (m *DatabaseMock) GetDisplayNameCommand(userId int) (displayName *string, err error) {
	args := m.Called(userId)
	displayName = args.Get(0).(*string)
	err = args.Error(1)
	return
}

func (m *DatabaseMock) ChangeDisplayNameCommand(userId int, displayName string) error {
	args := m.Called(userId, displayName)
	return args.Error(0)
}
