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
	ActionTypeEnumAlert  string = "alert"
)

// This section is the json data model we use for parsing/masrshaling

type ActionContainer struct {
	ActionType string

	// Strongly typed action data
	// All nil except the one aligning to actionType
	BannerAction *BannerAction
	AlertAction  *AlertAction

	// generalized interface all of above, for functions we need for all actions types
	// Typically a pointer to the one value above that is populated
	actionData ActionTypeInterface
}

type jsonActionContainer struct {
	ActionType    string          `json:"actionType"`
	RawActionData json.RawMessage `json:"actionData"`
}

// To be implemented by client libaray (eg: iOS SDK)
type ActionBindings interface {
	// Actions
	ShowBanner(banner *BannerAction) error
	ShowAlert(alert *AlertAction) error
}

type ActionTypeInterface interface {
	AllEmbeddedThemeNames() ([]string, error)
	AllEmbeddedActionNames() ([]string, error)
	ValidateReturningUserReadableIssue() string
	PerformAction(ActionBindings) error
}

var (
	actionTypeRegistry = map[string]func(json.RawMessage, *ActionContainer) (ActionTypeInterface, error){
		ActionTypeEnumBanner: unpackBannerFromJson,
		ActionTypeEnumAlert:  unpackAlertFromJson,
	}
)

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

	unpacker, ok := actionTypeRegistry[jac.ActionType]
	if !ok || unpacker == nil {
		return NewUserPresentableError(fmt.Sprintf("Unsupported action type: \"%v\"", jac.ActionType))
	}

	actionData, err := unpacker(jac.RawActionData, ac)
	if err != nil {
		return NewUserPresentableErrorWSource(fmt.Sprintf("Issue unpacking type \"%v\"", jac.ActionType), err)
	}
	ac.actionData = actionData
	ac.ActionType = jac.ActionType
	return nil
}

func (ac *ActionContainer) ValidateReturningUserReadableIssue() string {
	if ac.ActionType == "" {
		return "Empty actionType"
	}
	// Check the type hasn't been changed to something unsupported
	_, ok := actionTypeRegistry[ac.ActionType]
	if !ok {
		return "Internal error. Code 776232923."
	}

	// TODO check it's registered in system

	if ac.actionData == nil {
		// the action type data interface should be set after unmarshaling.
		// This is a code issue if it occurs, not a data issue
		return fmt.Sprintf("Action type %v has internal issues", ac.ActionType)
	}

	// TODO: losing the validation that one of the strong pointers is populated. Tests enough?

	return ac.actionData.ValidateReturningUserReadableIssue()
}

func (ac *ActionContainer) PerformAction(ab ActionBindings) error {
	if ac.actionData == nil {
		return errors.New("Attempted to perform action without AD interface")
	}
	return ac.actionData.PerformAction(ab)
}
