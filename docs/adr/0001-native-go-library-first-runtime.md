# 0001: Native-Go, library-first runtime with optional adapters

- **Status:** Proposed
- **Date:** 2026-10-22
- **Supersedes:** —
- **Affects:** `FND-001`, `FND-005`, `FND-008`, `EXT-006`, `EXT-007`

## Context

JasperReports runs on the JVM and commonly relies on a report server. The
project wants to render JRXML templates from Go applications without requiring a
JVM or a long-lived service, while keeping the door open for a hosted option.
The bootstrap already provides a library-shaped module with a metadata-only CLI
and no server.

## Decision

The production runtime is **pure Go**. The primary deliverable is a **library**
with a small stable public API. The CLI is a thin consumer. Any HTTP/gRPC
service is an **optional, separate module** built on the same library and is
never a dependency of the library. Java may appear **only** in an isolated
comparison/test harness if adopted later.

## Alternatives considered

1. **JVM-based engine with Go bindings** — rejected: violates no-JVM goal and
   complicates deployment.
2. **Mandatory report server** — rejected: heavy operational footprint.
3. **Library + first-class built-in server** — deferred: couples the library to
   transport concerns; a separate module is cleaner.
4. **CLI-only tool** — rejected: embedding is the primary use case.

## Consequences

- Public API (`report`, proposed) must avoid transport assumptions.
- Server work is tracked separately (`EXT-006`, `EXT-007`, `FND-008`).
- Reference/differential testing against JasperReports is isolated from the
  production path (`TST-006`, `OPS-009`).

## References

- [`../VISION.md`](../VISION.md)
- [`../ARCHITECTURE.md`](../ARCHITECTURE.md)
- [`../architecture/package-boundaries.md`](../architecture/package-boundaries.md)
