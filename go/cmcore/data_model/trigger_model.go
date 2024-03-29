package datamodel

import (
	"encoding/json"
	"fmt"
)

type Trigger struct {
	EventName  string
	ActionName string
	Condition  *Condition
}

type jsonTrigger struct {
	EventName  string     `json:"eventName"`
	ActionName string     `json:"actionName"`
	Condition  *Condition `json:"condition"`
}

func (t *Trigger) UnmarshalJSON(data []byte) error {
	var jt jsonTrigger
	err := json.Unmarshal(data, &jt)
	if err != nil {
		return NewUserPresentableErrorWSource("Unable to parse the json of a trigger. Check the format, variable names, and types (eg float vs int).", err)
	}

	t.ActionName = jt.ActionName
	t.EventName = jt.EventName
	t.Condition = jt.Condition

	if validationIssue := t.ValidateReturningUserReadableIssue(); validationIssue != "" {
		return NewUserPresentableError(validationIssue)
	}

	return nil
}

func (t *Trigger) Validate() bool {
	return t.ValidateReturningUserReadableIssue() == ""
}

func (t *Trigger) ValidateReturningUserReadableIssue() string {
	if t.EventName == "" {
		return "All triggers require an event"
	}
	if t.ActionName == "" {
		return "All triggers require an action name"
	}
	if t.Condition != nil {
		if err := t.Condition.Validate(); err != nil {
			return fmt.Sprintf("Condition in trigger is not valid: [[%v]]", t.Condition.conditionString)
		}
	}
	return ""
}
