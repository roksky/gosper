# TST-002-001: Golden fixture harness

- **Status:** Not Started
- **Feature:** [TST-002](../../FEATURE_MATRIX.md#tst--test-infrastructure-and-correctness-suites)
- **Milestone:** M1
- **Priority:** P0
- **Dependencies:** TST-001-001, TST-001-002
- **Blocked by:** none
- **Primary owner/package:** `internal/testutil`
- **Likely touched files:** `internal/testutil/golden.go`, `internal/testutil/golden_test.go`
- **Collision risk:** Low — new helper file; used by many later tests.
- **Parallel-safe with:** `internal/jrxml` tasks, `internal/text` tasks
- **Estimated size:** S

## Objective

Provide a reusable golden-file comparison helper with an explicit `-update`
flag, deterministic output, and readable diffs.

## Current state / prerequisites

Environment helpers and fixture conventions exist.

## Scope (in)

- `Golden(t, got []byte, path)` helper; `-update` regenerates goldens.
- Line/byte diff on mismatch.
- Self-test using a temp directory.

## Out of scope

- Semantic/pagination assertions (`TST-003`, `TST-004`).

## Proposed contracts and integration points

Used by exporter/layout/parser tests.

## Implementation checklist (future agent only)

- [ ] Implement helper + `-update` flag.
- [ ] Self-test.
- [ ] Run `make check`.

## Acceptance criteria (observable)

- First run with `-update` writes the golden; subsequent runs compare.
- A deliberate mismatch fails with a diff.

## Tests and fixtures to add (future agent only)

- `internal/testutil/golden_test.go`.

## Failure cases / edge cases

- Missing golden without `-update` ⇒ clear failure message.

## Definition of done

Helper merged/tested; used by at least one exporter test.

## References and downstream dependencies

- Feature: `TST-002`; plan: [`../../TEST_STRATEGY.md`](../../TEST_STRATEGY.md).
- Unblocks: `TST-003`, `TST-004`, `TST-005`.
