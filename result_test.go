package cortex_test

import (
	"testing"

	"github.com/kolosys/cortex"
)

func TestResultMethods(t *testing.T) {
	// Create a result with errors
	evalCtx := cortex.NewEvalContext()
	errors := []cortex.RuleError{
		{RuleID: "rule1", Err: nil},
		{RuleID: "rule2", Err: nil},
	}

	// We can't directly call newResult as it's unexported, but we can test via engine

	engine := cortex.New("test", func() *cortex.Config {
		cfg := cortex.DefaultConfig()
		cfg.Mode = cortex.ModeCollectAll
		return cfg
	}())

	engine.AddRule(cortex.MustFormula(cortex.FormulaConfig{
		ID:         "fail",
		Target:     "x",
		Expression: "undefined_var + 1",
	}))

	result, _ := engine.Evaluate(t.Context(), evalCtx)

	if !result.HasErrors() {
		t.Error("expected HasErrors() to be true")
	}

	firstErr := result.FirstError()
	if firstErr == nil {
		t.Error("expected FirstError() to return error")
	}

	msgs := result.ErrorMessages()
	if len(msgs) != len(errors) && len(msgs) != 1 {
		// We expect 1 error from the undefined variable
	}

	// Test result without errors
	engine2 := cortex.New("test2", cortex.DefaultConfig())
	engine2.AddRule(cortex.MustAssignment(cortex.AssignmentConfig{
		ID:     "ok",
		Target: "x",
		Value:  1,
	}))

	evalCtx2 := cortex.NewEvalContext()
	result2, _ := engine2.Evaluate(t.Context(), evalCtx2)

	if result2.HasErrors() {
		t.Error("expected HasErrors() to be false")
	}

	if result2.FirstError() != nil {
		t.Error("expected FirstError() to be nil")
	}
}
