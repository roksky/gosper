package fill

import (
	"context"
	"github.com/roksky/gosper/internal/datasource"
	"github.com/roksky/gosper/internal/expression"
	"github.com/roksky/gosper/internal/limits"
	"testing"
)

func compiled(t *testing.T, src string) *expression.Program {
	t.Helper()
	p, err := expression.Compile(src)
	if err != nil {
		t.Fatal(err)
	}
	return p
}
func TestReportCountPerRecord(t *testing.T) {
	tmpl := &Template{Fields: []string{"x"}}
	ds := datasource.FromSlice([]map[string]int{{"x": 3}, {"x": 4}, {"x": 5}})
	filled, err := Run(context.Background(), tmpl, ds, limits.Budget{})
	if err != nil {
		t.Fatal(err)
	}
	p := compiled(t, "$V{REPORT_COUNT} + $F{x}")
	for i, r := range filled.Records {
		v, ok := r.Variable("REPORT_COUNT")
		if !ok || v != expression.Int(i+1) {
			t.Errorf("record %d count=%v,%v", i, v, ok)
		}
		got, err := p.Eval(context.Background(), expression.Scope{r})
		if err != nil || got != expression.Int(2*i+4) {
			t.Errorf("record %d expression=%v,%v", i, got, err)
		}
	}
}
func TestParameterDefaultsAndOverrides(t *testing.T) {
	tmpl := &Template{Parameters: []ParameterDefault{{Name: "base", Default: compiled(t, "2")}, {Name: "total", Default: compiled(t, "$P{base} + 3")}, {Name: "unset"}}}
	ds := datasource.FromSlice([]map[string]int{{}, {}})
	filled, err := RunWithParams(context.Background(), tmpl, ds, nil, limits.Budget{})
	if err != nil {
		t.Fatal(err)
	}
	if filled.Parameters["total"] != expression.Int(5) || filled.Parameters["unset"] != expression.NullValue {
		t.Fatalf("defaults=%v", filled.Parameters)
	}
	for _, r := range filled.Records {
		if v, ok := r.Parameter("total"); !ok || v != expression.Int(5) {
			t.Errorf("record total=%v,%v", v, ok)
		}
	}
	filled, err = RunWithParams(context.Background(), tmpl, ds, map[string]expression.Value{"base": expression.Int(10)}, limits.Budget{})
	if err != nil || filled.Parameters["total"] != expression.Int(13) {
		t.Fatalf("override total=%v,%v", filled, err)
	}
	filled, err = RunWithParams(context.Background(), tmpl, ds, map[string]expression.Value{"base": nil}, limits.Budget{})
	if err != nil || filled.Parameters["base"] != expression.NullValue || filled.Parameters["total"] != expression.NullValue {
		t.Fatalf("null override result=%v, %v", filled, err)
	}
}
func TestParameterDefaultOrderAndErrors(t *testing.T) {
	ds := datasource.FromSlice([]map[string]int(nil))
	tmpl := &Template{Parameters: []ParameterDefault{{Name: "later", Default: compiled(t, "$P{earlier} + 1")}, {Name: "earlier", Default: compiled(t, "4")}}}
	got, err := Run(context.Background(), tmpl, ds, limits.Budget{})
	if err != nil {
		t.Fatal(err)
	}
	if got.Parameters["later"] != expression.NullValue {
		t.Fatalf("forward reference=%v", got.Parameters["later"])
	}
	tmpl.Parameters = []ParameterDefault{{Name: "x"}, {Name: "x"}}
	if _, err := Run(context.Background(), tmpl, ds, limits.Budget{}); err == nil {
		t.Fatal("duplicate parameter accepted")
	}
}
