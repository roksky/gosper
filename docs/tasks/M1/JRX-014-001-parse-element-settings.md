# JRX-014-001: Parse common element settings

- **Status:** Not Started
- **Feature:** [JRX-014](../../FEATURE_MATRIX.md#jrx--jrxml-syntaxversion-compatibility-and-compilation)
- **Milestone:** M1
- **Priority:** P0
- **Dependencies:** JRX-010-001
- **Blocked by:** none
- **Primary owner/package:** `internal/jrxml`
- **Likely touched files:** `internal/jrxml/element_attrs.go`, `internal/jrxml/element_attrs_test.go`
- **Collision risk:** Medium — owns a separate file (`element_attrs.go`) from `elements.go`; agree the `Element` struct fields with `JRX-010-001` first.
- **Parallel-safe with:** `JRX-003-001`, `JRX-004-001`
- **Estimated size:** S

## Objective

Parse common element attributes: `x`, `y`, `width`, `height`, `isStretchWith
Overflow`, `printWhenExpression`, `evaluationTime`, `evaluationGroup`, `uuid`
(ignored).

## Current state / prerequisites

Element content parsing (`JRX-010-001`) defines `Element`.

## Scope (in)

- Numeric position/size attributes into `Unit`.
- Expression attributes stored as text (`printWhenExpression`).
- `evaluationTime` enum captured.
- Tests asserting attribute values.

## Out of scope

- Applying stretch/printWhen (`FIL-008`, `LAY-004`).

## Proposed contracts and integration points

Adds fields to `Element`: coordinates/size, `Stretch bool`,
`PrintWhenExpr string`, `EvalTime string`.

## Implementation checklist (future agent only)

- [ ] Parse/convert attributes.
- [ ] Tests for defaults and explicit values.
- [ ] Run `make check`.

## Acceptance criteria (observable)

- Coordinates/size in `Unit` equal the fixture values.
- `evaluationTime` defaults to `Now` when absent.

## Tests and fixtures to add (future agent only)

- `testdata/jrxml/element-settings.jrxml` + `element_attrs_test.go`.

## Failure cases / edge cases

- Negative sizes; unknown `evaluationTime` value → diagnostic.

## Definition of done

Parsed/tested; compatibility row `JRX-014` updated.

## References and downstream dependencies

- Feature: `JRX-014`.
- Unblocks: `LAY-002-001`, `FIL-008`.
