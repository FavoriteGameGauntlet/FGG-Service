package typesysparams

type SystemParameter struct {
	Id              int
	Name            string
	Value           string
	ShouldShowToApp bool
}

const (
	ParamMinimumNumberOfUnplayedGames      = "MinimumNumberOfUnplayedGames"
	ParamTimerDurationInS                  = "TimerDurationInS"
	ParamTerritoryHoursIncreaseByTimer     = "TerritoryHoursIncreaseByTimer"
	ParamExperiencePointsIncreaseByTimer   = "ExperiencePointsIncreaseByTimer"
	ParamExperiencePointsLevelUp           = "ExperiencePointsLevelUp"
	ParamTerritoryHoursSeizeDecreaseSlice  = "TerritoryHoursSeizeDecreaseSlice"
	ParamFreePointsMinimum                 = "FreePointsMinimum"
	ParamShouldLimitFreePoints             = "ShouldLimitFreePoints"
	ParamMinimumWheelEffectsForRoll        = "MinimumWheelEffectsForRoll"
	ParamSeizePenaltyPoints                = "SeizePenaltyPoints"
	ParamTimerFinisherSchedulerIntervalInS = "TimerFinisherSchedulerIntervalInS"
	ParamMinimumRollsCountForEffectRoll    = "MinimumRollsCountForEffectRoll"
	ParamMaximumRollsCountForTimer         = "MaximumRollsCountForTimer"
)
