package appcore

import (
	datamodel "github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model"
)

// To be implemented by client libaray (eg: iOS SDK)
type LibBindings interface {
	// Themes
	SetDefaultTheme(theme *datamodel.Theme) error

	// Actions
	datamodel.ActionBindings
}
