package datamodel

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func arraysEqualOrderInsensitive(a []string, b []string) bool {
	less := func(aa, bb string) bool { return aa < bb }
	return "" == cmp.Diff(a, b, cmpopts.SortSlices(less))
}

func extractVarsTestHelper(s string) ([]string, error) {
	// Not using the contstructor here, because some tests have invalid vars
	c := Condition{
		conditionString: s,
	}
	return c.ExtractVariables()
}

func TestConditionConstructor(t *testing.T) {
	c, err := NewCondition("")
	if err == nil || c != nil {
		t.Fatal("Empty strings not valid conditions")
	}

	c, err = NewCondition("bad_var > 2")
	if err == nil || c != nil {
		t.Fatal("Unknown var not valid conditions")
	}

	c, err = NewCondition("true")
	if err != nil || c.String() != "true" {
		t.Fatal("Valid condition failed")
	}
}

func TestConditionVariableExtraction(t *testing.T) {
	code := "(a > 5555) && b && 'constantString' == c && 2 in [d, 3, 4]"
	variables, err := extractVarsTestHelper(code)
	if err != nil {
		t.Fatal(err)
	}
	if !arraysEqualOrderInsensitive(variables, []string{"a", "b", "c", "d"}) {
		t.Fatal("Extract variables failed")
	}

	code = "a && b.startsWith('constString') && c + d > 3"
	variables, err = extractVarsTestHelper(code)
	if err != nil {
		t.Fatal(err)
	}
	if !arraysEqualOrderInsensitive(variables, []string{"a", "b", "c", "d"}) {
		t.Fatal("Extract variables failed")
	}

	// It can optimize out the unneeded vars
	code = "a || (false && b + c + d > 0)"
	variables, err = extractVarsTestHelper(code)
	if err != nil {
		t.Fatal(err)
	}
	if !arraysEqualOrderInsensitive(variables, []string{"a"}) {
		t.Fatal("Extract variables failed")
	}

	// don't optimize out needed var
	code = "(a || false)"
	variables, err = extractVarsTestHelper(code)
	if err != nil {
		t.Fatal(err)
	}
	if !arraysEqualOrderInsensitive(variables, []string{"a"}) {
		t.Fatalf("Extract variables failed: %v", variables)
	}

	// unregistered method names should be included (ab),
	// registered ones should not (versionNumberComponent)
	// build in methods (startsWith) should not
	// repeated var a should only be listed once
	code = "a || ab() || versionNumberComponent(1) > 1 || a startsWith 'hello'"
	variables, err = extractVarsTestHelper(code)
	if err != nil {
		t.Fatal(err)
	}
	if !arraysEqualOrderInsensitive(variables, []string{"a", "ab"}) {
		t.Fatalf("Extract variables failed: %v", variables)
	}
}

func validateTestHelper(s string) error {
	c, err := NewCondition(s)
	if err != nil {
		return err
	}
	return c.Validate()
}
func TestValidateProps(t *testing.T) {
	err := validateTestHelper("1 < 2")
	if err != nil {
		t.Fatal("Simple case failed prop validation")
	}

	err = validateTestHelper("not_a_supported_prop > 1")
	if err == nil {
		t.Fatal("Invalid prop passed validation")
	}

	err = validateTestHelper("AddTwo(1) > 1")
	if err == nil {
		t.Fatal("Unrecognized method passed validation")
	}

	err = validateTestHelper("versionNumberComponent('1.2.3', 1) == 1")
	if err != nil {
		t.Fatal("Valid method failed validation")
	}

	err = validateTestHelper("platform == 'iOS'")
	if err != nil {
		t.Fatal("Valid required property failed validation")
	}

	err = validateTestHelper("screen_scale > 2.0")
	if err != nil {
		t.Fatal("Valid well known property failed validation")
	}

	err = validateTestHelper("app_version == 'iPhone13,3' && versionNumberComponent(os_version, 1) >= 15")
	if err != nil {
		t.Fatal("Valid version strings failed validation")
	}
}
