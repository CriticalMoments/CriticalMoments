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
		return NewUserPresentableErrorWSource("Unable to parse the json of a button.", err)
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
		errString := fmt.Sprintf("Invalid button style: \"%v\"", b.Style)
		if StrictDatamodelParsing {
			return NewUserPresentableError(errString)
		} else {
			fmt.Printf("CriticalMoments: %v. Will fallback to normal style.\n", errString)
			b.Style = ButtonStyleEnumNormal
		}
	}

	if validationIssue := b.ValidateReturningUserReadableIssue(); validationIssue != "" {
		return NewUserPresentableError(validationIssue)
	}

	return nil
}

func (b *Button) ValidateReturningUserReadableIssue() string {
	if b.Title == "" {
		return "Button title can not be empty."
	}

	if !slices.Contains(buttonStyles, b.Style) {
		return fmt.Sprintf("Invalid button style: \"%v\"", b.Style)
	}

	return ""
}
