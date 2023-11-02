package conditions

import (
	"testing"
	"time"

	"github.com/antonmedv/expr/builtin"
)

func TestUnixTimeHelpers(t *testing.T) {
	// Test nanoseconds
	r1 := UnixTimeNanoseconds(1136189045999999999)
	if r1.UnixMilli() != 1136189045999 {
		t.Fatal("UnixTimeNanoseconds failed")
	}
	if r1.UnixNano() != 1136189045999999999 {
		t.Fatal("UnixTimeNanoseconds failed")
	}

	// Test milliseconds
	r2 := UnixTimeMilliseconds(1136189045999)
	if r2.UnixMilli() != 1136189045999 {
		t.Fatal("UnixTimeMilliseconds failed")
	}

	// Test seconds
	r3 := UnixTimeSeconds(1136189045)
	if r3.UnixMilli() != 1136189045000 {
		t.Fatal("UnixTimeSeconds failed")
	}
}

func TestParseDatetime(t *testing.T) {
	builtin := builtin.Builtins
	var rawDateFunc func(args ...any) (any, error)
	for _, bi := range builtin {
		if bi.Name == "date" {
			rawDateFunc = bi.Func
		}
	}
	if rawDateFunc == nil {
		t.Fatal("date function not found")
	}
	dateFunc := func(args ...any) (time.Time, error) {
		r, err := rawDateFunc(args...)
		if err != nil {
			return time.Time{}, err
		}
		return r.(time.Time), nil
	}

	tr, err := dateFunc("2006-01-02T15:04:05.9997+07:00")
	if err != nil || tr.UnixMilli() != 1136189045999 {
		t.Fatal("Datetime parsing failed")
	}

	// Test timezone parsing, 7h offset
	r2, err := dateFunc("2006-01-02T15:04:05.9997+00:00")
	if err != nil || r2.UnixMilli() != 1136214245999 || r2.Sub(tr) != time.Hour*7 {
		t.Fatal("Datetime parsing failed")
	}

	// Ensure works without decimal Seconds
	r3, err := dateFunc("2006-01-02T15:04:05+00:00")
	if err != nil || r3.UnixMilli() != 1136214245000 {
		t.Fatal("Datetime parsing failed")
	}

	// Ensure works with Z for UTC
	r4, err := dateFunc("2006-01-02T15:04:05Z")
	if err != nil || r3.UnixMilli() != r4.UnixMilli() {
		t.Fatal("Date parsing with TZ failed")
	}

	// Date Only
	r5, err := dateFunc("2000-01-29Z", DateWithTzFormat)
	r6, err2 := dateFunc("2000-01-29+00:00", DateWithTzFormat)
	if err != nil || err2 != nil || r5.UnixMilli() != 949104000000 || r5.UnixMilli() != r6.UnixMilli() {
		t.Fatal("Date parsing failed")
	}

	// No timezone provided, should parse in local TZ
	local := time.Local
	defer func() {
		time.Local = local
	}()
	time.Local, _ = time.LoadLocation("America/Toronto")
	r7, err := dateFunc("2023-01-01 11:11:11", DateAndTimeFormat, "Local") // -5 from local offset
	r8, err2 := dateFunc("2023-01-01T16:11:11Z")                           // UTC adjusted by hand
	if err != nil || err2 != nil || r7.UnixMilli() != r8.UnixMilli() {
		t.Fatal("Datetime parsing failed for local TZ")
	}

	// No TZ provided should parse in local date TZ
	r9, err := dateFunc("2023-01-01", DateFormat, "Local") // -5 from local offset
	r10, err2 := dateFunc("2023-01-01T05:00:00Z")          // UTC adjusted by hand
	if err != nil || err2 != nil || r9.UnixMilli() != r10.UnixMilli() {
		t.Fatal("Date parsing failed for local TZ")
	}

	// UTC default when not specified
	r11, err := dateFunc("2023-01-01", DateFormat) // UTC
	r12, err2 := dateFunc("2023-01-01T00:00:00Z")  // UTC
	if err != nil || err2 != nil || r11.UnixMilli() != r12.UnixMilli() {
		t.Fatal("Date parsing didn't default to UTC")
	}

	time.Local = local
}

func TestTimeFormat(t *testing.T) {
	ti := time.UnixMilli(1698881571000)

	// Test golang formats
	r := TimeFormat(ti, "Mon", "America/Toronto")
	if r != "Wed" {
		t.Fatal("Failed to parse dow string")
	}
	r = TimeFormat(ti, "1", "America/Toronto")
	if r != 11 {
		t.Fatal("Failed to parse moy int")
	}
	r = TimeFormat(ti, "Jan", "America/Toronto")
	if r != "Nov" {
		t.Fatal("Failed to parse moy string")
	}
	r = TimeFormat(ti, "January", "America/Toronto")
	if r != "November" {
		t.Fatal("Failed to parse moy string")
	}
	r = TimeFormat(ti, "2006", "America/Toronto")
	if r != 2023 {
		t.Fatal("Failed to parse year")
	}

	// Test our short hands
	r = TimeFormat(ti, "dow", "America/Toronto")
	if r != 3 {
		t.Fatal("Failed to parse dow string")
	}
	r = TimeFormat(ti, "dow_short", "America/Toronto")
	if r != "Wed" {
		t.Fatal("Failed to parse dow string")
	}
	r = TimeFormat(ti, "dow_long", "America/Toronto")
	if r != "Wednesday" {
		t.Fatal("Failed to parse dow string")
	}
	r = TimeFormat(ti, "dom", "America/Toronto")
	if r != 1 {
		t.Fatal("Failed to parse dom int")
	}
	r = TimeFormat(ti, "month", "America/Toronto")
	if r != 11 {
		t.Fatal("Failed to parse moy int")
	}
	r = TimeFormat(ti, "month_short", "America/Toronto")
	if r != "Nov" {
		t.Fatal("Failed to parse moy string")
	}
	r = TimeFormat(ti, "month_long", "America/Toronto")
	if r != "November" {
		t.Fatal("Failed to parse moy string")
	}
	r = TimeFormat(ti, "hod", "America/Toronto")
	if r != 19 {
		t.Fatal("Failed to parse hod int")
	}
	r = TimeFormat(ti, "moh", "America/Toronto")
	if r != 32 {
		t.Fatal("Failed to parse moh int")
	}
	r = TimeFormat(ti, "ampm", "America/Toronto")
	if r != "PM" {
		t.Fatal("Failed to parse ampm string")
	}
	r = TimeFormat(ti, "year", "America/Toronto")
	if r != 2023 {
		t.Fatal("Failed to parse year")
	}

	// constants in format
	r = TimeFormat(ti, "nada", "America/Toronto")
	if r != "nada" {
		t.Fatal("Failed to return nil for invalid format")
	}
	r = TimeFormat(ti, "nada Mon", "America/Toronto")
	if r != "nada Wed" {
		t.Fatal("Failed to return nil for invalid format")
	}

	// Test UTC
	r = TimeFormat(ti, "hod", "UTC")
	if r != 23 {
		t.Fatal("Failed to parse hod int in UTC")
	}

	// Test TZ: local, explicit and default (to Local)
	local := time.Local
	defer func() {
		time.Local = local
	}()
	time.Local, _ = time.LoadLocation("America/St_Johns")
	ti = time.UnixMilli(1698881571000)

	testTz := []string{"Local", "America/St_Johns", ""}
	for _, tz := range testTz {
		r = TimeFormat(ti, "hod", tz)
		if r != 21 {
			t.Fatal("Failed to parse hod int in another TZ")
		}
		r = TimeFormat(ti, "moh", tz)
		if r != 2 {
			t.Fatal("Failed to parse moh int in another TZ")
		}
	}
}
