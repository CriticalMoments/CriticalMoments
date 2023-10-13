package events

import (
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
				// Parameter type+count checking is done with the Types signature
				count, err := em.db.EventCountByName(params[0].(string))
				if err != nil {
					return nil, err
				}
				return count, nil
			},
			Types: []any{new(func(string) int)},
		},
		"eventCountWithLimit": {
			Function: func(params ...any) (any, error) {
				// Parameter type+count checking is done the Types signature
				count, err := em.db.EventCountByNameWithLimit(params[0].(string), params[1].(int))
				if err != nil {
					return nil, err
				}
				return count, nil
			},
			Types: []any{new(func(string, int) int)},
		},
	}
}
