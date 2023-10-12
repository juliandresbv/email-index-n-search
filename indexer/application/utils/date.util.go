package utils

import (
	"time"
)

func DatasetDateToISO8601Date(date string) (string, error) {
	datasetDateFormat := "Mon, 2 Jan 2006 15:04:05 -0700 (MST)"
	iso8601DateFormat := "2006-01-02T15:04:05.000Z07:00"

	t, err := time.Parse(datasetDateFormat, date)

	if err != nil {
		return "", err
	}

	formattedDate := t.Format(iso8601DateFormat)

	return formattedDate, nil
}
