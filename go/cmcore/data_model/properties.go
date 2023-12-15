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

func BuiltInPropertyTypes() map[string]*CMPropertyConfig {
	return map[string]*CMPropertyConfig{
		"platform":                  requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"os_version":                requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"device_manufacturer":       requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"device_model":              requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"device_model_class":        requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"locale_language_code":      requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"locale_country_code":       requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"locale_currency_code":      requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"locale_language_direction": requiredPropertyConfig(reflect.String, CMPropertySampleTypeDoNotSample),
		"app_version":               requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"user_interface_idiom":      requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"app_id":                    requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"screen_width_pixels":       requiredPropertyConfig(reflect.Int, CMPropertySampleTypeAppStart),
		"screen_height_pixels":      requiredPropertyConfig(reflect.Int, CMPropertySampleTypeAppStart),
		"screen_width_points":       requiredPropertyConfig(reflect.Int, CMPropertySampleTypeAppStart),
		"screen_height_points":      requiredPropertyConfig(reflect.Int, CMPropertySampleTypeAppStart),
		"screen_scale":              requiredPropertyConfig(reflect.Float64, CMPropertySampleTypeAppStart),
		"device_battery_state":      requiredPropertyConfig(reflect.String, CMPropertySampleTypeOnUse),
		"device_battery_level":      requiredPropertyConfig(reflect.Float64, CMPropertySampleTypeOnUse),
		"device_low_power_mode":     requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeAppStart),
		"device_orientation":        requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"interface_orientation":     requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"dark_mode":                 requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeAppStart),
		"network_connection_type":   requiredPropertyConfig(reflect.String, CMPropertySampleTypeOnUse),
		"has_wifi_connection":       requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeOnUse),
		"has_cell_connection":       requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeOnUse),
		"has_active_network":        requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeOnUse),
		"expensive_network":         requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeOnUse),
		"cm_version":                requiredPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"foreground":                requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeDoNotSample),
		"app_state":                 requiredPropertyConfig(reflect.String, CMPropertySampleTypeDoNotSample),
		"app_install_date":          requiredPropertyConfig(CMTimeKind, CMPropertySampleTypeDoNotSample),
		"timezone_gmt_offset":       requiredPropertyConfig(reflect.Int, CMPropertySampleTypeAppStart),
		"has_watch":                 requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeOnUse),
		"screen_brightness":         requiredPropertyConfig(reflect.Float64, CMPropertySampleTypeAppStart),
		"screen_captured":           requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeAppStart),

		// Audio
		"other_audio_playing": requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeAppStart),
		"has_headphones":      requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeOnUse),
		"has_bt_headphones":   requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeOnUse),
		"has_bt_headset":      requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeOnUse),
		"has_wired_headset":   requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeOnUse),
		"has_car_audio":       requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeOnUse),
		"on_call":             requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeOnUse),

		// Location
		"location_permission":          requiredPropertyConfig(reflect.Bool, CMPropertySampleTypeOnUse),
		"location_permission_detailed": requiredPropertyConfig(reflect.String, CMPropertySampleTypeOnUse),
		"location_latitude":            requiredPropertyConfig(reflect.Float64, CMPropertySampleTypeOnUse),
		"location_longitude":           requiredPropertyConfig(reflect.Float64, CMPropertySampleTypeOnUse),
		"location_city":                requiredPropertyConfig(reflect.String, CMPropertySampleTypeOnUse),
		"location_region":              requiredPropertyConfig(reflect.String, CMPropertySampleTypeOnUse),
		"location_country":             requiredPropertyConfig(reflect.String, CMPropertySampleTypeOnUse),
		"location_approx_city":         requiredPropertyConfig(reflect.String, CMPropertySampleTypeOnUse),
		"location_approx_region":       requiredPropertyConfig(reflect.String, CMPropertySampleTypeOnUse),
		"location_approx_country":      requiredPropertyConfig(reflect.String, CMPropertySampleTypeOnUse),
		"location_approx_latitude":     requiredPropertyConfig(reflect.Float64, CMPropertySampleTypeOnUse),
		"location_approx_longitude":    requiredPropertyConfig(reflect.Float64, CMPropertySampleTypeOnUse),

		// Permissions
		"notifications_permission": requiredPropertyConfig(reflect.String, CMPropertySampleTypeOnUse),
		"microphone_permission":    requiredPropertyConfig(reflect.String, CMPropertySampleTypeOnUse),
		"camera_permission":        requiredPropertyConfig(reflect.String, CMPropertySampleTypeOnUse),
		"contacts_permission":      requiredPropertyConfig(reflect.String, CMPropertySampleTypeOnUse),
		"photo_library_permission": requiredPropertyConfig(reflect.String, CMPropertySampleTypeOnUse),
		"add_photo_permission":     requiredPropertyConfig(reflect.String, CMPropertySampleTypeOnUse),
		"calendar_permission":      requiredPropertyConfig(reflect.String, CMPropertySampleTypeOnUse),
		"reminders_permission":     requiredPropertyConfig(reflect.String, CMPropertySampleTypeOnUse),
		"bluetooth_permission":     requiredPropertyConfig(reflect.String, CMPropertySampleTypeOnUse),

		// Optional built in props
		"device_model_version": optionalPropertyConfig(reflect.String, CMPropertySampleTypeAppStart),
		"low_data_mode":        optionalPropertyConfig(reflect.Bool, CMPropertySampleTypeOnUse),

		// Well known properties - client should provide
		"user_signup_date":      wellKnownPropertyConfig(CMTimeKind, CMPropertySampleTypeOnCustomSet),
		"user_signed_in":        wellKnownPropertyConfig(reflect.Bool, CMPropertySampleTypeOnCustomSet),
		"have_user_email":       wellKnownPropertyConfig(reflect.Bool, CMPropertySampleTypeOnCustomSet),
		"user_email_validated":  wellKnownPropertyConfig(reflect.Bool, CMPropertySampleTypeOnCustomSet),
		"have_user_phone":       wellKnownPropertyConfig(reflect.Bool, CMPropertySampleTypeOnCustomSet),
		"user_age":              wellKnownPropertyConfig(reflect.Int, CMPropertySampleTypeOnCustomSet),
		"user_approx_age":       wellKnownPropertyConfig(reflect.Int, CMPropertySampleTypeOnCustomSet),
		"user_pronouns":         wellKnownPropertyConfig(reflect.String, CMPropertySampleTypeOnCustomSet),
		"user_gender":           wellKnownPropertyConfig(reflect.String, CMPropertySampleTypeOnCustomSet),
		"user_inferred_gender":  wellKnownPropertyConfig(reflect.String, CMPropertySampleTypeOnCustomSet),
		"has_paid_subscription": wellKnownPropertyConfig(reflect.Bool, CMPropertySampleTypeOnCustomSet),
		"ever_subscribed":       wellKnownPropertyConfig(reflect.Bool, CMPropertySampleTypeOnCustomSet),
		"has_purchased":         wellKnownPropertyConfig(reflect.Bool, CMPropertySampleTypeOnCustomSet),
		"purchase_count":        wellKnownPropertyConfig(reflect.Int, CMPropertySampleTypeOnCustomSet),
		"total_purchase_value":  wellKnownPropertyConfig(reflect.Float64, CMPropertySampleTypeOnCustomSet),
		"referral_source":       wellKnownPropertyConfig(reflect.String, CMPropertySampleTypeOnCustomSet),
		"referral_id":           wellKnownPropertyConfig(reflect.String, CMPropertySampleTypeOnCustomSet),
		"user_was_referred":     wellKnownPropertyConfig(reflect.Bool, CMPropertySampleTypeOnCustomSet),
		"user_referral_count":   wellKnownPropertyConfig(reflect.Int, CMPropertySampleTypeOnCustomSet),
		"session_source":        wellKnownPropertyConfig(reflect.String, CMPropertySampleTypeOnCustomSet),
	}
}
