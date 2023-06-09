package datamodel

import (
	"errors"
	"fmt"
	"strings"
)

type eventTypeEnum string

const (
	// Built in events only allowed to be fired by SDK
	eventTypeBuiltIn eventTypeEnum = "builtInEventType"
	// well known event type like (is signed in), that apps provide but SDK is aware of
	eventTypeWellKnown eventTypeEnum = "wellKnownEventType"
	// Competely custom events from app the SDK is not aware of
	eventTypeCustom eventTypeEnum = "customEventType"
)

const (
	buildInEventNamespace   = "io.criticalmoments.events.built_in."
	wellKnownEventNamespace = "io.criticalmoments.events.well_known."
)

// Enum type would be nice, but doesn't play well with gomobile exports
const (
	AppStartBuiltInEvent string = buildInEventNamespace + "app_start"
)

var (
	allBuiltInEventTypes = map[string]bool{
		AppStartBuiltInEvent: true,
	}
)

// Enum type would be nice, but doesn't play well with gomobile exports
const (
	SignedInEventAppStart string = wellKnownEventNamespace + "signed_in"
)

var (
	allWellKnownEventTypes = map[string]bool{
		SignedInEventAppStart: true,
	}
)

type Event struct {
	Name string

	// Event type is internal to cmcore
	eventType eventTypeEnum
}

func NewBuiltInEventWithName(name string) (*Event, error) {
	// Ensure this is a built in event we recognize
	if !allBuiltInEventTypes[name] {
		return nil, errors.New(fmt.Sprintf("Unknown built in event: %v", name))
	}
	if !strings.HasPrefix(name, buildInEventNamespace) {
		return nil, errors.New(fmt.Sprintf("Built in event outside namespace: %v", name))
	}

	e := Event{
		Name:      name,
		eventType: eventTypeBuiltIn,
	}
	return &e, nil
}

func NewWellKnownEventWithName(name string) (*Event, error) {
	// Ensure this is a well known event we recognize
	if !allWellKnownEventTypes[name] {
		return nil, errors.New(fmt.Sprintf("Unknown well known event: %v", name))
	}
	if !strings.HasPrefix(name, wellKnownEventNamespace) {
		return nil, errors.New(fmt.Sprintf("Well known event outside namespace: %v", name))
	}

	e := Event{
		Name:      name,
		eventType: eventTypeWellKnown,
	}
	return &e, nil
}

func NewCustomEventWithName(name string) (*Event, error) {
	if strings.HasPrefix(name, wellKnownEventNamespace) || strings.HasPrefix(name, buildInEventNamespace) {
		return nil, errors.New(fmt.Sprintf("Attempted to log custom event matching built in or well known event: %v", name))
	}

	e := Event{
		Name:      name,
		eventType: eventTypeCustom,
	}
	return &e, nil
}
