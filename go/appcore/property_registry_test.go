package appcore

import (
	"reflect"
	"testing"

	"github.com/CriticalMoments/CriticalMoments/go/cmcore"
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
	pr.wellKnownPropertyTypes = map[string]reflect.Kind{"a": cmcore.CMKindVersionNumber, "b": reflect.Bool, "c_version": cmcore.CMKindVersionNumber, "d_version": cmcore.CMKindVersionNumber}

	// invalid version should error
	if err := pr.registerStaticVersionNumberProperty("a", ""); err == nil {
		t.Fatal("Must provide a valid version number")
	}
	if err := pr.registerStaticVersionNumberProperty("a", "a.b.c"); err == nil {
		t.Fatal("Must provide a valid version number")
	}
	if err := pr.registerStaticVersionNumberProperty("a", "1.1.0x45"); err == nil {
		t.Fatal("Must provide a valid version number")
	}
	if err := pr.registerStaticVersionNumberProperty("a", "1.1.a"); err == nil {
		t.Fatal("Must provide a valid version number")
	}
	if err := pr.registerStaticVersionNumberProperty("", "1.1"); err == nil {
		t.Fatal("Must provide a prefix")
	}
	if pr.propertyValue("a_major") != nil {
		t.Fatal("Invalid version numbers saved component")
	}
	if pr.propertyValue("a_string") != "1.1.a" {
		t.Fatal("Failed version number failed to at least save string version")
	}

	// Invalid type
	if err := pr.registerStaticVersionNumberProperty("b", "1.2.3"); err == nil {
		t.Fatal("Allowed registering a version number to a bool type")
	}
	if pr.propertyValue("b_major") != nil {
		t.Fatal("Invalid version numbers saved component")
	}
	if pr.propertyValue("b_string") != nil {
		t.Fatal("Saved version for type mismatch")
	}

	// Valid
	if err := pr.registerStaticVersionNumberProperty("c_version", "1.2.3"); err != nil {
		t.Fatal("Valid version number failed to save")
	}
	if pr.propertyValue("c_version_string") != "1.2.3" {
		t.Fatal("Valid version number failed to save component")
	}
	if pr.propertyValue("c_version_major") != 1 {
		t.Fatal("Valid version number failed to save component")
	}
	if pr.propertyValue("c_version_minor") != 2 {
		t.Fatal("Valid version number failed to save component")
	}
	if pr.propertyValue("c_version_patch") != 3 {
		t.Fatal("Valid version number failed to save component")
	}
	if pr.propertyValue("c_version_micro") != nil {
		t.Fatal("Valid version saved extra component")
	}

	// Very long -- should save 7 deep and not error
	if err := pr.registerStaticVersionNumberProperty("d_version", "1.2.3.4.5.6.7.8.9.10.11.12"); err != nil {
		t.Fatal("Valid version number failed to save")
	}
	if pr.propertyValue("d_version_smol") != 7 {
		t.Fatal("Valid version number failed to save component")
	}
}

type testPropertyProvider struct {
	val int
}

func (p *testPropertyProvider) Type() int {
	return LibPropertyProviderTypeInt
}
func (p *testPropertyProvider) IntValue() int {
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
	if pr.propertyValue("a") != 1 {
		t.Fatal("dynamic property doesn't work")
	}
	if pr.propertyValue("a") != 2 {
		t.Fatal("dynamic property not dynamic")
	}
	if pr.propertyValue("a") != 3 {
		t.Fatal("dynamic property not dynamic")
	}
}

func TestPropertyRegistryConditionEval(t *testing.T) {
	pr := newPropertyRegistry()
	pr.requiredPropertyTypes = map[string]reflect.Kind{}
	pr.wellKnownPropertyTypes = map[string]reflect.Kind{
		"a_val": reflect.String,
		"a_int": reflect.Int,
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
		t.Fatal("Allowed invalid conditions")
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

	result, err = pr.evaluateCondition("1 + 2 + 3")
	if err == nil {
		t.Fatal("Allowed condition with non bool result")
	}

	result, err = pr.evaluateCondition("a_val ^#$%")
	if err == nil {
		t.Fatal("Allowed invalid condition")
	}

	result, err = pr.evaluateCondition("a_int > 99")
	if err != nil && result {
		t.Fatal("false condition passed")
	}

	result, err = pr.evaluateCondition("a_int > 2")
	if err != nil && !result {
		t.Fatal("true condition failed")
	}
}
