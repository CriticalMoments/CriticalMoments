package testing

import (
	"encoding/json"
	"os"
	"testing"

	datamodel "github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model"
)

func TestInvalidBannerMissingField(t *testing.T) {
	b := datamodel.BannerAction{}

	if b.Valid() {
		t.Fatal("Banners should require body")
	}
	b.Body = "Banner body"
	if !b.Valid() {
		t.Fatal("Minimal banner failed validation")
	}
	b.MaxLineCount = -2
	if b.Valid() {
		t.Fatal("Banner allowed negative max line count")
	}
	b.MaxLineCount = 4
	if !b.Valid() {
		t.Fatal("Minimal banner failed validation")
	}
	b.PreferredPosition = "invalid"
	if b.Valid() {
		t.Fatal("Banner allowed invalid position")
	}
	b.PreferredPosition = datamodel.BannerPositionBottom
	if !b.Valid() {
		t.Fatal("Banner disallowed valid position")
	}
}

func TestJsonParsingInvalidBanners(t *testing.T) {
	basePath := "./testdata/actions/banner/invalid"
	files, err := os.ReadDir(basePath)
	if err != nil {
		t.Fatal()
	}
	for _, file := range files {
		testFileData, err := os.ReadFile(basePath + "/" + file.Name())
		if err != nil {
			t.Fatal()
		}
		var ac datamodel.ActionContainer
		err = json.Unmarshal(testFileData, &ac)
		if err == nil {
			t.Fatalf("Parsed action when invalid: %v", file.Name())
		}
		// All errors should be user readable! We want to be able to tell user what was wrong
		_, ok := err.(*datamodel.UserPresentableError)
		if !ok {
			t.Fatalf("Banner parsing issue didn't return user presentable error: %v", file.Name())
		}
		if ac.ActionType != "" {
			t.Fatalf("Set type on invalid json. Invalid should not set type. %v", file.Name())
		}
		if ac.BannerAction != nil {
			t.Fatalf("Set BannerAction on invalid json: %v", file.Name())
		}
	}
}
func TestJsonParsingMinimalFieldsBanner(t *testing.T) {
	testFileData, err := os.ReadFile("./testdata/actions/banner/valid/minimalValid.json")
	if err != nil {
		t.Fatal()
	}
	var ac datamodel.ActionContainer
	err = json.Unmarshal(testFileData, &ac)
	if err != nil {
		t.Fatal()
	}

	if ac.ActionType != datamodel.ActionTypeEnumBanner {
		t.Fatal()
	}
	banner := ac.BannerAction
	if banner == nil || !banner.Valid() {
		t.Fatal()
	}
	if banner.Body != "Hello world, but on a banner!" {
		t.Fatal()
	}
	if banner.MaxLineCount != -1 {
		t.Fatal()
	}
	if banner.TapActionName != "" {
		t.Fatal()
	}
	if banner.CustomThemeName != "" {
		t.Fatal()
	}
	if banner.ShowDismissButton != true {
		t.Fatal()
	}
}

func TestJsonParsingAllFieldsBanner(t *testing.T) {
	testFileData, err := os.ReadFile("./testdata/actions/banner/valid/maximalValid.json")
	if err != nil {
		t.Fatal()
	}
	var ac datamodel.ActionContainer
	err = json.Unmarshal(testFileData, &ac)
	if err != nil {
		t.Fatal()
	}

	if ac.ActionType != datamodel.ActionTypeEnumBanner {
		t.Fatal()
	}
	banner := ac.BannerAction
	if banner == nil || !banner.Valid() {
		t.Fatal()
	}
	if banner.Body != "Hello world, but on a banner!" {
		t.Fatal()
	}
	if banner.MaxLineCount != 1 {
		t.Fatal()
	}
	if banner.TapActionName != "customAction" {
		t.Fatal()
	}
	if banner.CustomThemeName != "navy" {
		t.Fatal()
	}
	if banner.ShowDismissButton == true {
		t.Fatal()
	}
	if banner.PreferredPosition != datamodel.BannerPositionTop {
		t.Fatal()
	}
}
