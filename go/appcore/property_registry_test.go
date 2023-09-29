package appcore

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	datamodel "github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model"
)

func TestPropertyRegistrySetGet(t *testing.T) {
	pr := newPropertyRegistry()
	pr.wellKnownPropertyTypes = map[string]reflect.Kind{"a": reflect.String, "b": reflect.Int, "c": reflect.Float64, "d": reflect.Bool}
	err := pr.registerStaticProperty("a", "a")
	if err != nil {
		t.Fatal(err)
	}
	if pr.propertyValue("a") != "a" {
		t.Fatal("Property registry failed for string")
	}
	err = pr.registerStaticProperty("b", 2)
	if err != nil {
		t.Fatal(err)
	}
	if pr.propertyValue("b") != 2 {
		t.Fatal("Property registry failed for int")
	}
	err = pr.registerStaticProperty("c", 3.3)
	if err != nil {
		t.Fatal(err)
	}
	if pr.propertyValue("c") != 3.3 {
		t.Fatal("Property registry failed for int")
	}
	err = pr.registerStaticProperty("d", true)
	if err != nil {
		t.Fatal(err)
	}
	if pr.propertyValue("d") != true {
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
	if pr.propertyValue("a") != nil {
		t.Fatal("Property registry allowed invalid")
	}
	err = pr.registerStaticProperty("b", []string{}) // invalid type
	if err == nil {
		t.Fatal("Allowed invalid type")
	}
	if pr.propertyValue("b") != nil {
		t.Fatal("Property registry allowed invalid")
	}
	err = pr.registerStaticProperty("c", 3.3) // unexpected key
	if err == nil {
		t.Fatal("Allowed unexpected key")
	}
	if pr.propertyValue("c") != nil {
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
		t.Fatal(fmt.Sprintf("Condition in test is not valid %v", s))
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
	if pr.propertyValue("os_version") != "1.2.3" {
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
	if pr.propertyValue("app_version") != "1.b.3" {
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
	if pr.propertyValue("a").(int64) != 1 {
		t.Fatal("dynamic property doesn't work")
	}
	if pr.propertyValue("a").(int64) != 2 {
		t.Fatal("dynamic property not dynamic")
	}
	if pr.propertyValue("a").(int64) != 3 {
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
	if pr.propertyValue("app_version") != "hello" {
		t.Fatal("property not set")
	}
	if pr.propertyValue("screen_width_pixels") != 42 {
		t.Fatal("property not set")
	}

	// Need relections to make an invalid condition, but want to keep test case
	badCondition := testHelperNewCondition("true", t)
	v := reflect.ValueOf(badCondition).Elem()
	cf := v.FieldByName("conditionString")
	cf = reflect.NewAt(cf.Type(), unsafe.Pointer(cf.UnsafeAddr())).Elem()
	cf.SetString("a > 2")
	_, err := pr.evaluateCondition(badCondition)
	if err == nil {
		t.Fatal("Allowed invalid conditions: nil > 2")
	}

	result, err := pr.evaluateCondition(testHelperNewCondition("3 > 2", t))
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

	_, err = pr.evaluateCondition(testHelperNewCondition("1 + 2 + 3", t))
	if err == nil {
		t.Fatal("Allowed condition with non bool result")
	}

	_, err = datamodel.NewCondition("app_version ^#$%")
	if err == nil {
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
	result, err := pr.evaluateCondition(testHelperNewCondition("now() < 1688764455000", t))
	if err != nil || result {
		t.Fatal("back to the future isn't real: now() gave incorrect time")
	}
	// time in 2050
	result, err = pr.evaluateCondition(testHelperNewCondition("now() > 2540841255000", t))
	if err != nil || result {
		t.Fatal("back to the future 2 isn't real: now() gave incorrect time")
	}

	verifyDurationCondition := testHelperNewCondition(`
		(seconds(1) == 1000) && 
		(seconds(9) == 9000) &&
		(minutes(1) == seconds(60)) &&
		(minutes(2) == 120000) &&
		(hours(1) == minutes(60)) &&
		(hours(2) == 7200000) &&
		(days(1) == hours(24)) &&
		(days(2) == 172800000)
	`, t)
	result, err = pr.evaluateCondition(verifyDurationCondition)
	if err != nil || !result {
		t.Fatal("Duration functions failed e2e test")
	}

	// Verify parseDate function is wired up properly. Actual testing of the
	// function is in date_functions_test.go
	result, err = pr.evaluateCondition(testHelperNewCondition("parseDate('2006-01-02T15:04:05.9997+07:00') == 1136189045999", t))
	if err != nil || !result {
		t.Fatal("parseDate did not work inside a condition")
	}

	// Check invalid returns errors though the stack
	result, err = pr.evaluateCondition(testHelperNewCondition("parseDate('invalid') == 1136189045999", t))
	if err == nil || result {
		t.Fatal("invalid date didn't error", err)
	}
}
