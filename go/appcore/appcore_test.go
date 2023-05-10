package appcore

import (
	"fmt"
	"path/filepath"
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
	defaultTheme     *datamodel.Theme
}

func (lb *testLibBindings) ShowBanner(b *datamodel.BannerAction) error {
	lb.lastBannerAction = b
	return nil
}
func (lb *testLibBindings) SetDefaultTheme(theme *datamodel.Theme) error {
	lb.defaultTheme = theme
	return nil
}

func testBuildValidTestAppCore(t *testing.T) (*Appcore, error) {
	ac := Appcore{}
	configPath, err := filepath.Abs("../cmcore/data_model/test/testdata/primary_config/valid/maximalValid.json")
	if err != nil {
		t.Fatal(err)
	}
	configUrl := fmt.Sprintf("file://%v", configPath)
	err = ac.SetConfigUrl(configUrl)
	if err != nil {
		t.Fatal(err)
	}
	lb := testLibBindings{}
	ac.RegisterLibraryBindings(&lb)
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

func TestPerformNamedAction(t *testing.T) {
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
	// should fire bannerAction1
	ac.SendEvent("custom_event")
	if ac.libBindings.(*testLibBindings).lastBannerAction.Body != "Hello world, but on a banner!" {
		t.Fatal("last banner action should be nil on new appcore test binding")
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
