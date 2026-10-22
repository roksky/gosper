# TST-001-001: Deterministic test environment helpers

- **Status:** Not Started
- **Feature:** [TST-001](../../FEATURE_MATRIX.md#tst--test-infrastructure-and-correctness-suites)
- **Milestone:** M0
- **Priority:** P0
- **Dependencies:** none
- **Blocked by:** none
- **Primary owner/package:** `internal/testutil` (proposed)
- **Likely touched files:** `internal/testutil/env.go`, `internal/testutil/env_test.go`, `internal/testutil/doc.go`
- **Collision risk:** Low — new package.
- **Parallel-safe with:** all other M0 tasks
- **Estimated size:** S

## Objective

Provide deterministic helpers for time/clock, timezone, locale, and map-order
avoidance so tests produce stable output.

## Current state / prerequisites

No test helpers; existing tests are self-contained. Baseline `make check`/CI
already run tests.

## Scope (in)

- Fixed clock helper and fixed `*time.Location`.
- A documented default locale/timezone for tests.
- A helper to iterate maps in sorted order.
- Tests validating determinism.

## Out of scope

- Font fixtures (`LAY-005`).
- Golden comparison (`TST-002-001`).

## Proposed contracts and integration points

```go
// Proposed — non-binding.
func FixedClock(t time.Time) func() time.Time
func UTC() *time.Location
func SortedKeys[V any](m map[string]V) []string
```

Integration: used by all future engine tests.

## Implementation checklist (future agent only)

- [ ] Add helpers with docs.
- [ ] Tests assert stable behavior across calls.
- [ ] Run `make check` (including `-race`).

## Acceptance criteria (observable)

- `FixedClock` returns the same instant on every call.
- `SortedKeys` yields lexicographically ordered keys regardless of insertion.

## Tests and fixtures to add (future agent only)

- `internal/testutil/env_test.go`.

## Failure cases / edge cases

- Concurrent use of helpers must be race-free (test with `-race`).

## Definition of done

Helpers available and tested; referenced from
[`../../TEST_STRATEGY.md`](../../TEST_STRATEGY.md).

## References and downstream dependencies

- Feature: `TST-001`.
- Unblocks: `TST-001-002`, `TST-002-001`, date-related `EXP-*`/`LAY-*` tests.
