package datamodel

import (
	"encoding/json"
	"os"
	"testing"
)

func TestNotificationActionValidators(t *testing.T) {
	// valid
	a := Notification{
		Title:      "Title",
		Body:       "Body",
		ActionName: "ActionName",
		ID:         "io.criticalmoments.test",
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
	testFileData, err := os.ReadFile("./test/testdata/actions/notifications/valid/minimalValidAlert.json")
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
}

func TestJsonParsingMaxFieldsNotif(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/notifications/valid/maximalValidAlert.json")
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
