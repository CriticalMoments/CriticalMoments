package conditions

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type versionNumber struct {
	components []int
	postfix    string
}

func VersionFromVersionString(s string) (*versionNumber, error) {
	if s == "" {
		return nil, errors.New("invalid version: empty string")
	}

	// trim leading "v" if it exists
	if s[0] == 'v' {
		s = s[1:]
	}

	var version versionNumber
	// trim trailing "-beta" or similar if it exists
	dashIndex := strings.Index(s, "-")
	if dashIndex >= 0 {
		version.postfix = s[dashIndex:]
		s = s[:dashIndex]
	}

	if s == "" {
		return nil, errors.New("invalid version number")
	}

	// Parse string in format "16.4.1" to get each components
	components := strings.Split(s, ".")
	intComponents := make([]int, len(components))
	for i, component := range components {
		intComponent, err := strconv.Atoi(component)
		if err != nil {
			return nil, fmt.Errorf("invalid version component: %v", component)
		}
		intComponents[i] = intComponent
	}
	version.components = intComponents

	return &version, nil
}

func VersionNumberComponent(versionString string, index int) interface{} {
	// Parse string in format "16.4.1" to get a specific component
	v, err := VersionFromVersionString(versionString)
	if err != nil || v == nil {
		fmt.Printf("CriticalMoments: Invalid version number format: \"%v\"\n", versionString)
		return nil
	}

	if index >= len(v.components) {
		return nil
	}
	return v.components[index]
}

func VersionGreaterThan(a string, b string) bool {
	return versionCompareExpecting(a, b, 1)
}

func VersionLessThan(a string, b string) bool {
	return versionCompareExpecting(a, b, -1)
}

func VersionEqual(a string, b string) bool {
	return versionCompareExpecting(a, b, 0)
}

func versionCompareExpecting(a string, b string, target int) bool {
	r, err := versionCompare(a, b)
	if err != nil {
		fmt.Printf("CriticalMoments: Invalid version number format: \"%v\" or \"%v\"\n", a, b)
		return false
	}
	return r == target
}

// 0: equal
// 1: a > b
// -1: a < b
func versionCompare(a string, b string) (int, error) {
	av, err := VersionFromVersionString(a)
	if err != nil {
		return 0, err
	}
	bv, err := VersionFromVersionString(b)
	if err != nil {
		return 0, err
	}

	for i, aComponent := range av.components {
		if i >= len(bv.components) {
			// Matched up until this point, but a has more digits, a wins
			return 1, nil
		}
		bComponent := bv.components[i]
		if aComponent > bComponent {
			return 1, nil
		}
		if bComponent > aComponent {
			return -1, nil
		}
	}

	// Everything in a matches everything in b, but if b keeps going it wins
	if len(bv.components) > len(av.components) {
		return -1, nil
	}

	// All components are equal, but if one has a postfix, the one without is greater (v1 > v1-beta)
	if av.postfix == "" && bv.postfix != "" {
		return 1, nil
	}
	if bv.postfix == "" && av.postfix != "" {
		return -1, nil
	}

	// Equal
	return 0, nil
}
