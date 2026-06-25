package typesysparams

type SystemParameter struct {
	Id              int
	Name            string
	Value           string
	ShouldShowToApp bool
}

const (
	ParamMinimumNumberOfWishlistGames = "MinimumNumberOfWishlistGames"

	ParamTimerDurationInS                  = "TimerDurationInS"
	ParamTimerFinisherSchedulerIntervalInS = "TimerFinisherSchedulerIntervalInS"
	ParamMaximumAvailableRollCountForTimer = "MaximumAvailableRollCountForTimer"

	ParamAvailableRollChangeByTimer      = "AvailableRollChangeByTimer"
	ParamAvailableRollChangeByRoll       = "AvailableRollChangeByRoll"
	ParamTerritoryHourChangeByTimer      = "TerritoryHourChangeByTimer"
	ParamExperiencePointChangeByTimer    = "ExperiencePointChangeByTimer"
	ParamExperiencePointByLevelUp        = "ExperiencePointChangeByLevelUp"
	ParamTerritoryHourChangeBySeizeSlice = "TerritoryHourChangeBySeizeSlice"
	ParamFreePointsMinimum               = "FreePointsMinimum"
	ParamShouldLimitFreePoints           = "ShouldLimitFreePoints"
	ParamSeizePenaltyPoints              = "SeizePenaltyPoints"

	ParamMinimumAvailableRollCountForRoll    = "MinimumAvailableRollCountForRoll"
	ParamMinimumAvailableWheelEffectsForRoll = "MinimumAvailableWheelEffectsForRoll"
)
