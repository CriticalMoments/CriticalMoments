package appcore

import (
	"fmt"
	"reflect"
	"time"

	datamodel "github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model"
)

type EventManager struct {
	lastSessionStartTime *time.Time
}

func (em *EventManager) SendEvent(e *datamodel.Event, ac *Appcore) error {
	em.processEvent(e, ac)
	err := ac.db.InsertEvent(e)
	if err != nil {
		return err
	}
	return nil
}

func (em *EventManager) processEvent(e *datamodel.Event, ac *Appcore) {
	if e.EventType == datamodel.EventTypeBuiltIn && e.Name == datamodel.AppEnteredForegroundBuiltInEvent {
		err := em.updateSessionForForeground(ac)
		if err != nil {
			fmt.Printf("CriticalMoments: Error processing enter foreground event: %v", err)
		}
		return
	}
}

var SessionGapDuration = time.Minute * 10

// Update session when entering foreground
func (em *EventManager) updateSessionForForeground(ac *Appcore) error {
	// Fail fast: if session started in last 10 minutes, we know we're still in that session
	if em.lastSessionStartTime != nil {
		if time.Since(*em.lastSessionStartTime) < SessionGapDuration {
			// Continue current sessions
			return nil
		}
	}

	lastEnterBackgroundTime, err := ac.db.LatestEventTimeByName(datamodel.AppEnteredBackgroundBuiltInEvent)
	if err != nil {
		return err
	}

	if lastEnterBackgroundTime == nil || time.Since(*lastEnterBackgroundTime) > SessionGapDuration {
		// No prior backgrounds or it's been 10 mins, this is new session
		return em.startSession(ac)
	}

	// If we've never entered foreground (besides now), start a session.
	// The test harness does send "enter background" before ever sending foreground, so this is real case
	lastEnterForegroundTime, err := ac.db.LatestEventTimeByName(datamodel.AppEnteredForegroundBuiltInEvent)
	if err != nil {
		return err
	}
	if lastEnterForegroundTime == nil {
		return em.startSession(ac)
	}

	return nil
}

func (em *EventManager) startSession(ac *Appcore) error {
	// Start new session: fire event, set session_start_time, and remember last timestamp
	err := ac.SendBuiltInEvent(datamodel.SessionStartBuiltInEvent)
	if err != nil {
		return err
	}
	now := time.Now()
	em.lastSessionStartTime = &now

	return nil
}

// Property provider for session start time

type SessionStartTimePropertyProvider struct {
	eventManager *EventManager
}

func (s SessionStartTimePropertyProvider) Value() interface{} {
	if s.eventManager.lastSessionStartTime == nil {
		return time.Now()
	}
	return *s.eventManager.lastSessionStartTime
}

func (s SessionStartTimePropertyProvider) Kind() reflect.Kind {
	return datamodel.CMTimeKind
}
