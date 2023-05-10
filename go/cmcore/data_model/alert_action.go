package datamodel

import (
	"encoding/json"
	"fmt"
)

/*

Default is an alert with a title, body, and okay button which dismisses but takes no other action.

 - A Title and/or body is required. Neither is not allowed. Both are suggested.
 - showOkayButton defaults to yes if omitted.
 - okButtonActionName defaults to no action of omitted
 - You can add a standard cancel button with 1 property: showCancelButton:true. It never performs an action other than dismiss. It is not shown by default
 - Buttons are ordered by platform convention. If you desire a separate order use all custom buttons.
    - ios order: Custom Buttons, Cancel, Ok.
	- other platforms:
 - While it's supported, you probably shouldn't have "Ok" and custom buttons. Logically "Ok" works alone for an informational message, or paired with cancel for a confimation message. Ok paired with several custom options is usually confusing to the user.
 - Button styles are automatic for Ok/Cancel, and manual for custom buttons
   - The Ok button will get a treatment following the platform guidelines. On iOS that means "preferred" if paired with cancel, and plain if solo.
   - Cancel button will get the plain visual treatment.
   - Custom buttons specify their visual treatment: normal, desructive (red), primary.
 - Alert style is based on the platform
   - the default is dialog, which is UIAlertControllerStyleAlert on iOS and a Material dialog style: https://m3.material.io/components/dialogs/specs#23e479cf-c5a6-4a8b-87b3-1202d51855ac
   - large is UIAlertControllerStyleActionSheet on iOS and the material fullscreen style: https://m3.material.io/components/dialogs/specs#bbf1acde-f8d2-4ae1-9d51-343e96c4ac20
 - You must have at least 1 button: Ok, Cancel, or a valid custom button
 - There is no theme support for alerts, they use the system native alert look

https://developer.apple.com/design/human-interface-guidelines/alerts

*/

const (
	AlertActionStyleEnumDialog string = "dialog"
	AlertActionStyleEnumLarge  string = "large"
)

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
	// IOS =, Android = Neutral
	AlertActionButtonStyleEnumNormal      string = "normal"
	AlertActionButtonStyleEnumDestructive string = "destructive"
	AlertActionButtonStyleEnumPrimary     string = "primary"
)

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
	Style              string                   `json:"style,omitempty"`
	CustomButtons      *[]jsonAlertCustomButton `json:"customButtons,omitempty"`
}

type jsonAlertCustomButton struct {
	Label      string `json:"label"`
	ActionName string `json:"actionName"`
	Style      string `json:"style"`
}

func (a *AlertAction) Validate() bool {
	return a.ValidateReturningUserReadableIssue() == ""
}

func (a *AlertAction) ValidateReturningUserReadableIssue() string {
	if a.Title == "" && a.Message == "" {
		return "Alerts must have a title and/or a message. Both can not be blank."
	}
	if a.Style != AlertActionStyleEnumDialog && a.Style != AlertActionStyleEnumLarge {
		return "Alert style must be 'dialog' or 'large'"
	}
	if !a.ShowOkButton && a.OkButtonActionName != "" {
		return "For an alert, the okay button is hidden via 'showOkButton:false' but an Ok action name is specified. Either show the okay button or remove the action."
	}
	if !a.ShowOkButton && len(a.CustomButtons) == 0 {
		return "Alert must have an ok button and/or custom buttons."
	}
	for i, customButtom := range a.CustomButtons {
		if customButtonIssue := customButtom.ValidateReturningUserReadableIssue(); customButtonIssue != "" {
			return fmt.Sprintf("For an alert, button at index %v had issue \"%v\"", i, customButtonIssue)
		}
	}

	return ""
}

func (b *AlertActionCustomButton) Validate() bool {
	return b.ValidateReturningUserReadableIssue() == ""
}

func (b *AlertActionCustomButton) ValidateReturningUserReadableIssue() string {
	if b.Label == "" {
		return "Custom alert buttons must have a label"
	}
	if b.Style != AlertActionButtonStyleEnumNormal && b.Style != AlertActionButtonStyleEnumPrimary && b.Style != AlertActionButtonStyleEnumDestructive {
		return "Custom alert buttons must have a valid style: normal, destuctive, primary"
	}

	return ""
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
	if ja.Style != "" {
		alertStyle = ja.Style
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
			b := customButtonFromJson(&customButtonJson)
			customButtons = append(customButtons, b)
		}
	}
	a.CustomButtons = customButtons

	if validationIssue := a.ValidateReturningUserReadableIssue(); validationIssue != "" {
		return NewUserPresentableError(validationIssue)
	}

	return nil
}

func customButtonFromJson(jb *jsonAlertCustomButton) *AlertActionCustomButton {
	buttonStyle := AlertActionButtonStyleEnumNormal
	if jb.Style != "" {
		buttonStyle = jb.Style
	}

	return &AlertActionCustomButton{
		Label:      jb.Label,
		ActionName: jb.ActionName,
		Style:      buttonStyle,
	}
}
