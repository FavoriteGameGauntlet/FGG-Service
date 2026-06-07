package dbwheeleffectsmock

import (
	"FGG-Service/src/wheeleffects/types"

	"github.com/stretchr/testify/mock"
)

type DatabaseMock struct {
	mock.Mock
}

func (m *DatabaseMock) GetAvailableRollsCountCommand(userId int) (count int, err error) {
	args := m.Called(userId)
	count = args.Get(0).(int)
	err = args.Error(1)
	return
}

func (m *DatabaseMock) GetAvailableEffectsCommand(userId int) (effects typewheeleffects.WheelEffects, err error) {
	args := m.Called(userId)
	effects = args.Get(0).(typewheeleffects.WheelEffects)
	err = args.Error(1)
	return
}

func (m *DatabaseMock) GetEffectHistoryCommand(userId int) (effects typewheeleffects.RolledWheelEffects, err error) {
	args := m.Called(userId)
	effects = args.Get(0).(typewheeleffects.RolledWheelEffects)
	err = args.Error(1)
	return
}

func (m *DatabaseMock) MakeEffectRollCommand(userId int) (effects typewheeleffects.WheelEffects, err error) {
	args := m.Called(userId)
	effects = args.Get(0).(typewheeleffects.WheelEffects)
	err = args.Error(1)
	return
}

func (m *DatabaseMock) DecreaseAvailableRollsValueCommand(userId int) error {
	args := m.Called(userId)
	return args.Error(0)
}

func (m *DatabaseMock) AddLastRolledWheelEffectsCommand(userId int, effects typewheeleffects.WheelEffects) error {
	args := m.Called(userId, effects)
	return args.Error(0)
}

func (m *DatabaseMock) GetLastRolledWheelEffectsCommand(userId int) (effects typewheeleffects.RolledWheelEffects, err error) {
	args := m.Called(userId)
	effects = args.Get(0).(typewheeleffects.RolledWheelEffects)
	err = args.Error(1)
	return
}

func (m *DatabaseMock) GetEffectHistoryByEffectNameCommand(userId int, effectName string) (effect typewheeleffects.RolledWheelEffect, err error) {
	args := m.Called(userId, effectName)
	effect = args.Get(0).(typewheeleffects.RolledWheelEffect)
	err = args.Error(1)
	return
}

func (m *DatabaseMock) MarkLastWheelEffectAppliedCommand(userId int, wheelEffectId int) error {
	args := m.Called(userId, wheelEffectId)
	return args.Error(0)
}

func (m *DatabaseMock) AddWheelEffectHistoryCommand(userId int, wheelEffectId int) error {
	args := m.Called(userId, wheelEffectId)
	return args.Error(0)
}
