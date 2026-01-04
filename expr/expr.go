// Package expr provides a simple expression DSL for cortex formulas.
//
// Supported operations:
//   - Arithmetic: +, -, *, /, %
//   - Comparison: ==, !=, <, >, <=, >=
//   - Logical: &&, ||, !
//   - Functions: min, max, abs, floor, ceil, round, if, sqrt, pow
//
// Example expressions:
//
//	"base_salary * tax_rate"
//	"if(age >= 65, senior_discount, 0)"
//	"round(total * 0.0825, 2)"
//	"min(calculated, max_amount)"
package expr

import (
	"context"
)

// Expression represents a compiled expression.
type Expression struct {
	raw       string
	ast       Node
	evaluator *Evaluator
}

// Compile parses and compiles an expression string.
func Compile(input string) (*Expression, error) {
	ast, err := Parse(input)
	if err != nil {
		return nil, err
	}

	return &Expression{
		raw:       input,
		ast:       ast,
		evaluator: NewEvaluator(),
	}, nil
}

// MustCompile compiles an expression, panicking on error.
func MustCompile(input string) *Expression {
	e, err := Compile(input)
	if err != nil {
		panic(err)
	}
	return e
}

// Raw returns the original expression string.
func (e *Expression) Raw() string {
	return e.raw
}

// Eval evaluates the expression against a value getter.
func (e *Expression) Eval(ctx context.Context, getter ValueGetter) (any, error) {
	return e.evaluator.Eval(ctx, e.ast, getter)
}

// EvalFloat64 evaluates the expression and returns a float64.
func (e *Expression) EvalFloat64(ctx context.Context, getter ValueGetter) (float64, error) {
	result, err := e.Eval(ctx, getter)
	if err != nil {
		return 0, err
	}
	return toFloat(result)
}

// EvalBool evaluates the expression and returns a bool.
func (e *Expression) EvalBool(ctx context.Context, getter ValueGetter) (bool, error) {
	result, err := e.Eval(ctx, getter)
	if err != nil {
		return false, err
	}
	b, ok := result.(bool)
	if !ok {
		return false, nil
	}
	return b, nil
}

// RegisterFunc registers a custom function for this expression.
func (e *Expression) RegisterFunc(name string, fn Func) {
	e.evaluator.RegisterFunc(name, fn)
}

// mapGetter wraps a map[string]any as a ValueGetter.
type mapGetter map[string]any

func (m mapGetter) Get(key string) (any, bool) {
	v, ok := m[key]
	return v, ok
}

// EvalWithMap evaluates the expression using a map as the value source.
func (e *Expression) EvalWithMap(ctx context.Context, values map[string]any) (any, error) {
	return e.Eval(ctx, mapGetter(values))
}
