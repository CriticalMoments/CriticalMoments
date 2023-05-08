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
		} else {
			pc.namedThemes = make(map[string]Theme)
		}
	}

	// Actions
	// TODO test No named Actions
	// TODO invalid/empty names/keys?
	if jpc.ActionsConfig != nil && jpc.ActionsConfig.NamedActions != nil {
		pc.namedActions = jpc.ActionsConfig.NamedActions
	} else {
		pc.namedActions = map[string]ActionContainer{}
	}

	// Triggers
	if jpc.TriggerConfig != nil && jpc.TriggerConfig.NamedTriggers != nil {
		pc.namedTriggers = jpc.TriggerConfig.NamedTriggers
	} else {
		pc.namedTriggers = make(map[string]Trigger)
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

	if pc.namedActions == nil || pc.namedThemes == nil || pc.namedTriggers == nil {
		return "Internal issue initializing config: missing map for actions, themes or triggers"
	}

	if actionIssue := pc.validateEmbeddedActionsExistReturningUserReadable(); actionIssue != "" {
		return actionIssue
	}
	if themeIssue := pc.validateThemeNamesExistReturningUserReadable(); themeIssue != "" {
		return themeIssue
	}
	if emptyKeyIssue := pc.validateMapsDontContainEmptyStringReturningUserReadable(); emptyKeyIssue != "" {
		return emptyKeyIssue
	}

	return ""
}

func (pc PrimaryConfig) validateMapsDontContainEmptyStringReturningUserReadable() string {
	_, themeFound := pc.namedThemes[""]
	_, actionFound := pc.namedActions[""]
	_, triggerFound := pc.namedTriggers[""]
	if triggerFound || actionFound || themeFound {
		return "The empty string \"\" is not a valid name for actions/themes/triggers"
	}
	return ""
}

func (pc PrimaryConfig) validateThemeNamesExistReturningUserReadable() string {
	for sourceActionName, action := range pc.namedActions {
		if action.ThemeName != "" {
			_, ok := pc.namedThemes[action.ThemeName]
			if !ok {
				return fmt.Sprintf("Action \"%v\" included named theme \"%v\", which doesn't exist", sourceActionName, action.ThemeName)
			}
		}
	}

	return ""
}

func (pc PrimaryConfig) validateEmbeddedActionsExistReturningUserReadable() string {
	// Validate the actions in the trigger actually exist
	for tName, t := range pc.namedTriggers {
		_, ok := pc.namedActions[t.ActionName]
		if !ok {
			return fmt.Sprintf("Trigger \"%v\" included named action \"%v\", which doesn't exist", tName, t.ActionName)
		}
	}

	// validate any named actions embedded in other actions actually exist
	for sourceActionName, action := range pc.namedActions {
		actionList, err := action.AllEmbeddedActionNames()
		if err != nil || actionList == nil {
			return fmt.Sprintf("Unknown issue confirming all named actions in action \"%v\"exist in config", sourceActionName)
		}
		for _, actionName := range actionList {
			_, ok := pc.namedActions[actionName]
			if !ok {
				return fmt.Sprintf("Action \"%v\" included named action \"%v\", which doesn't exist", sourceActionName, actionName)
			}
		}
	}

	return ""
}
