package appcore

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/CriticalMoments/CriticalMoments/go/appcore/events"
	"github.com/CriticalMoments/CriticalMoments/go/cmcore"
	datamodel "github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model"
	"github.com/CriticalMoments/CriticalMoments/go/cmcore/signing"
)

func GoPing() string {
	return "AppcorePong->" + cmcore.CmCorePing()
}

type Appcore struct {
	started bool

	// Library binding/delegate
	libBindings LibBindings

	// API Key
	apiKey *signing.ApiKey

	// Primary configuration
	configUrlString string
	config          *datamodel.PrimaryConfig

	// Cache
	cache *cache

	// Event handler
	eventManager *events.EventManager

	// Properties
	propertyRegistry *propertyRegistry

	// Dev Mode namedCondition conflict check
	seenNamedConditions map[string]string
}

func NewAppcore() *Appcore {
	return &Appcore{
		propertyRegistry:    newPropertyRegistry(),
		seenNamedConditions: map[string]string{},
	}
}

// Hopefully no one wants http (no TLS) in 2023... but given the importance of the config file we can't open this up to injection attacks
const filePrefix = "file://"
const httpsPrefix = "https://"

func (ac *Appcore) SetConfigUrl(configUrl string) error {
	if !strings.HasPrefix(configUrl, filePrefix) && !strings.HasPrefix(configUrl, httpsPrefix) {
		return errors.New("Config URL must start with https:// or file://")
	}
	ac.configUrlString = configUrl

	return nil
}

func (ac *Appcore) SetApiKey(apiKey string, bundleID string) error {
	key, err := signing.ParseApiKey(apiKey)
	if err != nil {
		return errors.New("Invalid API Key. Please make sure you get your key from criticalmoments.io")
	}
	if v, err := key.Valid(); err != nil || !v {
		return errors.New("Invalid API Key. Please make sure you get your key from criticalmoments.io")
	}
	if key.BundleId() != bundleID {
		return errors.New(fmt.Sprintf("This API key isn't valid for this app. API key is for %s, but this app has bundle ID %s", key.BundleId(), bundleID))
	}
	ac.apiKey = key
	return nil
}

func (ac *Appcore) SetDataDirPath(dataDirPath string) error {
	cache, err := newCacheWithBaseDir(dataDirPath)
	if err != nil {
		return err
	}
	ac.cache = cache

	eventManager, err := events.NewEventManager(dataDirPath)
	if err != nil {
		return err
	}
	ac.eventManager = eventManager

	dbOperations := eventManager.EventManagerConditionFunctions()
	ac.propertyRegistry.RegisterDynamicFunctions(dbOperations)

	return nil
}

func (ac *Appcore) SetTimezoneGMTOffset(gmtOffset int) {
	tzName := fmt.Sprintf("UTCOffsetS:%v", gmtOffset)
	tz := time.FixedZone(tzName, gmtOffset)
	time.Local = tz

	ac.propertyRegistry.registerStaticProperty("timezone_gmt_offset", gmtOffset)
}

func (ac *Appcore) CheckNamedConditionCollision(name string, conditionString string) error {
	if name == "" {
		return nil
	}
	// in debug mode, track each built-in condition we see, and make sure the developer isn't reusing names
	// If they use the same name twice for different things, they won't be able to override in the future
	priorSeen := ac.seenNamedConditions[name]
	if priorSeen == "" {
		ac.seenNamedConditions[name] = conditionString
	} else if priorSeen != conditionString {
		return errors.New(fmt.Sprintf("The named condition \"%v\" is being used in multiple places in this codebase, with different fallback conditions (\"%v\" and \"%v\"). This will make it impossible to override each usage independently from remote configuration. Please use unique names for each named condition.", name, priorSeen, conditionString))
	}
	return nil
}

func (ac *Appcore) CheckNamedCondition(name string, conditionString string) (bool, error) {
	if !ac.started {
		return false, errors.New("Appcore not started")
	}
	if name == "" {
		return false, errors.New("CheckNamedCondition requires a non-empty name")
	}

	// lookup name for override, prefering the condition from the config when available
	condition := ac.config.ConditionWithName(name)

	if condition == nil {
		// Use provided condition, since config doesn't have an override
		pCond, err := datamodel.NewCondition(conditionString)
		if err != nil {
			return false, err
		}
		condition = pCond
	}

	return ac.propertyRegistry.evaluateCondition(condition)
}

func (ac *Appcore) RegisterLibraryBindings(lb LibBindings) {
	ac.libBindings = lb
}

func (ac *Appcore) Start() error {
	if ac.started {
		return errors.New("Appcore already started. Start should only be called once")
	}

	if ac.apiKey == nil {
		return errors.New("An API Key must be provided before starting critical moments")
	}
	if ac.configUrlString == "" {
		return errors.New("A config URL must be provided before starting critical moments")
	}
	if ac.libBindings == nil {
		return errors.New("The SDK must register LibBindings before calling start")
	}
	if ac.cache == nil {
		return errors.New("The SDK must register a cache directory before calling start")
	}
	if err := ac.propertyRegistry.validateProperties(); err != nil {
		return err
	}

	var configFilePath string
	var err error

	if strings.HasPrefix(ac.configUrlString, filePrefix) {
		// Strip file:// prefix
		configFilePath = ac.configUrlString[len(filePrefix):]
	} else if strings.HasPrefix(ac.configUrlString, httpsPrefix) {
		configFilePath, err = ac.cache.verifyOrFetchRemoteConfigFile(ac.configUrlString, "primary")
		if err != nil {
			return err
		}
	}
	if configFilePath == "" {
		return errors.New("CriticalMoments: Invalid config url")
	}

	configFileData, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}
	var pc datamodel.PrimaryConfig
	err = json.Unmarshal(configFileData, &pc)
	if err != nil {
		return err
	}
	ac.config = &pc
	err = ac.postConfigSetup()
	if err != nil {
		return err
	}

	ac.started = true
	return nil
}

func (ac *Appcore) postConfigSetup() error {
	if ac.config.DefaultTheme != nil {
		err := ac.libBindings.SetDefaultTheme(ac.config.DefaultTheme)
		if err != nil {
			fmt.Println("CriticalMoments: there was an issue setting up the default theme from config")
			return err
		}
	}

	return nil
}

func (ac *Appcore) SendEvent(name string) error {
	if !ac.started {
		return errors.New("Appcore not started")
	}

	event, err := datamodel.NewEventWithName(name)
	if err != nil {
		return errors.New(fmt.Sprintf("SendEvent error for \"%v\"", name))
	}

	err = ac.eventManager.SendEvent(event)
	if err != nil {
		return err
	}

	// Perform any actions for this event
	actions := ac.config.ActionsForEvent(name)
	var lastErr error
	for _, action := range actions {
		err := ac.PerformAction(&action)
		if err != nil {
			// return an error, but don't stop processing
			lastErr = errors.New(fmt.Sprintf("CriticalMoments: there was an issue performing action for event \"%v\". Error: %v\n", name, err))
		}
	}
	return lastErr
}

func (ac *Appcore) PerformNamedAction(actionName string) error {
	if !ac.started {
		return errors.New("Appcore not started")
	}
	action := ac.config.ActionWithName(actionName)
	if action == nil {
		return errors.New(fmt.Sprintf("No action found named %v", actionName))
	}
	return ac.PerformAction(action)
}

func (ac *Appcore) PerformAction(action *datamodel.ActionContainer) error {
	if !ac.started {
		return errors.New("Appcore not started")
	}
	if action.Condition != nil {
		conditionResult, err := ac.propertyRegistry.evaluateCondition(action.Condition)
		if err != nil {
			return err
		}
		if !conditionResult {
			// failing conditions are not errors
			return nil
		}
	}
	ad := actionDispatcher{
		appcore: ac,
	}
	return action.PerformAction(&ad)
}

func (ac *Appcore) ThemeForName(themeName string) *datamodel.Theme {
	if !ac.started {
		return nil
	}
	return ac.config.ThemeWithName(themeName)
}

// Repeitive, but gomobile doesn't allow for `interface{}`
func (ac *Appcore) RegisterStaticStringProperty(key string, value string) error {
	return ac.propertyRegistry.registerStaticProperty(key, value)
}
func (ac *Appcore) RegisterStaticIntProperty(key string, value int) error {
	return ac.propertyRegistry.registerStaticProperty(key, value)
}
func (ac *Appcore) RegisterStaticFloatProperty(key string, value float64) error {
	return ac.propertyRegistry.registerStaticProperty(key, value)
}
func (ac *Appcore) RegisterStaticBoolProperty(key string, value bool) error {
	return ac.propertyRegistry.registerStaticProperty(key, value)
}
func (ac *Appcore) RegisterLibPropertyProvider(key string, dpp LibPropertyProvider) error {
	return ac.propertyRegistry.registerLibPropertyProvider(key, dpp)
}
