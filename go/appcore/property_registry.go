package appcore

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

const cmKindVersionNumber reflect.Kind = math.MaxInt

type propertyRegistry struct {
	providers              map[string]propertyProvider
	requiredPropertyTypes  map[string]reflect.Kind
	wellKnownPropertyTypes map[string]reflect.Kind
}

func newPropertyRegistry() *propertyRegistry {
	return &propertyRegistry{
		providers: make(map[string]propertyProvider),
		requiredPropertyTypes: map[string]reflect.Kind{
			"platform":              reflect.String,
			"os_version":            cmKindVersionNumber,
			"device_manufacturer":   reflect.String,
			"device_model":          reflect.String,
			"locale_language_code":  reflect.String,
			"locale_country_code":   reflect.String,
			"locale_currency_code":  reflect.String,
			"app_version":           cmKindVersionNumber,
			"user_interface_idiom":  reflect.String,
			"app_id":                reflect.String,
			"screen_width_pixels":   reflect.Int,
			"screen_height_pixels":  reflect.Int,
			"device_battery_state":  reflect.String,
			"device_battery_level":  reflect.Float64,
			"device_low_power_mode": reflect.Bool,
		},
		wellKnownPropertyTypes: map[string]reflect.Kind{
			"user_signed_in":       reflect.Bool,
			"device_model_class":   reflect.String,
			"device_model_version": cmKindVersionNumber,
			"screen_width_points":  reflect.Int,
			"screen_height_points": reflect.Int,
			"screen_scale":         reflect.Float64,
		},
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

	if expectedKind != cmKindVersionNumber {
		provider, ok = p.providers[propName]
	} else {
		// custom validation for version numbers, expect a string key with _string postfix
		provider, ok = p.providers[fmt.Sprintf(versionNumberStringKeyFormat, propName)]
		expectedKind = reflect.String
	}

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

const versionNumberStringKeyFormat = "%v_string"

func (p *propertyRegistry) registerStaticVersionNumberProperty(prefix string, versionString string) error {
	componentNames := []string{"major", "minor", "patch", "mini", "micro", "nano", "smol"}

	if prefix == "" {
		return errors.New("Prefix required for version property")
	}

	expectedType := p.expectedTypeForKey(prefix)
	if expectedType != cmKindVersionNumber {
		return errors.New("Not expecting a version number for key: " + prefix)
	}

	// Save string even if we can't parse the rest. Can target using exact strings worst case.
	stringProperty := staticPropertyProvider{
		value: versionString,
	}
	p.providers[fmt.Sprintf(versionNumberStringKeyFormat, prefix)] = &stringProperty

	components := strings.Split(versionString, ".")
	intComponents := make([]int, len(components))
	for i, component := range components {
		intComponent, err := strconv.Atoi(component)
		if err != nil {
			return errors.New(fmt.Sprintf("Invalid version number format: \"%v\"", versionString))
		}
		intComponents[i] = intComponent
	}

	for i := 0; i < len(intComponents) && i < len(componentNames); i++ {
		componentProperty := staticPropertyProvider{
			value: intComponents[i],
		}
		p.providers[fmt.Sprintf("%v_%v", prefix, componentNames[i])] = &componentProperty
	}

	return nil
}
