# 0005: Public API lives in package `report` at the module root

- **Status:** Proposed
- **Date:** 2026-10-22
- **Supersedes:** —
- **Affects:** `FND-001`, `FND-005`, `FND-008`, OQ-6

## Context

Callers need one stable, importable package for the compile/fill/export
lifecycle. Two questions were open: (a) whether the public package should sit at
the module root, and (b) whether orchestration should live in the public package
or in a dedicated `internal/exec` package (OQ-6).

## Decision

1. The public API is a **single package named `report` at the module root**
   (import path `github.com/roksky/gosper/report`). No `pkg/` hierarchy is
   introduced.
2. For now, **orchestration lives in `report`**: the entry points coordinate the
   `internal/*` stages directly. A dedicated `internal/exec` package is
   **deferred** until orchestration is complex enough to justify splitting it
   out (caching, retries, service adapters). If introduced later, it must not be
   imported by any other `internal/*` package and must not contain the public
   types.

## Alternatives considered

1. **`pkg/report` or a versioned subpackage** — rejected: adds an indirection
   layer with no benefit for a single public package; Go modules already provide
   versioning.
2. **`internal/exec` from the start** — rejected: an extra boundary with a single
   consumer; premature given the empty lifecycle.
3. **Top-level package `gosper`** — rejected: less descriptive than `report` and
   collides conceptually with the module name.

## Consequences

- `report` may import `internal/*`; no `internal/*` package may import `report`
  (enforced by the build).
- The public `DataSource`/`Cursor`/`Record` names are aliases of the
  `internal/datasource` contract so callers need not import an internal package.
- CLI and any future server adapter are thin consumers of `report`.
- Revisiting (2) is cheap: move orchestration behind an interface if/when a
  second consumer appears.

## References

- [`../architecture/package-boundaries.md`](../architecture/package-boundaries.md)
- [`../architecture/report-lifecycle.md`](../architecture/report-lifecycle.md)
- Task `FND-001-001`; OQ-6 in
  [`../FEATURE_MATRIX.md`](../FEATURE_MATRIX.md#open-questions).
