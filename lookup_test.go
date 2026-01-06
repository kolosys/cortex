package cortex_test

import (
	"context"
	"math"
	"testing"

	"github.com/kolosys/cortex"
)

func TestMapLookup(t *testing.T) {
	lookup := cortex.NewMapLookup("codes", map[string]int{
		"active":   1,
		"inactive": 0,
		"pending":  2,
	})

	if lookup.Name() != "codes" {
		t.Errorf("expected name 'codes', got %q", lookup.Name())
	}

	val, ok := lookup.Get("active")
	if !ok {
		t.Error("expected to find 'active'")
	}
	if val != 1 {
		t.Errorf("expected 1, got %v", val)
	}

	_, ok = lookup.Get("missing")
	if ok {
		t.Error("expected not to find 'missing'")
	}

	// Wrong key type
	_, ok = lookup.Get(123)
	if ok {
		t.Error("expected not to find with wrong key type")
	}
}

func TestRangeLookup(t *testing.T) {
	lookup := cortex.NewRangeLookup("brackets", []cortex.RangeEntry[float64]{
		{Min: 0, Max: 50000, Value: 0.10},
		{Min: 50000, Max: 100000, Value: 0.22},
		{Min: 100000, Max: math.Inf(1), Value: 0.35},
	})

	if lookup.Name() != "brackets" {
		t.Errorf("expected name 'brackets', got %q", lookup.Name())
	}

	tests := []struct {
		input    float64
		expected float64
		found    bool
	}{
		{25000, 0.10, true},
		{50000, 0.22, true}, // boundary - 50000 is in second range
		{75000, 0.22, true},
		{100000, 0.35, true},
		{150000, 0.35, true},
		{-100, 0, false}, // negative value not in any range
	}

	for _, tt := range tests {
		val, ok := lookup.Get(tt.input)
		if ok != tt.found {
			t.Errorf("input %f: expected found=%v, got %v", tt.input, tt.found, ok)
			continue
		}
		if ok && val != tt.expected {
			t.Errorf("input %f: expected %f, got %v", tt.input, tt.expected, val)
		}
	}
}

func TestRangeLookupWrongType(t *testing.T) {
	lookup := cortex.NewRangeLookup("test", []cortex.RangeEntry[float64]{
		{Min: 0, Max: 100, Value: 1.0},
	})

	_, ok := lookup.Get("not a number")
	if ok {
		t.Error("expected not to find with string key")
	}
}

func TestLookupRule(t *testing.T) {
	rule := cortex.MustLookup(cortex.LookupConfig{
		ID:     "get-rate",
		Table:  "rates",
		Key:    "amount",
		Target: "rate",
	})

	if rule.Table() != "rates" {
		t.Errorf("expected table 'rates', got %q", rule.Table())
	}

	evalCtx := cortex.NewEvalContext()
	evalCtx.RegisterLookup(cortex.NewMapLookup("rates", map[float64]float64{
		100.0: 0.05,
		200.0: 0.10,
	}))
	evalCtx.Set("amount", 100.0)

	err := rule.Evaluate(context.Background(), evalCtx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rate, _ := evalCtx.GetFloat64("rate")
	if rate != 0.05 {
		t.Errorf("expected rate=0.05, got %f", rate)
	}
}

func TestLookupRuleDefault(t *testing.T) {
	rule := cortex.MustLookup(cortex.LookupConfig{
		ID:      "get-rate",
		Table:   "rates",
		Key:     "amount",
		Target:  "rate",
		Default: 0.01,
	})

	evalCtx := cortex.NewEvalContext()
	evalCtx.RegisterLookup(cortex.NewMapLookup("rates", map[float64]float64{
		100.0: 0.05,
	}))
	evalCtx.Set("amount", 999.0) // not in lookup

	err := rule.Evaluate(context.Background(), evalCtx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rate, _ := evalCtx.GetFloat64("rate")
	if rate != 0.01 {
		t.Errorf("expected rate=0.01 (default), got %f", rate)
	}
}

func TestLookupRuleRequired(t *testing.T) {
	rule := cortex.MustLookup(cortex.LookupConfig{
		ID:       "get-rate",
		Table:    "rates",
		Key:      "amount",
		Target:   "rate",
		Required: true,
	})

	evalCtx := cortex.NewEvalContext()
	evalCtx.RegisterLookup(cortex.NewMapLookup("rates", map[float64]float64{
		100.0: 0.05,
	}))
	evalCtx.Set("amount", 999.0) // not in lookup

	err := rule.Evaluate(context.Background(), evalCtx)
	if err == nil {
		t.Error("expected error for required lookup not found")
	}
}

func TestLookupRuleMissingTable(t *testing.T) {
	rule := cortex.MustLookup(cortex.LookupConfig{
		ID:     "get-rate",
		Table:  "missing_table",
		Key:    "amount",
		Target: "rate",
	})

	evalCtx := cortex.NewEvalContext()
	evalCtx.Set("amount", 100.0)

	err := rule.Evaluate(context.Background(), evalCtx)
	if err == nil {
		t.Error("expected error for missing table")
	}
}

func TestLookupRuleMissingKey(t *testing.T) {
	rule := cortex.MustLookup(cortex.LookupConfig{
		ID:     "get-rate",
		Table:  "rates",
		Key:    "amount",
		Target: "rate",
	})

	evalCtx := cortex.NewEvalContext()
	evalCtx.RegisterLookup(cortex.NewMapLookup("rates", map[float64]float64{}))
	// Don't set amount

	err := rule.Evaluate(context.Background(), evalCtx)
	if err == nil {
		t.Error("expected error for missing key")
	}
}

func TestTaxBracketLookup(t *testing.T) {
	lookup := cortex.NewTaxBracketLookup("tax", []cortex.TaxBracket{
		{Min: 0, Max: 50000, Rate: 0.10},
		{Min: 50000, Max: 100000, Rate: 0.22},
		{Min: 100000, Max: 0, Rate: 0.35}, // 0 = infinity
	})

	val, ok := lookup.Get(75000.0)
	if !ok || val != 0.22 {
		t.Errorf("expected 0.22, got %v (found=%v)", val, ok)
	}

	val, ok = lookup.Get(150000.0)
	if !ok || val != 0.35 {
		t.Errorf("expected 0.35, got %v (found=%v)", val, ok)
	}
}

func TestLookupValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  cortex.LookupConfig
		wantErr bool
	}{
		{"missing ID", cortex.LookupConfig{Table: "t", Key: "k", Target: "t"}, true},
		{"missing table", cortex.LookupConfig{ID: "id", Key: "k", Target: "t"}, true},
		{"missing key", cortex.LookupConfig{ID: "id", Table: "t", Target: "t"}, true},
		{"missing target", cortex.LookupConfig{ID: "id", Table: "t", Key: "k"}, true},
		{"valid", cortex.LookupConfig{ID: "id", Table: "t", Key: "k", Target: "t"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cortex.NewLookup(tt.config)
			if tt.wantErr && err == nil {
				t.Error("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
