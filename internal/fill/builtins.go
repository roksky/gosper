package fill

import (
	"context"
	"fmt"

	"github.com/roksky/gosper/internal/expression"
)

// ParameterDefault declares a parameter and its optional compiled default.
// Defaults evaluate in declaration order. A supplied override, including an
// explicit null, wins over the default.
type ParameterDefault struct {
	Name    string
	Default *expression.Program
}

type parameterSymbols map[string]expression.Value

func (s parameterSymbols) Field(string) (expression.Value, bool)    { return nil, false }
func (s parameterSymbols) Variable(string) (expression.Value, bool) { return nil, false }
func (s parameterSymbols) Parameter(name string) (expression.Value, bool) {
	v, ok := s[name]
	return v, ok
}

func cloneValues(src map[string]expression.Value) map[string]expression.Value {
	dst := make(map[string]expression.Value, len(src))
	for name, v := range src {
		dst[name] = v
	}
	return dst
}

func resolveParameters(ctx context.Context, decls []ParameterDefault, overrides map[string]expression.Value) (map[string]expression.Value, error) {
	params := cloneValues(overrides)
	for name, v := range params {
		if v == nil {
			params[name] = expression.NullValue
		}
	}
	seen := make(map[string]bool, len(decls))
	for _, decl := range decls {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if decl.Name == "" || seen[decl.Name] {
			return nil, fmt.Errorf("fill: empty or duplicate parameter name %q", decl.Name)
		}
		seen[decl.Name] = true
		if _, supplied := params[decl.Name]; supplied {
			continue
		}
		if decl.Default == nil {
			params[decl.Name] = expression.NullValue
			continue
		}
		v, err := decl.Default.Eval(ctx, expression.Scope{parameterSymbols(params)})
		if err != nil {
			return nil, fmt.Errorf("parameter %q default: %w", decl.Name, err)
		}
		params[decl.Name] = v
	}
	return params, nil
}
