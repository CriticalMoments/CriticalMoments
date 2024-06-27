package datamodel

import (
	"encoding/json"
)

type Notification struct {
	ID         string
	Title      string
	Body       string
	ActionName string
}

type jsonNotification struct {
	Title      string `json:"title,omitempty"`
	Body       string `json:"body,omitempty"`
	ActionName string `json:"actionName,omitempty"`
}

func (a *Notification) Validate() bool {
	return a.ValidateReturningUserReadableIssue() == ""
}

func (n *Notification) ValidateReturningUserReadableIssue() string {
	return n.ValidateReturningUserReadableIssueIgnoreID(false)
}

func (n *Notification) ValidateReturningUserReadableIssueIgnoreID(ignoreID bool) string {
	if !ignoreID && n.ID == "" {
		return "Notification must have ID"
	}
	if n.Title == "" {
		return "Notifications must have a title."
	}
	return ""
}

func (a *Notification) UnmarshalJSON(data []byte) error {
	var jn jsonNotification
	err := json.Unmarshal(data, &jn)
	if err != nil {
		return NewUserPresentableErrorWSource("Unable to parse the json of an notification with type=notification. Check the format, variable names, and types (eg float vs int).", err)
	}

	a.Title = jn.Title
	a.Body = jn.Body
	a.ActionName = jn.ActionName

	// ignore ID which is set later from primary config
	if validationIssue := a.ValidateReturningUserReadableIssueIgnoreID(true); validationIssue != "" {
		return NewUserPresentableError(validationIssue)
	}

	return nil
}
