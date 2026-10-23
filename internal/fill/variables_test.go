package fill

import (
	"context"
	"github.com/roksky/gosper/internal/datasource"
	"github.com/roksky/gosper/internal/expression"
	"github.com/roksky/gosper/internal/limits"
	"testing"
)

func TestVariableSnapshotsAndGroupReset(t *testing.T) {
	tmpl := &Template{Fields: []string{"class", "amount", "name"}, FieldTypes: map[string]string{"amount": "java.math.BigDecimal"}, Groups: []Group{{Name: "Class", Expr: compiled(t, "$F{class}")}}, Variables: []Variable{
		{Name: "groupSum", Calculation: "Sum", ResetType: "Group", ResetGroup: "Class", Expr: compiled(t, "$F{amount}")},
		{Name: "total", Calculation: "Sum", ResetType: "Report", Expr: compiled(t, "$F{amount}")},
		{Name: "count", Calculation: "Count", ResetType: "Report", Expr: compiled(t, "$F{name}")},
	}}
	rows := []map[string]any{{"class": "One", "amount": "0.10", "name": "Alice"}, {"class": "One", "amount": "0.20", "name": "Bob"}, {"class": "Two", "amount": "1.00", "name": nil}}
	filled, err := Run(context.Background(), tmpl, datasource.FromSlice(rows), limits.Budget{})
	if err != nil {
		t.Fatal(err)
	}
	for i, want := range []expression.Value{expression.Decimal("0.1"), expression.Decimal("0.3"), expression.Int(1)} {
		if got := filled.Records[i].Variables["groupSum"]; got != want {
			t.Errorf("record %d groupSum=%v want %v", i, got, want)
		}
	}
	final := filled.Records[2].Variables
	if final["total"] != expression.Decimal("1.3") || final["count"] != expression.Int(2) || final["Class_COUNT"] != expression.Int(1) {
		t.Fatalf("final variables %v", final)
	}
	if filled.Records[0].Variables["total"] != expression.Decimal("0.1") {
		t.Fatal("earlier snapshot mutated")
	}
}
