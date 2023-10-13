package appcore

import (
	"errors"
	"fmt"
	"reflect"

	datamodel "github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model"
	"github.com/antonmedv/expr"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type propertyRegistry struct {
	providers              map[string]propertyProvider
	requiredPropertyTypes  map[string]reflect.Kind
	wellKnownPropertyTypes map[string]reflect.Kind
	dynamicFunctionNames   []string
	dynamicFunctionOps     []expr.Option
	mapFunctions           map[string]interface{}
}

func newPropertyRegistry() *propertyRegistry {
	pr := &propertyRegistry{
		providers:              make(map[string]propertyProvider),
		requiredPropertyTypes:  datamodel.RequiredPropertyTypes(),
		wellKnownPropertyTypes: datamodel.WellKnownPropertyTypes(),
		dynamicFunctionNames:   []string{},
		dynamicFunctionOps:     []expr.Option{},
	}

	// register static/map functions
	pr.mapFunctions = datamodel.ConditionEnvWithHelpers()

	return pr
}

func (pr *propertyRegistry) RegisterDynamicFunctions(newFuncs map[string]*datamodel.ConditionDynamicFunction) error {
	for k, v := range newFuncs {
		pr.dynamicFunctionNames = append(pr.dynamicFunctionNames, k)
		pr.dynamicFunctionOps = append(pr.dynamicFunctionOps, expr.Function(k, v.Function, v.Types...))
	}
	return nil
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

var errPropertyNotFound = errors.New("Property not found")

func (p *propertyRegistry) propertyValue(key string) (interface{}, error) {
	v, ok := p.providers[key]
	if !ok {
		return nil, errPropertyNotFound
	}
	return v.Value(), nil
}

func (p *propertyRegistry) buildPropertyMapForCondition(fields *datamodel.ConditionFields) (map[string]interface{}, error) {
	// Extract only the used variables from the condition. Property evaluation isn't free, so
	// only evaluate those we need
	propsEnv := make(map[string]interface{})
	for _, v := range fields.Variables {
		if _, ok := propsEnv[v]; !ok {
			pv, err := p.propertyValue(v)
			if err != nil && err != errPropertyNotFound {
				return nil, err
			}
			if err == errPropertyNotFound {
				// set not-found variables to nil. Likely new var names from future SDK runing on an old SDK.
				// We want the condition string to be able to check for nil for backwards compatibility (typically "?? true" or "?? false")
				propsEnv[v] = nil
			} else {
				propsEnv[v] = pv
			}
		}
	}
	return propsEnv, nil
}

// Any unrecoginized method should return nil (not the default error)
// This is because we want to allow for backwards compatibility when newer SDKs add functions (old SDKs shouldn't fail, should return nil)
func (p *propertyRegistry) nilMethodsForUnknownFunctions(fields *datamodel.ConditionFields) ([]expr.Option, error) {
	existingFunctions := p.allFunctionNamesRegistered()
	nilFunctions := []expr.Option{}
	for _, m := range fields.Methods {
		if !slices.Contains(existingFunctions, m) {
			nfunc := expr.Function(m, func(params ...any) (interface{}, error) {
				return nil, nil
			})
			nilFunctions = append(nilFunctions, nfunc)
		}
	}

	return nilFunctions, nil
}

func (p *propertyRegistry) allFunctionNamesRegistered() []string {
	functions := []string{}
	functions = append(functions, maps.Keys(p.mapFunctions)...)
	functions = append(functions, p.dynamicFunctionNames...)

	return functions
}

func (p *propertyRegistry) evaluateCondition(condition *datamodel.Condition) (bool, error) {
	// Parse the condition, extract variable and method names
	fields, err := condition.ExtractIdentifiers()
	if err != nil {
		return false, err
	}

	// Build a map of all properties(variables) used in this condition, and their values
	envMap, err := p.buildPropertyMapForCondition(fields)
	if err != nil {
		return false, err
	}

	// Add all the static functions to the envidonment map
	maps.Copy(envMap, p.mapFunctions)

	// Build nil function handlers for any missing functions (backwards compatibility)
	nilOps, err := p.nilMethodsForUnknownFunctions(fields)
	if err != nil {
		return false, err
	}

	mergedOptions := []expr.Option{}
	mergedOptions = append(mergedOptions, p.dynamicFunctionOps...)
	mergedOptions = append(mergedOptions, expr.Env(envMap))
	mergedOptions = append(mergedOptions, nilOps...)

	program, err := condition.CompileWithEnv(mergedOptions...)
	if err != nil {
		return false, err
	}
	result, err := expr.Run(program, envMap)
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
	provider, ok := p.providers[propName]

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
