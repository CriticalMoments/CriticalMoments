package cmcore

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func arraysEqualOrderInsensitive(a []string, b []string) bool {
	less := func(aa, bb string) bool { return aa < bb }
	return "" == cmp.Diff(a, b, cmpopts.SortSlices(less))
}

func TestConditionVariableExtraction(t *testing.T) {
	code := "(a > 5555) && b && 'constantString' == c && 2 in [d, 3, 4]"
	variables, err := extractVariablesFromCode(code)
	if err != nil {
		t.Fatal(err)
	}
	if !arraysEqualOrderInsensitive(variables, []string{"a", "b", "c", "d"}) {
		t.Fatal("Extract variables failed")
	}

	code = "a && b.startsWith('constString') && c + d > 3"
	variables, err = extractVariablesFromCode(code)
	if err != nil {
		t.Fatal(err)
	}
	if !arraysEqualOrderInsensitive(variables, []string{"a", "b", "c", "d"}) {
		t.Fatal("Extract variables failed")
	}

	// It can optimize out the unneeded vars
	code = "a || (false && b + c + d > 0)"
	variables, err = extractVariablesFromCode(code)
	if err != nil {
		t.Fatal(err)
	}
	if !arraysEqualOrderInsensitive(variables, []string{"a"}) {
		t.Fatal("Extract variables failed")
	}

	// don't optimize out needed var
	code = "(a || false)"
	variables, err = extractVariablesFromCode(code)
	if err != nil {
		t.Fatal(err)
	}
	if !arraysEqualOrderInsensitive(variables, []string{"a"}) {
		t.Fatalf("Extract variables failed: %v", variables)
	}

	// unregistered method names should be included (ab), registered ones should not (AddOne)
	code = "a || ab() || AddOne(1) > 1"
	variables, err = extractVariablesFromCode(code)
	if err != nil {
		t.Fatal(err)
	}
	if !arraysEqualOrderInsensitive(variables, []string{"a", "ab"}) {
		t.Fatalf("Extract variables failed: %v", variables)
	}
}

func TestValidateProps(t *testing.T) {
	err := validateCodeCompatibleWithCMProps("1 < 2")
	if err != nil {
		t.Fatal("Simple case failed prop validation")
	}

	err = validateCodeCompatibleWithCMProps("not_a_supported_prop > 1")
	if err == nil {
		t.Fatal("Invalid prop passed validation")
	}

	err = validateCodeCompatibleWithCMProps("AddTwo(1) > 1")
	if err == nil {
		t.Fatal("Unrecognized method passed validation")
	}

	err = validateCodeCompatibleWithCMProps("AddOne(1) > 1")
	if err != nil {
		t.Fatal("Valid method failed validation")
	}

	err = validateCodeCompatibleWithCMProps("platform == 'iOS'")
	if err != nil {
		t.Fatal("Valid required property failed validation")
	}

	err = validateCodeCompatibleWithCMProps("screen_scale > 2.0")
	if err != nil {
		t.Fatal("Valid well known property failed validation")
	}

}
