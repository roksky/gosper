# PRT-001-002: Print document model

- **Status:** Done
- **Feature:** [PRT-001](../../FEATURE_MATRIX.md#prt--print-document-and-storage)
- **Milestone:** M1
- **Priority:** P0
- **Dependencies:** PRT-001-001
- **Blocked by:** none
- **Primary owner/package:** `internal/printdoc`
- **Likely touched files:** `internal/printdoc/document.go`, `internal/printdoc/document_test.go`
- **Collision risk:** Medium — implements types agreed in `PRT-001-001`; `PRT-002-001` extends style handling.
- **Parallel-safe with:** `internal/text`, `internal/jrxml` tasks
- **Estimated size:** M

## Objective

Implement the concrete print document model: pages, text/image/shape elements,
page metadata (number, size, count) and iteration for exporters.

## Current state / prerequisites

Contract/types agreed in `PRT-001-001`.

## Scope (in)

- `Document`, `Page`, element structs, `Metadata`.
- Page number/size/count fields (filled by layout).
- Convenience iteration over pages/elements.
- Tests constructing/inspecting documents.

## Out of scope

- Rendering (`OUT-*`); serialization (`PRT-006`, Deferred).

## Proposed contracts and integration points

Consumed by `internal/exporter`; produced by `internal/layout`.

## Implementation checklist (future agent only)

- [x] Implement model + helpers.
- [x] Tests for construction and metadata.
- [x] Run `make check`.

## Acceptance criteria (observable)

- A document with 2 pages reports `PageCount == 2` and correct sizes.
- Element bounds round-trip exactly.

## Tests and fixtures to add (future agent only)

- `internal/printdoc/document_test.go`.

## Failure cases / edge cases

- Empty document (0 pages) handled by exporters.

## Definition of done

Implemented/tested; linked from architecture.

## References and downstream dependencies

- Feature: `PRT-001`, `PRT-007`.
- Unblocks: `LAY-001-001`, `OUT-001-001`, `PRT-002-001`.
