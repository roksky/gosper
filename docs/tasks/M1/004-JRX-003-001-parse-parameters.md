# JRX-003-001: Parse report parameters

- **Status:** Done
- **Feature:** [JRX-003](../../FEATURE_MATRIX.md#jrx--jrxml-syntaxversion-compatibility-and-compilation)
- **Milestone:** M1
- **Priority:** P0
- **Dependencies:** JRX-001-002
- **Blocked by:** none
- **Primary owner/package:** `internal/jrxml`
- **Likely touched files:** `internal/jrxml/parameters.go`, `internal/jrxml/parameters_test.go`
- **Collision risk:** Medium — handler file using the shared seam.
- **Parallel-safe with:** `JRX-004-001`, `JRX-010-001`
- **Estimated size:** S

## Objective

Parse `<parameter>` declarations (name, class type, default expression, prompt
flags ignored for now) into the template model.

## Current state / prerequisites

Decode skeleton in place.

## Scope (in)

- `name`, `class`, `nestedType`, and `<defaultValueExpression>` text.
- Store raw expression text (compilation happens in `internal/compiler`).
- Tests asserting parameter list and defaults.

## Out of scope

- Evaluating default expressions (`FIL-002-001`).
- Compiler binding.

## Proposed contracts and integration points

Adds `Template.Parameters []Parameter{Name, Type, DefaultExpr string}`.

## Implementation checklist (future agent only)

- [x] Handler for `<parameter>` and its default expression.
- [x] Tests with a fixture containing typed and untyped parameters.
- [x] Run `make check`.

## Acceptance criteria (observable)

- Parameter names/types/defaults match the fixture.
- Missing `class` defaults to `java.lang.String` (documented).

## Tests and fixtures to add (future agent only)

- `testdata/jrxml/parameters.jrxml` + `parameters_test.go`.

## Failure cases / edge cases

- Duplicate parameter names → diagnostic.

## Definition of done

Parsed/tested; compatibility row `JRX-003` updated.

## References and downstream dependencies

- Feature: `JRX-003`.
- Unblocks: `FIL-002-001`, `EXP-001-002`.
