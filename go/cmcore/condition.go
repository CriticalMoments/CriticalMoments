package cmcore

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/checker"
	"github.com/antonmedv/expr/conf"
	"github.com/antonmedv/expr/optimizer"
	"github.com/antonmedv/expr/parser"
	"golang.org/x/exp/maps"
)

func RequiredPropertyTypes() map[string]reflect.Kind {
	return map[string]reflect.Kind{
		"platform":              reflect.String,
		"os_version":            reflect.String,
		"device_manufacturer":   reflect.String,
		"device_model":          reflect.String,
		"locale_language_code":  reflect.String,
		"locale_country_code":   reflect.String,
		"locale_currency_code":  reflect.String,
		"app_version":           reflect.String,
		"user_interface_idiom":  reflect.String,
		"app_id":                reflect.String,
		"screen_width_pixels":   reflect.Int,
		"screen_height_pixels":  reflect.Int,
		"device_battery_state":  reflect.String,
		"device_battery_level":  reflect.Float64,
		"device_low_power_mode": reflect.Bool,
		"device_orientation":    reflect.String,
		"interface_orientation": reflect.String,
	}
}

func WellKnownPropertyTypes() map[string]reflect.Kind {
	return map[string]reflect.Kind{
		"device_model_class":   reflect.String,
		"device_model_version": reflect.String,
		"screen_width_points":  reflect.Int,
		"screen_height_points": reflect.Int,
		"screen_scale":         reflect.Float64,
	}
}

func ConditionEnvWithHelpers() map[string]interface{} {
	return map[string]interface{}{
		"versionNumberComponent": versionNumberComponent,
	}
}

func versionNumberComponent(versionString string, index int) interface{} {
	// Parse string in format "16.4.1" to get a specific component
	components := strings.Split(versionString, ".")
	intComponents := make([]int, len(components))
	for i, component := range components {
		intComponent, err := strconv.Atoi(component)
		if err != nil {
			fmt.Printf("CriticalMoments: Invalid version number format: \"%v\"\n", versionString)
			return nil
		}
		intComponents[i] = intComponent
	}

	if index >= len(intComponents) {
		return nil
	}
	return intComponents[index]
}

// An AST walker we use to analyize code, to see if it's compatible with CM
type cmAnalysisVisitor struct {
	variables map[string]bool
}

func (v *cmAnalysisVisitor) Visit(n *ast.Node) {
	if node, ok := (*n).(*ast.IdentifierNode); ok {
		// exclude methods
		helperMethod := ConditionEnvWithHelpers()[node.Value]
		if helperMethod == nil {
			v.variables[node.Value] = true
		}
	}
}

func ExtractVariablesFromCondition(code string) ([]string, error) {
	tree, err := parser.Parse(code)
	if err != nil {
		return nil, err
	}

	config := conf.New(conf.CreateNew())
	config.Strict = false
	_, err = checker.Check(tree, config)
	if err != nil {
		return nil, err
	}
	err = optimizer.Optimize(&tree.Node, config)
	if err != nil {
		return nil, err
	}

	visitor := &cmAnalysisVisitor{
		variables: make(map[string]bool),
	}
	ast.Walk(&tree.Node, visitor)
	return maps.Keys(visitor.variables), nil
}

func ValidateCondition(code string) error {
	variables, err := ExtractVariablesFromCondition(code)
	if err != nil {
		return err
	}

	allValidVariables := make(map[string]reflect.Kind)
	maps.Copy(allValidVariables, RequiredPropertyTypes())
	maps.Copy(allValidVariables, WellKnownPropertyTypes())

	for _, varName := range variables {
		if _, ok := allValidVariables[varName]; !ok {
			return errors.New(fmt.Sprintf("Variable included in expression which isn't recognized: %v", varName))
		}
	}
	return nil
}
