package appcore

import (
	"fmt"
	"reflect"
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
