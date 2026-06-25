package dbpointsmock

import (
	typepoints "FGG-Service/src/points/type"

	"github.com/stretchr/testify/mock"
)

type DatabaseMock struct {
	mock.Mock
}

func (m *DatabaseMock) ChangeAvailableRollsCommand(userId int, changeValue int) error {
	args := m.Called(userId, changeValue)
	return args.Error(0)
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

func (m *DatabaseMock) GetFreePointHistoryCommand(userId int) (history typepoints.FreePointChangeHistories, err error) {
	args := m.Called(userId)
	history = args.Get(0).(typepoints.FreePointChangeHistories)
	err = args.Error(1)
	return
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

func (m *DatabaseMock) ChangeTerritoryPointsCommand(userId int, changeValue int) error {
	args := m.Called(userId, changeValue)
	return args.Error(0)
}

func (m *DatabaseMock) AddTerritoryPointHistoryCommand(userId int, sourceUserId int, changeSource string, changeValue int, actualChangeValue int, finalValue int) error {
	args := m.Called(userId, sourceUserId, changeSource, changeValue, actualChangeValue, finalValue)
	return args.Error(0)
}

func (m *DatabaseMock) GetTerritoryPointHistoryCommand(userId int) (history typepoints.TerritoryPointChangeHistories, err error) {
	args := m.Called(userId)
	history = args.Get(0).(typepoints.TerritoryPointChangeHistories)
	err = args.Error(1)
	return
}

func (m *DatabaseMock) GetPointInfoCommand(userId int) (info typepoints.PointInfo, err error) {
	args := m.Called(userId)
	info = args.Get(0).(typepoints.PointInfo)
	err = args.Error(1)
	return
}

func (m *DatabaseMock) GetAllPointInfoCommand() (infos typepoints.PointInfoByLogins, err error) {
	args := m.Called()
	infos = args.Get(0).(typepoints.PointInfoByLogins)
	err = args.Error(1)
	return
}
