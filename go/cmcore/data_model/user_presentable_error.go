package datamodel

import "fmt"

// Interface we can check against, to see if error is okay to present to user
// returns from this method should be plain english, and not refer to internals,
// only the errors of inputs the user controls
type UserPresentableErrorI interface {
	UserErrorString() string
}

// Implements UserPresentableErrorI and `error`
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

func (err *UserPresentableError) UserErrorString() string {
	sourcePresentableError, ok := interface{}(err.SourceError).(UserPresentableErrorI)
	if ok {
		return fmt.Sprintf("%v (from error: \"%v\")", err.userReadableErrorString, sourcePresentableError.UserErrorString())
	} else {
		return err.userReadableErrorString
	}
}

func (err *UserPresentableError) Error() string {
	if err.SourceError == nil {
		return err.userReadableErrorString
	}
	return fmt.Sprintf("%v (Source Error: %v)", err.userReadableErrorString, err.SourceError)
}
