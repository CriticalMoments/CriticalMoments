package datamodel

import (
	"encoding/json"
)

type ReviewAction struct {
}

func unpackReviewFromJson(rawJson json.RawMessage, ac *ActionContainer) (ActionTypeInterface, error) {
	return &ReviewAction{}, nil
}

func (l *ReviewAction) ValidateReturningUserReadableIssue() UserPresentableErrorInterface {
	return nil
}

func (r *ReviewAction) AllEmbeddedThemeNames() ([]string, error) {
	return []string{}, nil
}

func (r *ReviewAction) AllEmbeddedActionNames() ([]string, error) {
	return []string{}, nil
}

func (r *ReviewAction) AllEmbeddedConditions() ([]*Condition, error) {
	return []*Condition{}, nil
}

func (r *ReviewAction) PerformAction(ab ActionBindings, actionName string) error {
	return ab.ShowReviewPrompt()
}
