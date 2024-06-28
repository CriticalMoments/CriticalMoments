package datamodel

import (
	"encoding/json"
	"os"
	"slices"
	"testing"
	"time"
)

func TestNotificationActionValidators(t *testing.T) {
	// valid
	a := Notification{
		Title:              "Title",
		Body:               "Body",
		ActionName:         "ActionName",
		ID:                 "io.criticalmoments.test",
		DeliveryDaysOfWeek: []time.Weekday{time.Monday, time.Tuesday},
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
}

func TestJsonParsingInvalidNotif(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/notifications/invalid/invalid.json")
	if err != nil {
		t.Fatal()
	}
	var n Notification
	err = json.Unmarshal(testFileData, &n)
	if err == nil {
		t.Fatal("Allowed invalid json")
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
