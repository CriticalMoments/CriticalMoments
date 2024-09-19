package appcore

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	datamodel "github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"golang.org/x/exp/maps"
)

func TestPing(t *testing.T) {
	pingResponse := GoPing()
	if pingResponse != "AppcorePong->PongCmCore" {
		t.Fatalf("appcore ping failure: %v", pingResponse)
	}
}

func TestUrlValidation(t *testing.T) {
	ac := &Appcore{}
	err := ac.SetConfigUrl("http://asdf.com")
	if err == nil {
		t.Fatal("Allowed http (no s) url")
	}
	err = ac.SetConfigUrl("ftp://192.168.99.99")
	if err == nil {
		t.Fatal("Allowed invalid url")
	}
	err = ac.SetConfigUrl("https://asdf.com")
	if err != nil {
		t.Fatal("Disallowed valid url")
	}
	err = ac.SetConfigUrl("file://Users/criticalmoments/config.json")
	if err != nil {
		t.Fatal("Disallowed valid url")
	}
}

func buildTestBuiltInProps(propTypes map[string]*datamodel.CMPropertyConfig) map[string]*datamodel.CMPropertyConfig {
	// These are populated by Appcore, so should be included even if not used in test
	props := map[string]*datamodel.CMPropertyConfig{
		"app_start_time":     {Type: datamodel.CMTimeKind, Source: datamodel.CMPropertySourceLib, Optional: false, SampleType: datamodel.CMPropertySampleTypeDoNotSample},
		"session_start_time": {Type: datamodel.CMTimeKind, Source: datamodel.CMPropertySourceLib, Optional: false, SampleType: datamodel.CMPropertySampleTypeDoNotSample},
		"is_debug_build":     {Type: reflect.Bool, Source: datamodel.CMPropertySourceLib, Optional: false, SampleType: datamodel.CMPropertySampleTypeDoNotSample},
	}

	for k, v := range propTypes {
		props[k] = v
	}

	return props
}

type testLibBindings struct {
	lastBannerAction     *datamodel.BannerAction
	lastAlertAction      *datamodel.AlertAction
	lastLinkAction       *datamodel.LinkAction
	reviewCount          int
	defaultTheme         *datamodel.Theme
	libThemeName         string
	lastModal            *datamodel.ModalAction
	lastNotificationPlan *NotificationPlan
}

func (lb *testLibBindings) ShowBanner(b *datamodel.BannerAction, actionName string) error {
	lb.lastBannerAction = b
	return nil
}
func (lb *testLibBindings) ShowAlert(a *datamodel.AlertAction, actionName string) error {
	lb.lastAlertAction = a
	return nil
}
func (lb *testLibBindings) ShowLink(l *datamodel.LinkAction) error {
	lb.lastLinkAction = l
	return nil
}
func (lb *testLibBindings) SetDefaultTheme(theme *datamodel.Theme) error {
	lb.defaultTheme = theme
	return nil
}
func (lb *testLibBindings) SetDefaultThemeByLibaryThemeName(name string) error {
	lb.libThemeName = name
	return nil
}
func (lb *testLibBindings) ShowReviewPrompt() error {
	lb.reviewCount += 1
	return nil
}
func (lb *testLibBindings) ShowModal(modal *datamodel.ModalAction, actionName string) error {
	lb.lastModal = modal
	return nil
}
func (lb *testLibBindings) CanOpenURL(url string) bool {
	return false
}
func (lb *testLibBindings) AppVersion() string {
	return "1.2.3"
}
func (lb *testLibBindings) CMVersion() string {
	return "2.3.4"
}
func (lb *testLibBindings) IsTestBuild() bool {
	return true
}
func (lb *testLibBindings) UpdateNotificationPlan(notifPlan *NotificationPlan) error {
	lb.lastNotificationPlan = notifPlan
	return nil
}

func testBuildValidTestAppCore(t *testing.T) (*Appcore, error) {
	return buildTestAppCoreWithPath("../cmcore/data_model/test/testdata/primary_config/valid/maximalValid.json", t)
}

func buildTestAppCoreWithPath(path string, t *testing.T) (*Appcore, error) {
	ac := NewAppcore()
	configPath, err := filepath.Abs(path)
	if err != nil {
		t.Fatal(err)
	}
	configUrl := fmt.Sprintf("file://%v", configPath)
	err = ac.SetConfigUrl(configUrl)
	if err != nil {
		t.Fatal(err)
	}
	baseDataPath := fmt.Sprintf("/tmp/criticalmoments/test-temp-%v", rand.Int())
	os.MkdirAll(baseDataPath, os.ModePerm)
	err = ac.SetDataDirPath(baseDataPath)
	if err != nil {
		t.Fatal(err)
	}
	if ac.db == nil || ac.eventManager == nil || ac.db.PropertyHistoryManager() == nil || ac.cache == nil {
		t.Fatal("db, event manager, prop history manager, or cache not setup")
	}
	if ac.propertyRegistry.phm != ac.db.PropertyHistoryManager() {
		t.Fatal("property history manager not set to the correct DB instance via NewAppcore")
	}
	lb := testLibBindings{}
	ac.RegisterLibraryBindings(&lb)

	ac.SetApiKey("CM1-aGVsbG86d29ybGQ=-Yjppby5jcml0aWNhbG1vbWVudHMuZGVtbw==-MEUCIQCUfx6xlmQ0kdYkuw3SMFFI6WXrCWKWwetXBrXXG2hjAwIgWBPIMrdM1ET0HbpnXlnpj/f+VXtjRTqNNz9L/AOt4GY=", "io.criticalmoments.demo")

	// Clear required properties, for easier setup
	ac.propertyRegistry.builtInPropertyTypes = buildTestBuiltInProps(map[string]*datamodel.CMPropertyConfig{})
	return ac, nil
}

func TestAppcoreStart(t *testing.T) {
	ac, err := testBuildValidTestAppCore(t)
	if err != nil {
		t.Fatal(err)
	}
	err = ac.Start(true)
	if err != nil {
		t.Fatal(err)
	}
	// Check it loaded the config (more detailed test of parsing in cmcore)
	if ac.config.DefaultTheme() == nil {
		t.Fatal("Failed to load config in Appcore setup")
	}
}

func TestAppcoreStartMissingConfig(t *testing.T) {
	ac, err := testBuildValidTestAppCore(t)
	if err != nil {
		t.Fatal(err)
	}
	ac.configUrlString = ""
	err = ac.Start(true)
	if err == nil {
		t.Fatal("Should not start without config")
	}
	if ac.config != nil {
		t.Fatal("Loaded config from empty url")
	}
}

func TestAppcoreStartMissingBindings(t *testing.T) {
	ac, err := testBuildValidTestAppCore(t)
	if err != nil {
		t.Fatal(err)
	}
	ac.libBindings = nil
	err = ac.Start(true)
	if err == nil {
		t.Fatal("Should not start without bindings")
	}
	if ac.config != nil {
		t.Fatal("Loaded config without bindings")
	}
}

func TestAppcoreStartBadConfig(t *testing.T) {
	ac, err := testBuildValidTestAppCore(t)
	if err != nil {
		t.Fatal(err)
	}
	ac.configUrlString = "file:///Not/A/Real/Path"
	err = ac.Start(true)
	if err == nil {
		t.Fatal("Should not start with bad config")
	}
	if ac.config != nil {
		t.Fatal("Loaded config from bad url")
	}
}

func TestSendEvent(t *testing.T) {
	ac, err := testBuildValidTestAppCore(t)
	if err != nil {
		t.Fatal(err)
	}
	err = ac.Start(true)
	if err != nil {
		t.Fatal(err)
	}

	// built in should error through client API
	err = ac.SendClientEvent(datamodel.AppStartBuiltInEvent)
	if err == nil {
		t.Fatal("invalid build in event did not error")
	}

	// built in should work thorough internal API
	err = ac.SendBuiltInEvent(datamodel.AppStartBuiltInEvent)
	if err != nil {
		t.Fatal("valid build in event errored", err)
	}

	// well known should fail thorough built in API
	err = ac.SendBuiltInEvent(datamodel.SignedInEvent)
	if err == nil {
		t.Fatal("invalid well known event did not error")
	}

	// Well known should work though client
	err = ac.SendClientEvent(datamodel.SignedInEvent)
	if err != nil {
		t.Fatal(err)
	}

	// custom events with no actions should work
	err = ac.SendClientEvent("net.scosman.asdf")
	if err != nil {
		t.Fatal("valid custom event errored", err)
	}
}

func TestNotificationPlanCallback(t *testing.T) {
	ac, err := testBuildValidTestAppCore(t)
	if err != nil {
		t.Fatal(err)
	}

	lb := testLibBindings{}
	ac.RegisterLibraryBindings(&lb)
	if lb.lastNotificationPlan != nil {
		t.Fatal("NP binding set too soon")
	}

	err = ac.Start(true)
	if err != nil {
		t.Fatal(err)
	}

	if lb.lastNotificationPlan == nil {
		t.Fatal("didn't call notification plan callback")
	}
}

func TestPerformingAction(t *testing.T) {
	ac, err := testBuildValidTestAppCore(t)
	if err != nil {
		t.Fatal(err)
	}
	err = ac.Start(true)
	if err != nil {
		t.Fatal(err)
	}
	if ac.libBindings.(*testLibBindings).lastBannerAction != nil {
		t.Fatal("last banner action should be nil on new appcore test binding")
	}
	// should fire bannerAction1 via a trigger
	err = ac.SendClientEvent("custom_event")
	if err != nil {
		t.Fatal(err)
	}
	if ac.libBindings.(*testLibBindings).lastBannerAction.Body != "Hello world, but on a banner!" {
		t.Fatal("last banner action should be nil on new appcore test binding")
	}

	latestAlertTime, err := ac.db.LatestEventTimeByName("action:alertAction")
	if err != nil {
		t.Fatal(err)
	}
	if latestAlertTime != nil {
		t.Fatal("Alert action should not have fired yet")
	}
	if ac.libBindings.(*testLibBindings).lastAlertAction != nil {
		t.Fatal("last alert action should be nil on new appcore test binding")
	}
	// condition should stop it from firing
	err = ac.PerformNamedAction("alertActionWithFailingCondition")
	if err != nil {
		// Specifically, no not found error
		t.Fatal(err)
	}
	if ac.libBindings.(*testLibBindings).lastAlertAction != nil {
		t.Fatal("event fired when condition false")
	}
	err = ac.PerformNamedAction("alertAction")
	if err != nil {
		t.Fatal(err)
	}
	if ac.libBindings.(*testLibBindings).lastAlertAction == nil {
		t.Fatal("alert didn't fire")
	}
	latestAlertTime, err = ac.db.LatestEventTimeByName("action:alertAction")
	if err != nil {
		t.Fatal(err)
	}
	if latestAlertTime == nil {
		t.Fatal("Alert action should have fired now")
	}

	err = ac.PerformNamedAction("reviewAction")
	if err != nil {
		t.Fatal(err)
	}
	if ac.libBindings.(*testLibBindings).reviewCount != 1 {
		t.Fatal("review action didn't fire")
	}

	if ac.libBindings.(*testLibBindings).lastModal != nil {
		t.Fatal("modal fired too soon")
	}
	err = ac.PerformNamedAction("modalAction")
	if err != nil {
		t.Fatal(err)
	}
	if ac.libBindings.(*testLibBindings).lastModal == nil {
		t.Fatal("modal event didn't fire")
	}

	err = ac.PerformNamedAction("unknownActionTypeFutureProof")
	if err == nil || !strings.Contains(err.Error(), "does not support this action type") {
		t.Fatal("Unknown action didn't error")
	}

	err = ac.PerformNamedAction("nestedFutureTypeFail")
	if err == nil || !strings.Contains(err.Error(), "does not support this action type") {
		t.Fatal("Nested unknown actions didn't error up the stack")
	}

	// Verify fallback from future to alert, both single level and deep nested
	fallbackActions := []string{"futureAction", "nestedFutureTypeSuccess"} // add_test_count
	for _, actionName := range fallbackActions {
		ac.libBindings.(*testLibBindings).lastAlertAction = nil
		if ac.libBindings.(*testLibBindings).lastAlertAction != nil {
			t.Fatal("test not initialized")
		}
		err = ac.PerformNamedAction(actionName)
		if err != nil {
			t.Fatal(err)
		}
		if ac.libBindings.(*testLibBindings).lastAlertAction == nil {
			t.Fatal("alert event didn't fire as fallback")
		}
	}
}

func TestConditionalActionAndTriggerDispatching(t *testing.T) {
	ac, err := testBuildValidTestAppCore(t)
	if err != nil {
		t.Fatal(err)
	}
	err = ac.Start(true)
	if err != nil {
		t.Fatal(err)
	}
	if ac.libBindings.(*testLibBindings).lastBannerAction != nil {
		t.Fatal("last action should be nil on new appcore test binding")
	}
	if ac.libBindings.(*testLibBindings).lastAlertAction != nil {
		t.Fatal("last action should be nil on new appcore test binding")
	}
	if ac.libBindings.(*testLibBindings).lastLinkAction != nil {
		t.Fatal("last action should be nil on new appcore test binding")
	}
	err = ac.PerformNamedAction("conditionalWithTrueCondition")
	if err != nil {
		t.Fatal(err)
	}
	if ac.libBindings.(*testLibBindings).lastBannerAction != nil {
		t.Fatal("last action should be nil after condition run 1")
	}
	if ac.libBindings.(*testLibBindings).lastAlertAction == nil {
		t.Fatal("last alert action should not be nil after condiiton run 1")
	}
	if ac.libBindings.(*testLibBindings).lastLinkAction != nil {
		t.Fatal("last action should be nil after condition run 1")
	}
	err = ac.PerformNamedAction("conditionalWithFalseCondition")
	if err != nil {
		t.Fatal(err)
	}
	if ac.libBindings.(*testLibBindings).lastBannerAction != nil {
		t.Fatal("last action should be nil after condition run 2")
	}
	if ac.libBindings.(*testLibBindings).lastAlertAction == nil {
		t.Fatal("last alert action should not be nil after condiiton run 2")
	}
	if ac.libBindings.(*testLibBindings).lastLinkAction == nil {
		t.Fatal("last action should not be nil after condition run 2")
	}
	err = ac.PerformNamedAction("conditionalWithoutFalseAction")
	if err != nil {
		t.Fatal(err)
	}
	if ac.libBindings.(*testLibBindings).lastBannerAction != nil {
		t.Fatal("last action should be nil after condition run 3")
	}
	if ac.libBindings.(*testLibBindings).lastAlertAction == nil {
		t.Fatal("last alert action should not be nil after condiiton run 3")
	}
	if ac.libBindings.(*testLibBindings).lastLinkAction == nil {
		t.Fatal("last action should not be nil after condition run 3")
	}

	ac.libBindings.(*testLibBindings).lastAlertAction = nil
	ac.libBindings.(*testLibBindings).lastLinkAction = nil
	err = ac.SendClientEvent("custom_event_conditional_false")
	if err != nil {
		t.Fatal(err)
	}
	if ac.libBindings.(*testLibBindings).lastBannerAction != nil {
		t.Fatal("last action should be nil when condition false")
	}
	if ac.libBindings.(*testLibBindings).lastAlertAction != nil {
		t.Fatal("last action should be nil when condition false")
	}
	if ac.libBindings.(*testLibBindings).lastLinkAction != nil {
		t.Fatal("last action should be nil when condition false")
	}

	err = ac.SendClientEvent("custom_event_conditional_true")
	if err != nil {
		t.Fatal(err)
	}
	if ac.libBindings.(*testLibBindings).lastBannerAction != nil {
		t.Fatal("last action should be nil after trigger")
	}
	if ac.libBindings.(*testLibBindings).lastAlertAction == nil {
		t.Fatal("last alert action should be fired when condition true")
	}
	if ac.libBindings.(*testLibBindings).lastLinkAction != nil {
		t.Fatal("last action should be nil after trigger")
	}
}

func TestSetDefaultTheme(t *testing.T) {
	ac, err := testBuildValidTestAppCore(t)
	if err != nil {
		t.Fatal(err)
	}
	if ac.libBindings.(*testLibBindings).defaultTheme != nil {
		t.Fatal("Theme should be nil until started")
	}
	err = ac.Start(true)
	if err != nil {
		t.Fatal(err)
	}
	defaultTheme := ac.libBindings.(*testLibBindings).defaultTheme
	if defaultTheme == nil || defaultTheme.BannerBackgroundColor != "#ffffff" {
		t.Fatal("Default theme not set after start")
	}
	if ac.libBindings.(*testLibBindings).libThemeName != "" {
		t.Fatal("Default theme set after start")
	}
}

func TestSetDefaultLibraryTheme(t *testing.T) {
	ac, err := buildTestAppCoreWithPath("../cmcore/data_model/test/testdata/primary_config/valid/builtInLibraryTheme.json", t)
	if err != nil {
		t.Fatal(err)
	}
	if ac.libBindings.(*testLibBindings).defaultTheme != nil {
		t.Fatal("Theme should be nil until started")
	}
	err = ac.Start(true)
	if err != nil {
		t.Fatal(err)
	}
	defaultThemeName := ac.libBindings.(*testLibBindings).libThemeName
	if defaultThemeName != "system_dark" {
		t.Fatal("Default theme not set after start")
	}
	defaultTheme := ac.libBindings.(*testLibBindings).defaultTheme
	if defaultTheme != nil {
		t.Fatal("Default theme set after start")
	}
}

func TestNamedConditions(t *testing.T) {
	ac, err := testBuildValidTestAppCore(t)
	if err != nil {
		t.Fatal(err)
	}
	err = ac.Start(true)
	if err != nil {
		t.Fatal(err)
	}

	// conditions without config file entry should return false
	r, err := ac.CheckNamedCondition("conditionNotInConfig")
	if err == nil || r {
		t.Fatal("missing named condition should error and return false")
	}

	falseConditionEvent := "ff_false:falseCondition"
	latestFalseTime, err := ac.db.LatestEventTimeByName(falseConditionEvent)
	if err != nil {
		t.Fatal(err)
	}
	if latestFalseTime != nil {
		t.Fatal("false condition should not have fired yet")
	}
	// falseCondition should return false
	r, err = ac.CheckNamedCondition("falseCondition")
	if err != nil || r {
		t.Fatal("false conditions failed")
	}
	latestFalseTime, err = ac.db.LatestEventTimeByName(falseConditionEvent)
	if err != nil {
		t.Fatal(err)
	}
	if latestFalseTime == nil {
		t.Fatal("false condition should have fired")
	}

	trueConditionEvent := "ff_true:trueCondition"
	latestTrueTime, err := ac.db.LatestEventTimeByName(trueConditionEvent)
	if err != nil {
		t.Fatal(err)
	}
	if latestTrueTime != nil {
		t.Fatal("true condition should not have fired yet")
	}
	// trueCondition should return true
	r, err = ac.CheckNamedCondition("trueCondition")
	if err != nil || !r {
		t.Fatal("false conditions failed")
	}
	latestTrueTime, err = ac.db.LatestEventTimeByName(trueConditionEvent)
	if err != nil {
		t.Fatal(err)
	}
	if latestTrueTime == nil {
		t.Fatal("true condition should have fired")
	}

	// Check name check
	_, err = ac.CheckNamedCondition("")
	if err == nil {
		t.Fatal("CheckNamedCondition requires name and didn't validate empty string")
	}
}

func TestInternalTestConditions(t *testing.T) {
	// ensure CheckTestCondition is false+error for incorrect bundle ID.
	// Other tests done in objc layer
	ac, err := testBuildValidTestAppCore(t)
	if err != nil {
		t.Fatal(err)
	}
	err = ac.Start(true)
	if err != nil {
		t.Fatal(err)
	}
	r, err := ac.CheckTestCondition("true")
	if err == nil || r {
		t.Fatal("Test condition failed to block invalid test bundle")
	}
}

func TestEndToEndEvents(t *testing.T) {
	ac, err := testBuildValidTestAppCore(t)
	if err != nil {
		t.Fatal(err)
	}
	err = ac.Start(true)
	if err != nil {
		t.Fatal(err)
	}

	c, err := datamodel.NewCondition("eventCount('test') == 0")
	if err != nil {
		t.Fatal(err)
	}
	r, err := ac.propertyRegistry.evaluateCondition(c)
	if err != nil {
		t.Fatal(err)
	}
	if !r {
		t.Fatal("eventCount should be 0")
	}

	ac.SendClientEvent("test")
	ac.SendClientEvent("test")
	ac.SendClientEvent("test2")

	c, err = datamodel.NewCondition("eventCount('test') == 2 && eventCount('test2') == 1")
	if err != nil {
		t.Fatal(err)
	}
	r, err = ac.propertyRegistry.evaluateCondition(c)
	if err != nil {
		t.Fatal(err)
	}
	if !r {
		t.Fatal("eventCount should be 3 (2 and 1)")
	}

	c, err = datamodel.NewCondition("eventCountWithLimit('test', 1) == 1")
	if err != nil {
		t.Fatal(err)
	}
	r, err = ac.propertyRegistry.evaluateCondition(c)
	if err != nil {
		t.Fatal(err)
	}
	if !r {
		t.Fatal("eventCountWithLimit should limit")
	}

	// Confirm we are checking signature
	invalidParams := []string{
		"eventCount() > 1",
		"eventCount('test', 'test2') > 1",
		"eventCount(1) > 1",
		"eventCount('test', 1) > 1",
		"eventCountWithLimit('test') > 1",
		"eventCountWithLimit('test', 'test2') > 1",
		"eventCountWithLimit() > 1",
	}
	for _, cs := range invalidParams {
		c, err = datamodel.NewCondition(cs)
		if err != nil {
			t.Fatal(err)
		}
		_, err = ac.propertyRegistry.evaluateCondition(c)
		if err == nil {
			t.Fatal("Allowed condition with invalid parameters", cs)
		}
	}
}

func arraysEqualOrderInsensitive(a []string, b []string) bool {
	less := func(aa, bb string) bool { return aa < bb }
	return cmp.Diff(a, b, cmpopts.SortSlices(less)) == ""
}

func TestValidateAllBuiltInFunctionsAreRegistered(t *testing.T) {
	// Verify on startup, all the functions we expect to support in cmcore.AllBuiltInDynamicFunctions
	// are actually registered after .start. If not, strict validation is either not strict enough
	// or too strict!

	ac, err := testBuildValidTestAppCore(t)
	if err != nil {
		t.Fatal(err)
	}
	err = ac.Start(true)
	if err != nil {
		t.Fatal(err)
	}

	registered := ac.propertyRegistry.dynamicFunctionNames
	expected := maps.Keys(datamodel.AllBuiltInDynamicFunctions)
	if !arraysEqualOrderInsensitive(registered, expected) {
		t.Fatal("Not all built in functions registered or too many registered")
	}
}

func TestLoadingSignedConfig(t *testing.T) {
	// Signed with prod signature
	ac, err := buildTestAppCoreWithPath("../cmcore/data_model/test/testdata/primary_config/valid/signedValid.cmconfig", t)
	if err != nil {
		t.Fatal(err)
	}
	err = ac.Start(false)
	if err != nil {
		t.Fatal(err)
	}
	if ac.config == nil || ac.config.ConfigVersion != "v1" {
		t.Fatal("Failed to load signed config")
	}
	// is_debug_build should be false
	result, err := ac.propertyRegistry.evaluateCondition(testHelperNewCondition("is_debug_build", t))
	if err != nil {
		t.Fatal(err)
	}
	if result != false {
		t.Fatal("is_debug_build should be false")
	}
}

func TestLoadingEmptyUnsignedConfig(t *testing.T) {
	// Signed special case: empty without signature should work
	ac, err := buildTestAppCoreWithPath("../cmcore/data_model/test/testdata/primary_config/valid/emptyValid.json", t)
	if err != nil {
		t.Fatal(err)
	}
	err = ac.Start(false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadingJsonOnlyAllowedInDebug(t *testing.T) {
	ac, err := testBuildValidTestAppCore(t)
	if err != nil {
		t.Fatal(err)
	}
	// Debug=false should not allow unsigned
	err = ac.Start(false)
	if err == nil || ac.config != nil {
		t.Fatal("Should not load json config unless in debug mode", err)
	}
	// Debug=true should load unsigned/json
	err = ac.Start(true)
	if err != nil || ac.config == nil || ac.config.AppId != "io.criticalmoments.demo" {
		t.Fatal("Should not load json config unless in debug mode")
	}
}

func TestDebugBuild(t *testing.T) {
	ac, err := testBuildValidTestAppCore(t)
	if err != nil {
		t.Fatal(err)
	}
	err = ac.Start(true)
	if err != nil {
		t.Fatal(err)
	}
	result, err := ac.propertyRegistry.evaluateCondition(testHelperNewCondition("is_debug_build", t))
	if err != nil {
		t.Fatal(err)
	}
	if result != true {
		t.Fatal("is_debug_build should be true")
	}
}

func TestChecksAppId(t *testing.T) {
	ac, err := buildTestAppCoreWithPath("../cmcore/data_model/test/testdata/primary_config/invalid/invalidAppId.json", t)
	if err != nil {
		t.Fatal(err)
	}
	err = ac.Start(true)
	if err == nil || !strings.Contains(err.Error(), "this config file isn't valid for this app") {
		t.Fatal("Allowed loading a config with a bundle ID mismatch")
	}
}

func TestStartupAndCustomPropsRecordPropHistory(t *testing.T) {
	ac, err := testBuildValidTestAppCore(t)
	if err != nil {
		t.Fatal(err)
	}

	ac.propertyRegistry.builtInPropertyTypes = buildTestBuiltInProps(map[string]*datamodel.CMPropertyConfig{
		"builtInString": {Type: reflect.String, Source: datamodel.CMPropertySourceLib, Optional: false, SampleType: datamodel.CMPropertySampleTypeAppStart},
		"builtInNever":  {Type: reflect.String, Source: datamodel.CMPropertySourceLib, Optional: false, SampleType: datamodel.CMPropertySampleTypeDoNotSample},
	})

	ac.RegisterClientIntProperty("testInt", 42)
	ac.RegisterStaticStringProperty("builtInString", "hello world")
	ac.RegisterStaticStringProperty("builtInNever", "never")

	err = ac.Start(true)
	if err != nil {
		t.Fatal(err)
	}

	if v, err := ac.db.LatestPropertyHistory("custom_testInt"); err != nil || v != int64(42) {
		t.Fatal("custom static property not recorded in history")
	}
	if v, err := ac.db.LatestPropertyHistory("builtInString"); err != nil || v != "hello world" {
		t.Fatal("built in static sample_on_start property not recorded in history")
	}
	if v, err := ac.db.LatestPropertyHistory("builtInNever"); err != sql.ErrNoRows || v != nil {
		t.Fatal("built in static sample_never property recorded in history")
	}

	// Verify property history dynamic function also works
	result, err := ac.propertyRegistry.evaluateCondition(testHelperNewCondition("propertyHistoryLatestValue('builtInString') == 'hello world' && propertyHistoryLatestValue('builtInNever') == nil && propertyHistoryLatestValue('custom_testInt') == 42", t))
	if err != nil || !result {
		t.Fatal("Property history not working through condition function")
	}

	// Verify property history dynamic value check function also works
	result, err = ac.propertyRegistry.evaluateCondition(testHelperNewCondition("propertyEver('builtInString', 'hello world') && !propertyEver('builtInString', 'hello world2')", t))
	if err != nil || !result {
		t.Fatal("Property history by value not working through condition function")
	}
}
func TestAppStartEvent(t *testing.T) {
	ac, err := testBuildValidTestAppCore(t)
	if err != nil {
		t.Fatal(err)
	}

	count, err := ac.db.EventCountByName(datamodel.AppStartBuiltInEvent)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatal("App start event should not be recorded before start")
	}

	err = ac.Start(true)
	if err != nil {
		t.Fatal(err)
	}

	count, err = ac.db.EventCountByName(datamodel.AppStartBuiltInEvent)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatal("App start event should have fired")
	}
}

func TestStableRandomOperator(t *testing.T) {
	ac, err := testBuildValidTestAppCore(t)
	if err != nil {
		t.Fatal(err)
	}

	result, err := ac.propertyRegistry.evaluateCondition(testHelperNewCondition("stableRand() != 0 && stableRand() != nil && stableRand() == stableRand()", t))
	if err != nil || !result {
		t.Fatal("failed to generate consistent stableRand()")
	}
}

func TestMinConfigVersionChecks(t *testing.T) {
	tests := map[string]bool{
		"../cmcore/data_model/test/testdata/primary_config/invalid/appVersionTooLow.json":     false,
		"../cmcore/data_model/test/testdata/primary_config/invalid/cmVersionTooLow.json":      false, // add_test_count
		"../cmcore/data_model/test/testdata/primary_config/invalid/cmVersionInvalid.json":     false, // add_test_count
		"../cmcore/data_model/test/testdata/primary_config/valid/cmVersionHighEnough.json":    true,  // add_test_count
		"../cmcore/data_model/test/testdata/primary_config/valid/appVersionHighEnough.json":   true,  // add_test_count
		"../cmcore/data_model/test/testdata/primary_config/invalid/cmVersionTooLowInt.json":   false, // add_test_count
		"../cmcore/data_model/test/testdata/primary_config/invalid/cmVersionInvalidInt.json":  false, // add_test_count
		"../cmcore/data_model/test/testdata/primary_config/valid/cmVersionHighEnoughInt.json": true,  // add_test_count
	}
	for path, shouldPass := range tests {
		ac, err := buildTestAppCoreWithPath(path, t)
		if err != nil {
			t.Fatal(err)
		}
		err = ac.Start(true)
		if (err == nil) != shouldPass {
			t.Fatalf("Config version check failed for %v", path)
		}
	}
}

func TestSetLogEvents(t *testing.T) {
	ac, err := testBuildValidTestAppCore(t)
	if err != nil {
		t.Fatal(err)
	}

	if ac.eventManager.logEvents {
		t.Fatal("logEvents should be false by default")
	}
	ac.SetLogEvents(true)
	if !ac.eventManager.logEvents {
		t.Fatal("logEvents should be true after setting")
	}
}
