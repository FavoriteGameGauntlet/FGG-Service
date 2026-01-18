package effect_service

import (
	"FGG-Service/common"
)

func GetAvailableRollsCount(userId int) (count int, err error) {
	return GetAvailableRollsCountCommand(userId)
}

func GetAvailableEffects(userId int) (common.Effects, error) {
	return GetAvailableEffectsCommand(userId)
}

func GetEffectHistory(userId int) (common.RolledEffects, error) {
	return GetEffectHistoryCommand(userId)
}

func MakeEffectRoll(userId int) (effects common.Effects, err error) {
	rollCount, err := GetAvailableRollsCountCommand(userId)

	if err != nil {
		return
	}

	if rollCount == 0 {
		err = common.NewAvailableRollsNotFoundError()
		return
	}

	effects, err = MakeEffectRollCommand(userId)

	if err != nil {
		return
	}

	err = DeleteAvailableRollCommand(userId)

	return
}
