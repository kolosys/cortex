package cortex

import "context"

// Rule represents any rule that can be evaluated against an EvalContext.
type Rule interface {
	// ID returns the unique identifier for this rule.
	ID() string

	// Evaluate executes the rule against the provided context.
	Evaluate(ctx context.Context, evalCtx *EvalContext) error
}

// RuleMetadata provides optional metadata about a rule.
type RuleMetadata interface {
	Rule

	// Name returns a human-readable name for the rule.
	Name() string

	// Description returns a description of what the rule does.
	Description() string

	// Dependencies returns the IDs of other rules this rule depends on.
	Dependencies() []string
}

// RuleType identifies the type of rule.
type RuleType string

const (
	RuleTypeAssignment RuleType = "assignment"
	RuleTypeFormula    RuleType = "formula"
	RuleTypeAllocation RuleType = "allocation"
	RuleTypeLookup     RuleType = "lookup"
	RuleTypeBuildup    RuleType = "buildup"
)

// baseRule provides common fields for all rule types.
type baseRule struct {
	id          string
	name        string
	description string
	deps        []string
}

func (r *baseRule) ID() string             { return r.id }
func (r *baseRule) Name() string           { return r.name }
func (r *baseRule) Description() string    { return r.description }
func (r *baseRule) Dependencies() []string { return r.deps }
