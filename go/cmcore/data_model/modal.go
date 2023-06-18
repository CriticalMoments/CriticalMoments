package datamodel

import (
	"encoding/json"
)

type ModalAction struct {
	Content         Page
	ShowCloseButton bool
	CustomThemeName string
}

type jsonModalAction struct {
	Content         Page   `json:"content"`
	ShowCloseButton *bool  `json:"showCloseButton,omitempty"`
	CustomThemeName string `json:"customThemeName,omitempty"`
}

func unpackModalFromJson(rawJson json.RawMessage, ac *ActionContainer) (ActionTypeInterface, error) {
	var modal ModalAction
	err := json.Unmarshal(rawJson, &modal)
	if err != nil {
		return nil, err
	}
	ac.ModalAction = &modal
	return &modal, nil
}

func (m *ModalAction) Validate() bool {
	return m.ValidateReturningUserReadableIssue() == ""
}

func (m *ModalAction) ValidateReturningUserReadableIssue() string {
	if contentErr := m.Content.ValidateReturningUserReadableIssue(); contentErr != "" {
		return contentErr
	}

	return ""
}

func (m *ModalAction) UnmarshalJSON(data []byte) error {
	var jm jsonModalAction
	err := json.Unmarshal(data, &jm)
	if err != nil {
		return NewUserPresentableErrorWSource("Unable to parse the json of an action with type=modal. Check the format, variable names, and types.", err)
	}

	// Defaults
	showCloseButton := true
	if jm.ShowCloseButton != nil {
		showCloseButton = *jm.ShowCloseButton
	}

	m.Content = jm.Content
	m.CustomThemeName = jm.CustomThemeName
	m.ShowCloseButton = showCloseButton

	if validationIssue := m.ValidateReturningUserReadableIssue(); validationIssue != "" {
		return NewUserPresentableError(validationIssue)
	}

	return nil
}

func (m *ModalAction) AllEmbeddedThemeNames() ([]string, error) {
	if m.CustomThemeName == "" {
		return []string{}, nil
	}
	return []string{m.CustomThemeName}, nil
}

func (m *ModalAction) AllEmbeddedActionNames() ([]string, error) {
	var embeddedActions []string = []string{}
	for _, button := range m.Content.Buttons {
		if button.ActionName != "" {
			embeddedActions = append(embeddedActions, button.ActionName)
		}
	}
	return embeddedActions, nil
}

func (m *ModalAction) PerformAction(ab ActionBindings) error {
	return ab.ShowModal(m)
}
