package appcore

import (
	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/checker"
	"github.com/antonmedv/expr/conf"
	"github.com/antonmedv/expr/optimizer"
	"github.com/antonmedv/expr/parser"
)

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

func extractVariablesFromCode(code string) ([]string, error) {
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
