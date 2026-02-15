package controller

import (
	"FGG-Service/api"
	"FGG-Service/src/common"
	"FGG-Service/src/effect/effect_service"
	"FGG-Service/src/game/game_service"
	"FGG-Service/src/validator"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetUnplayedGames (GET /games/unplayed)
func (Server) GetUnplayedGames(ctx echo.Context) error {
	userId, err := GetUserId(ctx)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	games, err := game_service.GetUnplayedGames(userId)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	gamesDto := convertWishlistGamesToDto(games)

	return ctx.JSON(http.StatusOK, gamesDto)
}

// AddUnplayedGames (POST /games/unplayed)
func (Server) AddUnplayedGames(ctx echo.Context) error {
	userId, err := GetUserId(ctx)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	var gamesDto api.WishlistGames
	err = ctx.Bind(&gamesDto)

	if err != nil {
		err = common.NewBadRequestError(err.Error())
		return SendJSONErrorResponse(ctx, err)
	}

	games := convertUnplayedGamesFromDto(gamesDto)

	for _, game := range games {
		err = validator.ValidateName(game.Name)

		if err != nil {
			return SendJSONErrorResponse(ctx, err)
		}
	}

	err = game_service.AddUnplayedGames(userId, games)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	return ctx.NoContent(http.StatusOK)
}

func convertUnplayedGamesFromDto(games api.WishlistGames) common.UnplayedGames {
	gamesDto := make(common.UnplayedGames, len(games))

	for i, g := range games {
		gamesDto[i] = common.UnplayedGame{
			Name: g.Name,
		}
	}

	return gamesDto
}

// GetAvailableEffects (GET /effects/available)
func (Server) GetAvailableEffects(ctx echo.Context) error {
	userId, err := GetUserId(ctx)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	effects, err := effect_service.GetAvailableEffects(userId)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	effectsDto := convertEffectsToDto(effects)

	return ctx.JSON(http.StatusOK, effectsDto)
}

// GetAvailableEffectsCount (GET /effects/available/count)
func (Server) GetAvailableEffectsCount(ctx echo.Context) error {
	userId, err := GetUserId(ctx)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	count, err := effect_service.GetAvailableRollsCount(userId)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, count)
}

// GetEffectHistory (GET /effects/history)
func (Server) GetEffectHistory(ctx echo.Context) error {
	userId, err := GetUserId(ctx)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	effects, err := effect_service.GetEffectHistory(userId)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	effectsDto := convertRolledEffectsToDto(effects)

	return ctx.JSON(http.StatusOK, effectsDto)
}

// MakeEffectRoll (POST /effects/roll)
func (Server) MakeEffectRoll(ctx echo.Context) error {
	userId, err := GetUserId(ctx)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	effects, err := effect_service.MakeEffectRoll(userId)

	if err != nil {
		return SendJSONErrorResponse(ctx, err)
	}

	effectsDto := convertEffectsToDto(effects)

	return ctx.JSON(http.StatusOK, effectsDto)
}
