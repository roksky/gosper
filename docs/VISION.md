# Vision

## One paragraph

`gosper` is a **pure-Go, JasperReports-compatible reporting engine**. It reuses
existing JRXML templates (initially an explicitly documented subset of
JasperReports 6.x-style templates, typically authored in Jaspersoft Studio) and
generates reports entirely in-process, **without a JVM and without a mandatory
long-lived report server**. It is library-first: a small stable Go API that can
be embedded anywhere, with an optional CLI and (later) an optional service
adapter built on the same core.

## Why this exists

Existing JasperReports workflows require a JVM and often a report server.
Teams that have standardized on Go for their services must otherwise run a
sidecar JVM, re-implement reports, or give up on their existing JRXML assets.
`gosper` aims to let those teams keep their templates and render them natively.

## Personas

| Persona | Need |
| --- | --- |
| **Go application/platform engineer** | Embed report generation as a library; deterministic output; context-aware; no JVM. |
| **Report author** | Existing JRXML from Jaspersoft Studio renders predictably for the supported subset. |
| **Platform/ops engineer** | Predictable CPU/memory, no hidden network access, cancellable runs, safe handling of untrusted templates. |
| **Data engineer** | Pluggable data sources (slices, JSON, CSV, SQL) with parameter binding. |
| **Maintainer/agent** | Small, independently reviewable PRs with unambiguous ownership and acceptance criteria. |

## Supported usage modes

1. **Embedded library (primary).** Import the module; call compile/fill/export
   in-process. (Not yet implemented.)
2. **CLI (secondary).** `gosper` command for local compile/fill/export and
   metadata. (Today: metadata only.)
3. **Service adapter (deferred, optional).** An HTTP/gRPC front end over the
   same library, provided as a **separate module/package**, never required by
   the library. (Not implemented; see `EXT-006`/`EXT-007`.)

## Design assumptions (non-negotiable)

1. **Native Go at production runtime.** Java may exist **only** in an isolated
   comparison/test harness if adopted later; never in the production path.
2. **Small stable public API.** `internal/` owns implementation detail;
   public surface is introduced only with a real consumer.
3. **Compilation and execution are distinct.** Compiled definitions are
   immutable and reusable; execution state is per-run.
4. **Per-execution isolation.** Fill state is not shared across runs and is
   safe for concurrent use.
5. **Replaceable stages.** Source format, expression engine, data sources,
   layout, fill, print IR and exporters are independently replaceable.
6. **Operational discipline.** Context cancellation, predictable resource
   limits, input validation and deterministic output are requirements.
7. **Decimal correctness.** A decimal strategy is mandatory for financial
   expressions; binary floating point is not silently acceptable.
8. **No overdesign.** Tentative abstractions are documented as *proposals*, not
   promised APIs or shipped interfaces.
9. **Honest compatibility.** All docs distinguish *Jasper-compatible behavior*,
   *supported subset*, *extensions*, *known deviations* and *unsupported*.
10. **Clean-room and licensing-aware.** No JasperReports implementation code is
    copied; licensing and interoperability reviews are tracked work items.

## Non-goals

- **No JVM runtime** in production.
- **No mandatory hosted/long-lived service.**
- **No unbounded/arbitrary expression or class execution.** No arbitrary Java
  or Go code execution from templates.
- **No binary `.jasper` compatibility** in early milestones (see
  [`COMPATIBILITY.md`](COMPATIBILITY.md)).
- **No committed dates or effort estimates.** Roadmap milestones use `M0…Mn`
  identifiers, not calendar dates.
- **No promise of full report-designer parity** in v1.
- **No charts/crosstabs/barcodes** in the core v1 path (tracked as `EXT-*`).

## Related

- [`ROADMAP.md`](ROADMAP.md) — how the vision is sequenced.
- [`FEATURE_MATRIX.md`](FEATURE_MATRIX.md) — the concrete feature register.
- [`ARCHITECTURE.md`](ARCHITECTURE.md) — the structural interpretation.
