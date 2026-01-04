package expr

// TokenType represents the type of a token.
type TokenType int

const (
	TokenEOF TokenType = iota
	TokenError

	// Literals
	TokenNumber
	TokenString
	TokenIdent
	TokenBool

	// Operators
	TokenPlus      // +
	TokenMinus     // -
	TokenStar      // *
	TokenSlash     // /
	TokenPercent   // %
	TokenEq        // ==
	TokenNe        // !=
	TokenLt        // <
	TokenLe        // <=
	TokenGt        // >
	TokenGe        // >=
	TokenAnd       // &&
	TokenOr        // ||
	TokenNot       // !

	// Delimiters
	TokenLParen // (
	TokenRParen // )
	TokenComma  // ,
)

func (t TokenType) String() string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenError:
		return "ERROR"
	case TokenNumber:
		return "NUMBER"
	case TokenString:
		return "STRING"
	case TokenIdent:
		return "IDENT"
	case TokenBool:
		return "BOOL"
	case TokenPlus:
		return "+"
	case TokenMinus:
		return "-"
	case TokenStar:
		return "*"
	case TokenSlash:
		return "/"
	case TokenPercent:
		return "%"
	case TokenEq:
		return "=="
	case TokenNe:
		return "!="
	case TokenLt:
		return "<"
	case TokenLe:
		return "<="
	case TokenGt:
		return ">"
	case TokenGe:
		return ">="
	case TokenAnd:
		return "&&"
	case TokenOr:
		return "||"
	case TokenNot:
		return "!"
	case TokenLParen:
		return "("
	case TokenRParen:
		return ")"
	case TokenComma:
		return ","
	default:
		return "UNKNOWN"
	}
}

// Token represents a lexical token.
type Token struct {
	Type    TokenType
	Literal string
	Pos     int
}
