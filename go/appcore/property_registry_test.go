package appcore

import (
	"reflect"
	"testing"
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

func TestPropertyRegistryVersionNumber(t *testing.T) {
	pr := newPropertyRegistry()
	pr.requiredPropertyTypes = map[string]reflect.Kind{}
	pr.wellKnownPropertyTypes = map[string]reflect.Kind{"c_version": reflect.String, "a_version": reflect.String}

	// Valid string -- saves and able to parse components with function
	if err := pr.registerStaticProperty("c_version", "1.2.3"); err != nil {
		t.Fatal("Valid version number failed to save")
	}
	if pr.propertyValue("c_version") != "1.2.3" {
		t.Fatal("Valid version number failed to save")
	}
	if r, err := pr.evaluateCondition("versionNumberComponent(c_version, 0) == 1"); err != nil || !r {
		t.Fatal("Valid version number failed to extract component")
	}
	if r, err := pr.evaluateCondition("versionNumberComponent(c_version, 0) == 2"); err != nil || r {
		t.Fatal("Valid version number failed to extract component that fails test")
	}
	if r, err := pr.evaluateCondition("versionNumberComponent(c_version, 1) == 2"); err != nil || !r {
		t.Fatal("Valid version number failed to extract component")
	}
	if r, err := pr.evaluateCondition("versionNumberComponent(c_version, 2) == 3"); err != nil || !r {
		t.Fatal("Valid version number failed to extract component")
	}
	if r, err := pr.evaluateCondition("versionNumberComponent(c_version, 3) == nil"); err != nil || !r {
		t.Fatal("Valid version number failed to extract component")
	}

	// Invalid version string
	if err := pr.registerStaticProperty("a_version", "1.b.3"); err != nil {
		t.Fatal("Invalid version number failed to save. Should still save as string for exact comparison")
	}
	if pr.propertyValue("a_version") != "1.b.3" {
		t.Fatal("Invalid version number failed to save. Should still save as string for exact comparison")
	}
	if r, err := pr.evaluateCondition("versionNumberComponent(a_version, 0) == nil"); err != nil || !r {
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
		"a_val":     reflect.String,
		"a_int":     reflect.Int,
		"a_missing": reflect.String,
	}

	pr.registerStaticProperty("a_val", "hello")
	pr.registerStaticProperty("a_int", 42)
	if pr.propertyValue("a_val") != "hello" {
		t.Fatal("property not set")
	}
	if pr.propertyValue("a_int") != 42 {
		t.Fatal("property not set")
	}

	_, err := pr.evaluateCondition("a > 2")
	if err == nil {
		t.Fatal("Allowed invalid conditions: nil > 2")
	}

	result, err := pr.evaluateCondition("3 > 2")
	if err != nil || !result {
		t.Fatal("Failed to eval simple true condition")
	}

	result, err = pr.evaluateCondition("1 > 2")
	if err != nil || result {
		t.Fatal("Failed to eval simple false condition")
	}

	result, err = pr.evaluateCondition("a_val == 'hello'")
	if err != nil || !result {
		t.Fatal("Failed to eval true condition")
	}

	result, err = pr.evaluateCondition("a_val startsWith 'hel'")
	if err != nil || !result {
		t.Fatal("Failed to eval true condition with builtin function")
	}

	result, err = pr.evaluateCondition("a_val == 'world'")
	if err != nil || result {
		t.Fatal("Failed to eval false condition with property")
	}

	_, err = pr.evaluateCondition("1 + 2 + 3")
	if err == nil {
		t.Fatal("Allowed condition with non bool result")
	}

	_, err = pr.evaluateCondition("a_val ^#$%")
	if err == nil {
		t.Fatal("Allowed invalid condition")
	}

	result, err = pr.evaluateCondition("a_int > 99")
	if err != nil || result {
		t.Fatal("false condition passed")
	}

	result, err = pr.evaluateCondition("a_int > 2")
	if err != nil || !result {
		t.Fatal("true condition failed")
	}

	result, err = pr.evaluateCondition("a_missing == nil")
	if err != nil || !result {
		t.Fatal("true condition for allowed missing var failed")
	}

	result, err = pr.evaluateCondition("a_missing ?? false")
	if err != nil || result {
		t.Fatal("nil condition did not eval to false")
	}

	result, err = pr.evaluateCondition("a_missing == false")
	if err != nil || result {
		t.Fatal("missing condition should be differentiateable from false bool")
	}
}
