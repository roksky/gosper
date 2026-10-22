# FND-009-001: Define the context and resource-budget contract

- **Status:** Not Started
- **Feature:** [FND-009](../../FEATURE_MATRIX.md#fnd--foundation-public-api-cli-infrastructure)
- **Milestone:** M0
- **Priority:** P0
- **Dependencies:** none
- **Blocked by:** none
- **Primary owner/package:** `internal/limits` (proposed)
- **Likely touched files:** `internal/limits/limits.go`, `internal/limits/limits_test.go`, `internal/limits/doc.go`
- **Collision risk:** Low — new package.
- **Parallel-safe with:** all other M0 contract tasks
- **Estimated size:** S

## Objective

Define how cancellation and resource budgets are expressed and threaded
through the engine, so every long-running stage honors them consistently.

## Current state / prerequisites

No limits package. Cancellation is not used anywhere yet.

## Scope (in)

- A `Budget` struct: `MaxRecords`, `MaxPages`, `MaxSteps`, `MaxBytes`.
- A `Counter`/`Guard` helper that checks a budget and returns a diagnostic when
  exceeded (depends on `FND-003-001` code `limit`).
- A documented rule: every long-running entry point takes `context.Context`
  first and returns `context.Canceled`/`DeadlineExceeded` unchanged.
- Tests for budget exceedance and cancellation passthrough.

## Out of scope

- Applying budgets inside fill/layout (later tasks `OPS-001-001`).
- Global default budgets.
- Metrics/telemetry.

## Proposed contracts and integration points

```go
// Proposed — non-binding.
type Budget struct { MaxRecords, MaxPages int; MaxSteps, MaxBytes int64 }
func (b Budget) Guard(ctx context.Context) *Guard
type Guard struct { /* counters */ }
func (g *Guard) Record() error   // one record
func (g *Guard) Page() error     // one page
func (g *Guard) Step(n int64) error
func (g *Guard) Bytes(n int64) error
```

Integration: `internal/fill` and `internal/layout` call the guard; `report`
constructs the budget from options.

## Implementation checklist (future agent only)

- [ ] Implement `Budget`, `Guard` with `context` checks on each call.
- [ ] Return a `diag` code `limit` on exceedance; return `ctx.Err()` on cancel.
- [ ] Zero-value `Budget` means "unlimited" (document explicitly).
- [ ] Table tests for each limit and for cancellation.
- [ ] Run `make check`.

## Acceptance criteria (observable)

- Exceeding `MaxRecords` returns a diagnostic with code `limit`.
- A canceled context returns `context.Canceled` via `errors.Is`.
- Unlimited (`Budget{}`) never trips limits.

## Tests and fixtures to add (future agent only)

- `internal/limits/limits_test.go`: per-limit and cancellation tables.

## Failure cases / edge cases

- Negative budget values: reject or treat as unlimited — pick one and test.
- Concurrent use of one `Guard`: document single-goroutine use.

## Definition of done

`make check` passes; contract documented; referenced from architecture docs.

## References and downstream dependencies

- Feature: `FND-009`; architecture:
  [`../../architecture/concurrency-and-limits.md`](../../architecture/concurrency-and-limits.md).
- Unblocks: `FND-001-001`, `OPS-001-001`, `TST-009-001`.
