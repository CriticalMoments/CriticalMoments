package datamodel

import (
	"encoding/json"
	"fmt"
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
	fields, err := c.ExtractIdentifiers()
	if err != nil {
		return nil, err
	}
	return fields.Variables, nil
}

func TestConditionConstructor(t *testing.T) {
	c, err := NewCondition("")
	if err == nil || c != nil {
		t.Fatal("Empty strings not valid conditions")
	}

	c, err = NewCondition("bad_var > 2")
	if err != nil || c.String() != "bad_var > 2" {
		t.Fatal("Unknown var not allowed in Strict=false")
	}

	c, err = NewCondition("true")
	if err != nil || c.String() != "true" {
		t.Fatal("Valid condition failed")
	}

	StrictDatamodelParsing = true
	defer func() {
		StrictDatamodelParsing = false
	}()

	c, err = NewCondition("bad_var > 2")
	if err == nil || c != nil {
		t.Fatal("Unknown var not valid conditions")
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

	// method names should not be included
	// build in methods (startsWith) should not
	// repeated var a should only be listed once
	code = "a || ab() || versionNumberComponent(1) > 1 || a startsWith 'hello'"
	variables, err = extractVarsTestHelper(code)
	if err != nil {
		t.Fatal(err)
	}
	if !arraysEqualOrderInsensitive(variables, []string{"a"}) {
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
	if err != nil {
		t.Fatal("Invalid prop didn't pass non strict validation")
	}

	// TODO: test this returns nil
	err = validateTestHelper("AddTwo(1) > 1")
	if err != nil {
		t.Fatal("Unrecognized method failed non-strict validation")
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

	StrictDatamodelParsing = true
	defer func() {
		StrictDatamodelParsing = false
	}()
	err = validateTestHelper("not_a_supported_prop > 1")
	if err == nil {
		t.Fatal("Invalid prop passed strict validation")
	}
	err = validateTestHelper("AddTwo(1) > 1")
	if err == nil {
		t.Fatal("Unrecognized method passed validation")
	}
}

func TestParseCondtion(t *testing.T) {
	var c Condition

	s := "true"
	err := json.Unmarshal([]byte(fmt.Sprintf("\"%v\"", s)), &c)
	if err != nil || c.String() != s {
		t.Fatal("Parse condition failed", err)
	}
	s = "true && false"
	err = json.Unmarshal([]byte(fmt.Sprintf("\"%v\"", s)), &c)
	if err != nil || c.String() != s {
		t.Fatal("Parse condition failed", err)
	}

	s = "true && false || 5 > 9"
	err = json.Unmarshal([]byte(fmt.Sprintf("\"%v\"", s)), &c)
	if err != nil || c.String() != s {
		t.Fatal("Parse condition failed", err)
	}

	// Unknown vars allowed in non-strict mode
	s = "unknown_var > 6"
	err = json.Unmarshal([]byte(fmt.Sprintf("\"%v\"", s)), &c)
	if err != nil || c.String() != s {
		t.Fatal("Parse condition failed", err)
	}

	// invalid conditions should fallback to false if not in strict mode
	c = Condition{}
	s = "'qwert' > 3"
	err = json.Unmarshal([]byte(fmt.Sprintf("\"%v\"", s)), &c)
	if err != nil || c.String() != "" {
		t.Fatal("Parse bad condition did not err")
	}
	c = Condition{}
	s = "app_version ^#$%"
	err = json.Unmarshal([]byte(fmt.Sprintf("\"%v\"", s)), &c)
	if err != nil || c.String() != "" {
		t.Fatal("Parse bad condition did not err")
	}
	c = Condition{}
	s = ""
	err = json.Unmarshal([]byte(fmt.Sprintf("\"%v\"", s)), &c)
	if err != nil || c.String() != "" {
		t.Fatal("Parse allowed non JSON formated string")
	}

	c = Condition{}
	// invalid json errors
	err = json.Unmarshal([]byte(""), &c)
	if err == nil || c.String() != "" {
		t.Fatal("Parse allowed non JSON formated string")
	}

	// Unknown vars allowed in non-strict mode
	c = Condition{}
	s = "unknown_var > 6"
	err = json.Unmarshal([]byte(fmt.Sprintf("\"%v\"", s)), &c)
	if err != nil || c.String() != s {
		t.Fatal("Parse condition failed", err)
	}

	StrictDatamodelParsing = true
	defer func() {
		StrictDatamodelParsing = false
	}()

	// Unknown vars not allowed in strict mode
	c = Condition{}
	s = "unknown_var > 6"
	err = json.Unmarshal([]byte(fmt.Sprintf("\"%v\"", s)), &c)
	if err == nil || c.String() != "" {
		t.Fatal("Strict mode ignored unknown var", err)
	}

	// invalid conditions should fail in strict mode
	c = Condition{}
	s = "'qwert' > 3"
	err = json.Unmarshal([]byte(fmt.Sprintf("\"%v\"", s)), &c)
	if err == nil || c.String() != "" {
		t.Fatal("Parse bad condition did not err")
	}

	c = Condition{}
	s = "app_version ^#$%"
	err = json.Unmarshal([]byte(fmt.Sprintf("\"%v\"", s)), &c)
	if err == nil || c.String() != "" {
		t.Fatal("Parse bad condition did not err")
	}

	c = Condition{}
	s = ""
	err = json.Unmarshal([]byte(fmt.Sprintf("\"%v\"", s)), &c)
	if err == nil || c.String() != "" {
		t.Fatal("Parse allowed non JSON formated string")
	}

	c = Condition{}
	err = json.Unmarshal([]byte(""), &c)
	if err == nil || c.String() != "" {
		t.Fatal("Parse allowed non JSON formated string")
	}
}

func TestExtractIdentifiers(t *testing.T) {
	c := Condition{
		conditionString: "func1() && func2(5) && both(1) && var1 == 3 && both < 4 && 4 > var2",
	}
	fields, err := c.ExtractIdentifiers()
	if err != nil {
		t.Fatal(err)
	}
	if !arraysEqualOrderInsensitive(fields.Methods, []string{"func1", "func2", "both"}) {
		t.Fatal("Extract methods failed")
	}
	if !arraysEqualOrderInsensitive(fields.Variables, []string{"var1", "var2", "both"}) {
		t.Fatal("Extract variables failed")
	}
}
