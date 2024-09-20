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

	if err := t.Check(); err != nil {
		return NewUserErrorForJsonIssue(data, err)
	}
	return nil
}

func (t *Trigger) Valid() bool {
	return t.Check() == nil
}

func (t *Trigger) Check() UserPresentableErrorInterface {
	if t.EventName == "" {
		return NewUserPresentableError("Triggers require an eventName")
	}
	if t.ActionName == "" {
		return NewUserPresentableError("Triggers require an actionName")
	}
	if t.Condition != nil {
		if err := t.Condition.Validate(); err != nil {
			return NewUserPresentableErrorWSource(fmt.Sprintf("Condition in trigger is not valid: [[%v]]", t.Condition.conditionString), err)
		}
	}
	return nil
}
