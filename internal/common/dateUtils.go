package common

import "time"

func ToDateString(t *time.Time) string {
	if t == nil {
		return new(time.Time).Format("200601021504")
	} else {
		return t.Format("200601021504")
	}
}
