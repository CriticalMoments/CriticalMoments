package conditions

import (
	"testing"
	"time"
)

func TestNowFunction(t *testing.T) {
	// time library was written
	if Now() < 1688764455000 {
		t.Fatal("Now() is in the past")
	}

	// time in 2050
	if Now() > 2540841255000 {
		t.Fatal("Now() is too far in the future")
	}
}

func TestDurationFunction(t *testing.T) {
	if Seconds(1) != 1000 || Seconds(9) != 9000 {
		t.Fatal("Seconds function not reutrning in milliSeconds")
	}

	if Minutes(1) != Seconds(60) ||
		Minutes(2) != 2*Minutes(1) {
		t.Fatal("Minutes function not 60 Seconds")
	}
	if Hours(1) != Minutes(60) ||
		Hours(2) != 2*Hours(1) {
		t.Fatal("Hours function not 60 mins")
	}
	if Days(1) != Hours(24) ||
		Days(2) != 2*Days(1) {
		t.Fatal("Days function not 24h")
	}
}

func TestParseDatetime(t *testing.T) {
	r, err := ParseDatetime("2006-01-02T15:04:05.9997+07:00")
	if err != nil || r != 1136189045999 {
		t.Fatal("Datetime parsing failed")
	}

	// Test timezone parsing, 7h offset
	r2, err := ParseDatetime("2006-01-02T15:04:05.9997+00:00")
	if err != nil || r2 != 1136214245999 || r2-r != 1000*60*60*7 {
		t.Fatal("Datetime parsing failed")
	}

	// Ensure works without decimal Seconds
	r3, err := ParseDatetime("2006-01-02T15:04:05+00:00")
	if err != nil || r3 != 1136214245000 {
		t.Fatal("Datetime parsing failed")
	}

	// Ensure works with Z for UTC
	r4, err := ParseDatetime("2006-01-02T15:04:05Z")
	if err != nil || r3 != r4 {
		t.Fatal("Date parsing with TZ failed")
	}

	// Date Only
	r5, err := ParseDatetime("2000-01-29Z")
	r6, err2 := ParseDatetime("2000-01-29+00:00")
	if err != nil || err2 != nil || r5 != 949104000000 || r5 != r6 {
		t.Fatal("Date parsing failed")
	}

	// No timezone provided, should parse in local TZ
	local := time.Local
	defer func() {
		time.Local = local
	}()
	time.Local, _ = time.LoadLocation("America/Toronto")
	r7, err := ParseDatetime("2023-01-01T11:11:11")   // -5 from local offset
	r8, err2 := ParseDatetime("2023-01-01T16:11:11Z") // UTC adjusted by hand
	if err != nil || err2 != nil || r7 != r8 {
		t.Fatal("Datetime parsing failed for local TZ")
	}

	// No TZ provided should parse in local date TZ
	r9, err := ParseDatetime("2023-01-01")             // -5 from local offset
	r10, err2 := ParseDatetime("2023-01-01T05:00:00Z") // UTC adjusted by hand
	if err != nil || err2 != nil || r9 != r10 {
		t.Fatal("Date parsing failed for local TZ")
	}

	time.Local = local
}
