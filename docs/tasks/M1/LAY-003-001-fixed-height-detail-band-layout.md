# LAY-003-001: Fixed-height detail band layout

- **Status:** Not Started
- **Feature:** [LAY-003](../../FEATURE_MATRIX.md#lay--layout-typography-pagination)
- **Milestone:** M1
- **Priority:** P0
- **Dependencies:** LAY-002-001, FIL-001-001
- **Blocked by:** none
- **Primary owner/package:** `internal/layout`
- **Likely touched files:** `internal/layout/detail.go`, `internal/layout/detail_test.go`
- **Collision risk:** Medium/High — the M1 end-to-end layout path; keep fixed-height only.
- **Parallel-safe with:** `internal/text` tasks, `internal/exporter` tasks
- **Estimated size:** M

## Objective

Lay out one detail band per filled record at a **fixed** height, crossing pages
when needed, producing a complete print document for the M1 slice.

## Current state / prerequisites

Geometry, positioning, and fill iteration exist.

## Scope (in)

- Emit `detail` band per record; include `title`/`pageHeader`/`pageFooter` if
  present.
- Fixed heights (no stretch); simple page overflow pagination.
- Text elements placed with measured widths (no wrapping).
- Tests: page count for N records, element positions on page 1 and page 2.

## Out of scope

- Wrapping, stretch, groups, columns, summary variables.

## Proposed contracts and integration points

Entry: `Layout(ctx, filled, metrics) (*printdoc.Document, error)`.

## Implementation checklist (future agent only)

- [ ] Implement the M1 layout loop.
- [ ] Tests for multi-page output.
- [ ] Run `make check`.

## Acceptance criteria (observable)

- 1 record fits page 1; enough records to overflow produce page 2 with correct
  coordinates.
- No panic on empty dataset.

## Tests and fixtures to add (future agent only)

- `testdata/jrxml/m1-basic.jrxml` + `internal/layout/detail_test.go`.

## Failure cases / edge cases

- Band taller than page; zero records; exactly-full page boundary.

## Definition of done

Implemented/tested; M1 walking skeleton reaches print IR.

## References and downstream dependencies

- Feature: `LAY-003`.
- Unblocks: `OUT-001-001`, `FND-001-002`.
