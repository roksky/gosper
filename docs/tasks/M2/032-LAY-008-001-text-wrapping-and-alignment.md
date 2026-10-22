# LAY-008-001: Text wrapping and alignment

- **Status:** Done
- **Feature:** [LAY-008](../../FEATURE_MATRIX.md#lay--layout-typography-pagination)
- **Milestone:** M2
- **Priority:** P1
- **Dependencies:** LAY-006-001
- **Blocked by:** none
- **Primary owner/package:** `internal/text` + `internal/layout`
- **Likely touched files:** `internal/text/wrap.go`, `internal/text/wrap_test.go`, `internal/layout/text.go`
- **Collision risk:** Medium — spans two packages; keep measurement in `internal/text` and placement in `internal/layout`.
- **Parallel-safe with:** `DAT-003-001`, `DAT-004-001`, `JRX-020-001`
- **Estimated size:** M

## Objective

Wrap text to a measured width and apply horizontal/vertical alignment, without
shaping/Unicode fallback.

## Current state / prerequisites

Deterministic measurement exists.

## Scope (in)

- Greedy word wrapping using measured widths; break on whitespace.
- Alignment: left/center/right; vertical top/middle/bottom.
- Rotation deferred to a later task (document).
- Tests: wrap points for known strings; alignment offsets.

## Out of scope

- Unicode line breaking (`LAY-007`); styled text (`LAY-009`).

## Proposed contracts and integration points

`Wrap(s, font, size, width) []string`; layout consumes lines.

## Implementation checklist (future agent only)

- [x] Implement wrapping (greedy) + alignment.
- [x] Tests for wrap points and offsets.
- [x] Run `make check`.

## Acceptance criteria (observable)

- No line exceeds the measured width (within one glyph policy).
- Alignment offsets match expected values.

## Tests and fixtures to add (future agent only)

- `internal/text/wrap_test.go`, `internal/layout/text_test.go`.

## Failure cases / edge cases

- Words longer than the width; empty text; single glyph.

## Definition of done

Implemented/tested; compatibility rows for wrapping/alignment updated.

## References and downstream dependencies

- Feature: `LAY-008`.
- Unblocks: `LAY-009`, `LAY-010`.

## Verification and compatibility

`GOCACHE=/tmp/gosper-go-cache make check` passes. `internal/text/wrap_test.go` covers measured greedy wrapping and long words; `internal/layout/text_test.go` covers horizontal and vertical offsets and height overflow. Rotation, clipping, and Unicode line breaking remain deferred.
