# Architecture Decision Records

ADRs capture durable decisions and the alternatives considered. A `Proposed`
ADR is **not settled**; do not treat proposed decisions as agreed.

## Status values

`Proposed` · `Accepted` · `Rejected` · `Superseded`

## Index

| ID | Title | Status | Date | Affects |
| --- | --- | --- | --- | --- |
| [0001](0001-native-go-library-first-runtime.md) | Native-Go, library-first runtime with optional adapters | Proposed | 2026-10-22 | `FND-001`, `FND-008`, `EXT-006`, `EXT-007` |
| [0002](0002-split-compiled-ir-from-print-ir.md) | Split compiled report IR from paginated print IR | Proposed | 2026-10-22 | `PRT-001`, `PRT-002`, `OUT-*` |
| [0003](0003-jrxml-6x-subset-first.md) | JRXML 6.x subset-first, conservative compatibility claims | Proposed | 2026-10-22 | `JRX-001`, `JRX-020`, `EXT-008`, `COMPATIBILITY.md` |

## Process

1. Propose an ADR by copying the structure below into the next sequential file.
2. Open it as `Proposed`; discuss in the PR that references the relevant
   feature/task IDs.
3. On agreement, change status to `Accepted` and update
   [`FEATURE_MATRIX.md`](../FEATURE_MATRIX.md) / architecture docs if the
   decision changes them.
4. To replace a decision, create a new ADR and mark the old one `Superseded by
   NNNN`. **Never delete or renumber an ADR.**
5. `Rejected` ADRs are kept with their rationale.

## Template

```markdown
# NNNN: Title

- **Status:** Proposed
- **Date:** YYYY-MM-DD
- **Supersedes:** —
- **Affects:** FEATURE-IDs

## Context
## Decision
## Alternatives considered
## Consequences
## References
```

See [`CONTRIBUTING_TASKS.md`](../CONTRIBUTING_TASKS.md) for how ADR proposals
interact with task tracking.
