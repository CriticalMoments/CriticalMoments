package datamodel

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/CriticalMoments/CriticalMoments/go/cmcore/conditions"
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

	Condition string

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
}

type jsonActionContainer struct {
	ActionType    string          `json:"actionType"`
	Condition     string          `json:"condition"`
	RawActionData json.RawMessage `json:"actionData"`
}

// To be implemented by client libaray (eg: iOS SDK or Appcore)
type ActionBindings interface {
	// Actions
	ShowBanner(banner *BannerAction) error
	ShowAlert(alert *AlertAction) error
	ShowLink(link *LinkAction) error
	PerformConditionalAction(conditionalAction *ConditionalAction) error
	ShowReviewPrompt() error
	ShowModal(modal *ModalAction) error
}

type ActionTypeInterface interface {
	AllEmbeddedThemeNames() ([]string, error)
	AllEmbeddedActionNames() ([]string, error)
	ValidateReturningUserReadableIssue() string
	PerformAction(ActionBindings) error
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
		typeErr := fmt.Sprintf("Unsupported action type: \"%v\" found in config file.", jac.ActionType)
		if StrictDatamodelParsing {
			return NewUserPresentableError(typeErr)
		} else {
			fmt.Printf("CriticalMoments: %v. Will proceed, but this action will be a no-op. If unexpected, check the CM config file.\n", typeErr)
			actionData = &UnknownAction{ActionType: jac.ActionType}
		}
	}

	if jac.Condition != "" {
		if err = conditions.ValidateCondition(jac.Condition); err != nil {
			return NewUserPresentableErrorWSource(fmt.Sprintf("Invalid condition: [[ %v ]]", jac.Condition), err)
		}
	}

	ac.actionData = actionData
	ac.ActionType = jac.ActionType
	ac.Condition = jac.Condition
	return nil
}

func (ac *ActionContainer) ValidateReturningUserReadableIssue() string {
	if ac.ActionType == "" {
		return "Empty actionType not permitted"
	}
	// Check the type hasn't been changed to something unsupported.
	_, ok := actionTypeRegistry[ac.ActionType]
	if !ok {
		_, ok := ac.actionData.(*UnknownAction)
		if !ok {
			return "Internal error. Code 776232923."
		}
	}

	if ac.actionData == nil {
		// the action type data interface should be set after unmarshaling.
		// This is a code issue if it occurs, not a data issue
		return fmt.Sprintf("Action type %v has internal issues", ac.ActionType)
	}

	return ac.actionData.ValidateReturningUserReadableIssue()
}

func (ac *ActionContainer) PerformAction(ab ActionBindings) error {
	if ac.actionData == nil {
		return errors.New("Attempted to perform action without AD interface")
	}
	return ac.actionData.PerformAction(ab)
}
