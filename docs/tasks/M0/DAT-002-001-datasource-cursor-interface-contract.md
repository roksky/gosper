# DAT-002-001: Define the DataSource cursor interface contract

- **Status:** Not Started
- **Feature:** [DAT-002](../../FEATURE_MATRIX.md#dat--data-sources-and-queries)
- **Milestone:** M0
- **Priority:** P0
- **Dependencies:** none
- **Blocked by:** none
- **Primary owner/package:** `internal/datasource`
- **Likely touched files:** `internal/datasource/datasource.go`, `internal/datasource/datasource_test.go`, `internal/datasource/doc.go`
- **Collision risk:** Low — new code in a package that only has a `doc.go` stub. M1 `DAT-001-001` will add a file next to it.
- **Parallel-safe with:** all other M0 contract tasks
- **Estimated size:** S

## Objective

Define the `DataSource`/`Cursor`/`Record` interfaces and their lifecycle
semantics (open, forward-only iteration, idempotent close, cancellation).

## Current state / prerequisites

`internal/datasource/doc.go` is a stub.

## Scope (in)

- Interfaces `DataSource`, `Cursor`, `Record`.
- Documented lifecycle rules and `Close` idempotency.
- A tiny in-package fake implementation used only by the contract test.
- Tests covering close-after-error and `Next` after `Close`.

## Out of scope

- Concrete slice/JSON/CSV/SQL sources (M1/M2 tasks).
- Type conversion policy (`DAT-007`).

## Proposed contracts and integration points

```go
// Proposed — non-binding.
type DataSource interface{ Open(ctx context.Context) (Cursor, error) }
type Cursor interface {
    Next(ctx context.Context) (Record, bool, error)
    Close() error
}
type Record interface{ Value(field string) (any, bool) }
```

Integration: `internal/fill` opens and iterates.

## Implementation checklist (future agent only)

- [ ] Define interfaces and document semantics in the package comment.
- [ ] Provide an in-test fake to validate the contract.
- [ ] Test: `Next` after `Close` returns `(nil, false, nil)` or a defined error
      — pick and document one.
- [ ] Test: double `Close` returns nil.
- [ ] Run `make check`.

## Acceptance criteria (observable)

- A fake implementing the interfaces compiles and passes lifecycle tests.
- Contract documented unambiguously (return values after close).

## Tests and fixtures to add (future agent only)

- `internal/datasource/datasource_test.go`.

## Failure cases / edge cases

- `Open` returning `(nil, nil)` must be treated as a bug; document forbidden.

## Definition of done

Interfaces defined, documented, tested; referenced by
[`../../architecture/data-and-resources.md`](../../architecture/data-and-resources.md).

## References and downstream dependencies

- Feature: `DAT-002`.
- Unblocks: `DAT-001-001`, `DAT-003-001`, `DAT-004-001`, `FIL-001-001`.
