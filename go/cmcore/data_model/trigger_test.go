package datamodel

import (
	"encoding/json"
	"os"
	"testing"
)

func TestTriggerJsonValidation(t *testing.T) {
	trigger := Trigger{}
	if trigger.Validate() {
		t.Fatal()
	}
	trigger.EventName = "my_event"
	if trigger.Validate() {
		t.Fatal()
	}
	trigger.ActionName = "my_action"
	if !trigger.Validate() {
		t.Fatal()
	}
	trigger.EventName = ""
	if trigger.Validate() {
		t.Fatal()
	}
}

func TestTriggerParsingValidTrigger(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/triggers/valid/validTrigger.json")
	var trigger Trigger
	err = json.Unmarshal(testFileData, &trigger)
	if err != nil {
		t.Fatal()
	}

	// Check defaults for values not included in json
	if trigger.ActionName != "my_action" {
		t.Fatal()
	}
	if trigger.EventName != "my_event" {
		t.Fatal()
	}
	if !trigger.Validate() {
		t.Fatal()
	}
}

func TestTriggerParsingInvalidTrigger(t *testing.T) {
	testFileData, _ := os.ReadFile("./test/testdata/triggers/invalid/empty.json")
	var trigger Trigger
	json.Unmarshal(testFileData, &trigger)

	// Check defaults for values not included in json
	if trigger.ActionName != "" {
		t.Fatal()
	}
	if trigger.EventName != "" {
		t.Fatal()
	}
	if trigger.Validate() {
		t.Fatal()
	}
}
