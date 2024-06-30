package datamodel

import (
	"encoding/json"
	"os"
	"slices"
	"testing"
	"time"
)

func TestNotificationActionValidators(t *testing.T) {
	var timestamp int64 = 1000
	a := Notification{
		Title:              "Title",
		Body:               "Body",
		ActionName:         "ActionName",
		ID:                 "io.criticalmoments.test",
		DeliveryDaysOfWeek: []time.Weekday{time.Monday, time.Tuesday},
		DeliveryTime: DeliveryTime{
			TimestampEpoch: &timestamp,
		},
	}
	if !a.Validate() {
		t.Fatal(a.ValidateReturningUserReadableIssue())
	}
	a.Body = ""
	if !a.Validate() {
		t.Fatal("Didn't allow empty body")
	}
	a.ActionName = ""
	if !a.Validate() {
		t.Fatal("Didn't allow empty action name")
	}
	a.Title = ""
	if a.Validate() {
		t.Fatal("Allowed empty title")
	}
	a.Title = "title"
	if !a.Validate() {
		t.Fatal("should be valid")
	}
	a.ID = ""
	if a.Validate() {
		t.Fatal("Allowed empty ID")
	}
}

func TestJsonParsingMinimalFieldsNotif(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/notifications/valid/minimalValidNotif.json")
	if err != nil {
		t.Fatal()
	}
	var n Notification
	err = json.Unmarshal(testFileData, &n)
	if err != nil {
		t.Fatal(err)
	}

	if n.Title != "title" {
		t.Fatal("failed to parse title")
	}
	if n.Body != "" {
		t.Fatal("failed to parse body as nil")
	}
	if n.ActionName != "" {
		t.Fatal("failed to parse actionName as nil")
	}
	if n.ID != "" {
		t.Fatal("ID should be nil")
	}
	if !slices.Equal(n.DeliveryDaysOfWeek, allDaysOfWeek) {
		t.Fatal("failed to parse default delivery days of week")
	}
	if n.DeliveryWindowLocalTimeOfDayStart != defaultDeliveryWindowLocalTimeStart {
		t.Fatal("failed to parse default delivery start time")
	}
	if n.DeliveryWindowLocalTimeOfDayEnd != defaultDeliveryWindowLocalTimeEnd {
		t.Fatal("failed to parse default delivery end time")
	}
	if n.IdealDevlieryConditions != nil {
		t.Fatal("failed to parse ideal delivery conditions")
	}
	if n.CancelationEvents != nil {
		t.Fatal("failed to parse cancelation condition")
	}
	if n.DeliveryTime.TimestampEpoch == nil || *n.DeliveryTime.TimestampEpoch != 1000 {
		t.Fatal("failed to parse delivery time")
	}
	if n.DeliveryTime.EventInstanceString != nil || n.DeliveryTime.EventInstance() != EventInstanceTypeLatest {
		t.Fatal("event instance should default to latest when not set")
	}
}

func TestJsonParsingMaxFieldsNotif(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/notifications/valid/maximalValidNotif.json")
	if err != nil {
		t.Fatal()
	}
	var n Notification
	err = json.Unmarshal(testFileData, &n)
	if err != nil {
		t.Fatal(err)
	}

	if n.Title != "title" {
		t.Fatal("failed to parse title")
	}
	if n.Body != "body" {
		t.Fatal("failed to parse body")
	}
	if n.ActionName != "actionName" {
		t.Fatal("failed to parse actionName")
	}
	if !slices.Equal(n.DeliveryDaysOfWeek, []time.Weekday{time.Monday, time.Tuesday}) {
		t.Fatal("failed to parse delivery days of week")
	}
	if n.DeliveryWindowLocalTimeOfDayStart != 60 {
		t.Fatal("failed to parse delivery start time")
	}
	if n.DeliveryWindowLocalTimeOfDayEnd != 120 {
		t.Fatal("failed to parse delivery end time")
	}
	if n.IdealDevlieryConditions == nil {
		t.Fatal("failed to parse ideal delivery conditions")
	}
	if n.IdealDevlieryConditions.Condition.conditionString != "true" {
		t.Fatal("failed to parse ideal delivery condition")
	}
	if n.IdealDevlieryConditions.MaxWaitTime != 10 {
		t.Fatal("failed to parse ideal delivery max wait time")
	}
	if n.CancelationEvents == nil {
		t.Fatal("failed to parse cancelation condition")
	}
	if !slices.Equal(*n.CancelationEvents, []string{"event1", "event2"}) {
		t.Fatal("failed to parse cancelation events")
	}
	if n.DeliveryTime.EventName == nil || *n.DeliveryTime.EventName != "some_event" || n.DeliveryTime.EventOffset == nil || *n.DeliveryTime.EventOffset != 300 {
		t.Fatal("failed to parse delivery time")
	}
	if n.DeliveryTime.EventInstanceString == nil || n.DeliveryTime.EventInstance() != EventInstanceTypeFirst {
		t.Fatal("event instance should set to first")
	}
}

func TestJsonParsingInvalidNotif(t *testing.T) {
	cases := []string{
		"./test/testdata/actions/notifications/invalid/invalid.json",
		"./test/testdata/actions/notifications/invalid/invalidCondition.json",        // add_test_case
		"./test/testdata/actions/notifications/invalid/invalidCancelationEvent.json", // add_test_case
	}

	for _, testFile := range cases {
		testFileData, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatal()
		}
		var n Notification
		err = json.Unmarshal(testFileData, &n)
		if err == nil {
			t.Fatalf("Allowed invalid json: %v", testFile)
		}
	}
}

func TestParseDowString(t *testing.T) {
	cases := map[string][]time.Weekday{
		"Monday,Tuesday": {time.Monday, time.Tuesday},
		// dupes
		"Monday,Monday,Tuesday": {time.Monday, time.Tuesday}, // add_test_case
		// Invalid data
		"Monday,Febuary,Tuesday": {time.Monday, time.Tuesday}, // add_test_case
		// robustness
		"Monday,,Tuesday": {time.Monday, time.Tuesday}, // add_test_case
		// fix order
		"Monday,Tuesday,Sunday": {time.Sunday, time.Monday, time.Tuesday}, // add_test_case
		// empty
		"": {}, // add_test_case
		// all
		"Sunday,Monday,Tuesday,Wednesday,Thursday,Friday,Saturday": {time.Sunday, time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday}, // add_test_case
	}

	for dowString, target := range cases {
		dowResult := parseDaysOfWeekString(dowString)
		if !slices.Equal(dowResult, target) {
			t.Fatalf("Day of week string failed parsing for: %v", dowResult)
		}
	}
}

func TestDeliveryTimeValidation(t *testing.T) {
	// Case: Both TimestampEpoch and EventName are nil
	dt := DeliveryTime{}
	issue := dt.ValidateReturningUserReadableIssue()
	if issue != "DeliveryTime must have either a Timestamp or an EventName defined." {
		t.Fatalf("Unexpected validation issue: %v", issue)
	}

	// Case: Both TimestampEpoch and EventName are defined
	timestamp := int64(1000)
	eventName := "event"
	dt = DeliveryTime{TimestampEpoch: &timestamp, EventName: &eventName}
	issue = dt.ValidateReturningUserReadableIssue()
	if issue != "DeliveryTime cannot have both a Timestamp and an EventName defined." {
		t.Fatalf("Unexpected validation issue: %v", issue)
	}

	// Case: Both TimestampEpoch and EventOffset are defined
	eventOffset := 30
	dt = DeliveryTime{TimestampEpoch: &timestamp, EventOffset: &eventOffset}
	issue = dt.ValidateReturningUserReadableIssue()
	if issue != "DeliveryTime cannot have both a Timestamp and an EventOffset defined." {
		t.Fatalf("Unexpected validation issue: %v", issue)
	}

	// Case: Valid TimestampEpoch
	dt = DeliveryTime{TimestampEpoch: &timestamp}
	issue = dt.ValidateReturningUserReadableIssue()
	if issue != "" {
		t.Fatalf("Unexpected validation issue: %v", issue)
	}

	// Case: Valid EventName
	dt = DeliveryTime{EventName: &eventName}
	issue = dt.ValidateReturningUserReadableIssue()
	if issue != "" {
		t.Fatalf("Unexpected validation issue: %v", issue)
	}

	// Case: Valid EventOffset with EventName
	dt = DeliveryTime{EventName: &eventName, EventOffset: &eventOffset}
	issue = dt.ValidateReturningUserReadableIssue()
	if issue != "" {
		t.Fatalf("Unexpected validation issue: %v", issue)
	}

	// Empty/missing should detfault to latest
	if dt.EventInstance() != EventInstanceTypeLatest {
		t.Fatal("failed to return latest")
	}
	s := ""
	dt.EventInstanceString = &s
	if dt.EventInstance() != EventInstanceTypeLatest {
		t.Fatal("failed to return latest")
	}
	s = "latest"
	dt.EventInstanceString = &s
	if dt.EventInstance() != EventInstanceTypeLatest {
		t.Fatal("failed to return latest")
	}
	s = "first"
	dt.EventInstanceString = &s
	if dt.EventInstance() != EventInstanceTypeFirst {
		t.Fatal("failed to return First")
	}
	if dt.ValidateReturningUserReadableIssue() != "" {
		t.Fatal("Errored on valid event instance string")
	}
	s = "invalid"
	dt.EventInstanceString = &s
	if dt.EventInstance() != EventInstanceTypeUnknown {
		t.Fatal("failed to return latest for unrecognized")
	}
	if dt.ValidateReturningUserReadableIssue() != "" {
		t.Fatal("Errored on invalid event instance string in not strict mode")
	}
	StrictDatamodelParsing = true
	defer func() {
		StrictDatamodelParsing = false
	}()
	if dt.ValidateReturningUserReadableIssue() == "" {
		t.Fatal("Did not error on invalid event instance string in not strict mode")
	}
	s = "latest"
	if dt.ValidateReturningUserReadableIssue() != "" {
		t.Fatal("Errored on valid event instance string in strict mode")
	}
}
