# Feature domains

Domain documents describe **technical expectations and feature-specific
acceptance tests**. They are **not** status registers — canonical status lives
in [`../FEATURE_MATRIX.md`](../FEATURE_MATRIX.md).

| Domain | Scope | Document |
| --- | --- | --- |
| FND | Foundation, public API, CLI, infrastructure | [`FND.md`](FND.md) |
| JRX | JRXML syntax/version compatibility and compilation | [`JRX.md`](JRX.md) |
| EXP | Expression language and type semantics | [`EXP.md`](EXP.md) |
| DAT | Data sources and queries | [`DAT.md`](DAT.md) |
| FIL | Fill engine, variables, groups, subreports | [`FIL.md`](FIL.md) |
| LAY | Layout, text, pagination, images, shapes | [`LAY.md`](LAY.md) |
| PRT | Print document and storage | [`PRT.md`](PRT.md) |
| OUT | Format exporters | [`OUT.md`](OUT.md) |
| TST | Test infrastructure and correctness suites | [`TST.md`](TST.md) |
| OPS | Security, performance, diagnostics, integration | [`OPS.md`](OPS.md) |
| EXT | Advanced extensions (charts, crosstabs, service, 7.x) | this file, below |

## EXT — Advanced extensions (deferred)

Tracked only in [`../FEATURE_MATRIX.md#ext--advanced-extensions-deferred`](../FEATURE_MATRIX.md).
These are explicitly out of core v1: `EXT-001` charts, `EXT-002` crosstabs,
`EXT-003` barcodes, `EXT-004` PDF/A, `EXT-005` accessibility/tagged PDF,
`EXT-006` HTTP service, `EXT-007` gRPC service, `EXT-008` JRXML 7.x, `EXT-009`
designer integration. They require an explicit scheduling decision before any
task is authored.

## Conventions

- Each domain doc lists its features with **expectation** and **acceptance
  test** columns.
- IDs are stable and never recycled; see
  [`../CONTRIBUTING_TASKS.md`](../CONTRIBUTING_TASKS.md).
- Acceptance tests reference planned fixtures; none run yet.
