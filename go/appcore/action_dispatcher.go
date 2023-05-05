package appcore

import datamodel "github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model"

// To be implemented by client libaray (eg: iOS SDK)
type LibActionDispatcher interface {
	ShowBanner(bannerAction *datamodel.BannerAction)
}

var primaryActionDispatcher LibActionDispatcher

func RegisterActionDispatcher(ad LibActionDispatcher) {
	primaryActionDispatcher = ad
}

// TODO: remove this prior to V1. Just for e2e testing for now.
func InternalDipatchBannerFromGo() {
	if primaryActionDispatcher == nil {
		return
	}

	banner := datamodel.BannerAction{
		Body:              "Banner defined in appcore!",
		ShowDismissButton: true,
	}

	primaryActionDispatcher.ShowBanner(&banner)
}
