package common

import (
	"time"
)

func ToDateString(t *time.Time) string {
	// format RFC3339
	if t == nil {
		return new(time.Time).Format("2006-01-02")
	} else {
		return t.Format("2006-01-02")
	}
}

func ToDateTimeString(t *time.Time) string {
	// format RFC3339
	if t == nil {
		return new(time.Time).Format(time.DateTime)
	} else {
		return t.Format(time.DateTime)
	}
}

func ToTime(s string) *time.Time {
	// format RFC3339
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil
	}
	return &t
}
