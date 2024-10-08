package datamodel

import "fmt"

type UnknownAction struct {
	ActionType string
}

func (u *UnknownAction) Check() UserPresentableErrorInterface {
	return nil
}

func (u *UnknownAction) AllEmbeddedThemeNames() ([]string, error) {
	return []string{}, nil
}

func (u *UnknownAction) AllEmbeddedActionNames() ([]string, error) {
	return []string{}, nil
}

func (u *UnknownAction) AllEmbeddedConditions() ([]*Condition, error) {
	return []*Condition{}, nil
}

func (u *UnknownAction) PerformAction(ab ActionBindings, actionName string) error {
	return fmt.Errorf("this version of critical moments does not support this action type (\"%v\"). \"%v\" action not performed", u.ActionType, actionName)
}
