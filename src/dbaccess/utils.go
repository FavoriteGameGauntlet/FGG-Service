package dbaccess

import "time"

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
