package datamodel

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"
)

func TestJsonErrorFormatting(t *testing.T) {
	var cases = map[string]string{
		"invalid1.json": "JSON Parsing Error on line '4' (offset 1): invalid character '}'",
		"invalid2.json": "JSON Parsing Error on line '1' (offset 6): invalid character ':'", // add_test_case
		// valid json but not a valid config
		"invalid3.json":      "Config must have a config version of",                                    // add_test_case
		"invalid_empty.json": "JSON Parsing Error on line '1' (offset 0): unexpected end of JSON input", // add_test_case
	}

	for file, expectedErr := range cases {
		testFileData, err := os.ReadFile("./test/testdata/json_error/" + file)
		if err != nil {
			t.Fatal(err)
		}
		var pc PrimaryConfig
		err = json.Unmarshal(testFileData, &pc)
		if err == nil {
			t.Fatal("Failed to error on invalid json")
		}
		userFriendlyErr := UserFriendlyJsonError(err, testFileData)
		if userFriendlyErr == nil {
			t.Fatal("Failed to error on invalid json")
		}
		if _, ok := userFriendlyErr.(*UserPresentableError); !ok {
			t.Fatalf("Failed to parse error message in file %s.\nExpected UserPresentableError\nGot '%s'\n", file, userFriendlyErr.Error())
		}
		if !strings.Contains(userFriendlyErr.Error(), expectedErr) {
			t.Fatalf("Failed to parse error message in file %s.\nExpected '%s'\nGot '%s'\n", file, expectedErr, userFriendlyErr.Error())
		}
	}
}

func TestJsonErrorPassthrough(t *testing.T) {
	err := errors.New("test error")
	userFriendlyErr := UserFriendlyJsonError(err, nil)
	if userFriendlyErr == nil {
		t.Fatal("Failed to error on invalid json")
	}
	if userFriendlyErr.Error() != "test error" {
		t.Fatal("Did not pass through non-json error")
	}
}
