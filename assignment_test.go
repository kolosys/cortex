package cortex_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kolosys/cortex"
)

func TestAssignmentStatic(t *testing.T) {
	rule := cortex.MustAssignment(cortex.AssignmentConfig{
		ID:     "set-value",
		Target: "x",
		Value:  42,
	})

	if rule.Target() != "x" {
		t.Errorf("expected target 'x', got %q", rule.Target())
	}

	evalCtx := cortex.NewEvalContext()
	err := rule.Evaluate(context.Background(), evalCtx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	val, ok := evalCtx.Get("x")
	if !ok {
		t.Error("expected to find 'x'")
	}
	if val != 42 {
		t.Errorf("expected 42, got %v", val)
	}
}

func TestAssignmentDynamic(t *testing.T) {
	rule := cortex.MustAssignment(cortex.AssignmentConfig{
		ID:     "calc-value",
		Target: "result",
		ValueFunc: func(ctx context.Context, evalCtx *cortex.EvalContext) (any, error) {
			x, _ := evalCtx.GetFloat64("x")
			return x * 2, nil
		},
	})

	evalCtx := cortex.NewEvalContext()
	evalCtx.Set("x", 10.0)

	err := rule.Evaluate(context.Background(), evalCtx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, _ := evalCtx.GetFloat64("result")
	if result != 20.0 {
		t.Errorf("expected 20, got %f", result)
	}
}

func TestAssignmentDynamicError(t *testing.T) {
	rule := cortex.MustAssignment(cortex.AssignmentConfig{
		ID:     "calc-value",
		Target: "result",
		ValueFunc: func(ctx context.Context, evalCtx *cortex.EvalContext) (any, error) {
			return nil, errors.New("intentional error")
		},
	})

	evalCtx := cortex.NewEvalContext()
	err := rule.Evaluate(context.Background(), evalCtx)
	if err == nil {
		t.Error("expected error")
	}
}

func TestAssignmentValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  cortex.AssignmentConfig
		wantErr bool
	}{
		{"missing ID", cortex.AssignmentConfig{Target: "x", Value: 1}, true},
		{"missing target", cortex.AssignmentConfig{ID: "id", Value: 1}, true},
		{"missing value", cortex.AssignmentConfig{ID: "id", Target: "x"}, true},
		{"valid static", cortex.AssignmentConfig{ID: "id", Target: "x", Value: 1}, false},
		{"valid dynamic", cortex.AssignmentConfig{
			ID:     "id",
			Target: "x",
			ValueFunc: func(ctx context.Context, evalCtx *cortex.EvalContext) (any, error) {
				return 1, nil
			},
		}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cortex.NewAssignment(tt.config)
			if tt.wantErr && err == nil {
				t.Error("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestAssignmentMustPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()

	cortex.MustAssignment(cortex.AssignmentConfig{}) // invalid config
}

func TestAssignmentMetadata(t *testing.T) {
	rule := cortex.MustAssignment(cortex.AssignmentConfig{
		ID:          "set-value",
		Name:        "Set Value",
		Description: "Sets the value to 42",
		Target:      "x",
		Value:       42,
		Deps:        []string{"dep1", "dep2"},
	})

	if rule.ID() != "set-value" {
		t.Errorf("expected ID 'set-value', got %q", rule.ID())
	}
	if rule.Name() != "Set Value" {
		t.Errorf("expected name 'Set Value', got %q", rule.Name())
	}
	if rule.Description() != "Sets the value to 42" {
		t.Errorf("expected description, got %q", rule.Description())
	}
	deps := rule.Dependencies()
	if len(deps) != 2 {
		t.Errorf("expected 2 dependencies, got %d", len(deps))
	}
}
