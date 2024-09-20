package datamodel

import (
	"encoding/json"
	"fmt"

	"golang.org/x/exp/slices"
)

/*

Default is an alert with a title, body, and okay button which dismisses but takes no other action.

 - A Title and/or body is required. Neither is not allowed. Both are suggested.
 - showOkayButton defaults to yes if omitted.
 - okButtonActionName defaults to no action of omitted
 - You can add a standard cancel button with 1 property: showCancelButton:true. It never performs an action other than dismiss. It is not shown by default. It's position is optimized for the UI standards of the platform, so if you need a cancel button it's preferred to use showCancelButton:true over a custom button with the word "Cancel"
 - Buttons are ordered by platform convention (example, cancel moves to bottom on sheets, and left on alerts). If you desire an exact order use all custom buttons.
 - While it's supported, you probably shouldn't have "Ok" and custom buttons. Logically "Ok" works alone for an informational message, or paired with cancel for a confimation message. Ok paired with several custom options is usually confusing to the user.
 - Button styles are automatic for Ok/Cancel, and manual for custom buttons
   - The Ok button will get a treatment following the platform guidelines. On iOS that means "preferred" if paired with cancel, and plain if solo.
   - Cancel button treatment depends on system UI standards. It's highlighting and position are controlled by iOS depending on the type of alert (dialog/sheet) and number of buttons.
   - Custom buttons specify their visual treatment: normal, destructive (red), primary. Only one button can be primary on iOS - the last button you specify as primary (including okay) will get the primary treatment.
 - Cancel and OK buttons are localized using the system UI localization
 - Alert style is based on the platform
   - the default is "dialog", which is UIAlertControllerStyleAlert on iOS and a Material dialog style: https://m3.material.io/components/dialogs/specs#23e479cf-c5a6-4a8b-87b3-1202d51855ac
   - "large" is UIAlertControllerStyleActionSheet on iOS and the material fullscreen style: https://m3.material.io/components/dialogs/specs#bbf1acde-f8d2-4ae1-9d51-343e96c4ac20
 - You must have at least 1 button: Ok, or a valid custom button. No buttons or cancel alone aren't valid
 - There is no theme support for alerts, they use the system native alert look

https://developer.apple.com/design/human-interface-guidelines/alerts

*/

const (
	AlertActionStyleEnumDialog string = "dialog"
	AlertActionStyleEnumLarge  string = "large"
)

var alertStyles []string = []string{AlertActionStyleEnumDialog, AlertActionStyleEnumLarge}

type AlertAction struct {
	Title              string
	Message            string
	ShowCancelButton   bool
	ShowOkButton       bool
	OkButtonActionName string
	Style              string // AlertActionStyleEnum
	CustomButtons      []*AlertActionCustomButton
}

const (
	AlertActionButtonStyleEnumDefault     string = "default"
	AlertActionButtonStyleEnumDestructive string = "destructive"
	AlertActionButtonStyleEnumPrimary     string = "primary"
)

var alertActionStyles []string = []string{AlertActionButtonStyleEnumDefault, AlertActionButtonStyleEnumDestructive, AlertActionButtonStyleEnumPrimary}

type AlertActionCustomButton struct {
	Label      string
	ActionName string
	Style      string // AlertActionButtonStyleEnum
}

type jsonAlertAction struct {
	Title              string                   `json:"title,omitempty"`
	Messsage           string                   `json:"message,omitempty"`
	ShowOkButton       *bool                    `json:"showOkButton,omitempty"`
	OkButtonActionName string                   `json:"okButtonActionName,omitempty"`
	ShowCancelButton   *bool                    `json:"showCancelButton,omitempty"`
	Style              *string                  `json:"style,omitempty"`
	CustomButtons      *[]jsonAlertCustomButton `json:"customButtons,omitempty"`
}

type jsonAlertCustomButton struct {
	Label      string  `json:"label"`
	ActionName string  `json:"actionName"`
	Style      *string `json:"style"`
}

func unpackAlertFromJson(rawJson json.RawMessage, ac *ActionContainer) (ActionTypeInterface, error) {
	var alert AlertAction
	err := json.Unmarshal(rawJson, &alert)
	if err != nil {
		return nil, err
	}
	ac.AlertAction = &alert
	return &alert, nil
}

func (a *AlertAction) Valid() bool {
	return a.Check() == nil
}

func (a *AlertAction) Check() UserPresentableErrorInterface {
	if a.Title == "" && a.Message == "" {
		return NewUserPresentableError("Alerts must have a title and/or a message. Both can not be blank.")
	}
	if !slices.Contains(alertStyles, a.Style) {
		return NewUserPresentableError("Alert style must be 'dialog' or 'large'")
	}
	if !a.ShowOkButton && a.OkButtonActionName != "" {
		return NewUserPresentableError("For an alert, the okay button is hidden via 'showOkButton:false' but an Ok action name is specified. Either show the okay button or remove the action.")
	}
	if !a.ShowOkButton && len(a.CustomButtons) == 0 {
		return NewUserPresentableError("Alert must have an ok button and/or custom buttons.")
	}
	for i, customButtom := range a.CustomButtons {
		if customButtonIssue := customButtom.Check(); customButtonIssue != nil {
			return NewUserPresentableErrorWSource(fmt.Sprintf("For an alert, button at index %v had issue.", i), customButtonIssue)
		}
	}

	return nil
}

func (b *AlertActionCustomButton) Valid() bool {
	return b.Check() == nil
}

func (b *AlertActionCustomButton) Check() UserPresentableErrorInterface {
	if b.Label == "" {
		return NewUserPresentableError("Custom alert buttons must have a label")
	}
	if !slices.Contains(alertActionStyles, b.Style) {
		return NewUserPresentableError(fmt.Sprintf("Custom alert buttons must have a valid style: default, primary, or destructive. \"%v\" is not valid.", b.Style))
	}

	return nil
}

func (a *AlertAction) UnmarshalJSON(data []byte) error {
	var ja jsonAlertAction
	err := json.Unmarshal(data, &ja)
	if err != nil {
		return NewUserPresentableErrorWSource("Unable to parse the json of an action with type=alert. Check the format, variable names, and types (eg float vs int).", err)
	}

	// Default Values for nullable options
	showOkButton := true
	if ja.ShowOkButton != nil {
		showOkButton = *ja.ShowOkButton
	}
	showCancelButton := false
	if ja.ShowCancelButton != nil {
		showCancelButton = *ja.ShowCancelButton
	}
	alertStyle := AlertActionStyleEnumDialog
	if ja.Style != nil {
		if slices.Contains(alertStyles, *ja.Style) {
			alertStyle = *ja.Style
		} else {
			// Backwards compatibility -- default to dialog if this client doesn't recognize the style
			if StrictDatamodelParsing {
				styleErr := fmt.Sprintf("invalid 'style' tag found in config file under an alert action: \"%v\"", *ja.Style)
				return NewUserPresentableError(styleErr)
			}
		}
	}

	a.Title = ja.Title
	a.Message = ja.Messsage
	a.ShowCancelButton = showCancelButton
	a.ShowOkButton = showOkButton
	a.OkButtonActionName = ja.OkButtonActionName
	a.Style = alertStyle

	customButtons := make([]*AlertActionCustomButton, 0)
	if ja.CustomButtons != nil {
		for _, customButtonJson := range *ja.CustomButtons {
			b, err := customButtonFromJson(&customButtonJson)
			if err != nil {
				return err
			}
			customButtons = append(customButtons, b)
		}
	}
	a.CustomButtons = customButtons

	if err := a.Check(); err != nil {
		return NewUserErrorForJsonIssue(data, err)
	}

	return nil
}

func customButtonFromJson(jb *jsonAlertCustomButton) (*AlertActionCustomButton, error) {
	buttonStyle := AlertActionButtonStyleEnumDefault
	if jb.Style != nil {
		if slices.Contains(alertActionStyles, *jb.Style) {
			buttonStyle = *jb.Style
		} else {
			// Backwards compatibility: fall back to default style if this isn't recognized by this client
			if StrictDatamodelParsing {
				btnStyleErr := fmt.Sprintf("invalid alert action 'style' tag found in config. \"%v\" is not a valid style", *jb.Style)
				return nil, NewUserPresentableError(btnStyleErr)
			}
		}
	}

	return &AlertActionCustomButton{
		Label:      jb.Label,
		ActionName: jb.ActionName,
		Style:      buttonStyle,
	}, nil
}

func (a *AlertAction) AllEmbeddedThemeNames() ([]string, error) {
	return []string{}, nil
}

func (a *AlertAction) AllEmbeddedActionNames() ([]string, error) {
	alertActions := []string{}
	if a.OkButtonActionName != "" {
		alertActions = append(alertActions, a.OkButtonActionName)
	}
	for _, customButton := range a.CustomButtons {
		if customButton.ActionName != "" {
			alertActions = append(alertActions, customButton.ActionName)
		}
	}
	return alertActions, nil
}

func (l *AlertAction) AllEmbeddedConditions() ([]*Condition, error) {
	return []*Condition{}, nil
}

func (a *AlertAction) PerformAction(ab ActionBindings, actionName string) error {
	return ab.ShowAlert(a, actionName)
}

func (a *AlertAction) CustomButtonsCount() int {
	return len(a.CustomButtons)
}

func (a *AlertAction) CustomButtonAtIndex(i int) *AlertActionCustomButton {
	return a.CustomButtons[i]
}
