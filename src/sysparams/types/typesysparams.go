package typesysparams

type SystemParameter struct {
	Id    int
	Name  string
	Value string
}

const (
	ParamMinimumNumberOfUnplayedGames     = "MinimumNumberOfUnplayedGames"
	ParamTimerDurationInS                 = "TimerDurationInS"
	ParamTerritoryHoursIncreaseByTimer    = "TerritoryHoursIncreaseByTimer"
	ParamExperiencePointsIncreaseByTimer  = "ExperiencePointsIncreaseByTimer"
	ParamExperiencePointsLevelUp          = "ExperiencePointsLevelUp"
	ParamTerritoryHoursSeizeDecreaseSlice = "TerritoryHoursSeizeDecreaseSlice"
	ParamFreePointsMinimum                = "FreePointsMinimum"
	ParamShouldLimitFreePoints            = "ShouldLimitFreePoints"
)
