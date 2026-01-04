# Installation

This guide will help you install and set up cortex in your Go project.

## Prerequisites

Before installing cortex, ensure you have:

- **Go 1.21** or later installed
- A Go module initialized in your project (run `go mod init` if needed)
- Access to the GitHub repository (for private repositories)

## Installation Steps

### Step 1: Install the Package

Use `go get` to install cortex:

```bash
go get github.com/kolosys/cortex
```

This will download the package and add it to your `go.mod` file.

### Step 2: Import in Your Code

Import the package in your Go source files:

```go
import "github.com/kolosys/cortex"
```

### Multiple Packages

cortex includes several packages. Import the ones you need:

```go
// Package cortex provides a rules engine for business logic evaluation.

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

import "github.com/kolosys/cortex"
```

```go
// Package expr provides a simple expression DSL for cortex formulas.

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

import "github.com/kolosys/cortex/expr"
```

### Step 3: Verify Installation

Create a simple test file to verify the installation:

```go
package main

import (
    "fmt"
    "github.com/kolosys/cortex"
)

func main() {
    fmt.Println("cortex installed successfully!")
}
```

Run the test:

```bash
go run main.go
```

## Updating the Package

To update to the latest version:

```bash
go get -u github.com/kolosys/cortex
```

To update to a specific version:

```bash
go get github.com/kolosys/cortex@v1.2.3
```

## Installing a Specific Version

To install a specific version of the package:

```bash
go get github.com/kolosys/cortex@v1.0.0
```

Check available versions on the [GitHub releases page](https://github.com/kolosys/cortex/releases).

## Development Setup

If you want to contribute or modify the library:

### Clone the Repository

```bash
git clone https://github.com/kolosys/cortex.git
cd cortex
```

### Install Dependencies

```bash
go mod download
```

### Run Tests

```bash
go test ./...
```

## Troubleshooting

### Module Not Found

If you encounter a "module not found" error:

1. Ensure your `GOPATH` is set correctly
2. Check that you have network access to GitHub
3. Try running `go clean -modcache` and reinstall

### Private Repository Access

For private repositories, configure Git to use SSH or a personal access token:

```bash
git config --global url."git@github.com:".insteadOf "https://github.com/"
```

Or set up GOPRIVATE:

```bash
export GOPRIVATE=github.com/kolosys/cortex
```

## Next Steps

Now that you have cortex installed, check out the [Quick Start Guide](quick-start.md) to learn how to use it.

## Additional Resources

- [Quick Start Guide](quick-start.md)
- [API Reference](../api-reference/)
- [Examples](../examples/README.md)

