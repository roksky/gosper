package expression

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"
)

var errBooleanOperand = errors.New("boolean operands required")

// Eval evaluates a program in a scope. Missing symbols evaluate to null.
// Evaluation visits at most the nodes in the bounded compiled program.
func (p *Program) Eval(ctx context.Context, scope Scope) (Value, error) {
	if p == nil || p.root == nil {
		return nil, fmt.Errorf("nil program")
	}
	if ctx == nil {
		return nil, fmt.Errorf("nil context")
	}
	return p.eval(ctx, scope, p.root)
}
func (p *Program) eval(ctx context.Context, scope Scope, n *node) (Value, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if n.op == "" {
		return n.value, nil
	}
	if n.op == "ref" {
		v, ok := n.ref.Resolve(scope)
		if !ok || v == nil {
			return NullValue, nil
		}
		return v, nil
	}
	if n.op == "?:" {
		condition, err := p.eval(ctx, scope, n.left)
		if err != nil {
			return nil, err
		}
		if b, ok := condition.(Bool); ok && bool(b) {
			return p.eval(ctx, scope, n.middle)
		}
		return p.eval(ctx, scope, n.right)
	}
	if strings.HasPrefix(n.op, "method:") {
		receiver, err := p.eval(ctx, scope, n.left)
		if err != nil {
			return nil, err
		}
		args := make([]Value, len(n.args))
		for i, arg := range n.args {
			v, err := p.eval(ctx, scope, arg)
			if err != nil {
				return nil, err
			}
			args[i] = v
		}
		return callMethod(strings.TrimPrefix(n.op, "method:"), receiver, args)
	}
	if n.left == nil {
		v, err := p.eval(ctx, scope, n.right)
		if err != nil {
			return nil, err
		}
		return unary(n.op, v, p.src, n.offset)
	}
	left, err := p.eval(ctx, scope, n.left)
	if err != nil {
		return nil, err
	}
	if n.op == "&&" || n.op == "||" {
		if b, ok := left.(Bool); ok && ((n.op == "&&" && !bool(b)) || (n.op == "||" && bool(b))) {
			return b, nil
		}
	}
	right, err := p.eval(ctx, scope, n.right)
	if err != nil {
		return nil, err
	}
	return binary(n.op, left, right, p.src, n.offset)
}
func unary(op string, v Value, src string, offset int) (Value, error) {
	if result, handled := nullUnary(v); handled {
		return result, nil
	}
	if _, ok := v.(Decimal); ok {
		number, valid := numberRat(v)
		if valid {
			if op == "u+" {
				return v, nil
			}
			if op == "u-" {
				return decimalValue(number.Neg(number)), nil
			}
		}
	}
	switch op {
	case "u!":
		if b, ok := v.(Bool); ok {
			return Bool(!bool(b)), nil
		}
	case "u+":
		if x, ok := v.(Int); ok {
			return x, nil
		}
	case "u-":
		if x, ok := v.(Int); ok {
			if x == Int(-1<<63) {
				return nil, runtimeError(src, offset, "integer overflow")
			}
			return -x, nil
		}
	}
	return nil, runtimeError(src, offset, "invalid unary operand")
}
func binary(op string, a, b Value, src string, offset int) (Value, error) {
	if op == "+" {
		if x, ok := a.(String); ok {
			return String(string(x) + scalarString(b)), nil
		}
		if y, ok := b.(String); ok {
			return String(scalarString(a) + string(y)), nil
		}
	}

	if result, handled, err := nullBinary(op, a, b); handled {
		if err != nil {
			return nil, runtimeError(src, offset, err.Error())
		}
		return result, nil
	}
	if op == "&&" || op == "||" {
		x, xok := a.(Bool)
		y, yok := b.(Bool)
		if !xok || !yok {
			return nil, runtimeError(src, offset, "boolean operands required")
		}
		if op == "&&" {
			return Bool(bool(x) && bool(y)), nil
		}
		return Bool(bool(x) || bool(y)), nil
	}
	if _, ok := a.(Decimal); ok {
		if _, valid := numberRat(b); valid {
			return numericBinary(op, a, b)
		}
	}
	if _, ok := b.(Decimal); ok {
		if _, valid := numberRat(a); valid {
			return numericBinary(op, a, b)
		}
	}
	if op == "==" || op == "!=" {
		equal, err := equalValues(a, b)
		if err != nil {
			return nil, runtimeError(src, offset, err.Error())
		}
		if op == "!=" {
			equal = !equal
		}
		return Bool(equal), nil
	}
	if op == "+" {
		if x, ok := a.(String); ok {
			return String(string(x) + scalarString(b)), nil
		}
		if y, ok := b.(String); ok {
			return String(scalarString(a) + string(y)), nil
		}
	}
	x, xok := a.(Int)
	y, yok := b.(Int)
	if xok && yok {
		switch op {
		case "+", "-", "*":
			z := big.NewInt(int64(x))
			switch op {
			case "+":
				z.Add(z, big.NewInt(int64(y)))
			case "-":
				z.Sub(z, big.NewInt(int64(y)))
			case "*":
				z.Mul(z, big.NewInt(int64(y)))
			}
			if !z.IsInt64() {
				return nil, runtimeError(src, offset, "integer overflow")
			}
			return Int(z.Int64()), nil
		case "/", "%":
			if y == 0 {
				return nil, runtimeError(src, offset, "division by zero")
			}
			if x == Int(-1<<63) && y == -1 {
				return nil, runtimeError(src, offset, "integer overflow")
			}
			if op == "/" {
				return x / y, nil
			}
			return x % y, nil
		case "<":
			return Bool(x < y), nil
		case "<=":
			return Bool(x <= y), nil
		case ">":
			return Bool(x > y), nil
		case ">=":
			return Bool(x >= y), nil
		}
	}
	if x, xok := a.(String); xok {
		if y, yok := b.(String); yok {
			switch op {
			case "<":
				return Bool(x < y), nil
			case "<=":
				return Bool(x <= y), nil
			case ">":
				return Bool(x > y), nil
			case ">=":
				return Bool(x >= y), nil
			}
		}
	}
	return nil, runtimeError(src, offset, "incompatible operands or unsupported decimal arithmetic")
}
func equalValues(a, b Value) (bool, error) {
	switch x := a.(type) {
	case Bool:
		y, ok := b.(Bool)
		return ok && x == y, nil
	case Int:
		y, ok := b.(Int)
		return ok && x == y, nil
	case String:
		y, ok := b.(String)
		return ok && x == y, nil
	case Decimal:
		y, ok := b.(Decimal)
		return ok && x == y, nil
	case Time:
		y, ok := b.(Time)
		return ok && time.Time(x).Equal(time.Time(y)), nil
	}
	return false, fmt.Errorf("incompatible equality operands")
}

func scalarString(v Value) string {
	switch x := v.(type) {
	case String:
		return string(x)
	case Decimal:
		return string(x)
	case Time:
		return time.Time(x).Format(time.RFC3339)
	case Null:
		return "null"
	default:
		return fmt.Sprint(v)
	}
}
func callMethod(name string, receiver Value, args []Value) (Value, error) {
	if isNull(receiver) {
		return nil, fmt.Errorf("method %s on null", name)
	}
	if (name == "formatDate" || name == "formatNumber") && len(args) == 1 {
		if pattern, ok := receiver.(String); ok {
			value, err := FormatValue(args[0], string(pattern))
			return String(value), err
		}
	}
	if name == "equals" && len(args) == 1 {
		eq, err := equalValues(receiver, args[0])
		return Bool(eq), err
	}
	if name == "toString" && len(args) == 0 {
		return String(scalarString(receiver)), nil
	}
	if s, ok := receiver.(String); ok && len(args) == 0 {
		switch name {
		case "trim":
			return String(strings.TrimSpace(string(s))), nil
		case "toLowerCase":
			return String(strings.ToLower(string(s))), nil
		case "toUpperCase":
			return String(strings.ToUpper(string(s))), nil
		}
	}
	if len(args) == 1 {
		ops := map[string]string{"add": "+", "subtract": "-", "multiply": "*", "divide": "/"}
		if op, ok := ops[name]; ok {
			return numericBinary(op, receiver, args[0])
		}
	}
	return nil, fmt.Errorf("unsupported %s method operands", name)
}
