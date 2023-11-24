package appcore

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/CriticalMoments/CriticalMoments/go/appcore/db"
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

	// database
	db *db.DB

	// Properties
	propertyRegistry *propertyRegistry

	// Dev Mode namedCondition conflict check
	seenNamedConditions map[string]string
}

func NewAppcore() *Appcore {
	ac := &Appcore{
		propertyRegistry:    newPropertyRegistry(),
		seenNamedConditions: map[string]string{},
		db:                  db.NewDB(),
	}
	// Connect the property registry to the db/proptery history manager
	ac.propertyRegistry.phm = ac.db.PropertyHistoryManager()
	return ac
}

// Hopefully no one wants http (no TLS) in 2023... but given the importance of the config file we can't open this up to injection attacks
const filePrefix = "file://"
const httpsPrefix = "https://"

func (ac *Appcore) SetConfigUrl(configUrl string) (returnErr error) {
	defer func() {
		// We never intentionally panic in CM, but we want to recover if we do
		if r := recover(); r != nil {
			returnErr = fmt.Errorf("panic in SetConfigUrl: %v", r)
		}
	}()

	if !strings.HasPrefix(configUrl, filePrefix) && !strings.HasPrefix(configUrl, httpsPrefix) {
		return errors.New("config URL must start with https:// or file://")
	}
	ac.configUrlString = configUrl

	return nil
}

func (ac *Appcore) ConfigUrl() string {
	return ac.configUrlString
}

func (ac *Appcore) SetApiKey(apiKey string, bundleID string) (returnErr error) {
	defer func() {
		// We never intentionally panic in CM, but we want to recover if we do
		if r := recover(); r != nil {
			returnErr = fmt.Errorf("panic in SetApiKey: %v", r)
		}
	}()

	key, err := signing.ParseApiKey(apiKey)
	if err != nil {
		return errors.New("invalid API Key. Please make sure you get your key from criticalmoments.io")
	}
	if v, err := key.Valid(); err != nil || !v {
		return errors.New("invalid API Key. Please make sure you get your key from criticalmoments.io")
	}
	if key.BundleId() != bundleID {
		return fmt.Errorf("this API key isn't valid for this app. API key is for %s, but this app has bundle ID %s", key.BundleId(), bundleID)
	}
	ac.apiKey = key
	return nil
}

func (ac *Appcore) ApiKey() string {
	if ac.apiKey == nil {
		return ""
	}
	return ac.apiKey.String()
}

func (ac *Appcore) SetDataDirPath(dataDirPath string) (returnErr error) {
	defer func() {
		// We never intentionally panic in CM, but we want to recover if we do
		if r := recover(); r != nil {
			returnErr = fmt.Errorf("panic in SetDataDirPath: %v", r)
		}
	}()

	cache, err := newCacheWithBaseDir(dataDirPath)
	if err != nil {
		return err
	}
	ac.cache = cache

	err = ac.db.StartWithPath(dataDirPath)
	if err != nil {
		return err
	}

	dbOperations := ac.db.EventManager().EventManagerConditionFunctions()
	ac.propertyRegistry.RegisterDynamicFunctions(dbOperations)

	return nil
}

func (ac *Appcore) SetTimezoneGMTOffset(gmtOffset int) {
	defer func() {
		// We never intentionally panic in CM, but we want to recover if we do
		if r := recover(); r != nil {
			fmt.Printf("CriticalMoments: panic in SetTimezoneGMTOffset: %v\n", r)
		}
	}()

	tzName := fmt.Sprintf("UTCOffsetS:%v", gmtOffset)
	tz := time.FixedZone(tzName, gmtOffset)
	time.Local = tz

	ac.propertyRegistry.registerStaticProperty("timezone_gmt_offset", gmtOffset)
}

func (ac *Appcore) CheckNamedConditionCollision(name string, conditionString string) (returnErr error) {
	defer func() {
		// We never intentionally panic in CM, but we want to recover if we do
		if r := recover(); r != nil {
			returnErr = fmt.Errorf("panic in CheckNamedConditionsCollision: %v", r)
		}
	}()

	if name == "" {
		return nil
	}
	// in debug mode, track each built-in condition we see, and make sure the developer isn't reusing names
	// If they use the same name twice for different things, they won't be able to override in the future
	priorSeen := ac.seenNamedConditions[name]
	if priorSeen == "" {
		ac.seenNamedConditions[name] = conditionString
	} else if priorSeen != conditionString {
		return fmt.Errorf("the named condition \"%v\" is being used in multiple places in this codebase, with different fallback conditions (\"%v\" and \"%v\"). This will make it impossible to override each usage independently from remote configuration. Please use unique names for each named condition", name, priorSeen, conditionString)
	}
	return nil
}

func (ac *Appcore) CheckNamedCondition(name string, conditionString string) (returnResult bool, returnErr error) {
	defer func() {
		// We never intentionally panic in CM, but we want to recover if we do
		if r := recover(); r != nil {
			returnResult = false
			returnErr = fmt.Errorf("panic in CheckNamedCondition: %v", r)
		}
	}()

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
	defer func() {
		// We never intentionally panic in CM, but we want to recover if we do
		if r := recover(); r != nil {
			fmt.Printf("CriticalMoments: panic in RegisterLibrayBindings: %v\n", r)
		}
	}()

	ac.libBindings = lb

	// connect iOS functions to condition system
	ac.propertyRegistry.RegisterDynamicFunctions(map[string]*datamodel.ConditionDynamicFunction{
		"canOpenUrl": {
			Function: func(params ...any) (any, error) {
				// Parameter type+count checking is done the Types signature
				return lb.CanOpenURL(params[0].(string)), nil
			},
			Types: []any{new(func(string) bool)},
		},
	})
}

func (ac *Appcore) Start(allowDebugLoad bool) (returnErr error) {
	defer func() {
		// We never intentionally panic in CM, but we want to recover if we do
		if r := recover(); r != nil {
			returnErr = fmt.Errorf("panic in Start: %v", r)
		}
	}()

	if ac.started {
		return errors.New("appcore already started. Start should only be called once")
	}

	if ac.apiKey == nil {
		return errors.New("an API Key must be provided before starting critical moments")
	}
	if ac.configUrlString == "" {
		return errors.New("a config URL must be provided before starting critical moments")
	}
	if ac.libBindings == nil {
		return errors.New("the SDK must register LibBindings before calling start")
	}
	if ac.cache == nil {
		return errors.New("the SDK must register a cache directory before calling start")
	}
	if err := ac.propertyRegistry.validateProperties(); err != nil {
		return err
	}

	err := ac.loadConfig(allowDebugLoad)
	if err != nil {
		return err
	}

	err = ac.propertyRegistry.samplePropertiesForStartup()
	if err != nil {
		fmt.Printf("CriticalMoments: there was an issue sampling properties for startup. Continuing as this error is non-fatal: %v\n", err)
	}

	ac.started = true
	return nil
}

func (ac *Appcore) loadConfig(allowDebugLoad bool) error {
	var configFilePath string
	var err error
	isFilePath := strings.HasPrefix(ac.configUrlString, filePrefix)

	if isFilePath {
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

	pc, err := datamodel.DecodePrimaryConfig(configFileData, signing.SharedSignUtil())
	if err != nil {
		// If we're in debug mode and the file is local, allow parsing unsigned config files
		allowParsingUnsigned := allowDebugLoad && isFilePath
		if !allowParsingUnsigned {
			return err
		}
		pc = &datamodel.PrimaryConfig{}
		err = json.Unmarshal(configFileData, &pc)
		if err != nil {
			return err
		}
		fmt.Println("CriticalMoments: DEVELOPMENT MODE. Loaded an unsigned config file. This is allowed for local development, but don't forget to sign your config file before releasing to app store.")
	}
	if pc.AppId != ac.apiKey.BundleId() {
		return fmt.Errorf("this config file isn't valid for this app. Config file is key is for app id '%s', but this app has bundle ID is '%s'", pc.AppId, ac.apiKey.BundleId())
	}
	ac.config = pc
	err = ac.postConfigSetup()
	if err != nil {
		return err
	}

	return nil
}

func (ac *Appcore) postConfigSetup() error {
	dt := ac.config.DefaultTheme()
	if dt != nil {
		err := ac.libBindings.SetDefaultTheme(dt)
		if err != nil {
			fmt.Println("CriticalMoments: there was an issue setting up the default theme from config")
			return err
		}
	}

	return nil
}

func (ac *Appcore) SendEvent(name string) (returnErr error) {
	defer func() {
		// We never intentionally panic in CM, but we want to recover if we do
		if r := recover(); r != nil {
			returnErr = fmt.Errorf("panic in SendEvent: %v", r)
		}
	}()

	if !ac.started {
		return errors.New("Appcore not started")
	}

	event, err := datamodel.NewEventWithName(name)
	if err != nil {
		return fmt.Errorf("SendEvent error for \"%v\"", name)
	}

	err = ac.db.EventManager().SendEvent(event)
	if err != nil {
		return err
	}

	// Perform any actions for this event
	actions := ac.config.ActionsForEvent(name)
	var lastErr error
	for _, action := range actions {
		err := ac.PerformAction(action)
		if err != nil {
			// return an error, but don't stop processing
			lastErr = fmt.Errorf("CriticalMoments: there was an issue performing action for event \"%v\". Error: %v", name, err)
		}
	}
	return lastErr
}

func (ac *Appcore) PerformNamedAction(actionName string) (returnErr error) {
	defer func() {
		// We never intentionally panic in CM, but we want to recover if we do
		if r := recover(); r != nil {
			returnErr = fmt.Errorf("panic in PerformNamedAction: %v", r)
		}
	}()

	if !ac.started {
		return errors.New("Appcore not started")
	}
	action := ac.config.ActionWithName(actionName)
	if action == nil {
		return fmt.Errorf("no action found named %v", actionName)
	}
	return ac.PerformAction(action)
}

func (ac *Appcore) PerformAction(action *datamodel.ActionContainer) (returnErr error) {
	defer func() {
		// We never intentionally panic in CM, but we want to recover if we do
		if r := recover(); r != nil {
			returnErr = fmt.Errorf("panic in PerformAction: %v", r)
		}
	}()

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

func (ac *Appcore) ThemeForName(themeName string) (resultTheme *datamodel.Theme) {
	defer func() {
		// We never intentionally panic in CM, but we want to recover if we do
		if r := recover(); r != nil {
			fmt.Printf("CriticalMoments: panic in ThemeForName: %v", r)
			resultTheme = nil
		}
	}()

	if !ac.started {
		return nil
	}
	return ac.config.ThemeWithName(themeName)
}

var errRegisterAfterStart = errors.New("Appcore already started. Properties must be registered before starting")

// Repeitive, but gomobile doesn't allow for `interface{}`
// Panic catching is one level down stack here, but still there.
func (ac *Appcore) RegisterStaticStringProperty(key string, value string) error {
	if ac.started {
		return errRegisterAfterStart
	}
	return ac.propertyRegistry.registerStaticProperty(key, value)
}
func (ac *Appcore) RegisterStaticIntProperty(key string, value int) error {
	if ac.started {
		return errRegisterAfterStart
	}
	return ac.propertyRegistry.registerStaticProperty(key, value)
}
func (ac *Appcore) RegisterStaticFloatProperty(key string, value float64) error {
	if ac.started {
		return errRegisterAfterStart
	}
	return ac.propertyRegistry.registerStaticProperty(key, value)
}
func (ac *Appcore) RegisterStaticBoolProperty(key string, value bool) error {
	if ac.started {
		return errRegisterAfterStart
	}
	return ac.propertyRegistry.registerStaticProperty(key, value)
}
func (ac *Appcore) RegisterStaticTimeProperty(key string, value int64) error {
	if ac.started {
		return errRegisterAfterStart
	}
	if value == LibPropertyProviderNilIntValue {
		return ac.propertyRegistry.registerStaticProperty(key, nil)
	}
	timeVal := time.UnixMilli(value)
	return ac.propertyRegistry.registerStaticProperty(key, timeVal)
}
func (ac *Appcore) RegisterClientStringProperty(key string, value string) error {
	if ac.started {
		return errRegisterAfterStart
	}
	return ac.propertyRegistry.registerClientProperty(key, value)
}
func (ac *Appcore) RegisterClientIntProperty(key string, value int) error {
	if ac.started {
		return errRegisterAfterStart
	}
	return ac.propertyRegistry.registerClientProperty(key, value)
}
func (ac *Appcore) RegisterClientFloatProperty(key string, value float64) error {
	if ac.started {
		return errRegisterAfterStart
	}
	return ac.propertyRegistry.registerClientProperty(key, value)
}
func (ac *Appcore) RegisterClientBoolProperty(key string, value bool) error {
	if ac.started {
		return errRegisterAfterStart
	}
	return ac.propertyRegistry.registerClientProperty(key, value)
}
func (ac *Appcore) RegisterClientTimeProperty(key string, value int64) error {
	if ac.started {
		return errRegisterAfterStart
	}
	if value == LibPropertyProviderNilIntValue {
		return ac.propertyRegistry.registerClientProperty(key, nil)
	}
	timeVal := time.UnixMilli(value)
	return ac.propertyRegistry.registerClientProperty(key, timeVal)
}

func (ac *Appcore) RegisterLibPropertyProvider(key string, dpp LibPropertyProvider) (returnErr error) {
	defer func() {
		// We never intentionally panic in CM, but we want to recover if we do
		if r := recover(); r != nil {
			returnErr = fmt.Errorf("panic in RegisterLibPropertyProvider: %v", r)
		}
	}()

	if ac.started {
		return errRegisterAfterStart
	}
	return ac.propertyRegistry.registerLibPropertyProvider(key, dpp)
}
func (ac *Appcore) RegisterClientPropertiesFromJson(jsonData []byte) (returnErr error) {
	defer func() {
		// We never intentionally panic in CM, but we want to recover if we do
		if r := recover(); r != nil {
			returnErr = fmt.Errorf("panic in RegisterClientPropertiesFromJson: %v", r)
		}
	}()

	if ac.started {
		return errRegisterAfterStart
	}
	return ac.propertyRegistry.registerClientPropertiesFromJson(jsonData)
}
