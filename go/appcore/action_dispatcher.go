package appcore

import datamodel "github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model"

// Action dispatcher wraps logic for dispatching actions. Some go straight to native libBindings, but some require some pre-work
type actionDispatcher struct {
	appcore *Appcore
}

func (ap *actionDispatcher) ShowBanner(banner *datamodel.BannerAction) error {
	return ap.appcore.libBindings.ShowBanner(banner)
}

func (ap *actionDispatcher) ShowAlert(alert *datamodel.AlertAction) error {
	return ap.appcore.libBindings.ShowAlert(alert)
}

func (ap *actionDispatcher) ShowLink(link *datamodel.LinkAction) error {
	return ap.appcore.libBindings.ShowLink(link)
}

func (ap *actionDispatcher) ShowReviewPrompt() error {
	return ap.appcore.libBindings.ShowReviewPrompt()
}

func (ap *actionDispatcher) ShowModal(modal *datamodel.ModalAction) error {
	return ap.appcore.libBindings.ShowModal(modal)
}

func (ap *actionDispatcher) PerformConditionalAction(ca *datamodel.ConditionalAction) error {
	passed, err := ap.appcore.propertyRegistry.evaluateCondition(ca.Condition)
	if err != nil {
		return err
	}
	if passed {
		return ap.appcore.PerformNamedAction(ca.PassedActionName)
	} else if ca.FailedActionName != "" {
		return ap.appcore.PerformNamedAction(ca.FailedActionName)
	}
	return nil
}

func (ap *actionDispatcher) PerformNamedAction(actionName string) error {
	return ap.appcore.PerformNamedAction(actionName)
}
