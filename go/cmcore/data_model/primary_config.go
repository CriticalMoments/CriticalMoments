package datamodel

import (
	"bufio"
	"bytes"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model/conditions"
	"github.com/CriticalMoments/CriticalMoments/go/cmcore/signing"
	"golang.org/x/exp/slices"
)

// Enables "Strict mode" validation for datamodel parsing
// Should only be enabled where we know the library is the latest version, which typically
// is not true. Will throw errors for unrecognized types, which could break forwards compatibility.
var StrictDatamodelParsing = false

type PrimaryConfig struct {
	// Metadata
	ContainerVersion string
	ConfigVersion    string
	AppId            string

	MinCMVersion         string // for SDK users
	MinCMVersionInternal string // for internal use
	MinAppVersion        string

	// Themes
	defaultTheme     *Theme
	LibraryThemeName string
	namedThemes      map[string]*Theme

	// Actions
	namedActions map[string]*ActionContainer

	// Triggers
	namedTriggers map[string]*Trigger

	// Conditions
	namedConditions map[string]*Condition

	// Notifications
	Notifications map[string]*Notification
}

func (pc *PrimaryConfig) DefaultTheme() *Theme {
	return pc.themeIteratingFallbacks(pc.defaultTheme)
}

func (pc *PrimaryConfig) ThemeWithName(name string) *Theme {
	theme, ok := pc.namedThemes[name]
	if ok {
		return pc.themeIteratingFallbacks(theme)
	}

	theme, err := builtInThemeByName(name)
	if theme != nil && err == nil {
		return pc.themeIteratingFallbacks(theme)
	}

	return nil
}

func (pc *PrimaryConfig) IncludesCustomThemes() bool {
	return len(pc.namedThemes) > 0
}

func (pc *PrimaryConfig) ActionWithName(name string) *ActionContainer {
	action, ok := pc.namedActions[name]
	if ok {
		return action
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

func (pc *PrimaryConfig) TriggersForEvent(eventName string) []*Trigger {
	triggers := make([]*Trigger, 0)
	for _, trigger := range pc.namedTriggers {
		if trigger.EventName == eventName {
			triggers = append(triggers, trigger)
		}
	}
	return triggers
}

// Container Decoding

const primaryConfigConfigPemBlock = "CONFIG"
const primaryConfigHeadPemBlock = "CM"
const primaryConfigConfigSignatureHeader = "Signature"
const primaryConfigHeadContainerVersion = "Container-Version"

func DecodePrimaryConfig(data []byte, signUtil *signing.SignUtil) (*PrimaryConfig, error) {
	pc := &PrimaryConfig{}

	var rest []byte
	rest = data

	var configSignature string
	var configBytes []byte

	for len(rest) > 0 {
		var block *pem.Block
		block, rest = pem.Decode(rest)
		if block == nil {
			break
		}

		switch block.Type {
		case primaryConfigConfigPemBlock:
			configBytes = block.Bytes
			configSignature = block.Headers[primaryConfigConfigSignatureHeader]
			err := json.Unmarshal(block.Bytes, pc)
			if err != nil {
				return nil, err
			}
		case primaryConfigHeadPemBlock:
			err := pc.ParseHeadBlock(block)
			if err != nil {
				return nil, err
			}
		}
	}

	// Validate CM block
	if pc.ContainerVersion == "" {
		return nil, NewUserPresentableError("Config file not signed: no valid CM block found in config file")
	}

	// Validate CONFIG block
	if len(configBytes) == 0 {
		return nil, NewUserPresentableError("No CONFIG block found in config file")
	}
	err := ValidateSignature(signUtil, configBytes, configSignature)
	if err != nil {
		return nil, err
	}
	configErr := pc.Check()
	if configErr != nil {
		return nil, configErr
	}

	return pc, nil
}

func (pc *PrimaryConfig) ParseHeadBlock(b *pem.Block) error {
	pc.ContainerVersion = b.Headers[primaryConfigHeadContainerVersion]
	// We bump container version to 2+ when we want to break backwards compatibility.
	// We try to parse all v1 versions (including new subversions)
	if pc.ContainerVersion != "v1" && !strings.HasPrefix(pc.ContainerVersion, "v1.") {
		return NewUserPresentableError("Unsupported container version")
	}
	return nil
}

func ValidateSignature(su *signing.SignUtil, configBytes []byte, sig string) error {
	if sig == "" {
		return NewUserPresentableError("Missing Config Signature. Please sign your config at https://criticalmoments.io")
	}
	valid, err := su.VerifyMessage(configBytes, sig)
	if err != nil {
		return err
	}
	if !valid {
		return NewUserPresentableError("Configuration file invalid. The signature does not match the JSON body. Please re-sign your config at https://criticalmoments.io")
	}

	return nil
}

func EncodeConfig(configBytes []byte, su *signing.SignUtil) ([]byte, error) {
	var b bytes.Buffer
	r := bufio.NewWriter(&b)

	// Parse the config data to ensure it's valid
	pc := &PrimaryConfig{}
	err := json.Unmarshal(configBytes, pc)
	if err != nil {
		return nil, err
	}

	// Header block
	headerBlock := &pem.Block{
		Type: primaryConfigHeadPemBlock,
		Headers: map[string]string{
			primaryConfigHeadContainerVersion: "v1",
		},
		Bytes: []byte{},
	}
	err = pem.Encode(r, headerBlock)
	if err != nil {
		return nil, err
	}

	// Config block
	sig, err := su.SignMessage(configBytes)
	if err != nil {
		return nil, err
	}
	configBlock := &pem.Block{
		Type: primaryConfigConfigPemBlock,
		Headers: map[string]string{
			primaryConfigConfigSignatureHeader: sig,
		},
		Bytes: configBytes,
	}
	err = pem.Encode(r, configBlock)
	if err != nil {
		return nil, err
	}

	err = r.Flush()
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// JSON

type jsonPrimaryConfig struct {
	ConfigVersion        string `json:"configVersion"`
	AppId                string `json:"appId"`
	MinAppVersion        string `json:"minAppVersion"`
	MinCMVersion         string `json:"minCMVersion"`
	MinCMVersionInternal string `json:"minCMVersionInternal"`

	// Themes
	ThemesConfig *jsonThemesSection `json:"themes"`

	// Actions
	ActionsConfig *jsonActionsSection `json:"actions"`

	// Triggers
	TriggerConfig *jsonTriggersSection `json:"triggers"`

	// Conditions
	ConditionsConfig *jsonConditionsSection `json:"conditions"`

	// Notifications
	Notifications map[string]*Notification `json:"notifications"`
}

type jsonThemesSection struct {
	DefaultThemeName string            `json:"defaultThemeName"`
	NamedThemes      map[string]*Theme `json:"namedThemes"`
}

type jsonActionsSection struct {
	NamedActions map[string]*ActionContainer `json:"namedActions"`
}

type jsonTriggersSection struct {
	NamedTriggers map[string]*Trigger `json:"namedTriggers"`
}

type jsonConditionsSection struct {
	NamedConditions map[string]*Condition `json:"namedConditions"`
}

func (pc *PrimaryConfig) UnmarshalJSON(data []byte) error {
	var jpc jsonPrimaryConfig
	err := json.Unmarshal(data, &jpc)
	if err != nil {
		return NewUserPresentableErrorWSource("Invalid Critical Moments config file: ", err)
	}

	pc.ConfigVersion = jpc.ConfigVersion
	pc.AppId = jpc.AppId
	pc.MinAppVersion = jpc.MinAppVersion
	pc.MinCMVersion = jpc.MinCMVersion
	pc.MinCMVersionInternal = jpc.MinCMVersionInternal

	// Themes
	if jpc.ThemesConfig != nil {
		if jpc.ThemesConfig.NamedThemes != nil {
			pc.namedThemes = jpc.ThemesConfig.NamedThemes
		}
		if jpc.ThemesConfig.DefaultThemeName != "" {
			isLibaryTheme := libraryThemeNames[jpc.ThemesConfig.DefaultThemeName]
			appcoreBuiltInTheme, _ := builtInThemeByName(jpc.ThemesConfig.DefaultThemeName)
			namedDefaultTheme := pc.namedThemes[jpc.ThemesConfig.DefaultThemeName]

			// Priority order: named, library, cmcore built-in
			if namedDefaultTheme != nil {
				pc.defaultTheme = namedDefaultTheme
			} else if isLibaryTheme {
				pc.LibraryThemeName = jpc.ThemesConfig.DefaultThemeName
			} else if appcoreBuiltInTheme != nil {
				pc.defaultTheme = appcoreBuiltInTheme
			} else if StrictDatamodelParsing {
				return NewUserPresentableError("Default theme specified in config doesn't exist")
			} else {
				fmt.Println("CriticalMoments: WARNING - Default theme specified in config doesn't exist. Will fallback to system default theme.")
			}
		}
	}
	if pc.namedThemes == nil {
		pc.namedThemes = make(map[string]*Theme)
	}

	// Actions
	if jpc.ActionsConfig != nil && jpc.ActionsConfig.NamedActions != nil {
		pc.namedActions = jpc.ActionsConfig.NamedActions
	} else {
		pc.namedActions = map[string]*ActionContainer{}
	}

	// Triggers
	if jpc.TriggerConfig != nil && jpc.TriggerConfig.NamedTriggers != nil {
		pc.namedTriggers = jpc.TriggerConfig.NamedTriggers
	} else {
		pc.namedTriggers = make(map[string]*Trigger)
	}

	// Conditions
	if jpc.ConditionsConfig != nil && jpc.ConditionsConfig.NamedConditions != nil {
		pc.namedConditions = jpc.ConditionsConfig.NamedConditions
	} else {
		pc.namedConditions = make(map[string]*Condition)
	}

	// Notifications
	if jpc.Notifications != nil {
		for notificationID, notification := range jpc.Notifications {
			notification.ID = notificationID
		}
		pc.Notifications = jpc.Notifications
	} else {
		pc.Notifications = make(map[string]*Notification)
	}

	return pc.Check()
}

func (pc *PrimaryConfig) NameForActionContainer(c *ActionContainer) string {
	for name, action := range pc.namedActions {
		if action == c {
			return name
		}
	}
	return ""
}

func (pc *PrimaryConfig) NameForCondition(c *Condition) string {
	for name, condition := range pc.namedConditions {
		if condition == c {
			return name
		}
	}
	return ""
}

func (pc *PrimaryConfig) themeIteratingFallbacks(t *Theme) *Theme {
	if t == nil {
		return nil
	}
	// Limit depth of search, to avoid infinite loops
	for i := 0; i < 50; i++ {
		if t.ValidateDisallowFallthoughReturningUserReadableIssue() == nil {
			return t
		}
		if t.FallbackThemeName == "" {
			return nil
		}
		nextTheme, ok := pc.namedThemes[t.FallbackThemeName]
		if !ok {
			return nil
		}
		t = nextTheme
	}

	return nil
}

func (pc *PrimaryConfig) Valid() bool {
	return pc.Check() == nil
}

func (pc *PrimaryConfig) Check() UserPresentableErrorInterface {
	// We bump version to 2+ when we want to break backwards compatibility.
	// We try to parse all v1 versions (including new subversions)
	if pc.ConfigVersion != "v1" && !strings.HasPrefix(pc.ConfigVersion, "v1.") {
		return NewUserPresentableError("Config must have a config version of v1")
	}

	if pc.AppId == "" {
		return NewUserPresentableError("Config must have an appId")
	}

	if pc.MinAppVersion != "" {
		if _, err := conditions.VersionFromVersionString(pc.MinAppVersion); err != nil {
			return NewUserPresentableError(fmt.Sprintf("Config had invalid minAppVersion: %v", pc.MinAppVersion))
		}
	}
	if pc.MinCMVersion != "" {
		if _, err := conditions.VersionFromVersionString(pc.MinCMVersion); err != nil {
			return NewUserPresentableError(fmt.Sprintf("Config had invalid minCMVersion: %v", pc.MinCMVersion))
		}
	}
	if pc.MinCMVersionInternal != "" {
		if _, err := conditions.VersionFromVersionString(pc.MinCMVersionInternal); err != nil {
			return NewUserPresentableError(fmt.Sprintf("Config had invalid minCMVersionInternal: %v", pc.MinCMVersionInternal))
		}
	}

	// These should be initialized to empty maps
	if pc.namedActions == nil || pc.namedThemes == nil || pc.namedTriggers == nil || pc.namedConditions == nil || pc.Notifications == nil {
		return NewUserPresentableError("Internal issue. Nil maps. Code 7842371152")
	}

	if actionIssue := pc.validateEmbeddedActionsExistReturningUserReadable(); actionIssue != "" {
		return NewUserPresentableError(actionIssue)
	}
	if themeIssue := pc.validateThemeNamesExistReturningUserReadable(); themeIssue != "" {
		return NewUserPresentableError(themeIssue)
	}
	if emptyKeyIssue := pc.validateMapsDontContainEmptyStringReturningUserReadable(); emptyKeyIssue != "" {
		return NewUserPresentableError(emptyKeyIssue)
	}
	if fallbackNameIssue := pc.validateFallbackNames(); fallbackNameIssue != "" {
		return NewUserPresentableError(fallbackNameIssue)
	}

	// Run nested validations
	return pc.checkNested()
}

func (pc *PrimaryConfig) checkNested() UserPresentableErrorInterface {
	if pc.defaultTheme != nil {
		if defaultThemeIssue := pc.defaultTheme.Check(); defaultThemeIssue != nil {
			return defaultThemeIssue
		}
	}
	for themeName, theme := range pc.namedThemes {
		if themeIssue := theme.Check(); themeIssue != nil {
			return NewUserPresentableErrorWSource(fmt.Sprintf("Theme \"%v\" had issue. ", themeName), themeIssue)
		}
	}
	for actionName, action := range pc.namedActions {
		if actionValidationIssue := action.Check(); actionValidationIssue != nil {
			return NewUserPresentableErrorWSource(fmt.Sprintf("Action \"%v\" had issue", actionName), actionValidationIssue)
		}
	}
	for triggerName, trigger := range pc.namedTriggers {
		if triggerIssue := trigger.ValidateReturningUserReadableIssue(); triggerIssue != "" {
			// TODO_P0
			return NewUserPresentableError(fmt.Sprintf("Trigger \"%v\" had issue: %v", triggerName, triggerIssue))
		}
	}
	for notificationID, notification := range pc.Notifications {
		if notifIssue := notification.Check(); notifIssue != nil {
			return NewUserPresentableErrorWSource(fmt.Sprintf("Notification \"%v\" had issue.", notificationID), notifIssue)
		}
	}

	return nil
}

func (pc *PrimaryConfig) validateMapsDontContainEmptyStringReturningUserReadable() string {
	_, themeFound := pc.namedThemes[""]
	_, actionFound := pc.namedActions[""]
	_, triggerFound := pc.namedTriggers[""]
	if triggerFound || actionFound || themeFound {
		return "The empty string \"\" is not a valid name for actions/themes/triggers"
	}
	return ""
}

func (pc *PrimaryConfig) validateThemeNamesExistReturningUserReadable() string {
	allBuiltInThemeNames := AllBuiltInThemeNames()
	for sourceActionName, action := range pc.namedActions {
		if action.ActionType == "" || action.actionData == nil {
			return "Internal issue. Code 15234328"
		}
		themeList, err := action.actionData.AllEmbeddedThemeNames()
		if err != nil || themeList == nil {
			return fmt.Sprintf("Internal issue for action \"%v\". Code: 88456198", sourceActionName)
		}
		for _, themeName := range themeList {
			isBuiltIn := slices.Contains(allBuiltInThemeNames, themeName)
			_, hasNamedTheme := pc.namedThemes[themeName]
			if !hasNamedTheme && !isBuiltIn {
				// New built in names may be added later. Older devices should ignore unknown names.
				if StrictDatamodelParsing {
					return fmt.Sprintf("Action \"%v\" specified named theme \"%v\", which doesn't exist", sourceActionName, themeName)
				} else {
					fmt.Println("CriticalMoments: WARNING - Action specified named theme that doesn't exist. Will fallback to system default theme.")
				}
			}
		}
	}

	return ""
}

func (pc *PrimaryConfig) validateEmbeddedActionsExistReturningUserReadable() string {
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

	// Validate actions in notifications exist
	for id, notif := range pc.Notifications {
		if notif.ActionName != "" {
			_, ok := pc.namedActions[notif.ActionName]
			if !ok {
				return fmt.Sprintf("Notification '%v' has action name '%v', which does not exist", id, notif.ActionName)
			}
		}
	}

	return ""
}

func (pc *PrimaryConfig) validateFallbackNames() string {
	for themeName, theme := range pc.namedThemes {
		if theme.FallbackThemeName != "" {
			_, ok := pc.namedThemes[theme.FallbackThemeName]
			if !ok {
				return fmt.Sprintf("Theme \"%v\" specified fallback theme \"%v\", which doesn't exist", themeName, theme.FallbackThemeName)
			}
		}
	}

	if pc.defaultTheme != nil && pc.defaultTheme.FallbackThemeName != "" {
		_, ok := pc.namedThemes[pc.defaultTheme.FallbackThemeName]
		if !ok {
			return fmt.Sprintf("defaultTheme specified fallback theme \"%v\", which doesn't exist", pc.defaultTheme.FallbackThemeName)
		}
	}

	for actionName, action := range pc.namedActions {
		if action.FallbackActionName != "" {
			_, ok := pc.namedActions[action.FallbackActionName]
			if !ok {
				return fmt.Sprintf("Action \"%v\" specified fallback action \"%v\", which doesn't exist", actionName, action.FallbackActionName)
			}
		}
	}

	return ""
}

func (pc *PrimaryConfig) NamedConditionCount() int {
	return len(pc.namedConditions)
}

func (pc *PrimaryConfig) NamedConditionsConditionals() []string {
	conditionals := make([]string, 0, len(pc.namedConditions))
	for _, c := range pc.namedConditions {
		conditionals = append(conditionals, c.conditionString)
	}
	return conditionals
}

func (pc *PrimaryConfig) AllConditions() ([]*Condition, error) {
	all := make([]*Condition, 0)
	for _, c := range pc.namedConditions {
		all = append(all, c)
	}

	for _, a := range pc.namedActions {
		if a.Condition != nil {
			all = append(all, a.Condition)
		}
		if a.actionData != nil {
			actionConditions, err := a.actionData.AllEmbeddedConditions()
			if err != nil {
				return nil, err
			}
			all = append(all, actionConditions...)
		}
	}

	for _, t := range pc.namedTriggers {
		condition := t.Condition
		if condition != nil {
			all = append(all, condition)
		}
	}

	for _, n := range pc.Notifications {
		if n.ScheduleCondition != nil {
			all = append(all, n.ScheduleCondition)
		}
	}

	return all, nil
}

func (pc *PrimaryConfig) AllActions() []*ActionContainer {
	all := make([]*ActionContainer, 0)
	for _, a := range pc.namedActions {
		all = append(all, a)
	}
	return all
}
