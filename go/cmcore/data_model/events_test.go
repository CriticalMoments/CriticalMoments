package datamodel

import (
	"testing"
)

func TestBuiltInEvent(t *testing.T) {
	e, err := NewBuiltInEventWithName(AppStartBuiltInEvent)
	if err != nil {
		t.Fatal()
	}
	if e.Name != string(AppStartBuiltInEvent) || e.eventType != eventTypeBuiltIn {
		t.Fatal()
	}
}

func TestInvalidBuiltInEvents(t *testing.T) {
	e, err := NewBuiltInEventWithName("not built in")
	if err == nil || e != nil {
		t.Fatal()
	}
	// well known != built in
	e, err = NewBuiltInEventWithName(SignedInEventAppStart)
	if err == nil || e != nil {
		t.Fatal()
	}
}

func TestWellKnownEvent(t *testing.T) {
	e, err := NewWellKnownEventWithName(SignedInEventAppStart)
	if err != nil {
		t.Fatal()
	}
	if e.Name != string(SignedInEventAppStart) || e.eventType != eventTypeWellKnown {
		t.Fatal()
	}
}

func TestInvalidWellKnownEvents(t *testing.T) {
	e, err := NewWellKnownEventWithName("not well known")
	if err == nil || e != nil {
		t.Fatal()
	}
	// well known != built in
	e, err = NewWellKnownEventWithName(AppStartBuiltInEvent)
	if err == nil || e != nil {
		t.Fatal()
	}
}
func TestCustomEventEvent(t *testing.T) {
	name := "net.scosman.built_thing"
	e, err := NewCustomEventWithName(name)
	if err != nil {
		t.Fatal()
	}
	if e.Name != name || e.eventType != eventTypeCustom {
		t.Fatal()
	}
}

func TestInvalidCustomEvents(t *testing.T) {
	// Built in shouldn't work in custom
	e, err := NewCustomEventWithName(AppStartBuiltInEvent)
	if err == nil || e != nil {
		t.Fatal()
	}
	// Well known shouldn't work in custom
	e, err = NewCustomEventWithName(SignedInEventAppStart)
	if err == nil || e != nil {
		t.Fatal()
	}
	// 2x namespace errors
	e, err = NewCustomEventWithName("io.criticalmoments.events.built_in.custom")
	if err == nil || e != nil {
		t.Fatal()
	}
	e, err = NewCustomEventWithName("io.criticalmoments.events.well_known.custom")
	if err == nil || e != nil {
		t.Fatal()
	}
}
