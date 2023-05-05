package datamodel

import "encoding/json"

type Theme struct {
	// Banners
	BannerBackgroundColor string // eg: "#ffffff"
	BannerForegroundColor string // eg: "#000000"

	// Fonts
	FontName                   string
	BoldFontName               string
	ScaleFontForUserPreference bool
	FontScale                  float64
}

// Currently very close to Theme model, but don't want to couple
// serialization and app data model
type jsonTheme struct {
	// Banners
	BannerBackgroundColor string `json:"bannerBackgroundColor"` // "#ffffff"
	BannerForegroundColor string `json:"bannerForegroundColor"` // "#000000"

	// Fonts
	FontName                   string   `json:"fontName,omitempty"`
	BoldFontName               string   `json:"boldFontName,omitempty"`
	ScaleFontForUserPreference *bool    `json:"scaleFontForUserPreference,omitempty"` // pointer == nullable
	FontScale                  *float64 `json:"fontScale,omitempty"`                  // pointer == nullable
}

var (
	elegantTheme = Theme{
		BannerBackgroundColor: "#ffffff", // White
		BannerForegroundColor: "#000000", // Black
		FontScale:             1.0,
	}
	// Not pretty, but for e2e integration tests through to clients
	testTheme = Theme{
		BannerBackgroundColor:      "#ff0000", // Red
		BannerForegroundColor:      "#00ff00", // Green
		FontScale:                  1.1,
		FontName:                   "Palatino-Roman",
		BoldFontName:               "Palatino-Bold",
		ScaleFontForUserPreference: false,
	}
)

func ElegantTheme() *Theme {
	return &elegantTheme
}

func TestTheme() *Theme {
	return &testTheme
}

func NewThemeFromJson(data []byte) (*Theme, error) {
	var jt jsonTheme
	err := json.Unmarshal(data, &jt)
	if err != nil {
		return nil, NewUserPresentableErrorWSource("Unable to parse the json of a theme. Check the format, variable names, and types (eg float vs int).", err)
	}

	// Default Values for nullable options
	scaleFontForPref := true
	if jt.ScaleFontForUserPreference != nil {
		scaleFontForPref = *jt.ScaleFontForUserPreference
	}
	fontScale := 1.0
	if jt.FontScale != nil {
		fontScale = *jt.FontScale
	}

	t := Theme{
		BannerBackgroundColor:      jt.BannerBackgroundColor,
		BannerForegroundColor:      jt.BannerForegroundColor,
		FontName:                   jt.FontName,
		BoldFontName:               jt.BoldFontName,
		ScaleFontForUserPreference: scaleFontForPref,
		FontScale:                  fontScale,
	}

	if validationIssue := t.ValidateReturningUserReadableIssue(); validationIssue != "" {
		return nil, NewUserPresentableError(validationIssue)
	}

	return &t, nil
}

func (t Theme) Validate() bool {
	return t.ValidateReturningUserReadableIssue() == ""
}

func (t Theme) ValidateReturningUserReadableIssue() string {
	// Banner colors requires, and non-optional
	if !stringColorIsValid(t.BannerBackgroundColor) {
		return "Banner background color isn't a valid color. Should be in format '#ffffff' (lower case only)"
	}
	if !stringColorIsValid(t.BannerForegroundColor) {
		return "Banner foreground color isn't a valid color. Should be in format '#ffffff' (lower case only)"
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
