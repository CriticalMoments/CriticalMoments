package datamodel

import "testing"

func TestBuiltInThemesValid(t *testing.T) {

	for themeName, theme := range builtInThemes {
		issue := theme.ValidateReturningUserReadableIssue()
		if issue != "" {
			t.Errorf("Theme %v is invalid: %v", themeName, issue)
		}
	}

	issue := testTheme.ValidateReturningUserReadableIssue()
	if issue != "" {
		t.Errorf("Test theme is invalid: %v", issue)
	}
}
