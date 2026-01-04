package cortex

import (
	"context"
	"fmt"
)

// ValueFunc computes a value dynamically based on context.
type ValueFunc func(ctx context.Context, evalCtx *EvalContext) (any, error)

// AssignmentRule sets a value on the evaluation context.
type AssignmentRule struct {
	baseRule
	target    string
	value     any
	valueFunc ValueFunc
}

// AssignmentConfig configures an assignment rule.
type AssignmentConfig struct {
	ID          string
	Name        string
	Description string
	Deps        []string

	// Target is the context key to set.
	Target string

	// Value is the static value to assign (mutually exclusive with ValueFunc).
	Value any

	// ValueFunc computes the value dynamically (mutually exclusive with Value).
	ValueFunc ValueFunc
}

// NewAssignment creates a new assignment rule.
func NewAssignment(cfg AssignmentConfig) (*AssignmentRule, error) {
	if cfg.ID == "" {
		return nil, fmt.Errorf("%w: assignment rule requires ID", ErrInvalidRule)
	}
	if cfg.Target == "" {
		return nil, fmt.Errorf("%w: assignment rule %q requires target", ErrInvalidRule, cfg.ID)
	}
	if cfg.Value == nil && cfg.ValueFunc == nil {
		return nil, fmt.Errorf("%w: assignment rule %q requires value or value function", ErrInvalidRule, cfg.ID)
	}

	return &AssignmentRule{
		baseRule: baseRule{
			id:          cfg.ID,
			name:        cfg.Name,
			description: cfg.Description,
			deps:        cfg.Deps,
		},
		target:    cfg.Target,
		value:     cfg.Value,
		valueFunc: cfg.ValueFunc,
	}, nil
}

// MustAssignment creates a new assignment rule, panicking on error.
func MustAssignment(cfg AssignmentConfig) *AssignmentRule {
	r, err := NewAssignment(cfg)
	if err != nil {
		panic(err)
	}
	return r
}

// Evaluate sets the value on the evaluation context.
func (r *AssignmentRule) Evaluate(ctx context.Context, evalCtx *EvalContext) error {
	var value any
	var err error

	if r.valueFunc != nil {
		value, err = r.valueFunc(ctx, evalCtx)
		if err != nil {
			return NewRuleError(r.id, string(RuleTypeAssignment), "evaluate", err)
		}
	} else {
		value = r.value
	}

	evalCtx.Set(r.target, value)
	return nil
}

// Target returns the target key for this assignment.
func (r *AssignmentRule) Target() string {
	return r.target
}
