# cortex API

Complete API documentation for the cortex package.

**Import Path:** `github.com/kolosys/cortex`

## Package Documentation

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


## Variables

**ErrRuleNotFound, ErrLookupNotFound, ErrKeyNotFound, ErrValueNotFound, ErrBuildupNotFound, ErrInvalidRule, ErrInvalidExpression, ErrTypeMismatch, ErrDivisionByZero, ErrAllocationSum, ErrCircularDep, ErrEvaluation, ErrShortCircuit, ErrEngineClosed, ErrTimeout, ErrNilContext, ErrDuplicateRule, ErrDuplicateLookup**

Sentinel errors for common failure cases.


```go
var ErrRuleNotFound = errors.New("cortex: rule not found")
var ErrLookupNotFound = errors.New("cortex: lookup table not found")
var ErrKeyNotFound = errors.New("cortex: key not found in lookup")
var ErrValueNotFound = errors.New("cortex: value not found in context")
var ErrBuildupNotFound = errors.New("cortex: buildup not found")
var ErrInvalidRule = errors.New("cortex: invalid rule configuration")
var ErrInvalidExpression = errors.New("cortex: invalid expression")
var ErrTypeMismatch = errors.New("cortex: type mismatch")
var ErrDivisionByZero = errors.New("cortex: division by zero")
var ErrAllocationSum = errors.New("cortex: allocation percentages must sum to 100")
var ErrCircularDep = errors.New("cortex: circular dependency detected")
var ErrEvaluation = errors.New("cortex: evaluation failed")
var ErrShortCircuit = errors.New("cortex: evaluation short-circuited")
var ErrEngineClosed = errors.New("cortex: engine is closed")
var ErrTimeout = errors.New("cortex: evaluation timeout")
var ErrNilContext = errors.New("cortex: nil evaluation context")
var ErrDuplicateRule = errors.New("cortex: duplicate rule ID")
var ErrDuplicateLookup = errors.New("cortex: duplicate lookup table name")
```

## Types

### AllocationConfig
AllocationConfig configures an allocation rule.

#### Example Usage

```go
// Create a new AllocationConfig
allocationconfig := AllocationConfig{
    ID: "example",
    Name: "example",
    Description: "example",
    Deps: [],
    Source: "example",
    Strategy: AllocationStrategy{},
    Targets: [],
    Remainder: "example",
    Precision: 42,
}
```

#### Type Definition

```go
type AllocationConfig struct {
    ID string
    Name string
    Description string
    Deps []string
    Source string
    Strategy AllocationStrategy
    Targets []AllocationTarget
    Remainder string
    Precision int
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| ID | `string` |  |
| Name | `string` |  |
| Description | `string` |  |
| Deps | `[]string` |  |
| Source | `string` | Source is the context key containing the value to allocate. |
| Strategy | `AllocationStrategy` | Strategy is the distribution method. |
| Targets | `[]AllocationTarget` | Targets are the allocation destinations. |
| Remainder | `string` | Remainder is an optional context key for the rounding remainder. |
| Precision | `int` | Precision is the decimal precision (default 2). |

### AllocationRule
AllocationRule distributes a value across multiple targets.

#### Example Usage

```go
// Create a new AllocationRule
allocationrule := AllocationRule{

}
```

#### Type Definition

```go
type AllocationRule struct {
}
```

### Constructor Functions

### MustAllocation

MustAllocation creates a new allocation rule, panicking on error.

```go
func MustAllocation(cfg AllocationConfig) *AllocationRule
```

**Parameters:**
- `cfg` (AllocationConfig)

**Returns:**
- *AllocationRule

### NewAllocation

NewAllocation creates a new allocation rule.

```go
func NewAllocation(cfg AllocationConfig) (*AllocationRule, error)
```

**Parameters:**
- `cfg` (AllocationConfig)

**Returns:**
- *AllocationRule
- error

## Methods

### Dependencies



```go
func (*baseRule) Dependencies() []string
```

**Parameters:**
  None

**Returns:**
- []string

### Description



```go
func (*baseRule) Description() string
```

**Parameters:**
  None

**Returns:**
- string

### Evaluate

Evaluate distributes the source value across targets.

```go
func (*BuildupRule) Evaluate(ctx context.Context, evalCtx *EvalContext) error
```

**Parameters:**
- `ctx` (context.Context)
- `evalCtx` (*EvalContext)

**Returns:**
- error

### ID



```go
func (*baseRule) ID() string
```

**Parameters:**
  None

**Returns:**
- string

### Name



```go
func (**ast.IndexExpr) Name() string
```

**Parameters:**
  None

**Returns:**
- string

### Source

Source returns the source key.

```go
func (*AllocationRule) Source() string
```

**Parameters:**
  None

**Returns:**
- string

### Strategy

Strategy returns the allocation strategy.

```go
func (*AllocationRule) Strategy() AllocationStrategy
```

**Parameters:**
  None

**Returns:**
- AllocationStrategy

### Targets

Targets returns the allocation targets.

```go
func (*AllocationRule) Targets() []AllocationTarget
```

**Parameters:**
  None

**Returns:**
- []AllocationTarget

### AllocationStrategy
AllocationStrategy defines how values are distributed.

#### Example Usage

```go
// Example usage of AllocationStrategy
var value AllocationStrategy
// Initialize with appropriate value
```

#### Type Definition

```go
type AllocationStrategy int
```

### Constructor Functions

### ParseAllocationStrategy

ParseAllocationStrategy parses a string into an AllocationStrategy.

```go
func ParseAllocationStrategy(s string) (AllocationStrategy, error)
```

**Parameters:**
- `s` (string)

**Returns:**
- AllocationStrategy
- error

## Methods

### String



```go
func (BuildupOperation) String() string
```

**Parameters:**
  None

**Returns:**
- string

### AllocationTarget
AllocationTarget specifies a single allocation destination.

#### Example Usage

```go
// Create a new AllocationTarget
allocationtarget := AllocationTarget{
    Key: "example",
    Amount: 3.14,
}
```

#### Type Definition

```go
type AllocationTarget struct {
    Key string
    Amount float64
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Key | `string` | context key to set |
| Amount | `float64` | percentage, fixed amount, weight, or ratio (based on strategy) |

### AssignmentConfig
AssignmentConfig configures an assignment rule.

#### Example Usage

```go
// Create a new AssignmentConfig
assignmentconfig := AssignmentConfig{
    ID: "example",
    Name: "example",
    Description: "example",
    Deps: [],
    Target: "example",
    Value: any{},
    ValueFunc: ValueFunc{},
}
```

#### Type Definition

```go
type AssignmentConfig struct {
    ID string
    Name string
    Description string
    Deps []string
    Target string
    Value any
    ValueFunc ValueFunc
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| ID | `string` |  |
| Name | `string` |  |
| Description | `string` |  |
| Deps | `[]string` |  |
| Target | `string` | Target is the context key to set. |
| Value | `any` | Value is the static value to assign (mutually exclusive with ValueFunc). |
| ValueFunc | `ValueFunc` | ValueFunc computes the value dynamically (mutually exclusive with Value). |

### AssignmentRule
AssignmentRule sets a value on the evaluation context.

#### Example Usage

```go
// Create a new AssignmentRule
assignmentrule := AssignmentRule{

}
```

#### Type Definition

```go
type AssignmentRule struct {
}
```

### Constructor Functions

### MustAssignment

MustAssignment creates a new assignment rule, panicking on error.

```go
func MustAssignment(cfg AssignmentConfig) *AssignmentRule
```

**Parameters:**
- `cfg` (AssignmentConfig)

**Returns:**
- *AssignmentRule

### NewAssignment

NewAssignment creates a new assignment rule.

```go
func NewAssignment(cfg AssignmentConfig) (*AssignmentRule, error)
```

**Parameters:**
- `cfg` (AssignmentConfig)

**Returns:**
- *AssignmentRule
- error

## Methods

### Dependencies



```go
func (*baseRule) Dependencies() []string
```

**Parameters:**
  None

**Returns:**
- []string

### Description



```go
func (*baseRule) Description() string
```

**Parameters:**
  None

**Returns:**
- string

### Evaluate

Evaluate sets the value on the evaluation context.

```go
func (*AllocationRule) Evaluate(ctx context.Context, evalCtx *EvalContext) error
```

**Parameters:**
- `ctx` (context.Context)
- `evalCtx` (*EvalContext)

**Returns:**
- error

### ID



```go
func (*baseRule) ID() string
```

**Parameters:**
  None

**Returns:**
- string

### Name



```go
func (*Engine) Name() string
```

**Parameters:**
  None

**Returns:**
- string

### Target

Target returns the target key for this assignment.

```go
func (*FormulaRule) Target() string
```

**Parameters:**
  None

**Returns:**
- string

### Buildup
Buildup represents a running accumulator.

#### Example Usage

```go
// Create a new Buildup
buildup := Buildup{
    Name: "example",
    Operation: BuildupOperation{},
}
```

#### Type Definition

```go
type Buildup struct {
    Name string
    Operation BuildupOperation
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Name | `string` |  |
| Operation | `BuildupOperation` |  |

## Methods

### Add

Add adds a value to the buildup.

```go
func (*Buildup) Add(value float64)
```

**Parameters:**
- `value` (float64)

**Returns:**
  None

### Count

Count returns the number of values added.

```go
func (*Buildup) Count() int64
```

**Parameters:**
  None

**Returns:**
- int64

### Current

Current returns the current accumulated value.

```go
func (*Buildup) Current() float64
```

**Parameters:**
  None

**Returns:**
- float64

### Reset

Reset resets the buildup to its initial state.

```go
func (*Buildup) Reset(initial float64)
```

**Parameters:**
- `initial` (float64)

**Returns:**
  None

### BuildupConfig
BuildupConfig configures a buildup rule.

#### Example Usage

```go
// Create a new BuildupConfig
buildupconfig := BuildupConfig{
    ID: "example",
    Name: "example",
    Description: "example",
    Deps: [],
    Buildup: "example",
    Operation: BuildupOperation{},
    Source: "example",
    Initial: 3.14,
    Target: "example",
}
```

#### Type Definition

```go
type BuildupConfig struct {
    ID string
    Name string
    Description string
    Deps []string
    Buildup string
    Operation BuildupOperation
    Source string
    Initial float64
    Target string
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| ID | `string` |  |
| Name | `string` |  |
| Description | `string` |  |
| Deps | `[]string` |  |
| Buildup | `string` | Buildup is the buildup accumulator name (created if not exists). |
| Operation | `BuildupOperation` | Operation is the accumulation operation. |
| Source | `string` | Source is the context key containing the value to add. |
| Initial | `float64` | Initial is the initial value for the buildup (if created). |
| Target | `string` | Target is an optional context key to write the current value after adding. |

### BuildupOperation
BuildupOperation defines how values are accumulated.

#### Example Usage

```go
// Example usage of BuildupOperation
var value BuildupOperation
// Initialize with appropriate value
```

#### Type Definition

```go
type BuildupOperation int
```

### Constructor Functions

### ParseBuildupOperation

ParseBuildupOperation parses a string into a BuildupOperation.

```go
func ParseBuildupOperation(s string) (BuildupOperation, error)
```

**Parameters:**
- `s` (string)

**Returns:**
- BuildupOperation
- error

## Methods

### String



```go
func (BuildupOperation) String() string
```

**Parameters:**
  None

**Returns:**
- string

### BuildupRule
BuildupRule accumulates values (running totals, aggregations).

#### Example Usage

```go
// Create a new BuildupRule
builduprule := BuildupRule{

}
```

#### Type Definition

```go
type BuildupRule struct {
}
```

### Constructor Functions

### MustBuildup

MustBuildup creates a new buildup rule, panicking on error.

```go
func MustBuildup(cfg BuildupConfig) *BuildupRule
```

**Parameters:**
- `cfg` (BuildupConfig)

**Returns:**
- *BuildupRule

### NewBuildup

NewBuildup creates a new buildup rule.

```go
func NewBuildup(cfg BuildupConfig) (*BuildupRule, error)
```

**Parameters:**
- `cfg` (BuildupConfig)

**Returns:**
- *BuildupRule
- error

## Methods

### BuildupName

BuildupName returns the buildup accumulator name.

```go
func (*BuildupRule) BuildupName() string
```

**Parameters:**
  None

**Returns:**
- string

### Dependencies



```go
func (*baseRule) Dependencies() []string
```

**Parameters:**
  None

**Returns:**
- []string

### Description



```go
func (*baseRule) Description() string
```

**Parameters:**
  None

**Returns:**
- string

### Evaluate

Evaluate adds to the buildup accumulator.

```go
func (*BuildupRule) Evaluate(ctx context.Context, evalCtx *EvalContext) error
```

**Parameters:**
- `ctx` (context.Context)
- `evalCtx` (*EvalContext)

**Returns:**
- error

### ID



```go
func (*baseRule) ID() string
```

**Parameters:**
  None

**Returns:**
- string

### Name



```go
func (**ast.IndexExpr) Name() string
```

**Parameters:**
  None

**Returns:**
- string

### Config
Config configures the engine behavior.

#### Example Usage

```go
// Create a new Config
config := Config{
    Mode: EvalMode{},
    Timeout: /* value */,
    ShortCircuit: true,
    EnableMetrics: true,
    MaxRules: 42,
}
```

#### Type Definition

```go
type Config struct {
    Mode EvalMode
    Timeout time.Duration
    ShortCircuit bool
    EnableMetrics bool
    MaxRules int
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Mode | `EvalMode` | Mode determines error handling behavior. |
| Timeout | `time.Duration` | Timeout for entire evaluation (0 = no timeout). |
| ShortCircuit | `bool` | ShortCircuit allows rules to halt evaluation early. |
| EnableMetrics | `bool` | EnableMetrics enables detailed metrics collection. |
| MaxRules | `int` | MaxRules limits the number of rules (0 = unlimited). |

### Constructor Functions

### DefaultConfig

DefaultConfig returns a Config with sensible defaults.

```go
func DefaultConfig() *Config
```

**Parameters:**
  None

**Returns:**
- *Config

## Methods

### Validate

Validate validates the configuration.

```go
func (*Config) Validate() error
```

**Parameters:**
  None

**Returns:**
- error

### Engine
Engine evaluates rules in sequence.

#### Example Usage

```go
// Create a new Engine
engine := Engine{

}
```

#### Type Definition

```go
type Engine struct {
}
```

### Constructor Functions

### New

New creates a new rules engine.

```go
func New(name string, config *Config) *Engine
```

**Parameters:**
- `name` (string)
- `config` (*Config)

**Returns:**
- *Engine

## Methods

### AddRule

AddRule adds a rule to the engine.

```go
func (*Engine) AddRule(rule Rule) error
```

**Parameters:**
- `rule` (Rule)

**Returns:**
- error

### AddRules

AddRules adds multiple rules to the engine.

```go
func (*Engine) AddRules(rules ...Rule) error
```

**Parameters:**
- `rules` (...Rule)

**Returns:**
- error

### Clone

Clone creates a copy of the engine with the same configuration and lookups, but without any rules.

```go
func (*Engine) Clone(name string) *Engine
```

**Parameters:**
- `name` (string)

**Returns:**
- *Engine

### Close

Close closes the engine.

```go
func (*Engine) Close() error
```

**Parameters:**
  None

**Returns:**
- error

### Config

Config returns the engine configuration.

```go
func (*Engine) Config() *Config
```

**Parameters:**
  None

**Returns:**
- *Config

### Evaluate

Evaluate runs all rules against the provided context.

```go
func (*AllocationRule) Evaluate(ctx context.Context, evalCtx *EvalContext) error
```

**Parameters:**
- `ctx` (context.Context)
- `evalCtx` (*EvalContext)

**Returns:**
- error

### Lookups

Lookups returns the number of registered lookups.

```go
func (*Engine) Lookups() int
```

**Parameters:**
  None

**Returns:**
- int

### Name

Name returns the engine name.

```go
func (*Engine) Name() string
```

**Parameters:**
  None

**Returns:**
- string

### RegisterLookup

RegisterLookup registers a lookup table.

```go
func (*Engine) RegisterLookup(lookup Lookup) error
```

**Parameters:**
- `lookup` (Lookup)

**Returns:**
- error

### RegisterLookups

RegisterLookups registers multiple lookup tables.

```go
func (*Engine) RegisterLookups(lookups ...Lookup) error
```

**Parameters:**
- `lookups` (...Lookup)

**Returns:**
- error

### Rules

Rules returns the number of rules.

```go
func (*Engine) Rules() int
```

**Parameters:**
  None

**Returns:**
- int

### WithObservability

WithObservability sets observability hooks.

```go
func (*Engine) WithObservability(obs *Observability) *Engine
```

**Parameters:**
- `obs` (*Observability)

**Returns:**
- *Engine

### EvalContext
EvalContext holds the state during rule evaluation. It provides thread-safe access to values, buildups, and lookups.

#### Example Usage

```go
// Create a new EvalContext
evalcontext := EvalContext{
    ID: "example",
}
```

#### Type Definition

```go
type EvalContext struct {
    ID string
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| ID | `string` |  |

### Constructor Functions

### NewEvalContext

NewEvalContext creates a new evaluation context.

```go
func NewEvalContext() *EvalContext
```

**Parameters:**
  None

**Returns:**
- *EvalContext

## Methods

### Clone

Clone creates a shallow copy of the context.

```go
func (*Engine) Clone(name string) *Engine
```

**Parameters:**
- `name` (string)

**Returns:**
- *Engine

### Delete

Delete removes a value from the context.

```go
func (*EvalContext) Delete(key string)
```

**Parameters:**
- `key` (string)

**Returns:**
  None

### Duration

Duration returns the time since context creation.

```go
func (*EvalContext) Duration() time.Duration
```

**Parameters:**
  None

**Returns:**
- time.Duration

### ErrorCount

ErrorCount returns the number of errors encountered.

```go
func (*EvalContext) ErrorCount() int64
```

**Parameters:**
  None

**Returns:**
- int64

### Get

Get retrieves a value from the context.

```go
func (*EvalContext) Get(key string) (any, bool)
```

**Parameters:**
- `key` (string)

**Returns:**
- any
- bool

### GetBool

GetBool retrieves a bool value from the context.

```go
func (*EvalContext) GetBool(key string) (bool, error)
```

**Parameters:**
- `key` (string)

**Returns:**
- bool
- error

### GetBuildup

GetBuildup returns a buildup accumulator for the given key.

```go
func (*EvalContext) GetBuildup(key string) (*Buildup, bool)
```

**Parameters:**
- `key` (string)

**Returns:**
- *Buildup
- bool

### GetFloat64

GetFloat64 retrieves a float64 value from the context.

```go
func (*EvalContext) GetFloat64(key string) (float64, error)
```

**Parameters:**
- `key` (string)

**Returns:**
- float64
- error

### GetInt

GetInt retrieves an int value from the context.

```go
func (*EvalContext) GetInt(key string) (int, error)
```

**Parameters:**
- `key` (string)

**Returns:**
- int
- error

### GetMetadata

GetMetadata retrieves a metadata value.

```go
func (*EvalContext) GetMetadata(key string) (string, bool)
```

**Parameters:**
- `key` (string)

**Returns:**
- string
- bool

### GetOrCreateBuildup

GetOrCreateBuildup returns an existing buildup or creates a new one.

```go
func (*EvalContext) GetOrCreateBuildup(key string, op BuildupOperation, initial float64) *Buildup
```

**Parameters:**
- `key` (string)
- `op` (BuildupOperation)
- `initial` (float64)

**Returns:**
- *Buildup

### GetString

GetString retrieves a string value from the context.

```go
func (*EvalContext) GetString(key string) (string, error)
```

**Parameters:**
- `key` (string)

**Returns:**
- string
- error

### Halt

Halt stops evaluation with the given rule ID.

```go
func (*EvalContext) Halt(ruleID string)
```

**Parameters:**
- `ruleID` (string)

**Returns:**
  None

### HaltedBy

HaltedBy returns the rule ID that halted evaluation.

```go
func (*EvalContext) HaltedBy() string
```

**Parameters:**
  None

**Returns:**
- string

### Has

Has checks if a key exists in the context.

```go
func (*EvalContext) Has(key string) bool
```

**Parameters:**
- `key` (string)

**Returns:**
- bool

### IsHalted

IsHalted returns true if evaluation has been halted.

```go
func (*EvalContext) IsHalted() bool
```

**Parameters:**
  None

**Returns:**
- bool

### Keys

Keys returns all keys in the context.

```go
func (*EvalContext) Keys() []string
```

**Parameters:**
  None

**Returns:**
- []string

### Lookup

Lookup performs a lookup in the specified table.

```go
func (*EvalContext) Lookup(tableName string, key any) (any, bool, error)
```

**Parameters:**
- `tableName` (string)
- `key` (any)

**Returns:**
- any
- bool
- error

### RegisterLookup

RegisterLookup registers a lookup table in the context.

```go
func (*EvalContext) RegisterLookup(lookup Lookup)
```

**Parameters:**
- `lookup` (Lookup)

**Returns:**
  None

### RulesEvaluated

RulesEvaluated returns the number of rules evaluated.

```go
func (*EvalContext) RulesEvaluated() int64
```

**Parameters:**
  None

**Returns:**
- int64

### Set

Set stores a value in the context.

```go
func (*EvalContext) Set(key string, value any)
```

**Parameters:**
- `key` (string)
- `value` (any)

**Returns:**
  None

### SetMetadata

SetMetadata sets a metadata key-value pair.

```go
func (*EvalContext) SetMetadata(key, value string)
```

**Parameters:**
- `key` (string)
- `value` (string)

**Returns:**
  None

### Values

Values returns a copy of all values in the context.

```go
func (*EvalContext) Values() map[string]any
```

**Parameters:**
  None

**Returns:**
- map[string]any

### EvalMode
EvalMode determines how errors are handled during evaluation.

#### Example Usage

```go
// Example usage of EvalMode
var value EvalMode
// Initialize with appropriate value
```

#### Type Definition

```go
type EvalMode int
```

## Methods

### String



```go
func (AllocationStrategy) String() string
```

**Parameters:**
  None

**Returns:**
- string

### FormulaConfig
FormulaConfig configures a formula rule.

#### Example Usage

```go
// Create a new FormulaConfig
formulaconfig := FormulaConfig{
    ID: "example",
    Name: "example",
    Description: "example",
    Deps: [],
    Target: "example",
    Inputs: [],
    Formula: FormulaFunc{},
    Expression: "example",
}
```

#### Type Definition

```go
type FormulaConfig struct {
    ID string
    Name string
    Description string
    Deps []string
    Target string
    Inputs []string
    Formula FormulaFunc
    Expression string
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| ID | `string` |  |
| Name | `string` |  |
| Description | `string` |  |
| Deps | `[]string` |  |
| Target | `string` | Target is the context key to store the result. |
| Inputs | `[]string` | Inputs are the required input keys (for dependency tracking). |
| Formula | `FormulaFunc` | Formula is the Go function for complex rules. |
| Expression | `string` | Expression is the expression string for config-driven rules. (Mutually exclusive with Formula) |

### FormulaFunc
FormulaFunc computes a value from the evaluation context.

#### Example Usage

```go
// Example usage of FormulaFunc
var value FormulaFunc
// Initialize with appropriate value
```

#### Type Definition

```go
type FormulaFunc func(ctx context.Context, evalCtx *EvalContext) (any, error)
```

### Constructor Functions

### Add

Add returns a FormulaFunc that adds two context values.

```go
func (nopMetrics) Add(name string, v float64, kv ...any)
```

**Parameters:**
- `name` (string)
- `v` (float64)
- `kv` (...any)

**Returns:**
  None

### Conditional

Conditional returns a FormulaFunc that returns thenVal if condition is true, else elseVal.

```go
func Conditional(condition string, thenVal, elseVal any) FormulaFunc
```

**Parameters:**
- `condition` (string)
- `thenVal` (any)
- `elseVal` (any)

**Returns:**
- FormulaFunc

### Divide

Divide returns a FormulaFunc that divides a by b.

```go
func Divide(a, b string) FormulaFunc
```

**Parameters:**
- `a` (string)
- `b` (string)

**Returns:**
- FormulaFunc

### Multiply

Multiply returns a FormulaFunc that multiplies two context values.

```go
func Multiply(a, b string) FormulaFunc
```

**Parameters:**
- `a` (string)
- `b` (string)

**Returns:**
- FormulaFunc

### Percentage

Percentage returns a FormulaFunc that calculates a percentage of a value.

```go
func Percentage(value string, percent float64) FormulaFunc
```

**Parameters:**
- `value` (string)
- `percent` (float64)

**Returns:**
- FormulaFunc

### Subtract

Subtract returns a FormulaFunc that subtracts b from a.

```go
func Subtract(a, b string) FormulaFunc
```

**Parameters:**
- `a` (string)
- `b` (string)

**Returns:**
- FormulaFunc

### FormulaRule
FormulaRule calculates a value using a function or expression.

#### Example Usage

```go
// Create a new FormulaRule
formularule := FormulaRule{

}
```

#### Type Definition

```go
type FormulaRule struct {
}
```

### Constructor Functions

### MustFormula

MustFormula creates a new formula rule, panicking on error.

```go
func MustFormula(cfg FormulaConfig) *FormulaRule
```

**Parameters:**
- `cfg` (FormulaConfig)

**Returns:**
- *FormulaRule

### NewFormula

NewFormula creates a new formula rule.

```go
func NewFormula(cfg FormulaConfig) (*FormulaRule, error)
```

**Parameters:**
- `cfg` (FormulaConfig)

**Returns:**
- *FormulaRule
- error

## Methods

### Dependencies



```go
func (*baseRule) Dependencies() []string
```

**Parameters:**
  None

**Returns:**
- []string

### Description



```go
func (*baseRule) Description() string
```

**Parameters:**
  None

**Returns:**
- string

### Evaluate

Evaluate computes and stores the formula result.

```go
func (*AssignmentRule) Evaluate(ctx context.Context, evalCtx *EvalContext) error
```

**Parameters:**
- `ctx` (context.Context)
- `evalCtx` (*EvalContext)

**Returns:**
- error

### Expression

Expression returns the expression string (if any).

```go
func (*FormulaRule) Expression() string
```

**Parameters:**
  None

**Returns:**
- string

### ID



```go
func (*baseRule) ID() string
```

**Parameters:**
  None

**Returns:**
- string

### Inputs

Inputs returns the required input keys.

```go
func (*FormulaRule) Inputs() []string
```

**Parameters:**
  None

**Returns:**
- []string

### Name



```go
func (*Engine) Name() string
```

**Parameters:**
  None

**Returns:**
- string

### SetFormulaFunc

SetFormulaFunc sets the formula function (used by expression parser).

```go
func (*FormulaRule) SetFormulaFunc(fn FormulaFunc)
```

**Parameters:**
- `fn` (FormulaFunc)

**Returns:**
  None

### Target

Target returns the target key for this formula.

```go
func (*FormulaRule) Target() string
```

**Parameters:**
  None

**Returns:**
- string

### Logger
Logger provides a simple logging interface.

#### Example Usage

```go
// Example implementation of Logger
type MyLogger struct {
    // Add your fields here
}

func (m MyLogger) Debug(param1 string, param2 ...any)  {
    // Implement your logic here
    return
}

func (m MyLogger) Info(param1 string, param2 ...any)  {
    // Implement your logic here
    return
}

func (m MyLogger) Warn(param1 string, param2 ...any)  {
    // Implement your logic here
    return
}

func (m MyLogger) Error(param1 string, param2 error, param3 ...any)  {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type Logger interface {
    Debug(msg string, kv ...any)
    Info(msg string, kv ...any)
    Warn(msg string, kv ...any)
    Error(msg string, err error, kv ...any)
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### Lookup
Lookup represents a lookup table.

#### Example Usage

```go
// Example implementation of Lookup
type MyLookup struct {
    // Add your fields here
}

func (m MyLookup) Name() string {
    // Implement your logic here
    return
}

func (m MyLookup) Get(param1 any) any {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type Lookup interface {
    Name() string
    Get(key any) (any, bool)
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### LookupConfig
LookupConfig configures a lookup rule.

#### Example Usage

```go
// Create a new LookupConfig
lookupconfig := LookupConfig{
    ID: "example",
    Name: "example",
    Description: "example",
    Deps: [],
    Table: "example",
    Key: "example",
    Target: "example",
    Default: any{},
    Required: true,
}
```

#### Type Definition

```go
type LookupConfig struct {
    ID string
    Name string
    Description string
    Deps []string
    Table string
    Key string
    Target string
    Default any
    Required bool
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| ID | `string` |  |
| Name | `string` |  |
| Description | `string` |  |
| Deps | `[]string` |  |
| Table | `string` | Table is the lookup table name (must be registered). |
| Key | `string` | Key is the context key to use as lookup key. |
| Target | `string` | Target is the context key to store the result. |
| Default | `any` | Default is the value if not found (ignored if Required is true). |
| Required | `bool` | Required causes an error if the lookup key is not found. |

### LookupRule
LookupRule retrieves a value from a lookup table.

#### Example Usage

```go
// Create a new LookupRule
lookuprule := LookupRule{

}
```

#### Type Definition

```go
type LookupRule struct {
}
```

### Constructor Functions

### MustLookup

MustLookup creates a new lookup rule, panicking on error.

```go
func MustLookup(cfg LookupConfig) *LookupRule
```

**Parameters:**
- `cfg` (LookupConfig)

**Returns:**
- *LookupRule

### NewLookup

NewLookup creates a new lookup rule.

```go
func NewLookup(cfg LookupConfig) (*LookupRule, error)
```

**Parameters:**
- `cfg` (LookupConfig)

**Returns:**
- *LookupRule
- error

## Methods

### Dependencies



```go
func (*baseRule) Dependencies() []string
```

**Parameters:**
  None

**Returns:**
- []string

### Description



```go
func (*baseRule) Description() string
```

**Parameters:**
  None

**Returns:**
- string

### Evaluate

Evaluate performs the lookup and sets the result.

```go
func (*Engine) Evaluate(ctx context.Context, evalCtx *EvalContext) (*Result, error)
```

**Parameters:**
- `ctx` (context.Context)
- `evalCtx` (*EvalContext)

**Returns:**
- *Result
- error

### ID



```go
func (*baseRule) ID() string
```

**Parameters:**
  None

**Returns:**
- string

### Name



```go
func (**ast.IndexExpr) Name() string
```

**Parameters:**
  None

**Returns:**
- string

### Table

Table returns the lookup table name.

```go
func (*LookupRule) Table() string
```

**Parameters:**
  None

**Returns:**
- string

### MapLookup
MapLookup provides a simple map-based lookup implementation.

#### Example Usage

```go
// Create a new MapLookup
maplookup := MapLookup{

}
```

#### Type Definition

```go
type MapLookup struct {
}
```

### Constructor Functions

### NewMapLookup

NewMapLookup creates a new map-based lookup table.

```go
func NewMapLookup(name string, items map[K]V) **ast.IndexListExpr
```

**Parameters:**
- `name` (string)
- `items` (map[K]V)

**Returns:**
- **ast.IndexListExpr

## Methods

### Get



```go
func (*EvalContext) Get(key string) (any, bool)
```

**Parameters:**
- `key` (string)

**Returns:**
- any
- bool

### Name



```go
func (**ast.IndexExpr) Name() string
```

**Parameters:**
  None

**Returns:**
- string

### Metrics
Metrics provides a simple metrics interface.

#### Example Usage

```go
// Example implementation of Metrics
type MyMetrics struct {
    // Add your fields here
}

func (m MyMetrics) Inc(param1 string, param2 ...any)  {
    // Implement your logic here
    return
}

func (m MyMetrics) Add(param1 string, param2 float64, param3 ...any)  {
    // Implement your logic here
    return
}

func (m MyMetrics) Histogram(param1 string, param2 float64, param3 ...any)  {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type Metrics interface {
    Inc(name string, kv ...any)
    Add(name string, v float64, kv ...any)
    Histogram(name string, v float64, kv ...any)
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### Observability
Observability holds observability hooks for the engine.

#### Example Usage

```go
// Create a new Observability
observability := Observability{
    Logger: Logger{},
    Metrics: Metrics{},
    Tracer: Tracer{},
}
```

#### Type Definition

```go
type Observability struct {
    Logger Logger
    Metrics Metrics
    Tracer Tracer
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Logger | `Logger` |  |
| Metrics | `Metrics` |  |
| Tracer | `Tracer` |  |

### RangeEntry
RangeEntry represents a single range in a range lookup.

#### Example Usage

```go
// Create a new RangeEntry
rangeentry := RangeEntry{
    Min: 3.14,
    Max: 3.14,
    Value: V{},
}
```

#### Type Definition

```go
type RangeEntry struct {
    Min float64
    Max float64
    Value V
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Min | `float64` | inclusive |
| Max | `float64` | exclusive (use math.Inf(1) for unbounded) |
| Value | `V` |  |

### RangeLookup
RangeLookup provides range-based lookup (e.g., tax brackets).

#### Example Usage

```go
// Create a new RangeLookup
rangelookup := RangeLookup{

}
```

#### Type Definition

```go
type RangeLookup struct {
}
```

### Constructor Functions

### NewRangeLookup

NewRangeLookup creates a new range-based lookup table.

```go
func NewRangeLookup(name string, ranges []*ast.IndexExpr) **ast.IndexExpr
```

**Parameters:**
- `name` (string)
- `ranges` ([]*ast.IndexExpr)

**Returns:**
- **ast.IndexExpr

### NewTaxBracketLookup

NewTaxBracketLookup creates a range lookup for tax brackets.

```go
func NewTaxBracketLookup(name string, brackets []TaxBracket) **ast.IndexExpr
```

**Parameters:**
- `name` (string)
- `brackets` ([]TaxBracket)

**Returns:**
- **ast.IndexExpr

## Methods

### Get



```go
func (*EvalContext) Get(key string) (any, bool)
```

**Parameters:**
- `key` (string)

**Returns:**
- any
- bool

### Name



```go
func (**ast.IndexExpr) Name() string
```

**Parameters:**
  None

**Returns:**
- string

### Result
Result contains the outcome of a rule evaluation.

#### Example Usage

```go
// Create a new Result
result := Result{
    ID: "example",
    Success: true,
    RulesEvaluated: 42,
    RulesFailed: 42,
    Errors: [],
    Duration: /* value */,
    HaltedBy: "example",
    Context: &EvalContext{}{},
}
```

#### Type Definition

```go
type Result struct {
    ID string
    Success bool
    RulesEvaluated int
    RulesFailed int
    Errors []RuleError
    Duration time.Duration
    HaltedBy string
    Context *EvalContext
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| ID | `string` | ID is the evaluation context ID. |
| Success | `bool` | Success indicates if all rules evaluated without error. |
| RulesEvaluated | `int` | RulesEvaluated is the number of rules that were run. |
| RulesFailed | `int` | RulesFailed is the number of rules that failed. |
| Errors | `[]RuleError` | Errors contains all collected errors (in CollectAll mode). |
| Duration | `time.Duration` | Duration is the total evaluation time. |
| HaltedBy | `string` | HaltedBy is the rule ID that halted evaluation (if any). |
| Context | `*EvalContext` | Context is the final evaluation context state. |

## Methods

### ErrorMessages

ErrorMessages returns all error messages.

```go
func (*Result) ErrorMessages() []string
```

**Parameters:**
  None

**Returns:**
- []string

### FirstError

FirstError returns the first error, or nil if none.

```go
func (*Result) FirstError() error
```

**Parameters:**
  None

**Returns:**
- error

### HasErrors

HasErrors returns true if any errors were collected.

```go
func (*Result) HasErrors() bool
```

**Parameters:**
  None

**Returns:**
- bool

### Rule
Rule represents any rule that can be evaluated against an EvalContext.

#### Example Usage

```go
// Example implementation of Rule
type MyRule struct {
    // Add your fields here
}

func (m MyRule) ID() string {
    // Implement your logic here
    return
}

func (m MyRule) Evaluate(param1 context.Context, param2 *EvalContext) error {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type Rule interface {
    ID() string
    Evaluate(ctx context.Context, evalCtx *EvalContext) error
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### RuleError
RuleError wraps an error with rule context.

#### Example Usage

```go
// Create a new RuleError
ruleerror := RuleError{
    RuleID: "example",
    RuleType: "example",
    Phase: "example",
    Err: error{},
}
```

#### Type Definition

```go
type RuleError struct {
    RuleID string
    RuleType string
    Phase string
    Err error
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| RuleID | `string` |  |
| RuleType | `string` |  |
| Phase | `string` | "validate", "evaluate" |
| Err | `error` |  |

### Constructor Functions

### NewRuleError

NewRuleError creates a new RuleError.

```go
func NewRuleError(ruleID, ruleType, phase string, err error) *RuleError
```

**Parameters:**
- `ruleID` (string)
- `ruleType` (string)
- `phase` (string)
- `err` (error)

**Returns:**
- *RuleError

## Methods

### Error



```go
func (nopLogger) Error(msg string, err error, kv ...any)
```

**Parameters:**
- `msg` (string)
- `err` (error)
- `kv` (...any)

**Returns:**
  None

### Unwrap



```go
func (*RuleError) Unwrap() error
```

**Parameters:**
  None

**Returns:**
- error

### RuleMetadata
RuleMetadata provides optional metadata about a rule.

#### Example Usage

```go
// Example implementation of RuleMetadata
type MyRuleMetadata struct {
    // Add your fields here
}

func (m MyRuleMetadata) Name() string {
    // Implement your logic here
    return
}

func (m MyRuleMetadata) Description() string {
    // Implement your logic here
    return
}

func (m MyRuleMetadata) Dependencies() []string {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type RuleMetadata interface {
    Rule
    Name() string
    Description() string
    Dependencies() []string
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### RuleType
RuleType identifies the type of rule.

#### Example Usage

```go
// Example usage of RuleType
var value RuleType
// Initialize with appropriate value
```

#### Type Definition

```go
type RuleType string
```

### TaxBracket
TaxBracket is a convenience type for common tax bracket lookups.

#### Example Usage

```go
// Create a new TaxBracket
taxbracket := TaxBracket{
    Min: 3.14,
    Max: 3.14,
    Rate: 3.14,
}
```

#### Type Definition

```go
type TaxBracket struct {
    Min float64
    Max float64
    Rate float64
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Min | `float64` |  |
| Max | `float64` |  |
| Rate | `float64` |  |

### Tracer
Tracer provides a simple tracing interface.

#### Example Usage

```go
// Example implementation of Tracer
type MyTracer struct {
    // Add your fields here
}

func (m MyTracer) Start(param1 context.Context, param2 string, param3 ...any) context.Context {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type Tracer interface {
    Start(ctx context.Context, name string, kv ...any) (context.Context, func(err error))
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### ValueFunc
ValueFunc computes a value dynamically based on context.

#### Example Usage

```go
// Example usage of ValueFunc
var value ValueFunc
// Initialize with appropriate value
```

#### Type Definition

```go
type ValueFunc func(ctx context.Context, evalCtx *EvalContext) (any, error)
```

## Functions

### GetTyped
GetTyped retrieves a typed value from the context.

```go
func GetTyped(e *EvalContext, key string) (T, bool)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `e` | `*EvalContext` | |
| `key` | `string` | |

**Returns:**
| Type | Description |
|------|-------------|
| `T` | |
| `bool` | |

**Example:**

```go
// Example usage of GetTyped
result := GetTyped(/* parameters */)
```

### SetTyped
SetTyped stores a typed value in the context.

```go
func SetTyped(e *EvalContext, key string, value T)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `e` | `*EvalContext` | |
| `key` | `string` | |
| `value` | `T` | |

**Returns:**
None

**Example:**

```go
// Example usage of SetTyped
result := SetTyped(/* parameters */)
```

## External Links

- [Package Overview](../packages/cortex.md)
- [pkg.go.dev Documentation](https://pkg.go.dev/github.com/kolosys/cortex)
- [Source Code](https://github.com/kolosys/cortex/tree/main/github.com/kolosys/cortex)
