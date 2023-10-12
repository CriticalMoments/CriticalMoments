package events

import (
	"errors"

	datamodel "github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model"
)

type EventManager struct {
	db *DB
}

func NewEventManager(dataDir string) (*EventManager, error) {
	db, err := NewDB(dataDir)
	if err != nil {
		return nil, err
	}

	return &EventManager{
		db: db,
	}, nil
}

func (em *EventManager) SendEvent(e *datamodel.Event) error {
	err := em.db.InsertEvent(e)
	if err != nil {
		return err
	}
	return nil
}

func (em *EventManager) EventManagerConditionFunctions() map[string]*datamodel.ConditionDynamicFunction {
	return map[string]*datamodel.ConditionDynamicFunction{
		"eventCount": {
			Function: func(params ...any) (any, error) {
				if len(params) != 1 {
					return nil, errors.New("eventCount requires one parameter")
				}
				eventName, ok := params[0].(string)
				if !ok {
					return nil, errors.New("eventCount requires a string parameter")
				}
				count, err := em.db.EventCountByName(eventName)
				if err != nil {
					return nil, err
				}
				return count, nil
			},
			Types: []any{new(func(string) int)},
		},
	}
}
