package cortex_test

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/kolosys/cortex"
)

func TestEngineBasic(t *testing.T) {
	engine := cortex.New("test", cortex.DefaultConfig())

	if engine.Name() != "test" {
		t.Errorf("expected name 'test', got %q", engine.Name())
	}
	if engine.Rules() != 0 {
		t.Errorf("expected 0 rules, got %d", engine.Rules())
	}
}

func TestEngineAddRule(t *testing.T) {
	engine := cortex.New("test", cortex.DefaultConfig())

	err := engine.AddRule(cortex.MustAssignment(cortex.AssignmentConfig{
		ID:     "test-rule",
		Target: "value",
		Value:  42,
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if engine.Rules() != 1 {
		t.Errorf("expected 1 rule, got %d", engine.Rules())
	}
}

func TestEngineDuplicateRule(t *testing.T) {
	engine := cortex.New("test", cortex.DefaultConfig())

	rule := cortex.MustAssignment(cortex.AssignmentConfig{
		ID:     "dupe",
		Target: "value",
		Value:  42,
	})

	if err := engine.AddRule(rule); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err := engine.AddRule(rule)
	if !errors.Is(err, cortex.ErrDuplicateRule) {
		t.Errorf("expected ErrDuplicateRule, got %v", err)
	}
}

func TestEngineEvaluate(t *testing.T) {
	engine := cortex.New("test", cortex.DefaultConfig())

	engine.AddRules(
		cortex.MustAssignment(cortex.AssignmentConfig{
			ID:     "set-x",
			Target: "x",
			Value:  10.0,
		}),
		cortex.MustAssignment(cortex.AssignmentConfig{
			ID:     "set-y",
			Target: "y",
			Value:  20.0,
		}),
		cortex.MustFormula(cortex.FormulaConfig{
			ID:         "calc-sum",
			Target:     "sum",
			Expression: "x + y",
		}),
	)

	ctx := context.Background()
	evalCtx := cortex.NewEvalContext()

	result, err := engine.Evaluate(ctx, evalCtx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.Success {
		t.Error("expected success")
	}
	if result.RulesEvaluated != 3 {
		t.Errorf("expected 3 rules evaluated, got %d", result.RulesEvaluated)
	}

	sum, err := evalCtx.GetFloat64("sum")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sum != 30.0 {
		t.Errorf("expected sum=30, got %f", sum)
	}
}

func TestEngineTimeout(t *testing.T) {
	config := cortex.DefaultConfig()
	config.Timeout = 1 * time.Millisecond

	engine := cortex.New("test", config)

	// Add a rule that uses a slow formula
	engine.AddRule(cortex.MustFormula(cortex.FormulaConfig{
		ID:     "slow",
		Target: "result",
		Formula: func(ctx context.Context, evalCtx *cortex.EvalContext) (any, error) {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(100 * time.Millisecond):
				return 42, nil
			}
		},
	}))

	ctx := context.Background()
	evalCtx := cortex.NewEvalContext()

	_, err := engine.Evaluate(ctx, evalCtx)
	// Either ErrTimeout or context.DeadlineExceeded (wrapped in RuleError)
	if !errors.Is(err, cortex.ErrTimeout) && !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected timeout error, got %v", err)
	}
}

func TestEngineFailFast(t *testing.T) {
	config := cortex.DefaultConfig()
	config.Mode = cortex.ModeFailFast

	engine := cortex.New("test", config)

	engine.AddRules(
		cortex.MustFormula(cortex.FormulaConfig{
			ID:     "fail",
			Target: "x",
			Formula: func(ctx context.Context, evalCtx *cortex.EvalContext) (any, error) {
				return nil, errors.New("intentional error")
			},
		}),
		cortex.MustAssignment(cortex.AssignmentConfig{
			ID:     "after-fail",
			Target: "y",
			Value:  42,
		}),
	)

	ctx := context.Background()
	evalCtx := cortex.NewEvalContext()

	result, err := engine.Evaluate(ctx, evalCtx)
	if err == nil {
		t.Fatal("expected error")
	}

	if result.RulesEvaluated != 0 {
		t.Errorf("expected 0 rules evaluated in fail-fast mode, got %d", result.RulesEvaluated)
	}

	if evalCtx.Has("y") {
		t.Error("rule after failure should not have executed")
	}
}

func TestEngineCollectAll(t *testing.T) {
	config := cortex.DefaultConfig()
	config.Mode = cortex.ModeCollectAll

	engine := cortex.New("test", config)

	engine.AddRules(
		cortex.MustFormula(cortex.FormulaConfig{
			ID:     "fail1",
			Target: "x",
			Formula: func(ctx context.Context, evalCtx *cortex.EvalContext) (any, error) {
				return nil, errors.New("error 1")
			},
		}),
		cortex.MustFormula(cortex.FormulaConfig{
			ID:     "fail2",
			Target: "y",
			Formula: func(ctx context.Context, evalCtx *cortex.EvalContext) (any, error) {
				return nil, errors.New("error 2")
			},
		}),
		cortex.MustAssignment(cortex.AssignmentConfig{
			ID:     "success",
			Target: "z",
			Value:  42,
		}),
	)

	ctx := context.Background()
	evalCtx := cortex.NewEvalContext()

	result, _ := engine.Evaluate(ctx, evalCtx)

	if len(result.Errors) != 2 {
		t.Errorf("expected 2 errors, got %d", len(result.Errors))
	}

	if result.RulesEvaluated != 1 {
		t.Errorf("expected 1 successful rule, got %d", result.RulesEvaluated)
	}

	z, _ := evalCtx.GetFloat64("z")
	if z != 42 {
		t.Errorf("expected z=42, got %f", z)
	}
}

func TestEngineLookup(t *testing.T) {
	engine := cortex.New("test", cortex.DefaultConfig())

	engine.RegisterLookup(cortex.NewRangeLookup("tax_brackets", []cortex.RangeEntry[float64]{
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
			ID:     "get-rate",
			Table:  "tax_brackets",
			Key:    "salary",
			Target: "tax_rate",
		}),
		cortex.MustFormula(cortex.FormulaConfig{
			ID:         "calc-tax",
			Target:     "tax",
			Expression: "salary * tax_rate",
		}),
	)

	ctx := context.Background()
	evalCtx := cortex.NewEvalContext()

	_, err := engine.Evaluate(ctx, evalCtx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rate, _ := evalCtx.GetFloat64("tax_rate")
	if rate != 0.22 {
		t.Errorf("expected rate=0.22, got %f", rate)
	}

	tax, _ := evalCtx.GetFloat64("tax")
	if tax != 16500.0 {
		t.Errorf("expected tax=16500, got %f", tax)
	}
}

func TestEngineAllocation(t *testing.T) {
	engine := cortex.New("test", cortex.DefaultConfig())

	engine.AddRules(
		cortex.MustAssignment(cortex.AssignmentConfig{
			ID:     "budget",
			Target: "total_budget",
			Value:  100000.0,
		}),
		cortex.MustAllocation(cortex.AllocationConfig{
			ID:       "split",
			Source:   "total_budget",
			Strategy: cortex.StrategyPercentage,
			Targets: []cortex.AllocationTarget{
				{Key: "eng", Amount: 50},
				{Key: "ops", Amount: 30},
				{Key: "admin", Amount: 20},
			},
		}),
	)

	ctx := context.Background()
	evalCtx := cortex.NewEvalContext()

	_, err := engine.Evaluate(ctx, evalCtx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	eng, _ := evalCtx.GetFloat64("eng")
	ops, _ := evalCtx.GetFloat64("ops")
	admin, _ := evalCtx.GetFloat64("admin")

	if eng != 50000.0 {
		t.Errorf("expected eng=50000, got %f", eng)
	}
	if ops != 30000.0 {
		t.Errorf("expected ops=30000, got %f", ops)
	}
	if admin != 20000.0 {
		t.Errorf("expected admin=20000, got %f", admin)
	}
}

func TestEngineBuildup(t *testing.T) {
	engine := cortex.New("test", cortex.DefaultConfig())

	engine.AddRules(
		cortex.MustAssignment(cortex.AssignmentConfig{
			ID:     "v1",
			Target: "value",
			Value:  10.0,
		}),
		cortex.MustBuildup(cortex.BuildupConfig{
			ID:        "add1",
			Buildup:   "total",
			Operation: cortex.BuildupSum,
			Source:    "value",
		}),
		cortex.MustAssignment(cortex.AssignmentConfig{
			ID:     "v2",
			Target: "value",
			Value:  20.0,
		}),
		cortex.MustBuildup(cortex.BuildupConfig{
			ID:        "add2",
			Buildup:   "total",
			Operation: cortex.BuildupSum,
			Source:    "value",
			Target:    "running_total",
		}),
	)

	ctx := context.Background()
	evalCtx := cortex.NewEvalContext()

	_, err := engine.Evaluate(ctx, evalCtx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	total, _ := evalCtx.GetFloat64("running_total")
	if total != 30.0 {
		t.Errorf("expected total=30, got %f", total)
	}
}

func TestEngineClose(t *testing.T) {
	engine := cortex.New("test", cortex.DefaultConfig())
	engine.Close()

	err := engine.AddRule(cortex.MustAssignment(cortex.AssignmentConfig{
		ID:     "test",
		Target: "x",
		Value:  1,
	}))

	if !errors.Is(err, cortex.ErrEngineClosed) {
		t.Errorf("expected ErrEngineClosed, got %v", err)
	}
}

func TestEngineNilContext(t *testing.T) {
	engine := cortex.New("test", cortex.DefaultConfig())

	_, err := engine.Evaluate(context.Background(), nil)
	if !errors.Is(err, cortex.ErrNilContext) {
		t.Errorf("expected ErrNilContext, got %v", err)
	}
}
