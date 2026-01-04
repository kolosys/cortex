package cortex

import "time"

// Result contains the outcome of a rule evaluation.
type Result struct {
	// ID is the evaluation context ID.
	ID string

	// Success indicates if all rules evaluated without error.
	Success bool

	// RulesEvaluated is the number of rules that were run.
	RulesEvaluated int

	// RulesFailed is the number of rules that failed.
	RulesFailed int

	// Errors contains all collected errors (in CollectAll mode).
	Errors []RuleError

	// Duration is the total evaluation time.
	Duration time.Duration

	// HaltedBy is the rule ID that halted evaluation (if any).
	HaltedBy string

	// Context is the final evaluation context state.
	Context *EvalContext
}

// HasErrors returns true if any errors were collected.
func (r *Result) HasErrors() bool {
	return len(r.Errors) > 0
}

// FirstError returns the first error, or nil if none.
func (r *Result) FirstError() error {
	if len(r.Errors) == 0 {
		return nil
	}
	return &r.Errors[0]
}

// ErrorMessages returns all error messages.
func (r *Result) ErrorMessages() []string {
	msgs := make([]string, len(r.Errors))
	for i, e := range r.Errors {
		msgs[i] = e.Error()
	}
	return msgs
}

// newResult creates a new Result from an EvalContext.
func newResult(evalCtx *EvalContext, errors []RuleError) *Result {
	return &Result{
		ID:             evalCtx.ID,
		Success:        len(errors) == 0 && !evalCtx.IsHalted(),
		RulesEvaluated: int(evalCtx.RulesEvaluated()),
		RulesFailed:    len(errors),
		Errors:         errors,
		Duration:       evalCtx.Duration(),
		HaltedBy:       evalCtx.HaltedBy(),
		Context:        evalCtx,
	}
}
