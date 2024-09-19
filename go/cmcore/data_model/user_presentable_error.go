package datamodel

import "fmt"

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

func NewUserPresentableErrorWSource(s string, sourceErr error) *UserPresentableError {
	return &UserPresentableError{
		userReadableErrorString: s,
		SourceError:             sourceErr,
	}
}

func (err *UserPresentableError) Error() string {
	if err.SourceError == nil {
		return err.userReadableErrorString
	}
	return fmt.Sprintf("%v (Source Error: %v)", err.userReadableErrorString, err.SourceError)
}
