package appcore

import (
	"encoding/json"
	"errors"
	"fmt"
)

type propertySet struct {
	values map[string]interface{}
}

func newPropertySetFromJson(data []byte) (*propertySet, error) {
	values := map[string]interface{}{}

	var objmap map[string]interface{}
	err := json.Unmarshal(data, &objmap)
	if err != nil {
		return &propertySet{values: values}, err
	}

	// Parse our supported types
	for k, v := range objmap {
		boolVal, ok := v.(bool)
		if ok {
			values[k] = boolVal
			continue
		}

		// Note: ints also parsed into float64 in encoding/json
		floatVal, ok := v.(float64)
		if ok {
			values[k] = floatVal
			continue
		}

		stringVal, ok := v.(string)
		if ok {
			values[k] = stringVal
			continue
		}

		err = errors.Join(err, fmt.Errorf("unsupported type for key: %s", k))
	}

	return &propertySet{
		values: values,
	}, err
}
