# TST-009-001: Race, cancellation and leak tests

- **Status:** Not Started
- **Feature:** [TST-009](../../FEATURE_MATRIX.md#tst--test-infrastructure-and-correctness-suites)
- **Milestone:** M1
- **Priority:** P1
- **Dependencies:** FIL-001-001, OPS-001-001
- **Blocked by:** none
- **Primary owner/package:** `internal/fill` (tests) + `internal/testutil`
- **Likely touched files:** `internal/fill/concurrency_test.go`, `internal/testutil/leak.go`
- **Collision risk:** Medium — tests may reveal needed fixes in `fill`/`limits`; coordinate with their owners.
- **Parallel-safe with:** `internal/jrxml`/`internal/text` tasks
- **Estimated size:** M

## Objective

Prove concurrency safety and resource hygiene: concurrent fills of one compiled
template are race-free; cancellation closes cursors and drains goroutines.

## Current state / prerequisites

Fill loop and limits exist.

## Scope (in)

- A race test: N goroutines fill the same compiled template concurrently under
  `-race`.
- A cancellation test: cancel mid-iteration; assert cursor closed and no leaked
  goroutines (goroutine-count check).
- Optional small leak helper in `internal/testutil`.

## Out of scope

- Exporter concurrency.

## Proposed contracts and integration points

Uses a fake data source that records `Close` and can block until cancel.

## Implementation checklist (future agent only)

- [ ] Add concurrency + cancellation + leak tests.
- [ ] Run `make test-race`.
- [ ] Fix any revealed leaks/races (small, isolated).

## Acceptance criteria (observable)

- `go test -race ./...` passes.
- Cancellation test observes `Close` and a stable goroutine count.

## Tests and fixtures to add (future agent only)

- `internal/fill/concurrency_test.go`.

## Failure cases / edge cases

- Blocking data source that ignores `ctx` — test must still terminate (timeout).

## Definition of done

Tests pass under `-race`; any fixes reviewed in the same PR.

## References and downstream dependencies

- Feature: `TST-009`; architecture:
  [`../../architecture/concurrency-and-limits.md`](../../architecture/concurrency-and-limits.md).
