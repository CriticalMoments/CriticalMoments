package datamodel

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model/conditions"
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/checker"
	"github.com/antonmedv/expr/conf"
	"github.com/antonmedv/expr/optimizer"
	"github.com/antonmedv/expr/parser"
	"github.com/antonmedv/expr/vm"
	"golang.org/x/exp/maps"
)

type ConditionDynamicFunction struct {
	Function func(params ...any) (any, error)
	Types    []any
}

type Condition struct {
	conditionString string
}

func NewCondition(s string) (*Condition, error) {
	c := Condition{
		conditionString: s,
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	return &c, nil
}

// Stringer Interface
func (c *Condition) String() string {
	return c.conditionString
}

const CMTimeKind = (reflect.Kind)(math.MaxUint)

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

		// Permisions
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

func StaticConditionHelperFunctions() map[string]interface{} {
	return map[string]interface{}{
		"versionNumberComponent": conditions.VersionNumberComponent,
		"versionGreaterThan":     conditions.VersionGreaterThan,
		"versionLessThan":        conditions.VersionLessThan,
		"versionEqual":           conditions.VersionEqual,

		"unixTimeNanoseconds":  conditions.UnixTimeNanoseconds,
		"unixTimeMilliseconds": conditions.UnixTimeMilliseconds,
		"unixTimeSeconds":      conditions.UnixTimeSeconds,
		"formatTime":           conditions.TimeFormat,

		"rand":        conditions.Random,
		"sessionRand": conditions.SessionRandom,
		"randForKey":  conditions.RandomForKey,
	}
}

func StaticConditionConstantProperties() map[string]interface{} {
	return map[string]interface{}{
		"RFC3339":              time.RFC3339Nano,
		"RFC822":               time.RFC822,
		"RFC850":               time.RFC850,
		"RFC1123":              time.RFC1123,
		"RFC822Z":              time.RFC822Z,
		"RFC1123Z":             time.RFC1123Z,
		"date_with_tz_format":  conditions.DateWithTzFormat,
		"date_and_time_format": conditions.DateAndTimeFormat,
		"date_format":          conditions.DateFormat,
	}
}

var AllBuiltInDynamicFunctions = map[string]bool{
	"eventCount":          true,
	"eventCountWithLimit": true,
	"canOpenUrl":          true,
}

type ConditionFields struct {
	Identifiers []string
	Variables   []string
	Methods     []string
}

// An AST walker we use to analyze code, to see if it's compatible with CM
type conditionWalker struct {
	condition   string
	identifiers map[string]bool
	variables   map[string]bool
	methods     map[string]bool
}

func (v *conditionWalker) Visit(n *ast.Node) {
	if node, ok := (*n).(*ast.IdentifierNode); ok {
		v.identifiers[node.Value] = true

		// Check if this is a variable or a method. Unfortunately .Method() on the node does not work
		// so we check for open paren immediately after the identifier
		isMethod := false
		parenLoc := (*n).Location().Column + len(node.Value)
		if parenLoc < len(v.condition) {
			paran := v.condition[parenLoc : parenLoc+1]
			if paran == "(" {
				isMethod = true
			}
		}
		if isMethod {
			v.methods[node.Value] = true
		} else {
			v.variables[node.Value] = true
		}
	}
}

func (c *Condition) ExtractIdentifiers() (returnFields *ConditionFields, returnError error) {
	// expr can panic, so catch it and return an error instead
	defer func() {
		if r := recover(); r != nil {
			returnFields = nil
			returnError = fmt.Errorf("panic in ExtractIdentifiers: %v", r)
		}
	}()

	// single line needed because we use the location offset
	singleLineCondition := strings.ReplaceAll(c.conditionString, "\n", " ")
	tree, err := parser.Parse(singleLineCondition)
	if err != nil {
		return nil, err
	}

	config := conf.New(conf.CreateNew())
	config.Strict = false
	_, err = checker.Check(tree, config)
	if err != nil {
		return nil, err
	}
	err = optimizer.Optimize(&tree.Node, config)
	if err != nil {
		return nil, err
	}

	visitor := &conditionWalker{
		condition:   singleLineCondition,
		identifiers: map[string]bool{},
		variables:   map[string]bool{},
		methods:     map[string]bool{},
	}
	ast.Walk(&tree.Node, visitor)

	results := ConditionFields{
		Identifiers: maps.Keys(visitor.identifiers),
		Variables:   maps.Keys(visitor.variables),
		Methods:     maps.Keys(visitor.methods),
	}
	return &results, nil
}

func (c *Condition) Validate() error {
	if c.conditionString == "" {
		return NewUserPresentableError("Condition is empty string (not allowed). Use 'true' or 'false' for minimal condition.")
	}

	// Run this even if not strict. It is checking the format of the condition as well
	fields, err := c.ExtractIdentifiers()
	if err != nil {
		return err
	}

	if StrictDatamodelParsing {
		// Don't check variable names. We support custom vars so every name is valid

		// Check we support all methods used if strict parsing
		for _, methodName := range fields.Methods {
			if _, ok := AllBuiltInDynamicFunctions[methodName]; !ok {
				if _, ok := StaticConditionHelperFunctions()[methodName]; !ok {
					return NewUserPresentableError(fmt.Sprintf("Method included in condition which isn't recognized: %v", methodName))
				}
			}
		}

	}

	return nil
}

func (c *Condition) CompileWithEnv(ops ...expr.Option) (resultProgram *vm.Program, returnError error) {
	// expr can panic, so catch it and return an error instead
	defer func() {
		if r := recover(); r != nil {
			resultProgram = nil
			returnError = fmt.Errorf("panic in CompileWithEnv: %v", r)
		}
	}()

	allOptions := append(ops, expr.AsBool())

	return expr.Compile(c.conditionString, allOptions...)
}

func (c *Condition) UnmarshalJSON(data []byte) error {
	var conditionString *string
	err := json.Unmarshal(data, &conditionString)
	if err != nil {
		return NewUserPresentableErrorWSource(fmt.Sprintf("Invalid Condition String [[ %s ]]", string(data)), err)
	}
	c.conditionString = *conditionString

	if err := c.Validate(); err != nil {
		// Fallback to returning empty on non-strict clients. Don't want entire config file to fail
		// Downstream during eval we return false and error
		c.conditionString = ""
		if StrictDatamodelParsing {
			return NewUserPresentableErrorWSource(fmt.Sprintf("Invalid Condition: [[ %v ]]", string(data)), err)
		}
	}

	return nil
}
