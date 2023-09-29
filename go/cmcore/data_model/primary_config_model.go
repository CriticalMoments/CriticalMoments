package datamodel

import (
	"encoding/json"
	"fmt"

	"github.com/CriticalMoments/CriticalMoments/go/cmcore"
)

// Enables "Strict mode" validation for datamodel parsing
// Should only be enabled where we know the library is the latest version, which typically
// is not true. Will throw errors for unrecognized types, which could break forwards compatibility.
var StrictDatamodelParsing = false

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

	// Conditions
	namedConditions map[string]*Condition
}

func (pc *PrimaryConfig) ThemeWithName(name string) *Theme {
	theme, ok := pc.namedThemes[name]
	if ok {
		return &theme
	}
	return nil
}

func (pc *PrimaryConfig) ActionWithName(name string) *ActionContainer {
	action, ok := pc.namedActions[name]
	if ok {
		return &action
	}
	return nil
}

func (pc *PrimaryConfig) ConditionWithName(name string) *Condition {
	c, ok := pc.namedConditions[name]
	if ok {
		return c
	}
	return nil
}

func (pc *PrimaryConfig) ActionsForEvent(eventName string) []ActionContainer {
	// TODO P2: don't iterate, use a map
	actions := make([]ActionContainer, 0)
	for _, trigger := range pc.namedTriggers {
		if trigger.EventName == eventName {
			action, ok := pc.namedActions[trigger.ActionName]
			if ok {
				actions = append(actions, action)
			}
		}
	}
	return actions
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

	// Conditions
	ConditionsConfig *jsonConditionsSection `json:"conditions"`
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

type jsonConditionsSection struct {
	NamedConditions map[string]string `json:"namedConditions"`
}

func (pc *PrimaryConfig) UnmarshalJSON(data []byte) error {
	var jpc jsonPrimaryConfig
	err := json.Unmarshal(data, &jpc)
	if err != nil {
		return cmcore.NewUserPresentableErrorWSource("Unable to parse config -- invalid json", err)
	}

	pc.ConfigVersion = jpc.ConfigVersion

	// Themes
	if jpc.ThemesConfig != nil {
		if jpc.ThemesConfig.DefaultTheme != nil {
			pc.DefaultTheme = jpc.ThemesConfig.DefaultTheme
		}
		if jpc.ThemesConfig.NamedThemes != nil {
			pc.namedThemes = jpc.ThemesConfig.NamedThemes
		}
	}
	if pc.namedThemes == nil {
		pc.namedThemes = make(map[string]Theme)
	}

	// Actions
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

	// Conditions
	pc.namedConditions = make(map[string]*Condition)
	if jpc.ConditionsConfig != nil && jpc.ConditionsConfig.NamedConditions != nil {
		for name, conditionString := range jpc.ConditionsConfig.NamedConditions {
			condition, err := NewCondition(conditionString)
			if err != nil && StrictDatamodelParsing {
				return err
			} else if err != nil {
				// Fallback to conditions that always evaluates to false.
				condition, _ = NewCondition("false")
			}
			pc.namedConditions[name] = condition
		}
	}

	if validationIssue := pc.ValidateReturningUserReadableIssue(); validationIssue != "" {
		return cmcore.NewUserPresentableError(validationIssue)
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
		return "Internal issue: code 7842371152"
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

	// Run nested validations
	return pc.validateNestedReturningUserReadableIssue()
}

func (pc PrimaryConfig) validateNestedReturningUserReadableIssue() string {
	if pc.DefaultTheme != nil {
		if defaultThemeIssue := pc.DefaultTheme.ValidateReturningUserReadableIssue(); defaultThemeIssue != "" {
			return defaultThemeIssue
		}
	}
	for themeName, theme := range pc.namedThemes {
		if themeIssue := theme.ValidateReturningUserReadableIssue(); themeIssue != "" {
			return fmt.Sprintf("Theme \"%v\" had issue: %v", themeName, themeIssue)
		}
	}
	for actionName, action := range pc.namedActions {
		if actionValidationIssue := action.ValidateReturningUserReadableIssue(); actionValidationIssue != "" {
			return fmt.Sprintf("Action \"%v\" had issue: %v", actionName, actionValidationIssue)
		}
	}
	for triggerName, trigger := range pc.namedTriggers {
		if triggerIssue := trigger.ValidateReturningUserReadableIssue(); triggerIssue != "" {
			return fmt.Sprintf("Trigger \"%v\" had issue: %v", triggerName, triggerIssue)
		}
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
		if action.ActionType == "" || action.actionData == nil {
			return "Internal issue. Code 15234328"
		}
		themeList, err := action.actionData.AllEmbeddedThemeNames()
		if err != nil || themeList == nil {
			return fmt.Sprintf("Internal issue for action \"%v\". Code: 88456198", sourceActionName)
		}
		for _, themeName := range themeList {
			_, ok := pc.namedThemes[themeName]
			if !ok {
				return fmt.Sprintf("Action \"%v\" specified named theme \"%v\", which doesn't exist", sourceActionName, themeName)
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
		if action.ActionType == "" || action.actionData == nil {
			return "Internal issue. Code 98347134"
		}
		actionList, err := action.actionData.AllEmbeddedActionNames()
		if err != nil || actionList == nil {
			return fmt.Sprintf("Internal issue for action \"%v\". Code: 798853616", sourceActionName)
		}
		for _, actionName := range actionList {
			_, ok := pc.namedActions[actionName]
			if !ok {
				return fmt.Sprintf("Action \"%v\" specified named action \"%v\", which doesn't exist", sourceActionName, actionName)
			}
		}
	}

	return ""
}
