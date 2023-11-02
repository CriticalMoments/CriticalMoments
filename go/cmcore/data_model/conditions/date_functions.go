package conditions

import (
	"strconv"
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

var timeFormats = map[string]string{
	// dow: 0-6, Sunday is 0
	"dow_short":   "Mon",     // string short day of week
	"dow_long":    "Monday",  // string long day of week
	"dom":         "02",      // int day of month
	"hod":         "15",      // int hour of day
	"moh":         "04",      // int minute of hour
	"ampm":        "PM",      // AM or PM
	"month":       "01",      // int month. Jan is 1, Dec is 12
	"month_short": "Jan",     // string short month of year
	"month_long":  "January", // string long month of year
	"year":        "2006",    // int year
}

func TimeFormat(t time.Time, format string, params ...string) interface{} {
	// Shortcut format convert to golang formats
	nfmt, ok := timeFormats[format]
	if ok {
		format = nfmt
	}

	// Convert timezone
	if len(params) > 1 {
		return nil // undefined behaviour
	}
	tz := ""
	if len(params) == 1 {
		tz = params[0]
	}
	if tz == "" {
		tz = "Local"
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return nil
	}
	t = t.In(loc)

	// Special case dow
	if format == "dow" {
		return int(t.Weekday())
	}

	s := t.Format(format)

	// return int if we can, nil if malformed, and otherwise string
	i, err := strconv.Atoi(s)
	if err == nil {
		return i
	}
	if len(s) > 0 {
		return s
	}
	return nil
}
