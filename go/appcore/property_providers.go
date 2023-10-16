package appcore

import (
	"fmt"
	"math"
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

const LibPropertyProviderNilStringValue = "io.criticalmoments.libpropertyprovider.nilstringvalue"
const LibPropertyProviderNilFloatValue = -math.MaxFloat64
const LibPropertyProviderNilIntValue = math.MinInt64

// An interface libraries can implement to provide dynamic properties.
// Not ideal interface in go, but gomobile won't map interface{}, reflect.Kind, or enum types
type LibPropertyProvider interface {
	Type() int
	IntValue() int64
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
		v := d.propertyProvider.FloatValue()
		if v == LibPropertyProviderNilFloatValue {
			return nil
		}
		return v
	case LibPropertyProviderTypeInt:
		v := d.propertyProvider.IntValue()
		if v == LibPropertyProviderNilIntValue {
			return nil
		}
		return v
	case LibPropertyProviderTypeString:
		v := d.propertyProvider.StringValue()
		if v == LibPropertyProviderNilStringValue {
			return nil
		}
		return v
	}
	fmt.Println("CriticalMoments: Invalid property type!")
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
	fmt.Println("CriticalMoments: Invalid property type!")
	return reflect.String
}
