# JRX-022-001: Compatibility matrix and fixture registry

- **Status:** Done
- **Feature:** [JRX-022](../../FEATURE_MATRIX.md#jrx--jrxml-syntaxversion-compatibility-and-compilation)
- **Milestone:** M0
- **Priority:** P1
- **Dependencies:** none
- **Blocked by:** none
- **Primary owner/package:** `docs/COMPATIBILITY.md` + `testdata/`
- **Likely touched files:** `docs/COMPATIBILITY.md`, `docs/tasks/INDEX.md` (link check), optional `testdata/registry.md`
- **Collision risk:** Low — documentation and registry convention.
- **Parallel-safe with:** `OPS-008-001`, `TST-001-002`
- **Estimated size:** XS

## Objective

Ensure every compatibility claim is traceable to a fixture and status, and
define how the registry is updated as features land.

## Current state / prerequisites

[`../../COMPATIBILITY.md`](../../COMPATIBILITY.md) exists with a full matrix and
status vocabulary (established in the planning pass).

## Scope (in)

- Verify the matrix covers all M0/M1 constructs and uses the defined statuses.
- Define the fixture-link convention (`testdata/<suite>/...`).
- Define the update rule: compatibility row + justifying fixture/test change in
  the **same PR**.

## Out of scope

- Creating fixtures (feature tasks).
- Changing compatibility policy (that is `ADR-0003`).

## Proposed contracts and integration points

Registry convention only; no code.

## Implementation checklist (future agent only)

- [x] Audit `COMPATIBILITY.md` rows against `FEATURE_MATRIX.md` construct IDs.
- [x] Add the fixture-link convention where rows can already reference paths.
- [x] State the same-PR rule clearly.

## Acceptance criteria (observable)

- Every matrix row names a feature ID; no row claims `Verified` without a test
  path.
- Update rule documented.

## Tests and fixtures to add (future agent only)

- None.

## Failure cases / edge cases

- Rows for `Unsupported` constructs need no fixture but must cite policy.

## Definition of done

Registry conventions documented; no unverifiable compatibility claims.

## References and downstream dependencies

- Feature: `JRX-022`; [`ADR-0003`](../../adr/0003-jrxml-6x-subset-first.md).
- Downstream: all parser/layout/exporter tasks that add fixtures.
