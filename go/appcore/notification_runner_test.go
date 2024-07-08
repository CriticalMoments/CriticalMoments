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

	// TODO_P0: test bg execution
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

	err = ac.Start(true)
	if err != nil {
		t.Fatal(err)
	}

	// No events, should be nothing scheduled
	plan := lb.lastNotificationPlan
	if plan == nil {
		t.Fatal("NP binding not set after startup")
	}
	if plan.ScheduledNotificationCount() != 0 {
		t.Fatalf("Expected 0 scheduled notifications, got %d", plan.ScheduledNotificationCount())
	}
	if plan.UnscheduledNotificationCount() != 2 {
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
	if plan.UnscheduledNotificationCount() != 1 {
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
	// Check cache is working
	if *ac.seenCancelationEvents["cancel2event"] != true {
		t.Fatal("Expected cancel2event to be in seenCancelationEvents")
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
	shiftedTime := shiftDeliveryTimeForAllowedWindows(&n, &testTime)
	if shiftedTime.Hour() != 10 || shiftedTime.Minute() != 0 || shiftedTime.Second() != 0 {
		t.Fatalf("Expected shifted time to be 10am, got %v", shiftedTime)
	}
	if shiftedTime.Day() != testTime.Day() || shiftedTime.Month() != testTime.Month() || shiftedTime.Year() != testTime.Year() {
		t.Fatalf("Expected shifted time to be same day, got %v", shiftedTime)
	}

	// make delivery window too soon. Should shift to next day start of window
	n.DeliveryWindowTODStartMinutes = 7 * 60
	n.DeliveryWindowTODEndMinutes = 8 * 60
	shiftedTime = shiftDeliveryTimeForAllowedWindows(&n, &testTime)
	if shiftedTime.Hour() != 7 || shiftedTime.Minute() != 0 || shiftedTime.Second() != 0 {
		t.Fatalf("Expected shifted time to be 7am, got %v", shiftedTime)
	}
	if shiftedTime.Day() != testTime.Day()+1 || shiftedTime.Month() != testTime.Month() || shiftedTime.Year() != testTime.Year() {
		t.Fatalf("Expected shifted time to be next day, got %v", shiftedTime)
	}

	// make delivery window okay, should not modify
	n.DeliveryWindowTODStartMinutes = 7 * 60
	n.DeliveryWindowTODEndMinutes = 11 * 60
	shiftedTime = shiftDeliveryTimeForAllowedWindows(&n, &testTime)
	if shiftedTime.Sub(testTime) != 0 {
		t.Fatalf("Expected shifted time to be same as original, got %v", shiftedTime)
	}

	// only allow weekends, should not modify as date is on weekend
	n.DeliveryDaysOfWeek = []time.Weekday{
		time.Sunday,
		time.Saturday,
	}
	shiftedTime = shiftDeliveryTimeForAllowedWindows(&n, &testTime)
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
	shiftedTime = shiftDeliveryTimeForAllowedWindows(&n, &testTime)
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
	shiftedTime = shiftDeliveryTimeForAllowedWindows(&n, &testTime)
	if shiftedTime.Hour() != 7 || shiftedTime.Minute() != 0 || shiftedTime.Second() != 0 {
		t.Fatalf("Expected shifted time to be 7am, got %v", shiftedTime)
	}
	if shiftedTime.Day() != testTime.Day()+1 || shiftedTime.Weekday() != time.Monday || shiftedTime.Month() != testTime.Month() || shiftedTime.Year() != testTime.Year() {
		t.Fatalf("Expected shifted time of day to be next day, got %v", shiftedTime)
	}

	// Only allow on Wednesdays, should shift from Sunday to Wednesday
	n.DeliveryDaysOfWeek = []time.Weekday{time.Wednesday}
	shiftedTime = shiftDeliveryTimeForAllowedWindows(&n, &testTime)
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
	shiftedTime = shiftDeliveryTimeForAllowedWindows(&n, &dstTime)
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
	shiftedTime = shiftDeliveryTimeForAllowedWindows(&n, &testTime)
	if shiftedTime.Sub(testTime) != 0 {
		t.Fatalf("Expected shifted time to be same as original, got %v", shiftedTime)
	}
	chicagoTimeZone, err := time.LoadLocation("America/Chicago")
	if err != nil {
		t.Fatal(err)
	}
	chicagoTime := testTime.In(chicagoTimeZone)
	shiftedTime = shiftDeliveryTimeForAllowedWindows(&n, &chicagoTime)
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
