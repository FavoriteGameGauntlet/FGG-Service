package wheel_effect_service

import (
	"FGG-Service/src/common"
	"FGG-Service/src/wheel-effects/whee_effect_db"
)

func GetAvailableRollsCount(userId int) (count int, err error) {
	return whee_effect_db.GetAvailableRollsCountCommand(userId)
}

func GetAvailableEffects(userId int) (common.Effects, error) {
	return whee_effect_db.GetAvailableEffectsCommand(userId)
}

func GetEffectHistory(userId int) (common.RolledEffects, error) {
	return whee_effect_db.GetEffectHistoryCommand(userId)
}

func MakeEffectRoll(userId int) (effects common.Effects, err error) {
	rollCount, err := whee_effect_db.GetAvailableRollsCountCommand(userId)

	if err != nil {
		return
	}

	if rollCount == 0 {
		err = common.NewAvailableRollsNotFoundError()
		return
	}

	effects, err = whee_effect_db.MakeEffectRollCommand(userId)

	if err != nil {
		return
	}

	err = whee_effect_db.DecreaseAvailableRollsValueCommand(userId)

	return
}
