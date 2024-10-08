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
		typeErr := "Page 'pageType' tag must be 'stack'"
		if StrictDatamodelParsing {
			return NewUserErrorForJsonIssue(data, NewUserPresentableError(typeErr))
		} else {
			fmt.Printf("CriticalMoments: %v. Ignoring, but if not expected check your config file.\n", typeErr)
		}
	}

	if err := p.Check(); err != nil {
		return NewUserErrorForJsonIssue(data, err)
	}

	return nil
}

func (p *Page) Check() UserPresentableErrorInterface {
	if len(p.Sections) == 0 {
		// back-compat: allow zero sections when not strict
		if StrictDatamodelParsing {
			return NewUserPresentableError("page with 0 sections is not valid")
		}
	}

	for i, section := range p.Sections {
		if valErr := section.Check(); valErr != nil {
			return NewUserPresentableErrorWSource(fmt.Sprintf("Page has an invalid section at index [%v]", i), valErr)
		}
	}

	for i, button := range p.Buttons {
		if valErr := button.Check(); valErr != nil {
			return NewUserPresentableErrorWSource(fmt.Sprintf("Page has an invalid button at index [%v]", i), valErr)
		}
	}

	return nil
}

const (
	SectionTypeEnumTitle    string = "title"
	SectionTypeEnumBodyText string = "body"
	SectionTypeEnumImage    string = "image"
)

type pageSectionTypeInterface interface {
	Check() UserPresentableErrorInterface
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
		if StrictDatamodelParsing {
			return NewUserErrorForJsonIssue(data, NewUserPresentableError(fmt.Sprintf("CriticalMoments: Unsupported page section 'type' tag: \"%v\" found in config file", s.PageSectionType)))
		} else {
			// back-compat -- fallback to unknown section type
			s.pageSectionData = UnknownSection{}
		}
	} else {
		pageSectionData, err := unpacker(js.RawSectionData, s)
		if err != nil {
			return err
		}
		s.pageSectionData = pageSectionData
	}

	if err := s.Check(); err != nil {
		return NewUserErrorForJsonIssue(data, err)
	}

	return nil
}

func (s *PageSection) Check() UserPresentableErrorInterface {
	if s.TopSpacingScale < 0 {
		return NewUserPresentableError("Top space scale for a page section must be >= 0")
	}

	if s.pageSectionData == nil {
		return NewUserPresentableError("Invalid section in page -- nil section data")
	}
	if verr := s.pageSectionData.Check(); verr != nil {
		return verr
	}

	if StrictDatamodelParsing && !slices.Contains(maps.Keys(pageSectionTypeRegistry), s.PageSectionType) {
		return NewUserPresentableError(fmt.Sprintf("Page section with unknown 'type' tag: %v", s.PageSectionType))
	}

	return nil
}

// Unknown section

type UnknownSection struct{}

func (u UnknownSection) Check() UserPresentableErrorInterface {
	return nil
}

// Title Section

type TitlePageSection struct {
	Title               string
	ScaleFactor         float64
	Bold                bool
	UsePrimaryTextColor bool
	CenterText          bool
	Width               float64
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

	// default to 0
	widthFloat, _ := data["width"].(float64)

	td := TitlePageSection{
		Title:               title,
		ScaleFactor:         scaleFactor,
		Bold:                bold,
		CenterText:          centerText,
		UsePrimaryTextColor: usePrimaryFontColor,
		Width:               widthFloat,
	}
	s.TitleData = &td

	return td, nil
}

func (t TitlePageSection) Check() UserPresentableErrorInterface {
	if t.Title == "" {
		return NewUserPresentableError("Page section of type title must have a title string.")
	}
	if t.ScaleFactor <= 0 {
		return NewUserPresentableError("Page section of type title must have a positive scaleFactor")
	}

	return nil
}

// Body Section

type BodyPageSection struct {
	BodyText            string
	ScaleFactor         float64
	Bold                bool
	UsePrimaryTextColor bool
	CenterText          bool
	Width               float64
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

	// default to 0
	widthFloat, _ := data["width"].(float64)

	bd := BodyPageSection{
		BodyText:            bodyText,
		ScaleFactor:         scaleFactor,
		Bold:                bold,
		UsePrimaryTextColor: usePrimaryFontColor,
		CenterText:          centerText,
		Width:               widthFloat,
	}
	s.BodyData = &bd

	return bd, nil
}

func (t BodyPageSection) Check() UserPresentableErrorInterface {
	if t.BodyText == "" {
		return NewUserPresentableError("Page section of type body must have a bodyText.")
	}
	if t.ScaleFactor <= 0 {
		return NewUserPresentableError("Page section of type body must have a positive scaleFactor")
	}

	return nil
}

// Image Section

func unpackImageSection(rawData json.RawMessage, s *PageSection) (pageSectionTypeInterface, error) {
	var i Image
	err := json.Unmarshal(rawData, &i)
	if err != nil {
		return nil, NewUserPresentableErrorWSource("Unable to parse the json of a page section (image).", err)
	}

	s.ImageData = &i

	return &i, nil
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
