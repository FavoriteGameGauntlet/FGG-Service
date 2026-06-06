package srvwheeleffectsmock

import (
	typewheeleffects "FGG-Service/src/wheeleffects/types"

	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func (m *ServiceMock) GetEffectHistoryByEffectName(userId int, effectName string) (effect *typewheeleffects.RolledWheelEffect, err error) {
	args := m.Called(userId, effectName)
	if args.Get(0) != nil {
		effect = args.Get(0).(*typewheeleffects.RolledWheelEffect)
	}
	err = args.Error(1)
	return
}
