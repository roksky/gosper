# LAY — Layout, typography, pagination

**Architecture:** [`../architecture/layout-and-pagination.md`](../architecture/layout-and-pagination.md).

## Purpose

Turn filled data into a laid-out, paginated print document: geometry, band
placement, text measurement, dynamic sizing, page breaking.

## Technical expectations

- Absolute coordinates within the printable area; internal unit pending
  (`OQ-2`).
- **Deterministic measurement** for fixed fonts/locale (`LAY-015`).
- Band split types and overflow continuation (`LAY-010`).
- Column/page breaks, keep-together, orphans/widows (`LAY-011`, `LAY-016`).
- Text: wrapping, clipping, alignment, rotation (`LAY-008`); shaping/Unicode
  later (`LAY-007`).

## Features

| ID | Expectation | Acceptance test (planned) |
| --- | --- | --- |
| LAY-001 | Page geometry/band placement. | Coordinate test. |
| LAY-002 | Absolute positioning. | Position golden. |
| LAY-003 | Fixed-height detail layout. | Layout golden. |
| LAY-004 | Relative stretching. | Stretch fixture. |
| LAY-005 | Font resolver. | Resolve bundled fonts. |
| LAY-006 | Text measurement. | Metrics test. |
| LAY-007 | Unicode/shaping/fallback. | Unicode fixture. |
| LAY-008 | Wrap/clip/align/rotate. | Wrap fixture. |
| LAY-009 | Styled text/markup subset. | Markup fixture. |
| LAY-010 | Variable height/split types. | Split fixture. |
| LAY-011 | Column/page breaks/keep-together. | Break fixture. |
| LAY-012 | Borders/bg/colors/transparency. | Style golden. |
| LAY-013 | Image scaling. | Image golden. |
| LAY-014 | Shape geometry. | Shape golden. |
| LAY-015 | Deterministic measurement. | Repeat-run equality. |
| LAY-016 | Orphans/widows/page reservation. | Orphan fixture. |
| LAY-017 | Bidirectional text. | Bidi fixture. |
| LAY-018 | Hyperlink targets. | Link golden. |

## Tasks

`LAY-001-001`, `LAY-002-001`, `LAY-003-001`, `LAY-005-001`, `LAY-006-001` (M1);
`LAY-008-001` (M2).

## Risks

- Pagination is the highest-risk area; couple golden + semantic + page-count
  assertions (`TST-004`).
