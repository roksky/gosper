package expression

import (
	"fmt"
	"math/big"
	"strings"
)

func numberRat(value Value) (*big.Rat, bool) {
	switch v := value.(type) {
	case Int:
		return new(big.Rat).SetInt64(int64(v)), true
	case Decimal:
		r, ok := new(big.Rat).SetString(string(v))
		return r, ok
	}
	return nil, false
}
func decimalValue(r *big.Rat) Value {
	if r.IsInt() && r.Num().IsInt64() {
		return Int(r.Num().Int64())
	}
	if r.IsInt() {
		return Decimal(r.Num().String())
	}
	// Finite base-10 results retain their full scale. Only repeating
	// quotients use a bounded 16-place representation.
	denominator := new(big.Int).Set(r.Denom())
	twos := int(denominator.TrailingZeroBits())
	denominator.Rsh(denominator, uint(twos))
	fives := 0
	five := big.NewInt(5)
	for denominator.Sign() > 0 {
		q, rem := new(big.Int), new(big.Int)
		q.QuoRem(denominator, five, rem)
		if rem.Sign() != 0 {
			break
		}
		denominator = q
		fives++
	}
	scale := 16
	if denominator.Cmp(big.NewInt(1)) == 0 {
		scale = max(twos, fives)
	}
	s := strings.TrimRight(strings.TrimRight(r.FloatString(scale), "0"), ".")
	return Decimal(s)
}

// AddNumbers implements exact rational accumulation for report sums.
func AddNumbers(a, b Value) (Value, error) { return numericBinary("+", a, b) }
func numericBinary(op string, a, b Value) (Value, error) {
	x, ok := numberRat(a)
	if !ok {
		return nil, fmt.Errorf("numeric operand required")
	}
	y, ok := numberRat(b)
	if !ok {
		return nil, fmt.Errorf("numeric operand required")
	}
	switch op {
	case "+":
		return decimalValue(new(big.Rat).Add(x, y)), nil
	case "-":
		return decimalValue(new(big.Rat).Sub(x, y)), nil
	case "*":
		return decimalValue(new(big.Rat).Mul(x, y)), nil
	case "/":
		if y.Sign() == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return decimalValue(new(big.Rat).Quo(x, y)), nil
	case "==":
		return Bool(x.Cmp(y) == 0), nil
	case "!=":
		return Bool(x.Cmp(y) != 0), nil
	case "<":
		return Bool(x.Cmp(y) < 0), nil
	case "<=":
		return Bool(x.Cmp(y) <= 0), nil
	case ">":
		return Bool(x.Cmp(y) > 0), nil
	case ">=":
		return Bool(x.Cmp(y) >= 0), nil
	}
	return nil, fmt.Errorf("unsupported numeric operation %s", op)
}
