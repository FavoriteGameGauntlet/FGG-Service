package srveffects

import (
	"FGG-Service/src/common"
	"FGG-Service/src/effects/database"
)

type Service struct {
	Database dbeffects.Database
}

func (s *Service) GetAvailableRollsCount(userId int) (count int, err error) {
	return s.Database.GetAvailableRollsCountCommand(userId)
}

func (s *Service) GetAvailableEffects(userId int) (common.Effects, error) {
	return s.Database.GetAvailableEffectsCommand(userId)
}

func (s *Service) GetEffectHistory(userId int) (common.RolledEffects, error) {
	return s.Database.GetEffectHistoryCommand(userId)
}

func (s *Service) MakeEffectRoll(userId int) (effects common.Effects, err error) {
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
