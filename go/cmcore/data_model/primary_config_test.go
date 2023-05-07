package datamodel

import (
	"os"
	"testing"
)

func TestPrimaryConfigJson(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/primary_config/valid/maximalValid.json")
	pc, err := NewPrimaryConfigFromJson(testFileData)
	if err != nil {
		t.Fatal()
	}
	if pc.DefaultTheme == nil {
		t.Fatal()
	}

	// Check defaults for values not included in json
}
