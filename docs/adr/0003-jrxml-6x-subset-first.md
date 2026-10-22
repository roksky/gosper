# 0003: JRXML 6.x subset-first, conservative compatibility claims

- **Status:** Proposed
- **Date:** 2026-10-22
- **Supersedes:** —
- **Affects:** `JRX-001`, `JRX-020`, `EXT-008`, `JRX-021`, `COMPATIBILITY.md`

## Context

JasperReports has a large JRXML surface, and 7.x differs from 6.x. Claiming
broad compatibility early would be misleading and hard to verify. Templates in
the field are commonly authored in Jaspersoft Studio against 6.x-style JRXML.

## Decision

Target an **explicitly documented subset of JasperReports 6.x-style JRXML**
first. JasperReports 7.x is a distinct later milestone (`EXT-008`). Do not claim
full compatibility. Compiled `.jasper` binaries are out of scope (`JRX-021`).
Every compatibility claim must map to a fixture and a status in
[`COMPATIBILITY.md`](../COMPATIBILITY.md).

## Alternatives considered

1. **Aim for full 6.x compatibility up front** — rejected: unverifiable and
   blocks incremental delivery.
2. **Target 7.x first** — rejected: needs 7.x reference behavior and is newer;
   start from the more common baseline.
3. **Define compatibility by feature count** — rejected: percentages are
   meaningless without fixtures and evidence.

## Consequences

- [`COMPATIBILITY.md`](../COMPATIBILITY.md) is the single compatibility
  register, with per-construct statuses and evidence rules.
- Unsupported constructs must produce clear diagnostics (`JRX-019`).
- 7.x adaptation is explicitly deferred (`EXT-008`).

## References

- [`../COMPATIBILITY.md`](../COMPATIBILITY.md)
- [`../architecture/expressions-and-types.md`](../architecture/expressions-and-types.md)
