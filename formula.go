package cortex

import (
	"context"
	"fmt"

	"github.com/kolosys/cortex/expr"
)

// FormulaFunc computes a value from the evaluation context.
type FormulaFunc func(ctx context.Context, evalCtx *EvalContext) (any, error)

// FormulaRule calculates a value using a function or expression.
type FormulaRule struct {
	baseRule
	target       string
	inputs       []string // required input keys (for dependency tracking)
	formula      FormulaFunc
	expression   string           // for config-driven rules
	compiledExpr *expr.Expression // compiled expression
}

// FormulaConfig configures a formula rule.
type FormulaConfig struct {
	ID          string
	Name        string
	Description string
	Deps        []string

	// Target is the context key to store the result.
	Target string

	// Inputs are the required input keys (for dependency tracking).
	Inputs []string

	// Formula is the Go function for complex rules.
	Formula FormulaFunc

	// Expression is the expression string for config-driven rules.
	// (Mutually exclusive with Formula)
	Expression string
}

// NewFormula creates a new formula rule.
func NewFormula(cfg FormulaConfig) (*FormulaRule, error) {
	if cfg.ID == "" {
		return nil, fmt.Errorf("%w: formula rule requires ID", ErrInvalidRule)
	}
	if cfg.Target == "" {
		return nil, fmt.Errorf("%w: formula rule %q requires target", ErrInvalidRule, cfg.ID)
	}
	if cfg.Formula == nil && cfg.Expression == "" {
		return nil, fmt.Errorf("%w: formula rule %q requires formula or expression", ErrInvalidRule, cfg.ID)
	}

	var compiledExpr *expr.Expression
	if cfg.Expression != "" && cfg.Formula == nil {
		var err error
		compiledExpr, err = expr.Compile(cfg.Expression)
		if err != nil {
			return nil, fmt.Errorf("%w: formula rule %q expression error: %v", ErrInvalidExpression, cfg.ID, err)
		}
	}

	return &FormulaRule{
		baseRule: baseRule{
			id:          cfg.ID,
			name:        cfg.Name,
			description: cfg.Description,
			deps:        cfg.Deps,
		},
		target:       cfg.Target,
		inputs:       cfg.Inputs,
		formula:      cfg.Formula,
		expression:   cfg.Expression,
		compiledExpr: compiledExpr,
	}, nil
}

// MustFormula creates a new formula rule, panicking on error.
func MustFormula(cfg FormulaConfig) *FormulaRule {
	r, err := NewFormula(cfg)
	if err != nil {
		panic(err)
	}
	return r
}

// Evaluate computes and stores the formula result.
func (r *FormulaRule) Evaluate(ctx context.Context, evalCtx *EvalContext) error {
	var result any
	var err error

	if r.formula != nil {
		result, err = r.formula(ctx, evalCtx)
	} else if r.compiledExpr != nil {
		result, err = r.compiledExpr.Eval(ctx, evalCtx)
	} else {
		return NewRuleError(r.id, string(RuleTypeFormula), "evaluate",
			fmt.Errorf("no formula or expression configured"))
	}

	if err != nil {
		return NewRuleError(r.id, string(RuleTypeFormula), "evaluate", err)
	}

	evalCtx.Set(r.target, result)
	return nil
}

// Target returns the target key for this formula.
func (r *FormulaRule) Target() string {
	return r.target
}

// Inputs returns the required input keys.
func (r *FormulaRule) Inputs() []string {
	return r.inputs
}

// Expression returns the expression string (if any).
func (r *FormulaRule) Expression() string {
	return r.expression
}

// SetFormulaFunc sets the formula function (used by expression parser).
func (r *FormulaRule) SetFormulaFunc(fn FormulaFunc) {
	r.formula = fn
}

// Common formula helpers

// Add returns a FormulaFunc that adds two context values.
func Add(a, b string) FormulaFunc {
	return func(ctx context.Context, evalCtx *EvalContext) (any, error) {
		va, err := evalCtx.GetFloat64(a)
		if err != nil {
			return nil, err
		}
		vb, err := evalCtx.GetFloat64(b)
		if err != nil {
			return nil, err
		}
		return va + vb, nil
	}
}

// Subtract returns a FormulaFunc that subtracts b from a.
func Subtract(a, b string) FormulaFunc {
	return func(ctx context.Context, evalCtx *EvalContext) (any, error) {
		va, err := evalCtx.GetFloat64(a)
		if err != nil {
			return nil, err
		}
		vb, err := evalCtx.GetFloat64(b)
		if err != nil {
			return nil, err
		}
		return va - vb, nil
	}
}

// Multiply returns a FormulaFunc that multiplies two context values.
func Multiply(a, b string) FormulaFunc {
	return func(ctx context.Context, evalCtx *EvalContext) (any, error) {
		va, err := evalCtx.GetFloat64(a)
		if err != nil {
			return nil, err
		}
		vb, err := evalCtx.GetFloat64(b)
		if err != nil {
			return nil, err
		}
		return va * vb, nil
	}
}

// Divide returns a FormulaFunc that divides a by b.
func Divide(a, b string) FormulaFunc {
	return func(ctx context.Context, evalCtx *EvalContext) (any, error) {
		va, err := evalCtx.GetFloat64(a)
		if err != nil {
			return nil, err
		}
		vb, err := evalCtx.GetFloat64(b)
		if err != nil {
			return nil, err
		}
		if vb == 0 {
			return nil, ErrDivisionByZero
		}
		return va / vb, nil
	}
}

// Percentage returns a FormulaFunc that calculates a percentage of a value.
func Percentage(value string, percent float64) FormulaFunc {
	return func(ctx context.Context, evalCtx *EvalContext) (any, error) {
		v, err := evalCtx.GetFloat64(value)
		if err != nil {
			return nil, err
		}
		return v * percent / 100, nil
	}
}

// Conditional returns a FormulaFunc that returns thenVal if condition is true, else elseVal.
func Conditional(condition string, thenVal, elseVal any) FormulaFunc {
	return func(ctx context.Context, evalCtx *EvalContext) (any, error) {
		cond, err := evalCtx.GetBool(condition)
		if err != nil {
			return nil, err
		}
		if cond {
			return thenVal, nil
		}
		return elseVal, nil
	}
}
