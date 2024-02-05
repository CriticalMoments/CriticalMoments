package datamodel

import (
	"fmt"
)

type eventTypeEnum int

const (
	// Built in events only allowed to be fired by SDK
	EventTypeBuiltIn eventTypeEnum = 0
	// well known event type like (is signed in), that apps provide but SDK is aware of
	EventTypeWellKnown eventTypeEnum = 1
	// Competely custom events from app the SDK is not aware of
	EventTypeCustom eventTypeEnum = 2
)

// Enum type would be nice, but doesn't play well with gomobile exports
const (
	AppStartBuiltInEvent             string = "app_start"
	SessionStartBuiltInEvent         string = "session_start"
	AppEnteredForgroundBuiltInEvent  string = "app_entered_foreground"
	AppEnteredBackgroundBuiltInEvent string = "app_entered_background"
	AppTerminatedBuiltInEvent        string = "app_terminated"
)

var (
	allBuiltInEventTypes = map[string]bool{
		AppStartBuiltInEvent:             true,
		SessionStartBuiltInEvent:         true,
		AppEnteredForgroundBuiltInEvent:  true,
		AppEnteredBackgroundBuiltInEvent: true,
		AppTerminatedBuiltInEvent:        true,
	}
)

// Enum type would be nice, but doesn't play well with gomobile exports
const (
	SignedInEvent         string = "signed_in"
	LogOutEvent           string = "log_out"
	SignUpEvent           string = "sign_up"
	InAppPurchaseEvent    string = "in_app_purchase"
	SubscribeEvent        string = "subscribe"
	UnsubscribeEvent      string = "unsubscribe"
	SearchEvent           string = "search"
	TutorialCompleteEvent string = "tutorial_complete"
	AddToCartEvent        string = "add_to_cart"
	PhysicalPurchaseEvent string = "physical_purchase"
	LevelUpEvent          string = "level_up"
)

var (
	allWellKnownEventTypes = map[string]bool{
		SignedInEvent:         true,
		LogOutEvent:           true,
		SignUpEvent:           true,
		InAppPurchaseEvent:    true,
		SubscribeEvent:        true,
		UnsubscribeEvent:      true,
		SearchEvent:           true,
		TutorialCompleteEvent: true,
		AddToCartEvent:        true,
		PhysicalPurchaseEvent: true,
		LevelUpEvent:          true,
	}
)

type Event struct {
	Name string

	// Event type is internal to cmcore
	EventType eventTypeEnum
}

// Clients can send well known or custom events, but not built in
func NewClientEventWithName(name string) (*Event, error) {
	if name == "" {
		return nil, fmt.Errorf("event name required")
	}
	if allBuiltInEventTypes[name] {
		return nil, fmt.Errorf("built in events cannot be fired by client")
	}

	isWellKnown := allWellKnownEventTypes[name]
	if isWellKnown {
		return &Event{
			Name:      name,
			EventType: EventTypeWellKnown,
		}, nil
	} else {
		return &Event{
			Name:      name,
			EventType: EventTypeCustom,
		}, nil
	}
}

func NewBuiltInEventWithName(name string) (*Event, error) {
	// Ensure this is a built in event we recognize
	if !allBuiltInEventTypes[name] {
		return nil, fmt.Errorf("unknown built in event: %v", name)
	}

	e := Event{
		Name:      name,
		EventType: EventTypeBuiltIn,
	}
	return &e, nil
}

func NewWellKnownEventWithName(name string) (*Event, error) {
	// Ensure this is a well known event we recognize
	if !allWellKnownEventTypes[name] {
		return nil, fmt.Errorf("unknown well known event: %v", name)
	}

	e := Event{
		Name:      name,
		EventType: EventTypeWellKnown,
	}
	return &e, nil
}

func NewCustomEventWithName(name string) (*Event, error) {
	e := Event{
		Name:      name,
		EventType: EventTypeCustom,
	}
	return &e, nil
}
