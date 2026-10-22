# DAT-004-001: CSV data source

- **Status:** Not Started
- **Feature:** [DAT-004](../../FEATURE_MATRIX.md#dat--data-sources-and-queries)
- **Milestone:** M2
- **Priority:** P2
- **Dependencies:** DAT-002-001
- **Blocked by:** none
- **Primary owner/package:** `internal/datasource`
- **Likely touched files:** `internal/datasource/csv.go`, `internal/datasource/csv_test.go`, `testdata/csv/basic.csv`
- **Collision risk:** Low — new file.
- **Parallel-safe with:** `DAT-003-001`, `LAY-008-001`, `JRX-020-001`
- **Estimated size:** S

## Objective

Implement a data source over CSV using `encoding/csv`: first row as header,
fields exposed by column name.

## Current state / prerequisites

Cursor interface exists.

## Scope in/out

- In: header-based field access; quoted fields; configurable delimiter.
- Out: type inference beyond string/number heuristics (`DAT-007` finalizes).

## Proposed contracts and integration points

```go
// Proposed — non-binding.
func FromCSV(r io.Reader, opts ...CSVOption) DataSource
```

## Implementation checklist (future agent only)

- [ ] Implement decode + record access.
- [ ] Tests against `testdata/csv/basic.csv`.
- [ ] Run `make check`.

## Acceptance criteria (observable)

- Rows iterate in order; fields addressable by header name.
- Ragged rows produce a diagnostic or documented behavior.

## Tests and fixtures to add (future agent only)

- `internal/datasource/csv_test.go` + fixture.

## Failure cases / edge cases

- BOM; duplicate headers; empty file.

## Definition of done

Implemented/tested; linked in `features/DAT.md`.

## References and downstream dependencies

- Feature: `DAT-004`.
