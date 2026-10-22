package limits

import (
	"context"
	"errors"
	"math"
	"testing"

	"github.com/roksky/gosper/internal/diag"
)

func TestBudgetLimits(t *testing.T) {
	cases := []struct {
		name    string
		budget  Budget
		trip    func(g *Guard) error
		wantErr bool
	}{
		{
			name:    "record within",
			budget:  Budget{MaxRecords: 2},
			trip:    func(g *Guard) error { return g.Record() },
			wantErr: false,
		},
		{
			name:    "record exceeded",
			budget:  Budget{MaxRecords: 1},
			trip:    func(g *Guard) error { g.Record(); return g.Record() },
			wantErr: true,
		},
		{
			name:    "page exceeded",
			budget:  Budget{MaxPages: 1},
			trip:    func(g *Guard) error { g.Page(); return g.Page() },
			wantErr: true,
		},
		{
			name:    "step exceeded",
			budget:  Budget{MaxSteps: 10},
			trip:    func(g *Guard) error { return g.Step(11) },
			wantErr: true,
		},
		{
			name:    "bytes exceeded",
			budget:  Budget{MaxBytes: 100},
			trip:    func(g *Guard) error { return g.Bytes(101) },
			wantErr: true,
		},
		{
			name:    "unlimited zero value",
			budget:  Budget{},
			trip:    func(g *Guard) error { g.Record(); g.Page(); g.Step(1e9); return g.Bytes(1 << 40) },
			wantErr: false,
		},
		{
			name:    "negative treated as unlimited",
			budget:  Budget{MaxRecords: -1, MaxPages: -5, MaxSteps: -10, MaxBytes: -100},
			trip:    func(g *Guard) error { g.Record(); g.Page(); g.Step(1e9); return g.Bytes(1 << 40) },
			wantErr: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g := tc.budget.Guard(context.Background())
			err := tc.trip(g)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected a limit error, got nil")
				}
				d, ok := diag.As(err)
				if !ok {
					t.Fatalf("error is not a diagnostic: %v", err)
				}
				if d.Code != diag.CodeLimit {
					t.Fatalf("code = %q, want %q", d.Code, diag.CodeLimit)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestExactLimitAllowed(t *testing.T) {
	g := Budget{MaxRecords: 3}.Guard(context.Background())
	for i := 0; i < 3; i++ {
		if err := g.Record(); err != nil {
			t.Fatalf("record %d unexpectedly failed: %v", i+1, err)
		}
	}
	if err := g.Record(); err == nil {
		t.Fatal("record beyond the limit should fail")
	}
}

func TestCancellationPassthrough(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	g := Budget{}.Guard(ctx)
	for _, err := range []error{g.Record(), g.Page(), g.Step(1), g.Bytes(1)} {
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("guard returned %v, want context.Canceled", err)
		}
	}
}

func TestDeadlinePassthrough(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	g := Budget{MaxRecords: 1}.Guard(ctx)
	if err := g.Record(); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation should take precedence over limits, got %v", err)
	}
}

func TestAccountingCannotWrapOrSubtract(t *testing.T) {
	g := Budget{MaxSteps: 5, MaxBytes: 5}.Guard(context.Background())
	if err := g.Step(math.MaxInt64); err == nil {
		t.Fatal("large step count bypassed limit")
	}
	if err := g.Bytes(math.MaxInt64); err == nil {
		t.Fatal("large byte count bypassed limit")
	}
	if err := g.Step(-1); err == nil {
		t.Fatal("negative step count accepted")
	}
	if err := g.Bytes(-1); err == nil {
		t.Fatal("negative byte count accepted")
	}
	unlimited := Budget{}.Guard(context.Background())
	if err := unlimited.Step(math.MaxInt64); err != nil {
		t.Fatal(err)
	}
	if err := unlimited.Step(1); err != nil {
		t.Fatal(err)
	}
	if err := unlimited.Bytes(math.MaxInt64); err != nil {
		t.Fatal(err)
	}
	if err := unlimited.Bytes(1); err != nil {
		t.Fatal(err)
	}
}
