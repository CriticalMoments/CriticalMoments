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
	namedThemes  map[string]*Theme

	// Actions
	namedActions map[string]ActionContainer

	// Triggers
	namedTriggers map[string]*Trigger
}

func (pc PrimaryConfig) ThemeWithName(name string) *Theme {
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
	TriggerConfig *jsonThemesSection `json:"triggers"`
}

type jsonThemesSection struct {
	DefaultTheme *jsonTheme           `json:"defaultTheme"`
	NamedThemes  map[string]jsonTheme `json:"namedThemes"`
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
		namedThemes:   make(map[string]*Theme),
		namedTriggers: make(map[string]*Trigger),
	}

	// Themes
	// TODO test no default theme field at all
	if jpc.ThemesConfig.DefaultTheme != nil {
		defaultTheme, err := NewThemeFromJsonTheme(jpc.ThemesConfig.DefaultTheme)
		if err != nil {
			return nil, NewUserPresentableErrorWSource("Default theme not a valid theme", err)
		}
		pc.DefaultTheme = defaultTheme
	}
	// TODO test no namedThemes field at all
	for themeName, themeJsonConfig := range jpc.ThemesConfig.NamedThemes {
		if themeName == "" {
			return nil, NewUserPresentableError("Named themes with empty name")
		}
		theme, err := NewThemeFromJsonTheme(&themeJsonConfig)
		if err != nil || theme == nil {
			errString := fmt.Sprintf("Named theme \"%v\" not a valid theme", themeName)
			return nil, NewUserPresentableErrorWSource(errString, err)
		}
		pc.namedThemes[themeName] = theme
	}

	// Actions
	/*for actionName, actionJsonConfig := range jpc.ActionsConfig.NamedActions {
		if actionName == "" {
			return nil, NewUserPresentableError("Named action with empty name")
		}
		action, err := NewAc
		theme, err := NewThemeFromJsonTheme(&themeJsonConfig)
		if err != nil || theme == nil {
			errString := fmt.Sprintf("Named theme \"%v\" not a valid theme", themeName)
			return nil, NewUserPresentableErrorWSource(errString, err)
		}
		pc.namedThemes[themeName] = theme
	}*/
	// TODO test No named Actions
	// TODO invalid/empty names/keys?
	if jpc.ActionsConfig != nil && jpc.ActionsConfig.NamedActions != nil {
		pc.namedActions = jpc.ActionsConfig.NamedActions
	}

	if validationIssue := pc.ValidateReturningUserReadableIssue(); validationIssue != "" {
		return nil, NewUserPresentableError(validationIssue)
	}

	// TODO map action.primaryAction to a real action, and validate name exists

	return &pc, nil
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
