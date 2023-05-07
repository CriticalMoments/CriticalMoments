package testing

import (
	"testing"

	"github.com/CriticalMoments/CriticalMoments/go/appcore"
)

func TestBuiltInEvent(t *testing.T) {
	e, err := appcore.NewBuiltInEventWithName(string(appcore.AppStartBuiltInEvent))
	if err != nil {
		t.Fatal()
	}
	if e.Name != string(appcore.AppStartBuiltInEvent) || e.EventType != appcore.EventTypeBuiltIn {
		t.Fatal()
	}
}

func TestInvalidBuiltInEvents(t *testing.T) {
	e, err := appcore.NewBuiltInEventWithName("not built in")
	if err == nil || e != nil {
		t.Fatal()
	}
	// well known != built in
	e, err = appcore.NewBuiltInEventWithName(string(appcore.SignedInEventAppStart))
	if err == nil || e != nil {
		t.Fatal()
	}
}

func TestWellKnownEvent(t *testing.T) {
	e, err := appcore.NewWellKnownEventWithName(string(appcore.SignedInEventAppStart))
	if err != nil {
		t.Fatal()
	}
	if e.Name != string(appcore.SignedInEventAppStart) || e.EventType != appcore.EventTypeWellKnown {
		t.Fatal()
	}
}

func TestInvalidWellKnownEvents(t *testing.T) {
	e, err := appcore.NewWellKnownEventWithName("not well known")
	if err == nil || e != nil {
		t.Fatal()
	}
	// well known != built in
	e, err = appcore.NewWellKnownEventWithName(string(appcore.AppStartBuiltInEvent))
	if err == nil || e != nil {
		t.Fatal()
	}
}
func TestCustomEventEvent(t *testing.T) {
	name := "net.scosman.built_thing"
	e, err := appcore.NewCustomEventWithName(name)
	if err != nil {
		t.Fatal()
	}
	if e.Name != name || e.EventType != appcore.EventTypeCustom {
		t.Fatal()
	}
}

func TestInvalidCustomEvents(t *testing.T) {
	// Built in shouldn't work in custom
	e, err := appcore.NewCustomEventWithName(string(appcore.AppStartBuiltInEvent))
	if err == nil || e != nil {
		t.Fatal()
	}
	// Well known shouldn't work in custom
	e, err = appcore.NewCustomEventWithName(string(appcore.SignedInEventAppStart))
	if err == nil || e != nil {
		t.Fatal()
	}
	// 2x namespace errors
	e, err = appcore.NewCustomEventWithName("io.criticalmoments.events.built_in.custom")
	if err == nil || e != nil {
		t.Fatal()
	}
	e, err = appcore.NewCustomEventWithName("io.criticalmoments.events.well_known.custom")
	if err == nil || e != nil {
		t.Fatal()
	}
}
