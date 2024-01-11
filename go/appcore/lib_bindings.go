package appcore

import (
	datamodel "github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model"
)

// To be implemented by client libaray (eg: iOS SDK)
type LibBindings interface {
	// Themes
	SetDefaultTheme(theme *datamodel.Theme) error

	// Actions
	ShowBanner(banner *datamodel.BannerAction) error
	ShowAlert(alert *datamodel.AlertAction) error
	ShowLink(link *datamodel.LinkAction) error
	ShowReviewPrompt() error
	ShowModal(modal *datamodel.ModalAction) error

	// Condition functions
	CanOpenURL(url string) bool

	// Version numbers
	AppVersion() string
	CMVersion() string
}
