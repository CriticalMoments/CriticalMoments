package datamodel

import (
	"encoding/json"
	"fmt"
)

type PrimaryConfig struct {
	// Version number
	ConfigVersion string

	// Themes
	DefaultTheme *Theme
	namedThemes  map[string]Theme

	// Actions
	namedActions map[string]ActionContainer

	// Triggers
	namedTriggers map[string]Trigger
}

func (pc PrimaryConfig) ThemeWithName(name string) Theme {
	return pc.namedThemes[name]
}

type jsonPrimaryConfig struct {
	// Version number
	ConfigVersion string `json:"configVersion"`

	// Themes
	ThemesConfig *jsonThemesSection `json:"themes"`

	// Actions
	ActionsConfig *jsonActionsSection `json:"actions"`

	// Triggers
	TriggerConfig *jsonTriggersSection `json:"triggers"`
}

type jsonThemesSection struct {
	DefaultTheme *Theme           `json:"defaultTheme"`
	NamedThemes  map[string]Theme `json:"namedThemes"`
}

type jsonActionsSection struct {
	NamedActions map[string]ActionContainer `json:"namedActions"`
}

type jsonTriggersSection struct {
	NamedTriggers map[string]Trigger `json:"namedTriggers"`
}

func (pc *PrimaryConfig) UnmarshalJSON(data []byte) error {
	var jpc jsonPrimaryConfig
	err := json.Unmarshal(data, &jpc)
	if err != nil {
		return NewUserPresentableErrorWSource("Unable to parse config -- invalid json", err)
	}

	pc.ConfigVersion = jpc.ConfigVersion

	// Themes
	// TODO test no default/named/root theme field at all
	if jpc.ThemesConfig != nil {
		if jpc.ThemesConfig.DefaultTheme != nil {
			pc.DefaultTheme = jpc.ThemesConfig.DefaultTheme
		}
		if jpc.ThemesConfig.NamedThemes != nil {
			pc.namedThemes = jpc.ThemesConfig.NamedThemes
		}
	}

	// Actions
	// TODO test No named Actions
	// TODO invalid/empty names/keys?
	if jpc.ActionsConfig != nil && jpc.ActionsConfig.NamedActions != nil {
		pc.namedActions = jpc.ActionsConfig.NamedActions
	}

	// Triggers
	if jpc.TriggerConfig != nil && jpc.TriggerConfig.NamedTriggers != nil {
		pc.namedTriggers = jpc.TriggerConfig.NamedTriggers
	}

	if validationIssue := pc.ValidateReturningUserReadableIssue(); validationIssue != "" {
		return NewUserPresentableError(validationIssue)
	}

	return nil
}

func (pc PrimaryConfig) Validate() bool {
	return pc.ValidateReturningUserReadableIssue() == ""
}

func (pc PrimaryConfig) ValidateReturningUserReadableIssue() string {
	if pc.ConfigVersion != "v1" {
		return "Config must have a config version of v1"
	}

	// Validate the actions in the trigger actually exist
	if pc.namedTriggers != nil {
		for tName, t := range pc.namedTriggers {
			_, ok := pc.namedActions[t.ActionName]
			// TODO: test all 3 cases
			if !ok {
				return fmt.Sprintf("Trigger included named action \"%v\", which doesn't exist", t.ActionName)
			}
			if tName == "" {
				return "Empty trigger name"
			}
			if t.EventName == "" {
				return "Empty/nil event name in trigger"
			}
		}
	}

	// TODO validate action.primaryAction that a real action by that name exists

	// Missing: empty or nil expected for each type/set? Add checks and/or tests

	return ""
}
