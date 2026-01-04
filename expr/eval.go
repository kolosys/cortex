package expr

import (
	"context"
	"fmt"
	"math"
)

// ValueGetter retrieves values by name (e.g., from EvalContext).
type ValueGetter interface {
	Get(key string) (any, bool)
}

// Evaluator evaluates an AST against a value getter.
type Evaluator struct {
	funcs map[string]Func
}

// Func is a built-in function type.
type Func func(args ...any) (any, error)

// NewEvaluator creates a new evaluator with built-in functions.
func NewEvaluator() *Evaluator {
	e := &Evaluator{
		funcs: make(map[string]Func),
	}
	e.registerBuiltins()
	return e
}

func (e *Evaluator) registerBuiltins() {
	e.funcs["min"] = funcMin
	e.funcs["max"] = funcMax
	e.funcs["abs"] = funcAbs
	e.funcs["floor"] = funcFloor
	e.funcs["ceil"] = funcCeil
	e.funcs["round"] = funcRound
	e.funcs["if"] = funcIf
	e.funcs["sqrt"] = funcSqrt
	e.funcs["pow"] = funcPow
}

// RegisterFunc registers a custom function.
func (e *Evaluator) RegisterFunc(name string, fn Func) {
	e.funcs[name] = fn
}

// Eval evaluates an AST node against a value getter.
func (e *Evaluator) Eval(ctx context.Context, node Node, getter ValueGetter) (any, error) {
	return e.eval(ctx, node, getter)
}

func (e *Evaluator) eval(ctx context.Context, node Node, getter ValueGetter) (any, error) {
	switch n := node.(type) {
	case *NumberLit:
		return n.Value, nil

	case *StringLit:
		return n.Value, nil

	case *BoolLit:
		return n.Value, nil

	case *Ident:
		val, ok := getter.Get(n.Name)
		if !ok {
			return nil, fmt.Errorf("undefined variable: %s", n.Name)
		}
		return val, nil

	case *UnaryExpr:
		val, err := e.eval(ctx, n.Expr, getter)
		if err != nil {
			return nil, err
		}
		return e.evalUnary(n.Op, val)

	case *BinaryExpr:
		left, err := e.eval(ctx, n.Left, getter)
		if err != nil {
			return nil, err
		}
		right, err := e.eval(ctx, n.Right, getter)
		if err != nil {
			return nil, err
		}
		return e.evalBinary(n.Op, left, right)

	case *CallExpr:
		fn, ok := e.funcs[n.Name]
		if !ok {
			return nil, fmt.Errorf("undefined function: %s", n.Name)
		}
		args := make([]any, len(n.Args))
		for i, arg := range n.Args {
			val, err := e.eval(ctx, arg, getter)
			if err != nil {
				return nil, err
			}
			args[i] = val
		}
		return fn(args...)

	default:
		return nil, fmt.Errorf("unknown node type: %T", node)
	}
}

func (e *Evaluator) evalUnary(op TokenType, val any) (any, error) {
	switch op {
	case TokenNot:
		b, ok := val.(bool)
		if !ok {
			return nil, fmt.Errorf("expected bool for !, got %T", val)
		}
		return !b, nil

	case TokenMinus:
		f, err := toFloat(val)
		if err != nil {
			return nil, err
		}
		return -f, nil

	default:
		return nil, fmt.Errorf("unknown unary operator: %s", op)
	}
}

func (e *Evaluator) evalBinary(op TokenType, left, right any) (any, error) {
	switch op {
	case TokenAnd:
		l, lok := left.(bool)
		r, rok := right.(bool)
		if !lok || !rok {
			return nil, fmt.Errorf("expected bool for &&")
		}
		return l && r, nil

	case TokenOr:
		l, lok := left.(bool)
		r, rok := right.(bool)
		if !lok || !rok {
			return nil, fmt.Errorf("expected bool for ||")
		}
		return l || r, nil

	case TokenEq:
		return equals(left, right), nil

	case TokenNe:
		return !equals(left, right), nil

	case TokenLt, TokenLe, TokenGt, TokenGe:
		lf, err := toFloat(left)
		if err != nil {
			return nil, err
		}
		rf, err := toFloat(right)
		if err != nil {
			return nil, err
		}
		switch op {
		case TokenLt:
			return lf < rf, nil
		case TokenLe:
			return lf <= rf, nil
		case TokenGt:
			return lf > rf, nil
		case TokenGe:
			return lf >= rf, nil
		}

	case TokenPlus:
		// Handle string concatenation
		if ls, ok := left.(string); ok {
			if rs, ok := right.(string); ok {
				return ls + rs, nil
			}
		}
		lf, err := toFloat(left)
		if err != nil {
			return nil, err
		}
		rf, err := toFloat(right)
		if err != nil {
			return nil, err
		}
		return lf + rf, nil

	case TokenMinus:
		lf, err := toFloat(left)
		if err != nil {
			return nil, err
		}
		rf, err := toFloat(right)
		if err != nil {
			return nil, err
		}
		return lf - rf, nil

	case TokenStar:
		lf, err := toFloat(left)
		if err != nil {
			return nil, err
		}
		rf, err := toFloat(right)
		if err != nil {
			return nil, err
		}
		return lf * rf, nil

	case TokenSlash:
		lf, err := toFloat(left)
		if err != nil {
			return nil, err
		}
		rf, err := toFloat(right)
		if err != nil {
			return nil, err
		}
		if rf == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return lf / rf, nil

	case TokenPercent:
		lf, err := toFloat(left)
		if err != nil {
			return nil, err
		}
		rf, err := toFloat(right)
		if err != nil {
			return nil, err
		}
		if rf == 0 {
			return nil, fmt.Errorf("modulo by zero")
		}
		return math.Mod(lf, rf), nil
	}

	return nil, fmt.Errorf("unknown binary operator: %s", op)
}

func toFloat(v any) (float64, error) {
	switch n := v.(type) {
	case float64:
		return n, nil
	case float32:
		return float64(n), nil
	case int:
		return float64(n), nil
	case int64:
		return float64(n), nil
	case int32:
		return float64(n), nil
	default:
		return 0, fmt.Errorf("expected number, got %T", v)
	}
}

func equals(a, b any) bool {
	af, aok := toFloat(a)
	bf, bok := toFloat(b)
	if aok == nil && bok == nil {
		return af == bf
	}
	return a == b
}

// Built-in functions

func funcMin(args ...any) (any, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("min requires at least 2 arguments")
	}
	result, err := toFloat(args[0])
	if err != nil {
		return nil, err
	}
	for _, arg := range args[1:] {
		f, err := toFloat(arg)
		if err != nil {
			return nil, err
		}
		if f < result {
			result = f
		}
	}
	return result, nil
}

func funcMax(args ...any) (any, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("max requires at least 2 arguments")
	}
	result, err := toFloat(args[0])
	if err != nil {
		return nil, err
	}
	for _, arg := range args[1:] {
		f, err := toFloat(arg)
		if err != nil {
			return nil, err
		}
		if f > result {
			result = f
		}
	}
	return result, nil
}

func funcAbs(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("abs requires 1 argument")
	}
	f, err := toFloat(args[0])
	if err != nil {
		return nil, err
	}
	return math.Abs(f), nil
}

func funcFloor(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("floor requires 1 argument")
	}
	f, err := toFloat(args[0])
	if err != nil {
		return nil, err
	}
	return math.Floor(f), nil
}

func funcCeil(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ceil requires 1 argument")
	}
	f, err := toFloat(args[0])
	if err != nil {
		return nil, err
	}
	return math.Ceil(f), nil
}

func funcRound(args ...any) (any, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, fmt.Errorf("round requires 1 or 2 arguments")
	}
	f, err := toFloat(args[0])
	if err != nil {
		return nil, err
	}
	if len(args) == 1 {
		return math.Round(f), nil
	}
	precision, err := toFloat(args[1])
	if err != nil {
		return nil, err
	}
	multiplier := math.Pow(10, precision)
	return math.Round(f*multiplier) / multiplier, nil
}

func funcIf(args ...any) (any, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("if requires 3 arguments (condition, then, else)")
	}
	cond, ok := args[0].(bool)
	if !ok {
		return nil, fmt.Errorf("if condition must be bool")
	}
	if cond {
		return args[1], nil
	}
	return args[2], nil
}

func funcSqrt(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sqrt requires 1 argument")
	}
	f, err := toFloat(args[0])
	if err != nil {
		return nil, err
	}
	return math.Sqrt(f), nil
}

func funcPow(args ...any) (any, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("pow requires 2 arguments")
	}
	base, err := toFloat(args[0])
	if err != nil {
		return nil, err
	}
	exp, err := toFloat(args[1])
	if err != nil {
		return nil, err
	}
	return math.Pow(base, exp), nil
}
