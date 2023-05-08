package datamodel

import "encoding/json"

const BannerMaxLineCountSystemDefault = -1
const BannerMaxLineCountSystemUnlimited = 0

type BannerAction struct {
	Body              string
	ShowDismissButton bool
	MaxLineCount      int
	TapActionName     string
}

type jsonBannerAction struct {
	Body              string `json:"body"`
	ShowDismissButton *bool  `json:"showDismissButton,omitempty"`
	MaxLineCount      *int   `json:"maxLineCount,omitempty"`
	TapActionName     string `json:"tapActionName,omitempty"`
}

func (ba BannerAction) Validate() bool {
	return ba.ValidateReturningUserReadableIssue() == ""
}

func (b BannerAction) ValidateReturningUserReadableIssue() string {
	if b.Body == "" {
		return "Banners must have body text"
	}
	if b.MaxLineCount != BannerMaxLineCountSystemDefault && b.MaxLineCount < 0 {
		// Technically -1 allowed, but that's an internal between cmcore and libraries
		// Not user facing or a value they should put in json or see in libraries
		return "Banner max line count must be a positive integer, or 0 for no limit"
	}

	return ""
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

	if validationIssue := banner.ValidateReturningUserReadableIssue(); validationIssue != "" {
		return NewUserPresentableError(validationIssue)
	}

	return nil
}
