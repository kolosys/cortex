package expr

import (
	"unicode"
)

// Lexer tokenizes an expression string.
type Lexer struct {
	input string
	pos   int
	ch    rune
}

// NewLexer creates a new lexer for the given input.
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.pos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = rune(l.input[l.pos])
	}
	l.pos++
}

func (l *Lexer) peekChar() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	return rune(l.input[l.pos])
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	pos := l.pos - 1
	var tok Token

	switch l.ch {
	case 0:
		tok = Token{Type: TokenEOF, Pos: pos}
	case '+':
		tok = Token{Type: TokenPlus, Literal: "+", Pos: pos}
		l.readChar()
	case '-':
		tok = Token{Type: TokenMinus, Literal: "-", Pos: pos}
		l.readChar()
	case '*':
		tok = Token{Type: TokenStar, Literal: "*", Pos: pos}
		l.readChar()
	case '/':
		tok = Token{Type: TokenSlash, Literal: "/", Pos: pos}
		l.readChar()
	case '%':
		tok = Token{Type: TokenPercent, Literal: "%", Pos: pos}
		l.readChar()
	case '(':
		tok = Token{Type: TokenLParen, Literal: "(", Pos: pos}
		l.readChar()
	case ')':
		tok = Token{Type: TokenRParen, Literal: ")", Pos: pos}
		l.readChar()
	case ',':
		tok = Token{Type: TokenComma, Literal: ",", Pos: pos}
		l.readChar()
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TokenEq, Literal: "==", Pos: pos}
		} else {
			tok = Token{Type: TokenError, Literal: "unexpected '='", Pos: pos}
		}
		l.readChar()
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TokenNe, Literal: "!=", Pos: pos}
			l.readChar()
		} else {
			tok = Token{Type: TokenNot, Literal: "!", Pos: pos}
			l.readChar()
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TokenLe, Literal: "<=", Pos: pos}
		} else {
			tok = Token{Type: TokenLt, Literal: "<", Pos: pos}
		}
		l.readChar()
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TokenGe, Literal: ">=", Pos: pos}
		} else {
			tok = Token{Type: TokenGt, Literal: ">", Pos: pos}
		}
		l.readChar()
	case '&':
		if l.peekChar() == '&' {
			l.readChar()
			tok = Token{Type: TokenAnd, Literal: "&&", Pos: pos}
		} else {
			tok = Token{Type: TokenError, Literal: "unexpected '&'", Pos: pos}
		}
		l.readChar()
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			tok = Token{Type: TokenOr, Literal: "||", Pos: pos}
		} else {
			tok = Token{Type: TokenError, Literal: "unexpected '|'", Pos: pos}
		}
		l.readChar()
	case '"', '\'':
		tok = l.readString(l.ch)
		tok.Pos = pos
	default:
		if isDigit(l.ch) {
			tok = l.readNumber()
			tok.Pos = pos
			return tok
		} else if isLetter(l.ch) {
			tok = l.readIdentifier()
			tok.Pos = pos
			return tok
		} else {
			tok = Token{Type: TokenError, Literal: string(l.ch), Pos: pos}
			l.readChar()
		}
	}

	return tok
}

func (l *Lexer) readNumber() Token {
	start := l.pos - 1
	hasDot := false

	for isDigit(l.ch) || (l.ch == '.' && !hasDot) {
		if l.ch == '.' {
			hasDot = true
		}
		l.readChar()
	}

	return Token{Type: TokenNumber, Literal: l.input[start : l.pos-1]}
}

func (l *Lexer) readIdentifier() Token {
	start := l.pos - 1

	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}

	literal := l.input[start : l.pos-1]

	// Check for boolean literals
	if literal == "true" || literal == "false" {
		return Token{Type: TokenBool, Literal: literal}
	}

	return Token{Type: TokenIdent, Literal: literal}
}

func (l *Lexer) readString(quote rune) Token {
	l.readChar() // consume opening quote
	start := l.pos - 1

	for l.ch != quote && l.ch != 0 {
		l.readChar()
	}

	if l.ch == 0 {
		return Token{Type: TokenError, Literal: "unterminated string"}
	}

	literal := l.input[start : l.pos-1]
	l.readChar() // consume closing quote

	return Token{Type: TokenString, Literal: literal}
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

// Tokenize returns all tokens from the input.
func Tokenize(input string) []Token {
	l := NewLexer(input)
	var tokens []Token

	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == TokenEOF || tok.Type == TokenError {
			break
		}
	}

	return tokens
}
