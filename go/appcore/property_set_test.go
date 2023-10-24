package appcore

import "testing"

func TestParsingProperSetFromJson(t *testing.T) {
	j := `{
		"stringKey": "stringVal",
		"boolKey": true,
		"intKey": 42,
		"floatKey": 3.3
	}`

	ps, err := newPropertySetFromJson([]byte(j))
	if err != nil {
		t.Fatal(err)
	}

	if ps.values["stringKey"] != "stringVal" {
		t.Fatal("Failed to parse string value")
	}
	if ps.values["boolKey"] != true {
		t.Fatal("Failed to parse bool value")
	}
	// Note: ints also parsed into float64
	if ps.values["intKey"] != 42.0 {
		t.Fatal("Failed to parse int value")
	}
	if ps.values["floatKey"] != 3.3 {
		t.Fatal("Failed to parse float value")
	}
}

func TestParsingPartialPropSetFromJson(t *testing.T) {
	j := `{
		"stringKey": "stringVal",
		"invalidKey": {},
		"boolKey": true,
		"intKey": 42,
		"floatKey": 3.3
	}`

	ps, err := newPropertySetFromJson([]byte(j))
	if err == nil {
		t.Fatal("Failed to throw error for invalid key")
	}

	if ps.values["stringKey"] != "stringVal" {
		t.Fatal("Failed to parse string value")
	}
	if ps.values["boolKey"] != true {
		t.Fatal("Failed to parse bool value")
	}
	// Note: ints also parsed into float64
	if ps.values["intKey"] != 42.0 {
		t.Fatal("Failed to parse int value")
	}
	if ps.values["floatKey"] != 3.3 {
		t.Fatal("Failed to parse float value")
	}

	// empty and invalid json objects
	ps, err = newPropertySetFromJson([]byte("{}"))
	if err != nil || len(ps.values) != 0 {
		t.Fatal("Failed on empty json")
	}
	ps, err = newPropertySetFromJson([]byte("{{{xcv"))
	if err == nil || len(ps.values) != 0 {
		t.Fatal("Didn't fail on invalid json")
	}
}
