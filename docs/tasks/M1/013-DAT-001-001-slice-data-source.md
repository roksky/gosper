# DAT-001-001: Slice data source

- **Status:** Not Started
- **Feature:** [DAT-001](../../FEATURE_MATRIX.md#dat--data-sources-and-queries)
- **Milestone:** M1
- **Priority:** P0
- **Dependencies:** DAT-002-001
- **Blocked by:** none
- **Primary owner/package:** `internal/datasource`
- **Likely touched files:** `internal/datasource/slice.go`, `internal/datasource/slice_test.go`
- **Collision risk:** Low — new file next to the interface; does not touch the interface file.
- **Parallel-safe with:** `LAY-005-001`, `EXP-003-001`, `JRX-002-001`
- **Estimated size:** S/M

## Objective

Implement a data source over Go slices of structs, maps, or `[]any`, exposing
fields by name for `$F{}` binding.

## Current state / prerequisites

Cursor interface contract in place.

## Scope (in)

- `Open` over a slice; `Next` yields a `Record`; field access by struct tag,
  field name, or map key.
- Reflection-based mapping with a documented field-name rule.
- Tests: structs, maps, cancellation mid-iteration.

## Out of scope

- Type conversion policy (`DAT-007`); JSON/CSV/SQL.

## Proposed contracts and integration points

```go
// Proposed — non-binding.
func FromSlice[T any](rows []T) DataSource
```

## Implementation checklist (future agent only)

- [ ] Implement slice/map/struct adaptation.
- [ ] Honor `ctx` in `Next`.
- [ ] Tests for each shape and cancellation.
- [ ] Run `make check`.

## Acceptance criteria (observable)

- Iterating a 3-element slice yields 3 records with correct field values.
- Cancellation returns `context.Canceled`.

## Tests and fixtures to add (future agent only)

- `internal/datasource/slice_test.go`.

## Failure cases / edge cases

- Unknown field name → `(nil,false)`; nil slice → zero records.

## Definition of done

Implemented/tested; documented; linked from `features/DAT.md`.

## References and downstream dependencies

- Feature: `DAT-001`.
- Unblocks: `FIL-001-001`.
