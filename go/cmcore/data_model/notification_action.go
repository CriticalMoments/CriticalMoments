package datamodel

import (
	"encoding/json"
)

type NotificationAction struct {
	Title      string
	Body       string
	ActionName string
}

type jsonNotificationAction struct {
	Title      string `json:"title,omitempty"`
	Body       string `json:"body,omitempty"`
	ActionName string `json:"actionName,omitempty"`
}

func unpackNotificationFromJson(rawJson json.RawMessage, ac *ActionContainer) (ActionTypeInterface, error) {
	var n NotificationAction
	err := json.Unmarshal(rawJson, &n)
	if err != nil {
		return nil, err
	}
	ac.NotificationAction = &n
	return &n, nil
}

func (a *NotificationAction) Validate() bool {
	return a.ValidateReturningUserReadableIssue() == ""
}

func (a *NotificationAction) ValidateReturningUserReadableIssue() string {
	if a.Title == "" {
		return "Notifications must have a title."
	}
	return ""
}

func (a *NotificationAction) UnmarshalJSON(data []byte) error {
	var jn jsonNotificationAction
	err := json.Unmarshal(data, &jn)
	if err != nil {
		return NewUserPresentableErrorWSource("Unable to parse the json of an notification with type=notification. Check the format, variable names, and types (eg float vs int).", err)
	}

	a.Title = jn.Title
	a.Body = jn.Body
	a.ActionName = jn.ActionName

	if validationIssue := a.ValidateReturningUserReadableIssue(); validationIssue != "" {
		return NewUserPresentableError(validationIssue)
	}

	return nil
}

func (a *NotificationAction) AllEmbeddedThemeNames() ([]string, error) {
	return []string{}, nil
}

func (a *NotificationAction) AllEmbeddedActionNames() ([]string, error) {
	if a.ActionName != "" {
		return []string{a.ActionName}, nil
	}

	return []string{}, nil
}

func (l *NotificationAction) AllEmbeddedConditions() ([]*Condition, error) {
	return []*Condition{}, nil
}

func (a *NotificationAction) PerformAction(ab ActionBindings, actionName string) error {
	// TODO
	return nil
}
