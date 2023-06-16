package testing

import (
	"io/ioutil"
	"os"
	"testing"

	datamodel "github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model"
)

func minValidTheme() datamodel.Theme {
	return datamodel.Theme{
		BannerBackgroundColor: "#ffffff",
		BannerForegroundColor: "#000000",
		FontScale:             1.0,
	}
}

func TestBuiltinThemesValid(t *testing.T) {
	if !datamodel.TestTheme().Validate() {
		t.Fatal()
	}
}

func TestColorValidation(t *testing.T) {
	theme := minValidTheme()
	if !theme.Validate() {
		t.Fatal()
	}
	// Too long
	theme.BannerBackgroundColor = "#fffffff"
	if theme.Validate() {
		t.Fatal()
	}
	// Too Short
	theme.BannerBackgroundColor = "#fffff"
	if theme.Validate() {
		t.Fatal()
	}
	// No #
	theme.BannerBackgroundColor = "ffffff"
	if theme.Validate() {
		t.Fatal()
	}
	// No # long
	theme.BannerBackgroundColor = "fffffff"
	if theme.Validate() {
		t.Fatal()
	}
	// invalid char: out of range
	theme.BannerBackgroundColor = "#00000g"
	if theme.Validate() {
		t.Fatal()
	}
	// invalid char: uppercase
	theme.BannerBackgroundColor = "#00000A"
	if theme.Validate() {
		t.Fatal()
	}
	// invalid char: out of range, position 1
	theme.BannerBackgroundColor = "#.00000"
	if theme.Validate() {
		t.Fatal()
	}
	// all valid chars part 1
	theme.BannerBackgroundColor = "#012345"
	if !theme.Validate() {
		t.Fatal()
	}
	// all valid chars part 2
	theme.BannerBackgroundColor = "#6789ab"
	if !theme.Validate() {
		t.Fatal()
	}
	// all valid chars part 3
	theme.BannerBackgroundColor = "#cdefff"
	if !theme.Validate() {
		t.Fatal()
	}

	// Each color should validate -- allows nil (above), allows valid, disallows invalid
	theme.BannerForegroundColor = "#x"
	if theme.Validate() {
		t.Fatal("allowed invalid color")
	}
	theme.BannerForegroundColor = "#000000"
	if !theme.Validate() {
		t.Fatal("disallowed valid color")
	}
	theme.PrimaryColor = "#x"
	if theme.Validate() {
		t.Fatal("allowed invalid color")
	}
	theme.PrimaryColor = "#000000"
	if !theme.Validate() {
		t.Fatal("disallowed valid color")
	}
	theme.PrimaryTextColor = "#x"
	if theme.Validate() {
		t.Fatal("allowed invalid color")
	}
	theme.PrimaryTextColor = "#000000"
	if !theme.Validate() {
		t.Fatal("disallowed valid color")
	}
	theme.SecondaryTextColor = "#x"
	if theme.Validate() {
		t.Fatal("allowed invalid color")
	}
	theme.SecondaryTextColor = "#000000"
	if !theme.Validate() {
		t.Fatal("disallowed valid color")
	}
	theme.BackgroundColor = "#x"
	if theme.Validate() {
		t.Fatal("allowed invalid color")
	}
	theme.BackgroundColor = "#000000"
	if !theme.Validate() {
		t.Fatal("disallowed valid color")
	}
}

func TestFontScaleValidation(t *testing.T) {
	theme := minValidTheme()
	validValues := []float64{1.0, 0.8, 0.9, 1.2, 2.0, 1.5, 0.5}
	for _, valid := range validValues {
		theme.FontScale = valid
		if !theme.Validate() {
			t.Fatalf("Font scale %v expected to be valid", valid)
		}
	}

	invalidValues := []float64{-1.0, 0.1, 0.499999, 2.000001, 0.0}
	for _, invalid := range invalidValues {
		theme.FontScale = invalid
		if theme.Validate() {
			t.Fatalf("Font scale %v expected to be invalid", invalid)
		}
	}
}

func TestJsonParsingValid(t *testing.T) {
	testJsonFolder("./testdata/themes/valid", true, t)
}

func TestJsonParsingInvalid(t *testing.T) {
	testJsonFolder("./testdata/themes/invalid", false, t)
}

func testJsonFolder(basePath string, expectSuccess bool, t *testing.T) {
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		t.Fatal()
	}
	for _, file := range files {
		testFileData, err := os.ReadFile(basePath + "/" + file.Name())
		if err != nil {
			t.Fatal()
		}
		var theme datamodel.Theme
		err = theme.UnmarshalJSON(testFileData)
		if expectSuccess {
			if err != nil {
				t.Fatalf("Theme failed to parse: %v", file.Name())
			}
			if !theme.Validate() {
				t.Fatalf("Theme failed to validate: %v", file.Name())
			}
		} else {
			if err == nil {
				t.Fatalf("Parsed theme when invalid: %v", file.Name())
			}
			// All errors should be user readable! We want to be able to tell user what was wrong
			_, ok := interface{}(err).(datamodel.UserPresentableErrorI)
			if !ok {
				t.Fatalf("Theme parsing issue didn't return user presentable error: %v", file.Name())
			}
		}
	}
}

func TestJsonParsingDefaultsTheme(t *testing.T) {
	testFileData, err := os.ReadFile("./testdata/themes/valid/minimalValidTheme.json")
	if err != nil {
		t.Fatal()
	}
	var theme datamodel.Theme
	err = theme.UnmarshalJSON(testFileData)
	if err != nil {
		t.Fatal()
	}

	// Check defaults for values not included in json
	if theme.FontScale != 1.0 {
		t.Fatal()
	}
	if theme.ScaleFontForUserPreference != true {
		t.Fatal()
	}
}
func TestJsonParsingAllFieldsTheme(t *testing.T) {
	testFileData, err := os.ReadFile("./testdata/themes/valid/maximalValidTheme.json")
	if err != nil {
		t.Fatal()
	}
	var theme datamodel.Theme
	err = theme.UnmarshalJSON(testFileData)
	if err != nil {
		t.Fatal()
	}

	if theme.FontScale != 1.1 {
		t.Fatal()
	}
	if theme.ScaleFontForUserPreference != false {
		t.Fatal()
	}
	if theme.BannerBackgroundColor != "#ffffff" {
		t.Fatal()
	}
	if theme.BannerForegroundColor != "#000000" {
		t.Fatal()
	}
	if theme.PrimaryColor != "#ff0000" {
		t.Fatal()
	}
	if theme.BackgroundColor != "#ffffff" {
		t.Fatal()
	}
	if theme.PrimaryTextColor != "#ff0000" {
		t.Fatal()
	}
	if theme.SecondaryTextColor != "#00ff00" {
		t.Fatal()
	}

	if theme.FontName != "AvenirNext-Regular" {
		t.Fatal()
	}
	if theme.BoldFontName != "AvenirNext-Bold" {
		t.Fatal()
	}

	if theme.DarkModeTheme.FontScale != 0.9 {
		t.Fatal()
	}
	if theme.DarkModeTheme.ScaleFontForUserPreference != false {
		t.Fatal()
	}
	if theme.DarkModeTheme.BannerBackgroundColor != "#000000" {
		t.Fatal()
	}
	if theme.DarkModeTheme.BannerForegroundColor != "#ffffff" {
		t.Fatal()
	}
	if theme.DarkModeTheme.FontName != "AvenirNext-RegularD" {
		t.Fatal()
	}
	if theme.DarkModeTheme.BoldFontName != "AvenirNext-BoldD" {
		t.Fatal()
	}
	if theme.DarkModeTheme.PrimaryColor != "#ff0000" {
		t.Fatal()
	}
	if theme.DarkModeTheme.BackgroundColor != "#ffffff" {
		t.Fatal()
	}
	if theme.DarkModeTheme.PrimaryTextColor != "#ff0000" {
		t.Fatal()
	}
	if theme.DarkModeTheme.SecondaryTextColor != "#00ff00" {
		t.Fatal()
	}
}

func TestJsonParsingInvalidNestedTheme(t *testing.T) {
	testFileData, err := os.ReadFile("./testdata/themes/invalid/invalidDarkTheme.json")
	if err != nil {
		t.Fatal()
	}
	var theme datamodel.Theme
	err = theme.UnmarshalJSON(testFileData)
	if err == nil {
		t.Fatal()
	}
	uerr, ok := err.(*datamodel.UserPresentableError)
	if !ok || uerr.SourceError == nil {
		t.Fatal("nested error not returned")
	}
	_, uok := uerr.SourceError.(*datamodel.UserPresentableError)
	if !uok {
		t.Fatal("Nested error not user presentable")
	}
}
