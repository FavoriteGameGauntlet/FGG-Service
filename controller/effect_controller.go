package controller

import (
	"FGG-Service/api"
	"FGG-Service/auth_service"
	"FGG-Service/effect_service"
	"context"
)

// GetAvailableEffects (POST /users/{userId}/effects/available)
func (Server) GetAvailableEffects(_ context.Context, _ api.GetAvailableEffectsRequestObject) (api.GetAvailableEffectsResponseObject, error) {
	return api.GetAvailableEffects200JSONResponse{}, nil
}

// CheckAvailableEffectRoll CheckEffectRoll (GET /users/{userId}/effects/has-roll)
func (Server) CheckAvailableEffectRoll(ctx context.Context, _ api.CheckAvailableEffectRollRequestObject) (api.CheckAvailableEffectRollResponseObject, error) {
	sessionId, ok := ctx.Value("session_id").(string)

	if !ok {
		return api.CheckAvailableEffectRoll401JSONResponse{Code: api.NOACTIVESESSION}, nil
	}

	userId, err := auth_service.GetUserId(sessionId)

	if err != nil {
		return api.CheckAvailableEffectRoll500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	doesExist, err := effect_service.CheckIfAvailableRollExists(userId)

	if err != nil {
		return api.CheckAvailableEffectRoll500JSONResponse{Code: api.UNEXPECTED, Message: err.Error()}, nil
	}

	return api.CheckAvailableEffectRoll200JSONResponse(doesExist), nil
}

// GetEffectHistory (GET /users/{userId}/effects/history)
func (Server) GetEffectHistory(_ context.Context, _ api.GetEffectHistoryRequestObject) (api.GetEffectHistoryResponseObject, error) {
	return api.GetEffectHistory200JSONResponse{}, nil
}

// MakeEffectRoll (POST /users/{userId}/effects/roll)
func (Server) MakeEffectRoll(_ context.Context, _ api.MakeEffectRollRequestObject) (api.MakeEffectRollResponseObject, error) {
	return api.MakeEffectRoll200JSONResponse{}, nil
}
