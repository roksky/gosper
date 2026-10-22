# JRX-010-001: Parse basic elements (staticText, textField)

- **Status:** Not Started
- **Feature:** [JRX-010](../../FEATURE_MATRIX.md#jrx--jrxml-syntaxversion-compatibility-and-compilation)
- **Milestone:** M1
- **Priority:** P0
- **Dependencies:** JRX-001-002
- **Blocked by:** none
- **Primary owner/package:** `internal/jrxml`
- **Likely touched files:** `internal/jrxml/elements.go`, `internal/jrxml/elements_test.go`
- **Collision risk:** Medium/High — element parsing is extended by later tasks (shapes, images). Keep the M1 element model minimal and additive.
- **Parallel-safe with:** `JRX-003-001`, `JRX-004-001`
- **Estimated size:** M

## Objective

Parse `<staticText>` (with `<text>`) and `<textField>` (with
`<textFieldExpression>`, `pattern`, `isBlankWhenNull`) into element nodes.

## Current state / prerequisites

Decode skeleton and `JRX-014-001` both concern elements; coordinate: this task
owns textual content, `JRX-014-001` owns common element attributes. Split files:
`elements.go` (content) vs `element_attrs.go` (attributes).

## Scope (in)

- Text content, lazily evaluated expression text (stored, not evaluated).
- `pattern`, `isBlankWhenNull`, `textFieldExpression`.
- Tests asserting element values.

## Out of scope

- Layout/measurement.
- Shapes/images.

## Proposed contracts and integration points

```go
// Proposed — non-binding.
type Element struct { Kind string; X, Y, W, H Unit; Text, Expr, Pattern string }
```

Integration: consumed by layout after fill binds values.

## Implementation checklist (future agent only)

- [ ] Parse text/static content and expression fields.
- [ ] Register with the band/element dispatcher.
- [ ] Tests for static and expression text.
- [ ] Run `make check`.

## Acceptance criteria (observable)

- A fixture's static text and expression text are captured exactly.
- Unknown child elements produce a diagnostic, not silence.

## Tests and fixtures to add (future agent only)

- `testdata/jrxml/basic-elements.jrxml` + `elements_test.go`.

## Failure cases / edge cases

- CDATA content; empty expression; whitespace handling.

## Definition of done

Parsed/tested; compatibility rows `JRX-010` updated.

## References and downstream dependencies

- Feature: `JRX-010`.
- Unblocks: `JRX-008-001`, `JRX-014-001`, `LAY-003-001`.
