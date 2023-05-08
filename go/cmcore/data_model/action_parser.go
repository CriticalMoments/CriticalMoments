package datamodel

import (
	"encoding/json"
	"errors"
	"fmt"
)

/*
System here for each new type
 - enum value of all types
 - pure data models for each action like "BannerAction" (in their own file, lots of valdation)
 - json representations of each action like jsonBannerAction, and parsers like NewBannerActionFromJson in this file
 - jsonActionContainer includes all possible action type pointers, and a raw "ActionData" blob
 - jsonActionContainer.UnmarshalJson first populates type, then populates one of the pure data models using parsers mentioned above
 - All errors should be user readable since this deals with user input
*/

const (
	ActionTypeEnumBanner string = "banner"
)

// This section is the json data model we use for parsing/masrshaling

type ActionContainer struct {
	ActionType string

	// All nil except the one aligning to actionType
	BannerAction *BannerAction
}

type jsonActionContainer struct {
	ActionType    string          `json:"actionType"`
	RawActionData json.RawMessage `json:"actionData"`
}

type jsonBannerAction struct {
	Body              string `json:"body"`
	ShowDismissButton *bool  `json:"showDismissButton,omitempty"`
	MaxLineCount      *int   `json:"maxLineCount,omitempty"`
	TapActionName     string `json:"tapActionName,omitempty"`
	Theme             string `json:"theme,omitempty"`
}

// TODO: this is right system. Data model should support UnmarshalJSON for other types too...
func (ac *ActionContainer) UnmarshalJSON(data []byte) error {
	// docs suggest no-op for empty data
	if data == nil {
		return nil
	}

	var jac jsonActionContainer
	err := json.Unmarshal(data, &jac)
	if err != nil {
		return err
	}

	switch jac.ActionType {
	case ActionTypeEnumBanner:
		var banner BannerAction
		err = json.Unmarshal(jac.RawActionData, &banner)
		if err != nil {
			return err
		}
		ac.BannerAction = &banner
		ac.ActionType = ActionTypeEnumBanner
	default:
		return NewUserPresentableError(fmt.Sprintf("Unsupported action type: \"%v\"", jac.ActionType))
	}

	return nil
}

func (ac *ActionContainer) AllEmbeddedActionNames() ([]string, error) {
	if ac.ActionType == "" {
		return nil, errors.New("AllEmbeddedActionNames called on an uninitialized action continer")
	}

	switch ac.ActionType {
	case ActionTypeEnumBanner:
		// TODO Test case
		if ac.BannerAction.TapActionName == "" {
			return []string{}, nil
		}
		return []string{ac.BannerAction.TapActionName}, nil
	default:
		return nil, NewUserPresentableError(fmt.Sprintf("Unsupported action type: \"%v\"", ac.ActionType))
	}
}

func (banner *BannerAction) UnmarshalJSON(data []byte) error {
	var ja jsonBannerAction
	err := json.Unmarshal(data, &ja)
	if err != nil {
		return NewUserPresentableErrorWSource("Unable to parse the json of an action with type=banner. Check the format, variable names, and types (eg float vs int).", err)
	}

	// Default Values for nullable options
	showDismissButton := true
	if ja.ShowDismissButton != nil {
		showDismissButton = *ja.ShowDismissButton
	}
	maxLineCount := BannerMaxLineCountSystemDefault // go requires a value
	if ja.MaxLineCount != nil {
		maxLineCount = *ja.MaxLineCount
	}

	banner.Body = ja.Body
	banner.ShowDismissButton = showDismissButton
	banner.MaxLineCount = maxLineCount
	banner.TapActionName = ja.TapActionName
	banner.Theme = ja.Theme

	if validationIssue := banner.ValidateReturningUserReadableIssue(); validationIssue != "" {
		return NewUserPresentableError(validationIssue)
	}

	return nil
}
