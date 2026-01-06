package cortex_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kolosys/cortex"
)

func TestAllocationPercentage(t *testing.T) {
	rule := cortex.MustAllocation(cortex.AllocationConfig{
		ID:       "alloc",
		Source:   "total",
		Strategy: cortex.StrategyPercentage,
		Targets: []cortex.AllocationTarget{
			{Key: "a", Amount: 50},
			{Key: "b", Amount: 30},
			{Key: "c", Amount: 20},
		},
	})

	evalCtx := cortex.NewEvalContext()
	evalCtx.Set("total", 1000.0)

	err := rule.Evaluate(context.Background(), evalCtx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	a, _ := evalCtx.GetFloat64("a")
	b, _ := evalCtx.GetFloat64("b")
	c, _ := evalCtx.GetFloat64("c")

	if a != 500.0 {
		t.Errorf("expected a=500, got %f", a)
	}
	if b != 300.0 {
		t.Errorf("expected b=300, got %f", b)
	}
	if c != 200.0 {
		t.Errorf("expected c=200, got %f", c)
	}
}

func TestAllocationFixed(t *testing.T) {
	rule := cortex.MustAllocation(cortex.AllocationConfig{
		ID:       "alloc",
		Source:   "total",
		Strategy: cortex.StrategyFixed,
		Targets: []cortex.AllocationTarget{
			{Key: "a", Amount: 100},
			{Key: "b", Amount: 200},
		},
		Remainder: "leftover",
	})

	evalCtx := cortex.NewEvalContext()
	evalCtx.Set("total", 1000.0)

	err := rule.Evaluate(context.Background(), evalCtx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	a, _ := evalCtx.GetFloat64("a")
	b, _ := evalCtx.GetFloat64("b")
	leftover, _ := evalCtx.GetFloat64("leftover")

	if a != 100.0 {
		t.Errorf("expected a=100, got %f", a)
	}
	if b != 200.0 {
		t.Errorf("expected b=200, got %f", b)
	}
	if leftover != 700.0 {
		t.Errorf("expected leftover=700, got %f", leftover)
	}
}

func TestAllocationWeighted(t *testing.T) {
	rule := cortex.MustAllocation(cortex.AllocationConfig{
		ID:       "alloc",
		Source:   "total",
		Strategy: cortex.StrategyWeighted,
		Targets: []cortex.AllocationTarget{
			{Key: "a", Amount: 2},
			{Key: "b", Amount: 3},
		},
	})

	evalCtx := cortex.NewEvalContext()
	evalCtx.Set("total", 1000.0)

	err := rule.Evaluate(context.Background(), evalCtx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	a, _ := evalCtx.GetFloat64("a")
	b, _ := evalCtx.GetFloat64("b")

	if a != 400.0 {
		t.Errorf("expected a=400, got %f", a)
	}
	if b != 600.0 {
		t.Errorf("expected b=600, got %f", b)
	}
}

func TestAllocationEqual(t *testing.T) {
	rule := cortex.MustAllocation(cortex.AllocationConfig{
		ID:       "alloc",
		Source:   "total",
		Strategy: cortex.StrategyEqual,
		Targets: []cortex.AllocationTarget{
			{Key: "a"},
			{Key: "b"},
			{Key: "c"},
		},
	})

	evalCtx := cortex.NewEvalContext()
	evalCtx.Set("total", 300.0)

	err := rule.Evaluate(context.Background(), evalCtx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	a, _ := evalCtx.GetFloat64("a")
	b, _ := evalCtx.GetFloat64("b")
	c, _ := evalCtx.GetFloat64("c")

	if a != 100.0 {
		t.Errorf("expected a=100, got %f", a)
	}
	if b != 100.0 {
		t.Errorf("expected b=100, got %f", b)
	}
	if c != 100.0 {
		t.Errorf("expected c=100, got %f", c)
	}
}

func TestAllocationRatio(t *testing.T) {
	rule := cortex.MustAllocation(cortex.AllocationConfig{
		ID:       "alloc",
		Source:   "total",
		Strategy: cortex.StrategyRatio,
		Targets: []cortex.AllocationTarget{
			{Key: "a", Amount: 1},
			{Key: "b", Amount: 2},
			{Key: "c", Amount: 2},
		},
	})

	evalCtx := cortex.NewEvalContext()
	evalCtx.Set("total", 500.0)

	err := rule.Evaluate(context.Background(), evalCtx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	a, _ := evalCtx.GetFloat64("a")
	b, _ := evalCtx.GetFloat64("b")
	c, _ := evalCtx.GetFloat64("c")

	if a != 100.0 {
		t.Errorf("expected a=100, got %f", a)
	}
	if b != 200.0 {
		t.Errorf("expected b=200, got %f", b)
	}
	if c != 200.0 {
		t.Errorf("expected c=200, got %f", c)
	}
}

func TestAllocationInvalidPercentage(t *testing.T) {
	_, err := cortex.NewAllocation(cortex.AllocationConfig{
		ID:       "alloc",
		Source:   "total",
		Strategy: cortex.StrategyPercentage,
		Targets: []cortex.AllocationTarget{
			{Key: "a", Amount: 50},
			{Key: "b", Amount: 30},
		},
	})

	if !errors.Is(err, cortex.ErrAllocationSum) {
		t.Errorf("expected ErrAllocationSum, got %v", err)
	}
}

func TestAllocationMissingSource(t *testing.T) {
	rule := cortex.MustAllocation(cortex.AllocationConfig{
		ID:       "alloc",
		Source:   "total",
		Strategy: cortex.StrategyPercentage,
		Targets: []cortex.AllocationTarget{
			{Key: "a", Amount: 100},
		},
	})

	evalCtx := cortex.NewEvalContext()
	// Don't set total

	err := rule.Evaluate(context.Background(), evalCtx)
	if err == nil {
		t.Error("expected error for missing source")
	}
}

func TestParseAllocationStrategy(t *testing.T) {
	tests := []struct {
		input    string
		expected cortex.AllocationStrategy
		wantErr  bool
	}{
		{"percentage", cortex.StrategyPercentage, false},
		{"fixed", cortex.StrategyFixed, false},
		{"weighted", cortex.StrategyWeighted, false},
		{"equal", cortex.StrategyEqual, false},
		{"ratio", cortex.StrategyRatio, false},
		{"invalid", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := cortex.ParseAllocationStrategy(tt.input)
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
