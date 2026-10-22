# OPS — Security, performance, diagnostics, integration

**Architecture:** [`../architecture/concurrency-and-limits.md`](../architecture/concurrency-and-limits.md).

## Purpose

Enforce resource limits, secure resource/SQL access, provide diagnostics and
telemetry, and manage release/licensing hygiene.

## Technical expectations

- **Budgets**: records, pages, evaluation steps, recursion depth, bytes
  (`OPS-001`, `OPS-007`).
- **Deny-by-default** resource URIs; rooted filesystem resolver (`OPS-002`).
- **Parameterized SQL only** (`OPS-003`).
- Decompression/recursion limits (`OPS-004`).
- Opt-in debug/tracing; OpenTelemetry hooks later (`OPS-005`, `OPS-006`).
- License/SBOM review before adding dependencies (`OPS-008`); JasperReports
  interoperability review (`OPS-009`).

## Features

| ID | Expectation | Acceptance test (planned) |
| --- | --- | --- |
| OPS-001 | Limit enforcement. | Each budget trips a diagnostic. |
| OPS-002 | URI allowlists. | `file://`/`http://` rejected by default. |
| OPS-003 | Safe SQL binding. | Injection fixture fails safely. |
| OPS-004 | Recursion/decompression limits. | Bomb/deep-recursion diagnostic. |
| OPS-005 | Debug/tracing controls. | Debug flag toggles output. |
| OPS-006 | OTel hooks. | Span emitted. |
| OPS-007 | Memory/CPU budgets. | Profile-bound large run. |
| OPS-008 | License/SBOM review. | Policy + tooling decision. |
| OPS-009 | JasperReports interop review. | Review record. |
| OPS-010 | Changelog/release guidance. | Changelog process. |
| OPS-011 | Examples/corpus. | Example builds and runs. |

## Tasks

`OPS-008-001` (M0); `OPS-001-001` (M1).

## Risks

- Adding dependencies without license review; enforce via `OPS-008`.
