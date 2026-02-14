package effect_service

import (
	"FGG-Service/src/common"
	"FGG-Service/src/effect/effect_db"
)

func GetAvailableRollsCount(userId int) (count int, err error) {
	return effect_db.GetAvailableRollsCountCommand(userId)
}

func GetAvailableEffects(userId int) (common.Effects, error) {
	return effect_db.GetAvailableEffectsCommand(userId)
}

func GetEffectHistory(userId int) (common.RolledEffects, error) {
	return effect_db.GetEffectHistoryCommand(userId)
}

func MakeEffectRoll(userId int) (effects common.Effects, err error) {
	rollCount, err := effect_db.GetAvailableRollsCountCommand(userId)

	if err != nil {
		return
	}

	if rollCount == 0 {
		err = common.NewAvailableRollsNotFoundError()
		return
	}

	effects, err = effect_db.MakeEffectRollCommand(userId)

	if err != nil {
		return
	}

	err = effect_db.DecreaseAvailableRollsValueCommand(userId)

	return
}
