package appcore

import (
	"errors"
	"fmt"
	"strings"
)

type EventTypeEnum string

const (
	// Built in events only allowed to be fired by SDK
	EventTypeBuiltIn EventTypeEnum = "builtInEventType"
	// well known event type like (is signed in), that apps provide but SDK is aware of
	EventTypeWellKnown EventTypeEnum = "wellKnownEventType"
	// Competely custom events from app the SDK is not aware of
	EventTypeCustom EventTypeEnum = "customEventType"
)

const (
	buildInEventNamespace   = "io.criticalmoments.events.built_in."
	wellKnownEventNamespace = "io.criticalmoments.events.well_known."
)

type BuiltInEventTypeEnum string

const (
	AppStartBuiltInEvent BuiltInEventTypeEnum = BuiltInEventTypeEnum(buildInEventNamespace + "app_start")
)

var (
	allBuiltInEventTypes = map[BuiltInEventTypeEnum]bool{
		AppStartBuiltInEvent: true,
	}
)

type WellKnownEventTypeEnum string

const (
	SignedInEventAppStart WellKnownEventTypeEnum = WellKnownEventTypeEnum(wellKnownEventNamespace + "signed_in")
)

var (
	allWellKnownEventTypes = map[WellKnownEventTypeEnum]bool{
		SignedInEventAppStart: true,
	}
)

type Event struct {
	Name      string
	EventType EventTypeEnum
}

func NewBuiltInEventWithName(name string) (*Event, error) {
	// Ensure this is a built in event we recognize
	if !allBuiltInEventTypes[BuiltInEventTypeEnum(name)] {
		return nil, errors.New(fmt.Sprintf("Unknown built in event: %v", name))
	}
	if !strings.HasPrefix(name, buildInEventNamespace) {
		return nil, errors.New(fmt.Sprintf("Built in event outside namespace: %v", name))
	}

	e := Event{
		Name:      name,
		EventType: EventTypeBuiltIn,
	}
	return &e, nil
}

func NewWellKnownEventWithName(name string) (*Event, error) {
	// Ensure this is a well known event we recognize
	if !allWellKnownEventTypes[WellKnownEventTypeEnum(name)] {
		return nil, errors.New(fmt.Sprintf("Unknown well known event: %v", name))
	}
	if !strings.HasPrefix(name, wellKnownEventNamespace) {
		return nil, errors.New(fmt.Sprintf("Well known event outside namespace: %v", name))
	}

	e := Event{
		Name:      name,
		EventType: EventTypeWellKnown,
	}
	return &e, nil
}

func NewCustomEventWithName(name string) (*Event, error) {
	if strings.HasPrefix(name, wellKnownEventNamespace) || strings.HasPrefix(name, buildInEventNamespace) {
		return nil, errors.New(fmt.Sprintf("Attempted to log custom event matching built in or well known event: %v", name))
	}

	e := Event{
		Name:      name,
		EventType: EventTypeCustom,
	}
	return &e, nil
}
