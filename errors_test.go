package cortex_test

import (
	"errors"
	"testing"

	"github.com/kolosys/cortex"
)

func TestRuleError(t *testing.T) {
	baseErr := errors.New("base error")

	ruleErr := cortex.NewRuleError("rule-1", "formula", "evaluate", baseErr)

	expectedMsg := `cortex: rule "rule-1" (formula) evaluate: base error`
	if ruleErr.Error() != expectedMsg {
		t.Errorf("expected %q, got %q", expectedMsg, ruleErr.Error())
	}

	if !errors.Is(ruleErr, baseErr) {
		t.Error("expected Unwrap to return base error")
	}
}

func TestRuleErrorWithoutType(t *testing.T) {
	baseErr := errors.New("base error")

	ruleErr := cortex.NewRuleError("rule-1", "", "evaluate", baseErr)

	expectedMsg := `cortex: rule "rule-1" evaluate: base error`
	if ruleErr.Error() != expectedMsg {
		t.Errorf("expected %q, got %q", expectedMsg, ruleErr.Error())
	}
}

func TestSentinelErrors(t *testing.T) {
	// Just verify they are defined and non-nil
	sentinels := []error{
		cortex.ErrRuleNotFound,
		cortex.ErrLookupNotFound,
		cortex.ErrKeyNotFound,
		cortex.ErrValueNotFound,
		cortex.ErrBuildupNotFound,
		cortex.ErrInvalidRule,
		cortex.ErrInvalidExpression,
		cortex.ErrTypeMismatch,
		cortex.ErrDivisionByZero,
		cortex.ErrAllocationSum,
		cortex.ErrCircularDep,
		cortex.ErrEvaluation,
		cortex.ErrShortCircuit,
		cortex.ErrEngineClosed,
		cortex.ErrTimeout,
		cortex.ErrNilContext,
		cortex.ErrDuplicateRule,
		cortex.ErrDuplicateLookup,
	}

	for _, err := range sentinels {
		if err == nil {
			t.Error("expected non-nil sentinel error")
		}
	}
}
