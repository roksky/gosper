# LAY-005-001: Font resolver

- **Status:** Not Started
- **Feature:** [LAY-005](../../FEATURE_MATRIX.md#lay--layout-typography-pagination)
- **Milestone:** M1
- **Priority:** P0
- **Dependencies:** none
- **Blocked by:** none
- **Primary owner/package:** `internal/text` (proposed)
- **Likely touched files:** `internal/text/font.go`, `internal/text/font_test.go`, `internal/text/doc.go`
- **Collision risk:** Low/Medium — new package; `LAY-006-001` adds next to it.
- **Parallel-safe with:** `internal/jrxml` tasks, `DAT-001-001`, `FND-006-001`
- **Estimated size:** M

## Objective

Resolve logical font names (e.g. `SansSerif`, `Serif`, `Monospaced`, and a
bundled default) to a metric provider, with deterministic fallback and no
system-font dependence in tests.

## Current state / prerequisites

No text package. OQ-3 (shaping backend) is deferred; M1 needs only basic
metrics.

## Scope (in)

- A small set of **bundled** core fonts (or a bundled metric table) with stable
  metrics.
- `Resolve(name, style, size) (Font, error)`; fallback rules documented.
- Tests: known names resolve; unknown name ⇒ fallback or diagnostic.

## Out of scope

- Unicode shaping (`LAY-007`), embedding (`OUT-004`), JRXML font declarations
  (`JRX-023`).

## Proposed contracts and integration points

```go
// Proposed — non-binding.
type Font interface{ Measure(s string, size float64) (width Unit, err error) }
type Resolver interface{ Resolve(name string, style Style, size float64) (Font, error) }
```

## Implementation checklist (future agent only)

- [ ] Bundle metric data (no network; license-check any asset via `OPS-008`).
- [ ] Implement resolver + fallback.
- [ ] Tests for resolution and stable widths.
- [ ] Run `make check`.

## Acceptance criteria (observable)

- Resolving a core font returns stable widths for a fixed string/size.
- Unknown font resolves via documented fallback without panic.

## Tests and fixtures to add (future agent only)

- `internal/text/font_test.go`.

## Failure cases / edge cases

- Font assets must not break reproducibility; prefer metrics over binaries at
  M1 if licensing is unclear.

## Definition of done

Implemented/tested; asset licensing noted; linked from architecture.

## References and downstream dependencies

- Feature: `LAY-005`; OQ-3.
- Unblocks: `LAY-006-001`, `OUT-002`, `JRX-023`.
