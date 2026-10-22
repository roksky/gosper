# PRT-001-001: Define the print document IR contract

- **Status:** Not Started
- **Feature:** [PRT-001](../../FEATURE_MATRIX.md#prt--print-document-and-storage)
- **Milestone:** M0
- **Priority:** P0
- **Dependencies:** none
- **Blocked by:** none
- **Primary owner/package:** `internal/printdoc`
- **Likely touched files:** `internal/printdoc/doc.go`, `internal/printdoc/model.go`, `internal/printdoc/model_test.go`
- **Collision risk:** Medium — `PRT-001-002` (M1) implements next to this; this task fixes the type shape first.
- **Parallel-safe with:** `DAT-002-001`, `EXP-001-001`, `JRX-001-001`
- **Estimated size:** S

## Objective

Define the print IR types (document, page, element kinds, bounds, page
metadata) and **settle the internal unit system** (OQ-2) so layout and
exporters agree.

## Current state / prerequisites

`internal/printdoc/doc.go` is a stub. Unit system is undecided.

## Scope (in)

- `Unit` representation (e.g. scaled integer twips) with conversion helpers.
- `Document`, `Page`, `Element` (text/image/shape), `Bounds`, `Metadata`.
- Decision + rationale for the unit system, recorded in the PR and, if
  durable, an ADR.
- Tests for unit conversions (no floating drift).

## Out of scope

- Layout producing documents (`LAY-*`).
- Renderers (`OUT-*`).
- Serialization (`PRT-006`, Deferred).

## Proposed contracts and integration points

```go
// Proposed — non-binding.
type Unit int64 // e.g. twips (1/1440 inch)
type Document struct { Pages []Page; Metadata Metadata }
type Page struct { Width, Height Unit; Elements []Element }
type Metadata struct { PageNumber, PageCount int }
```

Integration: `internal/layout` builds it; `internal/exporter` consumes it.

## Implementation checklist (future agent only)

- [ ] Define unit type and conversions (inch/mm/point ↔ unit) with rounding
      rules.
- [ ] Define document/page/element/metadata types.
- [ ] Add tests for conversion round-trips and boundary values.
- [ ] Record the unit decision (ADR or task note).
- [ ] Run `make check`.

## Acceptance criteria (observable)

- A test constructs a page with known dimensions and reads back exact units.
- Conversion tests show no cumulative error.

## Tests and fixtures to add (future agent only)

- `internal/printdoc/model_test.go`.

## Failure cases / edge cases

- Unit overflow on very large pages; negative/zero dimensions rejected.

## Definition of done

Types compile, tests pass, unit decision documented; referenced from
[`../../architecture/intermediate-representations.md`](../../architecture/intermediate-representations.md).

## References and downstream dependencies

- Feature: `PRT-001`; OQ-2 in [`../../FEATURE_MATRIX.md`](../../FEATURE_MATRIX.md#open-questions).
- Unblocks: `PRT-001-002`, `LAY-001-001`, `OUT-001-001`.
