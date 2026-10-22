# 0004: Print IR internal unit system is the twip

- **Status:** Proposed
- **Date:** 2026-10-22
- **Supersedes:** —
- **Affects:** `PRT-001`, `LAY-001`, `LAY-006`, `OUT-001`, OQ-2

## Context

Layout produces positions and sizes, and exporters consume them. A floating
point unit (inches or points) would accumulate rounding error across repeated
band placement and pagination, and would make visual golden comparisons
fragile. Prior to this decision the print IR unit was undecided (OQ-2).

## Decision

The print IR uses a single integral unit: the **twip** (1/1440 inch), stored as
`int64` (`printdoc.Unit`). Conversions from points/inches/millimetres round half
away from zero at the boundary and are exact for common values (1 in = 1440
twips, 1 pt = 20 twips).

## Alternatives considered

1. **Points (float64)** — rejected: floating-point drift, non-deterministic
   comparisons.
2. **Sub-point integers (e.g. 1/1000 pt)** — rejected: does not divide inches or
   points exactly, causing persistent rounding on common page sizes.
3. **Twips (chosen)** — divides both inches and points exactly and is integral.

## Consequences

- Layout arithmetic is integer and deterministic (`LAY-015`, `OUT-010`).
- Exporters convert twips to their native unit (e.g. PDF points) once, at the
  edge.
- Very large pages rely on `int64`; overflow is out of range for realistic page
  sizes but must be considered if arbitrary dimensions are ever accepted.
- The decimal/fractional precision of a unit is fixed at 1/1440 inch.

## References

- [`../architecture/intermediate-representations.md`](../architecture/intermediate-representations.md)
- `internal/printdoc/model.go`
- Task `PRT-001-001`; open question OQ-2 in
  [`../FEATURE_MATRIX.md`](../FEATURE_MATRIX.md#open-questions).
