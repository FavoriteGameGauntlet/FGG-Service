package srvwheeleffects

import (
	"FGG-Service/src/common"
	"FGG-Service/src/wheeleffects/database"
	"FGG-Service/src/wheeleffects/types"
)

type Service struct {
	Database dbwheeleffects.Database
}

func (s *Service) GetAvailableRollsCount(userId int) (count int, err error) {
	return s.Database.GetAvailableRollsCountCommand(userId)
}

func (s *Service) GetAvailableEffects(userId int) (typewheeleffects.WheelEffects, error) {
	return s.Database.GetAvailableEffectsCommand(userId)
}

func (s *Service) GetEffectHistory(userId int) (typewheeleffects.RolledWheelEffects, error) {
	return s.Database.GetEffectHistoryCommand(userId)
}

func (s *Service) MakeEffectRoll(userId int) (effects typewheeleffects.WheelEffects, err error) {
	rollCount, err := s.Database.GetAvailableRollsCountCommand(userId)

	if err != nil {
		return
	}

	if rollCount == 0 {
		err = common.NewAvailableRollsNotFoundError()
		return
	}

	effects, err = s.Database.MakeEffectRollCommand(userId)

	if err != nil {
		return
	}

	err = s.Database.DecreaseAvailableRollsValueCommand(userId)

	return
}
