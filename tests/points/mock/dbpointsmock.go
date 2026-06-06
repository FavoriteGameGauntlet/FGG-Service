package dbpointsmock

import "github.com/stretchr/testify/mock"

type DatabaseMock struct {
	mock.Mock
}

func (m *DatabaseMock) GetExperiencePointsCommand(userId int) (points int, err error) {
	args := m.Called(userId)
	points = args.Get(0).(int)
	err = args.Error(1)
	return
}

func (m *DatabaseMock) ChangeExperiencePointsCommand(userId int, changeValue int) error {
	args := m.Called(userId, changeValue)
	return args.Error(0)
}

func (m *DatabaseMock) GetFreePointsCommand(userId int) (points int, err error) {
	args := m.Called(userId)
	points = args.Get(0).(int)
	err = args.Error(1)
	return
}

func (m *DatabaseMock) ChangeFreePointsCommand(userId int, changeValue int) error {
	args := m.Called(userId, changeValue)
	return args.Error(0)
}

func (m *DatabaseMock) AddFreePointHistoryCommand(userId int, sourceUserId int, changeSource string, changeValue int, actualChangeValue int, finalValue int, wheelEffectId *int) error {
	args := m.Called(userId, sourceUserId, changeSource, changeValue, actualChangeValue, finalValue, wheelEffectId)
	return args.Error(0)
}

func (m *DatabaseMock) GetTerritoryHoursCommand(userId int) (points int, err error) {
	args := m.Called(userId)
	points = args.Get(0).(int)
	err = args.Error(1)
	return
}

func (m *DatabaseMock) ChangeTerritoryHoursCommand(userId int, changeValue int) error {
	args := m.Called(userId, changeValue)
	return args.Error(0)
}

func (m *DatabaseMock) GetTerritoryPointsCommand(userId int) (points int, err error) {
	args := m.Called(userId)
	points = args.Get(0).(int)
	err = args.Error(1)
	return
}
