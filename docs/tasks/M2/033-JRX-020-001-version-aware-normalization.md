# JRX-020-001: Version-aware normalization

- **Status:** Done
- **Feature:** [JRX-020](../../FEATURE_MATRIX.md#jrx--jrxml-syntaxversion-compatibility-and-compilation)
- **Milestone:** M2
- **Priority:** P1
- **Dependencies:** JRX-001-002, JRX-019-001
- **Blocked by:** none
- **Primary owner/package:** `internal/jrxml`
- **Likely touched files:** `internal/jrxml/normalize.go`, `internal/jrxml/normalize_test.go`
- **Collision risk:** Medium — post-parse pass over the model; separate file.
- **Parallel-safe with:** `DAT-003-001`, `LAY-008-001`
- **Estimated size:** M

## Objective

Normalize dialect/version differences in the parsed model so downstream stages
see a single stable shape, while warning on uncertain cases.

## Current state / prerequisites

Parse skeleton + unsupported diagnostics exist; M1 model stable.

## Scope (in)

- Apply documented normalizations (e.g. attribute defaults/family differences)
  to produce a canonical model.
- Emit warnings for constructs whose semantics differ by version.
- Tests: equivalent inputs across variants normalize to equal models.

## Out of scope

- JRXML 7.x semantics (`EXT-008`).

## Proposed contracts and integration points

`Normalize(*Template) diag.List`; called by `internal/compiler` before
compilation.

## Implementation checklist (future agent only)

- [x] Implement normalizations discovered from fixtures.
- [x] Tests asserting canonical equality.
- [x] Run `make check`.

## Acceptance criteria (observable)

- Two equivalent 6.x variants produce structurally equal normalized models.
- Uncertain constructs emit warnings, not errors.

## Tests and fixtures to add (future agent only)

- `testdata/jrxml/versions/` variants + `normalize_test.go`.

## Failure cases / edge cases

- Missing version attribute; conflicting defaults.

## Definition of done

Implemented/tested; compatibility notes updated.

## References and downstream dependencies

- Feature: `JRX-020`; [`ADR-0003`](../../adr/0003-jrxml-6x-subset-first.md).
- Unblocks: `FND-004` (cache keys) and M3 semantics.

## Verification and compatibility

`GOCACHE=/tmp/gosper-go-cache make check` passes. `internal/jrxml/normalize_test.go` compares `testdata/jrxml/versions/JRX-020-001-{legacy,modern}.jrxml` and checks warning cases; `report/normalize_test.go` verifies the active compile path. Normalization targets the supported 6.x subset: default Java language, canonical 6.x dialect, and legacy `isStretchWithOverflow` versus `textAdjust="StretchHeight"`. Unknown versions and uncertain text adjustments produce warnings. The exact source version should be read before normalization if needed for provenance.
