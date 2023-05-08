package datamodel

import (
	"encoding/json"
	"os"
	"testing"
)

func TestPrimaryConfigJson(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/primary_config/valid/maximalValid.json")
	if err != nil {
		t.Fatal()
	}
	var pc PrimaryConfig
	uperr := json.Unmarshal(testFileData, &pc)
	if uperr != nil {
		t.Fatal()
	}
	if pc.DefaultTheme == nil {
		t.Fatal()
	}

	// Check defaults for values not included in json
}
