# OUT-010-001: Deterministic output tests

- **Status:** Done
- **Feature:** [OUT-010](../../FEATURE_MATRIX.md#out--format-exporters)
- **Milestone:** M1
- **Priority:** P1
- **Dependencies:** OUT-001-001
- **Blocked by:** none
- **Primary owner/package:** `internal/exporter` (tests) + `internal/testutil`
- **Likely touched files:** `internal/exporter/determinism_test.go`, possibly minor fixes in `pdf.go`
- **Collision risk:** Medium — may require small changes to `pdf.go` owned by `OUT-001-001`; land after it.
- **Parallel-safe with:** `internal/jrxml` tasks
- **Estimated size:** S

## Objective

Guarantee and test that exporting identical input twice yields byte-identical
output (fixed/injected metadata timestamps, stable object ordering).

## Current state / prerequisites

Baseline PDF exporter exists.

## Scope (in)

- Inject a fixed creation date in tests; ensure no wall-clock or map-order
  dependence.
- Test: export twice ⇒ identical bytes.
- Document which bytes are allowed to vary (none, ideally).

## Out of scope

- Cross-version stability guarantees.

## Proposed contracts and integration points

Use `internal/testutil` fixed clock; may add an exporter option for metadata
timestamps.

## Implementation checklist (future agent only)

- [x] Add a clock/timestamp injection point if missing (not needed: PDF emits no timestamps).
- [x] Add the byte-equality test.
- [x] Run `make check`.

## Acceptance criteria (observable)

- Exporting the same document twice yields identical bytes.
- No wall-clock value appears in output (or is fixed).

## Tests and fixtures to add (future agent only)

- `internal/exporter/determinism_test.go`.

## Failure cases / edge cases

- Compression randomness (must be disabled or seeded).

## Definition of done

Determinism test passing; behavior documented.

## References and downstream dependencies

- Feature: `OUT-010`; plan: [`../../TEST_STRATEGY.md`](../../TEST_STRATEGY.md).
