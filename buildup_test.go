package cortex_test

import (
	"context"
	"testing"

	"github.com/kolosys/cortex"
)

func TestBuildupSum(t *testing.T) {
	evalCtx := cortex.NewEvalContext()
	buildup := evalCtx.GetOrCreateBuildup("total", cortex.BuildupSum, 0)

	buildup.Add(10)
	buildup.Add(20)
	buildup.Add(30)

	if buildup.Current() != 60 {
		t.Errorf("expected 60, got %f", buildup.Current())
	}
	if buildup.Count() != 3 {
		t.Errorf("expected count=3, got %d", buildup.Count())
	}
}

func TestBuildupMin(t *testing.T) {
	evalCtx := cortex.NewEvalContext()
	buildup := evalCtx.GetOrCreateBuildup("min", cortex.BuildupMin, 0)

	buildup.Add(30)
	buildup.Add(10)
	buildup.Add(20)

	if buildup.Current() != 10 {
		t.Errorf("expected 10, got %f", buildup.Current())
	}
}

func TestBuildupMax(t *testing.T) {
	evalCtx := cortex.NewEvalContext()
	buildup := evalCtx.GetOrCreateBuildup("max", cortex.BuildupMax, 0)

	buildup.Add(10)
	buildup.Add(30)
	buildup.Add(20)

	if buildup.Current() != 30 {
		t.Errorf("expected 30, got %f", buildup.Current())
	}
}

func TestBuildupAvg(t *testing.T) {
	evalCtx := cortex.NewEvalContext()
	buildup := evalCtx.GetOrCreateBuildup("avg", cortex.BuildupAvg, 0)

	buildup.Add(10)
	buildup.Add(20)
	buildup.Add(30)

	if buildup.Current() != 20 {
		t.Errorf("expected 20, got %f", buildup.Current())
	}
}

func TestBuildupCount(t *testing.T) {
	evalCtx := cortex.NewEvalContext()
	buildup := evalCtx.GetOrCreateBuildup("count", cortex.BuildupCount, 0)

	buildup.Add(100) // value doesn't matter for count
	buildup.Add(200)
	buildup.Add(300)

	if buildup.Current() != 3 {
		t.Errorf("expected 3, got %f", buildup.Current())
	}
}

func TestBuildupProduct(t *testing.T) {
	evalCtx := cortex.NewEvalContext()
	buildup := evalCtx.GetOrCreateBuildup("product", cortex.BuildupProduct, 1)

	buildup.Add(2)
	buildup.Add(3)
	buildup.Add(4)

	if buildup.Current() != 24 {
		t.Errorf("expected 24, got %f", buildup.Current())
	}
}

func TestBuildupReset(t *testing.T) {
	evalCtx := cortex.NewEvalContext()
	buildup := evalCtx.GetOrCreateBuildup("total", cortex.BuildupSum, 0)

	buildup.Add(10)
	buildup.Add(20)
	buildup.Reset(100)

	if buildup.Current() != 100 {
		t.Errorf("expected 100 after reset, got %f", buildup.Current())
	}
	if buildup.Count() != 0 {
		t.Errorf("expected count=0 after reset, got %d", buildup.Count())
	}
}

func TestBuildupRule(t *testing.T) {
	rule := cortex.MustBuildup(cortex.BuildupConfig{
		ID:        "add",
		Buildup:   "total",
		Operation: cortex.BuildupSum,
		Source:    "value",
		Target:    "running_total",
	})

	evalCtx := cortex.NewEvalContext()
	ctx := context.Background()

	// First value
	evalCtx.Set("value", 10.0)
	if err := rule.Evaluate(ctx, evalCtx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	total, _ := evalCtx.GetFloat64("running_total")
	if total != 10 {
		t.Errorf("expected 10, got %f", total)
	}

	// Second value
	evalCtx.Set("value", 20.0)
	if err := rule.Evaluate(ctx, evalCtx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	total, _ = evalCtx.GetFloat64("running_total")
	if total != 30 {
		t.Errorf("expected 30, got %f", total)
	}
}

func TestBuildupCountRule(t *testing.T) {
	rule := cortex.MustBuildup(cortex.BuildupConfig{
		ID:        "count",
		Buildup:   "item_count",
		Operation: cortex.BuildupCount,
		Target:    "count",
	})

	evalCtx := cortex.NewEvalContext()
	ctx := context.Background()

	for range 5 {
		if err := rule.Evaluate(ctx, evalCtx); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	count, _ := evalCtx.GetFloat64("count")
	if count != 5 {
		t.Errorf("expected 5, got %f", count)
	}
}

func TestParseBuildupOperation(t *testing.T) {
	tests := []struct {
		input    string
		expected cortex.BuildupOperation
		wantErr  bool
	}{
		{"sum", cortex.BuildupSum, false},
		{"min", cortex.BuildupMin, false},
		{"max", cortex.BuildupMax, false},
		{"avg", cortex.BuildupAvg, false},
		{"average", cortex.BuildupAvg, false},
		{"count", cortex.BuildupCount, false},
		{"product", cortex.BuildupProduct, false},
		{"invalid", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := cortex.ParseBuildupOperation(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}
