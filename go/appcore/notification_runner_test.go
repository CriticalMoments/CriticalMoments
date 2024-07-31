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
	if plan.UnscheduledNotificationCount() != 5 {
		t.Fatalf("Expected 5 unscheduled notifications, got %d", plan.UnscheduledNotificationCount())
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
	if plan.UnscheduledNotificationCount() != 4 {
		t.Fatalf("Expected 4 unscheduled notification, got %d", plan.UnscheduledNotificationCount())
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
	first := scheduleNotificationWithName(plan.scheduledNotifications, "event2Notification")
	latest := scheduleNotificationWithName(plan.scheduledNotifications, "event1Notification")
	if first == nil || latest == nil {
		t.Fatal("Expected event1 and event2 notifications to be scheduled")
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

	// Test notification 4 with latest-once event instance
	err = ac.SendClientEvent("event4")
	if err != nil {
		t.Fatal(err)
	}
	events, err := ac.db.AllEventTimesByName("event4")
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}
	plan = lb.lastNotificationPlan
	if plan.ScheduledNotificationCount() != 2 {
		t.Fatalf("Expected 2 scheduled notification, got %d", plan.ScheduledNotificationCount())
	}
	sn = scheduleNotificationWithName(plan.scheduledNotifications, "event4Notification")
	if sn == nil {
		t.Fatal("Expected event4 notification to be scheduled")
	}
	// Expect it to be 60 seconds from now, but not more than 50ms off
	expectedTime := time.Now().Add(60 * time.Second)
	diffTime := expectedTime.UnixMilli() - sn.ScheduledAtEpochMilliseconds()
	if diffTime > 50 || diffTime < -50 {
		t.Fatalf("Expected scheduledAt %v, got %v", expectedTime, sn.ScheduledAtEpochMilliseconds())
	}

	// Test event5 with latest-once event instance, but no offset. This is the same as "first".
	expectedTime = time.Now()
	err = ac.SendClientEvent("event5")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(25 * time.Millisecond)
	afterTime := time.Now()
	err = ac.SendClientEvent("event5")
	if err != nil {
		t.Fatal(err)
	}
	events, err = ac.db.AllEventTimesByName("event5")
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 2 {
		t.Fatalf("Expected 2 events, got %d", len(events))
	}
	plan = lb.lastNotificationPlan
	sn = scheduleNotificationWithName(plan.scheduledNotifications, "event5Notification")
	if sn == nil {
		t.Fatal("Expected event5 notification to be scheduled")
	}
	diffTime = expectedTime.UnixMilli() - sn.ScheduledAtEpochMilliseconds()
	// Should be within 20ms of expected time (the first time, since offset is 0)
	if diffTime > 20 || diffTime < -20 {
		t.Fatalf("Expected scheduledAt %v, got %v", expectedTime, sn.ScheduledAtEpochMilliseconds())
	}
	if sn.ScheduledAtEpochMilliseconds() > afterTime.UnixMilli() {
		t.Fatalf("Expected scheduledAt %v, got %v", expectedTime, sn.ScheduledAtEpochMilliseconds())
	}
}

func scheduleNotificationWithName(sns []*ScheduledNotification, name string) *ScheduledNotification {
	for _, sn := range sns {
		if sn.Notification.ID == name {
			return sn
		}
	}
	return nil
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

var allDays = []time.Weekday{
	time.Sunday,
	time.Monday,
	time.Tuesday,
	time.Wednesday,
	time.Thursday,
	time.Friday,
	time.Saturday,
}

func TestDateWindowShift(t *testing.T) {
	torontoTime, err := time.LoadLocation("America/Toronto")
	if err != nil {
		t.Fatal(err)
	}
	// Sunday, 9:19 am Toronto time
	var testTimeEpoch int64 = 1720358361
	testTime := time.Unix(testTimeEpoch, 0).In(torontoTime)

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
	// Set custom time for testing
	customTime := time.Date(2023, time.October, 10, 12, 0, 0, 0, time.UTC)
	customTimePlusBgDelay := customTime.Add(checkTimeDelay)

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
		bgCheckTime          *time.Time
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
		shiftedTime, bgCheckTime := ac.shiftDeliveryTimeForIdealWindow(&test.notification, test.nonIdealDeliveryTime, customTime)
		if (shiftedTime == nil && test.expectedShiftedTime != nil) || (shiftedTime != nil && test.expectedShiftedTime == nil) {
			t.Fatalf("Test %s: Expected shiftedTime %v, but got %v", test.name, test.expectedShiftedTime, shiftedTime)
		}
		if shiftedTime != nil && !shiftedTime.Equal(*test.expectedShiftedTime) {
			t.Fatalf("Test %s: Expected shiftedTime %v, but got %v", test.name, *test.expectedShiftedTime, *shiftedTime)
		}
		if bgCheckTime == nil && test.bgCheckTime != nil || (bgCheckTime != nil && test.bgCheckTime == nil) {
			t.Fatalf("Test %s: Expected bgCheckTime %v, but got %v", test.name, test.bgCheckTime, bgCheckTime)
		}
		if bgCheckTime != nil && !bgCheckTime.Equal(*test.bgCheckTime) {
			t.Fatalf("Test %s: Expected bgCheckTime %v, but got %v", test.name, *test.bgCheckTime, *bgCheckTime)
		}
	}

	runTest(testType{ // add_test_count
		name:                 "valid in ideal window",
		notification:         buildValidNotification(),
		nonIdealDeliveryTime: &customTime,
		expectedShiftedTime:  &customTime,
		bgCheckTime:          nil,
	})

	runTest(testType{ // add_test_count
		name:                 "nil nonIdealDeliveryTime",
		notification:         buildValidNotification(),
		nonIdealDeliveryTime: nil,
		expectedShiftedTime:  nil,
		bgCheckTime:          nil,
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
		// technically could detect invalid condition and not schedule BG, but this is correct time had it been valid condition
		bgCheckTime: &customTimePlusBgDelay,
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
		// technically could detect false condition and not schedule BG, but this is correct time had it been condition which value can change
		bgCheckTime: &customTimePlusBgDelay,
	})

	// Condition passes and in ideal window, but not in filters. Should still be at end of window, not now
	filterFailNotification := buildValidNotification()
	filterFailNotification.DeliveryWindowTODEndMinutes = 1
	runTest(testType{ // add_test_count
		name:                 "filters fail",
		notification:         filterFailNotification,
		nonIdealDeliveryTime: &customTime,
		expectedShiftedTime:  &endOfIdealWindow,
		// bgCheck time would be same as expectedShiftedTime, so no need for bg time
		bgCheckTime: nil,
	})

	// Wait forever with false condition should not schedule at end of window
	// But should schedule BG check time to check if it changes
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
		bgCheckTime:          &customTimePlusBgDelay,
	})
}

func TestLatestOnceDeliveryTimeFromEventList(t *testing.T) {
	eventName := "test"
	offset := 300
	dt := datamodel.DeliveryTime{
		EventName:          &eventName,
		EventOffsetSeconds: &offset,
	}

	// No events
	tm, err := latestOnceEventTimeFromEventList(&dt, []time.Time{})
	if err != nil {
		t.Fatal(err)
	}
	if tm != nil {
		t.Fatal("Expected nil time")
	}

	// Many events at end
	customTime := time.Now()
	times := []time.Time{
		customTime.Add(time.Second),
		customTime.Add(2 * time.Second),
		customTime.Add(3 * time.Second),
		customTime.Add(4 * time.Second),
	}
	tm, err = latestOnceEventTimeFromEventList(&dt, times)
	if err != nil {
		t.Fatal(err)
	}
	expectedTime := times[len(times)-1]
	if !tm.Equal(expectedTime) {
		t.Fatal("Expected last time returned")
	}

	// Offset 0 should return first
	offset = 0
	tm, err = latestOnceEventTimeFromEventList(&dt, times)
	if err != nil {
		t.Fatal(err)
	}
	expectedTime = times[0]
	if !tm.Equal(expectedTime) {
		t.Fatal("Expected first time returned as offset is 0")
	}

	// Many events with gap larger than offset
	offset = 300
	times = []time.Time{
		customTime.Add(time.Second),
		customTime.Add(2 * time.Second),
		customTime.Add(3 * time.Second), // i=2, after this there is a gap of 10 minutes (> offset)
		customTime.Add(10 * time.Minute),
		customTime.Add(11 * time.Minute),
		customTime.Add(12 * time.Minute),
	}
	tm, err = latestOnceEventTimeFromEventList(&dt, times)
	if err != nil {
		t.Fatal(err)
	}
	expectedTime = times[2]
	if !tm.Equal(expectedTime) {
		t.Fatal("Expected last time returned")
	}

	// adding another event shouldn't change the result
	times = append(times, customTime.Add(13*time.Minute))
	tm, err = latestOnceEventTimeFromEventList(&dt, times)
	if err != nil {
		t.Fatal(err)
	}
	if !tm.Equal(expectedTime) {
		t.Fatal("Expected last time returned")
	}
}

func TestNextBackgroundWorkTimeForNotifications(t *testing.T) {
	customTime := time.Date(2023, time.October, 10, 12, 0, 0, 0, time.UTC)
	customTimeBeforeDelay := customTime.Add(checkTimeDelay - time.Minute)
	customTimePlusDay := customTime.Add(24 * time.Hour)

	// nil notification should return nil
	bgCheckTime := bgCheckTimeForIdealDeliveryWindow(nil, customTime, &customTimePlusDay)
	if bgCheckTime != nil {
		t.Fatal("Expected bgCheckTime to be set")
	}

	// 10 mins out should not return BG check time, as the delivery time is before first possible check time (15 minutes in the future)
	notification := datamodel.Notification{
		IdealDevlieryConditions: &datamodel.IdealDevlieryConditions{
			MaxWaitTimeSeconds: 60 * 60,
		},
		DeliveryDaysOfWeek:            allDays,
		DeliveryWindowTODStartMinutes: 0,
		DeliveryWindowTODEndMinutes:   24*60 - 1,
	}
	bgCheckTime = bgCheckTimeForIdealDeliveryWindow(&notification, customTime, &customTimeBeforeDelay)
	if bgCheckTime != nil {
		t.Fatal("Expected bgCheckTime to be nil when delivery time is before first possible check time")
	}

	// Expect 15 mins out when delivery time is past then
	bgCheckTime = bgCheckTimeForIdealDeliveryWindow(&notification, customTime, &customTimePlusDay)
	expectedTime := customTime.Add(checkTimeDelay)
	if bgCheckTime == nil || !bgCheckTime.Equal(expectedTime) {
		t.Fatalf("Expected bgCheckTime to be %v, got %v", expectedTime, bgCheckTime)
	}

	// expect 30 mins out when filters require it, plus 2 minutes buffer
	notification.DeliveryWindowTODStartMinutes = 12*60 + 30 // 12:30, when customTime is 12:00
	bgCheckTime = bgCheckTimeForIdealDeliveryWindow(&notification, customTime, &customTimePlusDay)
	expectedTime = customTime.Add(30*time.Minute + filterTimeBuffer)
	if bgCheckTime == nil || !bgCheckTime.Equal(expectedTime) {
		t.Fatalf("Expected bgCheckTime to be %v, got %v", expectedTime, bgCheckTime)
	}

	// Exepect nil when delivery window pushes bgCheckTime past fallback time
	notification.DeliveryDaysOfWeek = []time.Weekday{time.Sunday}
	bgCheckTime = bgCheckTimeForIdealDeliveryWindow(&notification, customTime, &customTimePlusDay)
	if bgCheckTime != nil {
		t.Fatal("Expected bgCheckTime to be nil when delivery window pushes bgCheckTime past fallback time")
	}
}

func TestTwoIdealTimeBackgroundTimes(t *testing.T) {
	customTime := time.Date(2023, time.October, 10, 12, 0, 0, 0, time.UTC)

	ac, err := buildTestAppCoreWithPath("../cmcore/data_model/test/testdata/notifications/dualIdealNotifications.json", t)
	if err != nil {
		t.Fatal(err)
	}
	err = ac.Start(true)
	if err != nil {
		t.Fatal(err)
	}

	plan, err := ac.generateNotificationPlanForTime(customTime)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.scheduledNotifications) != 0 {
		t.Fatal("Expected scheduledNotifications to be 0 since no trigger events fired")
	}
	if plan.EarliestBgCheckTimeEpochSeconds != 0 {
		t.Fatal("Expected EarliestBgCheckTimeEpochSeconds to be 0 since no trigger events fired")
	}

	// fire event1
	err = ac.SendClientEvent("event1")
	if err != nil {
		t.Fatal(err)
	}
	plan, err = ac.generateNotificationPlanForTime(customTime)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.scheduledNotifications) != 1 {
		t.Fatal("Expected scheduledNotifications to be 1 since event1 fired")
	}
	// 15:00 is first possible time for bg check because of filters, plus 2 minutes buffer
	expectedBgCheckTime := time.Date(2023, time.October, 10, 15, 00, 0, 0, time.UTC).Add(filterTimeBuffer)
	if plan.EarliestBgCheckTimeEpochSeconds != expectedBgCheckTime.Unix() {
		t.Fatalf("Expected EarliestBgCheckTimeEpochSeconds to be %v, got %v", expectedBgCheckTime.Unix(), plan.EarliestBgCheckTimeEpochSeconds)
	}

	// fire event2 which doesn't have filters, so should check bg in 15 mins, and should select earlier of 2 bg check times
	err = ac.SendClientEvent("event2")
	if err != nil {
		t.Fatal(err)
	}
	plan, err = ac.generateNotificationPlanForTime(customTime)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.scheduledNotifications) != 2 {
		t.Fatal("Expected scheduledNotifications to be 2 since event2 fired")
	}
	expectedBgCheckTime = customTime.Add(checkTimeDelay)
	if plan.EarliestBgCheckTimeEpochSeconds != expectedBgCheckTime.Unix() {
		t.Fatalf("Expected EarliestBgCheckTimeEpochSeconds to be %v, got %v", expectedBgCheckTime.Unix(), plan.EarliestBgCheckTimeEpochSeconds)
	}
}
