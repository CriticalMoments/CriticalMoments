package datamodel

import (
	"encoding/json"
	"fmt"
)

type Page struct {
	Sections []*PageSection
	Buttons  []*Button
}

type jsonPage struct {
	Sections []*PageSection `json:"sections"`
	Buttons  []*Button      `json:"buttons"`
}

func (p *Page) UnmarshalJSON(data []byte) error {
	var jp jsonPage
	err := json.Unmarshal(data, &jp)
	if err != nil {
		return NewUserPresentableErrorWSource("Unable to parse the json of a page.", err)
	}

	p.Sections = jp.Sections
	p.Buttons = jp.Buttons

	if validationIssue := p.ValidateReturningUserReadableIssue(); validationIssue != "" {
		return NewUserPresentableError(validationIssue)
	}

	return nil
}

func (p *Page) ValidateReturningUserReadableIssue() string {
	if len(p.Sections) == 0 {
		return "Page with 0 sections is not valid"
	}

	for _, section := range p.Sections {
		if valErr := section.ValidateReturningUserReadableIssue(); valErr != "" {
			return valErr
		}
	}

	for _, button := range p.Buttons {
		if valErr := button.ValidateReturningUserReadableIssue(); valErr != "" {
			return valErr
		}
	}

	return ""
}

const (
	SectionTypeEnumTitle    string = "title"
	SectionTypeEnumBodyText string = "body"
)

type pageSectionTypeInterface interface {
	ValidateReturningUserReadableIssue() string
}

var (
	pageSectionTypeRegistry = map[string]func(map[string]interface{}, *PageSection) (pageSectionTypeInterface, error){
		SectionTypeEnumTitle:    unpackTitleSection,
		SectionTypeEnumBodyText: unpackBodySection,
	}
)

type PageSection struct {
	PageSectionType string
	TopSpacingScale float64

	// Section types, stronly typed for easy consumption
	TitleData *TitlePageSection
	BodyData  *BodyPageSection

	// generalized interface for functions we need for any section type.
	// Typically a pointer to the one value above that is populated.
	pageSectionData pageSectionTypeInterface
}

type jsonPageSection struct {
	PageSectionType string                 `json:"pageSectionType"`
	TopSpacingScale *float64               `json:"topSpacingScale,omitempty"`
	RawSectionData  map[string]interface{} `json:"pageSectionData,omitempty"`
}

func (s *PageSection) UnmarshalJSON(data []byte) error {
	var js jsonPageSection
	err := json.Unmarshal(data, &js)
	if err != nil {
		return NewUserPresentableErrorWSource("Unable to parse the json of a page.", err)
	}

	s.PageSectionType = js.PageSectionType

	s.TopSpacingScale = 1.0
	if js.TopSpacingScale != nil && *js.TopSpacingScale >= 0.0 {
		s.TopSpacingScale = *js.TopSpacingScale
	}

	unpacker, ok := pageSectionTypeRegistry[js.PageSectionType]
	if !ok {
		errString := fmt.Sprintf("CriticalMoments: Unsupported section type: \"%v\" found in config file. This section will be ignored. If unexpected, check the CM config file.\n", s.PageSectionType)
		if StrictDatamodelParsing {
			return NewUserPresentableError(errString)
		} else {
			// Forward compatibility: warn them the type is unrecognized in debug console, but could be newer config on older build so no hard error
			fmt.Println(errString)
			s.pageSectionData = UnknownSection{}
		}
	} else {
		pageSectionData, err := unpacker(js.RawSectionData, s)
		if err != nil {
			return err
		}
		s.pageSectionData = pageSectionData
	}

	if validationIssue := s.ValidateReturningUserReadableIssue(); validationIssue != "" {
		return NewUserPresentableError(validationIssue)
	}

	return nil
}

func (s *PageSection) ValidateReturningUserReadableIssue() string {
	if s.TopSpacingScale < 0 {
		return "Top space scale for a page section must be >= 0"
	}

	if s.pageSectionData == nil {
		return "Invalid section in page"
	}
	if verr := s.pageSectionData.ValidateReturningUserReadableIssue(); verr != "" {
		return verr
	}

	return ""
}

// Unknown section

type UnknownSection struct{}

func (u UnknownSection) ValidateReturningUserReadableIssue() string {
	return ""
}

// Title Section

type TitlePageSection struct {
	Title       string
	ScaleFactor float64
	Bold        bool
}

func unpackTitleSection(data map[string]interface{}, s *PageSection) (pageSectionTypeInterface, error) {
	title, ok := data["title"].(string)
	if !ok || title == "" {
		return nil, NewUserPresentableError("Page section of type title must have a title string.")
	}

	scaleFactor, ok := data["scaleFactor"].(float64)
	if !ok || scaleFactor <= 0 {
		scaleFactor = 1.0
	}

	bold, ok := data["bold"].(bool)
	if !ok {
		bold = true
	}

	td := TitlePageSection{
		Title:       title,
		ScaleFactor: scaleFactor,
		Bold:        bold,
	}
	s.TitleData = &td

	return td, nil
}

func (t TitlePageSection) ValidateReturningUserReadableIssue() string {
	if t.Title == "" {
		return "Page section of type title must have a title string."
	}
	if t.ScaleFactor <= 0 {
		return "Page section of type title must have a postive scaleFactor"
	}

	return ""
}

// Body Section

type BodyPageSection struct {
	BodyText            string
	ScaleFactor         float64
	Bold                bool
	UsePrimaryTextColor bool
	CenterText          bool
}

func unpackBodySection(data map[string]interface{}, s *PageSection) (pageSectionTypeInterface, error) {
	bodyText, ok := data["bodyText"].(string)
	if !ok || bodyText == "" {
		return nil, NewUserPresentableError("Page section of type body must have a bodyText string.")
	}

	scaleFactor, ok := data["scaleFactor"].(float64)
	if !ok || scaleFactor <= 0 {
		scaleFactor = 1.0
	}

	// default to false
	bold, _ := data["bold"].(bool)
	usePrimaryFontColor, _ := data["usePrimaryFontColor"].(bool)

	// default true
	centerText, ok := data["centerText"].(bool)
	if !ok {
		centerText = true
	}

	bd := BodyPageSection{
		BodyText:            bodyText,
		ScaleFactor:         scaleFactor,
		Bold:                bold,
		UsePrimaryTextColor: usePrimaryFontColor,
		CenterText:          centerText,
	}
	s.BodyData = &bd

	return bd, nil
}

func (t BodyPageSection) ValidateReturningUserReadableIssue() string {
	if t.BodyText == "" {
		return "Page section of type body must have a bodyText."
	}
	if t.ScaleFactor <= 0 {
		return "Page section of type body must have a postive scaleFactor"
	}

	return ""
}
