# expr API

Complete API documentation for the expr package.

**Import Path:** `github.com/kolosys/cortex/expr`

## Package Documentation

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


## Types

### BinaryExpr
BinaryExpr represents a binary expression.

#### Example Usage

```go
// Create a new BinaryExpr
binaryexpr := BinaryExpr{
    Op: TokenType{},
    Left: Node{},
    Right: Node{},
}
```

#### Type Definition

```go
type BinaryExpr struct {
    Op TokenType
    Left Node
    Right Node
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Op | `TokenType` |  |
| Left | `Node` |  |
| Right | `Node` |  |

### BoolLit
BoolLit represents a boolean literal.

#### Example Usage

```go
// Create a new BoolLit
boollit := BoolLit{
    Value: true,
}
```

#### Type Definition

```go
type BoolLit struct {
    Value bool
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Value | `bool` |  |

### CallExpr
CallExpr represents a function call.

#### Example Usage

```go
// Create a new CallExpr
callexpr := CallExpr{
    Name: "example",
    Args: [],
}
```

#### Type Definition

```go
type CallExpr struct {
    Name string
    Args []Node
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Name | `string` |  |
| Args | `[]Node` |  |

### Evaluator
Evaluator evaluates an AST against a value getter.

#### Example Usage

```go
// Create a new Evaluator
evaluator := Evaluator{

}
```

#### Type Definition

```go
type Evaluator struct {
}
```

### Constructor Functions

### NewEvaluator

NewEvaluator creates a new evaluator with built-in functions.

```go
func NewEvaluator() *Evaluator
```

**Parameters:**
  None

**Returns:**
- *Evaluator

## Methods

### Eval

Eval evaluates an AST node against a value getter.

```go
func (*Evaluator) Eval(ctx context.Context, node Node, getter ValueGetter) (any, error)
```

**Parameters:**
- `ctx` (context.Context)
- `node` (Node)
- `getter` (ValueGetter)

**Returns:**
- any
- error

### RegisterFunc

RegisterFunc registers a custom function.

```go
func (*Evaluator) RegisterFunc(name string, fn Func)
```

**Parameters:**
- `name` (string)
- `fn` (Func)

**Returns:**
  None

### Expression
Expression represents a compiled expression.

#### Example Usage

```go
// Create a new Expression
expression := Expression{

}
```

#### Type Definition

```go
type Expression struct {
}
```

### Constructor Functions

### Compile

Compile parses and compiles an expression string.

```go
func Compile(input string) (*Expression, error)
```

**Parameters:**
- `input` (string)

**Returns:**
- *Expression
- error

### MustCompile

MustCompile compiles an expression, panicking on error.

```go
func MustCompile(input string) *Expression
```

**Parameters:**
- `input` (string)

**Returns:**
- *Expression

## Methods

### Eval

Eval evaluates the expression against a value getter.

```go
func (*Expression) Eval(ctx context.Context, getter ValueGetter) (any, error)
```

**Parameters:**
- `ctx` (context.Context)
- `getter` (ValueGetter)

**Returns:**
- any
- error

### EvalBool

EvalBool evaluates the expression and returns a bool.

```go
func (*Expression) EvalBool(ctx context.Context, getter ValueGetter) (bool, error)
```

**Parameters:**
- `ctx` (context.Context)
- `getter` (ValueGetter)

**Returns:**
- bool
- error

### EvalFloat64

EvalFloat64 evaluates the expression and returns a float64.

```go
func (*Expression) EvalFloat64(ctx context.Context, getter ValueGetter) (float64, error)
```

**Parameters:**
- `ctx` (context.Context)
- `getter` (ValueGetter)

**Returns:**
- float64
- error

### EvalWithMap

EvalWithMap evaluates the expression using a map as the value source.

```go
func (*Expression) EvalWithMap(ctx context.Context, values map[string]any) (any, error)
```

**Parameters:**
- `ctx` (context.Context)
- `values` (map[string]any)

**Returns:**
- any
- error

### Raw

Raw returns the original expression string.

```go
func (*Expression) Raw() string
```

**Parameters:**
  None

**Returns:**
- string

### RegisterFunc

RegisterFunc registers a custom function for this expression.

```go
func (*Evaluator) RegisterFunc(name string, fn Func)
```

**Parameters:**
- `name` (string)
- `fn` (Func)

**Returns:**
  None

### Func
Func is a built-in function type.

#### Example Usage

```go
// Example usage of Func
var value Func
// Initialize with appropriate value
```

#### Type Definition

```go
type Func func(args ...any) (any, error)
```

### Ident
Ident represents an identifier (variable reference).

#### Example Usage

```go
// Create a new Ident
ident := Ident{
    Name: "example",
}
```

#### Type Definition

```go
type Ident struct {
    Name string
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Name | `string` |  |

### Lexer
Lexer tokenizes an expression string.

#### Example Usage

```go
// Create a new Lexer
lexer := Lexer{

}
```

#### Type Definition

```go
type Lexer struct {
}
```

### Constructor Functions

### NewLexer

NewLexer creates a new lexer for the given input.

```go
func NewLexer(input string) *Lexer
```

**Parameters:**
- `input` (string)

**Returns:**
- *Lexer

## Methods

### NextToken

NextToken returns the next token from the input.

```go
func (*Lexer) NextToken() Token
```

**Parameters:**
  None

**Returns:**
- Token

### Node
Node represents an AST node.

#### Example Usage

```go
// Example implementation of Node
type MyNode struct {
    // Add your fields here
}


```

#### Type Definition

```go
type Node interface {
}
```

## Methods

| Method | Description |
| ------ | ----------- |

### Constructor Functions

### Parse

Parse parses an expression string into an AST.

```go
func Parse(input string) (Node, error)
```

**Parameters:**
- `input` (string)

**Returns:**
- Node
- error

### NumberLit
NumberLit represents a numeric literal.

#### Example Usage

```go
// Create a new NumberLit
numberlit := NumberLit{
    Value: 3.14,
}
```

#### Type Definition

```go
type NumberLit struct {
    Value float64
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Value | `float64` |  |

### Parser
Parser parses an expression into an AST.

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

NewParser creates a new parser for the given input.

```go
func NewParser(input string) *Parser
```

**Parameters:**
- `input` (string)

**Returns:**
- *Parser

## Methods

### Errors

Errors returns any parsing errors.

```go
func (*Parser) Errors() []string
```

**Parameters:**
  None

**Returns:**
- []string

### Parse

Parse parses the expression and returns the AST root.

```go
func Parse(input string) (Node, error)
```

**Parameters:**
- `input` (string)

**Returns:**
- Node
- error

### StringLit
StringLit represents a string literal.

#### Example Usage

```go
// Create a new StringLit
stringlit := StringLit{
    Value: "example",
}
```

#### Type Definition

```go
type StringLit struct {
    Value string
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Value | `string` |  |

### Token
Token represents a lexical token.

#### Example Usage

```go
// Create a new Token
token := Token{
    Type: TokenType{},
    Literal: "example",
    Pos: 42,
}
```

#### Type Definition

```go
type Token struct {
    Type TokenType
    Literal string
    Pos int
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Type | `TokenType` |  |
| Literal | `string` |  |
| Pos | `int` |  |

### Constructor Functions

### Tokenize

Tokenize returns all tokens from the input.

```go
func Tokenize(input string) []Token
```

**Parameters:**
- `input` (string)

**Returns:**
- []Token

### TokenType
TokenType represents the type of a token.

#### Example Usage

```go
// Example usage of TokenType
var value TokenType
// Initialize with appropriate value
```

#### Type Definition

```go
type TokenType int
```

## Methods

### String



```go
func (TokenType) String() string
```

**Parameters:**
  None

**Returns:**
- string

### UnaryExpr
UnaryExpr represents a unary expression.

#### Example Usage

```go
// Create a new UnaryExpr
unaryexpr := UnaryExpr{
    Op: TokenType{},
    Expr: Node{},
}
```

#### Type Definition

```go
type UnaryExpr struct {
    Op TokenType
    Expr Node
}
```

### Fields

| Field | Type | Description |
| ----- | ---- | ----------- |
| Op | `TokenType` |  |
| Expr | `Node` |  |

### ValueGetter
ValueGetter retrieves values by name (e.g., from EvalContext).

#### Example Usage

```go
// Example implementation of ValueGetter
type MyValueGetter struct {
    // Add your fields here
}

func (m MyValueGetter) Get(param1 string) any {
    // Implement your logic here
    return
}


```

#### Type Definition

```go
type ValueGetter interface {
    Get(key string) (any, bool)
}
```

## Methods

| Method | Description |
| ------ | ----------- |

## External Links

- [Package Overview](../packages/expr.md)
- [pkg.go.dev Documentation](https://pkg.go.dev/github.com/kolosys/cortex/expr)
- [Source Code](https://github.com/kolosys/cortex/tree/main/expr)
