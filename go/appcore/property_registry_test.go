package appcore

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
	"unsafe"

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
	pr.wellKnownPropertyTypes = map[string]reflect.Kind{"a": reflect.String, "b": reflect.Int, "c": reflect.Float64, "d": reflect.Bool, "e": datamodel.CMTimeKind}
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
	pr.wellKnownPropertyTypes = map[string]reflect.Kind{"a": reflect.String, "b": reflect.Int}
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
	err = pr.registerStaticProperty("b", []string{}) // invalid type
	if err == nil {
		t.Fatal("Allowed invalid type")
	}
	if propertyValueOrNil(pr, "b") != nil {
		t.Fatal("Property registry allowed invalid")
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
	pr.builtInPropertyTypes = map[string]reflect.Kind{
		"platform": reflect.String,
	}
	pr.wellKnownPropertyTypes = map[string]reflect.Kind{}

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
	pr.builtInPropertyTypes = map[string]reflect.Kind{
		"low_data_mode": reflect.Bool, // this property is optional
	}
	pr.wellKnownPropertyTypes = map[string]reflect.Kind{}

	if pr.validateProperties() != nil {
		t.Fatal("Missing optional failed validation")
	}
	err := pr.registerStaticProperty("low_data_mode", 42)
	if err == nil {
		t.Fatal("Added with type mismatch")
	}
	err = pr.registerStaticProperty("low_data_mode", true)
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
	pr.builtInPropertyTypes = map[string]reflect.Kind{}
	pr.wellKnownPropertyTypes = map[string]reflect.Kind{}

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
	pr.builtInPropertyTypes = map[string]reflect.Kind{}
	pr.wellKnownPropertyTypes = map[string]reflect.Kind{"os_version": reflect.String, "app_version": reflect.String}

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
	pr.builtInPropertyTypes = map[string]reflect.Kind{}
	pr.wellKnownPropertyTypes = map[string]reflect.Kind{"a": reflect.Int}

	dp := testPropertyProvider{}
	err := pr.registerLibPropertyProvider("a", &dp)
	if err != nil {
		t.Fatal(err)
	}
	if propertyValueOrNil(pr, "a").(int64) != 1 {
		t.Fatal("dynamic property doesn't work")
	}
	if propertyValueOrNil(pr, "a").(int64) != 2 {
		t.Fatal("dynamic property not dynamic")
	}
	if propertyValueOrNil(pr, "a").(int64) != 3 {
		t.Fatal("dynamic property not dynamic")
	}
}

func TestPropertyRegistryConditionEval(t *testing.T) {
	pr := newPropertyRegistry()
	pr.builtInPropertyTypes = map[string]reflect.Kind{}
	pr.wellKnownPropertyTypes = map[string]reflect.Kind{
		"app_version":         reflect.String,       // using as a populated string
		"screen_width_pixels": reflect.Int,          // using as a populated in
		"os_version":          reflect.String,       // using as a nil string
		"test_time":           datamodel.CMTimeKind, // populated date
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
	pr.builtInPropertyTypes = map[string]reflect.Kind{}
	pr.wellKnownPropertyTypes = map[string]reflect.Kind{}

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
	pr.builtInPropertyTypes = map[string]reflect.Kind{}
	pr.wellKnownPropertyTypes = map[string]reflect.Kind{}

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
	pr.wellKnownPropertyTypes = map[string]reflect.Kind{}
	pr.builtInPropertyTypes = map[string]reflect.Kind{
		"platform": reflect.String,
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
	pr.wellKnownPropertyTypes = map[string]reflect.Kind{}
	pr.builtInPropertyTypes = map[string]reflect.Kind{}

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
	pr.wellKnownPropertyTypes = map[string]reflect.Kind{
		"well_known_v": reflect.String,
	}
	pr.builtInPropertyTypes = map[string]reflect.Kind{}

	// Test a custom properties with correct prefix
	err := pr.registerStaticProperty("custom_stringv", "hello")
	if err != nil {
		t.Fatal(err)
	}
	err = pr.registerStaticProperty("custom_boolv", false)
	if err != nil {
		t.Fatal(err)
	}
	err = pr.registerStaticProperty("custom_intv", 42)
	if err != nil {
		t.Fatal(err)
	}
	err = pr.registerStaticProperty("custom_floatv", 3.3)
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

	for k := range datamodel.WellKnownPropertyTypes() {
		if strings.HasPrefix(k, CustomPropertyPrefix) {
			t.Fatalf("Built in property has custom prefix: %s", k)
		}
	}
}

func TestValidateCustomPrefix(t *testing.T) {
	pr := newPropertyRegistry()
	pr.wellKnownPropertyTypes = map[string]reflect.Kind{}
	pr.builtInPropertyTypes = map[string]reflect.Kind{}

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
	pr.wellKnownPropertyTypes = map[string]reflect.Kind{
		"well_known":  reflect.String,
		"well_known2": reflect.String,
	}
	pr.builtInPropertyTypes = map[string]reflect.Kind{
		"built_in": reflect.String,
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

	// should be able to register well known
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
}

func TestDisallowInvalidPropertyNames(t *testing.T) {
	pr := newPropertyRegistry()
	pr.wellKnownPropertyTypes = map[string]reflect.Kind{}
	pr.builtInPropertyTypes = map[string]reflect.Kind{}

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
	for name := range datamodel.WellKnownPropertyTypes() {
		if !validPropertyName(name) {
			t.Fatalf("Invalid well known property name: %s", name)
		}
	}
	for name := range datamodel.BuiltInPropertyTypes() {
		if !validPropertyName(name) {
			t.Fatalf("Invalid well known property name: %s", name)
		}
	}
}

func TestClientPropertyJsonRegistration(t *testing.T) {
	pr := newPropertyRegistry()
	pr.wellKnownPropertyTypes = map[string]reflect.Kind{
		"stringKey": reflect.String,
	}
	pr.builtInPropertyTypes = map[string]reflect.Kind{}

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
