# TST-001-002: Fixture conventions and golden update workflow

- **Status:** Not Started
- **Feature:** [TST-001](../../FEATURE_MATRIX.md#tst--test-infrastructure-and-correctness-suites)
- **Milestone:** M0
- **Priority:** P0
- **Dependencies:** TST-001-001
- **Blocked by:** none
- **Primary owner/package:** `testdata/` + `docs/`
- **Likely touched files:** `testdata/README.md`, `docs/TEST_STRATEGY.md`, optional `testdata/<suite>/` layout docs
- **Collision risk:** Low — documentation plus directory layout.
- **Parallel-safe with:** `OPS-008-001`, `JRX-022-001`
- **Estimated size:** XS

## Objective

Fix fixture directory conventions and the golden-file update workflow so future
tests are consistent and reviewable.

## Current state / prerequisites

`testdata/README.md` exists (bootstrap). `TEST_STRATEGY.md` describes fixtures
in prose.

## Scope (in)

- Define `testdata/<suite>/` layout and naming (`<feature-id>-<NNN>-<slug>`).
- Define an explicit `-update` flag convention for regenerating goldens.
- Document review rules (golden diffs must be inspected deliberately).
- Update `testdata/README.md` accordingly.

## Out of scope

- Implementing the harness (`TST-002-001`).
- Adding real fixtures (feature tasks).

## Proposed contracts and integration points

Convention: `testdata/<suite>/{input.jrxml,data.json,golden/*.pdf}` and a
`-update` test flag. No code symbols introduced here.

## Implementation checklist (future agent only)

- [ ] Update `testdata/README.md` with the layout and naming rules.
- [ ] Cross-link from `docs/TEST_STRATEGY.md`.
- [ ] State that goldens are never auto-committed without review.

## Acceptance criteria (observable)

- Documentation unambiguously specifies where inputs and goldens live and how
  to regenerate them.

## Tests and fixtures to add (future agent only)

- None (documentation task).

## Failure cases / edge cases

- Case-insensitive filesystems: use consistent lowercase directory names.

## Definition of done

Conventions documented and cross-linked; no fictional fixtures added.

## References and downstream dependencies

- Feature: `TST-001`; plan: [`../../TEST_STRATEGY.md`](../../TEST_STRATEGY.md).
- Unblocks: `TST-002-001`, `JRX-022-001`.
