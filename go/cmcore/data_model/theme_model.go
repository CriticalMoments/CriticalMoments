package datamodel

import (
	"encoding/json"
	"fmt"

	"golang.org/x/exp/slices"
)

type Theme struct {
	// Banners
	BannerBackgroundColor string // eg: "#ffffff"
	BannerForegroundColor string // eg: "#000000"

	// Colors
	PrimaryColor       string // eg: "#ff0000"
	BackgroundColor    string // eg: "#ffffff"
	PrimaryTextColor   string // eg: "#000000"
	SecondaryTextColor string // eg: "#222222"

	// Fonts
	FontName                   string
	BoldFontName               string
	ScaleFontForUserPreference bool
	FontScale                  float64

	// Dark Mode Theme
	DarkModeTheme *Theme

	// Fallback Theme Name
	IsFallthoughTheme bool
	FallbackThemeName string
}

// Currently very close to Theme model, but don't want to couple
// serialization and app data model
type jsonTheme struct {
	// Banners
	BannerBackgroundColor string `json:"bannerBackgroundColor"` // "#ffffff"
	BannerForegroundColor string `json:"bannerForegroundColor"` // "#000000"

	// Colors
	PrimaryColor       string `json:"primaryColor,omitempty"`
	BackgroundColor    string `json:"backgroundColor,omitempty"`
	PrimaryTextColor   string `json:"primaryTextColor,omitempty"`
	SecondaryTextColor string `json:"secondaryTextColor,omitempty"`

	// Fonts
	FontName                   string   `json:"fontName,omitempty"`
	BoldFontName               string   `json:"boldFontName,omitempty"`
	ScaleFontForUserPreference *bool    `json:"scaleFontForUserPreference,omitempty"` // pointer == nullable
	FontScale                  *float64 `json:"fontScale,omitempty"`                  // pointer == nullable

	// Dark mode theme
	DarkModeTheme *Theme `json:"darkModeTheme,omitempty"`

	// Fallback Theme Name
	FallbackThemeName string `json:"fallback,omitempty"`
}

// These themes provided by libary level, depend on the system
var libraryThemeNames map[string]bool = map[string]bool{
	"system":       true,
	"system_light": true,
	"system_dark":  true,
}

var builtInThemes map[string]*Theme = map[string]*Theme{
	"elegant_light": {
		BannerBackgroundColor:      "#000000",
		BannerForegroundColor:      "#ffffff",
		PrimaryColor:               "#000000",
		BackgroundColor:            "#ffffff",
		PrimaryTextColor:           "#000000",
		SecondaryTextColor:         "#222222",
		FontName:                   "AvenirNext-Regular",
		BoldFontName:               "AvenirNext-Bold",
		FontScale:                  1.0,
		ScaleFontForUserPreference: true,
	},
	"elegant_dark": {
		BannerBackgroundColor:      "#ffffff",
		BannerForegroundColor:      "#000000",
		PrimaryColor:               "#ffffff",
		BackgroundColor:            "#000000",
		PrimaryTextColor:           "#ffffff",
		SecondaryTextColor:         "#dddddd",
		FontName:                   "AvenirNext-Regular",
		BoldFontName:               "AvenirNext-Bold",
		FontScale:                  1.0,
		ScaleFontForUserPreference: true,
	},
	// terminal MIT source: https://github.com/altercation/solarized
	"terminal_light": {
		BannerBackgroundColor:      "#fdf6e3",
		BannerForegroundColor:      "#859900",
		PrimaryColor:               "#859900",
		BackgroundColor:            "#fdf6e3",
		PrimaryTextColor:           "#859900",
		SecondaryTextColor:         "#b58900",
		FontName:                   "Courier",
		BoldFontName:               "Courier-Bold",
		FontScale:                  1.0,
		ScaleFontForUserPreference: false,
	},
	"terminal_dark": {
		BannerBackgroundColor:      "#001e27",
		BannerForegroundColor:      "#6cbe6c",
		PrimaryColor:               "#6cbe6c",
		BackgroundColor:            "#001e27",
		PrimaryTextColor:           "#6cbe6c",
		SecondaryTextColor:         "#51ef84",
		FontName:                   "Courier",
		BoldFontName:               "Courier-Bold",
		FontScale:                  1.0,
		ScaleFontForUserPreference: false,
	},
	"jazzy_dark": {
		BannerBackgroundColor:      "#3d088c",
		BannerForegroundColor:      "#ffffff",
		PrimaryColor:               "#c1316d",
		BackgroundColor:            "#3d088c",
		PrimaryTextColor:           "#ffffff",
		SecondaryTextColor:         "#e2deed",
		FontScale:                  1.0,
		ScaleFontForUserPreference: false,
	},
	"jazzy_light": {
		BannerBackgroundColor:      "#3d088c",
		BannerForegroundColor:      "#ffffff",
		PrimaryColor:               "#c1316d",
		BackgroundColor:            "#ffffff",
		PrimaryTextColor:           "#3d088c",
		SecondaryTextColor:         "#64478d",
		FontScale:                  1.0,
		ScaleFontForUserPreference: false,
	},
	"honey_dark": {
		BannerBackgroundColor:      "#f4d42b",
		BannerForegroundColor:      "#000000",
		PrimaryColor:               "#000000",
		BackgroundColor:            "#f4d42b",
		PrimaryTextColor:           "#000000",
		SecondaryTextColor:         "#222222",
		FontScale:                  1.0,
		ScaleFontForUserPreference: false,
	},
	"honey_light": {
		BannerBackgroundColor:      "#f4d42b",
		BannerForegroundColor:      "#000000",
		PrimaryColor:               "#ebc603",
		BackgroundColor:            "#ffffff",
		PrimaryTextColor:           "#000000",
		SecondaryTextColor:         "#222222",
		FontScale:                  1.0,
		ScaleFontForUserPreference: false,
	},
	"sea_dark": {
		BannerBackgroundColor:      "#ffffff",
		BannerForegroundColor:      "#112c5d",
		PrimaryColor:               "#ffffff",
		BackgroundColor:            "#112c5d",
		PrimaryTextColor:           "#ffffff",
		SecondaryTextColor:         "#8eabdc",
		FontScale:                  1.0,
		ScaleFontForUserPreference: false,
	},
	"sea_light": {
		BannerBackgroundColor:      "#112c5d",
		BannerForegroundColor:      "#ffffff",
		PrimaryColor:               "#112c5d",
		BackgroundColor:            "#ffffff",
		PrimaryTextColor:           "#112c5d",
		SecondaryTextColor:         "#244278",
		FontScale:                  1.0,
		ScaleFontForUserPreference: false,
	},
}

var combinedThemeNames []string = []string{
	"terminal",
	"elegant",
	"jazzy",
	"honey",
	"sea",
}

func AllBuiltInThemeNames() []string {
	themeNames := combinedThemeNames
	for themeName := range builtInThemes {
		themeNames = append(themeNames, themeName)
	}
	return themeNames
}

func builtInThemeByName(name string) (*Theme, error) {
	builtIn, ok := builtInThemes[name]
	if ok {
		return builtIn, nil
	}

	if slices.Contains(combinedThemeNames, name) {
		darkTheme, darkOk := builtInThemes[name+"_dark"]
		lightTheme, lightOk := builtInThemes[name+"_light"]
		if darkOk && lightOk {
			// copy and return combined
			d := *darkTheme
			l := *lightTheme
			l.DarkModeTheme = &d
			return &l, nil
		}
	}

	return nil, NewUserPresentableError(fmt.Sprintf("Theme name '%v' is not valid. Check the spelling and try again.", name))
}

var (
	// For integration tests through to clients
	testTheme = Theme{
		BannerBackgroundColor:      "#ff0000", // Red
		BannerForegroundColor:      "#00ff00", // Green
		FontScale:                  1.1,
		FontName:                   "Palatino-Roman",
		BoldFontName:               "Palatino-Bold",
		ScaleFontForUserPreference: false,
		PrimaryColor:               "#ff0000",
		BackgroundColor:            "#ffffff",
		PrimaryTextColor:           "#ff0000",
		SecondaryTextColor:         "#00ff00",
		DarkModeTheme: &Theme{
			BannerBackgroundColor: "#00ff00", // Green
			BannerForegroundColor: "#ff0000", // Red
		},
	}
)

func TestTheme() *Theme {
	return &testTheme
}

func (t *Theme) UnmarshalJSON(data []byte) error {
	var jt jsonTheme
	err := json.Unmarshal(data, &jt)
	if err != nil {
		return NewUserPresentableErrorWSource("Unable to parse the json of a theme. Check the format, variable names, and types (eg float vs int).", err)
	}

	uperr := parseThemeFromJsonTheme(t, &jt)
	if uperr != nil {
		return uperr
	}

	return nil
}

func parseThemeFromJsonTheme(t *Theme, jt *jsonTheme) *UserPresentableError {
	// Default Values for nullable options
	t.ScaleFontForUserPreference = true
	if jt.ScaleFontForUserPreference != nil {
		t.ScaleFontForUserPreference = *jt.ScaleFontForUserPreference
	}
	t.FontScale = 1.0
	if jt.FontScale != nil {
		t.FontScale = *jt.FontScale
	}

	// Passthough values
	t.BannerBackgroundColor = jt.BannerBackgroundColor
	t.BannerForegroundColor = jt.BannerForegroundColor
	t.PrimaryColor = jt.PrimaryColor
	t.BackgroundColor = jt.BackgroundColor
	t.PrimaryTextColor = jt.PrimaryTextColor
	t.SecondaryTextColor = jt.SecondaryTextColor
	t.FontName = jt.FontName
	t.BoldFontName = jt.BoldFontName
	t.DarkModeTheme = jt.DarkModeTheme
	t.FallbackThemeName = jt.FallbackThemeName

	if validationIssue := t.ValidateReturningUserReadableIssue(); validationIssue != "" {
		// parse time, we can ignore validation issues if fallback present.
		// the primary config will validate the fallback theme exists and is valid
		if StrictDatamodelParsing || t.FallbackThemeName == "" {
			return NewUserPresentableError(validationIssue)
		}
		t.IsFallthoughTheme = true
	}

	return nil
}

func (t *Theme) Validate() bool {
	return t.ValidateReturningUserReadableIssue() == ""
}

func (t *Theme) ValidateReturningUserReadableIssue() string {
	// Fallthough themes are valid as long as they have fallback name
	if t.IsFallthoughTheme && t.FallbackThemeName != "" {
		return ""
	}

	return t.ValidateDisallowFallthoughReturningUserReadableIssue()
}

func (t *Theme) ValidateDisallowFallthoughReturningUserReadableIssue() string {
	// Check all colors are valid, but allow empty
	colors := []string{t.BackgroundColor, t.BannerBackgroundColor, t.BannerForegroundColor, t.PrimaryColor, t.PrimaryTextColor, t.SecondaryTextColor}
	for _, color := range colors {
		if !stringColorIsValidAllowEmpty(color) {
			return fmt.Sprintf("Color isn't a valid color. Should be in format '#ffffff' (lower case only). Found \"%v\"", color)
		}
	}

	if t.FontScale < 0.5 || t.FontScale > 2.0 {
		return "Font scale must be in the range 0.5-2.0"
	}

	return ""
}

func stringColorIsValidAllowEmpty(c string) bool {
	if len(c) == 0 {
		return true
	}
	return stringColorIsValid(c)
}

func stringColorIsValid(colorString string) bool {
	if len(colorString) != 7 {
		return false
	}

	// Verify format #fff000
	if colorString[0] != '#' {
		return false
	}
	for _, c := range colorString[1:] {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			return false
		}
	}

	return true
}
