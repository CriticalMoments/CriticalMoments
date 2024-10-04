package datamodel

import "testing"

func TestBuiltInThemesValid(t *testing.T) {
	themeNames := AllBuiltInThemeNames()

	// test no typos
	if len(themeNames) != len(builtInThemes)+len(combinedThemeNames) {
		t.Fatal("Incorrect number of themes")
	}
	if len(builtInThemes) != 2*len(combinedThemeNames) {
		t.Fatal("Incorrect number of themes")
	}
	if BaseThemeCount() != len(combinedThemeNames)+1 { // +1 for "system"
		t.Fatal("Incorrect number of themes")
	}

	for _, themeName := range themeNames {
		theme, err := builtInThemeByName(themeName)
		if err != nil {
			t.Errorf("Failed to get built in theme by name %v: %v", themeName, err)
		}
		issue := theme.Check()
		if issue != nil {
			t.Errorf("Theme %v is invalid: %v", themeName, issue)
		}
	}

	issue := testTheme.Check()
	if issue != nil {
		t.Errorf("Test theme is invalid: %v", issue)
	}

	// check combining isn't permuting the themes
	terminal_light, err := builtInThemeByName("terminal_light")
	if err != nil {
		t.Fatal("Failed to get terminal_light")
	}
	if terminal_light.DarkModeTheme != nil {
		t.Fatal("light theme mutated")
	}
	terminal, err := builtInThemeByName("terminal")
	if err != nil {
		t.Fatal("Failed to get terminal")
	}
	if terminal.BackgroundColor != terminal_light.BackgroundColor ||
		terminal.BannerBackgroundColor != terminal_light.BannerBackgroundColor ||
		terminal.BannerForegroundColor != terminal_light.BannerForegroundColor ||
		terminal.BoldFontName != terminal_light.BoldFontName ||
		terminal.FontName != terminal_light.FontName ||
		terminal.FontScale != terminal_light.FontScale ||
		terminal.PrimaryColor != terminal_light.PrimaryColor ||
		terminal.PrimaryTextColor != terminal_light.PrimaryTextColor ||
		terminal.ScaleFontForUserPreference != terminal_light.ScaleFontForUserPreference ||
		terminal.SecondaryTextColor != terminal_light.SecondaryTextColor {
		t.Fatal("combined themes don't match")
	}

}
