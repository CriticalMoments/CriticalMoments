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
	ActionTypeEnumBanner      string = "banner"
	ActionTypeEnumAlert       string = "alert"
	ActionTypeEnumLink        string = "link"
	ActionTypeEnumConditional string = "conditional_action"
	ActionTypeEnumModal       string = "modal"
	ActionTypeEnumReview      string = "review_prompt"
)

// This section is the json data model we use for parsing/masrshaling

type ActionContainer struct {
	ActionType string

	Condition *Condition

	// Strongly typed action data
	// All nil except the one aligning to actionType
	BannerAction      *BannerAction
	AlertAction       *AlertAction
	LinkAction        *LinkAction
	ConditionalAction *ConditionalAction
	ModalAction       *ModalAction

	// generalized interface for functions we need for any actions type.
	// Typically a pointer to the one value above that is populated.
	actionData ActionTypeInterface

	// Fallback if any issues performing action
	FallbackActionName string
}

type jsonActionContainer struct {
	ActionType         string          `json:"actionType"`
	Condition          *Condition      `json:"condition"`
	FallbackActionName string          `json:"fallback"`
	RawActionData      json.RawMessage `json:"actionData"`
}

// To be implemented by client libaray (eg: iOS SDK or Appcore)
type ActionBindings interface {
	// Actions
	ShowBanner(banner *BannerAction, actionName string) error
	ShowAlert(alert *AlertAction, actionName string) error
	ShowLink(link *LinkAction) error
	PerformConditionalAction(conditionalAction *ConditionalAction) error
	PerformNamedAction(name string) error
	ShowReviewPrompt() error
	ShowModal(modal *ModalAction, actionName string) error
}

type ActionTypeInterface interface {
	AllEmbeddedThemeNames() ([]string, error)
	AllEmbeddedActionNames() ([]string, error)
	AllEmbeddedConditions() ([]*Condition, error)
	Check() UserPresentableErrorInterface
	PerformAction(ab ActionBindings, actionName string) error
}

var (
	actionTypeRegistry = map[string]func(json.RawMessage, *ActionContainer) (ActionTypeInterface, error){
		ActionTypeEnumBanner:      unpackBannerFromJson,
		ActionTypeEnumAlert:       unpackAlertFromJson,
		ActionTypeEnumLink:        unpackLinkFromJson,
		ActionTypeEnumConditional: unpackConditionalActionFromJson,
		ActionTypeEnumReview:      unpackReviewFromJson,
		ActionTypeEnumModal:       unpackModalFromJson,
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
	var actionData ActionTypeInterface
	if ok && unpacker != nil {
		actionData, err = unpacker(jac.RawActionData, ac)
		if err != nil {
			return NewUserPresentableErrorWSource(fmt.Sprintf("Issue unpacking type \"%v\"", jac.ActionType), err)
		}
	} else {
		// Allow backwards compatibility, defaulting to no-op
		if StrictDatamodelParsing {
			typeErr := fmt.Sprintf("unsupported action type found in config file: \"%v\"", jac.ActionType)
			return NewUserPresentableError(typeErr)
		} else {
			actionData = &UnknownAction{ActionType: jac.ActionType}
		}
	}

	ac.actionData = actionData
	ac.ActionType = jac.ActionType
	ac.Condition = jac.Condition
	ac.FallbackActionName = jac.FallbackActionName
	return nil
}

func (ac *ActionContainer) Check() UserPresentableErrorInterface {
	if ac.ActionType == "" {
		return NewUserPresentableError("Empty actionType not permitted")
	}
	// Check the type hasn't been changed to something unsupported.
	_, ok := actionTypeRegistry[ac.ActionType]
	if !ok {
		_, ok := ac.actionData.(*UnknownAction)
		if !ok {
			return NewUserPresentableError("Internal error. Code 776232923.")
		}
	}

	if ac.actionData == nil {
		// the action type data interface should be set after unmarshaling.
		// This is a code issue if it occurs, not a data issue
		return NewUserPresentableError(fmt.Sprintf("Action type %v has internal issues", ac.ActionType))
	}

	return ac.actionData.Check()
}

func (ac *ActionContainer) PerformAction(ab ActionBindings, actionName string) error {
	if ac.actionData == nil {
		return errors.New("attempted to perform action without AD interface")
	}
	err := ac.actionData.PerformAction(ab, actionName)

	if err != nil && ac.FallbackActionName != "" {
		err = ab.PerformNamedAction(ac.FallbackActionName)
	}

	return err
}
