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
	if pc.DefaultTheme() == nil || pc.DefaultTheme().BannerBackgroundColor != "#ffffff" {
		t.Fatal("Default theme not parsed")
	}
	if len(pc.namedThemes) != 3 {
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
	futureThemeWithFallback := pc.ThemeWithName("futureThemeWithFallback")
	if futureThemeWithFallback != blueTheme {
		t.Fatal("Theme that fails validation should fallback")
	}

	// Actions
	if len(pc.namedActions) != 14 {
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
	if failConditionAction == nil || failConditionAction.Condition.String() != "1 > 2" {
		t.Fatal("Didn't parse alert action with failing condition")
	}
	ca1 := pc.ActionWithName("conditionalWithTrueCondition")
	if ca1.ConditionalAction == nil || ca1.ConditionalAction.Condition.String() != "2 > 1" {
		t.Fatal("Didn't parse conditional action 1")
	}
	ca2 := pc.ActionWithName("conditionalWithFalseCondition")
	if ca2.ConditionalAction == nil || ca2.ConditionalAction.Condition.String() != "1 > 2" {
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
	ra := pc.ActionWithName("reviewAction")
	_, ok = ra.actionData.(*ReviewAction)
	if !ok {
		t.Fatal("Review action failed to parse")
	}
	ma := pc.ActionWithName("modalAction")
	if ma.ModalAction == nil || len(ma.ModalAction.Content.Sections) != 1 {
		t.Fatal("failed to parse modal action")
	}
	fa := pc.ActionWithName("futureAction")
	_, ok = fa.actionData.(*UnknownAction)
	if fa.ActionType != "future_action_type" || !ok || fa.FallbackActionName != "alertAction" {
		t.Fatal("unknown action failed to parse. Old client will break for future config files.")
	}
	nfa := pc.ActionWithName("nestedFutureTypeFail")
	_, ok = nfa.actionData.(*UnknownAction)
	if nfa.ActionType != "future_action_type" || !ok || nfa.FallbackActionName != "unknownActionTypeFutureProof" {
		t.Fatal("unknown action failed to parse with fallback. Old client will break for future config files.")
	}
	nfas := pc.ActionWithName("nestedFutureTypeSuccess")
	_, ok = nfas.actionData.(*UnknownAction)
	if nfas.ActionType != "future_action_type" || !ok || nfas.FallbackActionName != "futureAction" {
		t.Fatal("unknown action failed to parse with fallback. Old client will break for future config files.")
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

	// Conditions
	if len(pc.namedConditions) != 3 {
		t.Fatal("Wrong condition count")
	}
	c1 := pc.ConditionWithName("trueCondition")
	if c1 == nil || c1.String() != "true" {
		t.Fatal("Issue with true condition")
	}
	c2 := pc.ConditionWithName("falseCondition")
	if c2 == nil || c2.String() != "false" {
		t.Fatal("Issue with true condition")
	}
	c3 := pc.ConditionWithName("complexCondition")
	if c3 == nil || c3.String() != "4 > 3 && os_version =='123'" {
		t.Fatal("complex condition failed")
	}
	c3Var, err := c3.ExtractIdentifiers()
	if err != nil || len(c3Var.Variables) != 1 || c3Var.Variables[0] != "os_version" {
		t.Fatal("complex condition failed to parse")
	}
}

func TestFutureConditionStrictValidation(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/primary_config/invalid/invalidCondition.json")
	if err != nil {
		t.Fatal(err)
	}
	var pc PrimaryConfig
	err = json.Unmarshal(testFileData, &pc)
	if err != nil {
		t.Fatal(err)
	}

	if pc.ConditionWithName("trueCondition").String() != "true" ||
		pc.ConditionWithName("invalidCondition").String() != "" || // non strict parses to empty, which later evals false
		pc.ConditionWithName("backCompatCondition").String() != "future_feature > 3" {
		t.Fatal("Failed to allow unrecognized variable when not in strict mode")
	}

	// Strict mode should fail since we have an unknown var
	StrictDatamodelParsing = true
	defer func() {
		StrictDatamodelParsing = false
	}()
	err = json.Unmarshal(testFileData, &pc)
	if err == nil {
		t.Fatal("failed to error with invalid conditionand strict mode on")
	}
}

func TestFutureTypeStrictValidation(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/primary_config/invalid/strictInvalidActionName.json")
	if err != nil {
		t.Fatal(err)
	}
	var pc PrimaryConfig
	err = json.Unmarshal(testFileData, &pc)
	if err != nil {
		t.Fatal(err)
	}
	ua := pc.ActionWithName("unknownActionTypeFutureProof")
	_, ok := ua.actionData.(*UnknownAction)
	if ua.ActionType != "unknown_future_type" || !ok {
		t.Fatal("unknown action failed to parse. Old client will break for future config files.")
	}

	// Strict mode should fail since we have an unknown section
	StrictDatamodelParsing = true
	defer func() {
		StrictDatamodelParsing = false
	}()
	err = json.Unmarshal(testFileData, &pc)
	if err == nil {
		t.Fatal("Strict parsing allowed unknown action type")
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

	pc.defaultTheme = nil
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
	pc.namedThemes = make(map[string]*Theme)
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
	pc.namedTriggers = make(map[string]*Trigger)
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
	pc.namedActions = map[string]*ActionContainer{}
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

	pc.namedActions[""] = &ActionContainer{}
	if pc.Validate() {
		t.Fatal("Allowed empty key")
	}
	delete(pc.namedActions, "")

	pc.namedThemes[""] = &Theme{}
	if pc.Validate() {
		t.Fatal("Allowed empty key")
	}
	delete(pc.namedThemes, "")

	pc.namedTriggers[""] = &Trigger{}
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

	pc.namedActions["invalidAction"] = &ActionContainer{}
	if pc.Validate() {
		t.Fatal("actions not re-validated")
	}
}

func TestBreakNestedValidationTriggers(t *testing.T) {
	pc := testHelperBuildMaxPrimaryConfig(t)
	if !pc.Validate() {
		t.Fatal()
	}

	pc.namedTriggers["invalidTrigger"] = &Trigger{}
	if pc.Validate() {
		t.Fatal("trigger not re-validated")
	}
}

func TestBreakNestedTheme(t *testing.T) {
	pc := testHelperBuildMaxPrimaryConfig(t)
	if !pc.Validate() {
		t.Fatal()
	}

	pc.defaultTheme = &Theme{} // invalid
	if pc.Validate() {
		t.Fatal("default theme not re-validated")
	}
	pc.defaultTheme = nil // valid
	if !pc.Validate() {
		t.Fatal()
	}
	pc.namedThemes["newInvalidTheme"] = &Theme{}
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
func TestFallbackNameChecks(t *testing.T) {
	invalidJson := []string{
		"invalidFallbackAction1.json",
		"invalidFallbackTheme2.json", // add_test_count
		"invalidFallbackTheme1.json", // add_test_count
	}

	for _, invalidFile := range invalidJson {
		testFileData, err := os.ReadFile("./test/testdata/primary_config/invalid/" + invalidFile)
		if err != nil {
			t.Fatal(err)
		}
		pc := PrimaryConfig{}
		err = json.Unmarshal(testFileData, &pc)
		if err == nil {
			t.Fatal("Failed to disallow invalid fallback")
		}
	}
}

func TestFallbackDefaultTheme(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/primary_config/valid/validFallbackDefaultTheme.json")
	if err != nil {
		t.Fatal(err)
	}
	pc := PrimaryConfig{}
	err = json.Unmarshal(testFileData, &pc)
	if err != nil {
		t.Fatal(err)
	}
	if pc.DefaultTheme() == nil {
		t.Fatal("Failed to fallback to default theme")
	}
	if pc.DefaultTheme() != pc.ThemeWithName("fbTheme") {
		t.Fatal("Failed to fallback to named theme")
	}
	if pc.ThemeWithName("missingName") != nil {
		t.Fatal("missing name not nil")
	}
}
