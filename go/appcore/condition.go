package appcore

import (
	"errors"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"golang.org/x/exp/maps"
)

func extractVariablesFromCode(code string) ([]string, error) {
	emptyEnv := expr.Env(map[string]interface{}{})

	program, err := expr.Compile(code, emptyEnv, expr.AllowUndefinedVariables(), expr.AsBool())
	if err != nil {
		return nil, err
	}

	variableMap := map[string]bool{}
	for i, bytecode := range program.Bytecode {
		// TODO: other loads (field?)
		// OpLoadConst -> runtime.Fetch(env, program.Constants[arg])
		// OpLoadField -> runtime.FetchField(env, program.Constants[arg].(*runtime.Field))
		// OpFetchField -> runtime.FetchField(a, program.Constants[arg].(*runtime.Field)

		// find all opcodeloadfast which loads variable from env, get var name
		if bytecode == vm.OpLoadFast {
			if i >= len(program.Arguments) {
				return nil, errors.New("Unexpected issue extracting variables from condition")
			}
			arg := program.Arguments[i]
			if arg >= len(program.Constants) {
				return nil, errors.New("Unexpected issue extracting variables from condition")
			}
			varName, ok := program.Constants[arg].(string)
			if ok {
				variableMap[varName] = true
			} else {
				return nil, errors.New("Unexpected issue extracting variables from condition")
			}
		}
	}
	variables := maps.Keys(variableMap)
	return variables, nil
}
