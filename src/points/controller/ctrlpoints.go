package points

import (
	"FGG-Service/api/generated/points"
	"FGG-Service/src/points/service"
	"context"
)

type Controller struct {
	Service srvpoints.Service
}

// ChangeOwnExperiencePoints POST /points/self/experience-points
func (c *Controller) ChangeOwnExperiencePoints(ctx context.Context, req points.OptPointChange) (
	points.ChangeOwnExperiencePointsRes, error) {

	return nil, nil
}

// ChangeOwnFreePoints POST /points/self/free-points
func (c *Controller) ChangeOwnFreePoints(ctx context.Context, req points.OptPointChange) (
	points.ChangeOwnFreePointsRes, error) {

	return nil, nil
}

// ChangeOwnTerritoryHours POST /points/self/territory-hours
func (c *Controller) ChangeOwnTerritoryHours(ctx context.Context, req points.OptPointChange) (
	points.ChangeOwnTerritoryHoursRes, error) {

	return nil, nil
}

// ChangeOwnTerritoryPoints POST /points/self/territory-points
func (c *Controller) ChangeOwnTerritoryPoints(ctx context.Context, req points.OptPointChange) (
	points.ChangeOwnTerritoryPointsRes, error) {

	return nil, nil
}

// GetFreePointHistoryByLogin GET /points/{login}/free-points/history
func (c *Controller) GetFreePointHistoryByLogin(ctx context.Context, params points.GetFreePointHistoryByLoginParams) (
	points.GetFreePointHistoryByLoginRes, error) {

	return nil, nil
}

// GetOwnExperiencePoints GET /points/self/experience-points
func (c *Controller) GetOwnExperiencePoints(ctx context.Context) (
	points.GetOwnExperiencePointsRes, error) {

	return nil, nil
}

// GetOwnFreePointHistory GET /points/self/free-points/history
func (c *Controller) GetOwnFreePointHistory(ctx context.Context) (
	points.GetOwnFreePointHistoryRes, error) {

	return nil, nil
}

// GetOwnFreePoints GET /points/self/free-points
func (c *Controller) GetOwnFreePoints(ctx context.Context) (
	points.GetOwnFreePointsRes, error) {

	return nil, nil
}

// GetOwnPointInfo GET /points/self/info
func (c *Controller) GetOwnPointInfo(ctx context.Context) (
	points.GetOwnPointInfoRes, error) {

	return nil, nil
}

// GetOwnTerritoryHours GET /points/self/territory-hours
func (c *Controller) GetOwnTerritoryHours(ctx context.Context) (
	points.GetOwnTerritoryHoursRes, error) {

	return nil, nil
}

// GetOwnTerritoryPointHistory GET /points/self/territory-points/history
func (c *Controller) GetOwnTerritoryPointHistory(ctx context.Context) (
	points.GetOwnTerritoryPointHistoryRes, error) {

	return nil, nil
}

// GetOwnTerritoryPoints GET /points/self/territory-points
func (c *Controller) GetOwnTerritoryPoints(ctx context.Context) (
	points.GetOwnTerritoryPointsRes, error) {

	return nil, nil
}

// GetPointInfoByLogin GET /points/{login}/info
func (c *Controller) GetPointInfoByLogin(ctx context.Context, params points.GetPointInfoByLoginParams) (
	points.GetPointInfoByLoginRes, error) {

	return nil, nil
}

// GetTerritoryPointHistoryByLogin GET /points/{login}/territory-points/history
func (c *Controller) GetTerritoryPointHistoryByLogin(ctx context.Context, params points.GetTerritoryPointHistoryByLoginParams) (
	points.GetTerritoryPointHistoryByLoginRes, error) {

	return nil, nil
}
