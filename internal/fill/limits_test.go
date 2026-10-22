package fill

import (
	"context"
	"github.com/roksky/gosper/internal/datasource"
	"github.com/roksky/gosper/internal/diag"
	"github.com/roksky/gosper/internal/limits"
	"testing"
)

func TestFillBudgetBoundaries(t *testing.T) {
	src := datasource.FromSlice([]map[string]int{{"n": 1}, {"n": 2}})
	tmpl := &Template{Fields: []string{"n"}}
	for _, tc := range []struct {
		name      string
		budget    limits.Budget
		wantLimit bool
	}{{"unlimited", limits.Budget{}, false}, {"records exact", limits.Budget{MaxRecords: 2}, false}, {"records exceeded", limits.Budget{MaxRecords: 1}, true}, {"steps exact", limits.Budget{MaxSteps: 4}, false}, {"steps exceeded", limits.Budget{MaxSteps: 3}, true}, {"bytes exact", limits.Budget{MaxBytes: 18}, false}, {"bytes exceeded", limits.Budget{MaxBytes: 17}, true}} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Run(context.Background(), tmpl, src, tc.budget)
			if tc.wantLimit {
				d, ok := diag.As(err)
				if !ok || d.Code != diag.CodeLimit || got != nil {
					t.Fatalf("limit result=%v,%v", got, err)
				}
			} else if err != nil || len(got.Records) != 2 {
				t.Fatalf("allowed result=%v,%v", got, err)
			}
		})
	}
}
func TestParameterByteBudget(t *testing.T) {
	src := datasource.FromSlice([]map[string]any(nil))
	tmpl := &Template{Parameters: []ParameterDefault{{Name: "param"}}}
	_, err := Run(context.Background(), tmpl, src, limits.Budget{MaxBytes: 1})
	d, ok := diag.As(err)
	if !ok || d.Code != diag.CodeLimit {
		t.Fatalf("parameter limit=%v", err)
	}
}
