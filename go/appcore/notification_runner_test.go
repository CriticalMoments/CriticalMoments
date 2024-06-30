package appcore

import (
	"math"
	"testing"
	"time"
)

func TestNotificationList(t *testing.T) {
	ac, err := testBuildValidTestAppCore(t)
	if err != nil {
		t.Fatal(err)
	}
	err = ac.Start(true)
	if err != nil {
		t.Fatal(err)
	}

	plan, err := ac.NotificationPlan()
	if err != nil {
		t.Fatal(err)
	}
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
	err = ac.Start(true)
	if err != nil {
		t.Fatal(err)
	}

	// No events, should be nothing scheduled
	plan, err := ac.NotificationPlan()
	if err != nil {
		t.Fatal(err)
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
	plan, err = ac.NotificationPlan()
	if err != nil {
		t.Fatal(err)
	}
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
	// Expect it scheduled within 1s of now
	if math.Abs(float64(sn.ScheduledAtEpochMilliseconds()-time.Now().UnixMilli())) > 1000 {
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
	plan, err = ac.NotificationPlan()
	if err != nil {
		t.Fatal(err)
	}
	if plan.ScheduledNotificationCount() != 2 {
		t.Fatalf("Expected 2 scheduled notifications, got %d", plan.ScheduledNotificationCount())
	}
	sn = plan.ScheduledNotificationAtIndex(0)
	if sn.Notification.ID != "event2Notification" {
		// Index is interterminate
		sn = plan.ScheduledNotificationAtIndex(1)
	}
	if sn.Notification.ID != "event2Notification" {
		t.Fatalf("Expected ScheduledNotificationAtIndex to return event notification, got %s", sn.Notification.ID)
	}
	// Expect it scheduled within 1s of 60s from now
	if math.Abs(float64(sn.ScheduledAtEpochMilliseconds()-time.Now().UnixMilli()-60000)) > 1000 {
		t.Fatalf("Expected ScheduledAtEpoch to return now, got %d", sn.ScheduledAtEpochMilliseconds())
	}
}

// TODO_P0: test event processor: firing event should dispatch "schedule notification"
