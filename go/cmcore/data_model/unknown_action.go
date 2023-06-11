package datamodel

import "fmt"

type UnknownAction struct {
	ActionType string
}

func (u *UnknownAction) Validate() bool {
	return true
}

func (u *UnknownAction) ValidateReturningUserReadableIssue() string {
	return ""
}

func (u *UnknownAction) AllEmbeddedThemeNames() ([]string, error) {
	return []string{}, nil
}

func (u *UnknownAction) AllEmbeddedActionNames() ([]string, error) {
	return []string{}, nil
}

func (u *UnknownAction) PerformAction(ab ActionBindings) error {
	fmt.Printf("CriticalMoments: this version of critical moments does not support this action type (\"%v\"). Skipping action.", u.ActionType)
	return nil
}
