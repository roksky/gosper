# DAT — Data sources and queries

**Architecture:** [`../architecture/data-and-resources.md`](../architecture/data-and-resources.md).

## Purpose

Define the data-source/cursor contract and provide concrete sources (slice,
JSON, CSV, SQL) with safe parameter binding and context cancellation.

## Technical expectations

- Cursor lifecycle: `Open` → `Next` → `Close`; `Close` idempotent.
- **No string concatenation** of parameters into SQL (`OPS-003`).
- Field-to-value conversion centralized and consistent with `EXP` types
  (`DAT-007`).
- Cancellation honored in `Open`/`Next` (`DAT-008`).
- Forward-only iteration in early milestones.

### Proposed contract — non-binding

```go
type DataSource interface{ Open(ctx context.Context) (Cursor, error) }
type Cursor interface {
    Next(ctx context.Context) (Record, bool, error)
    Close() error
}
```

## Features

| ID | Expectation | Acceptance test (planned) |
| --- | --- | --- |
| DAT-001 | Slice/struct/map source. | Iterate slice fixture. |
| DAT-002 | Cursor interface contract. | Interface compile + lifecycle test. |
| DAT-003 | JSON source. | Parse JSON fixture. |
| DAT-004 | CSV source. | Parse CSV fixture. |
| DAT-005 | SQL via `database/sql`. | Query ephemeral test DB. |
| DAT-006 | Parameter binding. | Bound-parameter test. |
| DAT-007 | Field conversion/type mapping. | Type-map test. |
| DAT-008 | Context cancellation. | Cancel mid-iteration. |
| DAT-009 | Multiple datasets/param mapping. | Multi-dataset fixture. |
| DAT-010 | Large datasets/streaming. | Large-fixture memory test. |
| DAT-011 | HTTP/external sources (Deferred). | Allowlist test. |

## Tasks

`DAT-002-001` (M0); `DAT-001-001` (M1); `DAT-003-001`, `DAT-004-001` (M2).

## Risks

- Driver/DB availability in CI; use ephemeral or in-memory fixtures only.

`DAT-003` is implemented by `internal/datasource/json.go`; `testdata/json/basic.json` and `internal/datasource/json_test.go` cover ordered rows, typed values, nested values, and malformed input.

`DAT-004` is implemented by `internal/datasource/csv.go`; `testdata/csv/basic.csv` and `internal/datasource/csv_test.go` cover headers, quoted fields, BOM, delimiter options, and malformed rows. Values are strings.
