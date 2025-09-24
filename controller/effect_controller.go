package controller

import (
	"FGG-Service/api"
	"context"
)

// (GET /users/{userId}/effects/has-roll)
func (Server) CheckEffectRoll(ctx context.Context, request api.CheckEffectRollRequestObject) (api.CheckEffectRollResponseObject, error) {
	return api.CheckEffectRoll200JSONResponse(true), nil
}

// (GET /users/{userId}/effects/history)
func (Server) GetEffectHistory(ctx context.Context, request api.GetEffectHistoryRequestObject) (api.GetEffectHistoryResponseObject, error) {
	return api.GetEffectHistory200JSONResponse{}, nil
}

// (POST /users/{userId}/effects/roll)
func (Server) MakeEffectRoll(ctx context.Context, request api.MakeEffectRollRequestObject) (api.MakeEffectRollResponseObject, error) {
	return api.MakeEffectRoll200JSONResponse{}, nil
}
