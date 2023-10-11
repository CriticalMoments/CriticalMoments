package conditions

import (
	"time"
)

// Our conditions use millisecond unix time.
// We use int64 since expr library doesn't have a date type.
// This is a set of helpers we expose as functions in
// our condition system to make working with milliseconds easy.

func Now() int64 {
	return time.Now().UnixMilli()
}

func Days(c int64) int64 {
	// 24 hours == day is simplified concept, will document
	return c * 24 * time.Hour.Milliseconds()
}

func Hours(c int64) int64 {
	return c * time.Hour.Milliseconds()
}

func Minutes(c int64) int64 {
	return c * time.Minute.Milliseconds()
}

func Seconds(c int64) int64 {
	return c * time.Second.Milliseconds()
}

func ParseDatetime(s string) (int64, error) {
	// RFC3339
	t, err := time.Parse(time.RFC3339Nano, s)
	if err == nil {
		return t.UnixMilli(), nil
	}

	// Date and timezone
	t, err = time.Parse("2006-01-02Z07:00", s)
	if err == nil {
		return t.UnixMilli(), nil
	}

	// RFC3339 (TZ removed), in local timezone
	t, err = time.ParseInLocation("2006-01-02T15:04:05.999999999", s, time.Local)
	if err == nil {
		return t.UnixMilli(), nil
	}

	// Date in local timezone
	t, err = time.ParseInLocation(time.DateOnly, s, time.Local)
	if err == nil {
		return t.UnixMilli(), nil
	}

	return 0, err
}
