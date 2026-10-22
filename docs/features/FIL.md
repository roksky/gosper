# FIL — Fill engine and report semantics

**Architecture:** [`../architecture/filling-and-groups.md`](../architecture/filling-and-groups.md).

## Purpose

Execute a compiled report against a data source: iterate records, evaluate
expressions, manage variables/groups, and produce a per-run filled state.

## Technical expectations

- **One fill state per run**; never shared or stored in compiled IR (`FIL-001`).
- Deterministic evaluation order (report → group close/open → detail → band
  scopes) is a contract (`FIL-005`).
- Built-in variables (`PAGE_NUMBER`, `COLUMN_NUMBER`, `REPORT_COUNT`) and
  parameter defaults (`FIL-002`).
- Calculations: count/sum/avg/min/max/distinct (`FIL-003`, `FIL-004`).
- Deferred evaluation resolved after pagination (`FIL-010`).
- Subreports run independent fills; recursion bounded (`FIL-012`, `FIL-013`).

## Features

| ID | Expectation | Acceptance test (planned) |
| --- | --- | --- |
| FIL-001 | Record iteration + state. | Iteration count test. |
| FIL-002 | Built-ins + param defaults. | Built-in value test. |
| FIL-003 | Variable calculations. | Sum/avg/min/max fixture. |
| FIL-004 | Distinct count. | Distinct fixture. |
| FIL-005 | Eval/increment/reset scopes. | Scope-order test. |
| FIL-006 | Group detection/nesting. | Nested-group fixture. |
| FIL-007 | Group page breaks/reprinted headers. | Group-break fixture. |
| FIL-008 | printWhen evaluation. | Visibility fixture. |
| FIL-009 | Title/summary/noData lifecycle. | Lifecycle fixture. |
| FIL-010 | Delayed page/total values. | Total-pages fixture. |
| FIL-011 | Multi-column filling. | Multi-column fixture. |
| FIL-012 | Subreport execution. | Subreport fixture. |
| FIL-013 | Recursion limits. | Deep-recursion diagnostic. |
| FIL-014 | Sort/filter at fill. | Sorted output test. |

## Tasks

`FIL-001-001`, `FIL-002-001` (M1).

## Risks

- Group/reset ordering is subtle and layout-coupled; test at each band.
