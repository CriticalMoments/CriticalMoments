package datamodel

import (
	"encoding/json"
	"fmt"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// Model is specific to "Stack" page. Can add type/data here if we add more types, but unnecessary
// while there is only one.
type Page struct {
	Sections []*PageSection
	Buttons  []*Button
}

const PageTypeEnumStack string = "stack"

// Future proofing. Expect more page types in future so ensure the datamodel can be typed.
type jsonPage struct {
	PageType string         `json:"pageType"`
	PageData *jsonStackPage `json:"pageData"`
}

type jsonStackPage struct {
	Sections []*PageSection `json:"sections"`
	Buttons  []*Button      `json:"buttons"`
}

func (p *Page) UnmarshalJSON(data []byte) error {
	var jp jsonPage
	err := json.Unmarshal(data, &jp)
	if err != nil {
		return NewUserPresentableErrorWSource("Unable to parse the json of a page.", err)
	}

	if jp.PageType == PageTypeEnumStack {
		p.Sections = jp.PageData.Sections
		p.Buttons = jp.PageData.Buttons
	} else {
		typeErr := "pageType must be 'stack'"
		if StrictDatamodelParsing {
			return NewUserPresentableError(typeErr)
		} else {
			fmt.Printf("CriticalMoments: %v. Ignoring, but if not expected check your config file.\n", typeErr)
		}
	}

	if validationIssue := p.ValidateReturningUserReadableIssue(); validationIssue != "" {
		return NewUserPresentableError(validationIssue)
	}

	return nil
}

func (p *Page) ValidateReturningUserReadableIssue() string {
	if len(p.Sections) == 0 {
		if StrictDatamodelParsing {
			return "Page with 0 sections is not valid"
		} else {
			fmt.Printf("CriticalMoments: page with 0 sections not valid. Ignoring but if unexpected check your config file.")
		}
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
	SectionTypeEnumImage    string = "image"
)

type pageSectionTypeInterface interface {
	ValidateReturningUserReadableIssue() string
}

var (
	pageSectionTypeRegistry = map[string]func(json.RawMessage, *PageSection) (pageSectionTypeInterface, error){
		SectionTypeEnumTitle:    unpackTitleSection,
		SectionTypeEnumBodyText: unpackBodySection,
		SectionTypeEnumImage:    unpackImageSection,
	}
)

type PageSection struct {
	PageSectionType string
	TopSpacingScale float64

	// Section types, stronly typed for easy consumption
	TitleData *TitlePageSection
	BodyData  *BodyPageSection
	ImageData *Image

	// generalized interface for functions we need for any section type.
	// Typically a pointer to the one value above that is populated.
	pageSectionData pageSectionTypeInterface
}

type jsonPageSection struct {
	PageSectionType string          `json:"pageSectionType"`
	TopSpacingScale *float64        `json:"topSpacingScale,omitempty"`
	RawSectionData  json.RawMessage `json:"pageSectionData,omitempty"`
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

	if StrictDatamodelParsing && !slices.Contains(maps.Keys(pageSectionTypeRegistry), s.PageSectionType) {
		return fmt.Sprintf("Page section with unknown type: %v", s.PageSectionType)
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
	Title               string
	ScaleFactor         float64
	Bold                bool
	UsePrimaryTextColor bool
	CenterText          bool
}

func unpackTitleSection(rawData json.RawMessage, s *PageSection) (pageSectionTypeInterface, error) {
	var data map[string]interface{}
	err := json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, NewUserPresentableErrorWSource("Unable to parse the json of a page section (title).", err)
	}

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
	centerText, ok := data["centerText"].(bool)
	if !ok {
		centerText = true
	}
	usePrimaryFontColor, ok := data["usePrimaryFontColor"].(bool)
	if !ok {
		usePrimaryFontColor = true
	}

	td := TitlePageSection{
		Title:               title,
		ScaleFactor:         scaleFactor,
		Bold:                bold,
		CenterText:          centerText,
		UsePrimaryTextColor: usePrimaryFontColor,
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

func unpackBodySection(rawData json.RawMessage, s *PageSection) (pageSectionTypeInterface, error) {
	var data map[string]interface{}
	err := json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, NewUserPresentableErrorWSource("Unable to parse the json of a page section (body).", err)
	}

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

// Image Section

func unpackImageSection(rawData json.RawMessage, s *PageSection) (pageSectionTypeInterface, error) {
	var i Image
	err := json.Unmarshal(rawData, &i)
	if err != nil {
		return nil, NewUserPresentableErrorWSource("Unable to parse the json of a page section (image).", err)
	}

	s.ImageData = &i

	return i, nil
}

// Enumerators because go mobile doesn't support arrays...

func (p *Page) ButtonsCount() int {
	return len(p.Buttons)
}

func (p *Page) ButtonAtIndex(i int) *Button {
	return p.Buttons[i]
}

func (p *Page) SectionCount() int {
	return len(p.Sections)
}

func (p *Page) SectionAtIndex(i int) *PageSection {
	return p.Sections[i]
}
