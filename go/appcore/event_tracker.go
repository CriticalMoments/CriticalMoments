package appcore

type EventTracker struct {
}

var sharedEventTracker EventTracker = EventTracker{}

func SharedEventTracker() *EventTracker {
	return &sharedEventTracker
}

func (ec EventTracker) SendEvent(e string) {

}
