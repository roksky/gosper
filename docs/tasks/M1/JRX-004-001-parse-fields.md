# JRX-004-001: Parse report fields

- **Status:** Not Started
- **Feature:** [JRX-004](../../FEATURE_MATRIX.md#jrx--jrxml-syntaxversion-compatibility-and-compilation)
- **Milestone:** M1
- **Priority:** P0
- **Dependencies:** JRX-001-002
- **Blocked by:** none
- **Primary owner/package:** `internal/jrxml`
- **Likely touched files:** `internal/jrxml/fields.go`, `internal/jrxml/fields_test.go`
- **Collision risk:** Medium — handler file.
- **Parallel-safe with:** `JRX-003-001`, `JRX-010-001`
- **Estimated size:** S

## Objective

Parse `<field>` declarations (name, class, description) into the template model.

## Current state / prerequisites

Decode skeleton in place.

## Scope (in)

- `name`, `class`, optional `description`.
- Store field type for later conversion (`DAT-007`, `EXP-007`).
- Tests asserting the field list.

## Out of scope

- Data binding/conversion (`DAT-*`).
- Expression compilation.

## Proposed contracts and integration points

Adds `Template.Fields []Field{Name, Type string}`.

## Implementation checklist (future agent only)

- [ ] Handler for `<field>`.
- [ ] Tests with typed fields.
- [ ] Run `make check`.

## Acceptance criteria (observable)

- Field names/types match the fixture.

## Tests and fixtures to add (future agent only)

- `testdata/jrxml/fields.jrxml` + `fields_test.go`.

## Failure cases / edge cases

- Unknown `class` values: store verbatim; mapping deferred to `DAT-007`.

## Definition of done

Parsed/tested; compatibility row `JRX-004` updated.

## References and downstream dependencies

- Feature: `JRX-004`.
- Unblocks: `EXP-001-002`, `DAT-001-001`.
