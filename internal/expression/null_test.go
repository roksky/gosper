package expression

import (
	"context"
	"testing"
)

func TestNullDecisionTable(t *testing.T) {
	for _, tc := range []struct {
		src  string
		want Value
	}{
		{"null + 1", NullValue}, {"2 * null", NullValue}, {"null / 0", NullValue},
		{"null < 1", NullValue}, {"1 >= null", NullValue},
		{"null == null", Bool(true)}, {"null != null", Bool(false)},
		{"null == 1", Bool(false)}, {"1 != null", Bool(true)},
		{"!null", NullValue}, {"-null", NullValue},
		{"false && null", Bool(false)}, {"null && false", Bool(false)},
		{"true && null", NullValue}, {"null && true", NullValue},
		{"true || null", Bool(true)}, {"null || true", Bool(true)},
		{"false || null", NullValue}, {"null || false", NullValue},
		{"null && null", NullValue}, {"null || null", NullValue},
		{"missing + 1", NullValue}, {"$F{missing} == null", Bool(true)},
	} {
		p, err := Compile(tc.src)
		if err != nil {
			t.Errorf("compile %q: %v", tc.src, err)
			continue
		}
		got, err := p.Eval(context.Background(), nil)
		if err != nil || got != tc.want {
			t.Errorf("%q = %v, %v; want %v", tc.src, got, err, tc.want)
		}
	}
}
func TestNullBooleanShortCircuit(t *testing.T) {
	for _, tc := range []struct {
		src  string
		want Value
	}{
		{"false && 1 / 0 == 1", Bool(false)},
		{"true || 1 / 0 == 1", Bool(true)},
		{"null && false", Bool(false)},
		{"null || true", Bool(true)},
	} {
		p, _ := Compile(tc.src)
		got, err := p.Eval(context.Background(), nil)
		if err != nil || got != tc.want {
			t.Errorf("%q = %v, %v", tc.src, got, err)
		}
	}
}
func TestNullBooleanTypeError(t *testing.T) {
	for _, src := range []string{"null && 1", "null || \"x\""} {
		p, _ := Compile(src)
		if _, err := p.Eval(context.Background(), nil); err == nil {
			t.Errorf("%q should reject non-boolean", src)
		}
	}
}
