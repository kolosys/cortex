package expr

// Node represents an AST node.
type Node interface {
	node()
}

// NumberLit represents a numeric literal.
type NumberLit struct {
	Value float64
}

func (*NumberLit) node() {}

// StringLit represents a string literal.
type StringLit struct {
	Value string
}

func (*StringLit) node() {}

// BoolLit represents a boolean literal.
type BoolLit struct {
	Value bool
}

func (*BoolLit) node() {}

// Ident represents an identifier (variable reference).
type Ident struct {
	Name string
}

func (*Ident) node() {}

// BinaryExpr represents a binary expression.
type BinaryExpr struct {
	Op    TokenType
	Left  Node
	Right Node
}

func (*BinaryExpr) node() {}

// UnaryExpr represents a unary expression.
type UnaryExpr struct {
	Op   TokenType
	Expr Node
}

func (*UnaryExpr) node() {}

// CallExpr represents a function call.
type CallExpr struct {
	Name string
	Args []Node
}

func (*CallExpr) node() {}
