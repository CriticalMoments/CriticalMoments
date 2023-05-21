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
