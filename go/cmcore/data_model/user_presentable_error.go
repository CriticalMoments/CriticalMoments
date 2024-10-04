package datamodel

import (
	"fmt"
)

// error, but additional type we can check, and accessor to reason
type UserPresentableErrorInterface interface {
	Error() string
	UserReadableErrorString() string
}

// Implements `error` and `UserPresentableErrorInterface`
type userPresentableError struct {
	userReadableErrorString string
	SourceError             error
}

func NewUserPresentableError(s string) UserPresentableErrorInterface {
	return &userPresentableError{
		userReadableErrorString: s,
	}
}

func NewUserPresentableErrorWSource(s string, sourceErr error) UserPresentableErrorInterface {
	return &userPresentableError{
		userReadableErrorString: s,
		SourceError:             sourceErr,
	}
}

func (err *userPresentableError) Error() string {
	if err.SourceError == nil {
		return err.userReadableErrorString
	}
	// source and current inverted. The root is really the most important.
	return fmt.Sprintf("%v\n    Source Error: %v", err.SourceError, err.userReadableErrorString)
}

func (err *userPresentableError) UserReadableErrorString() string {
	return err.Error()
}

// Implements `error` and `userPresentableErrorInterface`
type jsonUserPresentableError struct {
	SourceError error
	jsonSource  string
}

func NewUserErrorForJsonIssue(data []byte, sourceErr error) UserPresentableErrorInterface {
	jsonString := string(data)
	// truncated to max 600 characters, adding ... to the end
	if len(jsonString) > 600 {
		jsonString = jsonString[:600] + "... [truncated]"
	}
	return &jsonUserPresentableError{
		SourceError: sourceErr,
		jsonSource:  jsonString,
	}
}

func (err *jsonUserPresentableError) Error() string {
	if err.SourceError == nil {
		return fmt.Sprintf("JSON Parsing Error in json: %v", err.jsonSource)
	}
	if err.jsonSource == "" {
		return err.SourceError.Error()
	}
	return fmt.Sprintf("%v\n    JSON section this error occurred in: %v", err.SourceError.Error(), err.jsonSource)
}

func (err *jsonUserPresentableError) UserReadableErrorString() string {
	return err.Error()
}
