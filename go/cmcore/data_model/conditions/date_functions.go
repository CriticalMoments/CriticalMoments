package conditions

import (
	"time"
)

func UnixTimeNanoseconds(u int64) time.Time {
	return time.Unix(u/1_000_000_000, u%1_000_000_000)
}

func UnixTimeMilliseconds(u int64) time.Time {
	return time.UnixMilli(u)
}

func UnixTimeSeconds(u int64) time.Time {
	return time.UnixMilli(u * 1000)
}

const DateWithTzFormat = "2006-01-02Z07:00"
const DateAndTimeFormat = "2006-01-02 15:04:05.999999999"
const DateFormat = "2006-01-02"
