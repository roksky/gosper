package expression

// isNull treats an absent Symbols value as null as well as the Null value.
func isNull(v Value) bool {
	if v == nil {
		return true
	}
	_, ok := v.(Null)
	return ok
}

// nullUnary applies null propagation before normal operand checking.
func nullUnary(v Value) (Value, bool) {
	if isNull(v) {
		return NullValue, true
	}
	return nil, false
}

// nullBinary implements the expression null decision table. A false handled
// result leaves ordinary typed evaluation to the caller.
func nullBinary(op string, a, b Value) (Value, bool, error) {
	an, bn := isNull(a), isNull(b)
	if !an && !bn {
		return nil, false, nil
	}
	switch op {
	case "==":
		return Bool(an && bn), true, nil
	case "!=":
		return Bool(an != bn), true, nil
	case "&&", "||":
		if !an {
			if _, ok := a.(Bool); !ok {
				return nil, true, errBooleanOperand
			}
		}
		if !bn {
			if _, ok := b.(Bool); !ok {
				return nil, true, errBooleanOperand
			}
		}
		if op == "&&" {
			if a == Bool(false) || b == Bool(false) {
				return Bool(false), true, nil
			}
		} else if a == Bool(true) || b == Bool(true) {
			return Bool(true), true, nil
		}
		return NullValue, true, nil
	default:
		return NullValue, true, nil
	}
}
