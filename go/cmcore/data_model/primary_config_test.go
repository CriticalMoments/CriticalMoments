package datamodel

import (
	"encoding/json"
	"os"
	"testing"
)

func testHelperBuildMaxPrimaryConfig(t *testing.T) *PrimaryConfig {
	testFileData, err := os.ReadFile("./test/testdata/primary_config/valid/maximalValid.json")
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
	if pc.DefaultTheme == nil {
		t.Fatal()
	}

	// TODO check all the fields -- full parse checker

	// Check defaults for values not included in json
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
		t.Fatal("Named themes map is empty, and an action uses a missing named theme")
	}

	// fix the broken name mapping above
	banner := pc.namedActions["bannerAction1"]
	banner.ThemeName = ""
	pc.namedActions["bannerAction1"] = banner
	if !pc.Validate() {
		t.Fatal("empty named themes should be allowed")
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

// TODO Min valid file
// TODO extra fields
// TODO breaking a sub-element's validation after parsing, and that "Validate" is recursive (not implmented)
