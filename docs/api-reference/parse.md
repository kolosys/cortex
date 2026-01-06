# parse API

Complete API documentation for the parse package.

**Import Path:** `github.com/kolosys/cortex/parse`

## Package Documentation

Package parse provides JSON parsing for cortex rule definitions.


## Types

### AllocationDef
AllocationDef is the config structure for allocation rules.

#### Example Usage

```go
// Create a new AllocationDef
allocationdef := AllocationDef{
    Source: "example",
    Strategy: "example",
    Targets: [],
    Remainder: "example",
    Precision: 42,
}
```

#### Type Definition

```go
type AllocationDef struct {
    Source string `json:"source"`
    Strategy string `json:"strategy"`
    Targets []AllocationTarget `json:"targets"`
    Remainder string `json:"remainder,omitempty"`
    Precision int `json:"precision,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Source | `string` |  |
| Strategy | `string` |  |
| Targets | `[]AllocationTarget` |  |
| Remainder | `string` |  |
| Precision | `int` |  |

### AllocationTarget
AllocationTarget defines an allocation destination.

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
    Key string `json:"key"`
    Amount float64 `json:"amount"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Key | `string` |  |
| Amount | `float64` |  |

### AssignmentDef
AssignmentDef is the config structure for assignment rules.

#### Example Usage

```go
// Create a new AssignmentDef
assignmentdef := AssignmentDef{
    Target: "example",
    Value: any{},
}
```

#### Type Definition

```go
type AssignmentDef struct {
    Target string `json:"target"`
    Value any `json:"value"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Target | `string` |  |
| Value | `any` |  |

### BuildupDef
BuildupDef is the config structure for buildup rules.

#### Example Usage

```go
// Create a new BuildupDef
buildupdef := BuildupDef{
    Buildup: "example",
    Operation: "example",
    Source: "example",
    Initial: 3.14,
    Target: "example",
}
```

#### Type Definition

```go
type BuildupDef struct {
    Buildup string `json:"buildup"`
    Operation string `json:"operation"`
    Source string `json:"source,omitempty"`
    Initial float64 `json:"initial,omitempty"`
    Target string `json:"target,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Buildup | `string` |  |
| Operation | `string` |  |
| Source | `string` |  |
| Initial | `float64` |  |
| Target | `string` |  |

### FormulaDef
FormulaDef is the config structure for formula rules.

#### Example Usage

```go
// Create a new FormulaDef
formuladef := FormulaDef{
    Target: "example",
    Expression: "example",
    Function: "example",
    Inputs: [],
}
```

#### Type Definition

```go
type FormulaDef struct {
    Target string `json:"target"`
    Expression string `json:"expression,omitempty"`
    Function string `json:"function,omitempty"`
    Inputs []string `json:"inputs,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Target | `string` |  |
| Expression | `string` |  |
| Function | `string` | named registered function |
| Inputs | `[]string` |  |

### LookupDef
LookupDef defines a lookup table in config.

#### Example Usage

```go
// Create a new LookupDef
lookupdef := LookupDef{
    Name: "example",
    Type: "example",
    Entries: [],
    Items: map[],
}
```

#### Type Definition

```go
type LookupDef struct {
    Name string `json:"name"`
    Type string `json:"type"`
    Entries []LookupEntry `json:"entries,omitempty"`
    Items map[string]any `json:"items,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Name | `string` |  |
| Type | `string` | "map" or "range" |
| Entries | `[]LookupEntry` |  |
| Items | `map[string]any` | for map type |

### LookupEntry
LookupEntry defines a single entry in a range lookup.

#### Example Usage

```go
// Create a new LookupEntry
lookupentry := LookupEntry{
    Min: 3.14,
    Max: &3.14{},
    Value: any{},
}
```

#### Type Definition

```go
type LookupEntry struct {
    Min float64 `json:"min"`
    Max *float64 `json:"max"`
    Value any `json:"value"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Min | `float64` |  |
| Max | `*float64` | nil means +infinity |
| Value | `any` |  |

### LookupRuleDef
LookupRuleDef is the config structure for lookup rules.

#### Example Usage

```go
// Create a new LookupRuleDef
lookupruledef := LookupRuleDef{
    Table: "example",
    Key: "example",
    Target: "example",
    Default: any{},
    Required: true,
}
```

#### Type Definition

```go
type LookupRuleDef struct {
    Table string `json:"table"`
    Key string `json:"key"`
    Target string `json:"target"`
    Default any `json:"default,omitempty"`
    Required bool `json:"required,omitempty"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Table | `string` |  |
| Key | `string` |  |
| Target | `string` |  |
| Default | `any` |  |
| Required | `bool` |  |

### Parser
Parser parses config into rules.

#### Example Usage

```go
// Create a new Parser
parser := Parser{

}
```

#### Type Definition

```go
type Parser struct {
}
```

### Constructor Functions

### NewParser

NewParser creates a new parser.

```go
func NewParser() *Parser
```

**Parameters:**
  None

**Returns:**
- *Parser

## Methods

### ParseAndBuildEngine

ParseAndBuildEngine parses JSON and builds an Engine.

```go
func (*Parser) ParseAndBuildEngine(name string, data []byte, config *cortex.Config) (*cortex.Engine, error)
```

**Parameters:**
- `name` (string)
- `data` ([]byte)
- `config` (*cortex.Config)

**Returns:**
- *cortex.Engine
- error

### ParseJSON

ParseJSON parses a rule set from JSON.

```go
func (*Parser) ParseJSON(data []byte) (*RuleSet, error)
```

**Parameters:**
- `data` ([]byte)

**Returns:**
- *RuleSet
- error

### RegisterFormula

RegisterFormula registers a named formula function for use in config.

```go
func (*Parser) RegisterFormula(name string, fn cortex.FormulaFunc)
```

**Parameters:**
- `name` (string)
- `fn` (cortex.FormulaFunc)

**Returns:**
  None

### ToLookups

ToLookups converts lookup definitions to Lookup instances.

```go
func (*Parser) ToLookups(rs *RuleSet) ([]cortex.Lookup, error)
```

**Parameters:**
- `rs` (*RuleSet)

**Returns:**
- []cortex.Lookup
- error

### ToRules

ToRules converts rule definitions to Rule instances.

```go
func (*Parser) ToRules(rs *RuleSet) ([]cortex.Rule, error)
```

**Parameters:**
- `rs` (*RuleSet)

**Returns:**
- []cortex.Rule
- error

### RuleDefinition
RuleDefinition is a config-driven rule.

#### Example Usage

```go
// Create a new RuleDefinition
ruledefinition := RuleDefinition{
    ID: "example",
    Type: "example",
    Name: "example",
    Description: "example",
    Deps: [],
    Disabled: true,
    Config: map[],
}
```

#### Type Definition

```go
type RuleDefinition struct {
    ID string `json:"id"`
    Type string `json:"type"`
    Name string `json:"name,omitempty"`
    Description string `json:"description,omitempty"`
    Deps []string `json:"deps,omitempty"`
    Disabled bool `json:"disabled,omitempty"`
    Config map[string]any `json:"config"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| ID | `string` |  |
| Type | `string` | assignment, formula, allocation, lookup, buildup |
| Name | `string` |  |
| Description | `string` |  |
| Deps | `[]string` |  |
| Disabled | `bool` |  |
| Config | `map[string]any` |  |

### RuleSet
RuleSet represents a collection of rules from config.

#### Example Usage

```go
// Create a new RuleSet
ruleset := RuleSet{
    Version: "example",
    Name: "example",
    Lookups: [],
    Rules: [],
}
```

#### Type Definition

```go
type RuleSet struct {
    Version string `json:"version"`
    Name string `json:"name"`
    Lookups []LookupDef `json:"lookups,omitempty"`
    Rules []RuleDefinition `json:"rules"`
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Version | `string` |  |
| Name | `string` |  |
| Lookups | `[]LookupDef` |  |
| Rules | `[]RuleDefinition` |  |

## Functions

### ParseAndBuild
ParseAndBuild is a convenience function that parses JSON and builds an Engine.

```go
func ParseAndBuild(name string, data []byte, config *cortex.Config) (*cortex.Engine, error)
```

**Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `name` | `string` | |
| `data` | `[]byte` | |
| `config` | `*cortex.Config` | |

**Returns:**
| Type | Description |
|------|-------------|
| `*cortex.Engine` | |
| `error` | |

**Example:**

```go
// Example usage of ParseAndBuild
result := ParseAndBuild(/* parameters */)
```

## External Links

- [Package Overview](../packages/parse.md)
- [pkg.go.dev Documentation](https://pkg.go.dev/github.com/kolosys/cortex/parse)
- [Source Code](https://github.com/kolosys/cortex/tree/main/parse)
