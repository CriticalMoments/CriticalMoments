package datamodel

import (
	"testing"
)

func TestBuiltInEvent(t *testing.T) {
	e, err := NewBuiltInEventWithName(AppStartBuiltInEvent)
	if err != nil {
		t.Fatal()
	}
	if e.Name != string(AppStartBuiltInEvent) || e.EventType != EventTypeBuiltIn {
		t.Fatal()
	}
}

func TestInvalidBuiltInEvents(t *testing.T) {
	e, err := NewBuiltInEventWithName("not built in")
	if err == nil || e != nil {
		t.Fatal()
	}
	// well known != built in
	e, err = NewBuiltInEventWithName(SignedInEvent)
	if err == nil || e != nil {
		t.Fatal()
	}
}

func TestWellKnownEvent(t *testing.T) {
	e, err := NewWellKnownEventWithName(SignedInEvent)
	if err != nil {
		t.Fatal()
	}
	if e.Name != string(SignedInEvent) || e.EventType != EventTypeWellKnown {
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
	if e.Name != name || e.EventType != EventTypeCustom {
		t.Fatal()
	}
}

func TestSharedConstructor(t *testing.T) {
	e, err := NewClientEventWithName(AppStartBuiltInEvent)
	if err == nil || e != nil {
		t.Fatal("Build in should not be able to be fired by client")
	}

	e, err = NewClientEventWithName(SignedInEvent)
	if err != nil || e.EventType != EventTypeWellKnown || e.Name != SignedInEvent {
		t.Fatal("Failed to parse well known event from client")
	}

	e, err = NewClientEventWithName("net.scosman.hello")
	if err != nil || e.EventType != EventTypeCustom || e.Name != "net.scosman.hello" {
		t.Fatal("Failed to parse custom event")
	}
}
