package cortex_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kolosys/cortex"
)

func TestFormulaWithFunction(t *testing.T) {
	rule := cortex.MustFormula(cortex.FormulaConfig{
		ID:     "calc",
		Target: "result",
		Formula: func(ctx context.Context, evalCtx *cortex.EvalContext) (any, error) {
			x, _ := evalCtx.GetFloat64("x")
			y, _ := evalCtx.GetFloat64("y")
			return x + y, nil
		},
	})

	evalCtx := cortex.NewEvalContext()
	evalCtx.Set("x", 10.0)
	evalCtx.Set("y", 20.0)

	err := rule.Evaluate(context.Background(), evalCtx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, _ := evalCtx.GetFloat64("result")
	if result != 30.0 {
		t.Errorf("expected 30, got %f", result)
	}
}

func TestFormulaWithExpression(t *testing.T) {
	rule := cortex.MustFormula(cortex.FormulaConfig{
		ID:         "calc",
		Target:     "result",
		Expression: "x * y + 5",
	})

	if rule.Target() != "result" {
		t.Errorf("expected target 'result', got %q", rule.Target())
	}
	if rule.Expression() != "x * y + 5" {
		t.Errorf("expected expression 'x * y + 5', got %q", rule.Expression())
	}

	evalCtx := cortex.NewEvalContext()
	evalCtx.Set("x", 3.0)
	evalCtx.Set("y", 4.0)

	err := rule.Evaluate(context.Background(), evalCtx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, _ := evalCtx.GetFloat64("result")
	if result != 17.0 {
		t.Errorf("expected 17, got %f", result)
	}
}

func TestFormulaFunctionError(t *testing.T) {
	rule := cortex.MustFormula(cortex.FormulaConfig{
		ID:     "calc",
		Target: "result",
		Formula: func(ctx context.Context, evalCtx *cortex.EvalContext) (any, error) {
			return nil, errors.New("intentional error")
		},
	})

	evalCtx := cortex.NewEvalContext()
	err := rule.Evaluate(context.Background(), evalCtx)
	if err == nil {
		t.Error("expected error")
	}
}

func TestFormulaInvalidExpression(t *testing.T) {
	_, err := cortex.NewFormula(cortex.FormulaConfig{
		ID:         "calc",
		Target:     "result",
		Expression: "invalid ( expression",
	})

	if !errors.Is(err, cortex.ErrInvalidExpression) {
		t.Errorf("expected ErrInvalidExpression, got %v", err)
	}
}

func TestFormulaValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  cortex.FormulaConfig
		wantErr bool
	}{
		{"missing ID", cortex.FormulaConfig{Target: "t", Expression: "1"}, true},
		{"missing target", cortex.FormulaConfig{ID: "id", Expression: "1"}, true},
		{"missing formula and expression", cortex.FormulaConfig{ID: "id", Target: "t"}, true},
		{"valid with expression", cortex.FormulaConfig{ID: "id", Target: "t", Expression: "1 + 2"}, false},
		{"valid with function", cortex.FormulaConfig{
			ID:      "id",
			Target:  "t",
			Formula: func(ctx context.Context, evalCtx *cortex.EvalContext) (any, error) { return 1, nil },
		}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cortex.NewFormula(tt.config)
			if tt.wantErr && err == nil {
				t.Error("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestFormulaInputs(t *testing.T) {
	rule := cortex.MustFormula(cortex.FormulaConfig{
		ID:         "calc",
		Target:     "result",
		Expression: "x + y",
		Inputs:     []string{"x", "y"},
	})

	inputs := rule.Inputs()
	if len(inputs) != 2 {
		t.Errorf("expected 2 inputs, got %d", len(inputs))
	}
}

func TestFormulaHelperAdd(t *testing.T) {
	fn := cortex.Add("a", "b")

	evalCtx := cortex.NewEvalContext()
	evalCtx.Set("a", 10.0)
	evalCtx.Set("b", 20.0)

	result, err := fn(context.Background(), evalCtx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 30.0 {
		t.Errorf("expected 30, got %v", result)
	}
}

func TestFormulaHelperSubtract(t *testing.T) {
	fn := cortex.Subtract("a", "b")

	evalCtx := cortex.NewEvalContext()
	evalCtx.Set("a", 30.0)
	evalCtx.Set("b", 10.0)

	result, err := fn(context.Background(), evalCtx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 20.0 {
		t.Errorf("expected 20, got %v", result)
	}
}

func TestFormulaHelperMultiply(t *testing.T) {
	fn := cortex.Multiply("a", "b")

	evalCtx := cortex.NewEvalContext()
	evalCtx.Set("a", 5.0)
	evalCtx.Set("b", 4.0)

	result, err := fn(context.Background(), evalCtx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 20.0 {
		t.Errorf("expected 20, got %v", result)
	}
}

func TestFormulaHelperDivide(t *testing.T) {
	fn := cortex.Divide("a", "b")

	evalCtx := cortex.NewEvalContext()
	evalCtx.Set("a", 20.0)
	evalCtx.Set("b", 4.0)

	result, err := fn(context.Background(), evalCtx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 5.0 {
		t.Errorf("expected 5, got %v", result)
	}
}

func TestFormulaHelperDivideByZero(t *testing.T) {
	fn := cortex.Divide("a", "b")

	evalCtx := cortex.NewEvalContext()
	evalCtx.Set("a", 20.0)
	evalCtx.Set("b", 0.0)

	_, err := fn(context.Background(), evalCtx)
	if !errors.Is(err, cortex.ErrDivisionByZero) {
		t.Errorf("expected ErrDivisionByZero, got %v", err)
	}
}

func TestFormulaHelperPercentage(t *testing.T) {
	fn := cortex.Percentage("value", 15)

	evalCtx := cortex.NewEvalContext()
	evalCtx.Set("value", 200.0)

	result, err := fn(context.Background(), evalCtx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 30.0 {
		t.Errorf("expected 30 (15%% of 200), got %v", result)
	}
}

func TestFormulaHelperConditional(t *testing.T) {
	fn := cortex.Conditional("is_senior", 0.15, 0.10)

	evalCtx := cortex.NewEvalContext()

	evalCtx.Set("is_senior", true)
	result, _ := fn(context.Background(), evalCtx)
	if result != 0.15 {
		t.Errorf("expected 0.15 for senior, got %v", result)
	}

	evalCtx.Set("is_senior", false)
	result, _ = fn(context.Background(), evalCtx)
	if result != 0.10 {
		t.Errorf("expected 0.10 for non-senior, got %v", result)
	}
}
