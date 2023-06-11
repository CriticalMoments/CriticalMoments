package datamodel

import (
	"encoding/json"
	"os"
	"testing"
)

func testHelperBuildMaxPrimaryConfig(t *testing.T) *PrimaryConfig {
	return testHelperBuilPrimaryConfigFromFile(t, "./test/testdata/primary_config/valid/maximalValid.json")
}

func testHelperBuilPrimaryConfigFromFile(t *testing.T, filePath string) *PrimaryConfig {
	testFileData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	var pc PrimaryConfig
	err = json.Unmarshal(testFileData, &pc)
	if err != nil {
		t.Fatal(err)
	}
	return &pc
}

func TestPrimaryConfigJson(t *testing.T) {
	pc := testHelperBuildMaxPrimaryConfig(t)

	if !pc.Validate() {
		t.Fatal(pc.ValidateReturningUserReadableIssue())
	}

	// Check parsing of top level structure, but individual parsers (Themes, Actions) have
	// their own dedicated test files

	// Version
	if pc.ConfigVersion != "v1" {
		t.Fatal("invalid config version parse")
	}

	// Themes
	if pc.DefaultTheme == nil || pc.DefaultTheme.BannerBackgroundColor != "#ffffff" {
		t.Fatal("Default theme not parsed")
	}
	if len(pc.namedThemes) != 2 {
		t.Fatal("Wrong number of named themes")
	}
	blueTheme := pc.ThemeWithName("blueTheme")
	if blueTheme == nil || blueTheme.BannerBackgroundColor != "#00ff00" {
		t.Fatal("Named theme not parsed")
	}
	greenTheme := pc.ThemeWithName("greenTheme")
	if greenTheme == nil || greenTheme.BannerBackgroundColor != "#0000ff" {
		t.Fatal("Named theme not parsed")
	}

	// Actions
	if len(pc.namedActions) != 9 {
		t.Fatal("Wrong number of named actions")
	}
	bannerAction1 := pc.ActionWithName("bannerAction1")
	if bannerAction1 == nil || bannerAction1.BannerAction.Body != "Hello world, but on a banner!" {
		t.Fatal("Didn't parse banner action 1")
	}
	bannerAction2 := pc.ActionWithName("bannerAction2")
	if bannerAction2 == nil || bannerAction2.BannerAction.Body != "Hello world 2, but on a banner!" {
		t.Fatal("Didn't parse banner action 2")
	}
	alertAction := pc.ActionWithName("alertAction")
	if alertAction == nil || alertAction.AlertAction.Title != "Alert title" {
		t.Fatal("Didn't parse alert action")
	}
	linkAction := pc.ActionWithName("linkAction")
	if linkAction == nil || linkAction.LinkAction.UrlString != "https://criticalmoments.io" {
		t.Fatal("Didn't parse link action")
	}
	failConditionAction := pc.ActionWithName("alertActionWithFailingCondition")
	if failConditionAction == nil || failConditionAction.Condition != "1 > 2" {
		t.Fatal("Didn't parse alert action with failing condition")
	}
	ca1 := pc.ActionWithName("conditionalWithTrueCondition")
	if ca1.ConditionalAction == nil || ca1.ConditionalAction.Condition != "2 > 1" {
		t.Fatal("Didn't parse conditional action 1")
	}
	ca2 := pc.ActionWithName("conditionalWithFalseCondition")
	if ca2.ConditionalAction == nil || ca2.ConditionalAction.Condition != "1 > 2" {
		t.Fatal("Didn't parse conditional action 2")
	}
	ca3 := pc.ActionWithName("conditionalWithoutFalseAction")
	if ca3.ConditionalAction == nil || ca3.ConditionalAction.FailedActionName != "" {
		t.Fatal("Didn't parse conditional action 3")
	}
	ua := pc.ActionWithName("unknownActionTypeFutureProof")
	_, ok := ua.actionData.(*UnknownAction)
	if ua.ActionType != "unknown_future_type" || !ok {
		t.Fatal("unknown action failed to parse. Old client will break for future config files.")
	}

	// Triggers
	if len(pc.namedTriggers) != 2 {
		t.Fatal("Wrong trigger count")
	}
	trigger1 := pc.namedTriggers["trigger1"]
	if trigger1.ActionName != "bannerAction1" || trigger1.EventName != "custom_event" {
		t.Fatal("Trigger 1 parsing failed")
	}
	trigger2 := pc.namedTriggers["trigger_alert"]
	if trigger2.ActionName != "alertAction" || trigger2.EventName != "custom_event_alert" {
		t.Fatal("Trigger 2 parsing failed")
	}
}

func TestInvalidConfigVersionTheme(t *testing.T) {
	pc := testHelperBuildMaxPrimaryConfig(t)

	pc.ConfigVersion = "v2"
	if pc.Validate() {
		t.Fatal("invalid config (v2) passed validation")
	}
	pc.ConfigVersion = ""
	if pc.Validate() {
		t.Fatal("invalid config ('') passed validation")
	}
}

func TestNoDefaultTheme(t *testing.T) {
	pc := testHelperBuildMaxPrimaryConfig(t)

	pc.DefaultTheme = nil
	if !pc.Validate() {
		t.Fatal("Not allowing nil default theme, which should be allowed")
	}
}

func TestNoNamedThemes(t *testing.T) {
	pc := testHelperBuildMaxPrimaryConfig(t)

	pc.namedThemes = nil
	if pc.Validate() {
		t.Fatal("Named themes map is nil, and validated")
	}
	pc.namedThemes = make(map[string]Theme)
	if pc.Validate() {
		t.Fatal("Named themes map is empty, and an action uses a missing named theme, but it improperly validated")
	}

	// fix the broken name mapping above
	banner := pc.namedActions["bannerAction1"]
	banner.BannerAction.CustomThemeName = ""
	pc.namedActions["bannerAction1"] = banner
	if !pc.Validate() {
		t.Fatal("empty named themes should be allowed if no one references them")
	}
}

func TestNoNamedTriggers(t *testing.T) {
	pc := testHelperBuildMaxPrimaryConfig(t)

	pc.namedTriggers = nil
	if pc.Validate() {
		t.Fatal("Named triggers map is nil, and validated")
	}
	pc.namedTriggers = make(map[string]Trigger)
	if !pc.Validate() {
		t.Fatal("empty triggers map should be allowed")
	}
}

func TestNoNamedActions(t *testing.T) {
	pc := testHelperBuildMaxPrimaryConfig(t)

	pc.namedActions = nil
	if pc.Validate() {
		t.Fatal("Named actions map is nil, and validated")
	}
	pc.namedActions = map[string]ActionContainer{}
	if pc.Validate() {
		t.Fatal("empty map should fail since still triggers referencing them")
	}

	delete(pc.namedTriggers, "trigger1")
	delete(pc.namedTriggers, "trigger_alert")
	if !pc.Validate() {
		t.Fatal("empty actions should be allowed")
	}
}

func TestEmptyKey(t *testing.T) {
	pc := testHelperBuildMaxPrimaryConfig(t)

	pc.namedActions[""] = ActionContainer{}
	if pc.Validate() {
		t.Fatal("Allowed empty key")
	}
	delete(pc.namedActions, "")

	pc.namedThemes[""] = Theme{}
	if pc.Validate() {
		t.Fatal("Allowed empty key")
	}
	delete(pc.namedThemes, "")

	pc.namedTriggers[""] = Trigger{}
	if pc.Validate() {
		t.Fatal("Allowed empty key")
	}
	delete(pc.namedTriggers, "")

	if !pc.Validate() {
		t.Fatal("Should be valid")
	}
}

func TestBreakNestedValidationActions(t *testing.T) {
	pc := testHelperBuildMaxPrimaryConfig(t)
	if !pc.Validate() {
		t.Fatal()
	}

	pc.namedActions["invalidAction"] = ActionContainer{}
	if pc.Validate() {
		t.Fatal("actions not re-validated")
	}
}

func TestBreakNestedValidationTriggers(t *testing.T) {
	pc := testHelperBuildMaxPrimaryConfig(t)
	if !pc.Validate() {
		t.Fatal()
	}

	pc.namedTriggers["invalidTrigger"] = Trigger{}
	if pc.Validate() {
		t.Fatal("trigger not re-validated")
	}
}

func TestBreakNestedTheme(t *testing.T) {
	pc := testHelperBuildMaxPrimaryConfig(t)
	if !pc.Validate() {
		t.Fatal()
	}

	pc.DefaultTheme = &Theme{} // invalid
	if pc.Validate() {
		t.Fatal("default theme not re-validated")
	}
	pc.DefaultTheme = nil // valid
	if !pc.Validate() {
		t.Fatal()
	}
	pc.namedThemes["newInvalidTheme"] = Theme{}
	if pc.Validate() {
		t.Fatal("named themes not re-validated")
	}
}

func TestPcAccessors(t *testing.T) {
	pc := testHelperBuildMaxPrimaryConfig(t)
	if !pc.Validate() {
		t.Fatal()
	}

	theme := pc.ThemeWithName("doesntExist")
	if theme != nil {
		t.Fatal("Found a theme that doesn't exist")
	}

	theme = pc.ThemeWithName("greenTheme")
	if theme == nil {
		t.Fatal("Couldn't find theme by name")
	}

	action := pc.ActionWithName("nada")
	if action != nil {
		t.Fatal("Found a action that doesn't exist")
	}

	action = pc.ActionWithName("bannerAction1")
	if action == nil {
		t.Fatal("Couldn't find action by name")
	}

	actions := pc.ActionsForEvent("nada")
	if len(actions) > 0 {
		t.Fatal("Found a action that doesn't exist")
	}

	actions = pc.ActionsForEvent("custom_event")
	if len(actions) != 1 || actions[0].ActionType != ActionTypeEnumBanner {
		t.Fatal("Didn't find action from event")
	}
}

func TestMinValidConfig(t *testing.T) {
	pc := testHelperBuilPrimaryConfigFromFile(t, "./test/testdata/primary_config/valid/minimalValid.json")
	if !pc.Validate() {
		t.Fatal(pc.ValidateReturningUserReadableIssue())
	}
	if pc.ConfigVersion != "v1" {
		t.Fatal("Failed to parse config version")
	}
}

func TestOddballValidConfig(t *testing.T) {
	pc := testHelperBuilPrimaryConfigFromFile(t, "./test/testdata/primary_config/valid/oddballValid.json")
	if !pc.Validate() {
		t.Fatal(pc.ValidateReturningUserReadableIssue())
	}
	if len(pc.namedActions) != 0 || len(pc.namedThemes) != 0 || len(pc.namedTriggers) != 0 {
		t.Fatal("Expected oddball with empty maps")
	}
	if pc.ConfigVersion != "v1" {
		t.Fatal("Failed to parse config version")
	}
}

func TestMinWithUnknownFieldConfig(t *testing.T) {
	// https://github.com/golang/go/issues/41144
	t.Skip("This test fails -- we would need to implement strict decoding separately, which is non trivial")
	filePath := "./test/testdata/primary_config/invalid/minimalWithUnknownField.json"
	reader, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}

	// Without strict it should work
	decoder := json.NewDecoder(reader)
	var pc PrimaryConfig
	err = decoder.Decode(&pc)
	if err != nil || !pc.Validate() {
		t.Fatal("Non strict parsing failed")
	}

	// with strict it should fail
	reader, err = os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}
	strictDecoder := json.NewDecoder(reader)
	strictDecoder.DisallowUnknownFields()
	var strictPc PrimaryConfig
	err = strictDecoder.Decode(&strictPc)
	if err == nil {
		t.Fatal("Strict parsing ignored extra field")
	}
}
