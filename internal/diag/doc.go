// Package diag defines the diagnostic and error value model shared across the
// engine: a stable code, a severity, a human-readable message, and an optional
// source span.
//
// A Diagnostic is an error, so it can be wrapped with fmt.Errorf and recovered
// with errors.As; As is a convenience for that. A List aggregates diagnostics
// and formats them deterministically in slice order.
package diag
