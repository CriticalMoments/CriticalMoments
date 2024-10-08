package datamodel

import (
	"encoding/json"
	"os"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestNotificationActionValidators(t *testing.T) {
	var timestamp int64 = 1000
	a := Notification{
		Title:              "Title",
		Body:               "Body",
		BadgeCount:         -1,
		ActionName:         "ActionName",
		ID:                 "io.criticalmoments.test",
		DeliveryDaysOfWeek: []time.Weekday{time.Monday, time.Tuesday},
		DeliveryTime: DeliveryTime{
			TimestampEpoch: &timestamp,
		},
	}
	if !a.Valid() {
		t.Fatal(a.Check())
	}
	a.Body = ""
	if !a.Valid() {
		t.Fatal("Didn't allow empty body")
	}
	a.ActionName = ""
	if !a.Valid() {
		t.Fatal("Didn't allow empty action name")
	}
	a.Title = ""
	if a.Valid() {
		t.Fatal("Allowed empty title and body")
	}
	a.BadgeCount = 1
	if !a.Valid() {
		t.Fatal("Should be valid with badge count")
	}
	a.BadgeCount = -1
	if a.Valid() {
		t.Fatal("Needs a title, body, or badge count")
	}
	a.Title = "title"
	if !a.Valid() {
		t.Fatal("should be valid")
	}
	a.Title = ""
	a.Body = "body"
	if !a.Valid() {
		t.Fatal("should be valid")
	}
	a.ID = ""
	if a.Valid() {
		t.Fatal("Allowed empty ID")
	}
	a.ID = "io.criticalmoments.test"
	if !a.Valid() {
		t.Fatal("should be valid")
	}
	rc := 0.5
	a.RelevanceScore = &rc
	if !a.Valid() {
		t.Fatal("should be valid")
	}
	rc = 1.1
	a.RelevanceScore = &rc
	if a.Valid() {
		t.Fatal("Allowed invalid relevance score")
	}
	rc = -0.0001
	a.RelevanceScore = &rc
	if a.Valid() {
		t.Fatal("Allowed invalid relevance score")
	}
	a.RelevanceScore = nil
	a.InterruptionLevel = "passive"
	if !a.Valid() {
		t.Fatal("should be valid")
	}
	a.InterruptionLevel = "futureUnknown"
	if !a.Valid() {
		t.Fatal("should not error since not strict")
	}
	StrictDatamodelParsing = true
	defer func() {
		StrictDatamodelParsing = false
	}()
	a.InterruptionLevel = "futureUnknown"
	if a.Valid() {
		t.Fatal("should not be valid if strict")
	}
	StrictDatamodelParsing = false
	a.InterruptionLevel = ""
	a.DeliveryWindowTODEndMinutes = 24 * 60
	if a.Valid() {
		t.Fatal("delivery window out of bounds")
	}
	a.DeliveryWindowTODEndMinutes = 23*60 + 59
	if !a.Valid() {
		t.Fatal("delivery window failed in bounds")
	}
	a.DeliveryWindowTODEndMinutes = 2
	a.DeliveryWindowTODStartMinutes = 3
	if a.Valid() {
		t.Fatal("Allowed start after end")
	}
	a.DeliveryWindowTODEndMinutes = 60
	if !a.Valid() {
		t.Fatal("delivery window failed in bounds")
	}
}

func TestJsonParsingMinimalFieldsNotif(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/notifications/valid/minimalValidNotif.json")
	if err != nil {
		t.Fatal(err)
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
	if n.BadgeCount != -1 {
		t.Fatal("failed to parse badge count as unset")
	}
	if n.Sound != "" {
		t.Fatal("failed to parse sound as nil")
	}
	if n.LaunchImageName != "" {
		t.Fatal("failed to parse launch image name as nil")
	}
	if n.RelevanceScore != nil {
		t.Fatal("failed to parse relevance score as nil")
	}
	if n.InterruptionLevel != "" {
		t.Fatal("failed to parse interruption level as nil")
	}
	if n.ScheduleCondition != nil {
		t.Fatal("faild to parse nil scheduled condition")
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
	if n.DeliveryWindowTODStartMinutes != defaultDeliveryWindowLocalTimeStart {
		t.Fatal("failed to parse default delivery start time")
	}
	if n.DeliveryWindowTODEndMinutes != defaultDeliveryWindowLocalTimeEnd {
		t.Fatal("failed to parse default delivery end time")
	}
	if n.IdealDeliveryConditions != nil {
		t.Fatal("failed to parse ideal delivery conditions")
	}
	if n.CancelationEvents != nil {
		t.Fatal("failed to parse cancelation condition")
	}
	if n.DeliveryTime.TimestampEpoch == nil || *n.DeliveryTime.TimestampEpoch != 1000 {
		t.Fatal("failed to parse delivery time")
	}
	if n.DeliveryTime.EventInstanceString != nil || n.DeliveryTime.EventInstance() != EventInstanceTypeLatestOnce {
		t.Fatal("event instance should default to latest once when not set")
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
	if n.BadgeCount != 10 {
		t.Fatal("failed to parse badge count")
	}
	if n.Sound != "default" {
		t.Fatal("failed to parse sound")
	}
	if n.LaunchImageName != "storyboard1" {
		t.Fatal("failed to parse launch image name")
	}
	if *n.RelevanceScore != 0.5 {
		t.Fatal("failed to parse relevance score")
	}
	if n.InterruptionLevel != "passive" {
		t.Fatal("failed to parse interruption level")
	}
	if n.ScheduleCondition.conditionString != "true" {
		t.Fatal("failed to parse schedule condition")
	}
	if n.ActionName != "actionName" {
		t.Fatal("failed to parse actionName")
	}
	if !slices.Equal(n.DeliveryDaysOfWeek, []time.Weekday{time.Monday, time.Tuesday}) {
		t.Fatal("failed to parse delivery days of week")
	}
	if n.DeliveryWindowTODStartMinutes != 90 {
		t.Fatal("failed to parse delivery start time")
	}
	if n.DeliveryWindowTODEndMinutes != 23*60+59 {
		t.Fatal("failed to parse delivery end time")
	}
	if n.IdealDeliveryConditions == nil {
		t.Fatal("failed to parse ideal delivery conditions")
	}
	if n.IdealDeliveryConditions.Condition.conditionString != "true" {
		t.Fatal("failed to parse ideal delivery condition")
	}
	if n.IdealDeliveryConditions.MaxWaitTime() != 10*time.Second {
		t.Fatal("failed to parse ideal delivery max wait time")
	}
	if n.CancelationEvents == nil {
		t.Fatal("failed to parse cancelation condition")
	}
	if !slices.Equal(*n.CancelationEvents, []string{"event1", "event2"}) {
		t.Fatal("failed to parse cancelation events")
	}
	if n.DeliveryTime.EventName == nil || *n.DeliveryTime.EventName != "some_event" || n.DeliveryTime.EventOffsetSeconds == nil || *n.DeliveryTime.EventOffsetSeconds != 300 {
		t.Fatal("failed to parse delivery time")
	}
	if n.DeliveryTime.EventInstanceString == nil || n.DeliveryTime.EventInstance() != EventInstanceTypeFirst {
		t.Fatal("event instance should set to first")
	}
}

func TestJsonParsingInvalidNotif(t *testing.T) {
	cases := map[string]string{
		"./test/testdata/actions/notifications/invalid/invalidCondition.json":        "Error parsing condition string: true > 1",
		"./test/testdata/actions/notifications/invalid/invalidCancelationEvent.json": "Notification has an empty string in cancelationEvents",      // add_test_case
		"./test/testdata/actions/notifications/invalid/invalidBadgeCount.json":       "Notification badgeCount must be greater than or equal to 0", // add_test_case
	}

	StrictDatamodelParsing = true
	defer func() {
		StrictDatamodelParsing = false
	}()

	for testFile, expectedError := range cases {
		testFileData, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatal()
		}
		var n Notification
		err = json.Unmarshal(testFileData, &n)
		if err == nil {
			t.Fatalf("Allowed invalid json: %v", testFile)
		}
		errorString := err.Error()
		if !strings.Contains(errorString, expectedError) {
			t.Fatalf("Allowed invalid: %v, expected error: %v, got error: %v", testFile, expectedError, errorString)
		}
	}
}

func TestParseDowString(t *testing.T) {
	cases := map[string][]time.Weekday{
		"Monday,Tuesday": {time.Monday, time.Tuesday},
		// dupes
		"Monday,Monday,Tuesday": {time.Monday, time.Tuesday}, // add_test_case
		// Invalid data
		"Monday,Pineapple,Tuesday": {time.Monday, time.Tuesday}, // add_test_case
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
	issue := dt.Check().Error()
	if issue != "Notification DeliveryTime must have either a Timestamp or an EventName defined." {
		t.Fatalf("Unexpected validation issue: %v", issue)
	}

	// Case: Both TimestampEpoch and EventName are defined
	timestamp := int64(1000)
	eventName := "event"
	dt = DeliveryTime{TimestampEpoch: &timestamp, EventName: &eventName}
	issue = dt.Check().Error()
	if issue != "Notification DeliveryTime cannot have both a Timestamp and an EventName defined." {
		t.Fatalf("Unexpected validation issue: %v", issue)
	}

	// Case: Both TimestampEpoch and EventOffset are defined
	eventOffsetSeconds := 30
	dt = DeliveryTime{TimestampEpoch: &timestamp, EventOffsetSeconds: &eventOffsetSeconds}
	issue = dt.Check().Error()
	if issue != "Notification DeliveryTime cannot have both a Timestamp and an EventOffset defined." {
		t.Fatalf("Unexpected validation issue: %v", issue)
	}

	// Case: Valid TimestampEpoch
	dt = DeliveryTime{TimestampEpoch: &timestamp}
	if dt.Check() != nil {
		t.Fatalf("Unexpected validation issue: %v", issue)
	}

	// Case: Valid EventName
	dt = DeliveryTime{EventName: &eventName}
	if dt.Check() != nil {
		t.Fatalf("Unexpected validation issue: %v", issue)
	}

	// Case: Valid EventOffset with EventName
	dt = DeliveryTime{EventName: &eventName, EventOffsetSeconds: &eventOffsetSeconds}
	if dt.Check() != nil {
		t.Fatalf("Unexpected validation issue: %v", issue)
	}

	// Empty/missing should detfault to latest
	if dt.EventInstance() != EventInstanceTypeLatestOnce {
		t.Fatal("failed to return latest for default")
	}
	s := ""
	dt.EventInstanceString = &s
	if dt.EventInstance() != EventInstanceTypeLatestOnce {
		t.Fatal("failed to return latest once for empty string")
	}
	s = "latest-once"
	dt.EventInstanceString = &s
	if dt.EventInstance() != EventInstanceTypeLatestOnce {
		t.Fatal("failed to return latest once")
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
	if dt.Check() != nil {
		t.Fatal("Errored on valid event instance string")
	}
	s = "invalid"
	dt.EventInstanceString = &s
	if dt.EventInstance() != EventInstanceTypeUnknown {
		t.Fatal("failed to return latest for unrecognized")
	}
	if dt.Check() != nil {
		t.Fatal("Errored on invalid event instance string in not strict mode")
	}
	StrictDatamodelParsing = true
	defer func() {
		StrictDatamodelParsing = false
	}()
	if dt.Check() == nil {
		t.Fatal("Did not error on invalid event instance string in not strict mode")
	}
	s = "latest"
	if dt.Check() != nil {
		t.Fatal("Errored on valid event instance string in strict mode")
	}
}

func TestNotificationUniqueID(t *testing.T) {
	n := Notification{
		ID: "test",
	}
	if n.UniqueID() != "io.criticalmoments.notifications.test" {
		t.Fatal("failed to generate unique ID")
	}
}

func TestHHMMStringParsing(t *testing.T) {
	cases := map[string]int{
		"00:00": 0,
		"01:00": 60,         // add_test_case
		"13:13": 13*60 + 13, // add_test_case
		"23:59": 23*60 + 59, // add_test_case
	}

	for hhmmString, seconds := range cases {
		secondsResult, err := parseMinutesFromHHMMString(hhmmString)
		if err != nil {
			t.Fatalf("Failed to parse %v", hhmmString)
		}
		if secondsResult != seconds {
			t.Fatalf("Failed to parse %v, got %v", hhmmString, secondsResult)
		}
	}

	errorCases := []string{
		"00:60",
		"24:00",    // add_test_case
		"24:01",    // add_test_case
		"25:00",    // add_test_case
		"00:00:00", // add_test_case
		"00:00:60", // add_test_case
		"",         // add_test_case
		"asd:sdf",  // add_test_case
		"asdf",     // add_test_case
	}

	for _, hhmmString := range errorCases {
		_, err := parseMinutesFromHHMMString(hhmmString)
		if err == nil {
			t.Fatalf("Failed to error on %v", hhmmString)
		}
	}
}

func TestAccessors(t *testing.T) {
	i := IdealDeliveryConditions{
		MaxWaitTimeSeconds: NotificationMaxIdealWaitTimeForever,
	}

	if !i.WaitForever() {
		t.Fatal("Wait forever incorrect")
	}
	if i.MaxWaitTime() < time.Hour*24*365*100 {
		t.Fatal("Fallback to huge time failed")
	}

	i.MaxWaitTimeSeconds = 10
	if i.WaitForever() {
		t.Fatal("waitforever incorrect")
	}
	if i.MaxWaitTime() != time.Second*10 {
		t.Fatal("conversion to duration failed")
	}
}
