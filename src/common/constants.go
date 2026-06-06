package common

const (
	SessionCookieName = "session_id"

	MinimumNumberOfUnplayedGames      = 6
	DefaultTimerDurationInS           = 30
	DefaultTerritoryHoursIncreasing   = 2
	DefaultExperiencePointsIncreasing = 2
	DefaultExperiencePointsLevelUp    = -10
)

var DefaultTerritoryHoursSeizeDecreaseSlice = []int{-2, -4}

var _freePointsMinimum = 0
var FreePointMinimum *int = &_freePointsMinimum
