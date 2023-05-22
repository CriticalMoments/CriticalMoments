package appcore

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type propertyProvider interface {
	Value() interface{}
	Kind() reflect.Kind
}

// Set once properties
type staticPropertyProvider struct {
	value interface{}
}

func (s *staticPropertyProvider) Value() interface{} {
	return s.value
}

func (s *staticPropertyProvider) Kind() reflect.Kind {
	return reflect.TypeOf(s.value).Kind()
}

// An interface libraries can implement to provide dynamic properties.
// Not ideal interface in go, but gomobile won't map interface{}, reflect.Kind, or enum types
type LibPropertyProvider interface {
	Type() int
	IntValue() int
	StringValue() string
	FloatValue() float64
	BoolValue() bool
}

const (
	LibPropertyProviderTypeString int = iota
	LibPropertyProviderTypeInt
	LibPropertyProviderTypeFloat
	LibPropertyProviderTypeBool
)

func newLibPropertyProviderWrapper(dpp LibPropertyProvider) dynamicPropertyProviderWrapper {
	return dynamicPropertyProviderWrapper{
		propertyProvider: dpp,
	}
}

type dynamicPropertyProviderWrapper struct {
	propertyProvider LibPropertyProvider
}

func (d *dynamicPropertyProviderWrapper) Value() interface{} {
	switch d.propertyProvider.Type() {
	case LibPropertyProviderTypeBool:
		return d.propertyProvider.BoolValue()
	case LibPropertyProviderTypeFloat:
		return d.propertyProvider.FloatValue()
	case LibPropertyProviderTypeInt:
		return d.propertyProvider.IntValue()
	case LibPropertyProviderTypeString:
		return d.propertyProvider.StringValue()
	}
	fmt.Println("Invalid property type!")
	return nil
}

func (d *dynamicPropertyProviderWrapper) Kind() reflect.Kind {
	switch d.propertyProvider.Type() {
	case LibPropertyProviderTypeBool:
		return reflect.Bool
	case LibPropertyProviderTypeFloat:
		return reflect.Float64
	case LibPropertyProviderTypeInt:
		return reflect.Int
	case LibPropertyProviderTypeString:
		return reflect.String
	}
	fmt.Println("Invalid property type!")
	return reflect.String
}

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
			"os_version_string":     reflect.String,
			"device_manufacturer":   reflect.String,
			"device_model":          reflect.String,
			"locale_language_code":  reflect.String,
			"locale_country_code":   reflect.String,
			"locale_currency_code":  reflect.String,
			"app_version_string":    reflect.String,
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
			"device_model_version": reflect.Float64,
			"screen_width_points":  reflect.Int,
			"screen_height_points": reflect.Int,
			"screen_scale":         reflect.Float64,
		},
	}
}

func (p *propertyRegistry) registerStaticProperty(key string, value interface{}) {
	// TODO Block type mismatches
	// TODO block not well known or required until we add "custom" properties, store those separately
	// figure out how a mismatch doesn't create a type error if user uses a key we later decide is
	// well known or required
	s := staticPropertyProvider{
		value: value,
	}
	p.providers[key] = &s
}

func (p *propertyRegistry) registerLibPropertyProvider(key string, dpp LibPropertyProvider) {
	// TODO: same comments above
	dw := newLibPropertyProviderWrapper(dpp)
	p.providers[key] = &dw
}

func (p *propertyRegistry) propertyValue(key string) interface{} {
	v, ok := p.providers[key]
	if !ok {
		return nil
	}
	return v.Value()
}

func (p *propertyRegistry) validatePropertiesReturningUserReadable() string {
	// TODO: same comments above
	// Check required
	for propName, expectedKind := range p.requiredPropertyTypes {
		provider, ok := p.providers[propName]
		if !ok {
			return fmt.Sprintf("Missing required property: %v", propName)
		}
		if provider.Kind() != expectedKind {
			return fmt.Sprintf("Property \"%v\" of wrong kind. Expected %v", propName, expectedKind.String())
		}
	}

	// check well known
	for propName, expectedKind := range p.wellKnownPropertyTypes {
		provider, ok := p.providers[propName]
		if !ok {
			// missing is okay for well known, they are not required
			continue
		}
		if provider.Kind() != expectedKind {
			return fmt.Sprintf("Property \"%v\" of wrong kind. Expected %v", propName, expectedKind.String())
		}
	}

	return ""
}

func (p *propertyRegistry) registerStaticVersionNumberProperty(prefix string, versionString string) error {
	componentNames := []string{"major", "minor", "patch", "mini", "micro", "nano", "smol"}

	if prefix == "" {
		return errors.New("Prefix required for version property")
	}

	// Save string even if we can't parse the rest. Can target using exact strings worst case.
	stringProperty := staticPropertyProvider{
		value: versionString,
	}
	p.providers[fmt.Sprintf("%v_version_string", prefix)] = &stringProperty

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
		p.providers[fmt.Sprintf("%v_version_%v", prefix, componentNames[i])] = &componentProperty
	}

	return nil
}
