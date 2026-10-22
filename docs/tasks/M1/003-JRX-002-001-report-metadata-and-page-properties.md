# JRX-002-001: Parse report metadata and page properties

- **Status:** Done
- **Feature:** [JRX-002](../../FEATURE_MATRIX.md#jrx--jrxml-syntaxversion-compatibility-and-compilation)
- **Milestone:** M1
- **Priority:** P0
- **Dependencies:** JRX-001-002
- **Blocked by:** none
- **Primary owner/package:** `internal/jrxml`
- **Likely touched files:** `internal/jrxml/report.go`, `internal/jrxml/report_test.go`
- **Collision risk:** Medium — adds a handler file next to others in the same package; uses the seam from `JRX-001-002`.
- **Parallel-safe with:** `JRX-003-001`, `JRX-004-001`, `JRX-010-001` (separate files, shared dispatcher — coordinate the one-line registration)
- **Estimated size:** S

## Objective

Parse page properties: size, orientation, margins, columns, column width/spacing,
and report name/metadata.

## Current state / prerequisites

Requires the decode skeleton and the `PRT-001-001` unit decision.

## Scope (in)

- `pageWidth`, `pageHeight`, `orientation`, `columnCount`, `columnWidth`,
  `columnSpacing`, all four margins, `name`, `language` (captured, unused).
- Convert units using `printdoc` unit helpers.
- Tests asserting parsed values equal fixture values.

## Out of scope

- Layout application of these values (`LAY-001-001`).
- Band content.

## Proposed contracts and integration points

Populates `Template.Page` (a `printdoc`-style geometry struct). Depends on
`PRT-001-001` types.

## Implementation checklist (future agent only)

- [x] Register a `jasperReport` attribute/property handler.
- [x] Parse and convert page properties.
- [x] Tests against a fixture template.
- [x] Run `make check`.

## Acceptance criteria (observable)

- Parsed page size/margins/columns equal the fixture's declared values.
- Defaults applied when attributes are absent.

## Tests and fixtures to add (future agent only)

- `testdata/jrxml/page-properties.jrxml` + `report_test.go`.

## Failure cases / edge cases

- Missing size → defaults (A4/Letter policy — document choice).
- Non-numeric values → diagnostic.

## Definition of done

Parsed correctly, tested; compatibility row `JRX-002` updated with the fixture.

## References and downstream dependencies

- Feature: `JRX-002`; unit: [`PRT-001-001`](../M0/PRT-001-001-print-document-ir-contract.md) (in M0).
- Unblocks: `JRX-008-001`, `LAY-001-001`.
