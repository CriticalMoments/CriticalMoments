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
	if !datamodel.ElegantTheme().Validate() {
		t.Fatal()
	}
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
	// empty: requires you set both colors
	theme.BannerForegroundColor = ""
	if theme.Validate() {
		t.Fatal()
	}
	// empty: requires you set both colors
	theme.BannerForegroundColor = "#000000"
	theme.BannerBackgroundColor = ""
	if theme.Validate() {
		t.Fatal()
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
		theme := datamodel.NewThemeFromJson(testFileData)
		if (theme == nil && expectSuccess) || (theme != nil && !expectSuccess) {
			t.Fatalf("Theme json parsing failure: %v (expected %v)", file.Name(), expectSuccess)
		}
		if expectSuccess && !theme.Validate() {
			t.Fatalf("Theme json file failed validating: %v (expected %v)", file.Name(), expectSuccess)
		}
	}
}

func TestJsonParsingDefaults(t *testing.T) {
	testFileData, err := os.ReadFile("./testdata/themes/valid/minimalValidTheme.json")
	if err != nil {
		t.Fatal()
	}
	theme := datamodel.NewThemeFromJson(testFileData)

	// Check defaults for values not included in json
	if theme.FontScale != 1.0 {
		t.Fatal()
	}
	if theme.ScaleFontForUserPreference != true {
		t.Fatal()
	}
}
func TestJsonParsingAllFields(t *testing.T) {
	testFileData, err := os.ReadFile("./testdata/themes/valid/maximalValidTheme.json")
	if err != nil {
		t.Fatal()
	}
	theme := datamodel.NewThemeFromJson(testFileData)

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
	if theme.FontName != "AvenirNext-Regular" {
		t.Fatal()
	}
	if theme.BoldFontName != "AvenirNext-Bold" {
		t.Fatal()
	}
}
