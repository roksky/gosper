# Task index

> **Canonical task-level state.** Status, ordering, dependencies and
> ownership/reservation live here. Feature-level status lives in
> [`../FEATURE_MATRIX.md`](../FEATURE_MATRIX.md). All tasks currently start
> `Not Started`; none are claimed.

## Status vocabulary

`Not Started` → `Ready` → `In Progress` → `In Review` → `Done`, plus `Blocked`,
`Deferred`, `Cancelled`. See
[`../CONTRIBUTING_TASKS.md`](../CONTRIBUTING_TASKS.md).

## M0 — Baseline and contracts

All M0 tasks are independent (no inter-task dependencies) and touch different
packages/docs, so they are broadly parallel-safe.

| Task | Title | Feature | Status | Deps | Package / files | File |
| --- | --- | --- | --- | --- | --- | --- |
| FND-003-001 | Diagnostic and error code model | FND-003 | Done | none | `internal/diag` | [file](M0/FND-003-001-diagnostic-error-code-model.md) |
| FND-009-001 | Context and resource-budget contract | FND-009 | Done | none | `internal/limits` | [file](M0/FND-009-001-context-and-resource-budget-contract.md) |
| FND-001-001 | Public SDK surface contract | FND-001 | Done | FND-003-001, FND-009-001 | `report` | [file](M0/FND-001-001-public-sdk-surface-contract.md) |
| JRX-001-001 | Parser I/O contract and version detection | JRX-001 | Done | FND-003-001 | `internal/jrxml` | [file](M0/JRX-001-001-parser-io-contract-and-version-detection.md) |
| DAT-002-001 | DataSource cursor interface contract | DAT-002 | Done | none | `internal/datasource` | [file](M0/DAT-002-001-datasource-cursor-interface-contract.md) |
| EXP-001-001 | Expression symbol/scope contract | EXP-001 | Done | none | `internal/expression` | [file](M0/EXP-001-001-expression-symbol-scope-contract.md) |
| EXP-019-001 | Expression sandbox policy | EXP-019 | Done | none | `internal/expression`, `docs/` | [file](M0/EXP-019-001-expression-sandbox-policy.md) |
| PRT-001-001 | Print document IR contract | PRT-001 | Done | none | `internal/printdoc` | [file](M0/PRT-001-001-print-document-ir-contract.md) |
| TST-001-001 | Deterministic test environment helpers | TST-001 | Done | none | `internal/testutil` | [file](M0/TST-001-001-deterministic-test-environment-helpers.md) |
| TST-001-002 | Fixture conventions and golden workflow | TST-001 | Done | TST-001-001 | `testdata/`, `docs/TEST_STRATEGY.md` | [file](M0/TST-001-002-fixture-conventions-and-golden-workflow.md) |
| OPS-008-001 | Dependency/license review policy and SBOM | OPS-008 | Done | none | `docs/`, `.github/` | [file](M0/OPS-008-001-dependency-and-license-review-policy.md) |
| JRX-022-001 | Compatibility matrix and fixture registry | JRX-022 | Done | none | `docs/COMPATIBILITY.md`, `testdata/` | [file](M0/JRX-022-001-compatibility-matrix-and-fixture-registry.md) |

## M1 — First vertical slice (walking skeleton)

The numeric filename prefixes give a sequential pick order through M1 and M2;
each listed task follows its task dependencies. Before picking `027-TST-002-001`,
complete the M0 prerequisite `TST-001-002`, which is still Not Started. Task
`023-OUT-001-001` also requires the PDF library choice (OQ-5) recorded in its
"Blocked by" field.

| Task | Title | Feature | Status | Deps | Package / files | File |
| --- | --- | --- | --- | --- | --- | --- |
| FND-002-001 | Template loading and resource resolution | FND-002 | Not Started | none | `internal/resources` | [file](M1/001-FND-002-001-template-loading-and-resource-resolution.md) |
| JRX-001-002 | Decode skeleton and version detection | JRX-001 | Not Started | JRX-001-001 | `internal/jrxml` | [file](M1/002-JRX-001-002-jrxml-decode-skeleton-and-version-detection.md) |
| JRX-002-001 | Report metadata and page properties | JRX-002 | Not Started | JRX-001-002 | `internal/jrxml` | [file](M1/003-JRX-002-001-report-metadata-and-page-properties.md) |
| JRX-003-001 | Parse parameters | JRX-003 | Not Started | JRX-001-002 | `internal/jrxml` | [file](M1/004-JRX-003-001-parse-parameters.md) |
| JRX-004-001 | Parse fields | JRX-004 | Not Started | JRX-001-002 | `internal/jrxml` | [file](M1/005-JRX-004-001-parse-fields.md) |
| JRX-010-001 | Parse basic elements | JRX-010 | Not Started | JRX-001-002 | `internal/jrxml` | [file](M1/006-JRX-010-001-parse-basic-elements.md) |
| JRX-008-001 | Parse core bands | JRX-008 | Not Started | JRX-002-001, JRX-010-001 | `internal/jrxml` | [file](M1/007-JRX-008-001-parse-core-bands.md) |
| JRX-014-001 | Parse element settings | JRX-014 | Not Started | JRX-010-001 | `internal/jrxml` | [file](M1/008-JRX-014-001-parse-element-settings.md) |
| JRX-019-001 | Unsupported-construct diagnostics | JRX-019 | Not Started | JRX-001-002, FND-003-001 | `internal/jrxml` | [file](M1/009-JRX-019-001-unsupported-construct-diagnostics.md) |
| EXP-001-002 | Symbol resolution | EXP-001 | Not Started | EXP-001-001 | `internal/expression` | [file](M1/010-EXP-001-002-symbol-resolution.md) |
| EXP-003-001 | Expression parser and evaluator | EXP-003 | Not Started | EXP-001-002, EXP-019-001 | `internal/expression` | [file](M1/011-EXP-003-001-expression-parser-and-evaluator.md) |
| EXP-006-001 | Null semantics | EXP-006 | Not Started | EXP-003-001 | `internal/expression` | [file](M1/012-EXP-006-001-null-semantics.md) |
| DAT-001-001 | Slice data source | DAT-001 | Not Started | DAT-002-001 | `internal/datasource` | [file](M1/013-DAT-001-001-slice-data-source.md) |
| FIL-001-001 | Fill record iteration | FIL-001 | Not Started | DAT-002-001, FND-001-001 | `internal/fill` | [file](M1/014-FIL-001-001-fill-record-iteration.md) |
| FIL-002-001 | Built-in variables and parameter defaults | FIL-002 | Not Started | FIL-001-001 | `internal/fill` | [file](M1/015-FIL-002-001-builtin-variables-and-parameter-defaults.md) |
| LAY-005-001 | Font resolver | LAY-005 | Not Started | none | `internal/text` | [file](M1/016-LAY-005-001-font-resolver.md) |
| LAY-006-001 | Text measurement baseline | LAY-006 | Not Started | LAY-005-001 | `internal/text` | [file](M1/017-LAY-006-001-text-measurement-baseline.md) |
| PRT-001-002 | Print document model | PRT-001 | Not Started | PRT-001-001 | `internal/printdoc` | [file](M1/018-PRT-001-002-print-document-model.md) |
| LAY-001-001 | Page geometry and band placement | LAY-001 | Not Started | PRT-001-002 | `internal/layout` | [file](M1/019-LAY-001-001-page-geometry-and-band-placement.md) |
| LAY-002-001 | Absolute positioning | LAY-002 | Not Started | LAY-001-001 | `internal/layout` | [file](M1/020-LAY-002-001-absolute-positioning.md) |
| LAY-003-001 | Fixed-height detail band layout | LAY-003 | Not Started | LAY-002-001, FIL-001-001 | `internal/layout` | [file](M1/021-LAY-003-001-fixed-height-detail-band-layout.md) |
| PRT-002-001 | Resolved style data | PRT-002 | Not Started | PRT-001-002 | `internal/printdoc` | [file](M1/022-PRT-002-001-resolved-style-data.md) |
| OUT-001-001 | Baseline PDF exporter | OUT-001 | Not Started | PRT-001-002, PRT-002-001 | `internal/exporter` | [file](M1/023-OUT-001-001-baseline-pdf-exporter.md) |
| OUT-010-001 | Deterministic output tests | OUT-010 | Not Started | OUT-001-001 | `internal/exporter` | [file](M1/024-OUT-010-001-deterministic-output-tests.md) |
| FND-001-002 | Public SDK walking skeleton | FND-001 | Not Started | PRT-001-002, FIL-001-001, OUT-001-001 | `report` (proposed) | [file](M1/025-FND-001-002-public-sdk-walking-skeleton.md) |
| FND-006-001 | Build and version metadata | FND-006 | Not Started | none | `report`, `cmd/gosper` | [file](M1/026-FND-006-001-build-and-version-metadata.md) |
| TST-002-001 | Golden fixture harness | TST-002 | Not Started | TST-001-001, TST-001-002 | `internal/testutil` | [file](M1/027-TST-002-001-golden-fixture-harness.md) |
| OPS-001-001 | Resource limits enforcement | OPS-001 | Not Started | FND-009-001 | `internal/limits`, `internal/fill` | [file](M1/028-OPS-001-001-resource-limits-enforcement.md) |
| TST-009-001 | Race, cancellation, leak tests | TST-009 | Not Started | FIL-001-001, OPS-001-001 | `internal/fill` tests | [file](M1/029-TST-009-001-race-cancellation-leak-tests.md) |

## M2 — Early next tasks (contracts clear)

| Task | Title | Feature | Status | Deps | Package / files | File |
| --- | --- | --- | --- | --- | --- | --- |
| DAT-003-001 | JSON data source | DAT-003 | Not Started | DAT-002-001 | `internal/datasource` | [file](M2/030-DAT-003-001-json-data-source.md) |
| DAT-004-001 | CSV data source | DAT-004 | Not Started | DAT-002-001 | `internal/datasource` | [file](M2/031-DAT-004-001-csv-data-source.md) |
| LAY-008-001 | Text wrapping and alignment | LAY-008 | Not Started | LAY-006-001 | `internal/text`, `internal/layout` | [file](M2/032-LAY-008-001-text-wrapping-and-alignment.md) |
| JRX-020-001 | Version-aware normalization | JRX-020 | Not Started | JRX-001-002, JRX-019-001 | `internal/jrxml` | [file](M2/033-JRX-020-001-version-aware-normalization.md) |

## Backlog (M3+) — tasks not yet authored

Features for M3–M6 are tracked in
[`../FEATURE_MATRIX.md`](../FEATURE_MATRIX.md) and are intentionally **not**
broken into tasks yet, because their contracts depend on M1/M2 outcomes. Author
tasks for a milestone only once its prerequisite contracts are `Done`.

## Critical path

`FND-003-001` → `JRX-001-001` → `JRX-001-002` → `JRX-008-001` /
`JRX-010-001` → `PRT-001-002` → `FIL-001-001` → `LAY-003-001` →
`OUT-001-001` → `FND-001-002` → (M2) `DAT-003-001` / `LAY-008-001` → (M3,
unauthored).

## Parallel-safe queues

- **Queue A (M0 contracts):** `FND-003-001`, `FND-009-001`, `DAT-002-001`,
  `EXP-001-001`, `EXP-019-001`, `PRT-001-001`, `TST-001-001`, `OPS-008-001`,
  `JRX-022-001` — disjoint packages/docs.
- **Queue B (after `JRX-001-001`):** `JRX-002-001`, `JRX-003-001`,
  `JRX-004-001`, `JRX-010-001` — same package, so serialize *file*-level; run
  one at a time or coordinate hunks.
- **Queue C (independent):** `LAY-005-001` and `FND-006-001` and
  `FND-002-001` touch different packages and can run alongside B.
