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

	// All nil except the one aligning to actionType
	BannerAction *BannerAction
	AlertAction  *AlertAction
}

type jsonActionContainer struct {
	ActionType    string          `json:"actionType"`
	RawActionData json.RawMessage `json:"actionData"`
}

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

	switch jac.ActionType {
	case ActionTypeEnumBanner:
		var banner BannerAction
		err = json.Unmarshal(jac.RawActionData, &banner)
		if err != nil {
			return err
		}
		ac.BannerAction = &banner
		ac.ActionType = ActionTypeEnumBanner
	case ActionTypeEnumAlert:
		var alert AlertAction
		err = json.Unmarshal(jac.RawActionData, &alert)
		if err != nil {
			return err
		}
		ac.AlertAction = &alert
		ac.ActionType = ActionTypeEnumAlert
	default:
		return NewUserPresentableError(fmt.Sprintf("Unsupported action type: \"%v\"", jac.ActionType))
	}

	return nil
}

func (ac *ActionContainer) AllEmbeddedThemeNames() ([]string, error) {
	if ac.ActionType == "" {
		return nil, errors.New("AllEmbeddedThemeNames called on an uninitialized action continer")
	}

	switch ac.ActionType {
	case ActionTypeEnumBanner:
		if ac.BannerAction.CustomThemeName == "" {
			return []string{}, nil
		}
		return []string{ac.BannerAction.CustomThemeName}, nil
	case ActionTypeEnumAlert:
		return []string{}, nil
	default:
		return nil, NewUserPresentableError(fmt.Sprintf("Unsupported action type: \"%v\"", ac.ActionType))
	}
}

func (ac *ActionContainer) AllEmbeddedActionNames() ([]string, error) {
	if ac.ActionType == "" {
		return nil, errors.New("AllEmbeddedActionNames called on an uninitialized action continer")
	}

	switch ac.ActionType {
	case ActionTypeEnumBanner:
		if ac.BannerAction.TapActionName == "" {
			return []string{}, nil
		}
		return []string{ac.BannerAction.TapActionName}, nil
	case ActionTypeEnumAlert:
		// TODO: test in alert_action_test once we generalize
		alertActions := []string{}
		if ac.AlertAction.OkButtonActionName != "" {
			alertActions = append(alertActions, ac.AlertAction.OkButtonActionName)
		}
		for _, customButton := range ac.AlertAction.CustomButtons {
			if customButton.ActionName != "" {
				alertActions = append(alertActions, customButton.ActionName)
			}
		}
		return alertActions, nil
	default:
		return nil, NewUserPresentableError(fmt.Sprintf("Unsupported action type: \"%v\"", ac.ActionType))
	}
}

func (ac *ActionContainer) ValidateReturningUserReadableIssue() string {
	if ac.ActionType == "" {
		return "Empty actionType"
	}

	switch ac.ActionType {
	case ActionTypeEnumBanner:
		if ac.BannerAction == nil {
			return "Missing valid banner action data when actionType=banner"
		}
		return ac.BannerAction.ValidateReturningUserReadableIssue()
	case ActionTypeEnumAlert:
		if ac.AlertAction == nil {
			return "Missing valid banner action data when actionType=alert"
		}
		return ac.AlertAction.ValidateReturningUserReadableIssue()
	default:
		return fmt.Sprintf("Unsupported action type: \"%v\"", ac.ActionType)
	}
}
