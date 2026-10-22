package expression

import (
	"context"
	"strings"
	"testing"
)

func TestEvaluation(t *testing.T) {
	scope := Scope{fakeSymbols{variables: map[string]Value{"a": Int(1), "b": Int(2), "c": Int(3)}, fields: map[string]Value{"x": Int(7)}}}
	for _, tc := range []struct {
		src  string
		want Value
	}{
		{"1 + 2 * 3", Int(7)}, {"(1 + 2) * 3", Int(9)}, {"a < b && b < c", Bool(true)},
		{"true or false and false", Bool(true)}, {"not false", Bool(true)},
		{"-5 % 3", Int(-2)}, {"$F{x} + 1", Int(8)}, {"\"a\" + \"b\"", String("ab")},
		{"10 / 2", Int(5)}, {"2 == 2", Bool(true)}, {"2 != 3", Bool(true)},
		{"false && 1 / 0 == 0", Bool(false)}, {"true || 1 / 0 == 0", Bool(true)},
	} {
		p, err := Compile(tc.src)
		if err != nil {
			t.Errorf("compile %q: %v", tc.src, err)
			continue
		}
		got, err := p.Eval(context.Background(), scope)
		if err != nil || got != tc.want {
			t.Errorf("%q = %v, %v; want %v", tc.src, got, err, tc.want)
		}
	}
}
func TestEvaluationFailures(t *testing.T) {
	for _, src := range []string{"1 / 0", "1 % 0", "9223372036854775807 + 1", "-(-9223372036854775807 - 1)", "1 + true", "1.2 + 2.3"} {
		p, err := Compile(src)
		if err != nil {
			t.Fatalf("compile %q: %v", src, err)
		}
		if _, err = p.Eval(context.Background(), nil); err == nil {
			t.Errorf("%q should fail", src)
		}
	}
	p, _ := Compile("1 + 2")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := p.Eval(ctx, nil); err == nil {
		t.Fatal("canceled evaluation succeeded")
	}
}
func TestEvaluationBoundedSource(t *testing.T) {
	if _, err := Compile(strings.Repeat("1", maxSourceBytes+1)); err == nil {
		t.Fatal("oversized source compiled")
	}
}
