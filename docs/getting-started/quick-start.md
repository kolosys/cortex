# Quick Start

This guide will help you get started with cortex quickly with a basic example.

## Basic Usage

Here's a simple example to get you started:

```go
package main

import (
    "fmt"
    "log"
    "github.com/kolosys/cortex"
    "github.com/kolosys/cortex/expr"
)

func main() {
    // Basic usage example
    fmt.Println("Welcome to cortex!")
    
    // TODO: Add your code here
}
```

## Common Use Cases

### Using cortex

**Import Path:** `github.com/kolosys/cortex`

Package cortex provides a rules engine for business logic evaluation.

Cortex supports five rule types:
  - Assignment: Set values directly on the context
  - Formula: Calculate values using expressions or functions
  - Allocation: Distribute values across multiple targets
  - Lookup: Retrieve values from lookup tables
  - Buildup: Accumulate/aggregate values (running totals, sums, etc.)

Example:

	engine := cortex.New("payroll", cortex.DefaultConfig())
	engine.RegisterLookup(cortex.NewRangeLookup("tax_brackets", brackets))
	engine.AddRules(
	    cortex.MustAssignment(cortex.AssignmentConfig{...}),
	    cortex.MustFormula(cortex.FormulaConfig{...}),
	)
	result, err := engine.Evaluate(ctx, cortex.NewEvalContext())


```go
package main

import (
    "fmt"
    "github.com/kolosys/cortex"
)

func main() {
    // Example usage of cortex
    fmt.Println("Using cortex package")
}
```

#### Available Types
- **AllocationConfig** - AllocationConfig configures an allocation rule.
- **AllocationRule** - AllocationRule distributes a value across multiple targets.
- **AllocationStrategy** - AllocationStrategy defines how values are distributed.
- **AllocationTarget** - AllocationTarget specifies a single allocation destination.
- **AssignmentConfig** - AssignmentConfig configures an assignment rule.
- **AssignmentRule** - AssignmentRule sets a value on the evaluation context.
- **Buildup** - Buildup represents a running accumulator.
- **BuildupConfig** - BuildupConfig configures a buildup rule.
- **BuildupOperation** - BuildupOperation defines how values are accumulated.
- **BuildupRule** - BuildupRule accumulates values (running totals, aggregations).
- **Config** - Config configures the engine behavior.
- **Engine** - Engine evaluates rules in sequence.
- **EvalContext** - EvalContext holds the state during rule evaluation. It provides thread-safe access to values, buildups, and lookups.
- **EvalMode** - EvalMode determines how errors are handled during evaluation.
- **FormulaConfig** - FormulaConfig configures a formula rule.
- **FormulaFunc** - FormulaFunc computes a value from the evaluation context.
- **FormulaRule** - FormulaRule calculates a value using a function or expression.
- **Logger** - Logger provides a simple logging interface.
- **Lookup** - Lookup represents a lookup table.
- **LookupConfig** - LookupConfig configures a lookup rule.
- **LookupRule** - LookupRule retrieves a value from a lookup table.
- **MapLookup** - MapLookup provides a simple map-based lookup implementation.
- **Metrics** - Metrics provides a simple metrics interface.
- **Observability** - Observability holds observability hooks for the engine.
- **RangeEntry** - RangeEntry represents a single range in a range lookup.
- **RangeLookup** - RangeLookup provides range-based lookup (e.g., tax brackets).
- **Result** - Result contains the outcome of a rule evaluation.
- **Rule** - Rule represents any rule that can be evaluated against an EvalContext.
- **RuleError** - RuleError wraps an error with rule context.
- **RuleMetadata** - RuleMetadata provides optional metadata about a rule.
- **RuleType** - RuleType identifies the type of rule.
- **TaxBracket** - TaxBracket is a convenience type for common tax bracket lookups.
- **Tracer** - Tracer provides a simple tracing interface.
- **ValueFunc** - ValueFunc computes a value dynamically based on context.

#### Available Functions
- **GetTyped** - GetTyped retrieves a typed value from the context.
- **SetTyped** - SetTyped stores a typed value in the context.

For detailed API documentation, see the [cortex API Reference](../api-reference/cortex.md).

### Using expr

**Import Path:** `github.com/kolosys/cortex/expr`

Package expr provides a simple expression DSL for cortex formulas.

Supported operations:
  - Arithmetic: +, -, *, /, %
  - Comparison: ==, !=, <, >, <=, >=
  - Logical: &&, ||, !
  - Functions: min, max, abs, floor, ceil, round, if, sqrt, pow

Example expressions:

	"base_salary * tax_rate"
	"if(age >= 65, senior_discount, 0)"
	"round(total * 0.0825, 2)"
	"min(calculated, max_amount)"


```go
package main

import (
    "fmt"
    "github.com/kolosys/cortex/expr"
)

func main() {
    // Example usage of expr
    fmt.Println("Using expr package")
}
```

#### Available Types
- **BinaryExpr** - BinaryExpr represents a binary expression.
- **BoolLit** - BoolLit represents a boolean literal.
- **CallExpr** - CallExpr represents a function call.
- **Evaluator** - Evaluator evaluates an AST against a value getter.
- **Expression** - Expression represents a compiled expression.
- **Func** - Func is a built-in function type.
- **Ident** - Ident represents an identifier (variable reference).
- **Lexer** - Lexer tokenizes an expression string.
- **Node** - Node represents an AST node.
- **NumberLit** - NumberLit represents a numeric literal.
- **Parser** - Parser parses an expression into an AST.
- **StringLit** - StringLit represents a string literal.
- **Token** - Token represents a lexical token.
- **TokenType** - TokenType represents the type of a token.
- **UnaryExpr** - UnaryExpr represents a unary expression.
- **ValueGetter** - ValueGetter retrieves values by name (e.g., from EvalContext).

For detailed API documentation, see the [expr API Reference](../api-reference/expr.md).

## Step-by-Step Tutorial

### Step 1: Import the Package

First, import the necessary packages in your Go file:

```go
import (
    "fmt"
    "github.com/kolosys/cortex"
    "github.com/kolosys/cortex/expr"
)
```

### Step 2: Initialize

Set up the basic configuration:

```go
func main() {
    // Initialize your application
    fmt.Println("Initializing cortex...")
}
```

### Step 3: Use the Library

Implement your specific use case:

```go
func main() {
    // Your implementation here
}
```

## Running Your Code

To run your Go program:

```bash
go run main.go
```

To build an executable:

```bash
go build -o myapp
./myapp
```

## Configuration Options

cortex can be configured to suit your needs. Check the [Core Concepts](../core-concepts/) section for detailed information about configuration options.

## Error Handling

Always handle errors appropriately:

```go
result, err := someFunction()
if err != nil {
    log.Fatalf("Error: %v", err)
}
```

## Best Practices

- Always handle errors returned by library functions
- Check the API documentation for detailed parameter information
- Use meaningful variable and function names
- Add comments to document your code

## Complete Example

Here's a complete working example:

```go
package main

import (
    "fmt"
    "log"
    "github.com/kolosys/cortex"
    "github.com/kolosys/cortex/expr"
)

func main() {
    fmt.Println("Starting cortex application...")
    
    // Add your implementation here
    
    fmt.Println("Application completed successfully!")
}
```

## Next Steps

Now that you've seen the basics, explore:

- **[Core Concepts](../core-concepts/)** - Understanding the library architecture
- **[API Reference](../api-reference/)** - Complete API documentation
- **[Examples](../examples/README.md)** - More detailed examples
- **[Advanced Topics](../advanced/)** - Performance tuning and advanced patterns

## Getting Help

If you run into issues:

1. Check the [API Reference](../api-reference/)
2. Browse the [Examples](../examples/README.md)
3. Visit the [GitHub Issues](https://github.com/kolosys/cortex/issues) page

