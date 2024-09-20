package datamodel

import (
	"encoding/json"
	"os"
	"testing"
)

func TestTriggerJsonValidation(t *testing.T) {
	trigger := Trigger{}
	if trigger.Valid() {
		t.Fatal()
	}
	trigger.EventName = "my_event"
	if trigger.Valid() {
		t.Fatal()
	}
	trigger.ActionName = "my_action"
	if !trigger.Valid() {
		t.Fatal()
	}
	trigger.EventName = ""
	if trigger.Valid() {
		t.Fatal()
	}
}

func TestTriggerParsingValidTrigger(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/triggers/valid/validTrigger.json")
	if err != nil {
		t.Fatal()
	}
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
	if trigger.Condition.conditionString != "3 > 2" {
		t.Fatal()
	}
	if !trigger.Valid() {
		t.Fatal()
	}
}

func TestTriggerParsingInvalidTrigger(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/triggers/invalid/empty.json")
	if err != nil {
		t.Fatal()
	}
	var trigger Trigger
	err = json.Unmarshal(testFileData, &trigger)
	if err == nil {
		t.Fatal("allowed invalid empty trigger")
	}
	if trigger.ActionName != "" || trigger.EventName != "" {
		t.Fatal("trigger parse issue")
	}
	if trigger.Valid() {
		t.Fatal("validated empty trigger")
	}
}
func TestTriggerParsingInvalidConditionTrigger(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/triggers/invalid/invalidTriggerCondition.json")
	if err != nil {
		t.Fatal()
	}
	var trigger Trigger
	err = json.Unmarshal(testFileData, &trigger)
	if err == nil {
		t.Fatal("allowed invalid condition in trigger")
	}
	if trigger.ActionName != "my_action" || trigger.EventName != "my_event" {
		t.Fatal("trigger parse issue")
	}
	if trigger.Valid() {
		t.Fatal("validated trigger with invalid condition")
	}
}
