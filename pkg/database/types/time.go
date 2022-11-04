package types

import (
	"time"
)

func StrToTime(input string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", input)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func StrToDateTime(input string) (time.Time, error) {
	t, err := time.Parse("2006-01-02 15:04:05", input)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func StrToDateTimeRef(input string) (*time.Time, error) {
	t, err := time.Parse("2006-01-02 15:04:05", input)
	if err != nil {
		return &time.Time{}, err
	}
	return &t, nil
}

func TimeToStr(time time.Time) string {
	return time.Format("2006-01-02")
}
