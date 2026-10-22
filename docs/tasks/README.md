# Tasks

Task specifications are **future PR-sized units of work**. Writing one here is
**not** implementing it. The canonical task state (status, ordering,
dependencies, ownership) is [`INDEX.md`](INDEX.md); see
[`../CONTRIBUTING_TASKS.md`](../CONTRIBUTING_TASKS.md) for governance.

## Layout

```
tasks/
├── README.md
├── INDEX.md
├── M0/   # baseline and contracts
├── M1/   # first vertical slice
└── M2/   # a small number of clearly-scoped next tasks
```

Later milestones (`M3`+) are tracked as **features** in
[`../FEATURE_MATRIX.md`](../FEATURE_MATRIX.md) until their contracts are clear
enough to author tasks.

## Naming

M0 files use `<MILESTONE>/<FEATURE-ID>-<NNN>-<slug>.md`, e.g.
`M0/JRX-001-001-parser-io-contract-and-version-detection.md`. M1 and M2 files
use `<MILESTONE>/<ORDER>-<FEATURE-ID>-<NNN>-<slug>.md`. The three-digit order
prefix is unique across M1 and M2, so sorting filenames within each milestone
puts task prerequisites first. It is not part of the task ID; the task ID
still appears in the filename and the `# ` heading. Check `INDEX.md` for task
status and dependencies before starting work.

## Required template

Every task file uses exactly these sections:

```markdown
# <TASK-ID>: <short imperative title>

- **Status:** Not Started
- **Feature:** <FEATURE-ID> (link)
- **Milestone:** M0 or M1
- **Priority:** P0/P1/P2/P3
- **Dependencies:** <task IDs, or none>
- **Blocked by:** <condition, or none>
- **Primary owner/package:** <existing or proposed repo path>
- **Likely touched files:** <small scoped list>
- **Collision risk:** Low/Medium/High; explain overlap
- **Parallel-safe with:** <task IDs, if known>
- **Estimated size:** XS/S/M/L (relative, non-time-based)

## Objective
## Current state / prerequisites
## Scope (in)
## Out of scope
## Proposed contracts and integration points
## Implementation checklist (future agent only)
## Acceptance criteria (observable)
## Tests and fixtures to add (future agent only)
## Failure cases / edge cases
## Definition of done
## References and downstream dependencies
```

## Rules

- One focused PR per task; prefer `S`/`M`. Split `L` at a testable contract.
- No cycles. Contract tasks precede implementations that depend on them.
- Explicit file ownership to minimize shared-file edits (`go.mod`, CI,
  `README.md`).
- Acceptance criteria must be observable (specific tests/values), not "works".
- Ambiguous behavior ⇒ create a research/ADR task instead of guessing.
- Never assign contributor names or claim a task as taken.
