package common

import (
	"fmt"
	"time"
)

func DurationToISO8601(duration time.Duration) string {
	if duration == 0 {
		return "PT0S"
	}

	sign := ""
	if duration < 0 {
		sign = "-"
		duration = -duration
	}

	totalSeconds := int64(duration.Seconds())
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	result := fmt.Sprintf("%sPT", sign)

	if hours > 0 {
		result += fmt.Sprintf("%dH", hours)
	}
	if minutes > 0 {
		result += fmt.Sprintf("%dM", minutes)
	}
	if seconds > 0 {
		result += fmt.Sprintf("%dS", seconds)
	}

	return result

}
