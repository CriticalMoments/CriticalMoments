package datamodel

import (
	"encoding/json"
	"os"
	"testing"
)

func TestImageTypeValidation(t *testing.T) {
	li := LocalImage{}
	if li.Check() == nil {
		t.Fatal("Local images require path")
	}
	li.Path = "image.jpg"
	if li.Check() != nil {
		t.Fatal("Local image failed validation")
	}

	si := SymbolImage{}
	if si.Check() == nil {
		t.Fatal("Symbol images require symbolName")
	}
	si.SymbolName = "upload"
	if si.Check() != nil {
		t.Fatal("Symbol image failed validation")
	}

	si.Weight = "invalid"
	if si.Check() != nil {
		t.Fatal("Symbol image failed validation when it should pass with strict=off")
	}
	StrictDatamodelParsing = true
	defer func() {
		StrictDatamodelParsing = false
	}()
	if si.Check() == nil {
		t.Fatal("Symbol images require valid weight in strict mode")
	}
	si.Weight = SystemSymbolWeightEnumBold
	if si.Check() != nil {
		t.Fatal("Symbol image failed validation")
	}
	StrictDatamodelParsing = false

	si.Mode = "invalid"
	if si.Check() != nil {
		t.Fatal("Symbol image failed validation when it should pass with strict=off")
	}
	StrictDatamodelParsing = true
	defer func() {
		StrictDatamodelParsing = false
	}()
	if si.Check() == nil {
		t.Fatal("Symbol images require valid mode in strict mode")
	}
	si.Mode = SystemSymbolModeEnumHierarchical
	if si.Check() != nil {
		t.Fatal("Symbol image failed validation")
	}
	StrictDatamodelParsing = false

	si.PrimaryColor = "#x"
	if si.Check() == nil {
		t.Fatal("Invalid passed validation")
	}
	si.PrimaryColor = "#ffffff"
	if si.Check() != nil {
		t.Fatal("Symbol image failed validation")
	}

	si.SecondaryColor = "#x"
	if si.Check() == nil {
		t.Fatal("Invalid passed validation")
	}
	si.SecondaryColor = "#ffffff"
	if si.Check() != nil {
		t.Fatal("Symbol image failed validation")
	}
}

func TestJsonParsingImages(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/image/maximalValid.json")
	if err != nil {
		t.Fatal()
	}
	var i Image
	err = json.Unmarshal(testFileData, &i)
	if err != nil {
		t.Fatal(err)
	}

	if i.ImageType != "local" || i.LocalImageData.Path != "image.jpg" {
		t.Fatal("Local image failed to parse")
	}

	i = *i.Fallback
	if i.ImageType != ImageTypeEnumSFSymbol || i.SymbolImageData.SymbolName != "upload" {
		t.Fatal("Symbol image failed to parse")
	}
	if i.SymbolImageData.Mode != "" || i.SymbolImageData.Weight != "" || i.SymbolImageData.PrimaryColor != "" || i.SymbolImageData.SecondaryColor != "" {
		t.Fatal("Symbol defaults failed parse check")
	}

	i = *i.Fallback
	if i.ImageType != ImageTypeEnumSFSymbol || i.SymbolImageData.SymbolName != "download" {
		t.Fatal("Symbol image failed to parse")
	}
	if i.SymbolImageData.Mode != "palette" || i.SymbolImageData.Weight != "light" || i.SymbolImageData.PrimaryColor != "#ff0000" || i.SymbolImageData.SecondaryColor != "#00ff00" {
		t.Fatal("Symbol defaults failed parse check")
	}
}

func TestJsonParsingFutureImages(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/image/futureproof.json")
	if err != nil {
		t.Fatal()
	}
	var i Image
	err = json.Unmarshal(testFileData, &i)
	if err != nil {
		t.Fatal(err)
	}

	if i.SymbolImageData.Mode != SystemSymbolModeEnumMono {
		t.Fatal("Future mode didn't fallback for older client")
	}
	if i.SymbolImageData.Weight != SystemSymbolWeightEnumRegular {
		t.Fatal("Future weight didn't fallback for older client")
	}

	// Strict mode should fail since we have an unknown types
	StrictDatamodelParsing = true
	defer func() {
		StrictDatamodelParsing = false
	}()
	err = json.Unmarshal(testFileData, &i)
	if err == nil {
		t.Fatal("allowed unknown types in symbol")
	}
}
