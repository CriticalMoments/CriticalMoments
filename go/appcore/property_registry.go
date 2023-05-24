package appcore

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/CriticalMoments/CriticalMoments/go/cmcore"
	"github.com/antonmedv/expr"
	"golang.org/x/exp/slices"
)

type propertyRegistry struct {
	providers              map[string]propertyProvider
	requiredPropertyTypes  map[string]reflect.Kind
	wellKnownPropertyTypes map[string]reflect.Kind
}

func newPropertyRegistry() *propertyRegistry {
	return &propertyRegistry{
		providers:              make(map[string]propertyProvider),
		requiredPropertyTypes:  cmcore.RequiredPropertyTypes(),
		wellKnownPropertyTypes: cmcore.WellKnownPropertyTypes(),
	}
}

func (pr *propertyRegistry) expectedTypeForKey(key string) reflect.Kind {
	expectedType, foundType := pr.requiredPropertyTypes[key]
	if foundType {
		return expectedType
	}
	expectedType, foundType = pr.wellKnownPropertyTypes[key]
	if foundType {
		return expectedType
	}
	return reflect.Invalid
}

func (pr *propertyRegistry) addProviderForKey(key string, pp propertyProvider) error {
	_, hasCurrent := pr.providers[key]
	if hasCurrent {
		fmt.Println("CriticalMoments Warning: Re-registering property provider for key: " + key)
	}

	expectedType := pr.expectedTypeForKey(key)
	if expectedType == reflect.Invalid {
		return errors.New("Invalid property registered. Properties must be required or well known. Arbitrary properties are not allowed.")
	}

	// cmKindVersionNumber is also valid, but is validated below
	validTypes := []reflect.Kind{reflect.Bool, reflect.String, reflect.Int, reflect.Float64}
	if !slices.Contains(validTypes, expectedType) {
		return errors.New("Invalid property type for key: " + key)
	}

	if pp.Kind() != expectedType {
		return errors.New("Property registered of wrong type (does not match expected type): " + key)
	}

	pr.providers[key] = pp
	return nil
}

func (p *propertyRegistry) registerStaticProperty(key string, value interface{}) error {
	s := staticPropertyProvider{
		value: value,
	}
	return p.addProviderForKey(key, &s)
}

func (p *propertyRegistry) registerLibPropertyProvider(key string, dpp LibPropertyProvider) error {
	dw := newLibPropertyProviderWrapper(dpp)
	return p.addProviderForKey(key, &dw)
}

func (p *propertyRegistry) propertyValue(key string) interface{} {
	v, ok := p.providers[key]
	if !ok {
		return nil
	}
	return v.Value()
}

func (p *propertyRegistry) evaluateCondition(condition string) (bool, error) {
	variables, err := cmcore.ExtractVariablesFromCondition(condition)
	if err != nil {
		return false, err
	}

	// Build env with helper functions and vars from props
	env := cmcore.ConditionEnvWithHelpers()
	for _, v := range variables {
		if _, ok := env[v]; !ok {
			env[v] = p.propertyValue(v)
		}
	}

	// TODO functions not bound here. bind to cmExprEnv if we add function support
	program, err := expr.Compile(condition, expr.Env(env), expr.AllowUndefinedVariables(), expr.AsBool())
	if err != nil {
		return false, err
	}
	result, err := expr.Run(program, env)
	if err != nil {
		return false, err
	}
	boolResult, ok := result.(bool)
	if !ok {
		return false, nil
	}
	return boolResult, nil
}

func (p *propertyRegistry) validateProperties() error {
	// Check required
	for propName, expectedKind := range p.requiredPropertyTypes {
		err := p.validateExpectedProvider(propName, expectedKind, false)
		if err != nil {
			return err
		}
	}

	// check well known
	for propName, expectedKind := range p.wellKnownPropertyTypes {
		err := p.validateExpectedProvider(propName, expectedKind, true)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *propertyRegistry) validateExpectedProvider(propName string, expectedKind reflect.Kind, allowMissing bool) error {
	var provider propertyProvider
	var ok bool

	provider, ok = p.providers[propName]

	if !ok && !allowMissing {
		return errors.New(fmt.Sprintf("Missing required property: %v", propName))
	}
	if !ok && allowMissing {
		return nil
	}
	if provider.Kind() != expectedKind {
		return errors.New(fmt.Sprintf("Property \"%v\" of wrong kind. Expected %v", propName, expectedKind.String()))
	}
	return nil
}
