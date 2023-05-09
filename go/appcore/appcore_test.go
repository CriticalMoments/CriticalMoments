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

type NilLibBindings struct{}

func (lb *NilLibBindings) ShowBanner(b *datamodel.BannerAction) error {
	return nil
}
func (lb *NilLibBindings) SetDefaultTheme(theme *datamodel.Theme) error {
	return nil
}

func TestAppcoreHardcode(t *testing.T) {
	ac := Appcore{}
	err := ac.Start()
	if err == nil {
		t.Fatal("Should not start without config")
	}
	configPath, err := filepath.Abs("../cmcore/data_model/test/testdata/primary_config/valid/maximalValid.json")
	if err != nil {
		t.Fatal(err)
	}
	configUrl := fmt.Sprintf("file://%v", configPath)
	err = ac.SetConfigUrl(configUrl)
	if err != nil {
		t.Fatal(err)
	}
	err = ac.Start()
	if err == nil {
		t.Fatal("Should not start without lib bindings")
	}
	lb := NilLibBindings{}
	ac.RegisterLibraryBindings(&lb)
	err = ac.Start()
	if err != nil {
		t.Fatal(err)
	}
}
