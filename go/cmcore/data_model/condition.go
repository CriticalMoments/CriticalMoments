package datamodel

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model/conditions"
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/checker"
	"github.com/antonmedv/expr/conf"
	"github.com/antonmedv/expr/optimizer"
	"github.com/antonmedv/expr/parser"
	"github.com/antonmedv/expr/vm"
	"golang.org/x/exp/maps"
)

type ConditionDynamicFunction struct {
	Function func(params ...any) (any, error)
	Types    []any
}

type Condition struct {
	conditionString string
}

func NewCondition(s string) (*Condition, error) {
	c := Condition{
		conditionString: s,
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	return &c, nil
}

// Stringer Interface
func (c *Condition) String() string {
	return c.conditionString
}

func StaticConditionHelperFunctions() map[string]interface{} {
	return map[string]interface{}{
		"versionNumberComponent": conditions.VersionNumberComponent,
		"versionGreaterThan":     conditions.VersionGreaterThan,
		"versionLessThan":        conditions.VersionLessThan,
		"versionEqual":           conditions.VersionEqual,

		"unixTimeNanoseconds":  conditions.UnixTimeNanoseconds,
		"unixTimeMilliseconds": conditions.UnixTimeMilliseconds,
		"unixTimeSeconds":      conditions.UnixTimeSeconds,
		"formatTime":           conditions.TimeFormat,

		"rand":        conditions.Random,
		"sessionRand": conditions.SessionRandom,
		"randForKey":  conditions.RandomForKey,
	}
}

func StaticConditionConstantProperties() map[string]interface{} {
	return map[string]interface{}{
		"RFC3339":              time.RFC3339Nano,
		"RFC822":               time.RFC822,
		"RFC850":               time.RFC850,
		"RFC1123":              time.RFC1123,
		"RFC822Z":              time.RFC822Z,
		"RFC1123Z":             time.RFC1123Z,
		"date_with_tz_format":  conditions.DateWithTzFormat,
		"date_and_time_format": conditions.DateAndTimeFormat,
		"date_format":          conditions.DateFormat,
	}
}

var AllBuiltInDynamicFunctions = map[string]bool{
	"eventCount":                 true,
	"eventCountWithLimit":        true,
	"canOpenUrl":                 true,
	"propertyHistoryLatestValue": true,
	"propertyEver":               true,
	"stableRand":                 true,
}

type ConditionFields struct {
	Identifiers []string
	Variables   []string
	Methods     []string
}

// An AST walker we use to analyze code, to see if it's compatible with CM
type conditionWalker struct {
	condition   string
	identifiers map[string]bool
	variables   map[string]bool
	methods     map[string]bool
}

func (v *conditionWalker) Visit(n *ast.Node) {
	if node, ok := (*n).(*ast.IdentifierNode); ok {
		v.identifiers[node.Value] = true

		// Check if this is a variable or a method. Unfortunately .Method() on the node does not work
		// so we check for open paren immediately after the identifier
		isMethod := false
		parenLoc := (*n).Location().Column + len(node.Value)
		if parenLoc < len(v.condition) {
			paran := v.condition[parenLoc : parenLoc+1]
			if paran == "(" {
				isMethod = true
			}
		}
		if isMethod {
			v.methods[node.Value] = true
		} else {
			v.variables[node.Value] = true
		}
	}
}

func (c *Condition) ExtractIdentifiers() (returnFields *ConditionFields, returnError error) {
	// expr can panic, so catch it and return an error instead
	defer func() {
		if r := recover(); r != nil {
			returnFields = nil
			returnError = fmt.Errorf("panic in ExtractIdentifiers: %v", r)
		}
	}()

	// single line needed because we use the location offset
	singleLineCondition := strings.ReplaceAll(c.conditionString, "\n", " ")
	tree, err := parser.Parse(singleLineCondition)
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

	visitor := &conditionWalker{
		condition:   singleLineCondition,
		identifiers: map[string]bool{},
		variables:   map[string]bool{},
		methods:     map[string]bool{},
	}
	ast.Walk(&tree.Node, visitor)

	results := ConditionFields{
		Identifiers: maps.Keys(visitor.identifiers),
		Variables:   maps.Keys(visitor.variables),
		Methods:     maps.Keys(visitor.methods),
	}
	return &results, nil
}

func (c *Condition) Validate() error {
	if c.conditionString == "" {
		return NewUserPresentableError("Condition is empty string (not allowed). Use 'true' or 'false' for minimal condition.")
	}

	// Run this even if not strict. It is checking the format of the condition as well
	fields, err := c.ExtractIdentifiers()
	if err != nil {
		return err
	}

	if StrictDatamodelParsing {
		// Don't check variable names. We support custom vars so every name is valid

		// Check we support all methods used if strict parsing
		for _, methodName := range fields.Methods {
			if _, ok := AllBuiltInDynamicFunctions[methodName]; !ok {
				if _, ok := StaticConditionHelperFunctions()[methodName]; !ok {
					return NewUserPresentableError(fmt.Sprintf("Method included in condition which isn't recognized: %v", methodName))
				}
			}
		}

	}

	return nil
}

func (c *Condition) CompileWithEnv(ops ...expr.Option) (resultProgram *vm.Program, returnError error) {
	// expr can panic, so catch it and return an error instead
	defer func() {
		if r := recover(); r != nil {
			resultProgram = nil
			returnError = fmt.Errorf("panic in CompileWithEnv: %v", r)
		}
	}()

	allOptions := append(ops, expr.AsBool())

	return expr.Compile(c.conditionString, allOptions...)
}

func (c *Condition) UnmarshalJSON(data []byte) error {
	var conditionString *string
	err := json.Unmarshal(data, &conditionString)
	if err != nil {
		return NewUserPresentableErrorWSource(fmt.Sprintf("Invalid Condition String [[ %s ]]", string(data)), err)
	}
	c.conditionString = *conditionString

	if err := c.Validate(); err != nil {
		// Fallback to returning empty on non-strict clients. Don't want entire config file to fail
		// Downstream during eval we return false and error
		c.conditionString = ""
		if StrictDatamodelParsing {
			return NewUserPresentableErrorWSource(fmt.Sprintf("Invalid Condition: [[ %v ]]", string(data)), err)
		}
	}

	return nil
}
