package dbaccess

import (
	"log/slog"
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
	notNilFinishDate, err := time.Parse(ISO8601, dateString)

	if err != nil {
		return
	}

	date = notNilFinishDate

	return
}

func LogDbResult(queryName string, result any, err error) {
	if err != nil {
		slog.Error(queryName, "error", err)
	} else {
		slog.Info(queryName, "result", result)
	}
}
