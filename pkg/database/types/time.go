// Package types provides basic operations with SQL types
package types

import (
	"time"
)

// StrToTime converts string to date
func StrToTime(input string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", input)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

// StrToDateTime converts string to datetime
func StrToDateTime(input string) (time.Time, error) {
	t, err := time.Parse("2006-01-02 15:04:05", input)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

// StrToDateTimeRef converts string to reference on datetime
func StrToDateTimeRef(input string) (*time.Time, error) {
	t, err := time.Parse("2006-01-02 15:04:05", input)
	if err != nil {
		return &time.Time{}, err
	}
	return &t, nil
}

// TimeToStr converts time.Time to string
func TimeToStr(time time.Time) string {
	return time.Format("2006-01-02")
}
