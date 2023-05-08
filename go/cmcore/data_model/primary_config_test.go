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

func TestNoDefaultThemeVersionTheme(t *testing.T) {
	pc := testHelperBuildMaxPrimaryConfig(t)

	pc.DefaultTheme = nil
	if !pc.Validate() {
		t.Fatal("Not allowing nil default theme, which should be allowed")
	}
}
