package datamodel

import (
	"encoding/json"
	"os"
	"testing"
)

func TestModalValidation(t *testing.T) {
	m := ModalAction{}
	if m.Validate() {
		t.Fatal("Invalid content passed validation")
	}

	s := PageSection{
		PageSectionType: SectionTypeEnumBodyText,
	}
	s.BodyData = &BodyPageSection{
		ScaleFactor: 1.0,
		BodyText:    "body",
	}
	s.pageSectionData = s.BodyData
	m.Content = &Page{
		Sections: []*PageSection{&s},
	}
	if !m.Validate() {
		t.Fatal("valid content failed validation")
	}

	// theme name extraction
	if themes, err := m.AllEmbeddedThemeNames(); err != nil || len(themes) != 0 {
		t.Fatal("Theme listed when none specified")
	}
	m.CustomThemeName = "theme1"
	if themes, err := m.AllEmbeddedThemeNames(); err != nil || len(themes) != 1 || themes[0] != "theme1" {
		t.Fatal("Theme not listed")
	}

	// button action extraction
	if actions, err := m.AllEmbeddedActionNames(); err != nil || len(actions) != 0 {
		t.Fatal("actions included when none specified")
	}
	m.Content.Buttons = []*Button{
		{
			Title:      "button1",
			ActionName: "action1",
		},
	}
	if actions, err := m.AllEmbeddedActionNames(); err != nil || len(actions) != 1 || actions[0] != "action1" {
		t.Fatal("actions included when none specified")
	}
}

func TestJsonParsingModal(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/modal/maximalValid.json")
	if err != nil {
		t.Fatal()
	}
	var m ModalAction
	err = json.Unmarshal(testFileData, &m)
	if err != nil {
		t.Fatal(err)
	}

	if m.ShowCloseButton || m.CustomThemeName != "theme1" || len(m.Content.Sections) != 1 || len(m.Content.Buttons) != 1 {
		t.Fatal("error parsing modal")
	}
}

func TestJsonParsingMinModal(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/modal/minimalValid.json")
	if err != nil {
		t.Fatal()
	}
	var m ModalAction
	err = json.Unmarshal(testFileData, &m)
	if err != nil {
		t.Fatal(err)
	}

	if !m.ShowCloseButton || m.CustomThemeName != "" || len(m.Content.Sections) != 1 || len(m.Content.Buttons) != 0 {
		t.Fatal("error parsing modal")
	}
}

func TestJsonParsingInvalidModal(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/modal/invalid.json")
	if err != nil {
		t.Fatal()
	}
	var m ModalAction
	err = json.Unmarshal(testFileData, &m)
	if err != nil {
		t.Fatal("should allow unrecognized content for backwards compat when not strict")
	}

	// Strict mode should fail
	StrictDatamodelParsing = true
	defer func() {
		StrictDatamodelParsing = false
	}()
	err = json.Unmarshal(testFileData, &m)
	if err == nil {
		t.Fatal("allowed invalid modal parsing in strict mode")
	}
}
