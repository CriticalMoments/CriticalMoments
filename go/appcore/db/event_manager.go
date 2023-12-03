package db

import (
	datamodel "github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model"
)

type EventManager struct {
	db *DB
}

func (em *EventManager) SendEvent(e *datamodel.Event) error {
	err := em.db.InsertEvent(e)
	if err != nil {
		return err
	}
	return nil
}
