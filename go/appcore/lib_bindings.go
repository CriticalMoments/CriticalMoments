package appcore

import (
	"errors"
	"fmt"

	datamodel "github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model"
)

// To be implemented by client libaray (eg: iOS SDK)
type LibBindings interface {
	// Themes
	SetDefaultTheme(theme *datamodel.Theme) error

	// Actions
	ShowBanner(banner *datamodel.BannerAction) error
}

func dispatchActionToLib(action *datamodel.ActionContainer, lb LibBindings) error {
	switch action.ActionType {
	case datamodel.ActionTypeEnumBanner:
		return lb.ShowBanner(action.BannerAction)
	default:
		return errors.New(fmt.Sprintf("Action Dispatcher doesn't support action type %v", action.ActionType))
	}
}
