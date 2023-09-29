package appcore

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	datamodel "github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model"
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

type testLibBindings struct {
	lastBannerAction *datamodel.BannerAction
	lastAlertAction  *datamodel.AlertAction
	lastLinkAction   *datamodel.LinkAction
	reviewCount      int
	defaultTheme     *datamodel.Theme
	lastModal        *datamodel.ModalAction
}

func (lb *testLibBindings) ShowBanner(b *datamodel.BannerAction) error {
	lb.lastBannerAction = b
	return nil
}
func (lb *testLibBindings) ShowAlert(a *datamodel.AlertAction) error {
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
func (lb *testLibBindings) ShowReviewPrompt() error {
	lb.reviewCount += 1
	return nil
}
func (lb *testLibBindings) ShowModal(modal *datamodel.ModalAction) error {
	lb.lastModal = modal
	return nil
}

func testBuildValidTestAppCore(t *testing.T) (*Appcore, error) {
	ac := newAppcore()
	configPath, err := filepath.Abs("../cmcore/data_model/test/testdata/primary_config/valid/maximalValid.json")
	if err != nil {
		t.Fatal(err)
	}
	configUrl := fmt.Sprintf("file://%v", configPath)
	err = ac.SetConfigUrl(configUrl)
	if err != nil {
		t.Fatal(err)
	}
	baseCachePath := fmt.Sprintf("/tmp/criticalmoments/test-temp-%v", rand.Int())
	os.MkdirAll(baseCachePath, os.ModePerm)
	err = ac.SetCacheDirPath(baseCachePath)
	if err != nil {
		t.Fatal(err)
	}
	lb := testLibBindings{}
	ac.RegisterLibraryBindings(&lb)

	ac.SetApiKey("CM1-aGVsbG86d29ybGQ=-Yjppby5jcml0aWNhbG1vbWVudHMuZGVtbw==-MEUCIQCUfx6xlmQ0kdYkuw3SMFFI6WXrCWKWwetXBrXXG2hjAwIgWBPIMrdM1ET0HbpnXlnpj/f+VXtjRTqNNz9L/AOt4GY=", "io.criticalmoments.demo")

	// Clear required properties, for easier setup
	ac.propertyRegistry = newPropertyRegistry()
	ac.propertyRegistry.requiredPropertyTypes = map[string]reflect.Kind{}
	return &ac, nil
}

func TestAppcoreStart(t *testing.T) {
	ac, err := testBuildValidTestAppCore(t)
	if err != nil {
		t.Fatal(err)
	}
	err = ac.Start()
	if err != nil {
		t.Fatal(err)
	}
	// Check it loaded the config (more detailed test of parsing in cmcore)
	if ac.config.DefaultTheme == nil {
		t.Fatal("Failed to load config in Appcore setup")
	}
}

func TestAppcoreStartMissingConfig(t *testing.T) {
	ac, err := testBuildValidTestAppCore(t)
	if err != nil {
		t.Fatal(err)
	}
	ac.configUrlString = ""
	err = ac.Start()
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
	err = ac.Start()
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
	err = ac.Start()
	if err == nil {
		t.Fatal("Should not start with bad config")
	}
	if ac.config != nil {
		t.Fatal("Loaded config from bad url")
	}
}

func TestPerformingAction(t *testing.T) {
	ac, err := testBuildValidTestAppCore(t)
	if err != nil {
		t.Fatal(err)
	}
	err = ac.Start()
	if err != nil {
		t.Fatal(err)
	}
	if ac.libBindings.(*testLibBindings).lastBannerAction != nil {
		t.Fatal("last banner action should be nil on new appcore test binding")
	}
	// should fire bannerAction1 via a trigger
	err = ac.SendEvent("custom_event")
	if err != nil {
		t.Fatal(err)
	}
	if ac.libBindings.(*testLibBindings).lastBannerAction.Body != "Hello world, but on a banner!" {
		t.Fatal("last banner action should be nil on new appcore test binding")
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
		t.Fatal("alert event didn't fire")
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
}

func TestConditionalActionDispatching(t *testing.T) {
	ac, err := testBuildValidTestAppCore(t)
	if err != nil {
		t.Fatal(err)
	}
	err = ac.Start()
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
}

func TestSetDefaultTheme(t *testing.T) {
	ac, err := testBuildValidTestAppCore(t)
	if err != nil {
		t.Fatal(err)
	}
	if ac.libBindings.(*testLibBindings).defaultTheme != nil {
		t.Fatal("Theme should be nil until started")
	}
	err = ac.Start()
	if err != nil {
		t.Fatal(err)
	}
	defaultTheme := ac.libBindings.(*testLibBindings).defaultTheme
	if defaultTheme == nil && defaultTheme.BannerBackgroundColor != "#ffffff" {
		t.Fatal("Default theme not set after start")
	}
}

func TestNamedConditions(t *testing.T) {
	ac, err := testBuildValidTestAppCore(t)
	if err != nil {
		t.Fatal(err)
	}
	err = ac.Start()
	if err != nil {
		t.Fatal(err)
	}

	// conditions without overrides should use provided condition
	r, err := ac.CheckNamedCondition("newCondition1", "false")
	if err != nil || r {
		t.Fatal("false conditions failed")
	}
	r, err = ac.CheckNamedCondition("newCondition2", "true")
	if err != nil || !r {
		t.Fatal("false conditions failed")
	}

	// falseCondition should override provided string
	r, err = ac.CheckNamedCondition("falseCondition", "true")
	if err != nil || r {
		t.Fatal("false conditions failed")
	}

	// trueCondition should override provided string
	r, err = ac.CheckNamedCondition("trueCondition", "false")
	if err != nil || !r {
		t.Fatal("false conditions failed")
	}

	// Check name check
	r, err = ac.CheckNamedCondition("", "false")
	if err == nil {
		t.Fatal("CheckNamedCondition requires name and didn't validate empty string")
	}

	// Check debug mode checker
	dmerr := ac.CheckNamedConditionCollision("uniqueName", "false")
	if dmerr != nil {
		t.Fatal("dev mode condition failed")
	}
	dmerr = ac.CheckNamedConditionCollision("uniqueName", "false")
	if dmerr != nil {
		t.Fatal("dev mode condition second time errored, but should pass with same condition")
	}
	dmerr = ac.CheckNamedConditionCollision("uniqueName", "true")
	if dmerr == nil {
		t.Fatal("unque condition with new value should return a dev warning")
	}

}
