package expr_test

import (
	"context"
	"testing"

	"github.com/kolosys/cortex/expr"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		input    string
		expected []expr.TokenType
	}{
		{"1 + 2", []expr.TokenType{expr.TokenNumber, expr.TokenPlus, expr.TokenNumber, expr.TokenEOF}},
		{"x * y", []expr.TokenType{expr.TokenIdent, expr.TokenStar, expr.TokenIdent, expr.TokenEOF}},
		{"a >= 10", []expr.TokenType{expr.TokenIdent, expr.TokenGe, expr.TokenNumber, expr.TokenEOF}},
		{"true && false", []expr.TokenType{expr.TokenBool, expr.TokenAnd, expr.TokenBool, expr.TokenEOF}},
		{"min(a, b)", []expr.TokenType{expr.TokenIdent, expr.TokenLParen, expr.TokenIdent, expr.TokenComma, expr.TokenIdent, expr.TokenRParen, expr.TokenEOF}},
		{"!flag", []expr.TokenType{expr.TokenNot, expr.TokenIdent, expr.TokenEOF}},
		{"x != y", []expr.TokenType{expr.TokenIdent, expr.TokenNe, expr.TokenIdent, expr.TokenEOF}},
		{"3.14", []expr.TokenType{expr.TokenNumber, expr.TokenEOF}},
		{`"hello"`, []expr.TokenType{expr.TokenString, expr.TokenEOF}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tokens := expr.Tokenize(tt.input)
			if len(tokens) != len(tt.expected) {
				t.Errorf("expected %d tokens, got %d", len(tt.expected), len(tokens))
				return
			}
			for i, tok := range tokens {
				if tok.Type != tt.expected[i] {
					t.Errorf("token %d: expected %v, got %v", i, tt.expected[i], tok.Type)
				}
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []string{
		"1 + 2",
		"x * y",
		"a + b * c",
		"(a + b) * c",
		"min(x, y)",
		"if(a > b, a, b)",
		"!flag",
		"-x",
		"a && b || c",
		"a == b",
		"x >= 10 && x <= 20",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := expr.Parse(input)
			if err != nil {
				t.Errorf("parse error: %v", err)
			}
		})
	}
}

func TestEval(t *testing.T) {
	tests := []struct {
		expr     string
		values   map[string]any
		expected any
	}{
		{"1 + 2", nil, 3.0},
		{"10 - 3", nil, 7.0},
		{"4 * 5", nil, 20.0},
		{"20 / 4", nil, 5.0},
		{"10 % 3", nil, 1.0},
		{"-5", nil, -5.0},
		{"x + y", map[string]any{"x": 10.0, "y": 20.0}, 30.0},
		{"x * y", map[string]any{"x": 3.0, "y": 4.0}, 12.0},
		{"a + b * c", map[string]any{"a": 1.0, "b": 2.0, "c": 3.0}, 7.0},
		{"(a + b) * c", map[string]any{"a": 1.0, "b": 2.0, "c": 3.0}, 9.0},
		{"true && false", nil, false},
		{"true || false", nil, true},
		{"!true", nil, false},
		{"!false", nil, true},
		{"5 == 5", nil, true},
		{"5 != 5", nil, false},
		{"10 > 5", nil, true},
		{"10 < 5", nil, false},
		{"10 >= 10", nil, true},
		{"10 <= 9", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			e, err := expr.Compile(tt.expr)
			if err != nil {
				t.Fatalf("compile error: %v", err)
			}
			result, err := e.EvalWithMap(context.Background(), tt.values)
			if err != nil {
				t.Fatalf("eval error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEvalFunctions(t *testing.T) {
	tests := []struct {
		expr     string
		values   map[string]any
		expected float64
	}{
		{"min(10, 5)", nil, 5.0},
		{"max(10, 5)", nil, 10.0},
		{"min(a, b, c)", map[string]any{"a": 5.0, "b": 3.0, "c": 8.0}, 3.0},
		{"max(a, b, c)", map[string]any{"a": 5.0, "b": 3.0, "c": 8.0}, 8.0},
		{"abs(-5)", nil, 5.0},
		{"abs(5)", nil, 5.0},
		{"floor(3.7)", nil, 3.0},
		{"ceil(3.2)", nil, 4.0},
		{"round(3.5)", nil, 4.0},
		{"round(3.14159, 2)", nil, 3.14},
		{"sqrt(16)", nil, 4.0},
		{"pow(2, 3)", nil, 8.0},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			e, err := expr.Compile(tt.expr)
			if err != nil {
				t.Fatalf("compile error: %v", err)
			}
			result, err := e.EvalFloat64(context.Background(), mapGetter(tt.values))
			if err != nil {
				t.Fatalf("eval error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEvalIf(t *testing.T) {
	tests := []struct {
		expr     string
		values   map[string]any
		expected any
	}{
		{"if(true, 1, 2)", nil, 1.0},
		{"if(false, 1, 2)", nil, 2.0},
		{"if(x > 10, x, 10)", map[string]any{"x": 15.0}, 15.0},
		{"if(x > 10, x, 10)", map[string]any{"x": 5.0}, 10.0},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			e, err := expr.Compile(tt.expr)
			if err != nil {
				t.Fatalf("compile error: %v", err)
			}
			result, err := e.EvalWithMap(context.Background(), tt.values)
			if err != nil {
				t.Fatalf("eval error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEvalErrors(t *testing.T) {
	tests := []struct {
		name   string
		expr   string
		values map[string]any
	}{
		{"undefined variable", "x + 1", nil},
		{"division by zero", "10 / 0", nil},
		{"undefined function", "foo(1)", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := expr.Compile(tt.expr)
			if err != nil {
				return // compile error is also acceptable
			}
			_, err = e.EvalWithMap(context.Background(), tt.values)
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	tests := []string{
		"1 +",
		"(1 + 2",
		"1 + + 2",
		"",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := expr.Parse(input)
			if err == nil {
				t.Error("expected parse error")
			}
		})
	}
}

func TestCustomFunction(t *testing.T) {
	e := expr.MustCompile("double(x)")
	e.RegisterFunc("double", func(args ...any) (any, error) {
		f, _ := args[0].(float64)
		return f * 2, nil
	})

	result, err := e.EvalWithMap(context.Background(), map[string]any{"x": 5.0})
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	if result != 10.0 {
		t.Errorf("expected 10, got %v", result)
	}
}

func TestStringConcat(t *testing.T) {
	e := expr.MustCompile(`"hello" + " " + "world"`)
	result, err := e.EvalWithMap(context.Background(), nil)
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	if result != "hello world" {
		t.Errorf("expected 'hello world', got %v", result)
	}
}

type mapGetter map[string]any

func (m mapGetter) Get(key string) (any, bool) {
	v, ok := m[key]
	return v, ok
}
