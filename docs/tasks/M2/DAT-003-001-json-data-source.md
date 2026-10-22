# DAT-003-001: JSON data source

- **Status:** Not Started
- **Feature:** [DAT-003](../../FEATURE_MATRIX.md#dat--data-sources-and-queries)
- **Milestone:** M2
- **Priority:** P1
- **Dependencies:** DAT-002-001
- **Blocked by:** none
- **Primary owner/package:** `internal/datasource`
- **Likely touched files:** `internal/datasource/json.go`, `internal/datasource/json_test.go`, `testdata/json/basic.json`
- **Collision risk:** Low — new file next to the interface.
- **Parallel-safe with:** `DAT-004-001`, `LAY-008-001`, `JRX-020-001`
- **Estimated size:** S/M

## Objective

Implement a data source over JSON: an array of objects (and later a pointer to
an array via a query string) exposing fields by name.

## Current state / prerequisites

Cursor interface and slice source exist.

## Scope (in)

- Decode a JSON array of objects; expose fields by key.
- Type mapping for JSON numbers/strings/bools/null.
- Tests: iteration, missing keys, nested value handling policy.

## Out of scope

- JSONPath queries; streaming huge files (`DAT-010`).

## Proposed contracts and integration points

```go
// Proposed — non-binding.
func FromJSON(r io.Reader) DataSource
```

## Implementation checklist (future agent only)

- [ ] Implement decoding + record access.
- [ ] Tests against `testdata/json/basic.json`.
- [ ] Run `make check`.

## Acceptance criteria (observable)

- Iterating the fixture yields the expected values in order.
- Missing key returns `(nil,false)`.

## Failure cases / edge cases

- Top-level non-array JSON ⇒ diagnostic; invalid JSON ⇒ diagnostic.

## Definition of done

Implemented/tested; compatibility note; linked in `features/DAT.md`.

## References and downstream dependencies

- Feature: `DAT-003`.
- Unblocks: M2 DSL/fill coverage; `TST-012`.
