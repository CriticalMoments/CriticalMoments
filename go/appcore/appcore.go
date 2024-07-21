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
	"github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model/conditions"
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

	// database and events
	db           *db.DB
	eventManager *EventManager

	// Properties
	propertyRegistry *propertyRegistry

	// Notifications
	notificationPlan      *NotificationPlan
	seenCancelationEvents map[string]*bool
}

func NewAppcore() *Appcore {
	ac := &Appcore{
		propertyRegistry:      newPropertyRegistry(),
		db:                    db.NewDB(),
		eventManager:          &EventManager{},
		seenCancelationEvents: make(map[string]*bool),
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

	dbOperations := ac.db.DbConditionFunctions()
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

// Internal use only
func (ac *Appcore) CheckTestCondition(conditionString string) (returnResult bool, returnErr error) {
	defer func() {
		// We never intentionally panic in CM, but we want to recover if we do
		if r := recover(); r != nil {
			returnResult = false
			returnErr = fmt.Errorf("panic in CheckTestCondition: %v", r)
		}
	}()

	appId, err := ac.propertyRegistry.propertyValue("app_id")
	if err != nil || appId == nil {
		return false, errors.New("CheckTestCondition only available in the test app. No app ID")
	}
	appIdString, ok := appId.(string)
	if !ok || (appIdString != "io.criticalmoments.demo-app" &&
		appIdString != "com.apple.dt.xctest.tool") {
		return false, errors.New("CheckTestCondition only available in the test app")
	}

	if ac.libBindings == nil || !ac.libBindings.IsTestBuild() {
		return false, errors.New("CheckTestCondition only available on a test build")
	}

	cond, err := datamodel.NewCondition(conditionString)
	if err != nil {
		return false, err
	}

	return ac.propertyRegistry.evaluateCondition(cond)
}

func (ac *Appcore) CheckNamedCondition(name string) (returnResult bool, returnErr error) {
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

	// lookup name for override, preferring the condition from the config when available
	condition := ac.config.ConditionWithName(name)

	if condition == nil {
		return false, fmt.Errorf("CheckNamedCondition: no condition found named '%v'", name)
	}

	condResult, condErr := ac.propertyRegistry.evaluateCondition(condition)
	ac.logEventForNamedCondition(condition, condResult, condErr)
	return condResult, condErr
}

func (ac *Appcore) logEventForNamedCondition(condition *datamodel.Condition, result bool, err error) {
	name := ac.config.NameForCondition(condition)
	if name == "" {
		return
	}

	if err != nil {
		ac.SendClientEvent(fmt.Sprintf("ff_error:%v", name))
	} else if result {
		ac.SendClientEvent(fmt.Sprintf("ff_true:%v", name))
	} else {
		ac.SendClientEvent(fmt.Sprintf("ff_false:%v", name))
	}
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

	err := ac.RegisterStaticTimeProperty("app_start_time", time.Now().UnixMilli())
	if err != nil {
		return err
	}
	err = ac.propertyRegistry.addProviderForKey("session_start_time", SessionStartTimePropertyProvider{eventManager: ac.eventManager})
	if err != nil {
		return err
	}

	if err := ac.propertyRegistry.validateProperties(); err != nil {
		return err
	}

	err = ac.loadConfig(allowDebugLoad)
	if err != nil {
		return err
	}

	err = ac.propertyRegistry.samplePropertiesForStartup()
	if err != nil {
		fmt.Printf("CriticalMoments: there was an issue sampling properties for startup. Continuing as this error is non-fatal: %v\n", err)
	}

	ac.started = true

	err = ac.SendBuiltInEvent(datamodel.AppStartBuiltInEvent)
	if err != nil {
		fmt.Printf("CriticalMoments: there was an issue sending the built in event \"%v\". Continuing as this error is non-fatal: %v\n", datamodel.AppStartBuiltInEvent, err)
	}

	// In practice, the event above probably already triggered this, but just in case
	err = ac.initializeNotificationPlan()
	if err != nil {
		fmt.Printf("CriticalMoments: there was an issue setting up notifications. Continuing as this error is non-fatal: %v\n", err)
	}

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
		if len(configFileData) == 2 && string(configFileData) == "{}" {
			// Special case: empty config does not require signing.
			pc = &datamodel.PrimaryConfig{
				AppId: ac.apiKey.BundleId(),
			}
		} else {
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
		}
	}
	if pc.AppId != ac.apiKey.BundleId() {
		return fmt.Errorf("this config file isn't valid for this app. Config file is key is for app id '%s', but this app has bundle ID is '%s'", pc.AppId, ac.apiKey.BundleId())
	}
	if err = ac.isClientTooOldForConfig(pc); err != nil {
		return err
	}
	ac.config = pc
	err = ac.postConfigSetup()
	if err != nil {
		return err
	}

	return nil
}

func (ac *Appcore) isClientTooOldForConfig(pc *datamodel.PrimaryConfig) error {
	if pc.MinAppVersion != "" {
		if conditions.VersionLessThan(ac.libBindings.AppVersion(), pc.MinAppVersion) {
			return fmt.Errorf("CriticalMoments: this version of the App (%v) is too old for this config file. The minimum version is %v", ac.libBindings.AppVersion(), pc.MinAppVersion)
		}
	}
	if pc.MinCMVersion != "" {
		if conditions.VersionLessThan(ac.libBindings.CMVersion(), pc.MinCMVersion) {
			return fmt.Errorf("CriticalMoments: this version of the CM SDK (%v) is too old for this config file. The minimum version is %v", ac.libBindings.CMVersion(), pc.MinCMVersion)
		}
	}
	if pc.MinCMVersionInternal != "" {
		if conditions.VersionLessThan(ac.libBindings.CMVersion(), pc.MinCMVersionInternal) {
			return fmt.Errorf("CriticalMoments: this version of the CM SDK (%v) is too old for this config file. The minimum version is %v", ac.libBindings.CMVersion(), pc.MinCMVersionInternal)
		}
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
	} else if ac.config.LibraryThemeName != "" {
		err := ac.libBindings.SetDefaultThemeByLibaryThemeName(ac.config.LibraryThemeName)
		if err != nil {
			// Non critical error. We can continue with default theme
			fmt.Println("CriticalMoments: there was an issue setting up the default library theme from config")
		}
	}

	return nil
}

func (ac *Appcore) SendClientEvent(name string) error {
	event, err := datamodel.NewClientEventWithName(name)
	if err != nil {
		return fmt.Errorf("SendEvent error for \"%v\"", name)
	}
	return ac.processEvent(event)
}

func (ac *Appcore) SendBuiltInEvent(name string) error {
	event, err := datamodel.NewBuiltInEventWithName(name)
	if err != nil {
		return fmt.Errorf("SendEvent error for \"%v\"", name)
	}
	return ac.processEvent(event)
}

func (ac *Appcore) processEvent(event *datamodel.Event) (returnErr error) {
	defer func() {
		// We never intentionally panic in CM, but we want to recover if we do
		if r := recover(); r != nil {
			returnErr = fmt.Errorf("panic in SendEvent: %v", r)
		}
	}()

	if !ac.started {
		return errors.New("Appcore not started")
	}
	if ac.eventManager == nil {
		return errors.New("Appcore EM not started")
	}

	err := ac.eventManager.SendEvent(event, ac)
	if err != nil {
		return err
	}

	performErr := ac.performActionsForEvent(event.Name)

	err = ac.notificationRunnerProcessEvent(event)
	if err != nil {
		fmt.Printf("CriticalMoments: there was an issue processing notifications for event '%v'. Error: %v\n", event.Name, err)
	}

	return performErr
}

func (ac *Appcore) performActionsForEvent(eventName string) error {
	triggers := ac.config.TriggersForEvent(eventName)
	var lastErr error
	for _, trigger := range triggers {
		if trigger.Condition != nil {
			conditionResult, err := ac.propertyRegistry.evaluateCondition(trigger.Condition)
			if err != nil {
				// return an error, but don't stop processing
				lastErr = err
				continue
			}
			if !conditionResult {
				continue
			}
		}
		err := ac.PerformNamedAction(trigger.ActionName)
		if err != nil {
			// return an error, but don't stop processing
			lastErr = fmt.Errorf("CriticalMoments: there was an issue performing action for event \"%v\". Error: %v", eventName, err)
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
	actionName := ac.config.NameForActionContainer(action)
	actionErr := action.PerformAction(&ad, actionName)
	ac.sendEventForPerformedAction(actionName, actionErr)
	return actionErr
}

func (ac *Appcore) sendEventForPerformedAction(actionName string, err error) {
	if actionName == "" {
		return
	}

	if err == nil {
		ac.SendClientEvent(fmt.Sprintf("action:%v", actionName))
	} else {
		ac.SendClientEvent(fmt.Sprintf("action_error:%v", actionName))
	}
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

func (ac *Appcore) SetLogEvents(logEvents bool) {
	ac.eventManager.logEvents = logEvents
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

func (ac *Appcore) ActionForNotification(notificationId string) error {
	for _, notification := range ac.config.Notifications {
		if notification.UniqueID() == notificationId && notification.ActionName != "" {
			return ac.PerformNamedAction(notification.ActionName)
		}
	}
	return nil
}

func (ac *Appcore) PerformBackgroundWork() error {
	return ac.performBackgroundWorkForNotifications()
}
