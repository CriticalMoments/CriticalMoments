package appcore

import (
	"errors"
	"reflect"
	"testing"
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
	pr.wellKnownPropertyTypes = map[string]reflect.Kind{"a": reflect.String, "b": reflect.Int, "c": reflect.Float64, "d": reflect.Bool}
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
	pr.requiredPropertyTypes = map[string]reflect.Kind{
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

func TestPropertyRegistryValidateWellKnown(t *testing.T) {
	pr := newPropertyRegistry()
	pr.requiredPropertyTypes = map[string]reflect.Kind{}
	pr.wellKnownPropertyTypes = map[string]reflect.Kind{
		"user_signed_in": reflect.Bool,
	}

	if pr.validateProperties() != nil {
		t.Fatal("Missing well known failed validation")
	}
	err := pr.registerStaticProperty("user_signed_in", 42)
	if err == nil {
		t.Fatal("Added with type mismatch")
	}
	pr.registerStaticProperty("user_signed_in", true)
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
	pr.requiredPropertyTypes = map[string]reflect.Kind{}
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
	pr.requiredPropertyTypes = map[string]reflect.Kind{}
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

func TestDynamicProperties(t *testing.T) {
	pr := newPropertyRegistry()
	pr.requiredPropertyTypes = map[string]reflect.Kind{}
	pr.wellKnownPropertyTypes = map[string]reflect.Kind{"a": reflect.Int}

	dp := testPropertyProvider{}
	pr.registerLibPropertyProvider("a", &dp)
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
	pr.requiredPropertyTypes = map[string]reflect.Kind{}
	pr.wellKnownPropertyTypes = map[string]reflect.Kind{
		"app_version":         reflect.String, // using as a populated string
		"screen_width_pixels": reflect.Int,    // using as a populated in
		"os_version":          reflect.String, // using as a nil string
	}

	pr.registerStaticProperty("app_version", "hello")
	pr.registerStaticProperty("screen_width_pixels", 42)
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
	pr.requiredPropertyTypes = map[string]reflect.Kind{}
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
}

func TestUnknownVarsInConditions(t *testing.T) {
	c, err := datamodel.NewCondition("unknown_var == nil")
	if err != nil {
		t.Fatal(err)
	}

	pr := newPropertyRegistry()
	pr.requiredPropertyTypes = map[string]reflect.Kind{}
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
	pr.requiredPropertyTypes = map[string]reflect.Kind{
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
	pr.requiredPropertyTypes = map[string]reflect.Kind{}

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
