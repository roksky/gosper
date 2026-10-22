# PRT-002-001: Resolved style data

- **Status:** Done
- **Feature:** [PRT-002](../../FEATURE_MATRIX.md#prt--print-document-and-storage)
- **Milestone:** M1
- **Priority:** P1
- **Dependencies:** PRT-001-002
- **Blocked by:** none
- **Primary owner/package:** `internal/printdoc`
- **Likely touched files:** `internal/printdoc/style.go`, `internal/printdoc/style_test.go`
- **Collision risk:** Medium — same package as model task; separate file.
- **Parallel-safe with:** `internal/text`, `internal/jrxml` tasks
- **Estimated size:** S

## Objective

Represent **resolved** style values (font, size, color, alignment, bold/italic,
padding) on print elements, rather than references that exporters would have to
resolve.

## Current state / prerequisites

Print model exists; M1 has a minimal style set (font/style/size/color/align).

## Scope (in)

- A `Style` value struct with the M1 subset and sensible defaults.
- Elements carry a resolved `Style`.
- Tests for defaults and explicit values.

## Out of scope

- Borders/transparency/conditional styles (`LAY-012`, `JRX-018`).

## Proposed contracts and integration points

```go
// Proposed — non-binding.
type Style struct { Font string; Size Unit; Bold, Italic bool; Color uint32; Align Align }
```

## Implementation checklist (future agent only)

- [x] Implement `Style` with defaults.
- [x] Tests for default and set values.
- [x] Run `make check`.

## Acceptance criteria (observable)

- Default style has documented values; explicit values round-trip.

## Tests and fixtures to add (future agent only)

- `internal/printdoc/style_test.go`.

## Failure cases / edge cases

- Unknown color format; out-of-range sizes.

## Definition of done

Implemented/tested; linked from architecture.

## References and downstream dependencies

- Feature: `PRT-002`.
- Unblocks: `OUT-001-001`, `LAY-012`.
