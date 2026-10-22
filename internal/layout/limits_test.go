package layout

import (
	"context"
	"github.com/roksky/gosper/internal/diag"
	"github.com/roksky/gosper/internal/fill"
	"github.com/roksky/gosper/internal/limits"
	"testing"
)

func TestLayoutPageBudget(t *testing.T) {
	spec := &DetailSpec{Page: PageGeometry{Width: 200, Height: 100, Top: 10, Bottom: 10, Left: 10, Right: 10}, Detail: Band{Height: 40}}
	filled := &fill.Filled{Records: []fill.FilledRecord{{}, {}, {}}}
	for _, tc := range []struct {
		limit     int
		wantLimit bool
	}{{0, false}, {1, true}, {2, false}} {
		doc, err := LayoutWithBudget(context.Background(), filled, spec, limits.Budget{MaxPages: tc.limit})
		if tc.wantLimit {
			d, ok := diag.As(err)
			if !ok || d.Code != diag.CodeLimit || doc != nil {
				t.Fatalf("limit %d: %v,%v", tc.limit, doc, err)
			}
		} else if err != nil || doc.Metadata.PageCount != 2 {
			t.Fatalf("limit %d: %v,%v", tc.limit, doc, err)
		}
	}
}
