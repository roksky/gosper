# EXP-006-001: Null semantics

- **Status:** Done
- **Feature:** [EXP-006](../../FEATURE_MATRIX.md#exp--expression-language-and-type-semantics)
- **Milestone:** M1
- **Priority:** P0
- **Dependencies:** EXP-003-001
- **Blocked by:** none
- **Primary owner/package:** `internal/expression`
- **Likely touched files:** `internal/expression/null.go`, `internal/expression/null_test.go`
- **Collision risk:** Medium — same package; keep null rules in this file.
- **Parallel-safe with:** `DAT-001-001`, `LAY-005-001`
- **Estimated size:** S

## Objective

Define and implement null propagation across operators and comparisons so
expressions behave reproducibly (Jasper-like) for missing/blank values.

## Current state / prerequisites

Evaluator exists (`EXP-003-001`).

## Scope (in)

- Null in arithmetic ⇒ null; comparison with null ⇒ defined result (decide:
  null or false) and document.
- Boolean short-circuit interaction with null.
- A decision table encoded as tests.

## Out of scope

- Full coercion (`EXP-007`); functions (`EXP-009`).

## Proposed contracts and integration points

Central helpers used by the evaluator; no new exported API.

## Implementation checklist (future agent only)

- [x] Implement null rules per decision table.
- [x] Encode the table as tests.
- [x] Document the choices in the package/architecture doc.
- [x] Run `make check`.

## Acceptance criteria (observable)

- Table-driven tests pass for null in arithmetic, comparison, boolean.
- Documented behavior matches JasperReports conventions or records a deviation.

## Tests and fixtures to add (future agent only)

- `internal/expression/null_test.go`.

## Failure cases / edge cases

- Null in division (avoid panic); null equality (`null == null`).

## Definition of done

Implemented/tested; documented; compatibility note if deviating.

## References and downstream dependencies

- Feature: `EXP-006`.
- Unblocks: `EXP-007`, `EXP-009`.
