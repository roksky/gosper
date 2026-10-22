package report_test

import (
	"context"
	"github.com/roksky/gosper/internal/diag"
	"github.com/roksky/gosper/report"
	"testing"
)

func TestPublicBudgetOption(t *testing.T) {
	compiled := fixture(t)
	source := report.FromSlice([]map[string]string{{"name": "a"}, {"name": "b"}, {"name": "c"}})
	for _, tc := range []struct {
		name      string
		budget    report.Budget
		wantLimit bool
	}{{"records", report.Budget{MaxRecords: 2}, true}, {"pages", report.Budget{MaxPages: 1}, true}, {"exact", report.Budget{MaxRecords: 3, MaxPages: 2}, false}, {"unlimited", report.Budget{}, false}} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := report.Fill(context.Background(), compiled, nil, source, report.WithBudget(tc.budget))
			if tc.wantLimit {
				d, ok := diag.As(err)
				if !ok || d.Code != diag.CodeLimit || got != nil {
					t.Fatalf("result=%v,%v", got, err)
				}
			} else if err != nil || got.PageCount() != 2 {
				t.Fatalf("result=%v,%v", got, err)
			}
		})
	}
}
