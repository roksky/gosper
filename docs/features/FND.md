# FND — Foundation, public API, CLI, infrastructure

**Architecture:** [`../architecture/package-boundaries.md`](../architecture/package-boundaries.md),
[`../architecture/report-lifecycle.md`](../architecture/report-lifecycle.md).

## Purpose

Define the public API shape, template loading, diagnostics, cancellation,
caching, CLI surface, and version/logging concerns. This is the seam every
other domain attaches to.

## Technical expectations

- **Small, stable public API** at the module root (package `report`,
  *proposed*). No premature exported types.
- **Compile/Fill/Export** are distinct calls; a convenience wrapper may combine
  them. All accept `context.Context` first.
- **Template loading** resolves inputs through an injected resolver
  (`internal/resources`, proposed) with deny-by-default URIs.
- **Diagnostics** carry a `diag.Code`, message, severity and optional source
  span (`FND-003`).
- **Cancellation** is honored at stage boundaries and inside iteration
  (`FND-009`, `OPS-001`).
- **Cache keys** include content hash + compiler schema version (`FND-004`).

### Proposed contract — non-binding

```go
// Proposed — non-binding.
func Compile(ctx context.Context, r io.Reader, opts ...Option) (*CompiledReport, error)
func Fill(ctx context.Context, c *CompiledReport, params Params, ds DataSource, opts ...Option) (*FilledReport, error)
func Export(ctx context.Context, f *FilledReport, w io.Writer, format Format) error
```

## Features

| ID | Expectation | Acceptance test (planned) |
| --- | --- | --- |
| FND-001 | Public API compiles/fills/exports with context. | API round-trip test with a trivial template. |
| FND-002 | Templates load via resolver; URIs deny-by-default. | Resolver rejects `file://`/`http://` by default. |
| FND-003 | Errors expose stable codes + spans. | Code/span assertions on a malformed template. |
| FND-004 | Cache keyed on content+schema version. | Changed flag ⇒ different compiled IR. |
| FND-005 | CLI reports status/exit codes. | CLI integration test with exit-code checks. |
| FND-006 | Build/version metadata exposed. | Version non-empty; override via ldflags. |
| FND-007 | Optional logging/telemetry hooks. | Hook invoked once per stage. |
| FND-008 | Service is a separate module. | Adapter depends on library, not vice versa. |
| FND-009 | Context + budgets on all long ops. | Cancel mid-run returns `context.Canceled`. |
| FND-010 | API stability policy documented. | Policy doc + compatibility notes. |
| FND-011 | CLI metadata (help/version/unknown). | `cmd/gosper/main_test.go` (passing). |

## Tasks

`FND-001-001`, `FND-003-001`, `FND-009-001` (M0); `FND-001-002`,
`FND-002-001`, `FND-006-001` (M1). See [`../tasks/INDEX.md`](../tasks/INDEX.md).

## Risks

- Freezing an API before its first consumer. Mitigate by marking proposals and
  deferring detail.
