package appcore

import (
	"reflect"
	"testing"
)

func TestPropertyRegistrySetGet(t *testing.T) {
	pr := newPropertyRegistry()
	pr.registerStaticProperty("a", "a")
	if pr.propertyValue("a") != "a" {
		t.Fatal("Property registry failed for string")
	}
	pr.registerStaticProperty("b", 2)
	if pr.propertyValue("b") != 2 {
		t.Fatal("Property registry failed for int")
	}
	pr.registerStaticProperty("c", 3.3)
	if pr.propertyValue("c") != 3.3 {
		t.Fatal("Property registry failed for int")
	}
}

func TestPropertyRegistryValidateRequired(t *testing.T) {
	pr := newPropertyRegistry()
	pr.requiredPropertyTypes = map[string]reflect.Kind{
		"platform": reflect.String,
	}
	pr.wellKnownPropertyTypes = map[string]reflect.Kind{}

	if pr.validatePropertiesReturningUserReadable() == "" {
		t.Fatal("Validated missing required properties")
	}
	pr.registerStaticProperty("platform", 42)
	if pr.validatePropertiesReturningUserReadable() == "" {
		t.Fatal("Validated with type mismatch")
	}
	pr.registerStaticProperty("platform", "ios")
	if pr.validatePropertiesReturningUserReadable() != "" {
		t.Fatal("Validation failed on valid type")
	}
}

func TestPropertyRegistryValidateWellKnown(t *testing.T) {
	pr := newPropertyRegistry()
	pr.requiredPropertyTypes = map[string]reflect.Kind{}
	pr.wellKnownPropertyTypes = map[string]reflect.Kind{
		"user_signed_in": reflect.Bool,
	}

	if pr.validatePropertiesReturningUserReadable() != "" {
		t.Fatal("Missing well known failed validation")
	}
	pr.registerStaticProperty("user_signed_in", 42)
	if pr.validatePropertiesReturningUserReadable() == "" {
		t.Fatal("Validated with type mismatch")
	}
	pr.registerStaticProperty("user_signed_in", true)
	if pr.validatePropertiesReturningUserReadable() != "" {
		t.Fatal("Validation failed on valid type")
	}
}

func TestPropertyRegistryVersionNumber(t *testing.T) {
	pr := newPropertyRegistry()
	pr.requiredPropertyTypes = map[string]reflect.Kind{}
	pr.wellKnownPropertyTypes = map[string]reflect.Kind{}

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
	if pr.propertyValue("a_version_major") != nil {
		t.Fatal("Failed version numbers saved")
	}
	if pr.propertyValue("a_version_string") != "1.1.a" {
		t.Fatal("Failed version number failed to at least save string version")
	}

	// Valid
	if err := pr.registerStaticVersionNumberProperty("b", "1.2.3"); err != nil {
		t.Fatal("Valid version number failed to save")
	}
	if pr.propertyValue("b_version_string") != "1.2.3" {
		t.Fatal("Valid version number failed to save component")
	}
	if pr.propertyValue("b_version_major") != 1 {
		t.Fatal("Valid version number failed to save component")
	}
	if pr.propertyValue("b_version_minor") != 2 {
		t.Fatal("Valid version number failed to save component")
	}
	if pr.propertyValue("b_version_patch") != 3 {
		t.Fatal("Valid version number failed to save component")
	}
	if pr.propertyValue("b_version_micro") != nil {
		t.Fatal("Valid version saved extra component")
	}

	// Very long -- should save 7 deep and not error
	if err := pr.registerStaticVersionNumberProperty("c", "1.2.3.4.5.6.7.8.9.10.11.12"); err != nil {
		t.Fatal("Valid version number failed to save")
	}
	if pr.propertyValue("c_version_smol") != 7 {
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
	pr.wellKnownPropertyTypes = map[string]reflect.Kind{}

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
