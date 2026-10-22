package fill

import "github.com/roksky/gosper/internal/expression"

// valueSize counts retained scalar payload bytes for the M1 byte budget.
// Map and slice allocator overhead is deliberately outside this estimate.
func valueSize(v expression.Value) int64 {
	switch x := v.(type) {
	case expression.String:
		return int64(len(x))
	case expression.Decimal:
		return int64(len(x))
	case expression.Int:
		return 8
	case expression.Bool:
		return 1
	case expression.Time:
		return 24
	default:
		return 0
	}
}
