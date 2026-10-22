# JRX-008-001: Parse core bands

- **Status:** Not Started
- **Feature:** [JRX-008](../../FEATURE_MATRIX.md#jrx--jrxml-syntaxversion-compatibility-and-compilation)
- **Milestone:** M1
- **Priority:** P0
- **Dependencies:** JRX-002-001, JRX-010-001
- **Blocked by:** none
- **Primary owner/package:** `internal/jrxml`
- **Likely touched files:** `internal/jrxml/bands.go`, `internal/jrxml/bands_test.go`
- **Collision risk:** High — assembles elements into bands; central to the model. Land after element/attr tasks to avoid churn.
- **Parallel-safe with:** `DAT-001-001`, `EXP-003-001`, `LAY-005-001`
- **Estimated size:** M

## Objective

Assemble the core bands (`title`, `pageHeader`, `columnHeader`, `detail`,
`columnFooter`, `pageFooter`, `summary`) in declared order with their elements
and `height`.

## Current state / prerequisites

Requires page properties and element parsing.

## Scope (in)

- Band ordering and `height` attribute.
- Elements attached to bands.
- Tests asserting band order and contained elements.

## Out of scope

- Advanced bands (`JRX-009`).
- Group bands (`JRX-015`).

## Proposed contracts and integration points

```go
// Proposed — non-binding.
type Band struct { Kind string; Height Unit; Elements []Element }
```

## Implementation checklist (future agent only)

- [ ] Handler mapping known band element names to `Band`.
- [ ] Preserve declared order.
- [ ] Tests for a multi-band fixture.
- [ ] Run `make check`.

## Acceptance criteria (observable)

- Band order and heights match the fixture.
- Unknown band names produce a diagnostic.

## Tests and fixtures to add (future agent only)

- `testdata/jrxml/core-bands.jrxml` + `bands_test.go`.

## Failure cases / edge cases

- Missing `detail` band (allowed); duplicate bands (diagnostic).

## Definition of done

Parsed/tested; compatibility rows for core bands updated.

## References and downstream dependencies

- Feature: `JRX-008`.
- Unblocks: `FIL-001-001`, `LAY-003-001`.
