package cortex

import (
	"errors"
	"fmt"
)

// Sentinel errors for common failure cases.
var (
	ErrRuleNotFound      = errors.New("cortex: rule not found")
	ErrLookupNotFound    = errors.New("cortex: lookup table not found")
	ErrKeyNotFound       = errors.New("cortex: key not found in lookup")
	ErrValueNotFound     = errors.New("cortex: value not found in context")
	ErrBuildupNotFound   = errors.New("cortex: buildup not found")
	ErrInvalidRule       = errors.New("cortex: invalid rule configuration")
	ErrInvalidExpression = errors.New("cortex: invalid expression")
	ErrTypeMismatch      = errors.New("cortex: type mismatch")
	ErrDivisionByZero    = errors.New("cortex: division by zero")
	ErrAllocationSum     = errors.New("cortex: allocation percentages must sum to 100")
	ErrCircularDep       = errors.New("cortex: circular dependency detected")
	ErrEvaluation        = errors.New("cortex: evaluation failed")
	ErrShortCircuit      = errors.New("cortex: evaluation short-circuited")
	ErrEngineClosed      = errors.New("cortex: engine is closed")
	ErrTimeout           = errors.New("cortex: evaluation timeout")
	ErrNilContext        = errors.New("cortex: nil evaluation context")
	ErrDuplicateRule     = errors.New("cortex: duplicate rule ID")
	ErrDuplicateLookup   = errors.New("cortex: duplicate lookup table name")
)

// RuleError wraps an error with rule context.
type RuleError struct {
	RuleID   string
	RuleType string
	Phase    string // "validate", "evaluate"
	Err      error
}

func (e *RuleError) Error() string {
	if e.RuleType != "" {
		return fmt.Sprintf("cortex: rule %q (%s) %s: %v", e.RuleID, e.RuleType, e.Phase, e.Err)
	}
	return fmt.Sprintf("cortex: rule %q %s: %v", e.RuleID, e.Phase, e.Err)
}

func (e *RuleError) Unwrap() error {
	return e.Err
}

// NewRuleError creates a new RuleError.
func NewRuleError(ruleID, ruleType, phase string, err error) *RuleError {
	return &RuleError{
		RuleID:   ruleID,
		RuleType: ruleType,
		Phase:    phase,
		Err:      err,
	}
}
