package controller

import (
	"FGG-Service/api"
	"FGG-Service/effect_service"
	"FGG-Service/user_service"
	"context"
)

// GetAvailableEffects (POST /users/{userId}/effects/available)
func (Server) GetAvailableEffects(_ context.Context, request api.GetAvailableEffectsRequestObject) (api.GetAvailableEffectsResponseObject, error) {
	return api.GetAvailableEffects200JSONResponse{}, nil
}

// CheckEffectRoll (GET /users/{userId}/effects/has-roll)
func (Server) CheckEffectRoll(_ context.Context, request api.CheckEffectRollRequestObject) (api.CheckEffectRollResponseObject, error) {
	doesExist, err := user_service.CheckIfUserExistsById(request.UserId)

	if err != nil {
		return api.CheckEffectRoll500JSONResponse{Code: api.UNEXPECTEDDATABASE, Message: err.Error()}, nil
	}

	if !doesExist {
		return api.CheckEffectRoll404JSONResponse{Code: api.USERNOTFOUND}, nil
	}

	doesExist, err = effect_service.CheckIfEffectRollExists(request.UserId)

	if err != nil {
		return api.CheckEffectRoll500JSONResponse{Code: api.UNEXPECTEDDATABASE, Message: err.Error()}, nil
	}

	return api.CheckEffectRoll200JSONResponse(doesExist), nil
}

// GetEffectHistory (GET /users/{userId}/effects/history)
func (Server) GetEffectHistory(_ context.Context, request api.GetEffectHistoryRequestObject) (api.GetEffectHistoryResponseObject, error) {
	doesExist, err := user_service.CheckIfUserExistsById(request.UserId)

	if err != nil {
		return api.GetEffectHistory500JSONResponse{Code: api.UNEXPECTEDDATABASE, Message: err.Error()}, nil
	}

	if !doesExist {
		return api.GetEffectHistory404JSONResponse{Code: api.USERNOTFOUND}, nil
	}

	effects, err := effect_service.GetEffectHistory(request.UserId)

	if err != nil {
		return api.GetEffectHistory500JSONResponse{Code: api.UNEXPECTEDDATABASE, Message: err.Error()}, nil
	}

	effectsDto := ConvertEffectsToDto(effects)

	return api.GetEffectHistory200JSONResponse(*effectsDto), nil
}

func ConvertEffectsToDto(effects *effect_service.Effects) *api.Effects {
	effectsDto := make(api.Effects, len(*effects))

	for i, effect := range *effects {
		effectsDto[i] = *ConvertEffectToDto(&effect)
	}

	return &effectsDto
}

// MakeEffectRoll (POST /users/{userId}/effects/roll)
func (Server) MakeEffectRoll(_ context.Context, request api.MakeEffectRollRequestObject) (api.MakeEffectRollResponseObject, error) {
	return api.MakeEffectRoll200JSONResponse{}, nil
}

func ConvertEffectToDto(effect *effect_service.Effect) *api.Effect {
	return &api.Effect{
		CreateDate:  effect.CreateDate,
		Description: effect.Description,
		GameName:    effect.GameName,
		Name:        effect.Name,
		RollDate:    effect.RollDate,
	}
}
