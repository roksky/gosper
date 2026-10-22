# JRX-001-002: JRXML decode skeleton and version detection

- **Status:** Done
- **Feature:** [JRX-001](../../FEATURE_MATRIX.md#jrx--jrxml-syntaxversion-compatibility-and-compilation)
- **Milestone:** M1
- **Priority:** P0
- **Dependencies:** JRX-001-001
- **Blocked by:** none
- **Primary owner/package:** `internal/jrxml`
- **Likely touched files:** `internal/jrxml/parse.go`, `internal/jrxml/parse_test.go`
- **Collision risk:** High — central file for the parser; later `internal/jrxml` tasks extend it. Keep the root decoding minimal and register extension points (handlers) to reduce conflicts.
- **Parallel-safe with:** `DAT-001-001`, `EXP-003-001`, `LAY-005-001` (different packages)
- **Estimated size:** M

## Objective

Implement the parser entry point: decode the root, detect/record the dialect,
and expose a handler-based structure so element/band/field parsing lands in
separate files.

## Current state / prerequisites

`JRX-001-001` defined `Parse`, `Template`, `Dialect`.

## Scope (in)

- Wire `encoding/xml` decoding of the root element and report attributes.
- Record version/dialect; warn on unknown.
- A registration seam so `JRX-002-001`, `JRX-003-001`, `JRX-004-001`,
  `JRX-010-001`, `JRX-008-001` add parsing without editing one another's code.

## Out of scope

- Any child-element semantics beyond the skeleton.

## Proposed contracts and integration points

Prefer `xml.Decoder` token stream with a dispatcher, or per-element
`UnmarshalXML`, so files are small and independently testable.

## Implementation checklist (future agent only)

- [x] Implement `Parse` over a token stream.
- [x] Emit diagnostics for malformed XML.
- [x] Add the dispatcher/handler seam and document it.
- [x] Tests for root detection and malformed input.
- [x] Run `make check`.

## Acceptance criteria (observable)

- Parsing a minimal report returns a `Template` with correct `Dialect`.
- Malformed XML yields a `parse` diagnostic; no panic.
- A test demonstrates adding a handler without modifying existing files.

## Tests and fixtures to add (future agent only)

- `internal/jrxml/parse_test.go`; first fixtures under `testdata/jrxml/`.

## Failure cases / edge cases

- BOM, namespaces, unknown root element.

## Definition of done

Skeleton parses/tests; handler seam documented; downstream tasks can add
parsers in isolations.

## References and downstream dependencies

- Feature: `JRX-001`; tasks `JRX-001-001`.
- Unblocks: `JRX-002-001`, `JRX-003-001`, `JRX-004-001`, `JRX-010-001`,
  `JRX-019-001`.
