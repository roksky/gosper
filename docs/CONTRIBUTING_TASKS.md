# Contributing tasks (agents and humans)

This document defines how planning work is tracked so multiple contributors —
including autonomous coding agents — can work incrementally without collisions.
It complements the code-level [`CONTRIBUTING.md`](../CONTRIBUTING.md).

## 1. Adding a feature (stable ID)

1. Check [`FEATURE_MATRIX.md`](FEATURE_MATRIX.md) for an existing entry.
2. If new, allocate the next free three-digit sequence **within its domain**
   (`FND`, `JRX`, `EXP`, `DAT`, `FIL`, `LAY`, `PRT`, `OUT`, `TST`, `OPS`,
   `EXT`). **Never recycle or renumber IDs.**
3. Add the row with **exactly one** priority, milestone, status
   (default `Not Started`), and explicit dependencies (`—` if none, `TBD` if
   unknown and flagged).
4. Add the technical expectations and acceptance tests to the domain doc under
   [`features/`](features/).
5. Do **not** mark it `Ready` until dependencies are satisfied and acceptance
   criteria are written.

## 2. Splitting a feature into tasks

- A task is **one focused PR**, usually one package or a few related files.
- ID = `<FEATURE-ID>-<NNN>`, e.g. `JRX-001-002`; keep the ID in the filename
  and the `# <TASK-ID>:` heading.
- Use the template in [`tasks/README.md`](tasks/README.md).
- Contract/interface tasks come **before** implementations that depend on them.
- If a task is too large (`L`), split at a testable intermediate contract.

## 3. Reserving a task (collision avoidance)

Markdown cannot provide transactional locks. Coordination rules:

1. Set the task's `Status` to `In Progress` in its own file and in
   [`tasks/INDEX.md`](tasks/INDEX.md) **in the same PR/branch**.
2. Claim ownership via the collaboration channel you actually use (e.g. a PR
   label, an issue assignment). The authoritative reservation is that channel;
   the Markdown status is only a reflection.
3. Declare the **file paths you will touch**. If two tasks list overlapping
   files, they are **not** parallel-safe; serialize them or re-scope.
4. Keep edits to shared files (`go.mod`, CI, root `README.md`) minimal and
   isolated; prefer a dedicated "chore" task for shared-file changes.

## 4. Checking dependencies and integration contracts

- Before starting, confirm every `Dependencies` task/feature is `Done`.
- If you must change a shared contract (e.g. an `internal/` interface), create
  or update a **contract task** first and note the change in the domain doc and
  architecture doc.
- Record breaking contract changes explicitly; do not silently re-shape an
  interface another task depends on.

## 5. Status transitions

`Not Started` → `Ready` → `In Progress` → `In Review` → `Done`.

- `Ready`: dependencies satisfied **and** acceptance criteria present.
- `Blocked`: use only with an explicit dependency/condition; never as a vague
  pause. Include what unblocks it.
- `Deferred` / `Cancelled`: include a reason.

## 6. Marking `Done`

A feature/task is `Done` only with **observable verification**:

- a linked PR/commit, and
- passing relevant tests (cite the test path/command), or equivalent evidence.

Never mark `Done` from prose assertion alone. See `FND-011` (Done, cites
`cmd/gosper/main_test.go`) as the reference pattern.

## 7. Partial completion and reprioritization

- Keep the item `In Progress` and add a follow-up task for the remainder.
- Reprioritize by editing the priority in the canonical matrix **once**.
- Never falsify progress to appear further along.

## 8. Adding subfeatures without renumbering

Append new IDs with higher sequences. To replace an item, add `supersedes` /
`superseded by` notes in both rows; keep the old ID.

## 9. Which trackers update together

| Change | Update in the same PR |
| --- | --- |
| Task status | task file + [`tasks/INDEX.md`](tasks/INDEX.md) |
| Feature status/priority/milestone | [`FEATURE_MATRIX.md`](FEATURE_MATRIX.md) |
| Compatibility claim | [`COMPATIBILITY.md`](COMPATIBILITY.md) + the justifying fixture/test |
| Architecture/contract change | relevant `architecture/*` doc (+ ADR if durable) |
| Code change | tests + any of the above it touches |

[`STATUS.md`](STATUS.md) is a **derived** snapshot; refresh it opportunistically,
never treat it as the source of truth.

## 10. ADRs

- Propose with status `Proposed`; do not treat as settled.
- On agreement, set `Accepted` and update affected docs.
- Supersede, never delete. See [`adr/README.md`](adr/README.md).

## 11. Critical path and parallel queues

- **Critical path** is listed in [`ROADMAP.md`](ROADMAP.md).
- **Parallel-safe** tasks have disjoint file ownership and no unmet shared
  dependencies. Contract tasks in M0 are parallel-safe with each other because
  they touch different packages.

## 12. Preventing documentation drift

**PR checklist:**

- [ ] Any status changed is updated in **both** the task file and
      `tasks/INDEX.md`.
- [ ] Any feature status/scope change is reflected in `FEATURE_MATRIX.md`.
- [ ] Links resolve; no new dangling IDs.
- [ ] Compatibility rows changed only with a justifying fixture/test.
- [ ] `STATUS.md` updated if the evidenced baseline changed.

**Periodic consistency audit:** scan for dangling task/feature references,
duplicate IDs, features with multiple priorities/milestones, and dependency
cycles; fix in a small "docs consistency" PR.

---

## Example: fictional task lifecycle (illustrative only — not real status)

> This example is fictional and does not describe any task's current state.

1. Feature `DAT-003` (JSON source) is `Not Started`, depends on `DAT-002`.
2. `DAT-002-001` (cursor contract) reaches `Done` with a passing interface test.
3. A contributor writes task `DAT-003-001`, sets `Ready`, reserves it via a PR
   label, and lists files: `internal/datasource/json.go`,
   `internal/datasource/json_test.go`, `testdata/json/basic.json`.
4. While implementing, they set `In Progress` in the task file and INDEX.
5. They open a PR; reviewer sets `In Review`.
6. CI passes, fixture test added; they set `Done`, cite the PR and test, and
   mark `DAT-003` `Done` in the matrix.

## Example: two parallel-safe tasks (fictional)

| Task | Files owned | Why parallel-safe |
| --- | --- | --- |
| `LAY-005-001` (font resolver) | `internal/text/font*.go`, `internal/text/font_test.go` | touches only `internal/text` |
| `DAT-002-001` (cursor contract) | `internal/datasource/datasource.go`, `..._test.go` | touches only `internal/datasource` |

They share no files and neither depends on the other, so they can proceed
simultaneously. If both needed to edit `go.mod`, they would **not** be
parallel-safe.
