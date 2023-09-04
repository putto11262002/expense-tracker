package utils

import (
	"time"
)

func ParseTimeFromISO8601(str string) (time.Time, error) {
	ISO8601Format := "2006-01-02T15:04:05Z07:00"
	parsedTime, err := time.Parse(ISO8601Format, str)
	if err != nil {
		return time.Time{}, err
	}
	return parsedTime, nil
}
