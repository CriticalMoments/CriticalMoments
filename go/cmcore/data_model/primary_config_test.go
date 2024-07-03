package datamodel

import (
	"encoding/json"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/CriticalMoments/CriticalMoments/go/cmcore/signing"
)

var testContainer = `
-----BEGIN CM-----
Container-Version: v1
additional-future-header: value2

-----END CM-----
-----BEGIN CONFIG-----
Signature: MD0CHQDVjg2+dxUL47DvlctnjAcObzCKvvM6mp2tT507Ahwsvr2Zs5vgE8BSgai6XIMGaFL3CZshtgFubvpq
ewogICAgImNvbmZpZ1ZlcnNpb24iOiAidjEiLAogICAgImFwcElkIjogImlvLmNyaXRpY2FsbW9tZW50cy5zYW1wbGUtYXBwIgp9
-----END CONFIG-----
-----BEGIN FUTUREBLOCK-----
aGVsbG8gd29ybGQ=
-----END FUTUREBLOCK-----
`

var testPrivKey = "MGgCAQEEHOEUmigOOoZ+STQ1jkYuXRN2hXLbxLKTvKdlXEygBwYFK4EEACGhPAM6AASDljuXqf/dic4vnAfRtqFsl/fQANciY+xACkgOOE9MGgvu+XIfTOqsqagLJ6ZUedbZus5FUa4awQ=="
var testPubKey = "ME4wEAYHKoZIzj0CAQYFK4EEACEDOgAEg5Y7l6n/3YnOL5wH0bahbJf30ADXImPsQApIDjhPTBoL7vlyH0zqrKmoCyemVHnW2brORVGuGsE="

func TestContainer(t *testing.T) {
	su, err := signing.NewSignUtilWithSerializedPublicKey(testPubKey)
	if err != nil {
		t.Fatal(err)
	}
	pc, err := DecodePrimaryConfig([]byte(testContainer), su)
	if err != nil {
		t.Fatal(err)
	}

	if pc.ContainerVersion != "v1" {
		t.Fatal("Failed to parse container version")
	}
	if pc.AppId != "io.criticalmoments.sample-app" {
		t.Fatal("Failed to parse app id")
	}
	if pc.ConfigVersion != "v1" {
		t.Fatal("Failed to parse config version from JSON")
	}
}

func TestContainerVersionCheck(t *testing.T) {
	su, err := signing.NewSignUtilWithSerializedPublicKey(testPubKey)
	if err != nil {
		t.Fatal(err)
	}
	subversion := strings.Replace(testContainer, "Container-Version: v1", "Container-Version: v1.2", 1)
	pc, err := DecodePrimaryConfig([]byte(subversion), su)
	if err != nil || pc == nil {
		t.Fatal(err)
	}

	v2 := strings.Replace(testContainer, "Container-Version: v1", "Container-Version: v2", 1)
	pc, err = DecodePrimaryConfig([]byte(v2), su)
	if err == nil || pc != nil {
		t.Fatal("Failed to error on containerVersion = v2")
	}
}

func TestContainerInvalidSig(t *testing.T) {
	// Signature is correct format, and signed by correct key, but does not match this body
	var testContainerInvalidSig = `
-----BEGIN CM-----
Container-Version: v1
aGVsbG8gd29ybGQ=
-----END CM-----
-----BEGIN CONFIG-----
Signature: MD4CHQDWjw/kUoBUF5C4M1rLtYSHdcpkLBkH0vGYfSrRAh0AyV/+yosj2C2hqybZEsWYU/x4bPeP2soQF+2cIQ==
ewogICAgImNvbmZpZ1ZlcnNpb24iOiAidjEiLAogICAgImFwcElkIjogImlvLmNyaXRpY2FsbW9tZW50cy5zYW1wbGUtYXBwIgp9
-----END CONFIG-----
`

	su, err := signing.NewSignUtilWithSerializedPublicKey(testPubKey)
	if err != nil {
		t.Fatal(err)
	}
	pc, err := DecodePrimaryConfig([]byte(testContainerInvalidSig), su)
	if err == nil || pc != nil {
		t.Fatal("Failed to error on pc with invalid sig")
	}
}

func TestContainerEncode(t *testing.T) {
	su, err := signing.NewSignUtilWithSerializedPrivateKey(testPrivKey)
	if err != nil {
		t.Fatal(err)
	}

	// Invalid
	_, err = EncodeConfig([]byte("not json"), su)
	if err == nil {
		t.Fatal("Failed to error on invalid json")
	}

	// Valid
	c := []byte("{\"configVersion\": \"v1\",\"appId\": \"io.criticalmoments.demo\"}")
	b, err := EncodeConfig(c, su)
	if err != nil {
		t.Fatal(err)
	}
	pc, err := DecodePrimaryConfig(b, su)
	if err != nil || pc == nil {
		t.Fatal("Failed to decode encoded config", err)
	}
	if pc.ContainerVersion != "v1" || pc.ConfigVersion != "v1" || pc.AppId != "io.criticalmoments.demo" {
		t.Fatal("Failed to decode encoded config")
	}
}

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

func testAllConditionsContainsCondition(t *testing.T, pc *PrimaryConfig, c *Condition) {
	all, err := pc.AllConditions()
	if err != nil {
		t.Fatal(err)
	}

	for _, cond := range all {
		// Pointer check, check actual instance returned in case string matches
		if cond == c {
			return
		}
	}
	t.Fatalf("Failed to find condition in all conditions: %v", c.String())
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
	if pc.MinAppVersion != "1.0.0" {
		t.Fatal("invalid min app version parse")
	}
	if pc.MinCMVersion != "0.8.0" {
		t.Fatal("invalid min cm version parse")
	}
	if pc.MinCMVersionInternal != "0.7.0" {
		t.Fatal("invalid min cm version parse")
	}

	// Themes
	if pc.DefaultTheme() == nil || pc.DefaultTheme().BannerBackgroundColor != "#ffffff" {
		t.Fatal("Default theme not parsed")
	}
	if len(pc.namedThemes) != 4 {
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
	if bannerAction1.BannerAction.CustomThemeName != "blueTheme" {
		t.Fatal("Didn't parse banner action 1 theme")
	}
	bannerAction2 := pc.ActionWithName("bannerAction2")
	if bannerAction2 == nil || bannerAction2.BannerAction.Body != "Hello world 2, but on a banner!" {
		t.Fatal("Didn't parse banner action 2")
	}
	if bannerAction2.BannerAction.CustomThemeName != "elegant" {
		t.Fatal("Didn't parse banner action 2 theme")
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
	testAllConditionsContainsCondition(t, pc, failConditionAction.Condition)
	ca1 := pc.ActionWithName("conditionalWithTrueCondition")
	if ca1.ConditionalAction == nil || ca1.ConditionalAction.Condition.String() != "2 > 1" {
		t.Fatal("Didn't parse conditional action 1")
	}
	testAllConditionsContainsCondition(t, pc, ca1.ConditionalAction.Condition)
	ca2 := pc.ActionWithName("conditionalWithFalseCondition")
	if ca2.ConditionalAction == nil || ca2.ConditionalAction.Condition.String() != "1 > 2" {
		t.Fatal("Didn't parse conditional action 2")
	}
	testAllConditionsContainsCondition(t, pc, ca2.ConditionalAction.Condition)
	ca3 := pc.ActionWithName("conditionalWithoutFalseAction")
	if ca3.ConditionalAction == nil || ca3.ConditionalAction.FailedActionName != "" {
		t.Fatal("Didn't parse conditional action 3")
	}
	testAllConditionsContainsCondition(t, pc, ca3.ConditionalAction.Condition)
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
	if len(pc.AllActions()) != len(pc.namedActions) {
		t.Fatal("all actions count mismatch")
	}

	// Triggers
	if len(pc.namedTriggers) != 4 {
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
	trigger3 := pc.namedTriggers["conditional_trigger_true"]
	if trigger3.ActionName != "alertAction" || trigger3.EventName != "custom_event_conditional_true" || trigger3.Condition.String() != "2 > 1" {
		t.Fatal("Trigger 3 parsing failed")
	}
	testAllConditionsContainsCondition(t, pc, trigger3.Condition)
	trigger4 := pc.namedTriggers["conditional_trigger_false"]
	if trigger4.ActionName != "alertAction" || trigger4.EventName != "custom_event_conditional_false" || trigger4.Condition.String() != "2 > 3" {
		t.Fatal("Trigger 4 parsing failed")
	}
	testAllConditionsContainsCondition(t, pc, trigger4.Condition)

	// Conditions
	if len(pc.namedConditions) != 3 {
		t.Fatal("Wrong condition count")
	}
	if pc.NamedConditionCount() != 3 {
		t.Fatal("Named condition count mismatch")
	}
	if !slices.Contains(pc.NamedConditionsConditionals(), "4 > 3 && os_version =='123'") {
		t.Fatal("Named condition incorrect")
	}
	if len(pc.NamedConditionsConditionals()) != 3 {
		t.Fatal("Named condition count mismatch")
	}
	c1 := pc.ConditionWithName("trueCondition")
	if c1 == nil || c1.String() != "true" {
		t.Fatal("Issue with true condition")
	}
	testAllConditionsContainsCondition(t, pc, c1)
	c2 := pc.ConditionWithName("falseCondition")
	if c2 == nil || c2.String() != "false" {
		t.Fatal("Issue with true condition")
	}
	testAllConditionsContainsCondition(t, pc, c2)
	c3 := pc.ConditionWithName("complexCondition")
	if c3 == nil || c3.String() != "4 > 3 && os_version =='123'" {
		t.Fatal("complex condition failed")
	}
	testAllConditionsContainsCondition(t, pc, c3)
	c3Var, err := c3.ExtractIdentifiers()
	if err != nil || len(c3Var.Variables) != 1 || c3Var.Variables[0] != "os_version" {
		t.Fatal("complex condition failed to parse")
	}

	// Notifications - tested in more depth in notification_test.go
	if len(pc.Notifications) != 2 {
		t.Fatal("Failed to parse notifications")
	}
	if pc.Notifications["testNotification"] == nil {
		t.Fatal("Failed to parse notification id")
	}
	if pc.Notifications["testNotification"].ID != "testNotification" {
		t.Fatal("Failed to parse notification id")
	}
	if pc.Notifications["testNotification"].Title != "Notification title" {
		t.Fatal("Failed to parse notification title")
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
	StrictDatamodelParsing = true
	defer func() {
		StrictDatamodelParsing = false
	}()

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

	pc.namedTriggers = make(map[string]*Trigger)
	pc.Notifications = make(map[string]*Notification)
	if !pc.Validate() {
		t.Fatal("empty actions should be allowed when no triggers or notifications reference them")
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

	triggers := pc.TriggersForEvent("nada")
	if len(triggers) > 0 {
		t.Fatal("Found a action that doesn't exist")
	}

	triggers = pc.TriggersForEvent("custom_event")
	if len(triggers) != 1 || pc.ActionWithName(triggers[0].ActionName).ActionType != ActionTypeEnumBanner {
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
	if pc.AppId != "io.criticalmoments.demo" {
		t.Fatal("Failed to parse config version")
	}
}

func TestOddballValidConfig(t *testing.T) {
	pc := testHelperBuilPrimaryConfigFromFile(t, "./test/testdata/primary_config/valid/oddballvalid.json")
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
func TestInvalidsError(t *testing.T) {
	invalidJson := []string{
		// Fallback name checks
		"invalidFallbackAction1.json",
		"invalidFallbackTheme2.json", // add_test_count
		"invalidFallbackTheme1.json", // add_test_count
		// Invalid notification action name
		"invalidNotificationAction.json", // add_test_count
	}

	for _, invalidFile := range invalidJson {
		testFileData, err := os.ReadFile("./test/testdata/primary_config/invalid/" + invalidFile)
		if err != nil {
			t.Fatal(err)
		}
		pc := PrimaryConfig{}
		err = json.Unmarshal(testFileData, &pc)
		if err == nil {
			t.Fatalf("Failed to disallow invalid primary config: %v\n\n%v", invalidFile, err)
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

func TestDefaultThemeSelection(t *testing.T) {
	// Strict mode
	StrictDatamodelParsing = true
	defer func() {
		StrictDatamodelParsing = false
	}()

	testFileData, err := os.ReadFile("./test/testdata/primary_config/valid/builtInLibraryTheme.json")
	if err != nil {
		t.Fatal(err)
	}
	pc := PrimaryConfig{}
	err = json.Unmarshal(testFileData, &pc)
	if err != nil {
		t.Fatal(err)
	}
	if pc.DefaultTheme() != nil {
		t.Fatal("Libary theme should not set default")
	}
	if pc.LibraryThemeName != "system_dark" {
		t.Fatal("Failed to parse library theme name")
	}

	testFileData, err = os.ReadFile("./test/testdata/primary_config/valid/builtInDefaultTheme.json")
	if err != nil {
		t.Fatal(err)
	}
	pc = PrimaryConfig{}
	err = json.Unmarshal(testFileData, &pc)
	if err != nil {
		t.Fatal(err)
	}
	// Also checks ThemeWithName works on built in themes
	if pc.DefaultTheme() != pc.ThemeWithName("elegant_light") {
		t.Fatal("Libary theme should not set default")
	}
	if pc.LibraryThemeName != "" {
		t.Fatal("Set Library theme name when not needed")
	}

	// named config based default tested in maximal valid test case

	testFileData, err = os.ReadFile("./test/testdata/primary_config/invalid/invalidDefaultTheme.json")
	if err != nil {
		t.Fatal(err)
	}
	pc = PrimaryConfig{}
	err = json.Unmarshal(testFileData, &pc)
	if err == nil {
		t.Fatal("Failed to error on invalid default theme")
	}
}

func TestFutureBuiltInTheme(t *testing.T) {
	testFiles := []string{
		"./test/testdata/primary_config/valid/futureBuiltInActionTheme.json",
		"./test/testdata/primary_config/valid/futureBuiltInTheme.json", // add_test_count
	}
	for _, file := range testFiles {
		StrictDatamodelParsing = false
		testFileData, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		// Non strict should allow unknown theme names, could be future built in
		pc := PrimaryConfig{}
		err = json.Unmarshal(testFileData, &pc)
		if err != nil {
			t.Fatal(err)
		}

		// Strict mode
		StrictDatamodelParsing = true
		defer func() {
			StrictDatamodelParsing = false
		}()

		// strict mode should not allow unknown theme names
		pc = PrimaryConfig{}
		err = json.Unmarshal(testFileData, &pc)
		if err == nil {
			t.Fatal("allowed unknown theme name in strict mode")
		}
	}
}

func TestMinAppAndClientVersionNumberValidation(t *testing.T) {
	pc := testHelperBuilPrimaryConfigFromFile(t, "./test/testdata/primary_config/valid/minimalValid.json")
	if !pc.Validate() {
		t.Fatal(pc.ValidateReturningUserReadableIssue())
	}
	if pc.ConfigVersion != "v1" {
		t.Fatal("Failed to parse config version")
	}

	pc.MinAppVersion = "invalid"
	if pc.Validate() {
		t.Fatal("failed to validate MinAppVersion")
	}
	pc.MinAppVersion = ""
	pc.MinCMVersion = "invalid"
	if pc.Validate() {
		t.Fatal("failed to validate MinCMVersion")
	}

	pc.MinCMVersion = ""
	pc.MinCMVersionInternal = "invalid"
	if pc.Validate() {
		t.Fatal("failed to validate MinCMVersionInternal")
	}

	pc.MinAppVersion = "1.2.3.4.5"
	pc.MinCMVersion = "v34.234234.234.1123.32"
	pc.MinCMVersionInternal = "v34.234234.234.1123.32"
	if !pc.Validate() {
		t.Fatal(pc.ValidateReturningUserReadableIssue())
	}
}
