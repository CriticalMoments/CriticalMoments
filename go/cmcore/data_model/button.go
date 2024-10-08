package datamodel

import (
	"encoding/json"
	"fmt"

	"golang.org/x/exp/slices"
)

const (
	ButtonStyleEnumLarge     string = "large"
	ButtonStyleEnumNormal    string = "normal"
	ButtonStyleEnumSecondary string = "secondary"
	ButtonStyleEnumTertiary  string = "tertiary"
	ButtonStyleEnumInfo      string = "info"
	ButtonStyleEnumInfoSmall string = "info-small"
)

var buttonStyles = []string{
	ButtonStyleEnumLarge,
	ButtonStyleEnumNormal,
	ButtonStyleEnumSecondary,
	ButtonStyleEnumTertiary,
	ButtonStyleEnumInfo,
	ButtonStyleEnumInfoSmall,
}

type Button struct {
	Title          string
	Style          string
	ActionName     string
	PreventDefault bool
}

type jsonButton struct {
	Title          string `json:"title"`
	Style          string `json:"style,omitempty"`
	ActionName     string `json:"actionName,omitempty"`
	PreventDefault bool   `json:"preventDefault,omitempty"`
}

func (b *Button) UnmarshalJSON(data []byte) error {
	var jb jsonButton
	err := json.Unmarshal(data, &jb)
	if err != nil {
		return NewUserPresentableErrorWSource("Unable to parse the json of a modal button.", err)
	}

	b.Title = jb.Title
	b.ActionName = jb.ActionName
	b.PreventDefault = jb.PreventDefault

	// Style: default to normal if empty or not strict validation
	b.Style = jb.Style
	if b.Style == "" {
		b.Style = ButtonStyleEnumNormal
	}
	if !slices.Contains(buttonStyles, b.Style) {
		if StrictDatamodelParsing {
			return NewUserErrorForJsonIssue(data, NewUserPresentableError(fmt.Sprintf("invalid button 'style' tag found in config (for a modal button): \"%v\"", b.Style)))
		} else {
			// Backwards compatibility: fallback to normal if this client doesn't recognize the style
			b.Style = ButtonStyleEnumNormal
		}
	}

	if err := b.Check(); err != nil {
		return NewUserErrorForJsonIssue(data, err)
	}

	return nil
}

func (b *Button) Check() UserPresentableErrorInterface {
	if b.Title == "" {
		return NewUserPresentableError("Modal button 'title' tag can not be empty in config.")
	}

	if !slices.Contains(buttonStyles, b.Style) {
		return NewUserPresentableError(fmt.Sprintf("Modal button had invalid button 'style' tag: \"%v\"", b.Style))
	}

	return nil
}
