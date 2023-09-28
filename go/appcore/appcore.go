package appcore

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/CriticalMoments/CriticalMoments/go/cmcore"
	datamodel "github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model"
	"github.com/CriticalMoments/CriticalMoments/go/cmcore/signing"
)

func GoPing() string {
	return "AppcorePong->" + cmcore.CmCorePing()
}

type Appcore struct {
	// Library binding/delegate
	libBindings LibBindings

	// API Key
	apiKey *signing.ApiKey

	// Primary configuration
	configUrlString string
	config          *datamodel.PrimaryConfig

	// Cache
	cache *cache

	// Properties
	propertyRegistry *propertyRegistry
}

var sharedAppcore Appcore = newAppcore()

func SharedAppcore() *Appcore {
	return &sharedAppcore
}
func newAppcore() Appcore {
	return Appcore{
		propertyRegistry: newPropertyRegistry(),
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

func (ac *Appcore) SetCacheDirPath(cacheDirPath string) error {
	cache, err := newCacheWithBaseDir(cacheDirPath)
	if err != nil {
		return err
	}

	ac.cache = cache
	return nil
}

func (ac *Appcore) SetTimezoneGMTOffset(gmtOffset int) {
	tzName := fmt.Sprintf("UTCOffsetS:%v", gmtOffset)
	tz := time.FixedZone(tzName, gmtOffset)
	time.Local = tz
}

func (ac *Appcore) RegisterLibraryBindings(lb LibBindings) {
	ac.libBindings = lb
}

// TODO: guard against double start call
func (ac *Appcore) Start() error {
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

// TODO: events should be queued during setup, and run after postConfigSetup
func (ac *Appcore) SendEvent(e string) error {
	actions := ac.config.ActionsForEvent(e)
	if len(actions) == 0 {
		return errors.New(fmt.Sprintf("Event not found: %v", e))
	}
	var lastErr error
	for _, action := range actions {
		err := ac.PerformAction(&action)
		if err != nil {
			// return an error, but don't stop sending
			lastErr = errors.New(fmt.Sprintf("CriticalMoments: there was an issue performing action for event \"%v\". Error: %v\n", e, err))
		}
	}
	return lastErr
}

func (ac *Appcore) PerformNamedAction(actionName string) error {
	action := ac.config.ActionWithName(actionName)
	if action == nil {
		return errors.New(fmt.Sprintf("No action found named %v", actionName))
	}
	return ac.PerformAction(action)
}

func (ac *Appcore) PerformAction(action *datamodel.ActionContainer) error {
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
