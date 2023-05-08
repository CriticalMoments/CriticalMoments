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
	NamedTriggers map[string]jsonTrigger `json:"namedTriggers"`
}

func NewPrimaryConfigFromJson(data []byte) (*PrimaryConfig, error) {
	var jpc jsonPrimaryConfig
	err := json.Unmarshal(data, &jpc)
	if err != nil {
		return nil, NewUserPresentableErrorWSource("Unable to parse config -- invalid json", err)
	}

	pc := PrimaryConfig{
		ConfigVersion: jpc.ConfigVersion,
	}

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

	pc.buildTriggersFromJsonTriggers(jpc.TriggerConfig.NamedTriggers)

	if validationIssue := pc.ValidateReturningUserReadableIssue(); validationIssue != "" {
		return nil, NewUserPresentableError(validationIssue)
	}

	// TODO map action.primaryAction to a real action, and validate name exists

	return &pc, nil
}

func (pc PrimaryConfig) buildTriggersFromJsonTriggers(jsonTriggers map[string]jsonTrigger) *UserPresentableError {
	// TODO test missing and {} triggers
	// TODO decide if empty map or nil map is right for missing
	if jsonTriggers == nil {
		return nil
	}

	namedTriggers := make(map[string]Trigger)
	for tName, jt := range jsonTriggers {
		action, ok := pc.namedActions[jt.ActionName]
		// TODO: test all 3 cases
		if !ok {
			return NewUserPresentableError(fmt.Sprintf("Trigger included named action \"%v\", which doesn't exist", jt.ActionName))
		}
		if tName == "" {
			return NewUserPresentableError("Empty trigger name")
		}
		if jt.EventName == "" {
			return NewUserPresentableError("Empty event name in trigger")
		}
		trigger := Trigger{
			Action:    action,
			EventName: jt.EventName,
		}
		namedTriggers[tName] = trigger
	}

	pc.namedTriggers = namedTriggers
	return nil
}

func (pc PrimaryConfig) Validate() bool {
	return pc.ValidateReturningUserReadableIssue() == ""
}

func (pc PrimaryConfig) ValidateReturningUserReadableIssue() string {
	if pc.ConfigVersion != "v1" {
		return "Config must have a config version of v1"
	}

	return ""
}
