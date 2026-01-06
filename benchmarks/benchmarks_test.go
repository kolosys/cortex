package benchmarks_test

import (
	"context"
	"math"
	"testing"

	"github.com/kolosys/cortex"
)

func BenchmarkEngineEvaluate(b *testing.B) {
	engine := cortex.New("bench", cortex.DefaultConfig())

	engine.AddRules(
		cortex.MustAssignment(cortex.AssignmentConfig{
			ID:     "x",
			Target: "x",
			Value:  10.0,
		}),
		cortex.MustAssignment(cortex.AssignmentConfig{
			ID:     "y",
			Target: "y",
			Value:  20.0,
		}),
		cortex.MustFormula(cortex.FormulaConfig{
			ID:         "sum",
			Target:     "sum",
			Expression: "x + y",
		}),
	)

	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		evalCtx := cortex.NewEvalContext()
		engine.Evaluate(ctx, evalCtx)
	}
}

func BenchmarkEngineWithLookup(b *testing.B) {
	engine := cortex.New("bench", cortex.DefaultConfig())

	engine.RegisterLookup(cortex.NewRangeLookup("rates", []cortex.RangeEntry[float64]{
		{Min: 0, Max: 50000, Value: 0.10},
		{Min: 50000, Max: 100000, Value: 0.22},
		{Min: 100000, Max: math.Inf(1), Value: 0.35},
	}))

	engine.AddRules(
		cortex.MustAssignment(cortex.AssignmentConfig{
			ID:     "salary",
			Target: "salary",
			Value:  75000.0,
		}),
		cortex.MustLookup(cortex.LookupConfig{
			ID:     "rate",
			Table:  "rates",
			Key:    "salary",
			Target: "rate",
		}),
		cortex.MustFormula(cortex.FormulaConfig{
			ID:         "tax",
			Target:     "tax",
			Expression: "salary * rate",
		}),
	)

	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		evalCtx := cortex.NewEvalContext()
		engine.Evaluate(ctx, evalCtx)
	}
}

func BenchmarkExpressionEvaluate(b *testing.B) {
	engine := cortex.New("bench", cortex.DefaultConfig())

	engine.AddRules(
		cortex.MustAssignment(cortex.AssignmentConfig{
			ID:     "a",
			Target: "a",
			Value:  10.0,
		}),
		cortex.MustAssignment(cortex.AssignmentConfig{
			ID:     "b",
			Target: "b",
			Value:  20.0,
		}),
		cortex.MustAssignment(cortex.AssignmentConfig{
			ID:     "c",
			Target: "c",
			Value:  30.0,
		}),
		cortex.MustFormula(cortex.FormulaConfig{
			ID:         "result",
			Target:     "result",
			Expression: "(a + b) * c / 2 + min(a, b)",
		}),
	)

	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		evalCtx := cortex.NewEvalContext()
		engine.Evaluate(ctx, evalCtx)
	}
}

func BenchmarkAllocation(b *testing.B) {
	engine := cortex.New("bench", cortex.DefaultConfig())

	engine.AddRules(
		cortex.MustAssignment(cortex.AssignmentConfig{
			ID:     "total",
			Target: "total",
			Value:  1000000.0,
		}),
		cortex.MustAllocation(cortex.AllocationConfig{
			ID:       "split",
			Source:   "total",
			Strategy: cortex.StrategyPercentage,
			Targets: []cortex.AllocationTarget{
				{Key: "a", Amount: 25},
				{Key: "b", Amount: 25},
				{Key: "c", Amount: 25},
				{Key: "d", Amount: 25},
			},
		}),
	)

	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		evalCtx := cortex.NewEvalContext()
		engine.Evaluate(ctx, evalCtx)
	}
}

func BenchmarkBuildup(b *testing.B) {
	engine := cortex.New("bench", cortex.DefaultConfig())

	for i := range 10 {
		engine.AddRule(cortex.MustAssignment(cortex.AssignmentConfig{
			ID:     "v" + string(rune('0'+i)),
			Target: "value",
			Value:  float64(i * 100),
		}))
		engine.AddRule(cortex.MustBuildup(cortex.BuildupConfig{
			ID:        "add" + string(rune('0'+i)),
			Buildup:   "total",
			Operation: cortex.BuildupSum,
			Source:    "value",
		}))
	}

	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		evalCtx := cortex.NewEvalContext()
		engine.Evaluate(ctx, evalCtx)
	}
}

func BenchmarkEvalContextSet(b *testing.B) {
	evalCtx := cortex.NewEvalContext()

	b.ResetTimer()
	for b.Loop() {
		evalCtx.Set("key", 42.0)
	}
}

func BenchmarkEvalContextGet(b *testing.B) {
	evalCtx := cortex.NewEvalContext()
	evalCtx.Set("key", 42.0)

	b.ResetTimer()
	for b.Loop() {
		evalCtx.Get("key")
	}
}
