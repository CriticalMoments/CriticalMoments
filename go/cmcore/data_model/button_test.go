package datamodel

import (
	"encoding/json"
	"os"
	"testing"
)

func TestButtonValidators(t *testing.T) {
	b := Button{}
	b.Style = ButtonStyleEnumLarge
	if b.Check() == nil {
		t.Fatal("Button requires a title")
	}

	b.Title = "Title"
	b.Style = ""
	if b.Check() == nil {
		t.Fatal("Button with invalid style passed validation")
	}

	b.Style = ButtonStyleEnumLarge
	if b.Check() != nil {
		t.Fatal("Button with title failed validation")
	}

	b.Style = "invalidStyle"
	if b.Check() == nil {
		t.Fatal("Invalid button style passes validation")
	}

	for _, style := range buttonStyles {
		b.Style = style
		if b.Check() != nil {
			t.Fatal("Valid Button failed validation")
		}
	}
}

func TestButtonParsing(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/button/validButton.json")
	if err != nil {
		t.Fatal()
	}
	var b Button
	err = json.Unmarshal(testFileData, &b)
	if err != nil {
		t.Fatal(err)
	}
	if b.Title != "title" || b.Style != ButtonStyleEnumLarge || b.ActionName != "action1" || b.PreventDefault != true {
		t.Fatal("failed to parse button")
	}

	// Strict mode should pass
	StrictDatamodelParsing = true
	defer func() {
		StrictDatamodelParsing = false
	}()
	err = json.Unmarshal(testFileData, &b)
	if err != nil {
		t.Fatal(err)
	}
}

func TestButtonFutureParsing(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/button/futureButtonType.json")
	if err != nil {
		t.Fatal()
	}
	var b Button
	err = json.Unmarshal(testFileData, &b)
	if err != nil {
		t.Fatal(err)
	}
	if b.Title != "title" || b.Style != ButtonStyleEnumNormal {
		t.Fatal("failed to parse button with unknon type into normal")
	}

	// Strict mode should fail since type is unknon
	StrictDatamodelParsing = true
	defer func() {
		StrictDatamodelParsing = false
	}()
	err = json.Unmarshal(testFileData, &b)
	if err == nil {
		t.Fatal("Strict mode failed to detect known type")
	}
}

func TestInvalidButtonParsing(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/button/invalidButton.json")
	if err != nil {
		t.Fatal()
	}
	var b Button
	err = json.Unmarshal(testFileData, &b)
	if err == nil {
		t.Fatal("Invalid button parsed")
	}
}
