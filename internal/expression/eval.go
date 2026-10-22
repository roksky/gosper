package expression

import (
	"context"
	"fmt"
	"math/big"
	"time"
)

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
	if _, ok := v.(Null); ok {
		return NullValue, nil
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
	if _, ok := a.(Null); ok {
		return NullValue, nil
	}
	if _, ok := b.(Null); ok {
		return NullValue, nil
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
			if y, ok := b.(String); ok {
				return x + y, nil
			}
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
