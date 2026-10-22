# LAY-002-001: Absolute positioning

- **Status:** Done
- **Feature:** [LAY-002](../../FEATURE_MATRIX.md#lay--layout-typography-pagination)
- **Milestone:** M1
- **Priority:** P0
- **Dependencies:** LAY-001-001
- **Blocked by:** none
- **Primary owner/package:** `internal/layout`
- **Likely touched files:** `internal/layout/position.go`, `internal/layout/position_test.go`
- **Collision risk:** Medium — uses the canvas; separate file.
- **Parallel-safe with:** `internal/text` tasks
- **Estimated size:** S

## Objective

Translate element `x`/`y` (relative to its band) into absolute page coordinates
and append elements to the print IR.

## Current state / prerequisites

Canvas/geometry exists.

## Scope (in)

- Add band origin to element offsets; produce absolute `Bounds`.
- Emit text elements with resolved position (text content already bound).
- Tests asserting absolute coordinates.

## Out of scope

- Stretching (`LAY-004`); styles (`LAY-012`).

## Proposed contracts and integration points

Consumes `Element{X,Y,W,H}` from the model; writes `printdoc` elements.

## Implementation checklist (future agent only)

- [x] Implement offset translation.
- [x] Tests with known band origin + element offsets.
- [x] Run `make check`.

## Acceptance criteria (observable)

- Absolute coordinates equal band origin + element offset.

## Tests and fixtures to add (future agent only)

- `internal/layout/position_test.go`.

## Failure cases / edge cases

- Elements partially outside the printable area (clip or diagnostic policy).

## Definition of done

Implemented/tested; linked from architecture.

## References and downstream dependencies

- Feature: `LAY-002`.
- Unblocks: `LAY-003-001`.
