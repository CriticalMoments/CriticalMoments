package datamodel

import (
	"encoding/json"
	"os"
	"testing"
)

func testHelperNewCondition(s string, t *testing.T) *Condition {
	c, err := NewCondition(s)
	if err != nil {
		t.Fatal("Condition in test invalid", err)
	}
	return c
}

func TestConditionalActionValidators(t *testing.T) {
	c := ConditionalAction{}
	if c.Validate() {
		t.Fatal("Conditional actions require a condition")
	}
	c.Condition = testHelperNewCondition("(network_connection_type == 'wifi')", t)
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
	c.Condition.conditionString = "not_a_valid_var > 5"
	if !c.Validate() {
		t.Fatal("Conditional action should validate condition validity but non-strict okay")
	}
	c.Condition = nil
	if c.Validate() {
		t.Fatal("Conditional action require condition")
	}
	c.Condition = testHelperNewCondition("true", t)
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

	StrictDatamodelParsing = true
	defer func() {
		StrictDatamodelParsing = false
	}()
	// Check it calls nested validators. Can't construct a problematic condition without reflection
	c.Condition.conditionString = "not_a_valid_var > 5"
	if c.Validate() {
		t.Fatal("Conditional action should validate condition validity")
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
	if ac.ConditionalAction.Condition.String() != "(device_battery_state == 'charging' || device_battery_state == 'full')" {
		t.Fatal("Failed to parse condition")
	}
	if ac.ConditionalAction.PassedActionName != "conditional_true" || ac.ConditionalAction.FailedActionName != "conditional_false" {
		t.Fatal("Failed to parse action names")
	}
}

func TestJsonParsingInvalidConditionalAction(t *testing.T) {
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

func TestJsonParsingInvalidConditionalActionCondition(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/conditional_actions/invalid_condition.json")
	if err != nil {
		t.Fatal()
	}
	var ac ActionContainer
	err = json.Unmarshal(testFileData, &ac)
	if err == nil || ac.ActionType == ActionTypeEnumConditional {
		t.Fatal("Invalid conditionals should not parse")
	}
}

func TestJsonParsingInvalidStrictConditionalActionCondition(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/conditional_actions/invalid_strict_condition.json")
	if err != nil {
		t.Fatal()
	}
	var ac ActionContainer
	err = json.Unmarshal(testFileData, &ac)
	if err != nil || ac.ConditionalAction.Condition.String() != "(fake_var == 'charging' || device_battery_state == 'full')" {
		t.Fatal("should pass non strict validation")
	}
	StrictDatamodelParsing = true
	defer func() {
		StrictDatamodelParsing = false
	}()
	err = json.Unmarshal(testFileData, &ac)
	if err == nil {
		t.Fatal("should not pass strict validation")
	}
}
