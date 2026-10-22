// Package limits defines resource budgets and cancellation plumbing for
// long-running engine stages.
//
// Every long-running entry point accepts a context.Context and honors
// cancellation. A Budget caps records, pages, evaluation steps and bytes; the
// zero value means "unlimited". Guard is a single-goroutine accounting helper
// that enforces a Budget and surfaces diagnostics with code diag.CodeLimit.
package limits
