package parse_test

import (
	"context"
	"testing"

	"github.com/kolosys/cortex"
	"github.com/kolosys/cortex/parse"
)

const testRuleSetJSON = `{
	"version": "1.0",
	"name": "test-rules",
	"lookups": [
		{
			"name": "tax_rates",
			"type": "range",
			"entries": [
				{"min": 0, "max": 50000, "value": 0.10},
				{"min": 50000, "max": 100000, "value": 0.22},
				{"min": 100000, "max": null, "value": 0.35}
			]
		},
		{
			"name": "status_codes",
			"type": "map",
			"items": {
				"active": 1,
				"inactive": 0
			}
		}
	],
	"rules": [
		{
			"id": "set-salary",
			"type": "assignment",
			"name": "Set Base Salary",
			"config": {
				"target": "salary",
				"value": 75000
			}
		},
		{
			"id": "get-rate",
			"type": "lookup",
			"config": {
				"table": "tax_rates",
				"key": "salary",
				"target": "tax_rate"
			}
		},
		{
			"id": "calc-tax",
			"type": "formula",
			"config": {
				"target": "tax",
				"expression": "salary * tax_rate"
			}
		},
		{
			"id": "dept-split",
			"type": "allocation",
			"config": {
				"source": "salary",
				"strategy": "percentage",
				"targets": [
					{"key": "eng", "amount": 60},
					{"key": "ops", "amount": 40}
				]
			}
		},
		{
			"id": "add-total",
			"type": "buildup",
			"config": {
				"buildup": "total_payroll",
				"operation": "sum",
				"source": "salary",
				"target": "running_total"
			}
		}
	]
}`

func TestParseJSON(t *testing.T) {
	parser := parse.NewParser()
	rs, err := parser.ParseJSON([]byte(testRuleSetJSON))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if rs.Name != "test-rules" {
		t.Errorf("expected name 'test-rules', got %q", rs.Name)
	}
	if len(rs.Lookups) != 2 {
		t.Errorf("expected 2 lookups, got %d", len(rs.Lookups))
	}
	if len(rs.Rules) != 5 {
		t.Errorf("expected 5 rules, got %d", len(rs.Rules))
	}
}

func TestToLookups(t *testing.T) {
	parser := parse.NewParser()
	rs, _ := parser.ParseJSON([]byte(testRuleSetJSON))

	lookups, err := parser.ToLookups(rs)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if len(lookups) != 2 {
		t.Fatalf("expected 2 lookups, got %d", len(lookups))
	}

	// Test range lookup
	taxRates := lookups[0]
	if taxRates.Name() != "tax_rates" {
		t.Errorf("expected name 'tax_rates', got %q", taxRates.Name())
	}

	rate, ok := taxRates.Get(75000.0)
	if !ok {
		t.Error("expected to find rate for 75000")
	}
	if rate != 0.22 {
		t.Errorf("expected rate 0.22, got %v", rate)
	}

	// Test map lookup
	statusCodes := lookups[1]
	status, ok := statusCodes.Get("active")
	if !ok {
		t.Error("expected to find 'active'")
	}
	if status != float64(1) {
		t.Errorf("expected 1, got %v", status)
	}
}

func TestToRules(t *testing.T) {
	parser := parse.NewParser()
	rs, _ := parser.ParseJSON([]byte(testRuleSetJSON))

	rules, err := parser.ToRules(rs)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if len(rules) != 5 {
		t.Fatalf("expected 5 rules, got %d", len(rules))
	}

	// Check rule IDs
	expectedIDs := []string{"set-salary", "get-rate", "calc-tax", "dept-split", "add-total"}
	for i, rule := range rules {
		if rule.ID() != expectedIDs[i] {
			t.Errorf("rule %d: expected ID %q, got %q", i, expectedIDs[i], rule.ID())
		}
	}
}

func TestParseAndBuildEngine(t *testing.T) {
	parser := parse.NewParser()
	engine, err := parser.ParseAndBuildEngine("test", []byte(testRuleSetJSON), cortex.DefaultConfig())
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if engine.Rules() != 5 {
		t.Errorf("expected 5 rules, got %d", engine.Rules())
	}
	if engine.Lookups() != 2 {
		t.Errorf("expected 2 lookups, got %d", engine.Lookups())
	}

	// Run evaluation
	ctx := context.Background()
	evalCtx := cortex.NewEvalContext()

	result, err := engine.Evaluate(ctx, evalCtx)
	if err != nil {
		t.Fatalf("evaluation error: %v", err)
	}

	if !result.Success {
		t.Error("expected success")
	}

	// Check results
	salary, _ := evalCtx.GetFloat64("salary")
	if salary != 75000 {
		t.Errorf("expected salary=75000, got %f", salary)
	}

	rate, _ := evalCtx.GetFloat64("tax_rate")
	if rate != 0.22 {
		t.Errorf("expected tax_rate=0.22, got %f", rate)
	}

	tax, _ := evalCtx.GetFloat64("tax")
	if tax != 16500 {
		t.Errorf("expected tax=16500, got %f", tax)
	}

	eng, _ := evalCtx.GetFloat64("eng")
	if eng != 45000 {
		t.Errorf("expected eng=45000, got %f", eng)
	}

	total, _ := evalCtx.GetFloat64("running_total")
	if total != 75000 {
		t.Errorf("expected running_total=75000, got %f", total)
	}
}

func TestDisabledRules(t *testing.T) {
	json := `{
		"version": "1.0",
		"name": "test",
		"rules": [
			{
				"id": "enabled",
				"type": "assignment",
				"config": {"target": "x", "value": 1}
			},
			{
				"id": "disabled",
				"type": "assignment",
				"disabled": true,
				"config": {"target": "y", "value": 2}
			}
		]
	}`

	parser := parse.NewParser()
	engine, err := parser.ParseAndBuildEngine("test", []byte(json), nil)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if engine.Rules() != 1 {
		t.Errorf("expected 1 rule (disabled should be skipped), got %d", engine.Rules())
	}
}

func TestRegisteredFormula(t *testing.T) {
	json := `{
		"version": "1.0",
		"name": "test",
		"rules": [
			{
				"id": "set-value",
				"type": "assignment",
				"config": {"target": "x", "value": 10}
			},
			{
				"id": "double-it",
				"type": "formula",
				"config": {
					"target": "result",
					"function": "double"
				}
			}
		]
	}`

	parser := parse.NewParser()
	parser.RegisterFormula("double", func(ctx context.Context, evalCtx *cortex.EvalContext) (any, error) {
		x, _ := evalCtx.GetFloat64("x")
		return x * 2, nil
	})

	engine, err := parser.ParseAndBuildEngine("test", []byte(json), nil)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	ctx := context.Background()
	evalCtx := cortex.NewEvalContext()

	_, err = engine.Evaluate(ctx, evalCtx)
	if err != nil {
		t.Fatalf("evaluation error: %v", err)
	}

	result, _ := evalCtx.GetFloat64("result")
	if result != 20 {
		t.Errorf("expected result=20, got %f", result)
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{"invalid json", `{invalid`},
		{"unknown rule type", `{"version":"1.0","name":"test","rules":[{"id":"x","type":"unknown","config":{}}]}`},
		{"missing target", `{"version":"1.0","name":"test","rules":[{"id":"x","type":"assignment","config":{"value":1}}]}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := parse.NewParser()
			_, err := parser.ParseAndBuildEngine("test", []byte(tt.json), nil)
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}
