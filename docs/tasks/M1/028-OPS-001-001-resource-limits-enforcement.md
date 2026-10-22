# OPS-001-001: Resource limits enforcement

- **Status:** Done
- **Feature:** [OPS-001](../../FEATURE_MATRIX.md#ops--security-performance-diagnostics-integration)
- **Milestone:** M1
- **Priority:** P0
- **Dependencies:** FND-009-001
- **Blocked by:** none
- **Primary owner/package:** `internal/limits` + `internal/fill`
- **Likely touched files:** `internal/fill/limits.go`, `internal/fill/limits_test.go` (wiring), maybe `internal/limits/limits.go` (fixes)
- **Collision risk:** Medium — touches the fill loop owned by `FIL-001-001`; coordinate a small, clearly-scoped hunk.
- **Parallel-safe with:** `internal/jrxml`/`internal/text` tasks
- **Estimated size:** S

## Objective

Enforce the per-run budgets (records/pages/steps/bytes) inside fill so runaway
runs abort with a `limit` diagnostic instead of exhausting resources.

## Current state / prerequisites

Budget/Guard contract (`FND-009-001`) and the fill loop exist.

## Scope (in)

- Call the guard per record and per page; abort with `diag` code `limit`.
- Tests: a fixture exceeding `MaxRecords` and one exceeding `MaxPages` abort
  with the expected code; unlimited budget does not abort.

## Out of scope

- Page/step budgets beyond M1 scope; memory profiling (`OPS-007`).

## Proposed contracts and integration points

Wire `limits.Budget` from `report` options into fill.

## Implementation checklist (future agent only)

- [x] Wire guard calls into fill and page creation.
- [x] Tests for each tripped limit.
- [x] Run `make check`.

## Acceptance criteria (observable)

- Exceeding a budget returns error `errors.Is`/`diag.As` code `limit`.
- Unlimited budget fills normally.

## Tests and fixtures to add (future agent only)

- `internal/fill/limits_test.go`.

## Failure cases / edge cases

- Exactly-at-limit boundary (allowed); zero-valued budget (unlimited).

## Definition of done

Enforced/tested; linked from architecture.

## References and downstream dependencies

- Features: `OPS-001`, `FND-009`.
- Unblocks: `TST-009-001`, `OPS-007`.
