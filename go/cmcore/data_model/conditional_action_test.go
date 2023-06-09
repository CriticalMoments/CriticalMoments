package datamodel

import (
	"encoding/json"
	"os"
	"testing"
)

func TestConditionalActionValidators(t *testing.T) {
	c := ConditionalAction{}
	if c.Validate() {
		t.Fatal("Conditional actions require a condition")
	}
	c.Condition = "(network_connection_type == 'wifi')"
	if c.Validate() {
		t.Fatal("Conditional actions require a passed action")
	}
	c.PassedActionName = "pass_action"
	if !c.Validate() {
		t.Fatal("Conditional action should be valid")
	}
	an, err := c.AllEmbeddedActionNames()
	if err != nil && len(an) != 1 && an[0] != "pass_action" {
		t.Fatal("Failed to return action name for pass action")
	}
	c.Condition = "not_a_valid_var > 5"
	if c.Validate() {
		t.Fatal("Conditional action should validate expression validity")
	}
	c.Condition = ""
	if c.Validate() {
		t.Fatal("Conditional action require condition")
	}
	c.Condition = "true"
	if !c.Validate() {
		t.Fatal("Conditional action should be valid")
	}
	c.FailedActionName = "fail_action"
	if !c.Validate() {
		t.Fatal("Conditional action should be valid with or without failed_action")
	}
	an, err = c.AllEmbeddedActionNames()
	if err != nil && len(an) != 2 && an[0] != "pass_action" && an[1] != "fail_action" {
		t.Fatal("Failed to return action name for pass action")
	}
}

func TestJsonParsingValidConditional(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/conditional_actions/valid.json")
	if err != nil {
		t.Fatal()
	}
	var ac ActionContainer
	err = json.Unmarshal(testFileData, &ac)
	if err != nil {
		t.Fatal(err)
	}
	if ac.ActionType != ActionTypeEnumConditional {
		t.Fatal("Failed to parse valid conditional action")
	}
	if ac.ConditionalAction.Condition != "(device_battery_state == 'charging' || device_battery_state == 'full')" {
		t.Fatal("Failed to parse condition")
	}
	if ac.ConditionalAction.PassedActionName != "conditional_true" || ac.ConditionalAction.FailedActionName != "conditional_false" {
		t.Fatal("Failed to parse action names")
	}
}

func TestJsonParsingInvalidCondiationalAction(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/conditional_actions/invalid.json")
	if err != nil {
		t.Fatal()
	}
	var ac ActionContainer
	err = json.Unmarshal(testFileData, &ac)
	if err == nil || ac.ActionType == ActionTypeEnumConditional {
		t.Fatal("Invalid conditionals should not parse")
	}
}
