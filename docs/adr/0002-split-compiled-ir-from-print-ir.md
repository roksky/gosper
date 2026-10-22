# 0002: Split compiled report IR from paginated print IR

- **Status:** Proposed
- **Date:** 2026-10-22
- **Supersedes:** —
- **Affects:** `PRT-001`, `PRT-002`, `PRT-004`, `OUT-001`, `OUT-006`

## Context

A naive design uses a single report object from compilation through rendering.
That couples compilation to data and geometry, makes cached compilation
unsafe to reuse, and forces every exporter to understand report semantics.
Exporters are also not equivalent: PDF targets visual fidelity while CSV/XLSX
target data semantics.

## Decision

Keep two distinct intermediate representations:

1. **Compiled report IR** — immutable, data-free, geometry-free, cacheable.
2. **Print IR** — paginated, with absolute coordinates and resolved styles,
   produced per run and consumed by exporters.

Exporters depend only on the print IR.

## Alternatives considered

1. **Single mutable report object** — rejected: reuse/concurrency hazards and
   exporter coupling.
2. **Exporter-specific models generated from the compiled IR** — rejected:
   duplicates semantics and drifts.
3. **Serialized print IR as a public file format** — deferred (`PRT-006`).

## Consequences

- Layout is the only stage that needs text metrics and geometry.
- Exporter implementations stay thin and independently replaceable.
- Compiled IR caching becomes straightforward (`FND-004`).
- A stable internal unit system must be chosen (`PRT-001-001`, OQ-2).

## References

- [`../architecture/intermediate-representations.md`](../architecture/intermediate-representations.md)
- [`../architecture/exporters.md`](../architecture/exporters.md)
