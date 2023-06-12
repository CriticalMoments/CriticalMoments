package datamodel

import (
	"encoding/json"
)

type ReviewAction struct {
}

func unpackReviewFromJson(rawJson json.RawMessage, ac *ActionContainer) (ActionTypeInterface, error) {
	return &ReviewAction{}, nil
}

func (l *ReviewAction) ValidateReturningUserReadableIssue() string {
	return ""
}

func (r *ReviewAction) AllEmbeddedThemeNames() ([]string, error) {
	return []string{}, nil
}

func (r *ReviewAction) AllEmbeddedActionNames() ([]string, error) {
	return []string{}, nil
}

func (r *ReviewAction) PerformAction(ab ActionBindings) error {
	return ab.ShowReviewPrompt()
}
