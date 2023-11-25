package appcore

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/CriticalMoments/CriticalMoments/go/appcore/db"
	datamodel "github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model"
)

func propertyValueOrNil(pr *propertyRegistry, key string) interface{} {
	v, err := pr.propertyValue(key)
	if err != nil {
		return nil
	}
	return v
}

func TestPropertyRegistrySetGet(t *testing.T) {
	pr := newPropertyRegistry()
	pr.builtInPropertyTypes = map[string]*datamodel.CMPropertyConfig{
		"a": {Type: reflect.String, Source: datamodel.CMPropertySourceLib, Optional: true},
		"b": {Type: reflect.Int, Source: datamodel.CMPropertySourceLib, Optional: true},
		"c": {Type: reflect.Float64, Source: datamodel.CMPropertySourceLib, Optional: true},
		"d": {Type: reflect.Bool, Source: datamodel.CMPropertySourceLib, Optional: true},
		"e": {Type: datamodel.CMTimeKind, Source: datamodel.CMPropertySourceLib, Optional: true},
	}

	err := pr.registerStaticProperty("a", "a")
	if err != nil {
		t.Fatal(err)
	}
	if propertyValueOrNil(pr, "a") != "a" {
		t.Fatal("Property registry failed for string")
	}
	err = pr.registerStaticProperty("b", 2)
	if err != nil {
		t.Fatal(err)
	}
	if propertyValueOrNil(pr, "b") != 2 {
		t.Fatal("Property registry failed for int")
	}
	err = pr.registerStaticProperty("c", 3.3)
	if err != nil {
		t.Fatal(err)
	}
	if propertyValueOrNil(pr, "c") != 3.3 {
		t.Fatal("Property registry failed for int")
	}
	err = pr.registerStaticProperty("d", true)
	if err != nil {
		t.Fatal(err)
	}
	if propertyValueOrNil(pr, "d") != true {
		t.Fatal("Property registry failed for bool")
	}
	now := time.Now()
	err = pr.registerStaticProperty("e", now)
	if err != nil {
		t.Fatal(err)
	}
	if propertyValueOrNil(pr, "e").(time.Time) != now {
		t.Fatal("Property registry failed for time")
	}
}

func TestPropertyRegistrySetInvalid(t *testing.T) {
	pr := newPropertyRegistry()
	pr.builtInPropertyTypes = map[string]*datamodel.CMPropertyConfig{
		"a": {Type: reflect.String, Source: datamodel.CMPropertySourceLib, Optional: true},
		"b": {Type: reflect.Int, Source: datamodel.CMPropertySourceLib, Optional: true},
	}
	err := pr.registerStaticProperty("a", 1) // type mismatch from expected
	if err == nil {
		t.Fatal("Allowed type mismatch")
	}
	if propertyValueOrNil(pr, "a") != nil {
		t.Fatal("Property registry allowed invalid")
	}
	err = pr.registerStaticProperty("a", time.Now()) // type mismatch from expected
	if err == nil {
		t.Fatal("Allowed type mismatch")
	}
	if propertyValueOrNil(pr, "a") != nil {
		t.Fatal("Property registry allowed invalid")
	}
	err = pr.registerStaticProperty("a", "aval") // correct type
	if err != nil {
		t.Fatal(err)
	}
	if propertyValueOrNil(pr, "a") != "aval" {
		t.Fatal("Failed to set with valid type")
	}

	err = pr.registerStaticProperty("b", []string{}) // invalid type
	if err == nil {
		t.Fatal("Allowed invalid type")
	}
	if propertyValueOrNil(pr, "b") != nil {
		t.Fatal("Property registry allowed invalid")
	}
	err = pr.registerStaticProperty("b", 42) // correct type
	if err != nil {
		t.Fatal(err)
	}
	if propertyValueOrNil(pr, "b") != 42 {
		t.Fatal("Failed to set with valid type")
	}

	err = pr.registerStaticProperty("c", 3.3) // unexpected key
	if err == nil {
		t.Fatal("Allowed unexpected key")
	}
	if propertyValueOrNil(pr, "c") != nil {
		t.Fatal("Property registry allowed invalid")
	}
}

func TestPropertyRegistryValidateRequired(t *testing.T) {
	pr := newPropertyRegistry()
	pr.builtInPropertyTypes = map[string]*datamodel.CMPropertyConfig{
		"platform": {Type: reflect.String, Source: datamodel.CMPropertySourceLib, Optional: false},
	}

	if pr.validateProperties() == nil {
		t.Fatal("Validated missing required properties")
	}
	pr.registerStaticProperty("platform", 42)
	if pr.validateProperties() == nil {
		t.Fatal("Validated with type mismatch")
	}
	pr.registerStaticProperty("platform", "ios")
	if pr.validateProperties() != nil {
		t.Fatal("Validation failed on valid type")
	}
}

func TestPropertyRegistryValidateOptional(t *testing.T) {
	pr := newPropertyRegistry()
	pr.builtInPropertyTypes = map[string]*datamodel.CMPropertyConfig{
		"optional_bool": {Type: reflect.Bool, Source: datamodel.CMPropertySourceLib, Optional: true},
	}

	if pr.validateProperties() != nil {
		t.Fatal("Missing optional failed validation")
	}
	err := pr.registerStaticProperty("optional_bool", 42)
	if err == nil {
		t.Fatal("Added with type mismatch")
	}
	err = pr.registerStaticProperty("optional_bool", true)
	if err != nil {
		t.Fatal(err)
	}
	if pr.validateProperties() != nil {
		t.Fatal("Validation failed on valid type")
	}
}

func testHelperNewCondition(s string, t *testing.T) *datamodel.Condition {
	c, err := datamodel.NewCondition(s)
	if err != nil {
		t.Fatalf("Condition in test is not valid %v", s)
	}
	return c
}

func TestPropertyRegistryVersionNumberHelpers(t *testing.T) {
	pr := newPropertyRegistry()
	pr.builtInPropertyTypes = map[string]*datamodel.CMPropertyConfig{}

	versionConditions := testHelperNewCondition(`
		(versionGreaterThan('invalid', '1.0') == false) && 
		(versionGreaterThan('1.1', '1.0') == true) && 
		(versionGreaterThan('1.0', '1.0') == false) && 
		(versionGreaterThan('1.0', '2.0') == false) && 
		(versionLessThan('invalid', '1.0') == false) && 
		(versionLessThan('1.1', '1.0') == false) && 
		(versionLessThan('1.0', '1.0') == false) && 
		(versionLessThan('1.0', '2.0') == true) && 
		(versionEqual('invalid', '1') == false) && 
		(versionEqual('v1.2.3', '1.2.3') == true) && 
		(versionEqual('v2', 'v1') == false) 
	`, t)
	if r, err := pr.evaluateCondition(versionConditions); err != nil || !r {
		t.Fatalf("Version helpers failed: %v", err)
	}
}

func TestPropertyRegistryVersionNumberComponent(t *testing.T) {
	pr := newPropertyRegistry()
	pr.builtInPropertyTypes = map[string]*datamodel.CMPropertyConfig{
		"os_version":  {Type: reflect.String, Source: datamodel.CMPropertySourceLib, Optional: true},
		"app_version": {Type: reflect.String, Source: datamodel.CMPropertySourceLib, Optional: true},
	}

	// Valid string -- saves and able to parse components with function
	if err := pr.registerStaticProperty("os_version", "1.2.3"); err != nil {
		t.Fatal("Valid version number failed to save")
	}
	if propertyValueOrNil(pr, "os_version") != "1.2.3" {
		t.Fatal("Valid version number failed to save")
	}
	if r, err := pr.evaluateCondition(testHelperNewCondition("versionNumberComponent(os_version, 0) == 1", t)); err != nil || !r {
		t.Fatal("Valid version number failed to extract component")
	}
	if r, err := pr.evaluateCondition(testHelperNewCondition("versionNumberComponent(os_version, 0) == 2", t)); err != nil || r {
		t.Fatal("Valid version number failed to extract component that fails test")
	}
	if r, err := pr.evaluateCondition(testHelperNewCondition("versionNumberComponent(os_version, 1) == 2", t)); err != nil || !r {
		t.Fatal("Valid version number failed to extract component")
	}
	if r, err := pr.evaluateCondition(testHelperNewCondition("versionNumberComponent(os_version, 2) == 3", t)); err != nil || !r {
		t.Fatal("Valid version number failed to extract component")
	}
	if r, err := pr.evaluateCondition(testHelperNewCondition("versionNumberComponent(os_version, 3) == nil", t)); err != nil || !r {
		t.Fatal("Valid version number failed to extract component")
	}

	// Invalid version string
	if err := pr.registerStaticProperty("app_version", "1.b.3"); err != nil {
		t.Fatal("Invalid version number failed to save. Should still save as string for exact comparison")
	}
	if propertyValueOrNil(pr, "app_version") != "1.b.3" {
		t.Fatal("Invalid version number failed to save. Should still save as string for exact comparison")
	}
	if r, err := pr.evaluateCondition(testHelperNewCondition("versionNumberComponent(app_version, 0) == nil", t)); err != nil || !r {
		t.Fatal("Invalid version failed to return nil for component")
	}
}

type testPropertyProvider struct {
	val int64
}

func (p *testPropertyProvider) Type() int {
	return LibPropertyProviderTypeInt
}
func (p *testPropertyProvider) IntValue() int64 {
	p.val = p.val + 1
	return p.val
}
func (p *testPropertyProvider) StringValue() string {
	return ""
}
func (p *testPropertyProvider) FloatValue() float64 {
	return 0.0
}
func (p *testPropertyProvider) BoolValue() bool {
	return false
}

const testTimestampUnixMilli = 1698767307000

func (p *testPropertyProvider) TimeEpochMilliseconds() int64 {
	return testTimestampUnixMilli
}

func TestDynamicProperties(t *testing.T) {
	pr := newPropertyRegistry()
	pr.builtInPropertyTypes = map[string]*datamodel.CMPropertyConfig{
		"screen_width_pixels": {Type: reflect.Int, Source: datamodel.CMPropertySourceLib, Optional: false},
		"a":                   {Type: reflect.Int, Source: datamodel.CMPropertySourceClient, Optional: true},
	}

	dp := testPropertyProvider{}
	err := pr.registerLibPropertyProvider("a", &dp)
	if err == nil {
		t.Fatal("allowed registering a 'Lib' property provider that isn't built in")
	}
	err = pr.registerLibPropertyProvider("b", &dp)
	if err == nil {
		t.Fatal("allowed registering a 'Lib' property provider that isn't built in")
	}

	err = pr.registerLibPropertyProvider("screen_width_pixels", &dp)
	if err != nil {
		t.Fatal(err)
	}
	if propertyValueOrNil(pr, "screen_width_pixels").(int64) != 1 {
		t.Fatal("dynamic property doesn't work")
	}
	if propertyValueOrNil(pr, "screen_width_pixels").(int64) != 2 {
		t.Fatal("dynamic property not dynamic")
	}
	if propertyValueOrNil(pr, "screen_width_pixels").(int64) != 3 {
		t.Fatal("dynamic property not dynamic")
	}
}

func TestPropertyRegistryConditionEval(t *testing.T) {
	pr := newPropertyRegistry()
	pr.builtInPropertyTypes = map[string]*datamodel.CMPropertyConfig{
		"app_version":         {Type: reflect.String, Source: datamodel.CMPropertySourceLib, Optional: true},       // using as a populated string
		"screen_width_pixels": {Type: reflect.Int, Source: datamodel.CMPropertySourceLib, Optional: true},          // using as a populated int
		"os_version":          {Type: reflect.String, Source: datamodel.CMPropertySourceLib, Optional: true},       // using as a nil string
		"test_time":           {Type: datamodel.CMTimeKind, Source: datamodel.CMPropertySourceLib, Optional: true}, // populated date
	}

	pr.registerStaticProperty("app_version", "hello")
	pr.registerStaticProperty("screen_width_pixels", 42)
	pr.registerStaticProperty("test_time", time.UnixMilli(testTimestampUnixMilli))
	if propertyValueOrNil(pr, "app_version") != "hello" {
		t.Fatal("property not set")
	}
	if propertyValueOrNil(pr, "screen_width_pixels") != 42 {
		t.Fatal("property not set")
	}

	result, err := pr.evaluateCondition(testHelperNewCondition("a > 2", t))
	if err == nil || result {
		t.Fatal("Allowed invalid conditions: nil > 2")
	}

	result, err = pr.evaluateCondition(testHelperNewCondition("3 > 2", t))
	if err != nil || !result {
		t.Fatal("Failed to eval simple true condition")
	}

	result, err = pr.evaluateCondition(testHelperNewCondition("1 > 2", t))
	if err != nil || result {
		t.Fatal("Failed to eval simple false condition")
	}

	result, err = pr.evaluateCondition(testHelperNewCondition("app_version == 'hello'", t))
	if err != nil || !result {
		t.Fatal("Failed to eval true condition")
	}

	result, err = pr.evaluateCondition(testHelperNewCondition("app_version startsWith 'hel'", t))
	if err != nil || !result {
		t.Fatal("Failed to eval true condition with builtin function")
	}

	result, err = pr.evaluateCondition(testHelperNewCondition("app_version == 'world'", t))
	if err != nil || result {
		t.Fatal("Failed to eval false condition with property")
	}

	result, err = pr.evaluateCondition(testHelperNewCondition("1 + 2 + 3", t))
	if err == nil || result {
		t.Fatal("Allowed condition with non bool result")
	}

	timeCondition := fmt.Sprintf("test_time == unixTimeMilliseconds(%d) && test_time > unixTimeMilliseconds(%d)", testTimestampUnixMilli, testTimestampUnixMilli-1)
	result, err = pr.evaluateCondition(testHelperNewCondition(timeCondition, t))
	if err != nil || !result {
		t.Fatal("Failed conduition with timestamps")
	}

	// Need reflection to created these invalid cases.
	// Empty strings may exist in non-strict mode parsing so really important case that it returns false and err
	con := &datamodel.Condition{}
	v := reflect.ValueOf(con).Elem()
	cf := v.FieldByName("conditionString")
	cf = reflect.NewAt(cf.Type(), unsafe.Pointer(cf.UnsafeAddr())).Elem()
	cf.SetString("")
	result, err = pr.evaluateCondition(con)
	if err == nil || result {
		t.Fatal("Allowed empty condition")
	}
	cf.SetString("app_version ^#$%")
	result, err = pr.evaluateCondition(con)
	if err == nil || result {
		t.Fatal("Allowed invalid condition")
	}

	result, err = pr.evaluateCondition(testHelperNewCondition("screen_width_pixels > 99", t))
	if err != nil || result {
		t.Fatal("false condition passed")
	}

	result, err = pr.evaluateCondition(testHelperNewCondition("screen_width_pixels > 2", t))
	if err != nil || !result {
		t.Fatal("true condition failed")
	}

	result, err = pr.evaluateCondition(testHelperNewCondition("os_version == nil", t))
	if err != nil || !result {
		t.Fatal("true condition for allowed missing var failed")
	}

	// Failing, used to be nil. Not not nil
	result, err = pr.evaluateCondition(testHelperNewCondition("os_version ?? false", t))
	if err != nil || result {
		t.Fatal("nil condition did not eval to false")
	}

	result, err = pr.evaluateCondition(testHelperNewCondition("os_version == false", t))
	if err != nil || result {
		t.Fatal("missing condition should be differentiateable from false bool")
	}
}

func TestDateFunctionsInConditions(t *testing.T) {
	pr := newPropertyRegistry()
	pr.builtInPropertyTypes = map[string]*datamodel.CMPropertyConfig{}

	// time library created
	result, err := pr.evaluateCondition(testHelperNewCondition("now() < unixTimeMilliseconds(1688764455000)", t))
	if err != nil || result {
		t.Fatal("back to the future isn't real: now() gave incorrect time")
	}
	// time in 2050
	result, err = pr.evaluateCondition(testHelperNewCondition("now() > unixTimeMilliseconds(2540841255000)", t))
	if err != nil || result {
		t.Fatal("back to the future 2 isn't real: now() gave incorrect time")
	}

	verifyDurationCondition := testHelperNewCondition(`
		(duration('2s') == 2 * duration('1s')) &&
		(duration('1m') == 60 * duration('1s')) &&
		(duration('1h') == 60 * duration('1m'))
	`, t)
	result, err = pr.evaluateCondition(verifyDurationCondition)
	if err != nil || !result {
		t.Fatal("Duration functions failed e2e test")
	}

	// Verify parseDate function is wired up properly.
	result, err = pr.evaluateCondition(testHelperNewCondition("date('2006-01-02T15:04:05.999+07:00', RFC3339) == unixTimeMilliseconds(1136189045999)", t))
	if err != nil || !result {
		t.Fatal("date() did not work inside a condition")
	}

	// Verify RFC3339Nano also works without the subsecond component
	result, err = pr.evaluateCondition(testHelperNewCondition("date('2006-01-02T15:04:05+07:00', RFC3339) == unixTimeSeconds(1136189045)", t))
	if err != nil || !result {
		t.Fatal("date() did not work inside a condition")
	}

	// Check invalid returns errors though the stack
	result, err = pr.evaluateCondition(testHelperNewCondition("date('invalid') == unixTimeMilliseconds(1136189045999)", t))
	if err == nil || result {
		t.Fatal("invalid date didn't error", err)
	}

	// Check date formatting
	local := time.Local
	defer func() {
		time.Local = local
	}()
	time.Local, _ = time.LoadLocation("America/St_Johns")
	formatConditionCases := []string{
		"formatTime(unixTimeMilliseconds(1698881571000), 'hod', 'America/Toronto') == 19",
		"formatTime(unixTimeMilliseconds(1698881571000), 'hod', 'America/St_Johns') == 21",     // add_test_count
		"formatTime(unixTimeMilliseconds(1698881571000), 'hod', 'Local') == 21",                // add_test_count
		"formatTime(unixTimeMilliseconds(1698881571000), 'hod') == 21",                         // add_test_count
		"formatTime(unixTimeMilliseconds(1698881571000), 'hod', '') == 21",                     // add_test_count
		"formatTime(unixTimeMilliseconds(1698881571000), 'hod', 'Local', 'extraParam') == nil", // add_test_count
	}
	for _, c := range formatConditionCases {
		result, err = pr.evaluateCondition(testHelperNewCondition(c, t))
		if err != nil || !result {
			t.Fatal("Failed to use time format function: ", c)
		}
	}
}

func TestUnknownVarsInConditions(t *testing.T) {
	c, err := datamodel.NewCondition("unknown_var == nil")
	if err != nil {
		t.Fatal(err)
	}

	pr := newPropertyRegistry()
	pr.builtInPropertyTypes = map[string]*datamodel.CMPropertyConfig{}

	result, err := pr.evaluateCondition(c)
	if err != nil {
		t.Fatal(err)
	}
	if result != true {
		t.Fatal("failed to return nil for unknown var")
	}

	c, err = datamodel.NewCondition("unknown_var > 6")
	if err != nil {
		t.Fatal(err)
	}
	result, err = pr.evaluateCondition(c)
	if err == nil {
		t.Fatal(err)
	}
	if result != false {
		t.Fatal("failed to return false for error")
	}

	c, err = datamodel.NewCondition("unknownFunction() == nil")
	if err != nil {
		t.Fatal(err)
	}
	result, err = pr.evaluateCondition(c)
	if err != nil || !result {
		t.Fatal("unknown function didn't return nil")
	}
}

func TestDynamicMethods(t *testing.T) {
	pr := newPropertyRegistry()
	pr.builtInPropertyTypes = map[string]*datamodel.CMPropertyConfig{
		"platform": {Type: reflect.String, Source: datamodel.CMPropertySourceLib, Optional: false},
	}
	pr.registerStaticProperty("platform", "ios")
	pr.RegisterDynamicFunctions(map[string]*datamodel.ConditionDynamicFunction{
		"testFunc": {
			Function: func(params ...any) (any, error) {
				if len(params) != 1 {
					return nil, errors.New("eventCount requires one parameter")
				}
				_, ok := params[0].(string)
				if !ok {
					return nil, errors.New("eventCount requires a string parameter")
				}
				return 99, nil
			},
			Types: []any{new(func(string) int)},
		},
		"isNintyNine": {
			Function: func(params ...any) (any, error) {
				v, ok := params[0].(int)
				if !ok {
					return nil, errors.New("notZero requires an int parameter")
				}
				return v != 0, nil
			},
			Types: []any{new(func(int) bool)},
		},
	})

	// Check we a method that doesn't exist returns nil
	result, err := pr.evaluateCondition(testHelperNewCondition("notRealFunc('test') == nil", t))
	if err != nil || !result {
		t.Fatal("Invalid method didn't return nil")
	}

	// Check mixing methods that exist and don't exist
	result, err = pr.evaluateCondition(testHelperNewCondition("notRealFunc('test') == nil && testFunc('asfd') == 99", t))
	if err != nil || !result {
		t.Fatal("Invalid method didn't return nil")
	}

	// Check we can't call a method with the wrong number of params
	result, err = pr.evaluateCondition(testHelperNewCondition("testFunc() == 99", t))
	if err == nil || result {
		t.Fatal("Allowed invalid method")
	}
	result, err = pr.evaluateCondition(testHelperNewCondition("testFunc('test', 'test') == 99", t))
	if err == nil || result {
		t.Fatal("Allowed invalid method")
	}

	// Check we can't call a method with the wrong param type
	result, err = pr.evaluateCondition(testHelperNewCondition("testFunc(1) == 99", t))
	if err == nil || result {
		t.Fatal("Allowed invalid method")
	}

	// Check we can call a method with the right number of params
	result, err = pr.evaluateCondition(testHelperNewCondition("testFunc('test') == 99", t))
	if err != nil || !result {
		t.Fatal("Failed to call method")
	}

	// Check we can call a method with a property
	result, err = pr.evaluateCondition(testHelperNewCondition("testFunc(platform) == 99", t))
	if err != nil || !result {
		t.Fatal("Failed to call method")
	}

	// Check we a method and vars that don't exist, right up to last charater to check string parsing
	result, err = pr.evaluateCondition(testHelperNewCondition("notRealFunc('test') == nil && notRealFunc2 == nil && nil == notRealEnd", t))
	if err != nil || !result {
		t.Fatal("Invalid method/vars didn't return nil")
	}

	// Check we can chain methods and properties
	result, err = pr.evaluateCondition(testHelperNewCondition("isNintyNine(testFunc(platform))", t))
	if err != nil || !result {
		t.Fatal("Failed to call method chain")
	}
}

func TestRandomInConditions(t *testing.T) {
	pr := newPropertyRegistry()
	pr.builtInPropertyTypes = map[string]*datamodel.CMPropertyConfig{}

	result, err := pr.evaluateCondition(testHelperNewCondition("rand() >= 0 && rand() <= 2^63 && rand() != rand()", t))
	if err != nil || !result {
		t.Fatal("random generation not in range or is stable (not random)")
	}

	result, err = pr.evaluateCondition(testHelperNewCondition("sessionRand() >= 0 && sessionRand() <= 2^63 && sessionRand() == sessionRand()", t))
	if err != nil || !result {
		t.Fatal("session random generation not in range, or changing")
	}

	result, err = pr.evaluateCondition(testHelperNewCondition("randForKey('key1', 1) == 292785326893130985 && randForKey('x', sessionRand()) >= 0", t))
	if err != nil || !result {
		t.Fatal("randForKey not stable, or can't be seeded with sessionRand")
	}
}

func TestCustomProperties(t *testing.T) {
	pr := newPropertyRegistry()
	pr.builtInPropertyTypes = map[string]*datamodel.CMPropertyConfig{
		"well_known_v": {Type: reflect.String, Source: datamodel.CMPropertySourceClient, Optional: true},
	}

	// Test a custom properties with correct prefix
	err := pr.registerStaticPropertyWithSource("custom_stringv", datamodel.CMPropertySourceClient, "hello")
	if err != nil {
		t.Fatal(err)
	}
	err = pr.registerStaticPropertyWithSource("custom_boolv", datamodel.CMPropertySourceClient, false)
	if err != nil {
		t.Fatal(err)
	}
	err = pr.registerStaticPropertyWithSource("custom_intv", datamodel.CMPropertySourceClient, 42)
	if err != nil {
		t.Fatal(err)
	}
	err = pr.registerStaticPropertyWithSource("custom_floatv", datamodel.CMPropertySourceClient, 3.3)
	if err != nil {
		t.Fatal(err)
	}
	// test accessing with full name
	result, err := pr.evaluateCondition(testHelperNewCondition("custom_stringv == 'hello' && custom_boolv == false && custom_intv == 42 && custom_floatv == 3.3 && custom_nilv == nil", t))
	if err != nil || !result {
		t.Fatal("custom properties failed")
	}
	// test accessing with short name
	result, err = pr.evaluateCondition(testHelperNewCondition("stringv == 'hello' && boolv == false && intv == 42 && floatv == 3.3 && nilv == nil", t))
	if err != nil || !result {
		t.Fatal("custom properties failed with short names")
	}

	// test without prefix
	err = pr.registerStaticProperty("no_prefix_v", "hello")
	if err == nil {
		t.Fatal("Allowed custom property without prefix")
	}

	// Test an invalid type property (float32)
	err = pr.registerStaticProperty("custom_float32v", float32(3.3))
	if err == nil {
		t.Fatal("Allowed custom property with invalid type")
	}
}

func TestNoPrefixCollision(t *testing.T) {
	// Built in and well known can't use custom prefix
	for k := range datamodel.BuiltInPropertyTypes() {
		if strings.HasPrefix(k, CustomPropertyPrefix) {
			t.Fatalf("Built in property has custom prefix: %s", k)
		}
	}
}

func TestValidateCustomPrefix(t *testing.T) {
	pr := newPropertyRegistry()
	pr.builtInPropertyTypes = map[string]*datamodel.CMPropertyConfig{}

	err := pr.validateProperties()
	if err != nil {
		t.Fatal(err)
	}

	s := staticPropertyProvider{
		value: 42,
	}

	pr.providers["custom_stringv"] = &s
	err = pr.validateProperties()
	if err != nil {
		t.Fatal(err)
	}

	pr.providers["not_custom_stringv"] = &s
	err = pr.validateProperties()
	if err == nil {
		t.Fatal("allowed non custom property")
	}
}

func TestClientPropertyRegistration(t *testing.T) {
	pr := newPropertyRegistry()
	pr.builtInPropertyTypes = map[string]*datamodel.CMPropertyConfig{
		"built_in":    {Type: reflect.String, Source: datamodel.CMPropertySourceLib, Optional: false},
		"well_known":  {Type: reflect.String, Source: datamodel.CMPropertySourceClient, Optional: true},
		"well_known2": {Type: reflect.String, Source: datamodel.CMPropertySourceClient, Optional: true},
	}

	// Should not allow registering well known with wrong type
	err := pr.registerClientProperty("well_known", 42)
	if err == nil || pr.providers["well_known"] != nil {
		t.Fatal("Allowed registering well known with wrong type")
	}

	// Should not allow registering built in
	err = pr.registerClientProperty("built_in", "hello")
	if err == nil || pr.providers["built_in"] != nil {
		t.Fatal("Allowed registering built in")
	}

	// Should not allow registering built in through other API
	err = pr.registerStaticPropertyWithSource("built_in", datamodel.CMPropertySourceClient, "hello")
	if err == nil || pr.providers["built_in"] != nil {
		t.Fatal("Allowed registering built in")
	}

	// should be able to register well known with correct type
	err = pr.registerClientProperty("well_known", "hello")
	if err != nil {
		t.Fatal(err)
	}
	if pr.providers["well_known"] == nil {
		t.Fatal("Failed to register well known without a prefix")
	}
	if v, err := pr.propertyValue("well_known"); v != "hello" || err != nil {
		t.Fatal("Failed to register well known")
	}

	// should be able to register custom
	err = pr.registerClientProperty("customv", "hello2")
	if err != nil {
		t.Fatal(err)
	}
	if pr.providers["custom_customv"] == nil {
		t.Fatal("Failed to register well known without a prefix")
	}
	if pr.providers["customv"] != nil {
		t.Fatal("registered custom without a prefix")
	}
	if v, err := pr.propertyValue("customv"); v != "hello2" || err != nil {
		t.Fatal("Failed to access custom via short hand")
	}
	if v, err := pr.propertyValue("custom_customv"); v != "hello2" || err != nil {
		t.Fatal("Failed to access custom via full name")
	}

	// should not be able to regsiter nil
	err = pr.registerClientProperty("well_known2", nil)
	if err == nil || pr.providers["well_known2"] != nil {
		t.Fatal("Allowed nil value")
	}

	// Library register method should not be able to register non-built in
	err = pr.registerStaticProperty("customv", "hello3")
	if err == nil {
		t.Fatal("Allowed library to register custom")
	}
	err = pr.registerStaticPropertyWithSource("customv", datamodel.CMPropertySourceLib, "hello3")
	if err == nil {
		t.Fatal("Allowed library to register custom")
	}
	// old value from above, not new one.
	if v, err := pr.propertyValue("customv"); v != "hello2" || err != nil {
		t.Fatal("Failed to access custom via full name")
	}
}

func TestDisallowInvalidPropertyNames(t *testing.T) {
	pr := newPropertyRegistry()
	pr.builtInPropertyTypes = map[string]*datamodel.CMPropertyConfig{}

	invalidNames := []string{
		"with.period",
		"With space",
		"with-dash",
		"withEmoji😀",
		"withPunctuation!",
		"",
	}

	for _, n := range invalidNames {
		err := pr.registerClientProperty(n, "hello2")
		if err == nil {
			t.Fatal("allowed non alphanumeric property name" + n)
		}
	}
	if !validPropertyName("edgesAZaz09_") {
		t.Fatal("valid property name failed")
	}
	for name := range datamodel.BuiltInPropertyTypes() {
		if !validPropertyName(name) {
			t.Fatalf("Invalid well known property name: %s", name)
		}
	}
}

func TestClientPropertyJsonRegistration(t *testing.T) {
	pr := newPropertyRegistry()
	pr.builtInPropertyTypes = map[string]*datamodel.CMPropertyConfig{
		"stringsKey": {Type: reflect.String, Source: datamodel.CMPropertySourceClient, Optional: true},
	}

	j := `{
		"stringKey": "stringVal",
		"invalidKey": {},
		"boolKey": true,
		"intKey": 42,
		"floatKey": 3.3
	}`

	err := pr.registerClientPropertiesFromJson(([]byte)(j))
	if err == nil {
		t.Fatal("json registration failed to error on invalid")
		// but we still expect some to succeed
	}
	if v, err := pr.propertyValue("stringKey"); v != "stringVal" || err != nil {
		t.Fatal("Failed to register json properties")
	}
	if v, err := pr.propertyValue("boolKey"); v != true || err != nil {
		t.Fatal("Failed to register json properties")
	}
	if v, err := pr.propertyValue("intKey"); v != 42.0 || err != nil {
		t.Fatal("Failed to register json properties")
	}
	if v, err := pr.propertyValue("floatKey"); v != 3.3 || err != nil {
		t.Fatal("Failed to register json properties")
	}
	if _, err := pr.propertyValue("invalidKey"); err == nil {
		t.Fatal("invalid registered")
	}
}

func testBuildTestDb(t *testing.T) *db.DB {
	dataPath := fmt.Sprintf("/tmp/criticalmoments/test-temp-%v", rand.Int())
	err := os.MkdirAll(dataPath, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	db := db.NewDB()
	err = db.StartWithPath(dataPath)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestSetPropertyHistoryOnStartup(t *testing.T) {
	pr := newPropertyRegistry()
	pr.builtInPropertyTypes = map[string]*datamodel.CMPropertyConfig{
		"on_start_prop":     {Type: reflect.String, Source: datamodel.CMPropertySourceLib, Optional: false, SampleType: datamodel.CMPropertySampleTypeAppStart},
		"on_access_prop":    {Type: reflect.String, Source: datamodel.CMPropertySourceLib, Optional: false, SampleType: datamodel.CMPropertySampleTypeOnUse},
		"never_sample_prop": {Type: reflect.String, Source: datamodel.CMPropertySourceLib, Optional: false, SampleType: datamodel.CMPropertySampleTypeDoNotSample},
		"int_prop":          {Type: reflect.Int, Source: datamodel.CMPropertySourceLib, Optional: false, SampleType: datamodel.CMPropertySampleTypeAppStart},
		"float_prop":        {Type: reflect.Float64, Source: datamodel.CMPropertySourceLib, Optional: false, SampleType: datamodel.CMPropertySampleTypeAppStart},
		"bool_prop":         {Type: reflect.Bool, Source: datamodel.CMPropertySourceLib, Optional: false, SampleType: datamodel.CMPropertySampleTypeAppStart},
		"date_prop":         {Type: datamodel.CMTimeKind, Source: datamodel.CMPropertySourceLib, Optional: false, SampleType: datamodel.CMPropertySampleTypeAppStart},
	}

	err := pr.registerStaticProperty("on_start_prop", "onstart")
	if err != nil {
		t.Fatal(err)
	}
	err = pr.registerStaticProperty("on_access_prop", "onaccess")
	if err != nil {
		t.Fatal(err)
	}
	err = pr.registerStaticProperty("never_sample_prop", "never")
	if err != nil {
		t.Fatal(err)
	}
	err = pr.registerStaticProperty("int_prop", 42)
	if err != nil {
		t.Fatal(err)
	}
	err = pr.registerStaticProperty("float_prop", 3.3)
	if err != nil {
		t.Fatal(err)
	}
	err = pr.registerStaticProperty("bool_prop", true)
	if err != nil {
		t.Fatal(err)
	}
	err = pr.registerStaticProperty("date_prop", time.UnixMilli(testTimestampUnixMilli))
	if err != nil {
		t.Fatal(err)
	}

	// Create test DB and connect property manager
	db := testBuildTestDb(t)
	pr.phm = db.PropertyHistoryManager()

	err = pr.samplePropertiesForStartup()
	if err != nil {
		t.Fatal(err)
	}

	// Check that the on start properties were set, and others were not
	if v, err := db.LatestPropertyHistory("on_start_prop"); err != nil || v != "onstart" {
		t.Fatal("on start property not set in db")
	}
	if v, err := db.LatestPropertyHistory("on_access_prop"); err != sql.ErrNoRows || v != nil {
		t.Fatal("on access property set in db even though it has not been accessed")
	}
	if v, err := db.LatestPropertyHistory("never_sample_prop"); err != sql.ErrNoRows || v != nil {
		t.Fatal("never sample prop in db")
	}
	if v, err := db.LatestPropertyHistory("int_prop"); err != nil || v != int64(42) {
		t.Fatal("int prop in db")
	}
	if v, err := db.LatestPropertyHistory("float_prop"); err != nil || v != 3.3 {
		t.Fatal("float prop in db")
	}
	if v, err := db.LatestPropertyHistory("bool_prop"); err != nil || v != true {
		t.Fatal("bool prop in db")
	}
	if v, err := db.LatestPropertyHistory("date_prop"); err != nil || (v.(time.Time).UnixMilli() != testTimestampUnixMilli) {
		t.Fatal("date prop in db")
	}

	result, err := pr.evaluateCondition(testHelperNewCondition("on_start_prop == 'onstart' && on_access_prop == 'onaccess' && never_sample_prop == 'never' && int_prop == 42 && float_prop == 3.3 && bool_prop", t))
	if err != nil || !result {
		t.Fatal("Properties not set correctly")
	}

	if v, err := db.LatestPropertyHistory("on_start_prop"); err != nil || v != "onstart" {
		t.Fatal("on start property not set in db")
	}
	// Check that the on access property was set
	if v, err := db.LatestPropertyHistory("on_access_prop"); err != nil || v != "onaccess" {
		t.Fatal("on access property was not set after access")
	}
	// Check we don't set the never sample property, even after access
	if v, err := db.LatestPropertyHistory("never_sample_prop"); err != sql.ErrNoRows || v != nil {
		t.Fatal("never sample prop in db")
	}

	// Check property history check method works for all types
	historyChecks := map[string]any{
		"on_start_prop":  "onstart",
		"on_access_prop": "onaccess",                             // add_test_count
		"int_prop":       42,                                     // add_test_count
		"float_prop":     3.3,                                    // add_test_count
		"bool_prop":      true,                                   // add_test_count
		"date_prop":      time.UnixMilli(testTimestampUnixMilli), // add_test_count
	}
	for k, v := range historyChecks {
		has, err := db.PropertyHistoryEverHadValue(k, v)
		if err != nil {
			t.Fatal(err)
		}
		if !has {
			t.Fatal("Property history check failed for " + k)
		}
	}

	// Property value it has never had
	has, err := db.PropertyHistoryEverHadValue("on_start_prop", "asdf")
	if err != nil {
		t.Fatal(err)
	}
	if has {
		t.Fatal("Property history check failed for mismatched value")
	}
}
