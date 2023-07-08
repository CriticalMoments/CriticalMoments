package conditions

import (
	"testing"
	"time"
)

func TestNowFunction(t *testing.T) {
	// time library was written
	if now() < 1688764455000 {
		t.Fatal("now() is in the past")
	}

	// time in 2050
	if now() > 2540841255000 {
		t.Fatal("now() is too far in the future")
	}
}

func TestDurationFunction(t *testing.T) {
	if seconds(1) != 1000 || seconds(9) != 9000 {
		t.Fatal("seconds function not reutrning in milliseconds")
	}

	if minutes(1) != seconds(60) ||
		minutes(2) != 2*minutes(1) {
		t.Fatal("minutes function not 60 seconds")
	}
	if hours(1) != minutes(60) ||
		hours(2) != 2*hours(1) {
		t.Fatal("hours function not 60 mins")
	}
	if days(1) != hours(24) ||
		days(2) != 2*days(1) {
		t.Fatal("days function not 24h")
	}
}

func TestParseDatetime(t *testing.T) {
	r, err := parseDatetime("2006-01-02T15:04:05.9997+07:00")
	if err != nil || r != 1136189045999 {
		t.Fatal("Datetime parsing failed")
	}

	// Test timezone parsing, 7h offset
	r2, err := parseDatetime("2006-01-02T15:04:05.9997+00:00")
	if err != nil || r2 != 1136214245999 || r2-r != 1000*60*60*7 {
		t.Fatal("Datetime parsing failed")
	}

	// Ensure works without decimal seconds
	r3, err := parseDatetime("2006-01-02T15:04:05+00:00")
	if err != nil || r3 != 1136214245000 {
		t.Fatal("Datetime parsing failed")
	}

	// Ensure works with Z for UTC
	r4, err := parseDatetime("2006-01-02T15:04:05Z")
	if err != nil || r3 != r4 {
		t.Fatal("Date parsing with TZ failed")
	}

	// Date Only
	r5, err := parseDatetime("2000-01-29Z")
	r6, err2 := parseDatetime("2000-01-29+00:00")
	if err != nil || err2 != nil || r5 != 949104000000 || r5 != r6 {
		t.Fatal("Date parsing failed")
	}

	// No timezone provided, should parse in local TZ
	local := time.Local
	defer func() {
		time.Local = local
	}()
	time.Local, _ = time.LoadLocation("America/Toronto")
	r7, err := parseDatetime("2023-01-01T11:11:11")   // -5 from local offset
	r8, err2 := parseDatetime("2023-01-01T16:11:11Z") // UTC adjusted by hand
	if err != nil || err2 != nil || r7 != r8 {
		t.Fatal("Datetime parsing failed for local TZ")
	}

	// No TZ provided should parse in local date TZ
	r9, err := parseDatetime("2023-01-01")             // -5 from local offset
	r10, err2 := parseDatetime("2023-01-01T05:00:00Z") // UTC adjusted by hand
	if err != nil || err2 != nil || r9 != r10 {
		t.Fatal("Date parsing failed for local TZ")
	}

	time.Local = local
}
