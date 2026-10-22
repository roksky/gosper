# FIL-001-001: Fill record iteration

- **Status:** Done
- **Feature:** [FIL-001](../../FEATURE_MATRIX.md#fil--fill-engine-and-report-semantics)
- **Milestone:** M1
- **Priority:** P0
- **Dependencies:** DAT-002-001, FND-001-001
- **Blocked by:** none
- **Primary owner/package:** `internal/fill`
- **Likely touched files:** `internal/fill/fill.go`, `internal/fill/fill_test.go`
- **Collision risk:** High — core fill loop other FIL tasks extend. Keep the M1 loop minimal (no groups/variables).
- **Parallel-safe with:** `LAY-005-001`, `EXP-003-001`, `JRX-008-001`
- **Estimated size:** M

## Objective

Implement the fill entry point: create per-run state, open the data source,
iterate records, and produce a filled structure (record values bound to fields)
for layout.

## Current state / prerequisites

Cursor contract, public contract, and (ideally) `EXP-003-001` evaluation.

## Scope (in)

- Open cursor; for each record collect field values via `$F{}` resolution.
- Guarantee `Close` on all paths (including cancellation/error).
- One state object per run; no shared mutable state.
- Tests: record count, close-on-error, cancellation.

## Out of scope

- Variables/groups/subreports (`FIL-003`+).
- Expressions beyond field references (compose with `EXP-003-001`).

## Proposed contracts and integration points

```go
// Proposed — non-binding.
type State struct { /* per-run */ }
func Run(ctx context.Context, tmpl *Template, ds DataSource, params Params, b limits.Budget) (*Filled, error)
```

## Implementation checklist (future agent only)

- [x] Implement the loop with `defer cursor.Close()` and `ctx` checks.
- [x] Bind field values per record.
- [x] Tests for count, error, cancellation, isolation across runs.
- [x] Run `make check` and `-race`.

## Acceptance criteria (observable)

- N records ⇒ N filled records.
- Data-source error closes the cursor and returns the error.
- Two concurrent `Run` calls on the same compiled template are race-free.

## Tests and fixtures to add (future agent only)

- `internal/fill/fill_test.go` with a fake data source.

## Failure cases / edge cases

- Empty data source; record with missing fields.

## Definition of done

Implemented/tested; documented; linked from architecture.

## References and downstream dependencies

- Feature: `FIL-001`; architecture:
  [`../../architecture/filling-and-groups.md`](../../architecture/filling-and-groups.md).
- Unblocks: `FIL-002-001`, `LAY-003-001`, `TST-009-001`.
