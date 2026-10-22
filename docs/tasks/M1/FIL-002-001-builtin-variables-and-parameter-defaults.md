# FIL-002-001: Built-in variables and parameter defaults

- **Status:** Not Started
- **Feature:** [FIL-002](../../FEATURE_MATRIX.md#fil--fill-engine-and-report-semantics)
- **Milestone:** M1
- **Priority:** P1
- **Dependencies:** FIL-001-001
- **Blocked by:** none
- **Primary owner/package:** `internal/fill`
- **Likely touched files:** `internal/fill/builtins.go`, `internal/fill/builtins_test.go`
- **Collision risk:** Medium — extends fill state; coordinate struct fields with `FIL-001-001`.
- **Parallel-safe with:** `LAY-005-001`, `EXP-006-001`
- **Estimated size:** S

## Objective

Provide `REPORT_COUNT` and apply parameter default expressions when a parameter
is not supplied.

## Current state / prerequisites

Fill loop and expression evaluation exist.

## Scope in/out

- In: `REPORT_COUNT`; parameter default evaluation using `EXP-003-001`.
- Out: `PAGE_NUMBER`/`COLUMN_NUMBER` (require layout; `FIL-010`).

## Proposed contracts and integration points

Expose built-ins through the `Symbols` implementation used by fill.

## Implementation checklist (future agent only)

- [ ] Implement `REPORT_COUNT` increment per record.
- [ ] Evaluate parameter defaults when missing.
- [ ] Tests for both.
- [ ] Run `make check`.

## Acceptance criteria (observable)

- `$V{REPORT_COUNT}` reflects record count.
- A defaulted parameter takes its default value when absent.

## Tests and fixtures to add (future agent only)

- `internal/fill/builtins_test.go`.

## Failure cases / edge cases

- Default expression referencing another parameter (evaluation order).

## Definition of done

Implemented/tested; documented.

## References and downstream dependencies

- Feature: `FIL-002`.
- Unblocks: `FIL-010`.
