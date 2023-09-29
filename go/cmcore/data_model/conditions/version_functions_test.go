package conditions

import (
	"strings"
	"testing"
)

func TestVersionParsing(t *testing.T) {
	v, err := versionFromVersionString("v1")
	if err != nil || v == nil || v.components[0] != 1 || len(v.components) != 1 {
		t.Fatal("Failed to parse version string v1")
	}

	v, err = versionFromVersionString("v1-beta")
	if err != nil || v == nil || v.components[0] != 1 || len(v.components) != 1 || v.postfix != "-beta" {
		t.Fatal("Failed to parse version string v1")
	}

	v, err = versionFromVersionString("1.a.b")
	if err == nil || v != nil {
		t.Fatal("Failed to error on invalid version string 1.a.b")
	}

	v, err = versionFromVersionString("")
	if err == nil || v != nil {
		t.Fatal("Failed to error on invalid version string empty")
	}

	v, err = versionFromVersionString("v")
	if err == nil || v != nil {
		t.Fatal("Failed to error on invalid version string v")
	}

	v, err = versionFromVersionString("-beta")
	if err == nil || v != nil {
		t.Fatal("Failed to error on invalid version string -beta")
	}

	v, err = versionFromVersionString("v-beta")
	if err == nil || v != nil {
		t.Fatal("Failed to error on invalid version string v-beta")
	}

	v, err = versionFromVersionString("1.2.3")
	if err != nil || v.components[0] != 1 || v.components[1] != 2 || v.components[2] != 3 || len(v.components) != 3 || v.postfix != "" {
		t.Fatal("Failed to parse version string 1.2.3")
	}

	v, err = versionFromVersionString("v1.2.3")
	if err != nil || v.components[0] != 1 || v.components[1] != 2 || v.components[2] != 3 || len(v.components) != 3 || v.postfix != "" {
		t.Fatal("Failed to parse version string v1.2.3")
	}

	v, err = versionFromVersionString("v1.2.3-alpha.12-2")
	if err != nil || v.components[0] != 1 || v.components[1] != 2 || v.components[2] != 3 || len(v.components) != 3 || v.postfix != "-alpha.12-2" {
		t.Fatal("Failed to parse version string v1.2.3")
	}
}

func TestGetComponent(t *testing.T) {
	s := "1.2.3"
	r1 := VersionNumberComponent(s, 0)
	r2 := VersionNumberComponent(s, 1)
	r3 := VersionNumberComponent(s, 2)
	r4 := VersionNumberComponent(s, 3)
	if r1 != 1 || r2 != 2 || r3 != 3 || r4 != nil {
		t.Fatal("Failed to extract version number components")
	}
}

func TestVersionComparisonHelpers(t *testing.T) {
	s := "v1"
	l := "v2"

	if VersionGreaterThan(l, s) != true || VersionGreaterThan(s, l) != false {
		t.Fatal("Version Greater Than fails")
	}
	if VersionLessThan(s, l) != true || VersionLessThan(l, s) != false {
		t.Fatal("Version Less Than fails")
	}
	if VersionEqual(s, s) != true || VersionEqual(s, l) != false {
		t.Fatal("Version Greater Than fails")
	}
	if VersionGreaterThan("", "") || VersionEqual("", "") || VersionLessThan("", "") {
		t.Fatal("Invalid version not returning false")
	}
}

func TestVersionComparisons(t *testing.T) {
	// white box testing: each branch. 2=nil/error
	// no less thans because we test inverse of each
	cases := map[string]int{
		"invalid&v1":     2,
		"v1&v1":          0,
		"v1.2.3&1.2.3":   0,
		"v1.2&v1.2-beta": 1,
		"v1.2.3&v1.1":    1,
		"v1.3&v1.2":      1,
		"v1.1&v1":        1,
	}
	for c, expected := range cases {
		parts := strings.Split(c, "&")
		if len(parts) != 2 {
			t.Fatal("Invalid test case")
		}
		a := parts[0]
		b := parts[1]

		compareResult, err := versionCompare(a, b)
		if expected == 2 && err == nil {
			t.Fatalf("Invalid case did not return error: %v", c)
		}
		if expected != 2 && err != nil {
			t.Fatalf("valid case did return error: %v, %v", c, err)
		}
		if expected != 2 && expected != compareResult {
			t.Fatalf("valid case did return expected result: %v", c)
		}

		invertCompareResult, err := versionCompare(b, a)
		if expected == 2 && err == nil {
			t.Fatalf("Invalid case did not return error: %v", c)
		}
		if expected != 2 && err != nil {
			t.Fatalf("valid case did return error: %v", c)
		}
		if expected == 1 && invertCompareResult != -1 {
			t.Fatalf("valid case did return expected result: %v", c)
		}
		if expected == 0 && invertCompareResult != 0 {
			t.Fatalf("valid case did return expected result: %v", c)
		}
	}
}
