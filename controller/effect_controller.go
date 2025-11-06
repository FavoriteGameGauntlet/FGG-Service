package controller

import (
	"FGG-Service/api"
	"FGG-Service/effect_service"
	"FGG-Service/user_service"
	"context"
)

// CheckEffectRoll (GET /users/{userId}/effects/has-roll)
func (Server) CheckEffectRoll(_ context.Context, request api.CheckEffectRollRequestObject) (api.CheckEffectRollResponseObject, error) {
	return api.CheckEffectRoll200JSONResponse(true), nil
}

// GetEffectHistory (GET /users/{userId}/effects/history)
func (Server) GetEffectHistory(_ context.Context, request api.GetEffectHistoryRequestObject) (api.GetEffectHistoryResponseObject, error) {
	doesUserExist, err := user_service.CheckIfUserExistsById(request.UserId)

	if err != nil {
		return api.GetEffectHistory503JSONResponse{Code: api.CHECKUSER, Message: err.Error()}, nil
	}

	if !doesUserExist {
		return api.GetEffectHistory404JSONResponse{Code: api.USERNOTFOUND}, nil
	}

	effects, err := effect_service.GetEffectHistory(request.UserId)

	if err != nil {
		return api.GetEffectHistory503JSONResponse{Code: api.GETEFFECTHISTORY, Message: err.Error()}, nil
	}

	effectsDto := ConvertEffectsToDto(effects)

	return api.GetEffectHistory200JSONResponse(*effectsDto), nil
}

func ConvertEffectsToDto(effects *effect_service.Effects) *api.EffectsDto {
	effectsDto := make(api.EffectsDto, len(*effects))

	for i, effect := range *effects {
		effectsDto[i] = *ConvertEffectToDto(&effect)
	}

	return &effectsDto
}

// MakeEffectRoll (POST /users/{userId}/effects/roll)
func (Server) MakeEffectRoll(_ context.Context, request api.MakeEffectRollRequestObject) (api.MakeEffectRollResponseObject, error) {
	return api.MakeEffectRoll200JSONResponse{}, nil
}

func ConvertEffectToDto(effect *effect_service.Effect) *api.EffectDto {
	return &api.EffectDto{
		CreateDate:  effect.CreateDate,
		Description: effect.Description,
		GameName:    effect.GameName,
		Name:        effect.Name,
		RollDate:    effect.RollDate,
	}
}
