package datamodel

import (
	"encoding/json"
	"fmt"
)

type ConditionalAction struct {
	Condition        *Condition
	PassedActionName string
	FailedActionName string
}

type jsonConditionalAction struct {
	Condition        *Condition `json:"condition"`
	PassedActionName string     `json:"passedActionName"`
	FailedActionName string     `json:"failedActionName,omitempty"`
}

func unpackConditionalActionFromJson(rawJson json.RawMessage, ac *ActionContainer) (ActionTypeInterface, error) {
	var condition ConditionalAction
	err := json.Unmarshal(rawJson, &condition)
	if err != nil {
		return nil, err
	}
	ac.ConditionalAction = &condition
	return &condition, nil
}

func (c *ConditionalAction) Validate() bool {
	return c.ValidateReturningUserReadableIssue() == ""
}

func (c *ConditionalAction) ValidateReturningUserReadableIssue() string {
	if c.Condition == nil {
		return "Conditional actions must have a condition"
	}
	if err := c.Condition.Validate(); err != nil {
		return fmt.Sprintf("Condition in conditional action is not valid: [[%v]]", c.Condition)
	}
	if c.PassedActionName == "" {
		return "Conditional actions must include a passedActionName to run if condition passes (failedActionName is optional)"
	}

	return ""
}
func (c *ConditionalAction) UnmarshalJSON(data []byte) error {
	var jc jsonConditionalAction
	err := json.Unmarshal(data, &jc)
	if err != nil {
		return NewUserPresentableErrorWSource("Unable to parse the json of an action with type=conditional_action. Check the format, variable names, and types (eg float vs int).", err)
	}

	c.Condition = jc.Condition
	c.PassedActionName = jc.PassedActionName
	c.FailedActionName = jc.FailedActionName

	if validationIssue := c.ValidateReturningUserReadableIssue(); validationIssue != "" {
		return NewUserPresentableError(validationIssue)
	}

	return nil
}

func (c *ConditionalAction) AllEmbeddedThemeNames() ([]string, error) {
	return []string{}, nil
}

func (c *ConditionalAction) AllEmbeddedActionNames() ([]string, error) {
	conditionActions := []string{c.PassedActionName}
	if c.FailedActionName != "" {
		conditionActions = append(conditionActions, c.FailedActionName)
	}
	return conditionActions, nil
}

func (ca *ConditionalAction) AllEmbeddedConditions() ([]*Condition, error) {
	return []*Condition{ca.Condition}, nil
}

func (c *ConditionalAction) PerformAction(ab ActionBindings, actionName string) error {
	return ab.PerformConditionalAction(c)
}
