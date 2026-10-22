# FND-003-001: Define the diagnostic and error code model

- **Status:** Not Started
- **Feature:** [FND-003](../../FEATURE_MATRIX.md#fnd--foundation-public-api-cli-infrastructure)
- **Milestone:** M0
- **Priority:** P0
- **Dependencies:** none
- **Blocked by:** none
- **Primary owner/package:** `internal/diag` (proposed)
- **Likely touched files:** `internal/diag/diag.go`, `internal/diag/diag_test.go`, `internal/diag/doc.go`
- **Collision risk:** Low — new package, no other task touches it.
- **Parallel-safe with:** all other M0 contract tasks
- **Estimated size:** S

## Objective

Establish the error/diagnostic value model used across the engine: a stable
code, severity, message, optional source span, and wrapping so callers can
attach context without inventing new codes.

## Current state / prerequisites

No error model exists. `internal/compiler`, `internal/jrxml`, etc. are `doc.go`
stubs. Nothing imports a shared error type.

## Scope (in)

- A `Code` type (string constant) with a small initial set (e.g. `parse`,
  `unsupported`, `validation`, `canceled`, `limit`).
- A `Diagnostic` type: code, severity, message, optional `Span{Line, Col, Offset}`.
- A `List` aggregator and a function to extract diagnostics from an error tree
  via `errors.As`.
- Unit tests for code round-tripping, span formatting and aggregation.

## Out of scope

- Wiring into parser/compiler (later tasks).
- Localization.
- Any report logic.

## Proposed contracts and integration points

```go
// Proposed — non-binding.
type Code string
type Severity int // Error, Warning
type Span struct { Line, Col, Offset int }
type Diagnostic struct { Code Code; Sev Severity; Msg string; Span *Span }
func (d *Diagnostic) Error() string
func As(err error) (*Diagnostic, bool)      // via errors.As
type List []*Diagnostic
func (l List) Err() error
```

Integration: `internal/jrxml` and `internal/compiler` produce `*Diagnostic`;
`report` and CLI render them.

## Implementation checklist (future agent only)

- [ ] Add `internal/diag` package with the types above.
- [ ] Ensure `Error()` format is stable and documented.
- [ ] Implement `errors.As` compatibility (use `Unwrap`/`Is` as needed).
- [ ] Add table-driven tests.
- [ ] Run `make check`.

## Acceptance criteria (observable)

- `diag.As` returns the original code and span for a wrapped diagnostic.
- A `List` with two diagnostics stringifies deterministically (stable order).
- Unknown codes are permitted (forward-compatible) without panics.

## Tests and fixtures to add (future agent only)

- `internal/diag/diag_test.go`: code/span/aggregation tables.
- No `testdata/` fixtures needed.

## Failure cases / edge cases

- Nil span → formatting must not panic.
- Empty message → still a valid error string.

## Definition of done

Package compiles, tests pass via `make check`, exported contract documented in
the package comment, and referenced from
[`../../architecture/package-boundaries.md`](../../architecture/package-boundaries.md).

## References and downstream dependencies

- Feature: `FND-003` in [`../../FEATURE_MATRIX.md`](../../FEATURE_MATRIX.md).
- Unblocks: `JRX-001-001`, `JRX-019-001`.
