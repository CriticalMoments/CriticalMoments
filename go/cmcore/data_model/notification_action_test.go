package datamodel

import (
	"encoding/json"
	"os"
	"testing"
)

func TestNotificationActionValidators(t *testing.T) {
	// valid
	a := NotificationAction{
		Title:      "Title",
		Body:       "Body",
		ActionName: "ActionName",
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
}

func TestJsonParsingMinimalFieldsNotif(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/notifications/valid/minimalValidAlert.json")
	if err != nil {
		t.Fatal()
	}
	var ac ActionContainer
	err = json.Unmarshal(testFileData, &ac)
	if err != nil {
		t.Fatal(err)
	}

	if ac.ActionType != ActionTypeEnumNotification {
		t.Fatal("wrong action type")
	}
	na := ac.NotificationAction

	if na.Title != "title" {
		t.Fatal("failed to parse title")
	}
	if na.Body != "" {
		t.Fatal("failed to parse body as nil")
	}
	if na.ActionName != "" {
		t.Fatal("failed to parse actionName as nil")
	}
}

func TestJsonParsingMaxFieldsNotif(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/notifications/valid/maximalValidAlert.json")
	if err != nil {
		t.Fatal()
	}
	var ac ActionContainer
	err = json.Unmarshal(testFileData, &ac)
	if err != nil {
		t.Fatal(err)
	}

	if ac.ActionType != ActionTypeEnumNotification {
		t.Fatal("wrong action type")
	}
	na := ac.NotificationAction

	if na.Title != "title" {
		t.Fatal("failed to parse title")
	}
	if na.Body != "body" {
		t.Fatal("failed to parse body")
	}
	if na.ActionName != "actionName" {
		t.Fatal("failed to parse actionName")
	}
}

func TestJsonParsingInvalidNotif(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/notifications/invalid/invalid.json")
	if err != nil {
		t.Fatal()
	}
	var ac ActionContainer
	err = json.Unmarshal(testFileData, &ac)
	if err == nil {
		t.Fatal("Allowed invalid json")
	}
}
