package datamodel

import (
	"math"
	"reflect"
	"time"

	"golang.org/x/exp/slices"
)

const CMTimeKind = (reflect.Kind)(math.MaxUint)

var ValidPropertyTypes = []reflect.Kind{
	reflect.Bool,
	reflect.String,
	reflect.Int,
	reflect.Float64,
	CMTimeKind,
}

func CMTypeFromValue(v interface{}) reflect.Kind {
	if v == nil {
		return reflect.Invalid
	}
	if _, ok := v.(time.Time); ok {
		return CMTimeKind
	}
	k := reflect.TypeOf(v).Kind()
	if slices.Contains(ValidPropertyTypes, k) {
		return k
	}
	return reflect.Invalid
}

type CMPropertySampleType int

const (
	CMPropertySampleTypeAppStart    CMPropertySampleType = 1
	CMPropertySampleTypeOnUse       CMPropertySampleType = 2
	CMPropertySampleTypeOnCustomSet CMPropertySampleType = 3
	CMPropertySampleTypeDoNotSample CMPropertySampleType = 4
)

type CMPropertySource int

const (
	// Lib properties are provided by CM library, and only CM library
	CMPropertySourceLib CMPropertySource = iota
	// Client properties are provided by the client, and only the client
	CMPropertySourceClient
)

type CMPropertyConfig struct {
	Type       reflect.Kind
	Source     CMPropertySource
	Optional   bool
	SampleType CMPropertySampleType
}

func requiredPropertyConfig(t reflect.Kind, sampleType CMPropertySampleType) *CMPropertyConfig {
	return &CMPropertyConfig{
		Type:       t,
		Source:     CMPropertySourceLib,
		Optional:   false,
		SampleType: sampleType,
	}
}
func optionalPropertyConfig(t reflect.Kind, sampleType CMPropertySampleType) *CMPropertyConfig {
	return &CMPropertyConfig{
		Type:       t,
		Source:     CMPropertySourceLib,
		Optional:   true,
		SampleType: sampleType,
	}
}
func wellKnownPropertyConfig(t reflect.Kind, sampleType CMPropertySampleType) *CMPropertyConfig {
	return &CMPropertyConfig{
		Type:       t,
		Source:     CMPropertySourceClient,
		Optional:   true,
		SampleType: sampleType,
	}
}

// TODO: audit the sample types
func BuiltInPropertyTypes() map[string]*CMPropertyConfig {
	return map[string]*CMPropertyConfig{
		"platform":                requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"os_version":              requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"device_manufacturer":     requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"device_model":            requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"device_model_class":      requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"locale_language_code":    requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"locale_country_code":     requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"locale_currency_code":    requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"app_version":             requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"user_interface_idiom":    requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"app_id":                  requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"screen_width_pixels":     requiredPropertyConfig(reflect.Int, CMPropertySampleTypeAppStart),
		"screen_height_pixels":    requiredPropertyConfig(reflect.Int, CMPropertySampleTypeAppStart),
		"screen_width_points":     requiredPropertyConfig(reflect.Int, CMPropertySampleTypeAppStart),
		"screen_height_points":    requiredPropertyConfig(reflect.Int, CMPropertySampleTypeAppStart),
		"screen_scale":            requiredPropertyConfig(reflect.Float64, CMPropertySampleTypeAppStart),
		"device_battery_state":    requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"device_battery_level":    requiredPropertyConfig(reflect.Float64, CMPropertySampleTypeAppStart),
		"device_low_power_mode":   requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeAppStart),
		"device_orientation":      requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"interface_orientation":   requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"dark_mode":               requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeAppStart),
		"network_connection_type": requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"has_wifi_connection":     requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeAppStart),
		"has_cell_connection":     requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeAppStart),
		"has_active_network":      requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeAppStart),
		"expensive_network":       requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeAppStart),
		"cm_version":              requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"foreground":              requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeAppStart),
		"app_install_date":        requiredPropertyConfig(CMTimeKind, CMPropertySampleTypeAppStart),
		"timezone_gmt_offset":     requiredPropertyConfig(reflect.Int, CMPropertySampleTypeAppStart),
		"app_state":               requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"has_watch":               requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeAppStart),
		"screen_brightness":       requiredPropertyConfig(reflect.Float64, CMPropertySampleTypeAppStart),
		"screen_captured":         requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeAppStart),

		// Audio
		"other_audio_playing": requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeAppStart),
		"has_headphones":      requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeAppStart),
		"has_bt_headphones":   requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeAppStart),
		"has_bt_headset":      requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeAppStart),
		"has_wired_headset":   requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeAppStart),
		"has_car_audio":       requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeAppStart),
		"on_call":             requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeAppStart),

		// Location
		"location_permission":          requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeAppStart),
		"location_permission_detailed": requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"location_latitude":            requiredPropertyConfig(reflect.Float64, CMPropertySampleTypeAppStart),
		"location_longitude":           requiredPropertyConfig(reflect.Float64, CMPropertySampleTypeAppStart),
		"location_city":                requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"location_region":              requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"location_country":             requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"location_approx_city":         requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"location_approx_region":       requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"location_approx_country":      requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"location_approx_latitude":     requiredPropertyConfig(reflect.Float64, CMPropertySampleTypeAppStart),
		"location_approx_longitude":    requiredPropertyConfig(reflect.Float64, CMPropertySampleTypeAppStart),

		// Permissions
		"notifications_permission": requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"microphone_permission":    requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"camera_permission":        requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"contacts_permission":      requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"photo_library_permission": requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"add_photo_permission":     requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"calendar_permission":      requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"reminders_permission":     requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"bluetooth_permission":     requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),

		// Optional built in props
		"device_model_version": optionalPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"low_data_mode":        optionalPropertyConfig(reflect.Bool, CMPropertySampleTypeAppStart),

		// Well known properties - client should provide
		"user_signup_date": wellKnownPropertyConfig(CMTimeKind, CMPropertySampleTypeOnCustomSet),
	}
}
