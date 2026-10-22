# OUT — Format exporters

**Architecture:** [`../architecture/exporters.md`](../architecture/exporters.md).

## Purpose

Render a print document to concrete formats. PDF first (visual), then XLSX,
HTML, CSV (data/grid).

## Technical expectations

- Exporters depend only on `internal/printdoc`.
- **PDF is visual/strict**; **XLSX/CSV are data-oriented** and do not preserve
  visual geometry the same way.
- Output streams to `io.Writer` (`OUT-009`).
- **Deterministic** output for identical input (`OUT-010`); injected
  timestamps in tests.
- Unsupported elements produce diagnostics, never silent drops.

## Features

| ID | Expectation | Acceptance test (planned) |
| --- | --- | --- |
| OUT-001 | Baseline PDF. | PDF golden, valid structure. |
| OUT-002 | PDF text/fonts. | Text golden. |
| OUT-003 | PDF vector/images. | Shape golden. |
| OUT-004 | PDF font embedding. | Embedded-font golden. |
| OUT-005 | PDF metadata/links. | Metadata/link test. |
| OUT-006 | XLSX. | XLSX golden. |
| OUT-007 | HTML. | HTML golden. |
| OUT-008 | CSV. | CSV exact match. |
| OUT-009 | Streaming writer. | Write to buffer/stream. |
| OUT-010 | Deterministic output. | Byte-equal repeat runs. |

## Tasks

`OUT-001-001`, `OUT-010-001` (M1).

## Risks

- PDF library choice (`OQ-5`); if a dependency is added, record an ADR and
  license review (`OPS-008`).
