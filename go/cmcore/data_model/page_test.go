package datamodel

import (
	"encoding/json"
	"os"
	"testing"
)

func TestPageParsing(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/page/maximalValid.json")
	if err != nil {
		t.Fatal()
	}
	var p Page
	err = json.Unmarshal(testFileData, &p)
	if err != nil {
		t.Fatal(err)
	}

	if len(p.Sections) != 5 {
		t.Fatal("Parsing page failed")
	}

	s1 := p.Sections[0]
	if s1.PageSectionType != SectionTypeEnumTitle || s1.TopSpacingScale != 1.1 {
		t.Fatal("Parsing page failed")
	}
	s1d := s1.TitleData
	if s1d.Title != "title1" || s1d.ScaleFactor != 1.2 || s1d.Bold != false {
		t.Fatal("Parsing page failed")
	}

	s2 := p.Sections[1]
	if s2.PageSectionType != SectionTypeEnumTitle || s2.TopSpacingScale != 1.0 {
		t.Fatal("Parsing page failed")
	}
	s2d := s2.TitleData
	if s2d.Title != "title2" || s2d.ScaleFactor != 1 || s2d.Bold != true {
		t.Fatal("Parsing page failed")
	}

	if _, ok := p.Sections[2].pageSectionData.(UnknownSection); !ok {
		t.Fatal("failed to parse future action name to unknown")
	}

	s4 := p.Sections[3]
	if s4.PageSectionType != SectionTypeEnumBodyText || s4.TopSpacingScale != 1.0 {
		t.Fatal("Parsing page failed")
	}
	s4d := s4.BodyData
	if s4d.BodyText != "body1" || s4d.ScaleFactor != 1 || s4d.Bold != false || s4d.CenterText != true || s4d.UsePrimaryTextColor != false {
		t.Fatal("Parsing page failed")
	}

	s5 := p.Sections[4]
	if s5.PageSectionType != SectionTypeEnumBodyText || s5.TopSpacingScale != 1.0 {
		t.Fatal("Parsing page failed")
	}
	s5d := s5.BodyData
	if s5d.BodyText != "body2" || s5d.ScaleFactor != 1.1 || s5d.Bold != true || s5d.CenterText != false || s5d.UsePrimaryTextColor != true {
		t.Fatal("Parsing page failed")
	}

	// Strict mode should fail since we have an unknown section at index=2
	StrictDatamodelParsing = true
	defer func() {
		StrictDatamodelParsing = false
	}()
	err = json.Unmarshal(testFileData, &p)
	if err == nil {
		t.Fatal("Failed to error in strict mode, when unknown type present")
	}
	StrictDatamodelParsing = false

	// buttons
	if len(p.Buttons) != 7 {
		t.Fatal("Page button parsing failed")
	}

	b1 := p.Buttons[0]
	if b1.Title != "button1" || b1.ActionName != "action1" || b1.PreventDefault != true || b1.Style != "large" {
		t.Fatal("Button failed to parse")
	}

	b2 := p.Buttons[1]
	if b2.Title != "button2" || b2.ActionName != "action2" || b2.PreventDefault != false || b2.Style != "normal" {
		t.Fatal("Button failed to parse")
	}

	b3 := p.Buttons[2]
	if b3.Title != "button3" || b3.ActionName != "" || b3.PreventDefault != false || b3.Style != "normal" {
		t.Fatal("Button failed to parse")
	}

	b4 := p.Buttons[3]
	if b4.Title != "button4" || b4.Style != "secondary" {
		t.Fatal("Button failed to parse")
	}

	b5 := p.Buttons[4]
	if b5.Title != "button5" || b5.Style != "tertiary" {
		t.Fatal("Button failed to parse")
	}

	b6 := p.Buttons[5]
	if b6.Title != "button6" || b6.Style != "info" {
		t.Fatal("Button failed to parse")
	}

	b7 := p.Buttons[6]
	if b7.Title != "button7" || b7.Style != "info-small" {
		t.Fatal("Button failed to parse")
	}
}

func TestPageParsingMinimal(t *testing.T) {
	testFileData, err := os.ReadFile("./test/testdata/actions/page/minimalValid.json")
	if err != nil {
		t.Fatal()
	}
	var p Page
	err = json.Unmarshal(testFileData, &p)
	if err != nil {
		t.Fatal(err)
	}

	if len(p.Sections) != 1 {
		t.Fatal("Parsing page failed")
	}

	s1 := p.Sections[0]
	if s1.PageSectionType != SectionTypeEnumTitle || s1.TopSpacingScale != 1.0 || s1.TitleData.Title != "title1" {
		t.Fatal("Parsing page failed")
	}

	// Strict mode should succeed since all types are known
	StrictDatamodelParsing = true
	defer func() {
		StrictDatamodelParsing = false
	}()
	err = json.Unmarshal(testFileData, &p)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPageValidation(t *testing.T) {
	p := Page{}
	if p.ValidateReturningUserReadableIssue() == "" {
		t.Fatal("Allowed page with no sections")
	}

	p.Sections = []*PageSection{
		{
			PageSectionType: SectionTypeEnumTitle,
			TitleData: &TitlePageSection{
				Title:       "title",
				ScaleFactor: 1.0,
			},
		},
	}
	p.Sections[0].pageSectionData = p.Sections[0].TitleData
	if verr := p.ValidateReturningUserReadableIssue(); verr != "" {
		t.Fatal("Didn't allow valid page: " + verr)
	}

	// Needs to be future proof, client accept/ignore types it doesn't recognize
	p.Sections[0].PageSectionType = "invalidType"
	if p.ValidateReturningUserReadableIssue() != "" {
		t.Fatal("Didn't allow valid page")
	}

	p.Sections[0].PageSectionType = SectionTypeEnumTitle
	p.Sections[0].TitleData.Title = ""
	if p.ValidateReturningUserReadableIssue() == "" {
		t.Fatal("Allowed page with invalid subsection")
	}
	p.Sections[0].TitleData.Title = "title"

	// test buttons
	p.Buttons = []*Button{
		{},
	}
	if p.ValidateReturningUserReadableIssue() == "" {
		t.Fatal("Allowed page with invalid button")
	}
	p.Buttons = []*Button{
		{
			Title: "button1",
			Style: ButtonStyleEnumNormal,
		},
	}
	if p.ValidateReturningUserReadableIssue() != "" {
		t.Fatal("Disallowed valid button")
	}
}

func TestTitlePageSectionValidation(t *testing.T) {
	s := PageSection{
		PageSectionType: SectionTypeEnumTitle,
	}
	if s.ValidateReturningUserReadableIssue() == "" {
		t.Fatal("Allowed title section with no title data")
	}
	s.TitleData = &TitlePageSection{
		ScaleFactor: 1.0,
	}
	s.pageSectionData = s.TitleData
	if s.ValidateReturningUserReadableIssue() == "" {
		t.Fatal("Allowed title section with no title string")
	}
	s.TitleData.Title = "title"
	if s.ValidateReturningUserReadableIssue() != "" {
		t.Fatal("Valid title failed validation")
	}
	s.TitleData.ScaleFactor = -1.0
	if s.ValidateReturningUserReadableIssue() == "" {
		t.Fatal("Allowed title section with negative scale")
	}
}

func TestBodyPageSectionValidation(t *testing.T) {
	s := PageSection{
		PageSectionType: SectionTypeEnumBodyText,
	}
	if s.ValidateReturningUserReadableIssue() == "" {
		t.Fatal("Allowed body section with no data")
	}
	s.BodyData = &BodyPageSection{
		ScaleFactor: 1.0,
	}
	s.pageSectionData = s.BodyData
	if s.ValidateReturningUserReadableIssue() == "" {
		t.Fatal("Allowed body section with no body")
	}
	s.BodyData.BodyText = "body"
	if s.ValidateReturningUserReadableIssue() != "" {
		t.Fatal("Valid body failed validation")
	}
	s.BodyData.ScaleFactor = -1.0
	if s.ValidateReturningUserReadableIssue() == "" {
		t.Fatal("Allowed section with negative scale")
	}
}
