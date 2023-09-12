package utils

import "time"

func ParseStringISODate(stringDate string) *time.Time {
	if stringDate == "" {
		return nil
	}
	date, _ := time.Parse(time.RFC3339, stringDate)
	return &date
}
