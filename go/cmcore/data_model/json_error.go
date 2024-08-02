package datamodel

import (
	"encoding/json"
	"fmt"
)

func UserFriendlyJsonError(err error, jsonData []byte) error {
	if jsonErr, ok := err.(*json.SyntaxError); ok {
		line := 1
		var lineStartOffset int64 = 0
		for i, b := range jsonData {
			if b == '\n' {
				line++
				lineStartOffset = int64(i) + 1
			}
			if int64(i) == jsonErr.Offset {
				break
			}
		}

		errString := fmt.Sprintf("JSON Parsing Error on line '%d' (offset %d): %s", line, jsonErr.Offset-lineStartOffset, err.Error())
		return NewUserPresentableError(errString)
	}

	return err
}
