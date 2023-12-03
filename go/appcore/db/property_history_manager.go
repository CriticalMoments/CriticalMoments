package db

import (
	"errors"

	datamodel "github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model"
)

type propHistoryValue struct {
	value       interface{}
	sample_type datamodel.CMPropertySampleType
}

type PropertyHistoryManager struct {
	db *DB

	preStartPropsCache map[string]propHistoryValue
}

func newPropertyHistoryManager(db *DB) *PropertyHistoryManager {
	return &PropertyHistoryManager{
		db:                 db,
		preStartPropsCache: map[string]propHistoryValue{},
	}
}

func (phm *PropertyHistoryManager) TrackPropertyHistoryForStartup(appStartValues map[string]interface{}) error {
	// keep processing on error, but return all errors at end
	errorSet := []error{}

	// Set values for properties that were set before startup
	for name, prop := range phm.preStartPropsCache {
		err := phm.db.InsertPropertyHistory(name, prop.value, prop.sample_type)
		if err != nil {
			errorSet = append(errorSet, err)
		}
	}
	phm.preStartPropsCache = map[string]propHistoryValue{}

	// Set the startup values (used for built in props with sample type= CMPropertySampleTypeAppStart)
	for name, val := range appStartValues {
		err := phm.db.InsertPropertyHistory(name, val, datamodel.CMPropertySampleTypeAppStart)
		if err != nil {
			errorSet = append(errorSet, err)
		}
	}

	if len(errorSet) > 0 {
		return errors.Join(errorSet...)
	}

	return nil
}

func (phm *PropertyHistoryManager) CustomPropertySet(name string, val interface{}) error {
	return phm.setPropertyHistory(name, val, datamodel.CMPropertySampleTypeOnCustomSet)
}

func (phm *PropertyHistoryManager) UpdateHistoryForPropertyAccessed(name string, val interface{}) error {
	return phm.setPropertyHistory(name, val, datamodel.CMPropertySampleTypeOnUse)
}

func (phm *PropertyHistoryManager) setPropertyHistory(name string, val interface{}, sampleType datamodel.CMPropertySampleType) error {
	if name == "" {
		return nil
	}

	if phm.db.started {
		err := phm.db.InsertPropertyHistory(name, val, datamodel.CMPropertySampleTypeOnCustomSet)
		if err != nil {
			return err
		}
	} else {
		// Not started, cache
		phm.preStartPropsCache[name] = propHistoryValue{
			value:       val,
			sample_type: sampleType,
		}
	}

	return nil
}
