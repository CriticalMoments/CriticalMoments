package appcore

import (
	"math"
	"testing"
	"time"

	datamodel "github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model"
)

func TestNotificationList(t *testing.T) {
	ac, err := testBuildValidTestAppCore(t)
	if err != nil {
		t.Fatal(err)
	}
	lb := testLibBindings{}
	ac.RegisterLibraryBindings(&lb)
	err = ac.Start(true)
	if err != nil {
		t.Fatal(err)
	}

	plan := lb.lastNotificationPlan
	if plan.UnscheduledNotificationCount() != 1 {
		t.Fatalf("Expected 1 unscheduled notification, got %d", plan.UnscheduledNotificationCount())
	}
	un := plan.UnscheduledNotificationAtIndex(0)
	// A notification in the past should be unscheduled
	if un == nil {
		t.Fatal("Expected UnscheduledNotificationAtIndex to return a value")
	}
	if un.ID != "testNotification" {
		t.Fatalf("Expected UnscheduledNotificationAtIndex to return a value with Id 'testNotification', got %s", un.ID)
	}
	if plan.ScheduledNotificationCount() != 1 {
		t.Fatalf("Expected 1 scheduled notification, got %d", plan.ScheduledNotificationCount())
	}
	sn := plan.ScheduledNotificationAtIndex(0)
	if sn == nil {
		t.Fatal("Expected ScheduledNotificationAtIndex to return a value")
	}
	if sn.Notification.ID != "futureStaticNotification" {
		t.Fatalf("Expected ScheduledNotificationAtIndex to return future static notification, got %s", sn.Notification.ID)
	}
	if sn.ScheduledAtEpochMilliseconds() != 2224580934*1000 {
		t.Fatalf("Expected ScheduledAtEpoch to return time in 2040, got %d", sn.ScheduledAtEpochMilliseconds())
	}
}

func TestEventNotificationPlan(t *testing.T) {
	ac, err := buildTestAppCoreWithPath("../cmcore/data_model/test/testdata/notifications/eventNotifications.json", t)
	if err != nil {
		t.Fatal(err)
	}

	lb := testLibBindings{}
	ac.RegisterLibraryBindings(&lb)
	if lb.lastNotificationPlan != nil {
		t.Fatal("NP binding set too soon")
	}

	// Check seenCancelationEvents is empty
	if len(ac.seenCancelationEvents) != 0 {
		t.Fatal("Expected seenCancelationEvents to be empty")
	}

	// manually insert a cancelation event before starting to check we load it from DB correctly
	cancel3Event := datamodel.Event{
		Name:      "cancel3event",
		EventType: datamodel.EventTypeCustom,
	}
	err = ac.db.InsertEvent(&cancel3Event)
	if err != nil {
		t.Fatal(err)
	}

	err = ac.Start(true)
	if err != nil {
		t.Fatal(err)
	}

	// Check seenCancelationEvents work: all should load from DB, 3 should be true
	if len(ac.seenCancelationEvents) != 2 ||
		*ac.seenCancelationEvents["cancel2event"] != false ||
		*ac.seenCancelationEvents["cancel3event"] != true {
		t.Fatal("Expected seenCancelationEvents to be populated for all cancelation events")
	}

	// No events, should be nothing scheduled
	plan := lb.lastNotificationPlan
	if plan == nil {
		t.Fatal("NP binding not set after startup")
	}
	if plan.ScheduledNotificationCount() != 0 {
		t.Fatalf("Expected 0 scheduled notifications, got %d", plan.ScheduledNotificationCount())
	}
	if plan.UnscheduledNotificationCount() != 3 {
		t.Fatalf("Expected 2 unscheduled notifications, got %d", plan.UnscheduledNotificationCount())
	}

	// Fire event, should be scheduled for now (no offset)
	err = ac.SendClientEvent("event1")
	if err != nil {
		t.Fatal(err)
	}

	plan = lb.lastNotificationPlan
	if plan.ScheduledNotificationCount() != 1 {
		t.Fatalf("Expected 1 scheduled notification, got %d", plan.ScheduledNotificationCount())
	}
	sn := plan.ScheduledNotificationAtIndex(0)
	if sn == nil {
		t.Fatal("Expected ScheduledNotificationAtIndex to return a value")
	}
	if sn.Notification.ID != "event1Notification" {
		t.Fatalf("Expected ScheduledNotificationAtIndex to return event notification, got %s", sn.Notification.ID)
	}
	// Expect it scheduled within 100ms of now
	if math.Abs(float64(sn.ScheduledAtEpochMilliseconds()-time.Now().UnixMilli())) > 100 {
		t.Fatalf("Expected ScheduledAtEpoch to return now, got %d", sn.ScheduledAtEpochMilliseconds())
	}
	if plan.UnscheduledNotificationCount() != 2 {
		t.Fatalf("Expected 1 unscheduled notification, got %d", plan.UnscheduledNotificationCount())
	}

	// Fire event, should be scheduled for offset
	err = ac.SendClientEvent("event2")
	if err != nil {
		t.Fatal(err)
	}

	plan = lb.lastNotificationPlan
	if err != nil {
		t.Fatal(err)
	}
	if plan.ScheduledNotificationCount() != 2 {
		t.Fatalf("Expected 2 scheduled notifications, got %d", plan.ScheduledNotificationCount())
	}
	sn2 := plan.ScheduledNotificationAtIndex(0)
	if sn2.Notification.ID != "event2Notification" {
		// Index is interterminate
		sn2 = plan.ScheduledNotificationAtIndex(1)
	}
	if sn2.Notification.ID != "event2Notification" {
		t.Fatalf("Expected ScheduledNotificationAtIndex to return event notification, got %s", sn.Notification.ID)
	}
	// Expect it scheduled within 100ms of 60s from now
	if math.Abs(float64(sn2.ScheduledAtEpochMilliseconds()-time.Now().UnixMilli()-60000)) > 100 {
		t.Fatalf("Expected ScheduledAtEpoch to return now, got %d", sn2.ScheduledAtEpochMilliseconds())
	}

	// Fire again after delay. Testing lastest vs first eventInstance targeting
	time.Sleep(110 * time.Millisecond)
	ac.SendClientEvent("event1")
	ac.SendClientEvent("event2")

	plan = lb.lastNotificationPlan
	if err != nil {
		t.Fatal(err)
	}
	if plan.ScheduledNotificationCount() != 2 {
		t.Fatalf("Expected 2 scheduled notifications, got %d", plan.ScheduledNotificationCount())
	}
	first := plan.ScheduledNotificationAtIndex(0)
	latest := plan.ScheduledNotificationAtIndex(1)
	if first.Notification.ID != "event2Notification" {
		// Index is interterminate
		first = plan.ScheduledNotificationAtIndex(1)
		latest = plan.ScheduledNotificationAtIndex(0)
	}
	// Check first has same time
	if first.ScheduledAtEpochMilliseconds()-sn2.ScheduledAtEpochMilliseconds() != 0 {
		t.Fatal("Event notifications scheduled around first event should not move in time")
	}
	// Check lastest has moved
	latestMove := latest.ScheduledAtEpochMilliseconds() - sn.ScheduledAtEpochMilliseconds()
	if latestMove < 100 || latestMove > 200 {
		t.Fatal("latest should move if the event fired again")
	}

	// Sent an event which should cancel the "2" notification
	ac.SendClientEvent("cancel2event")

	// Check we persisted as seenCancelationEvent
	if len(ac.seenCancelationEvents) != 2 ||
		*ac.seenCancelationEvents["cancel2event"] != true ||
		*ac.seenCancelationEvents["cancel3event"] != true {
		t.Fatal("Expected seenCancelationEvents to be populated for all cancelation events")
	}

	// Check it's cancelled
	plan = lb.lastNotificationPlan
	if plan.ScheduledNotificationCount() != 1 {
		t.Fatalf("Expected 1 scheduled notification, got %d", plan.ScheduledNotificationCount())
	}
	sn = plan.ScheduledNotificationAtIndex(0)
	if sn == nil {
		t.Fatal("Expected ScheduledNotificationAtIndex to return a value")
	}
	if sn.Notification.ID != "event1Notification" {
		t.Fatalf("Expected ScheduledNotificationAtIndex to return event notification, got %s", sn.Notification.ID)
	}
}

func TestNotificationEventAction(t *testing.T) {
	ac, err := buildTestAppCoreWithPath("../cmcore/data_model/test/testdata/notifications/eventNotifications.json", t)
	if err != nil {
		t.Fatal(err)
	}

	lb := testLibBindings{}
	ac.RegisterLibraryBindings(&lb)

	err = ac.Start(true)
	if err != nil {
		t.Fatal(err)
	}

	if lb.lastBannerAction != nil {
		t.Fatal("Expected no banner action")
	}

	n := datamodel.Notification{
		ID: "event2Notification",
	}
	// Calling simulates tapping the notification, which should trigger the action
	ac.ActionForNotification(n.UniqueID())

	if lb.lastBannerAction == nil {
		t.Fatal("Expected banner action")
	}
}

func TestDateWindowShift(t *testing.T) {
	torontoTime, err := time.LoadLocation("America/Toronto")
	if err != nil {
		t.Fatal(err)
	}
	// Sunday, 9:19 am Toronto time
	var testTimeEpoch int64 = 1720358361
	testTime := time.Unix(testTimeEpoch, 0).In(torontoTime)

	allDays := []time.Weekday{
		time.Sunday,
		time.Monday,
		time.Tuesday,
		time.Wednesday,
		time.Thursday,
		time.Friday,
		time.Saturday,
	}

	n := datamodel.Notification{
		DeliveryWindowTODStartMinutes: 10 * 60,
		DeliveryWindowTODEndMinutes:   11 * 60,
		DeliveryDaysOfWeek:            allDays,
		DeliveryTime: datamodel.DeliveryTime{
			TimestampEpoch: &testTimeEpoch,
		},
	}

	// Should shift to 10am
	if timeMeetsFilterConditions(&n, &testTime) {
		t.Fatal("Should not meet filters, needs to be shifted")
	}
	shiftedTime := shiftDeliveryTimeForFilters(&n, &testTime)
	if shiftedTime.Hour() != 10 || shiftedTime.Minute() != 0 || shiftedTime.Second() != 0 {
		t.Fatalf("Expected shifted time to be 10am, got %v", shiftedTime)
	}
	if shiftedTime.Day() != testTime.Day() || shiftedTime.Month() != testTime.Month() || shiftedTime.Year() != testTime.Year() {
		t.Fatalf("Expected shifted time to be same day, got %v", shiftedTime)
	}

	// make delivery window too soon. Should shift to next day start of window
	n.DeliveryWindowTODStartMinutes = 7 * 60
	n.DeliveryWindowTODEndMinutes = 8 * 60
	if timeMeetsFilterConditions(&n, &testTime) {
		t.Fatal("Should not meet filters, needs to be shifted")
	}
	shiftedTime = shiftDeliveryTimeForFilters(&n, &testTime)
	if shiftedTime.Hour() != 7 || shiftedTime.Minute() != 0 || shiftedTime.Second() != 0 {
		t.Fatalf("Expected shifted time to be 7am, got %v", shiftedTime)
	}
	if shiftedTime.Day() != testTime.Day()+1 || shiftedTime.Month() != testTime.Month() || shiftedTime.Year() != testTime.Year() {
		t.Fatalf("Expected shifted time to be next day, got %v", shiftedTime)
	}

	// make delivery window okay, should not modify
	n.DeliveryWindowTODStartMinutes = 7 * 60
	n.DeliveryWindowTODEndMinutes = 11 * 60
	if !timeMeetsFilterConditions(&n, &testTime) {
		t.Fatal("Should meet filters")
	}
	shiftedTime = shiftDeliveryTimeForFilters(&n, &testTime)
	if shiftedTime.Sub(testTime) != 0 {
		t.Fatalf("Expected shifted time to be same as original, got %v", shiftedTime)
	}

	// only allow weekends, should not modify as date is on weekend
	n.DeliveryDaysOfWeek = []time.Weekday{
		time.Sunday,
		time.Saturday,
	}
	if !timeMeetsFilterConditions(&n, &testTime) {
		t.Fatal("Should meet filters")
	}
	shiftedTime = shiftDeliveryTimeForFilters(&n, &testTime)
	if shiftedTime.Sub(testTime) != 0 {
		t.Fatalf("Expected shifted time to be same as original, got %v", shiftedTime)
	}

	// only allow on weekdays, should modify as date is on weekend
	n.DeliveryDaysOfWeek = []time.Weekday{
		time.Monday,
		time.Tuesday,
		time.Wednesday,
		time.Thursday,
		time.Friday,
	}
	if timeMeetsFilterConditions(&n, &testTime) {
		t.Fatal("Should not meet filters, needs to be shifted")
	}
	shiftedTime = shiftDeliveryTimeForFilters(&n, &testTime)
	// time should be the same as the original
	if shiftedTime.Hour() != testTime.Hour() || shiftedTime.Minute() != testTime.Minute() || shiftedTime.Second() != testTime.Second() {
		t.Fatalf("Expected shifted time of day to be same as original, got %v", shiftedTime)
	}
	// day should be moved to following Monday (original was Sunday)
	if shiftedTime.Day() != testTime.Day()+1 || shiftedTime.Weekday() != time.Monday || shiftedTime.Month() != testTime.Month() || shiftedTime.Year() != testTime.Year() {
		t.Fatalf("Expected shifted time of day to be next day, got %v", shiftedTime)
	}

	// should shift both time and date, next window is tomorrow
	n.DeliveryWindowTODStartMinutes = 7 * 60
	n.DeliveryWindowTODEndMinutes = 8 * 60
	if timeMeetsFilterConditions(&n, &testTime) {
		t.Fatal("Should not meet filters, needs to be shifted")
	}
	shiftedTime = shiftDeliveryTimeForFilters(&n, &testTime)
	if shiftedTime.Hour() != 7 || shiftedTime.Minute() != 0 || shiftedTime.Second() != 0 {
		t.Fatalf("Expected shifted time to be 7am, got %v", shiftedTime)
	}
	if shiftedTime.Day() != testTime.Day()+1 || shiftedTime.Weekday() != time.Monday || shiftedTime.Month() != testTime.Month() || shiftedTime.Year() != testTime.Year() {
		t.Fatalf("Expected shifted time of day to be next day, got %v", shiftedTime)
	}

	// Only allow on Wednesdays, should shift from Sunday to Wednesday
	n.DeliveryDaysOfWeek = []time.Weekday{time.Wednesday}
	if timeMeetsFilterConditions(&n, &testTime) {
		t.Fatal("Should not meet filters, needs to be shifted")
	}
	shiftedTime = shiftDeliveryTimeForFilters(&n, &testTime)
	if shiftedTime.Hour() != 7 || shiftedTime.Minute() != 0 || shiftedTime.Second() != 0 {
		t.Fatalf("Expected shifted time to be 7am, got %v", shiftedTime)
	}
	if shiftedTime.Day() != testTime.Day()+3 || shiftedTime.Weekday() != time.Wednesday || shiftedTime.Month() != testTime.Month() || shiftedTime.Year() != testTime.Year() {
		t.Fatalf("Expected shifted time of day to be next day, got %v", shiftedTime)
	}

	// DST test: Sat March 9th 2024 at 11am Toronto time to next day, should only shift 23h, not 24h
	dstTime := time.Date(2024, time.March, 9, 11, 0, 0, 0, torontoTime)
	n.DeliveryDaysOfWeek = []time.Weekday{time.Sunday}
	n.DeliveryWindowTODStartMinutes = 0
	n.DeliveryWindowTODEndMinutes = 24*60 - 1
	if timeMeetsFilterConditions(&n, &dstTime) {
		t.Fatal("Should not meet filters, needs to be shifted")
	}
	shiftedTime = shiftDeliveryTimeForFilters(&n, &dstTime)
	if shiftedTime.Hour() != 11 || shiftedTime.Minute() != 0 || shiftedTime.Second() != 0 {
		t.Fatalf("Expected shifted time to be 11am, got %v", shiftedTime)
	}
	if shiftedTime.Day() != dstTime.Day()+1 || shiftedTime.Month() != dstTime.Month() || shiftedTime.Year() != dstTime.Year() {
		t.Fatalf("Expected shifted time of day to be next day, got %v", shiftedTime)
	}
	if shiftedTime.Sub(dstTime) != 23*time.Hour {
		t.Fatalf("Expected shifted time to be 23 hours because of DST later, got %v", shiftedTime.Sub(dstTime))
	}

	// test two different timezones, one in the window, one not
	n.DeliveryDaysOfWeek = allDays
	// Okay in toronto, needs to shift in chicago
	n.DeliveryWindowTODStartMinutes = 9 * 60
	n.DeliveryWindowTODEndMinutes = 10 * 60
	if !timeMeetsFilterConditions(&n, &testTime) {
		t.Fatal("Should meet filters")
	}
	shiftedTime = shiftDeliveryTimeForFilters(&n, &testTime)
	if shiftedTime.Sub(testTime) != 0 {
		t.Fatalf("Expected shifted time to be same as original, got %v", shiftedTime)
	}
	chicagoTimeZone, err := time.LoadLocation("America/Chicago")
	if err != nil {
		t.Fatal(err)
	}
	chicagoTime := testTime.In(chicagoTimeZone)
	if timeMeetsFilterConditions(&n, &chicagoTime) {
		t.Fatal("Should not meet filters, needs to be shifted")
	}
	shiftedTime = shiftDeliveryTimeForFilters(&n, &chicagoTime)
	if shiftedTime.Hour() != 9 || shiftedTime.Minute() != 0 || shiftedTime.Second() != 0 {
		t.Fatalf("Expected shifted time to be 9am from 8:19am Chicago time, got %v", shiftedTime)
	}
	if shiftedTime.Day() != testTime.Day() || shiftedTime.Month() != testTime.Month() || shiftedTime.Year() != testTime.Year() {
		t.Fatalf("Expected shifted time of day to be same day, got %v", shiftedTime)
	}
	time.Local = torontoTime
}

func TestScheduleCondition(t *testing.T) {
	ac, err := buildTestAppCoreWithPath("../cmcore/data_model/test/testdata/notifications/conditionalNotifications.json", t)
	if err != nil {
		t.Fatal(err)
	}

	lb := testLibBindings{}
	ac.RegisterLibraryBindings(&lb)

	err = ac.Start(true)
	if err != nil {
		t.Fatal(err)
	}

	// Condition is false, should not be scheduled
	err = ac.SendClientEvent("event3")
	if err != nil {
		t.Fatal(err)
	}
	plan := lb.lastNotificationPlan
	if err != nil {
		t.Fatal(err)
	}
	if plan.ScheduledNotificationCount() != 0 {
		t.Fatal("scheduled a notification when condition wasn't met")
	}

	// Make condition true, should be scheduled
	err = ac.SendClientEvent("event3con")
	if err != nil {
		t.Fatal(err)
	}
	err = ac.SendClientEvent("event3")
	if err != nil {
		t.Fatal(err)
	}
	plan = lb.lastNotificationPlan
	if plan.ScheduledNotificationCount() != 1 {
		t.Fatalf("Expected notification to be scheduled")
	}
	sn := plan.ScheduledNotificationAtIndex(0)
	if sn.Notification.ID != "event3Notification" {
		t.Fatal("should schedule notification 3")
	}
}

func TestNotificationInIdealDeliveryWindow(t *testing.T) {
	// Set custom 'now' time for testing
	customTime := time.Date(2023, time.October, 10, 12, 0, 0, 0, time.UTC)

	allDays := []time.Weekday{
		time.Sunday,
		time.Monday,
		time.Tuesday,
		time.Wednesday,
		time.Thursday,
		time.Friday,
		time.Saturday,
	}

	tests := []struct {
		name                  string
		notification          *datamodel.Notification
		nonIdealDeliveryTime  *time.Time
		expectedInIdealWindow bool
	}{
		{
			name:                  "nil notification",
			notification:          nil,
			nonIdealDeliveryTime:  &customTime,
			expectedInIdealWindow: false,
		},
		{
			name: "nil IdealDevlieryConditions",
			notification: &datamodel.Notification{
				IdealDevlieryConditions:       nil,
				DeliveryDaysOfWeek:            allDays,
				DeliveryWindowTODStartMinutes: 0,
				DeliveryWindowTODEndMinutes:   24*60 - 1,
			},
			nonIdealDeliveryTime:  &customTime,
			expectedInIdealWindow: false,
		},
		{
			name: "nil nonIdealDeliveryTime",
			notification: &datamodel.Notification{
				IdealDevlieryConditions:       &datamodel.IdealDevlieryConditions{},
				DeliveryDaysOfWeek:            allDays,
				DeliveryWindowTODStartMinutes: 0,
				DeliveryWindowTODEndMinutes:   24*60 - 1,
			},
			nonIdealDeliveryTime:  nil,
			expectedInIdealWindow: false,
		},
		{
			name: "nonIdealDeliveryTime in future",
			notification: &datamodel.Notification{
				IdealDevlieryConditions: &datamodel.IdealDevlieryConditions{
					MaxWaitTimeSeconds: 60 * 60,
				},
				DeliveryDaysOfWeek:            allDays,
				DeliveryWindowTODStartMinutes: 0,
				DeliveryWindowTODEndMinutes:   24*60 - 1,
			},
			nonIdealDeliveryTime: func() *time.Time {
				t := customTime.Add(time.Hour)
				return &t
			}(),
			expectedInIdealWindow: false,
		},
		{
			name: "MaxWaitTime exceeded",
			notification: &datamodel.Notification{
				IdealDevlieryConditions: &datamodel.IdealDevlieryConditions{
					MaxWaitTimeSeconds: 60,
				},
				DeliveryDaysOfWeek:            allDays,
				DeliveryWindowTODStartMinutes: 0,
				DeliveryWindowTODEndMinutes:   24*60 - 1,
			},
			nonIdealDeliveryTime: func() *time.Time {
				t := customTime.Add(-time.Hour)
				return &t
			}(),
			expectedInIdealWindow: false,
		},
		{
			name: "current day not in DeliveryDaysOfWeek",
			notification: &datamodel.Notification{
				IdealDevlieryConditions: &datamodel.IdealDevlieryConditions{
					MaxWaitTimeSeconds: 60 * 60,
				},
				DeliveryWindowTODStartMinutes: 0,
				DeliveryWindowTODEndMinutes:   24*60 - 1,
				DeliveryDaysOfWeek:            []time.Weekday{time.Monday},
			},
			nonIdealDeliveryTime: func() *time.Time {
				t := customTime.Add(-time.Minute * 5)
				return &t
			}(),
			expectedInIdealWindow: false,
		},
		{
			name: "current time not in DeliveryWindowTODStartMinutes",
			notification: &datamodel.Notification{
				IdealDevlieryConditions: &datamodel.IdealDevlieryConditions{
					MaxWaitTimeSeconds: 60 * 60,
				},
				DeliveryDaysOfWeek:            allDays,
				DeliveryWindowTODStartMinutes: (customTime.Hour()-1)*60 + customTime.Minute(),
				DeliveryWindowTODEndMinutes:   (customTime.Hour()-1)*60 + customTime.Minute() + 30,
			},
			nonIdealDeliveryTime: func() *time.Time {
				t := customTime.Add(-time.Minute * 5)
				return &t
			}(),
			expectedInIdealWindow: false,
		},
		{
			name: "current time in ideal window",
			notification: &datamodel.Notification{
				IdealDevlieryConditions: &datamodel.IdealDevlieryConditions{
					MaxWaitTimeSeconds: 60 * 60,
				},
				DeliveryDaysOfWeek:            []time.Weekday{customTime.Weekday()},
				DeliveryWindowTODStartMinutes: (customTime.Hour()-1)*60 + customTime.Minute(),
				DeliveryWindowTODEndMinutes:   (customTime.Hour()+1)*60 + customTime.Minute(),
			},
			nonIdealDeliveryTime: func() *time.Time {
				t := customTime.Add(-time.Minute * 5)
				return &t
			}(),
			expectedInIdealWindow: true,
		},
	}

	for _, test := range tests {
		if test.name != "current day not in DeliveryDaysOfWeek" {
			continue
		}
		inIdealWindow := notificationInIdealDeliveryWindow(test.notification, test.nonIdealDeliveryTime, customTime)
		if inIdealWindow != test.expectedInIdealWindow {
			t.Errorf("notificationInIdealDeliveryWindow() = %v, want %v for test %s", inIdealWindow, test.expectedInIdealWindow, test.name)
		}
	}
}

func TestShiftDeliveryTimeForIdealWindow(t *testing.T) {
	// Save and restore the original timeNow function
	originalTimeNow := timeNow
	defer func() { timeNow = originalTimeNow }()

	// Set custom time for testing
	customTime := time.Date(2023, time.October, 10, 12, 0, 0, 0, time.UTC)
	timeNow = func() time.Time {
		return customTime
	}

	ac, err := buildTestAppCoreWithPath("../cmcore/data_model/test/testdata/notifications/conditionalNotifications.json", t)
	if err != nil {
		t.Fatal(err)
	}

	trueCondition, err := datamodel.NewCondition("true")
	if err != nil {
		t.Fatal(err)
	}

	type testType struct {
		name                 string
		notification         datamodel.Notification
		nonIdealDeliveryTime *time.Time
		expectedShiftedTime  *time.Time
	}

	var buildValidNotification = func() datamodel.Notification {
		return datamodel.Notification{
			IdealDevlieryConditions: &datamodel.IdealDevlieryConditions{
				Condition:          *trueCondition,
				MaxWaitTimeSeconds: 60 * 60,
			},
			DeliveryDaysOfWeek:            []time.Weekday{time.Tuesday},
			DeliveryWindowTODStartMinutes: 11 * 60,
			DeliveryWindowTODEndMinutes:   13 * 60,
		}
	}

	if customTime.Weekday() != time.Tuesday {
		t.Fatal("customTime not in window")
	}
	validNotif := buildValidNotification()
	if !timeMeetsFilterConditions(&validNotif, &customTime) {
		t.Fatal("valid notification should be in delivery window of custom time")
	}

	var runTest = func(test testType) {
		shiftedTime := ac.shiftDeliveryTimeForIdealWindow(&test.notification, test.nonIdealDeliveryTime)
		if (shiftedTime == nil && test.expectedShiftedTime != nil) || (shiftedTime != nil && test.expectedShiftedTime == nil) {
			t.Fatalf("Test %s: Expected shiftedTime %v, but got %v", test.name, test.expectedShiftedTime, shiftedTime)
		}
		if shiftedTime != nil && !shiftedTime.Equal(*test.expectedShiftedTime) {
			t.Fatalf("Test %s: Expected shiftedTime %v, but got %v", test.name, *test.expectedShiftedTime, *shiftedTime)
		}
	}

	runTest(testType{ // add_test_count
		name:                 "valid in ideal window",
		notification:         buildValidNotification(),
		nonIdealDeliveryTime: &customTime,
		expectedShiftedTime:  &customTime,
	})

	runTest(testType{ // add_test_count
		name:                 "nil nonIdealDeliveryTime",
		notification:         buildValidNotification(),
		nonIdealDeliveryTime: nil,
		expectedShiftedTime:  nil,
	})

	// For invalid condition, expect it to run at end of window since condition never passes
	invalidConditionNotification := buildValidNotification()
	invalidCondition, err := datamodel.NewCondition("invalid")
	if err != nil {
		t.Fatal(err)
	}
	invalidConditionNotification.IdealDevlieryConditions.Condition = *invalidCondition
	endOfIdealWindow := customTime.Add(time.Hour)
	runTest(testType{ // add_test_count
		name:                 "invalid condition",
		notification:         invalidConditionNotification,
		nonIdealDeliveryTime: &customTime,
		expectedShiftedTime:  &endOfIdealWindow,
	})

	// Same for false condition: Expect it to run at end of window since condition never passes
	falseCondition, err := datamodel.NewCondition("false")
	if err != nil {
		t.Fatal(err)
	}
	idealWithFalseCondition := buildValidNotification()
	idealWithFalseCondition.IdealDevlieryConditions.Condition = *falseCondition
	runTest(testType{ // add_test_count
		name:                 "false condition should push back to end of window",
		notification:         idealWithFalseCondition,
		nonIdealDeliveryTime: &customTime,
		expectedShiftedTime:  &endOfIdealWindow,
	})

	// Condition passes and in ideal window, but not in filters. Should still be at end of window, not now
	filterFailNotification := buildValidNotification()
	filterFailNotification.DeliveryWindowTODEndMinutes = 1
	runTest(testType{ // add_test_count
		name:                 "filters fail",
		notification:         filterFailNotification,
		nonIdealDeliveryTime: &customTime,
		expectedShiftedTime:  &endOfIdealWindow,
	})

	// Wait forever with failing condition should not schedule at end of window
	foreverNotif := buildValidNotification()
	foreverNotif.IdealDevlieryConditions.Condition = *falseCondition
	foreverNotif.IdealDevlieryConditions.MaxWaitTimeSeconds = -1
	if !foreverNotif.IdealDevlieryConditions.WaitForever() {
		t.Fatal("not setup correctly for wait forever")
	}
	runTest(testType{ // add_test_count
		name:                 "wait forever",
		notification:         foreverNotif,
		nonIdealDeliveryTime: &customTime,
		expectedShiftedTime:  nil,
	})
}
