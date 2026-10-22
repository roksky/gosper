# gosper documentation hub

This is the entry point for all planning, architecture and tracking
documentation for **gosper**, a pure-Go, JasperReports-compatible reporting
engine (in bootstrap/scaffold state — no engine functionality is implemented).

> **Status: scaffold only.** See [`STATUS.md`](STATUS.md) for an evidenced
> snapshot of what exists today. Do not treat planned capabilities as
> implemented.

## Where to find things

| Document | Responsibility |
| --- | --- |
| [`VISION.md`](VISION.md) | Goals, personas, supported usage modes, non-goals. |
| [`STATUS.md`](STATUS.md) | Derived snapshot of evidenced baseline + blockers. |
| [`ROADMAP.md`](ROADMAP.md) | Ordered milestones, gates, dependency DAG, release scope. |
| [`ARCHITECTURE.md`](ARCHITECTURE.md) | System overview, invariants, intended structure. |
| [`FEATURE_MATRIX.md`](FEATURE_MATRIX.md) | **Canonical** feature register (IDs, priority, milestone, status). |
| [`COMPATIBILITY.md`](COMPATIBILITY.md) | JasperReports/JRXML compatibility policy + matrix. |
| [`TEST_STRATEGY.md`](TEST_STRATEGY.md) | Unit/contract/golden/differential/perf/security testing plan. |
| [`CONTRIBUTING_TASKS.md`](CONTRIBUTING_TASKS.md) | Rules for agents: IDs, reservations, status transitions. |
| [`development.md`](development.md) | Local setup, commands, CI, module conventions. |
| [`DEPENDENCY_POLICY.md`](DEPENDENCY_POLICY.md) | How third-party dependencies are justified, licensed and tracked. |
| [`LICENSING.md`](LICENSING.md) | Licensing decision (currently **pending**). |

### Architecture deep dives — [`architecture/`](architecture/)

- [`package-boundaries.md`](architecture/package-boundaries.md)
- [`report-lifecycle.md`](architecture/report-lifecycle.md)
- [`intermediate-representations.md`](architecture/intermediate-representations.md)
- [`expressions-and-types.md`](architecture/expressions-and-types.md)
- [`filling-and-groups.md`](architecture/filling-and-groups.md)
- [`layout-and-pagination.md`](architecture/layout-and-pagination.md)
- [`data-and-resources.md`](architecture/data-and-resources.md)
- [`exporters.md`](architecture/exporters.md)
- [`concurrency-and-limits.md`](architecture/concurrency-and-limits.md)

### Feature domains — [`features/`](features/)

Start at [`features/README.md`](features/README.md); it indexes the domain
documents (`FND`, `JRX`, `EXP`, `DAT`, `FIL`, `LAY`, `PRT`, `OUT`, `TST`,
`OPS`, `EXT`). These hold technical expectations and feature-specific
acceptance tests. The **canonical status register** remains
[`FEATURE_MATRIX.md`](FEATURE_MATRIX.md).

### Tasks — [`tasks/`](tasks/)

- [`tasks/README.md`](tasks/README.md) — task template and assignment rules.
- [`tasks/INDEX.md`](tasks/INDEX.md) — canonical task queue, status and
  dependencies.

### Decisions — [`adr/`](adr/)

- [`adr/README.md`](adr/README.md) — ADR index and status.
- [`adr/0001-native-go-library-first-runtime.md`](adr/0001-native-go-library-first-runtime.md)
- [`adr/0002-split-compiled-ir-from-print-ir.md`](adr/0002-split-compiled-ir-from-print-ir.md)
- [`adr/0003-jrxml-6x-subset-first.md`](adr/0003-jrxml-6x-subset-first.md)
- [`adr/0004-print-ir-unit-system.md`](adr/0004-print-ir-unit-system.md)

## Source-of-truth rules

To avoid divergent registers, each fact has exactly one canonical home:

| Fact | Canonical document |
| --- | --- |
| Feature identity, scope, priority, milestone, status, dependencies | [`FEATURE_MATRIX.md`](FEATURE_MATRIX.md) |
| Task status, ordering, ownership/reservation | [`tasks/INDEX.md`](tasks/INDEX.md) |
| Milestone sequence, exit gates, deferred scope | [`ROADMAP.md`](ROADMAP.md) |
| Feature-specific expectations/acceptance | [`features/*.md`](features/) |
| Evidenced baseline + curated blockers | [`STATUS.md`](STATUS.md) |
| Proposed structure, invariants, rationale | [`ARCHITECTURE.md`](ARCHITECTURE.md), [`architecture/*`](architecture/) |
| Durable decisions and alternatives | [`adr/*`](adr/) |

When status must be repeated for readability, label it a *derived snapshot* and
link to the canonical source. See [`CONTRIBUTING_TASKS.md`](CONTRIBUTING_TASKS.md)
for the update procedure.

## Project README and contribution guide

- [Project README](../README.md)
- [How to contribute code](../CONTRIBUTING.md)
- [How to contribute/track planning work](CONTRIBUTING_TASKS.md)
