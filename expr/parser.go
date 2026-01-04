package expr

import (
	"fmt"
	"strconv"
)

// Parser parses an expression into an AST.
type Parser struct {
	lexer   *Lexer
	current Token
	peek    Token
	errors  []string
}

// NewParser creates a new parser for the given input.
func NewParser(input string) *Parser {
	p := &Parser{lexer: NewLexer(input)}
	p.advance()
	p.advance()
	return p
}

func (p *Parser) advance() {
	p.current = p.peek
	p.peek = p.lexer.NextToken()
}

func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, msg)
}

// Errors returns any parsing errors.
func (p *Parser) Errors() []string {
	return p.errors
}

// Parse parses the expression and returns the AST root.
func (p *Parser) Parse() (Node, error) {
	node := p.parseExpression(0)

	if p.current.Type != TokenEOF {
		return nil, fmt.Errorf("unexpected token %s at position %d", p.current.Type, p.current.Pos)
	}

	if len(p.errors) > 0 {
		return nil, fmt.Errorf("parse errors: %v", p.errors)
	}

	return node, nil
}

// Precedence levels
const (
	precLowest  = 0
	precOr      = 1
	precAnd     = 2
	precCompare = 3
	precSum     = 4
	precProduct = 5
)

func precedence(t TokenType) int {
	switch t {
	case TokenOr:
		return precOr
	case TokenAnd:
		return precAnd
	case TokenEq, TokenNe, TokenLt, TokenLe, TokenGt, TokenGe:
		return precCompare
	case TokenPlus, TokenMinus:
		return precSum
	case TokenStar, TokenSlash, TokenPercent:
		return precProduct
	default:
		return precLowest
	}
}

func (p *Parser) parseExpression(prec int) Node {
	left := p.parseUnary()

	for prec < precedence(p.current.Type) {
		op := p.current.Type
		opPrec := precedence(op)
		p.advance()
		right := p.parseExpression(opPrec)
		left = &BinaryExpr{Op: op, Left: left, Right: right}
	}

	return left
}

func (p *Parser) parseUnary() Node {
	if p.current.Type == TokenNot || p.current.Type == TokenMinus {
		op := p.current.Type
		p.advance()
		return &UnaryExpr{Op: op, Expr: p.parseUnary()}
	}
	return p.parsePrimary()
}

func (p *Parser) parsePrimary() Node {
	switch p.current.Type {
	case TokenNumber:
		val, err := strconv.ParseFloat(p.current.Literal, 64)
		if err != nil {
			p.addError(fmt.Sprintf("invalid number: %s", p.current.Literal))
			return nil
		}
		p.advance()
		return &NumberLit{Value: val}

	case TokenString:
		val := p.current.Literal
		p.advance()
		return &StringLit{Value: val}

	case TokenBool:
		val := p.current.Literal == "true"
		p.advance()
		return &BoolLit{Value: val}

	case TokenIdent:
		name := p.current.Literal
		p.advance()

		// Check if it's a function call
		if p.current.Type == TokenLParen {
			return p.parseCall(name)
		}

		return &Ident{Name: name}

	case TokenLParen:
		p.advance()
		node := p.parseExpression(precLowest)
		if p.current.Type != TokenRParen {
			p.addError("expected ')'")
			return nil
		}
		p.advance()
		return node

	default:
		p.addError(fmt.Sprintf("unexpected token: %s", p.current.Type))
		return nil
	}
}

func (p *Parser) parseCall(name string) Node {
	p.advance() // consume '('

	var args []Node

	if p.current.Type != TokenRParen {
		args = append(args, p.parseExpression(precLowest))

		for p.current.Type == TokenComma {
			p.advance()
			args = append(args, p.parseExpression(precLowest))
		}
	}

	if p.current.Type != TokenRParen {
		p.addError("expected ')'")
		return nil
	}
	p.advance()

	return &CallExpr{Name: name, Args: args}
}

// Parse parses an expression string into an AST.
func Parse(input string) (Node, error) {
	p := NewParser(input)
	return p.Parse()
}
