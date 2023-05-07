package datamodel

import (
	"encoding/json"
	"os"
	"testing"
)

func TestTriggerJsonValidation(t *testing.T) {
	tj := jsonTrigger{}
	if tj.Validate() {
		t.Fatal()
	}
	tj.EventName = "my_event"
	if tj.Validate() {
		t.Fatal()
	}
	tj.ActionName = "my_action"
	if !tj.Validate() {
		t.Fatal()
	}
	tj.EventName = ""
	if tj.Validate() {
		t.Fatal()
	}
}
func TestJsonParsingTrigger(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/triggers/valid/validTrigger.json")
	var jt jsonTrigger
	err = json.Unmarshal(testFileData, &jt)
	if err != nil {
		t.Fatal()
	}

	// Check defaults for values not included in json
	if jt.ActionName != "my_action" {
		t.Fatal()
	}
	if jt.EventName != "my_event" {
		t.Fatal()
	}
	if !jt.Validate() {
		t.Fatal()
	}
}

func TestJsonParsingInvalidTrigger(t *testing.T) {
	testFileData, _ := os.ReadFile("./test/testdata/triggers/invalid/empty.json")
	var jt jsonTrigger
	json.Unmarshal(testFileData, &jt)

	// Check defaults for values not included in json
	if jt.ActionName != "" {
		t.Fatal()
	}
	if jt.EventName != "" {
		t.Fatal()
	}
	if jt.Validate() {
		t.Fatal()
	}
}
