package typesysparams

type SystemParameter struct {
	Id    int
	Name  string
	Value string
}

const (
	ParamMinimumNumberOfUnplayedGames      = "MinimumNumberOfUnplayedGames"
	ParamDefaultTimerDurationInS           = "DefaultTimerDurationInS"
	ParamDefaultTerritoryHoursIncreasing   = "DefaultTerritoryHoursIncreasing"
	ParamDefaultExperiencePointsIncreasing = "DefaultExperiencePointsIncreasing"
	ParamDefaultExperiencePointsLevelUp    = "DefaultExperiencePointsLevelUp"
	ParamTerritoryHoursSeizeDecreaseSlice  = "TerritoryHoursSeizeDecreaseSlice"
	ParamFreePointMinimum                  = "FreePointMinimum"
	ParamShouldLimitFreePoints             = "ShouldLimitFreePoints"
)
