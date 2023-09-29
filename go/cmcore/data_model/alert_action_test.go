package datamodel

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"golang.org/x/exp/slices"
)

func TestAlertActionValidators(t *testing.T) {
	// valid
	a := AlertAction{
		Title:        "Title",
		Message:      "Message",
		ShowOkButton: true,
		Style:        AlertActionStyleEnumDialog,
	}
	if !a.Validate() {
		t.Fatal(a.ValidateReturningUserReadableIssue())
	}
	a.Style = ""
	if a.Validate() {
		t.Fatal("Allowed empty style")
	}
	a.Style = "asdf"
	if a.Validate() {
		t.Fatal("Allowed invalid style")
	}
	a.Style = AlertActionStyleEnumLarge
	if !a.Validate() {
		t.Fatal(a.ValidateReturningUserReadableIssue())
	}
	a.Title = ""
	if !a.Validate() {
		t.Fatal("Should allow empty title if message still provided")
	}
	a.Message = ""
	if a.Validate() {
		t.Fatal("Should not allow empty title and message")
	}
	a.Title = "New title"
	if !a.Validate() {
		t.Fatal("Should allow title and not message")
	}
	a.ShowOkButton = false
	a.OkButtonActionName = "action"
	if a.Validate() {
		t.Fatal("Should not allow an okay actio when ok button hidden")
	}
	a.ShowOkButton = true
	a.OkButtonActionName = ""
	if !a.Validate() {
		t.Fatal("Should allow okay without an action")
	}
	cb := AlertActionCustomButton{}
	a.CustomButtons = []*AlertActionCustomButton{&cb}
	if a.Validate() {
		t.Fatal("Should vaidate buttons as well")
	}
	a.CustomButtons = []*AlertActionCustomButton{}
	if !a.Validate() {
		t.Fatal()
	}
	a.ShowOkButton = false
	if a.Validate() {
		t.Fatal("Alert requires ok or custom buttons, but allowed neither")
	}
}

func TestCustomButtonValidation(t *testing.T) {
	b := AlertActionCustomButton{
		Label: "Label",
		Style: AlertActionButtonStyleEnumPrimary,
	}
	if !b.Validate() {
		t.Fatal("Valid button fails validation")
	}
	b.Style = ""
	if b.Validate() {
		t.Fatal("Empty button style should not validate")
	}
	b.Style = "adsf"
	if b.Validate() {
		t.Fatal("INvalid style should not validate")
	}
	b.Style = AlertActionButtonStyleEnumDestructive
	if !b.Validate() {
		t.Fatal("Valid button style fails validation")
	}
	b.Style = AlertActionButtonStyleEnumDefault
	if !b.Validate() {
		t.Fatal("Valid button style fails validation")
	}
	b.Label = ""
	if b.Validate() {
		t.Fatal("Buttons require a label")
	}
}

func TestJsonParsingMaximalFieldsAlert(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/alert/valid/maximalValidAlert.json")
	if err != nil {
		t.Fatal()
	}
	var ac ActionContainer
	err = json.Unmarshal(testFileData, &ac)
	if err != nil {
		t.Fatal(err)
	}

	if ac.ActionType != ActionTypeEnumAlert {
		t.Fatal()
	}
	if ac.Condition.String() != "platform == 'iOS'" {
		t.Fatal()
	}
	a := ac.AlertAction
	if a == nil || !a.Validate() {
		t.Fatal()
	}
	if a.Title != "For real?" {
		t.Fatal()
	}
	if a.Message != "Are you sure you want to?" {
		t.Fatal()
	}
	if !a.ShowCancelButton || !a.ShowOkButton {
		t.Fatal()
	}
	if a.OkButtonActionName != "custom_event" {
		t.Fatal()
	}
	if a.Style != AlertActionStyleEnumLarge {
		t.Fatal()
	}
	cb1 := a.CustomButtons[0]
	if cb1.Label != "Custom 1" || cb1.ActionName != "event1" || cb1.Style != AlertActionButtonStyleEnumPrimary {
		t.Fatal()
	}
	cb2 := a.CustomButtons[1]
	if cb2.Label != "Custom 2" || cb2.ActionName != "event2" || cb2.Style != AlertActionButtonStyleEnumDestructive {
		t.Fatal()
	}
	cb3 := a.CustomButtons[2]
	if cb3.Label != "Custom 3" || cb3.ActionName != "event3" || cb3.Style != AlertActionButtonStyleEnumDefault {
		t.Fatal()
	}

	// Theme names
	themes, err := a.AllEmbeddedThemeNames()
	if err != nil || len(themes) > 0 {
		t.Fatal("alerts don't have themes!")
	}
	// Embedded action names
	actions, err := a.AllEmbeddedActionNames()
	if err != nil {
		t.Fatal(err)
	}
	expectedActions := []string{"custom_event", "event1", "event2", "event3"}
	for _, expected := range expectedActions {
		if !slices.Contains(actions, expected) {
			t.Fatal(fmt.Sprintf("Expected %v but missing", expected))
		}
	}
	if len(actions) != len(expectedActions) {
		t.Fatal()
	}
}

func TestJsonParsingMinimalFieldsAlert(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/alert/valid/minimalValidAlert.json")
	if err != nil {
		t.Fatal()
	}
	var ac ActionContainer
	err = json.Unmarshal(testFileData, &ac)
	if err != nil {
		t.Fatal(err)
	}

	if ac.ActionType != ActionTypeEnumAlert {
		t.Fatal()
	}
	if ac.Condition != nil {
		t.Fatal()
	}
	a := ac.AlertAction
	if a == nil || !a.Validate() {
		t.Fatal()
	}
	if a.Title != "For real?" {
		t.Fatal()
	}
	if a.Message != "" {
		t.Fatal()
	}
	if a.ShowCancelButton || !a.ShowOkButton {
		t.Fatal()
	}
	if a.OkButtonActionName != "" {
		t.Fatal()
	}
	if a.Style != AlertActionStyleEnumDialog {
		t.Fatal()
	}
	if len(a.CustomButtons) != 0 {
		t.Fatal()
	}
}

func TestJsonParsingOkayDisabledAlert(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/alert/valid/okDisabled.json")
	if err != nil {
		t.Fatal()
	}
	var ac ActionContainer
	err = json.Unmarshal(testFileData, &ac)
	if err != nil {
		t.Fatal(err)
	}

	if ac.ActionType != ActionTypeEnumAlert {
		t.Fatal()
	}
	a := ac.AlertAction
	if a == nil || !a.Validate() {
		t.Fatal()
	}
	if a.ShowOkButton {
		t.Fatal("failed to parse showOkayButton")
	}
}

func TestParsingInvalidConditionAlert(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/alert/invalid/invalidCondition.json")
	if err != nil {
		t.Fatal()
	}
	var ac ActionContainer
	err = json.Unmarshal(testFileData, &ac)
	if err == nil {
		t.Fatal("invalid condition should return error")
	}

}

func TestJsonParsingInvalidAlert(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/alert/invalid/invalid.json")
	if err != nil {
		t.Fatal()
	}
	var ac ActionContainer
	err = json.Unmarshal(testFileData, &ac)
	if err == nil {
		t.Fatal("invalid json should error")
	}
}

func TestJsonParsingFuture(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/alert/futureproof.json")
	if err != nil {
		t.Fatal()
	}
	var ac ActionContainer
	err = json.Unmarshal(testFileData, &ac)
	if err != nil {
		t.Fatal(err)
	}

	a := ac.AlertAction
	if a == nil || !a.Validate() {
		t.Fatal()
	}
	if a.Title != "hello from the future" {
		t.Fatal()
	}
	if a.Style != AlertActionStyleEnumDialog {
		t.Fatal("didn't fall back on unrecognized style")
	}

	// Strict mode should fail
	StrictDatamodelParsing = true
	defer func() {
		StrictDatamodelParsing = false
	}()
	err = json.Unmarshal(testFileData, &ac)
	if err == nil {
		t.Fatal("Strict parsing allowed unknown style")
	}
}

func TestJsonParsingFutureButton(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/alert/futureproofButtonStyle.json")
	if err != nil {
		t.Fatal()
	}
	var ac ActionContainer
	err = json.Unmarshal(testFileData, &ac)
	if err != nil {
		t.Fatal(err)
	}

	a := ac.AlertAction
	if a == nil || !a.Validate() {
		t.Fatal()
	}
	if a.CustomButtons[0].Style != AlertActionButtonStyleEnumDefault {
		t.Fatal("failed to fallback to default style")
	}

	// Strict mode should fail
	StrictDatamodelParsing = true
	defer func() {
		StrictDatamodelParsing = false
	}()
	err = json.Unmarshal(testFileData, &ac)
	if err == nil {
		t.Fatal("Strict parsing allowed unknown style")
	}
}
