# LAY-006-001: Text measurement baseline

- **Status:** Not Started
- **Feature:** [LAY-006](../../FEATURE_MATRIX.md#lay--layout-typography-pagination)
- **Milestone:** M1
- **Priority:** P0
- **Dependencies:** LAY-005-001
- **Blocked by:** none
- **Primary owner/package:** `internal/text`
- **Likely touched files:** `internal/text/measure.go`, `internal/text/measure_test.go`
- **Collision risk:** Medium — same package as font resolver; separate file.
- **Parallel-safe with:** `internal/jrxml` tasks, `DAT-001-001`
- **Estimated size:** S

## Objective

Provide deterministic text measurement (width for a string at a size) plus line
height, sufficient for fixed-height M1 layout without shaping.

## Current state / prerequisites

Font resolver exists.

## Scope (in)

- Width measurement per font/size; ascent/descent/line height.
- Determinism guarantee (same inputs ⇒ same result).
- Tests with fixed strings.

## Out of scope

- Wrapping/alignment (`LAY-008-001`), kerning/shaping (`LAY-007`).

## Proposed contracts and integration points

`Measure(s, font, size) (width Unit)`; used by layout and exporters.

## Implementation checklist (future agent only)

- [ ] Implement measurement over the font metric provider.
- [ ] Tests: widths for known strings; determinism across calls.
- [ ] Run `make check`.

## Acceptance criteria (observable)

- Repeated measurement yields identical widths.
- Widths scale monotonically with size.

## Tests and fixtures to add (future agent only)

- `internal/text/measure_test.go`.

## Failure cases / edge cases

- Empty string; very long string (bound any caching).

## Definition of done

Implemented/tested; determinism documented; linked from architecture.

## References and downstream dependencies

- Feature: `LAY-006`.
- Unblocks: `LAY-008-001`, `OUT-002`, `LAY-015`.
