# LAY-001-001: Page geometry and band placement

- **Status:** Done
- **Feature:** [LAY-001](../../FEATURE_MATRIX.md#lay--layout-typography-pagination)
- **Milestone:** M1
- **Priority:** P0
- **Dependencies:** PRT-001-002
- **Blocked by:** none
- **Primary owner/package:** `internal/layout`
- **Likely touched files:** `internal/layout/geometry.go`, `internal/layout/geometry_test.go`
- **Collision risk:** High — foundational layout file other LAY tasks extend.
- **Parallel-safe with:** `internal/text` tasks, `internal/jrxml` tasks
- **Estimated size:** M

## Objective

Compute the printable area from page size/margins and track the current band
placement cursor; create pages in the print IR as bands are emitted.

## Current state / prerequisites

Print document model (`PRT-001-002`) and page properties parsing
(`JRX-002-001`).

## Scope (in)

- Printable-area computation (width/height minus margins).
- Band cursor: place a band of known height at the current position.
- Start a new page when the cursor would overflow (minimal pagination).
- Tests asserting coordinates and page starts.

## Out of scope

- Stretching/splits (`LAY-004`, `LAY-010`); headers/footers (partially
  `JRX-008-001` provides bands, `LAY-003-001` places them).

## Proposed contracts and integration points

```go
// Proposed — non-binding.
type Canvas struct { Doc *printdoc.Document; /* cursor */ }
func (c *Canvas) Place(b *Band, h printdoc.Unit) (printdoc.Unit, error)
```

## Implementation checklist (future agent only)

- [x] Implement printable area + cursor.
- [x] Emit pages into `printdoc.Document`.
- [x] Tests for placement and page break.
- [x] Run `make check`.

## Acceptance criteria (observable)

- A band placed at cursor `y` appears at the expected coordinate.
- Exceeding page height starts a new page.

## Tests and fixtures to add (future agent only)

- `internal/layout/geometry_test.go`.

## Failure cases / edge cases

- Band taller than a full page (diagnostic; no infinite loop).

## Definition of done

Implemented/tested; linked from architecture.

## References and downstream dependencies

- Feature: `LAY-001`.
- Unblocks: `LAY-002-001`, `LAY-003-001`.
