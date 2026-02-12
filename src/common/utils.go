package common

import (
	"FGG-Service/src/db_access"
	"fmt"
	"time"
)

func ConvertToNullableDate(dateString *string) (date *time.Time, err error) {
	if dateString != nil {
		var notNilFinishDate time.Time
		notNilFinishDate, err = ConvertToDate(*dateString)

		if err != nil {
			return
		}

		date = &notNilFinishDate
	}

	return
}

func ConvertToDate(dateString string) (date time.Time, err error) {
	notNilFinishDate, err := time.Parse(db_access.ISO8601, dateString)

	if err != nil {
		return
	}

	date = notNilFinishDate

	return
}

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
