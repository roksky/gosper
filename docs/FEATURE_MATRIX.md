# Feature matrix

> **Canonical feature register.** This is the single source of truth for feature
> identity, scope, priority, milestone, status and dependencies. It is derived
> from evidenced repository state plus planned work; **"Not Started" is the
> default**, and no feature is `Done` without observable verification. See
> [`CONTRIBUTING_TASKS.md`](CONTRIBUTING_TASKS.md) for how to update this file.

## Legend

- **Priority:** `P0` first useful capability / correctness-safety · `P1`
  production v1 · `P2` next-release expansion · `P3` optional/experimental.
- **Milestone:** `M0`…`M6` per [`ROADMAP.md`](ROADMAP.md).
- **Status:** `Not Started` → `Ready` → `In Progress` → `In Review` → `Done`,
  plus `Blocked`, `Deferred`, `Cancelled`.
- **Dependencies:** `—` means none; `TBD` means an open question (see
  [Open questions](#open-questions)).
- **Specification:** technical expectations live in the domain doc.
- **Tasks:** PR-sized task specs; `INDEX` means the detailed task is not yet
  authored (see [`tasks/INDEX.md`](tasks/INDEX.md)).

Statuses shown are a snapshot; a `Done`/`In Progress` row cites evidence in the
final column.

## FND — Foundation, public API, CLI, infrastructure

Spec: [`features/FND.md`](features/FND.md) · Architect: [`architecture/package-boundaries.md`](architecture/package-boundaries.md)

| ID | Feature | Domain | Priority | Milestone | Status | Dependencies | Specification | Tasks | Verification evidence / note |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| FND-001 | Public SDK surface (compile/fill/export + convenience) | FND | P0 | M1 | In Progress | FND-003, FND-009 | [FND](features/FND.md) | [FND-001-001](tasks/M0/FND-001-001-public-sdk-surface-contract.md), [FND-001-002](tasks/M1/025-FND-001-002-public-sdk-walking-skeleton.md) | contract stubs in `report` (`api_test.go`); behavior pending FND-001-002 |
| FND-002 | Template loading and resource resolution | FND | P1 | M1 | **Done** | — | [FND](features/FND.md) | [FND-002-001](tasks/M1/001-FND-002-001-template-loading-and-resource-resolution.md) | `internal/resources/resolver_test.go`; branch `task/FND-002-001-template-loading-and-resource-resolution` |
| FND-003 | Error model and diagnostic codes with source spans | FND | P0 | M0 | **Done** | — | [FND](features/FND.md) | [FND-003-001](tasks/M0/FND-003-001-diagnostic-error-code-model.md) | `internal/diag/diag_test.go`; branch `task/M0-contracts` |
| FND-004 | Compilation cache with versioned keys | FND | P2 | M3 | Not Started | FND-001, JRX-020 | [FND](features/FND.md) | INDEX | planned: invalidation test |
| FND-005 | CLI report commands (compile/fill/export) | FND | P2 | M4 | Not Started | FND-001, OUT-001 | [FND](features/FND.md) | INDEX | planned: CLI integration test |
| FND-006 | Build/version metadata surfaced via API | FND | P2 | M1 | Not Started | — | [FND](features/FND.md) | [FND-006-001](tasks/M1/026-FND-006-001-build-and-version-metadata.md) | planned: version test |
| FND-007 | Logging and telemetry hooks | FND | P3 | M5 | Not Started | — | [FND](features/FND.md) | INDEX | planned: hook invocation test |
| FND-008 | Optional server transport module (deferred) | FND | P3 | M6 | Not Started | FND-001 | [FND](features/FND.md) | INDEX | planned: adapter smoke test |
| FND-009 | Context cancellation and resource-budget API | FND | P0 | M0 | **Done** | — | [FND](features/FND.md) | [FND-009-001](tasks/M0/FND-009-001-context-and-resource-budget-contract.md) | `internal/limits/limits_test.go`; branch `task/M0-contracts` |
| FND-010 | Public API stability and compatibility policy | FND | P1 | M5 | Not Started | FND-001 | [FND](features/FND.md) | INDEX | planned: doc + lint policy |
| FND-011 | CLI metadata (help/version/unknown-command) | FND | P1 | M0 | **Done** | — | [FND](features/FND.md) | — | `cmd/gosper/main_test.go`; commit `7a787f4` |

## JRX — JRXML syntax/version compatibility and compilation

Spec: [`features/JRX.md`](features/JRX.md) · Compat: [`COMPATIBILITY.md`](COMPATIBILITY.md) · Architect: [`architecture/intermediate-representations.md`](architecture/intermediate-representations.md)

| ID | Feature | Domain | Priority | Milestone | Status | Dependencies | Specification | Tasks | Verification evidence / note |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| JRX-001 | JRXML parsing and version detection | JRX | P0 | M1 | **Done** | FND-003 | [JRX](features/JRX.md) | [JRX-001-001](tasks/M0/JRX-001-001-parser-io-contract-and-version-detection.md), [JRX-001-002](tasks/M1/002-JRX-001-002-jrxml-decode-skeleton-and-version-detection.md) | `internal/jrxml/parse_test.go`; dispatch seam; branch `task/JRX-m1-parse-skeleton` |
| JRX-002 | Report metadata, page properties, margins, columns | JRX | P0 | M1 | **Done** | JRX-001 | [JRX](features/JRX.md) | [JRX-002-001](tasks/M1/003-JRX-002-001-report-metadata-and-page-properties.md) | `testdata/jrxml/JRX-002-001-page-properties.jrxml`; `report_test.go` |
| JRX-003 | Report parameter declarations | JRX | P0 | M1 | **Done** | JRX-001 | [JRX](features/JRX.md) | [JRX-003-001](tasks/M1/004-JRX-003-001-parse-parameters.md) | `testdata/jrxml/JRX-003-001-parameters.jrxml`; `parameters_test.go` |
| JRX-004 | Report field declarations | JRX | P0 | M1 | **Done** | JRX-001 | [JRX](features/JRX.md) | [JRX-004-001](tasks/M1/005-JRX-004-001-parse-fields.md) | `testdata/jrxml/JRX-004-001-fields.jrxml`; `fields_test.go` |
| JRX-005 | Report variable declarations | JRX | P1 | M3 | Not Started | JRX-001 | [JRX](features/JRX.md) | INDEX | planned: variable fixture |
| JRX-006 | Query element parsing | JRX | P1 | M2 | Not Started | JRX-001 | [JRX](features/JRX.md) | INDEX | planned: query fixture |
| JRX-007 | Report properties and basic styles | JRX | P1 | M3 | Not Started | JRX-001 | [JRX](features/JRX.md) | INDEX | planned: style fixture |
| JRX-008 | Core bands (title/page/column/detail/summary) | JRX | P0 | M1 | **Done** | JRX-002, JRX-010 | [JRX](features/JRX.md) | [JRX-008-001](tasks/M1/007-JRX-008-001-parse-core-bands.md) | `testdata/jrxml/JRX-008-001-core-bands.jrxml`; `bands_test.go` |
| JRX-009 | Advanced bands (background/lastPageFooter/noData) | JRX | P2 | M3 | Not Started | JRX-008 | [JRX](features/JRX.md) | INDEX | planned: band fixture |
| JRX-010 | Basic elements: staticText, textField | JRX | P0 | M1 | **Done** | JRX-001 | [JRX](features/JRX.md) | [JRX-010-001](tasks/M1/006-JRX-010-001-parse-basic-elements.md) | `testdata/jrxml/JRX-010-001-basic-elements.jrxml`; `elements_test.go` |
| JRX-011 | Drawable elements: line, rectangle, ellipse | JRX | P1 | M2 | Not Started | JRX-010 | [JRX](features/JRX.md) | INDEX | planned: shape fixture |
| JRX-012 | Image element parsing | JRX | P2 | M4 | Not Started | JRX-010, PRT-003 | [JRX](features/JRX.md) | INDEX | planned: image fixture |
| JRX-013 | Style inheritance | JRX | P1 | M3 | Not Started | JRX-007 | [JRX](features/JRX.md) | INDEX | planned: inheritance fixture |
| JRX-014 | Element position/stretch/printWhen/evaluationTime | JRX | P0 | M1 | **Done** | JRX-010 | [JRX](features/JRX.md) | [JRX-014-001](tasks/M1/008-JRX-014-001-parse-element-settings.md) | `testdata/jrxml/JRX-014-001-element-settings.jrxml`; `element_attrs_test.go` |
| JRX-015 | Group declarations and reset behavior | JRX | P0 | M3 | Not Started | JRX-005 | [JRX](features/JRX.md) | INDEX | planned: group fixture |
| JRX-016 | Dataset runs and subdatasets | JRX | P2 | M4 | Not Started | JRX-006, DAT-009 | [JRX](features/JRX.md) | INDEX | planned: subdataset fixture |
| JRX-017 | Subreport element | JRX | P2 | M4 | Not Started | JRX-016 | [JRX](features/JRX.md) | INDEX | planned: subreport fixture |
| JRX-018 | Conditional styles and hyperlinks | JRX | P3 | M4 | Not Started | JRX-013 | [JRX](features/JRX.md) | INDEX | planned: conditional style fixture |
| JRX-019 | Unsupported-construct diagnostics and extension points | JRX | P1 | M1 | **Done** | JRX-001, FND-003 | [JRX](features/JRX.md) | [JRX-019-001](tasks/M1/009-JRX-019-001-unsupported-construct-diagnostics.md) | `testdata/jrxml/JRX-019-001-unsupported.jrxml`; `unsupported_test.go` |
| JRX-020 | Version-aware normalization | JRX | P1 | M2 | Not Started | JRX-001, JRX-019 | [JRX](features/JRX.md) | [JRX-020-001](tasks/M2/033-JRX-020-001-version-aware-normalization.md) | planned: version fixture |
| JRX-021 | Compiled `.jasper` binary loading | JRX | P3 | — | **Deferred** | — | [JRX](features/JRX.md) | — | Out of scope; see compat policy |
| JRX-022 | Compatibility matrix and fixture registry | JRX | P1 | M0 | **Done** | — | [JRX](features/JRX.md) | [JRX-022-001](tasks/M0/JRX-022-001-compatibility-matrix-and-fixture-registry.md) | [`COMPATIBILITY.md`](COMPATIBILITY.md) fixture registry |
| JRX-023 | Font declarations in JRXML | JRX | P2 | M4 | Not Started | LAY-005 | [JRX](features/JRX.md) | INDEX | planned: font fixture |
| JRX-024 | Component elements (charts/crosstabs/barcodes) | JRX | P3 | M6 | Not Started | JRX-001 | [JRX](features/JRX.md) | INDEX | planned: component fixture |

## EXP — Expression language and type semantics

Spec: [`features/EXP.md`](features/EXP.md) · Architect: [`architecture/expressions-and-types.md`](architecture/expressions-and-types.md)

| ID | Feature | Domain | Priority | Milestone | Status | Dependencies | Specification | Tasks | Verification evidence / note |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| EXP-001 | `$F`/`$P`/`$V` references and lexical scope | EXP | P0 | M1 | Done | — | [EXP](features/EXP.md) | [EXP-001-001](tasks/M0/EXP-001-001-expression-symbol-scope-contract.md), [EXP-001-002](tasks/M1/010-EXP-001-002-symbol-resolution.md) | `TestReferenceResolution`, `TestReferenceScopeAndNull` |
| EXP-002 | Literals and constants | EXP | P0 | M1 | Done | — | [EXP](features/EXP.md) | [EXP-003-001](tasks/M1/011-EXP-003-001-expression-parser-and-evaluator.md) | `TestEvaluation` |
| EXP-003 | Operator grammar and precedence | EXP | P0 | M1 | Done | — | [EXP](features/EXP.md) | [EXP-003-001](tasks/M1/011-EXP-003-001-expression-parser-and-evaluator.md) | `TestEvaluation`, `TestParseErrorsHavePositions` |
| EXP-004 | Boolean logic and comparisons | EXP | P0 | M1 | Done | EXP-003 | [EXP](features/EXP.md) | [EXP-003-001](tasks/M1/011-EXP-003-001-expression-parser-and-evaluator.md) | `TestEvaluation` |
| EXP-005 | Conditional (ternary) expressions | EXP | P1 | M1 | Not Started | EXP-004 | [EXP](features/EXP.md) | INDEX | planned: ternary test |
| EXP-006 | Null semantics and missing values | EXP | P0 | M1 | Done | — | [EXP](features/EXP.md) | [EXP-006-001](tasks/M1/012-EXP-006-001-null-semantics.md) | `TestNullDecisionTable` |
| EXP-007 | Type system and coercion rules | EXP | P0 | M2 | Not Started | EXP-002 | [EXP](features/EXP.md) | INDEX | planned: coercion test |
| EXP-008 | Decimal/BigDecimal-compatible arithmetic | EXP | P0 | M2 | Not Started | EXP-007 | [EXP](features/EXP.md) | INDEX | planned: decimal test |
| EXP-009 | String functions | EXP | P1 | M2 | Not Started | EXP-015 | [EXP](features/EXP.md) | INDEX | planned: function test |
| EXP-010 | Numeric conversions and rounding | EXP | P1 | M2 | Not Started | EXP-008 | [EXP](features/EXP.md) | INDEX | planned: rounding test |
| EXP-011 | Date/time handling, locale, timezone | EXP | P1 | M2 | Not Started | EXP-007 | [EXP](features/EXP.md) | INDEX | planned: date test |
| EXP-012 | Date formatting | EXP | P2 | M3 | Not Started | EXP-011 | [EXP](features/EXP.md) | INDEX | planned: format test |
| EXP-013 | Number/currency formatting patterns | EXP | P2 | M3 | Not Started | EXP-008 | [EXP](features/EXP.md) | INDEX | planned: format test |
| EXP-014 | Java-style method translation subset | EXP | P1 | M2 | Not Started | EXP-015 | [EXP](features/EXP.md) | INDEX | planned: translation test |
| EXP-015 | Safe function registry | EXP | P1 | M2 | Not Started | EXP-019 | [EXP](features/EXP.md) | INDEX | planned: registry test |
| EXP-016 | Compile-time type validation | EXP | P1 | M2 | Not Started | EXP-007 | [EXP](features/EXP.md) | INDEX | planned: validation test |
| EXP-017 | Expression caching | EXP | P2 | M3 | Not Started | EXP-003 | [EXP](features/EXP.md) | INDEX | planned: cache test |
| EXP-018 | Alternative expression languages (Groovy/JS) | EXP | P3 | — | **Deferred** | — | [EXP](features/EXP.md) | — | Not supported; see compat policy |
| EXP-019 | Sandbox: no arbitrary code execution | EXP | P0 | M1 | **Done** | — | [EXP](features/EXP.md) | [EXP-019-001](tasks/M0/EXP-019-001-expression-sandbox-policy.md) | sandbox policy in `architecture/expressions-and-types.md` |

## DAT — Data sources and queries

Spec: [`features/DAT.md`](features/DAT.md) · Architect: [`architecture/data-and-resources.md`](architecture/data-and-resources.md)

| ID | Feature | Domain | Priority | Milestone | Status | Dependencies | Specification | Tasks | Verification evidence / note |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| DAT-001 | Slice/struct/map data source | DAT | P0 | M1 | Done | DAT-002 | [DAT](features/DAT.md) | [DAT-001-001](tasks/M1/013-DAT-001-001-slice-data-source.md) | `TestSliceSourceRows`, `TestSliceCancellationAndClose` |
| DAT-002 | DataSource cursor lifecycle interface | DAT | P0 | M0 | **Done** | — | [DAT](features/DAT.md) | [DAT-002-001](tasks/M0/DAT-002-001-datasource-cursor-interface-contract.md) | `internal/datasource/datasource_test.go`; branch `task/DAT-002-001` |
| DAT-003 | JSON data source | DAT | P1 | M2 | Not Started | DAT-002 | [DAT](features/DAT.md) | [DAT-003-001](tasks/M2/030-DAT-003-001-json-data-source.md) | planned: JSON fixture |
| DAT-004 | CSV data source | DAT | P2 | M2 | Not Started | DAT-002 | [DAT](features/DAT.md) | [DAT-004-001](tasks/M2/031-DAT-004-001-csv-data-source.md) | planned: CSV fixture |
| DAT-005 | SQL via `database/sql` | DAT | P1 | M2 | Not Started | DAT-002, DAT-006 | [DAT](features/DAT.md) | INDEX | planned: SQL fixture (test DB) |
| DAT-006 | Query parameter binding | DAT | P1 | M2 | Not Started | DAT-002 | [DAT](features/DAT.md) | INDEX | planned: binding test |
| DAT-007 | Field conversion and type mapping | DAT | P1 | M2 | Not Started | DAT-002, EXP-007 | [DAT](features/DAT.md) | INDEX | planned: conversion test |
| DAT-008 | Context cancellation in data access | DAT | P0 | M2 | Not Started | DAT-002, FND-009 | [DAT](features/DAT.md) | INDEX | planned: cancellation test |
| DAT-009 | Multiple datasets and parameter mapping | DAT | P2 | M3 | Not Started | DAT-001, DAT-007 | [DAT](features/DAT.md) | INDEX | planned: dataset test |
| DAT-010 | Large datasets and streaming | DAT | P2 | M5 | Not Started | DAT-002, PRT-004 | [DAT](features/DAT.md) | INDEX | planned: large-fixture test |
| DAT-011 | HTTP/external data sources | DAT | P3 | M6 | Not Started | OPS-002 | [DAT](features/DAT.md) | INDEX | planned: allowlist test |

## FIL — Fill engine and report semantics

Spec: [`features/FIL.md`](features/FIL.md) · Architect: [`architecture/filling-and-groups.md`](architecture/filling-and-groups.md)

| ID | Feature | Domain | Priority | Milestone | Status | Dependencies | Specification | Tasks | Verification evidence / note |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| FIL-001 | Report execution state and record iteration | FIL | P0 | M1 | Done | DAT-002, FND-001 | [FIL](features/FIL.md) | [FIL-001-001](tasks/M1/014-FIL-001-001-fill-record-iteration.md) | `TestRunSnapshotsAndMissing`, `TestRunConcurrentIsolation` |
| FIL-002 | Built-in variables and parameter defaults | FIL | P1 | M1 | Done | FIL-001 | [FIL](features/FIL.md) | [FIL-002-001](tasks/M1/015-FIL-002-001-builtin-variables-and-parameter-defaults.md) | `TestReportCountPerRecord`, `TestParameterDefaultsAndOverrides` |
| FIL-003 | Variable calculations (count/sum/avg/min/max) | FIL | P1 | M3 | Not Started | FIL-005 | [FIL](features/FIL.md) | INDEX | planned: calculation test |
| FIL-004 | Distinct-count calculation | FIL | P2 | M3 | Not Started | FIL-003 | [FIL](features/FIL.md) | INDEX | planned: distinct test |
| FIL-005 | Evaluation/increment/reset scopes and timing | FIL | P1 | M3 | Not Started | FIL-001 | [FIL](features/FIL.md) | INDEX | planned: scope test |
| FIL-006 | Group start/end detection and nesting | FIL | P1 | M3 | Not Started | FIL-005, JRX-015 | [FIL](features/FIL.md) | INDEX | planned: group fixture |
| FIL-007 | Group page breaks and reprinted headers | FIL | P2 | M3 | Not Started | FIL-006, LAY-011 | [FIL](features/FIL.md) | INDEX | planned: group-break fixture |
| FIL-008 | Conditional visibility (printWhen) evaluation | FIL | P1 | M2 | Not Started | EXP-004, FIL-001 | [FIL](features/FIL.md) | INDEX | planned: printWhen test |
| FIL-009 | Title/summary/noData lifecycle | FIL | P1 | M3 | Not Started | FIL-001, JRX-009 | [FIL](features/FIL.md) | INDEX | planned: lifecycle test |
| FIL-010 | Delayed evaluation (page/total pages) | FIL | P2 | M3 | Not Started | FIL-002, LAY-016 | [FIL](features/FIL.md) | INDEX | planned: total-pages test |
| FIL-011 | Multi-column filling and column counters | FIL | P2 | M3 | Not Started | JRX-002, LAY-011 | [FIL](features/FIL.md) | INDEX | planned: multi-column fixture |
| FIL-012 | Subreport execution and parameter passing | FIL | P2 | M4 | Not Started | JRX-017, FIL-001 | [FIL](features/FIL.md) | INDEX | planned: subreport fixture |
| FIL-013 | Report recursion limits | FIL | P2 | M4 | Not Started | FIL-012 | [FIL](features/FIL.md) | INDEX | planned: recursion test |
| FIL-014 | Sort/filter at fill time | FIL | P3 | M3 | Not Started | DAT-009 | [FIL](features/FIL.md) | INDEX | planned: sort/ filter test |

## LAY — Layout, typography, pagination

Spec: [`features/LAY.md`](features/LAY.md) · Architect: [`architecture/layout-and-pagination.md`](architecture/layout-and-pagination.md)

| ID | Feature | Domain | Priority | Milestone | Status | Dependencies | Specification | Tasks | Verification evidence / note |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| LAY-001 | Page geometry and band placement | LAY | P0 | M1 | Done | JRX-002, PRT-001 | [LAY](features/LAY.md) | [LAY-001-001](tasks/M1/019-LAY-001-001-page-geometry-and-band-placement.md) | `TestPrintableAreaAndBandPlacement`, `TestGeometryInvalidAndOversizeBand` |
| LAY-002 | Absolute/float positioning | LAY | P0 | M1 | Not Started | LAY-001 | [LAY](features/LAY.md) | [LAY-002-001](tasks/M1/020-LAY-002-001-absolute-positioning.md) | planned: position test |
| LAY-003 | Fixed-height detail band layout | LAY | P0 | M1 | Not Started | LAY-002, FIL-001 | [LAY](features/LAY.md) | [LAY-003-001](tasks/M1/021-LAY-003-001-fixed-height-detail-band-layout.md) | planned: layout golden |
| LAY-004 | Relative stretching | LAY | P1 | M2 | Not Started | LAY-003 | [LAY](features/LAY.md) | INDEX | planned: stretch fixture |
| LAY-005 | Font resolver (bundled core fonts) | LAY | P0 | M1 | Done | — | [LAY](features/LAY.md) | [LAY-005-001](tasks/M1/016-LAY-005-001-font-resolver.md) | `TestResolveCoreFontsAndFallback`, `TestMetricTableStable` |
| LAY-006 | Text measurement (baseline metrics) | LAY | P0 | M1 | Done | LAY-005 | [LAY](features/LAY.md) | [LAY-006-001](tasks/M1/017-LAY-006-001-text-measurement-baseline.md) | `TestMeasureFixedStrings`, `TestMeasureDeterminismAndScale` |
| LAY-007 | Unicode shaping and fallback | LAY | P2 | M5 | Not Started | LAY-006 | [LAY](features/LAY.md) | INDEX | planned: Unicode fixture |
| LAY-008 | Wrapping, clipping, alignment, rotation | LAY | P1 | M2 | Not Started | LAY-006 | [LAY](features/LAY.md) | [LAY-008-001](tasks/M2/032-LAY-008-001-text-wrapping-and-alignment.md) | planned: wrap fixture |
| LAY-009 | Styled text and markup subset | LAY | P2 | M3 | Not Started | LAY-008 | [LAY](features/LAY.md) | INDEX | planned: markup fixture |
| LAY-010 | Variable band height and split types | LAY | P2 | M3 | Not Started | LAY-004, LAY-006 | [LAY](features/LAY.md) | INDEX | planned: split fixture |
| LAY-011 | Column/page breaks and keep-together | LAY | P1 | M2 | Not Started | LAY-003 | [LAY](features/LAY.md) | INDEX | planned: break fixture |
| LAY-012 | Borders, backgrounds, colors, transparency | LAY | P1 | M2 | Not Started | PRT-002 | [LAY](features/LAY.md) | INDEX | planned: style golden |
| LAY-013 | Image scaling | LAY | P2 | M4 | Not Started | JRX-012, PRT-003 | [LAY](features/LAY.md) | INDEX | planned: image golden |
| LAY-014 | Shape geometry (line/rect/ellipse) | LAY | P1 | M2 | Not Started | JRX-011, PRT-001 | [LAY](features/LAY.md) | INDEX | planned: shape golden |
| LAY-015 | Deterministic measurement and convergence | LAY | P1 | M2 | Not Started | LAY-006 | [LAY](features/LAY.md) | INDEX | planned: determinism test |
| LAY-016 | Orphans/widows and page-count reservation | LAY | P2 | M3 | Not Started | LAY-010 | [LAY](features/LAY.md) | INDEX | planned: orphan fixture |
| LAY-017 | Bidirectional text | LAY | P3 | M5 | Not Started | LAY-007 | [LAY](features/LAY.md) | INDEX | planned: bidi fixture |
| LAY-018 | Hyperlink targets in layout | LAY | P3 | M3 | Not Started | LAY-008, PRT-002 | [LAY](features/LAY.md) | INDEX | planned: link golden |

## PRT — Print document and storage

Spec: [`features/PRT.md`](features/PRT.md) · Architect: [`architecture/intermediate-representations.md`](architecture/intermediate-representations.md)

| ID | Feature | Domain | Priority | Milestone | Status | Dependencies | Specification | Tasks | Verification evidence / note |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| PRT-001 | Typed print document model | PRT | P0 | M1 | Done | — | [PRT](features/PRT.md) | [PRT-001-001](tasks/M0/PRT-001-001-print-document-ir-contract.md), [PRT-001-002](tasks/M1/018-PRT-001-002-print-document-model.md) | `TestDocumentAppendAndIterate`, `TestDocumentRejectsInvalidElementsAndStopsOnError` |
| PRT-002 | Resolved style data in print IR | PRT | P1 | M1 | Not Started | PRT-001 | [PRT](features/PRT.md) | [PRT-002-001](tasks/M1/022-PRT-002-001-resolved-style-data.md) | planned: style test |
| PRT-003 | Resource references (fonts/images) | PRT | P1 | M2 | Not Started | PRT-001 | [PRT](features/PRT.md) | INDEX | planned: reference test |
| PRT-004 | Lazy/streaming print storage | PRT | P2 | M3 | Not Started | PRT-001 | [PRT](features/PRT.md) | INDEX | planned: streaming test |
| PRT-005 | Spill-to-disk storage for large documents | PRT | P3 | M5 | Not Started | PRT-004 | [PRT](features/PRT.md) | INDEX | planned: large-doc test |
| PRT-006 | Print IR serialization policy | PRT | P3 | — | **Deferred** | PRT-004 | [PRT](features/PRT.md) | — | No public IR format planned |
| PRT-007 | Page metadata (numbers, size, count) | PRT | P1 | M1 | Done | PRT-001 | [PRT](features/PRT.md) | [PRT-001-002](tasks/M1/018-PRT-001-002-print-document-model.md) | `TestDocumentAppendAndIterate`, `TestEmptyDocumentAndValidation` |

## OUT — Exporters

Spec: [`features/OUT.md`](features/OUT.md) · Architect: [`architecture/exporters.md`](architecture/exporters.md)

| ID | Feature | Domain | Priority | Milestone | Status | Dependencies | Specification | Tasks | Verification evidence / note |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| OUT-001 | Baseline PDF exporter | OUT | P0 | M1 | Not Started | PRT-001, PRT-007, OUT-009 | [OUT](features/OUT.md) | [OUT-001-001](tasks/M1/023-OUT-001-001-baseline-pdf-exporter.md) | planned: PDF golden |
| OUT-002 | PDF text and font rendering | OUT | P0 | M1 | Not Started | OUT-001 | [OUT](features/OUT.md) | [OUT-001-001](tasks/M1/023-OUT-001-001-baseline-pdf-exporter.md) | planned: text golden |
| OUT-003 | PDF vector geometry and images | OUT | P1 | M2 | Not Started | OUT-001, LAY-014 | [OUT](features/OUT.md) | INDEX | planned: shape golden |
| OUT-004 | PDF font embedding | OUT | P1 | M2 | Not Started | OUT-002, LAY-005 | [OUT](features/OUT.md) | INDEX | planned: embedding golden |
| OUT-005 | PDF metadata and links | OUT | P2 | M3 | Not Started | OUT-001, LAY-018 | [OUT](features/OUT.md) | INDEX | planned: metadata test |
| OUT-006 | XLSX exporter | OUT | P2 | M4 | Not Started | PRT-002 | [OUT](features/OUT.md) | INDEX | planned: XLSX golden |
| OUT-007 | HTML exporter | OUT | P2 | M4 | Not Started | PRT-002 | [OUT](features/OUT.md) | INDEX | planned: HTML golden |
| OUT-008 | CSV exporter | OUT | P2 | M4 | Not Started | PRT-002, FIL-001 | [OUT](features/OUT.md) | INDEX | planned: CSV golden |
| OUT-009 | Streaming `io.Writer` output | OUT | P0 | M1 | Not Started | PRT-001 | [OUT](features/OUT.md) | [OUT-001-001](tasks/M1/023-OUT-001-001-baseline-pdf-exporter.md) | planned: streaming test |
| OUT-010 | Deterministic output | OUT | P1 | M1 | Not Started | OUT-001 | [OUT](features/OUT.md) | [OUT-010-001](tasks/M1/024-OUT-010-001-deterministic-output-tests.md) | planned: byte-stability test |

## TST — Test infrastructure and correctness suites

Spec: [`features/TST.md`](features/TST.md) · Plan: [`TEST_STRATEGY.md`](TEST_STRATEGY.md)

| ID | Feature | Domain | Priority | Milestone | Status | Dependencies | Specification | Tasks | Verification evidence / note |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| TST-001 | Reproducible baseline and fixture infrastructure | TST | P0 | M0 | **Done** | — | [TST](features/TST.md) | [TST-001-001](tasks/M0/TST-001-001-deterministic-test-environment-helpers.md), [TST-001-002](tasks/M0/TST-001-002-fixture-conventions-and-golden-workflow.md) | baseline (`Makefile`, CI) + `internal/testutil` + fixture conventions (`testdata/README.md`) |
| TST-002 | Golden fixture harness | TST | P0 | M1 | Not Started | TST-001 | [TST](features/TST.md) | [TST-002-001](tasks/M1/027-TST-002-001-golden-fixture-harness.md) | planned: harness test |
| TST-003 | Semantic assertions | TST | P1 | M2 | Not Started | TST-002 | [TST](features/TST.md) | INDEX | planned: assertion helpers |
| TST-004 | Pagination/page-count assertions | TST | P1 | M2 | Not Started | TST-002, LAY-011 | [TST](features/TST.md) | INDEX | planned: page-count test |
| TST-005 | Visual regression testing | TST | P2 | M5 | Not Started | TST-002 | [TST](features/TST.md) | INDEX | planned: render diff |
| TST-006 | Differential JasperReports harness (test-only) | TST | P2 | M5 | Not Started | TST-003 | [TST](features/TST.md) | INDEX | planned: harness (isolated) |
| TST-007 | Fuzzing (XML, expressions) | TST | P2 | M2 | Not Started | JRX-001, EXP-003 | [TST](features/TST.md) | INDEX | planned: fuzz targets |
| TST-008 | Property-based tests | TST | P3 | M5 | Not Started | TST-007 | [TST](features/TST.md) | INDEX | planned: property suite |
| TST-009 | Race, cancellation and leak tests | TST | P1 | M2 | Not Started | FND-009 | [TST](features/TST.md) | [TST-009-001](tasks/M1/029-TST-009-001-race-cancellation-leak-tests.md) | planned: race/leak tests |
| TST-010 | Benchmarks and profiling | TST | P2 | M5 | Not Started | FND-001 | [TST](features/TST.md) | INDEX | planned: benchmarks |
| TST-011 | Cross-platform smoke tests | TST | P2 | M5 | Not Started | OUT-001 | [TST](features/TST.md) | INDEX | planned: OS matrix |
| TST-012 | Compatibility fixture corpus | TST | P1 | M2 | Not Started | TST-002 | [TST](features/TST.md) | INDEX | planned: corpus |

## OPS — Security, performance, diagnostics, integration

Spec: [`features/OPS.md`](features/OPS.md) · Architect: [`architecture/concurrency-and-limits.md`](architecture/concurrency-and-limits.md)

| ID | Feature | Domain | Priority | Milestone | Status | Dependencies | Specification | Tasks | Verification evidence / note |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| OPS-001 | Resource limit enforcement | OPS | P0 | M1 | Not Started | FND-009 | [OPS](features/OPS.md) | [OPS-001-001](tasks/M1/028-OPS-001-001-resource-limits-enforcement.md) | planned: limit tests |
| OPS-002 | Resource URI allowlists (SSRF/path traversal) | OPS | P1 | M4 | Not Started | FND-002 | [OPS](features/OPS.md) | INDEX | planned: allowlist tests |
| OPS-003 | SQL permissions and safe binding | OPS | P1 | M2 | Not Started | DAT-005, DAT-006 | [OPS](features/OPS.md) | INDEX | planned: injection tests |
| OPS-004 | Recursion/decompression/archive limits | OPS | P2 | M4 | Not Started | FIL-013 | [OPS](features/OPS.md) | INDEX | planned: limit tests |
| OPS-005 | Debugging and tracing controls | OPS | P2 | M5 | Not Started | FND-003 | [OPS](features/OPS.md) | INDEX | planned: debug flag test |
| OPS-006 | OpenTelemetry hooks | OPS | P3 | M5 | Not Started | FND-007 | [OPS](features/OPS.md) | INDEX | planned: span test |
| OPS-007 | Memory/CPU budgets and predictable use | OPS | P1 | M3 | Not Started | OPS-001, PRT-004 | [OPS](features/OPS.md) | INDEX | planned: profile test |
| OPS-008 | Third-party license/SBOM review | OPS | P1 | M0 | **Done** | — | [OPS](features/OPS.md) | [OPS-008-001](tasks/M0/OPS-008-001-dependency-and-license-review-policy.md) | [`DEPENDENCY_POLICY.md`](DEPENDENCY_POLICY.md) |
| OPS-009 | JasperReports interoperability/licensing review | OPS | P1 | M5 | Not Started | TST-006 | [OPS](features/OPS.md) | INDEX | planned: review record |
| OPS-010 | Changelog and release guidance | OPS | P2 | M5 | Not Started | FND-006 | [OPS](features/OPS.md) | INDEX | planned: changelog check |
| OPS-011 | Public examples and compatibility corpus | OPS | P2 | M5 | Not Started | OUT-001, FND-005 | [OPS](features/OPS.md) | INDEX | planned: example build |

## EXT — Advanced extensions (deferred)

Spec: [`features/README.md`](features/README.md) (EXT group) · All P3, milestone M6 unless noted.

| ID | Feature | Domain | Priority | Milestone | Status | Dependencies | Specification | Tasks | Verification evidence / note |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| EXT-001 | Charts | EXT | P3 | M6 | Not Started | JRX-024, LAY-014 | [EXT](features/README.md) | INDEX | deferred |
| EXT-002 | Crosstabs | EXT | P3 | M6 | Not Started | JRX-024, FIL-006 | [EXT](features/README.md) | INDEX | deferred |
| EXT-003 | Barcodes | EXT | P3 | M6 | Not Started | JRX-024, LAY-014 | [EXT](features/README.md) | INDEX | deferred |
| EXT-004 | PDF/A | EXT | P3 | M6 | Not Started | OUT-004 | [EXT](features/README.md) | INDEX | deferred |
| EXT-005 | Accessibility / tagged PDF | EXT | P3 | M6 | Not Started | OUT-001 | [EXT](features/README.md) | INDEX | deferred |
| EXT-006 | Optional HTTP service adapter | EXT | P3 | M6 | Not Started | FND-008 | [EXT](features/README.md) | INDEX | deferred; separate module |
| EXT-007 | Optional gRPC service adapter | EXT | P3 | M6 | Not Started | EXT-006 | [EXT](features/README.md) | INDEX | deferred; separate module |
| EXT-008 | JRXML 7.x adaptation | EXT | P3 | M6 | Not Started | JRX-020 | [EXT](features/README.md) | INDEX | deferred |
| EXT-009 | Report designer integration | EXT | P3 | M6 | Not Started | EXT-001, EXT-002 | [EXT](features/README.md) | INDEX | deferred |

## Aggregate

| Metric | Count |
| --- | --- |
| Total features | 146 |
| Done | 17 (`FND-002`, `FND-003`, `FND-009`, `FND-011`, `JRX-001`, `JRX-002`, `JRX-003`, `JRX-004`, `JRX-008`, `JRX-010`, `JRX-014`, `JRX-019`, `JRX-022`, `DAT-002`, `EXP-019`, `OPS-008`, `TST-001`) |
| In Progress | 3 (`EXP-001`, `PRT-001`, `FND-001`) |
| Deferred | 3 (`JRX-021`, `EXP-018`, `PRT-006`) |
| Not Started | 123 |
| M0 features | 7 |
| M1 features | 35 |
| M2 features | 33 |
| M3 features | 26 |
| M4 features | 14 |
| M5 features | 16 |
| M6 features | 12 |

> Counts are a derived snapshot; the tables above are canonical. Milestone
> counts include features with that `Milestone` value; `Deferred` features have
> no milestone.

## Open questions

| # | Question | Affects | Owner/process |
| --- | --- | --- | --- |
| OQ-1 | Decimal backend (`math/big.Rat`, third-party decimal, scaled ints)? | `EXP-008` | ADR when `EXP-008` starts |
| OQ-2 | Internal unit system for print IR (twips/points/sub-point)? | `PRT-001`, `LAY-*` | **Resolved: twips** (`ADR-0004`, `PRT-001-001`) |
| OQ-3 | Text shaping backend (pure-Go vs dependency)? | `LAY-006`, `LAY-007` | ADR when `LAY-007` starts |
| OQ-4 | Are remote fonts/images ever permitted, and under what policy? | `OPS-002`, `DAT-011` | `OPS-002` |
| OQ-5 | PDF library choice (pure-Go vs dependency)? | `OUT-001` | `OUT-001-001` / ADR |
| OQ-6 | Does `internal/exec` exist, or does `report` orchestrate? | `FND-001` | **Resolved: `report` orchestrates; `internal/exec` deferred** (`ADR-0005`) |
