package datamodel

import (
	"fmt"
)

// error, but additional type we can check, and accessor to reason
type UserPresentableErrorInterface interface {
	Error() string
	UserReadableErrorString() string
}

// Implements `error`
type UserPresentableError struct {
	userReadableErrorString string
	SourceError             error
}

func NewUserPresentableError(s string) *UserPresentableError {
	return &UserPresentableError{
		userReadableErrorString: s,
	}
}

func NewUserPresentableErrorWSource(s string, sourceErr error) UserPresentableErrorInterface {
	return &UserPresentableError{
		userReadableErrorString: s,
		SourceError:             sourceErr,
	}
}

func NewUserErrorForJsonIssue(data []byte, sourceErr error) UserPresentableErrorInterface {
	jsonString := string(data)
	// truncated to max 600 characters, adding ... to the end
	if len(jsonString) > 600 {
		jsonString = jsonString[:600] + "... [truncated]"
	}
	return NewUserPresentableErrorWSource(fmt.Sprintf("Error parsing your config in the following section. See description of the source error below:\n%v", jsonString), sourceErr)
}

func (err *UserPresentableError) Error() string {
	if err.SourceError == nil {
		return err.userReadableErrorString
	}
	return fmt.Sprintf("%v\n  Source Error: %v", err.userReadableErrorString, err.SourceError)
}

func (err *UserPresentableError) UserReadableErrorString() string {
	return err.userReadableErrorString
}
