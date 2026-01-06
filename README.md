# Cortex

A general-purpose rules engine for Go supporting multiple rule types, sequential evaluation, and config-driven rule definitions.

![GoVersion](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-blue.svg)
![Zero Dependencies](https://img.shields.io/badge/Zero-Dependencies-green.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/kolosys/cortex.svg)](https://pkg.go.dev/github.com/kolosys/cortex)
[![Go Report Card](https://goreportcard.com/badge/github.com/kolosys/cortex)](https://goreportcard.com/report/github.com/kolosys/cortex)

## Features

- **Five Rule Types**: Assignment, Formula, Allocation, Lookup, Buildup
- **Expression DSL**: Simple expressions for config-driven formulas
- **Thread-Safe**: Concurrent-safe evaluation context
- **Observable**: Pluggable logging, metrics, and tracing
- **Zero Dependencies**: stdlib only

## Installation

```bash
go get github.com/kolosys/cortex
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "math"

    "github.com/kolosys/cortex"
)

func main() {
    engine := cortex.New("payroll", cortex.DefaultConfig())

    // Register lookup table
    engine.RegisterLookup(cortex.NewRangeLookup("tax_brackets", []cortex.RangeEntry[float64]{
        {Min: 0, Max: 50000, Value: 0.10},
        {Min: 50000, Max: 100000, Value: 0.22},
        {Min: 100000, Max: math.Inf(1), Value: 0.35},
    }))

    // Add rules
    engine.AddRules(
        cortex.MustAssignment(cortex.AssignmentConfig{
            ID: "salary", Target: "salary", Value: 75000.0,
        }),
        cortex.MustLookup(cortex.LookupConfig{
            ID: "rate", Table: "tax_brackets", Key: "salary", Target: "tax_rate",
        }),
        cortex.MustFormula(cortex.FormulaConfig{
            ID: "tax", Target: "tax", Expression: "salary * tax_rate",
        }),
        cortex.MustAllocation(cortex.AllocationConfig{
            ID: "split", Source: "salary", Strategy: cortex.StrategyPercentage,
            Targets: []cortex.AllocationTarget{
                {Key: "dept_eng", Amount: 60},
                {Key: "dept_ops", Amount: 40},
            },
        }),
    )

    // Evaluate
    evalCtx := cortex.NewEvalContext()
    result, _ := engine.Evaluate(context.Background(), evalCtx)

    fmt.Printf("Success: %v, Rules: %d\n", result.Success, result.RulesEvaluated)

    tax, _ := evalCtx.GetFloat64("tax")
    fmt.Printf("Tax: $%.2f\n", tax) // Tax: $16500.00
}
```

## Rule Types

### Assignment

Set values directly on the context:

```go
cortex.MustAssignment(cortex.AssignmentConfig{
    ID:     "set-salary",
    Target: "salary",
    Value:  75000.0,
})
```

### Formula

Calculate values using expressions or Go functions:

```go
// Expression-based
cortex.MustFormula(cortex.FormulaConfig{
    ID:         "calc-tax",
    Target:     "tax",
    Expression: "salary * tax_rate",
})

// Function-based
cortex.MustFormula(cortex.FormulaConfig{
    ID:     "calc-bonus",
    Target: "bonus",
    Formula: func(ctx context.Context, e *cortex.EvalContext) (any, error) {
        salary, _ := e.GetFloat64("salary")
        return salary * 0.10, nil
    },
})
```

### Allocation

Distribute values across targets:

```go
cortex.MustAllocation(cortex.AllocationConfig{
    ID:       "split",
    Source:   "budget",
    Strategy: cortex.StrategyPercentage,
    Targets: []cortex.AllocationTarget{
        {Key: "eng", Amount: 50},
        {Key: "ops", Amount: 30},
        {Key: "admin", Amount: 20},
    },
})
```

**Strategies**: `StrategyPercentage`, `StrategyFixed`, `StrategyWeighted`, `StrategyEqual`, `StrategyRatio`

### Lookup

Retrieve values from lookup tables:

```go
// Map lookup
engine.RegisterLookup(cortex.NewMapLookup("status", map[string]int{
    "active": 1, "inactive": 0,
}))

// Range lookup (tax brackets)
engine.RegisterLookup(cortex.NewRangeLookup("rates", []cortex.RangeEntry[float64]{
    {Min: 0, Max: 50000, Value: 0.10},
    {Min: 50000, Max: 100000, Value: 0.22},
}))

cortex.MustLookup(cortex.LookupConfig{
    ID:     "get-rate",
    Table:  "rates",
    Key:    "income",
    Target: "rate",
})
```

### Buildup

Accumulate values (running totals):

```go
cortex.MustBuildup(cortex.BuildupConfig{
    ID:        "add-total",
    Buildup:   "running_total",
    Operation: cortex.BuildupSum,
    Source:    "amount",
    Target:    "current_total",
})
```

**Operations**: `BuildupSum`, `BuildupMin`, `BuildupMax`, `BuildupAvg`, `BuildupCount`, `BuildupProduct`

## Expression DSL

Supported in config-driven formulas:

```
Arithmetic:  +, -, *, /, %
Comparison:  ==, !=, <, >, <=, >=
Logical:     &&, ||, !
Functions:   min, max, abs, floor, ceil, round, if, sqrt, pow

Examples:
  "base_salary * tax_rate"
  "if(age >= 65, senior_discount, 0)"
  "round(total * 0.0825, 2)"
```

## Config-Driven Rules (JSON)

```go
import "github.com/kolosys/cortex/parse"

json := `{
    "version": "1.0",
    "name": "payroll",
    "rules": [
        {"id": "salary", "type": "assignment", "config": {"target": "salary", "value": 75000}},
        {"id": "tax", "type": "formula", "config": {"target": "tax", "expression": "salary * 0.22"}}
    ]
}`

parser := parse.NewParser()
engine, _ := parser.ParseAndBuildEngine("payroll", []byte(json), nil)
```

## Evaluation Modes

- `ModeFailFast` (default): Stop on first error
- `ModeCollectAll`: Evaluate all rules, collect errors
- `ModeContinueOnError`: Log errors, continue evaluation

## License

MIT
