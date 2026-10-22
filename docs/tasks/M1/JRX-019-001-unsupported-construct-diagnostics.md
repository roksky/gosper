# JRX-019-001: Unsupported-construct diagnostics

- **Status:** Not Started
- **Feature:** [JRX-019](../../FEATURE_MATRIX.md#jrx--jrxml-syntaxversion-compatibility-and-compilation)
- **Milestone:** M1
- **Priority:** P1
- **Dependencies:** JRX-001-002, FND-003-001
- **Blocked by:** none
- **Primary owner/package:** `internal/jrxml`
- **Likely touched files:** `internal/jrxml/unsupported.go`, `internal/jrxml/unsupported_test.go`
- **Collision risk:** Medium — registers a catch-all handler in the dispatcher; coordinate registration with other handler tasks.
- **Parallel-safe with:** `DAT-001-001`, `LAY-005-001`
- **Estimated size:** S

## Objective

Ensure unknown or not-yet-supported JRXML constructs produce a clear, coded
diagnostic with a source span instead of being silently ignored.

## Current state / prerequisites

Diagnostics contract (`FND-003-001`) and the decode skeleton must exist.

## Scope (in)

- A default handler for unrecognized elements/attributes that emits a
  `unsupported` diagnostic.
- Allowlist of known-but-deferred constructs (e.g. `chart`, `crosstab`) to
  produce `unsupported` rather than `parse`.
- Tests for an unknown element and a deferred component.

## Out of scope

- Implementing the unsupported constructs.
- Version normalization (`JRX-020`).

## Proposed contracts and integration points

Integrates with the `internal/jrxml` dispatcher; emits into `diag.List`.

## Implementation checklist (future agent only)

- [ ] Add catch-all diagnostics with element name + span.
- [ ] Test an unknown element yields `unsupported`/`parse` with span.
- [ ] Run `make check`.

## Acceptance criteria (observable)

- Unknown element ⇒ diagnostic containing the element name and a non-nil span.
- Deferred known construct ⇒ `unsupported` diagnostic, not a panic.

## Tests and fixtures to add (future agent only)

- `testdata/jrxml/unsupported.jrxml` + `unsupported_test.go`.

## Failure cases / edge cases

- Namespaced unknown elements; deeply nested unknowns (do not recurse
  unboundedly).

## Definition of done

Diagnostics emitted/tested; linked from architecture.

## References and downstream dependencies

- Feature: `JRX-019`.
- Unblocks: `JRX-020-001`, `JRX-022-001` fixture rows.
