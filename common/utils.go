package common

import (
	"FGG-Service/db_access"
	"time"
)

func ConvertToNullableDate(dateString *string) (*time.Time, error) {
	var date *time.Time

	if dateString != nil {
		var notNilFinishDate time.Time
		notNilFinishDate, err := time.Parse(db_access.ISO8601, *dateString)

		if err != nil {
			return nil, err
		}

		date = &notNilFinishDate
	}

	return date, nil
}
