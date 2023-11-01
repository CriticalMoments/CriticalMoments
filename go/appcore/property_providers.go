package appcore

import (
	"fmt"
	"math"
	"reflect"
	"time"

	datamodel "github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model"
	"golang.org/x/exp/slices"
)

type propertyProvider interface {
	Value() interface{}
	Kind() reflect.Kind
}

var validPropertyTypes = []reflect.Kind{
	reflect.Bool,
	reflect.String,
	reflect.Int,
	reflect.Float64,
	datamodel.CMTimeKind,
}

func typeFromValue(v interface{}) reflect.Kind {
	if v == nil {
		return reflect.Invalid
	}
	if _, ok := v.(time.Time); ok {
		return datamodel.CMTimeKind
	}
	k := reflect.TypeOf(v).Kind()
	if slices.Contains(validPropertyTypes, k) {
		return k
	}
	return reflect.Invalid
}

// Set once properties
type staticPropertyProvider struct {
	value interface{}
}

func (s *staticPropertyProvider) Value() interface{} {
	return s.value
}

func (s *staticPropertyProvider) Kind() reflect.Kind {
	return typeFromValue(s.value)
}

// Nil/pointers not a native type for go-bind, so define constants for nil values.
// SDK wrapper layer should hide these from the consumer.
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
	TimeEpochMilliseconds() int64
	BoolValue() bool
}

const (
	LibPropertyProviderTypeString int = iota
	LibPropertyProviderTypeInt
	LibPropertyProviderTypeFloat
	LibPropertyProviderTypeTime
	LibPropertyProviderTypeBool
)

func newLibPropertyProviderWrapper(dpp LibPropertyProvider) *dynamicPropertyProviderWrapper {
	return &dynamicPropertyProviderWrapper{
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
	case LibPropertyProviderTypeTime:
		ems := d.propertyProvider.TimeEpochMilliseconds()
		if ems == LibPropertyProviderNilIntValue {
			return nil
		}
		return time.UnixMilli(ems)
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
	case LibPropertyProviderTypeTime:
		return datamodel.CMTimeKind
	}
	fmt.Println("CriticalMoments: Invalid property type!")
	return reflect.String
}
