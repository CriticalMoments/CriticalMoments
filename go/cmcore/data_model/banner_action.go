package datamodel

import (
	"encoding/json"
	"fmt"
)

const BannerMaxLineCountSystemDefault = -1
const BannerMaxLineCountSystemUnlimited = 0

const BannerPositionTop = "top"
const BannerPositionBottom = "bottom"

type BannerAction struct {
	Body              string
	ShowDismissButton bool
	MaxLineCount      int
	TapActionName     string
	CustomThemeName   string
	PreferredPosition string
}

type jsonBannerAction struct {
	Body              string  `json:"body"`
	ShowDismissButton *bool   `json:"showDismissButton,omitempty"`
	MaxLineCount      *int    `json:"maxLineCount,omitempty"`
	TapActionName     string  `json:"tapActionName,omitempty"`
	CustomThemeName   string  `json:"themeName,omitempty"`
	PreferredPosition *string `json:"preferredPosition,omitempty"`
}

func unpackBannerFromJson(rawJson json.RawMessage, ac *ActionContainer) (ActionTypeInterface, error) {
	var banner BannerAction
	err := json.Unmarshal(rawJson, &banner)
	if err != nil {
		return nil, err
	}
	ac.BannerAction = &banner
	return &banner, nil
}

func (ba *BannerAction) Validate() bool {
	return ba.ValidateReturningUserReadableIssue() == ""
}

func (b *BannerAction) ValidateReturningUserReadableIssue() string {
	if b.Body == "" {
		return "Banners must have body text"
	}
	if b.MaxLineCount != BannerMaxLineCountSystemDefault && b.MaxLineCount < 0 {
		// Technically -1 allowed, but that's an internal between cmcore and libraries
		// Not user facing or a value they should put in json or see in libraries
		return "Banner max line count must be a positive integer, or 0 for no limit"
	}
	if b.PreferredPosition != "" && b.PreferredPosition != BannerPositionTop && b.PreferredPosition != BannerPositionBottom {
		return fmt.Sprintf("Banner preferred position must be empty, top or bottom. \"%v\" is not valid", b.PreferredPosition)
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
	preferredPosition := ""
	if ja.PreferredPosition != nil {
		preferredPosition = *ja.PreferredPosition
	}

	banner.Body = ja.Body
	banner.ShowDismissButton = showDismissButton
	banner.MaxLineCount = maxLineCount
	banner.TapActionName = ja.TapActionName
	banner.CustomThemeName = ja.CustomThemeName
	banner.PreferredPosition = preferredPosition

	if validationIssue := banner.ValidateReturningUserReadableIssue(); validationIssue != "" {
		return NewUserPresentableError(validationIssue)
	}

	return nil
}

func (b *BannerAction) AllEmbeddedThemeNames() ([]string, error) {
	if b.CustomThemeName == "" {
		return []string{}, nil
	}
	return []string{b.CustomThemeName}, nil
}

func (b *BannerAction) AllEmbeddedActionNames() ([]string, error) {
	if b.TapActionName == "" {
		return []string{}, nil
	}
	return []string{b.TapActionName}, nil
}

func (l *BannerAction) AllEmbeddedConditions() ([]*Condition, error) {
	return []*Condition{}, nil
}

func (b *BannerAction) PerformAction(ab ActionBindings) error {
	return ab.ShowBanner(b)
}
