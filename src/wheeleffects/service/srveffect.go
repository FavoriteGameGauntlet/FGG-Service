package srvwheeleffects

import (
	"FGG-Service/src/common"
	"FGG-Service/src/points/service"
	"FGG-Service/src/points/type"
	"FGG-Service/src/sysparams/service"
	"FGG-Service/src/sysparams/types"
	"FGG-Service/src/wheeleffects/database"
	"FGG-Service/src/wheeleffects/types"
	"database/sql"
	"errors"
)

type IService interface {
	ApplyWheelEffectRoll(userId int, rollApply typewheeleffects.WheelEffectRollApply) (
		results typepoints.PointChangeResultByUserIds, err error)
	GetLastWheelEffectByName(userId int, effectName string) (effect typewheeleffects.RolledWheelEffect, err error)
}

type Service struct {
	Database         dbwheeleffects.IDatabase
	PointService     srvpoints.Service
	SysParamsService srvsysparams.IService
}

func NewService() *Service {
	db := new(dbwheeleffects.Database)
	ps := srvpoints.NewService()
	sp := srvsysparams.NewService()

	return &Service{
		db,
		*ps,
		sp,
	}
}

func (s *Service) GetAvailableRollsCount(userId int) (count int, err error) {
	return s.Database.GetAvailableRollsCountCommand(userId)
}

func (s *Service) GetAvailableEffects(userId int) (typewheeleffects.WheelEffects, error) {
	return s.Database.GetAvailableEffectsCommand(userId)
}

func (s *Service) GetEffectHistory(userId int) (typewheeleffects.RolledWheelEffectHistories, error) {
	return s.Database.GetEffectHistoryCommand(userId)
}

func (s *Service) GetEffectHistoryByEffectName(userId int, effectName string) (effect *typewheeleffects.RolledWheelEffect, err error) {
	notNilEffect, err := s.Database.GetEffectHistoryByEffectNameCommand(userId, effectName)

	if errors.Is(err, sql.ErrNoRows) {
		err = nil
		return
	}

	if err != nil {
		return
	}

	effect = &notNilEffect

	return
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

	minimumWheelEffectsForRoll, err := s.SysParamsService.GetInt(typesysparams.ParamMinimumWheelEffectsForRoll)

	if err != nil {
		return
	}

	if len(effects) < minimumWheelEffectsForRoll {
		err = common.NewNotEnoughAvailableWheelEffectsConflictError()
		return
	}

	err = s.Database.DecreaseAvailableRollsValueCommand(userId)

	if err != nil {
		return
	}

	err = s.Database.AddLastRolledWheelEffectsCommand(userId, effects)

	return
}

func (s *Service) GetLastRolledWheelEffects(userId int) (effects typewheeleffects.RolledWheelEffects, err error) {
	effects, err = s.Database.GetLastRolledWheelEffectsCommand(userId)

	if err == nil && len(effects) == 0 {
		err = common.NewLastWheelEffectsNotFoundError()
		return
	}

	return
}

func (s *Service) ApplyWheelEffectRoll(userId int, rollApply typewheeleffects.WheelEffectRollApply) (
	results typepoints.PointChangeResultByUserIds, err error) {

	effect, err := s.GetLastWheelEffectByName(userId, rollApply.WheelEffectName)

	if err != nil {
		return
	}

	results = make(typepoints.PointChangeResultByUserIds, len(rollApply.PointChangeByUserIds))

	for i, pointChange := range rollApply.PointChangeByUserIds {
		var changeResult typepoints.PointChangeResult
		changeResult, err = s.PointService.ChangeFreePoints(pointChange.UserId, pointChange.PointChange, &effect.Id)

		if err != nil {
			return
		}

		results[i] = typepoints.PointChangeResultByUserId{
			UserId:       pointChange.UserId,
			Login:        pointChange.Login,
			ChangeResult: changeResult,
		}
	}

	err = s.Database.MarkLastWheelEffectAppliedCommand(userId, effect.Id)

	if err != nil {
		return
	}

	err = s.Database.AddWheelEffectHistoryCommand(userId, effect.Id)

	return
}

func (s *Service) GetLastWheelEffectByName(userId int, effectName string) (
	effect typewheeleffects.RolledWheelEffect, err error) {

	lastEffects, err := s.Database.GetLastRolledWheelEffectsCommand(userId)

	if errors.Is(err, sql.ErrNoRows) {
		err = common.NewLastWheelEffectsNotFoundError()
		return
	}

	if err != nil {
		return
	}

	if len(lastEffects) == 0 {
		err = common.NewLastWheelEffectsNotFoundError()
		return
	}

	var foundLastEffect *typewheeleffects.RolledWheelEffect
	for _, lastEffect := range lastEffects {
		if lastEffect.Name == effectName {
			foundLastEffect = &lastEffect
			break
		}
	}

	if foundLastEffect == nil {
		err = common.NewWheelEffectNameNotFoundError()
		return
	}

	if foundLastEffect.IsApplied {
		err = common.NewWheelEffectRollAlreadyAppliedConflictError()
		return
	}

	effect = *foundLastEffect

	return
}
