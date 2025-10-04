package controller

import (
	"FGG-Service/api"
	"context"
)

// CheckEffectRoll (GET /users/{userId}/effects/has-roll)
func (Server) CheckEffectRoll(_ context.Context, request api.CheckEffectRollRequestObject) (api.CheckEffectRollResponseObject, error) {
	return api.CheckEffectRoll200JSONResponse(true), nil
}

// GetEffectHistory (GET /users/{userId}/effects/history)
func (Server) GetEffectHistory(_ context.Context, request api.GetEffectHistoryRequestObject) (api.GetEffectHistoryResponseObject, error) {
	return api.GetEffectHistory200JSONResponse{}, nil
}

// MakeEffectRoll (POST /users/{userId}/effects/roll)
func (Server) MakeEffectRoll(_ context.Context, request api.MakeEffectRollRequestObject) (api.MakeEffectRollResponseObject, error) {
	return api.MakeEffectRoll200JSONResponse{}, nil
}
