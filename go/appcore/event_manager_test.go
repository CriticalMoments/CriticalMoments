package appcore

import (
	"testing"
	"time"

	datamodel "github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model"
)

func TestSessionStart(t *testing.T) {
	ac, err := testBuildValidTestAppCore(t)
	if err != nil {
		t.Fatal(err)
	}
	err = ac.Start(true)
	if err != nil {
		t.Fatal(err)
	}

	// Test there isn't a session start event before foreground
	latestSessionStart, err := ac.db.LatestEventTimeByName(datamodel.SessionStartBuiltInEvent)
	if err != nil || latestSessionStart != nil {
		t.Fatalf("Unexpected session start before foreground: %v, %v", latestSessionStart, err)
	}

	ac.SendBuiltInEvent(datamodel.AppEnteredForegroundBuiltInEvent)

	// Test there is a session start event after foreground
	firstSessionStart, err := ac.db.LatestEventTimeByName(datamodel.SessionStartBuiltInEvent)
	if err != nil || firstSessionStart == nil {
		t.Fatalf("Expected session start after foreground: %v, %v", latestSessionStart, err)
	}
	if time.Since(*firstSessionStart) > 3*time.Second {
		t.Fatalf("Unexpected session start time: %v", latestSessionStart)
	}

	// confirm session start sent again after background
	ac.SendBuiltInEvent(datamodel.AppEnteredBackgroundBuiltInEvent)
	originalSessionGapDuration := SessionGapDuration
	SessionGapDuration = time.Millisecond * 100
	defer func() {
		SessionGapDuration = originalSessionGapDuration
	}()

	// Another foreground event within time window should not restart session
	ac.SendBuiltInEvent(datamodel.AppEnteredForegroundBuiltInEvent)
	latestSessionStart, err = ac.db.LatestEventTimeByName(datamodel.SessionStartBuiltInEvent)
	if err != nil || latestSessionStart.UnixMilli() != firstSessionStart.UnixMilli() {
		t.Fatalf("Unexpected session restart within time window: %v, %v", latestSessionStart, err)
	}
	ac.SendBuiltInEvent(datamodel.AppEnteredBackgroundBuiltInEvent)

	// Another foreground event outside time window should restart session
	time.Sleep(time.Millisecond * 101)
	ac.SendBuiltInEvent(datamodel.AppEnteredForegroundBuiltInEvent)
	latestSessionStart, err = ac.db.LatestEventTimeByName(datamodel.SessionStartBuiltInEvent)
	if err != nil || latestSessionStart.UnixMilli() == firstSessionStart.UnixMilli() {
		t.Fatalf("Expected session restart after foreground: %v, %v", latestSessionStart, err)
	}

	// confirm sesstion_start_time updated to match (5ms tolerance)
	r, err := ac.propertyRegistry.propertyValue("session_start_time")
	if err != nil {
		t.Fatal(err)
	}
	sessionStartResult := r.(time.Time)
	if err != nil {
		t.Fatal(err)
	}
	diff := sessionStartResult.UnixMilli() - latestSessionStart.UnixMilli()
	if diff > 5 || diff < -5 {
		t.Fatalf("Unexpected session start time: %v", sessionStartResult)
	}
}
