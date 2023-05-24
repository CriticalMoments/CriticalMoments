package cmcore

import (
	"errors"
	"fmt"
	"math"
	"reflect"

	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/checker"
	"github.com/antonmedv/expr/conf"
	"github.com/antonmedv/expr/optimizer"
	"github.com/antonmedv/expr/parser"
)

const CMKindVersionNumber reflect.Kind = math.MaxInt

func RequiredPropertyTypes() map[string]reflect.Kind {
	return map[string]reflect.Kind{
		"platform":              reflect.String,
		"os_version":            CMKindVersionNumber,
		"device_manufacturer":   reflect.String,
		"device_model":          reflect.String,
		"locale_language_code":  reflect.String,
		"locale_country_code":   reflect.String,
		"locale_currency_code":  reflect.String,
		"app_version":           CMKindVersionNumber,
		"user_interface_idiom":  reflect.String,
		"app_id":                reflect.String,
		"screen_width_pixels":   reflect.Int,
		"screen_height_pixels":  reflect.Int,
		"device_battery_state":  reflect.String,
		"device_battery_level":  reflect.Float64,
		"device_low_power_mode": reflect.Bool,
	}
}

func WellKnownPropertyTypes() map[string]reflect.Kind {
	return map[string]reflect.Kind{
		"user_signed_in":       reflect.Bool,
		"device_model_class":   reflect.String,
		"device_model_version": CMKindVersionNumber,
		"screen_width_points":  reflect.Int,
		"screen_height_points": reflect.Int,
		"screen_scale":         reflect.Float64,
	}
}

type cmExprEnv struct{}

// To add methods, add them to our ExprEnv. Example for testing:
func (cmExprEnv) AddOne(i int) int { return i + 1 }

// An AST walker we use to analyize code, to see if it's compatible with CM
type cmAnalysisVisitor struct {
	variables []string
}

func (v *cmAnalysisVisitor) Visit(n *ast.Node) {
	if node, ok := (*n).(*ast.IdentifierNode); ok {
		if !node.Method {
			v.variables = append(v.variables, node.Value)
		}
	}
}

func ExtractVariablesFromCondition(code string) ([]string, error) {
	tree, err := parser.Parse(code)
	if err != nil {
		return nil, err
	}

	config := conf.New(cmExprEnv{})
	config.Strict = false
	_, err = checker.Check(tree, config)
	if err != nil {
		return nil, err
	}
	err = optimizer.Optimize(&tree.Node, config)
	if err != nil {
		return nil, err
	}

	visitor := &cmAnalysisVisitor{}
	ast.Walk(&tree.Node, visitor)
	return visitor.variables, nil
}

func ValidateCondition(code string) error {
	variables, err := ExtractVariablesFromCondition(code)
	if err != nil {
		return err
	}

	required := RequiredPropertyTypes()
	wellKnown := WellKnownPropertyTypes()

	for _, v := range variables {
		// TODO: expand version strings
		if _, ok := required[v]; ok {
			continue
		}
		if _, ok := wellKnown[v]; ok {
			continue
		}
		return errors.New(fmt.Sprintf("Variable included in expression which isn't recognized: %v", v))
	}
	return nil
}
