package limits

import (
	"context"
	"math"

	"github.com/roksky/gosper/internal/diag"
)

// Budget caps the resources a single report run may consume. The zero value is
// unlimited: a limit that is zero or negative means "no limit" for that
// dimension.
type Budget struct {
	MaxRecords int
	MaxPages   int
	MaxSteps   int64
	MaxBytes   int64
}

// Guard accounts for resource use against a Budget for a single run.
//
// A Guard is not safe for concurrent use. Every method first checks the run
// context and returns ctx.Err() when it is done; otherwise it records the
// requested amount and returns a *diag.Diagnostic with code diag.CodeLimit when
// a limit is exceeded.
type Guard struct {
	budget Budget
	ctx    context.Context

	records int
	pages   int
	steps   int64
	bytes   int64
}

// Guard starts accounting for one run. ctx must be non-nil.
func (b Budget) Guard(ctx context.Context) *Guard {
	return &Guard{budget: b, ctx: ctx}
}

// check returns the context's error, if any.
func (g *Guard) check() error {
	if g.ctx == nil {
		return nil
	}
	return g.ctx.Err()
}

func limitErr(msg string) error {
	return diag.New(diag.CodeLimit, diag.SeverityError, msg, nil)
}

// Record accounts for one record and reports whether the record budget is
// exceeded.
func (g *Guard) Record() error {
	if err := g.check(); err != nil {
		return err
	}
	if g.budget.MaxRecords > 0 && g.records >= g.budget.MaxRecords {
		return limitErr("record budget exceeded")
	}
	if g.records < math.MaxInt {
		g.records++
	}
	return nil
}

// Page accounts for one page and reports whether the page budget is exceeded.
func (g *Guard) Page() error {
	if err := g.check(); err != nil {
		return err
	}
	if g.budget.MaxPages > 0 && g.pages >= g.budget.MaxPages {
		return limitErr("page budget exceeded")
	}
	if g.pages < math.MaxInt {
		g.pages++
	}
	return nil
}

// Step accounts for n evaluation steps and reports whether the step budget is
// exceeded.
func (g *Guard) Step(n int64) error {
	if err := g.check(); err != nil {
		return err
	}
	if n < 0 {
		return limitErr("negative evaluation step count")
	}
	if g.budget.MaxSteps > 0 && n > g.budget.MaxSteps-g.steps {
		return limitErr("evaluation step budget exceeded")
	}
	if n > math.MaxInt64-g.steps {
		g.steps = math.MaxInt64
	} else {
		g.steps += n
	}
	return nil
}

// Bytes accounts for n bytes of retained or produced data and reports whether
// the byte budget is exceeded.
func (g *Guard) Bytes(n int64) error {
	if err := g.check(); err != nil {
		return err
	}
	if n < 0 {
		return limitErr("negative byte count")
	}
	if g.budget.MaxBytes > 0 && n > g.budget.MaxBytes-g.bytes {
		return limitErr("byte budget exceeded")
	}
	if n > math.MaxInt64-g.bytes {
		g.bytes = math.MaxInt64
	} else {
		g.bytes += n
	}
	return nil
}
