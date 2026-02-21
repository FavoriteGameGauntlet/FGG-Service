package users

import (
	"FGG-Service/api/generated/users"
	"FGG-Service/src/users/service"
	"context"
)

type Controller struct {
	Service user_service.Service
}

// ChangeOwnDisplayName POST /users/self/display-name
func (c *Controller) ChangeOwnDisplayName(ctx context.Context, req users.OptName) (
	users.ChangeOwnDisplayNameRes, error) {

	return nil, nil
}

// GetOwnDisplayName GET /users/self/display-name
func (c *Controller) GetOwnDisplayName(ctx context.Context) (
	users.GetOwnDisplayNameRes, error) {

	return nil, nil
}

// GetUserInfos GET /users/names
func (c *Controller) GetUserInfos(ctx context.Context) (
	users.GetUserInfosRes, error) {

	return nil, nil
}
