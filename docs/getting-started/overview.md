# Overview

Kolosys Template for GO projects

## About cortex

This documentation provides comprehensive guidance for using cortex, a Go library designed to help you build better software.

## Project Information

- **Repository**: [https://github.com/kolosys/cortex](https://github.com/kolosys/cortex)
- **Import Path**: `github.com/kolosys/cortex`
- **License**: MIT
- **Version**: latest

## What You'll Find Here

This documentation is organized into several sections to help you find what you need:

- **[Getting Started](../getting-started/)** - Installation instructions and quick start guides
- **[Core Concepts](../core-concepts/)** - Fundamental concepts and architecture details
- **[Advanced Topics](../advanced/)** - Performance tuning and advanced usage patterns
- **[API Reference](../api-reference/)** - Complete API reference documentation
- **[Examples](../examples/)** - Working code examples and tutorials

## Project Features

cortex provides:
- **cortex** - Package cortex provides a rules engine for business logic evaluation.

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

- **expr** - Package expr provides a simple expression DSL for cortex formulas.

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


## Quick Links

- [Installation Guide](installation.md)
- [Quick Start Guide](quick-start.md)
- [API Reference](../api-reference/)
- [Examples](../examples/README.md)

## Community & Support

- **GitHub Issues**: [https://github.com/kolosys/cortex/issues](https://github.com/kolosys/cortex/issues)
- **Discussions**: [https://github.com/kolosys/cortex/discussions](https://github.com/kolosys/cortex/discussions)
- **Repository Owner**: [kolosys](https://github.com/kolosys)

## Getting Help

If you encounter any issues or have questions:

1. Check the [API Reference](../api-reference/) for detailed documentation
2. Browse the [Examples](../examples/README.md) for common use cases
3. Search existing [GitHub Issues](https://github.com/kolosys/cortex/issues)
4. Open a new issue if you've found a bug or have a feature request

## Next Steps

Ready to get started? Head over to the [Installation Guide](installation.md) to begin using cortex.

